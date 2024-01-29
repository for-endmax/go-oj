package judge

import (
	"context"
	"encoding/binary"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-uuid"
	"go.uber.org/zap"
	"io"
	"judge_srv/global"
	"judge_srv/message"
	"os"
	"path"
	"strconv"
	"time"
)

// Output 返回
type Output struct {
	msg      []byte
	exitCode int
}

// 解析日志
// header := [8]byte{STREAM_TYPE, 0, 0, 0, SIZE1, SIZE2, SIZE3, SIZE4}
func parseDockerLog(logs []byte) []byte {
	output := make([]byte, 0, len(logs))

	for i := 0; i < len(logs); {
		sizeBinary := logs[i+4 : i+8]
		i += 8

		size := int(binary.BigEndian.Uint32(sizeBinary))
		data := logs[i : i+size]
		output = append(output, data...)
		i += size
	}

	return output
}

// Result 一次判断的结果
type Result struct {
	ErrCode int32 // 0正常 1编译或运行错误 2超时 3测试用例不通过 4内存超限
	ErrMsg  string
	runTime int32 //ms
	runMem  int32 //KB
}

// Task 判题任务
type Task struct {
	uuid         string
	containerID  string
	option       Option
	dockerClient *client.Client
	absPathDir   string
	msgSend      message.MsgSend
}

// CreateTask 创建任务
func CreateTask(msgSend message.MsgSend) (task *Task, err error) {
	task = &Task{}
	if task.uuid, err = uuid.GenerateUUID(); err != nil {
		return nil, err
	}
	task.msgSend = msgSend
	task.option = langMap[task.msgSend.Lang]
	if task.dockerClient, err = client.NewClientWithOpts(); err != nil {
		zap.S().Info("创建客户端失败")
		return nil, err
	}

	task.absPathDir = path.Join("/home/endmax/code", task.uuid)
	if err := os.MkdirAll(task.absPathDir, 0755); err != nil {
		zap.S().Info("创建目录失败")
		return nil, err
	}

	if err := os.WriteFile(path.Join(task.absPathDir, task.option.fileName), []byte(task.msgSend.SubmitCode), 0755); err != nil {
		zap.S().Info("创建文件失败")
		return nil, err
	}
	//创建容器
	createContainerResp, err := task.dockerClient.ContainerCreate(context.Background(),
		&container.Config{
			Image:        task.option.ImageName,
			User:         "root",
			WorkingDir:   "/judge",
			Tty:          false,
			AttachStdout: true,
			AttachStderr: true,
			AttachStdin:  true,
			Env:          nil,
		},
		&container.HostConfig{
			NetworkMode: "none",
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: task.absPathDir,
					Target: "/judge",
				},
			},
			//Resources: container.Resources{
			//	NanoCPUs: 1000000000 * 0.5,  // 总共是10^9
			//	Memory:   1024 * 1024 * 100, // 100MB
			//},
		}, nil, nil, task.uuid) // uuid作为container的名称
	if err != nil {
		zap.S().Info("创建容器失败")
		return nil, err
	}
	task.containerID = createContainerResp.ID

	// 启动容器
	if err := task.dockerClient.ContainerStart(context.Background(), task.containerID, container.StartOptions{}); err != nil {
		zap.S().Info("启动容器失败")
		return nil, err
	}
	return task, nil
}

// Clean 清理
func (t *Task) Clean() {
	zap.S().Info("删除容器和目录")
	// 清除
	if err := t.dockerClient.ContainerKill(context.Background(), t.containerID, "9"); err != nil {
		zap.S().Infof("停止容器失败: %v", err)
	}

	if err := t.dockerClient.ContainerRemove(context.Background(), t.containerID, container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		zap.S().Infof("删除容器失败: %v", err)
	}
	err := os.RemoveAll(t.absPathDir)
	if err != nil {
		zap.S().Infof("删除目录失败: %v", err)
	}
}

// Run 运行判题机进行判题(一次)
func (t *Task) Run(testInput string, testOutput string, num int) (error, *Result) {
	if num == 0 {
		// 编译
		zap.S().Info("编译...")
		buildOutput, err := t.Exec(t.option.buildCmd, "")
		if err != nil {
			zap.S().Info("编译服务出错")
			return err, nil
		}
		if buildOutput.exitCode != 0 {
			zap.S().Info("编译错误")
			return nil, &Result{
				ErrCode: 1,
				ErrMsg:  string(buildOutput.msg),
			}
		}
	}

	// 执行
	zap.S().Info("执行...")
	zap.S().Infof("测试输入数据：\n%s测试输出数据:\n%s", testInput, testOutput)
	runOutput, err := t.Exec(t.option.RunCmd, testInput)
	//zap.S()().Info()(string(runOutput.msg))
	if err != nil {
		zap.S().Info("执行服务出错")
		return err, nil
	}

	if runOutput.exitCode != 0 {
		if runOutput.exitCode == 137 {
			zap.S().Info("内存超限")
			return nil, &Result{
				ErrCode: 4,
				ErrMsg:  string(runOutput.msg),
			}
		}
		if runOutput.exitCode == 2001 {
			zap.S().Info("超时")
			return nil, &Result{
				ErrCode: 2,
				ErrMsg:  string(runOutput.msg),
			}
		}
		zap.S().Info("运行时错误")
		return nil, &Result{
			ErrCode: 1,
			ErrMsg:  string(runOutput.msg),
		}
	}

	//获取时间 0.00
	runTime, runMem, realMsg, err := GetResource(runOutput.msg)
	if err != nil {
		return err, nil
	}
	zap.S().Infof("运行时间 %d ms,运行内存 %d KB", runTime, runMem)
	// 测试用例输出比较
	if testOutput != realMsg {
		zap.S().Info("测试用例不通过")
		return nil, &Result{
			ErrCode: 3,
			ErrMsg:  realMsg,
		}
	}
	zap.S().Info("测试用例通过")

	result := Result{
		runTime: runTime,
		runMem:  runMem,
	}

	return nil, &result
}

// GetResource 获取时间(ms),内存KB  -1000:0.00
func GetResource(msg []byte) (time int32, mem int32, realMsg string, err error) {
	n := len(msg)
	timeStr := ""
	memStr := ""
	tmp := []byte{}
	for i := n - 2; i >= 0; i-- {
		if msg[i] == ':' {
			timeStr = string(reverse(tmp))
			tmp = []byte{}
		} else if msg[i] == '-' {
			memStr = string(reverse(tmp))
			realMsg = string(msg[:i])
			break
		} else {
			tmp = append(tmp, msg[i])
		}
	}
	// 将字符串转换为浮点数
	f, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		zap.S().Info("转换失败:", err)
		return 0, 0, "", nil
	}
	time = int32(f * 1000)

	m, err := strconv.Atoi(memStr)
	if err != nil {
		zap.S().Info("转换失败:", err)
		return 0, 0, "", nil
	}
	mem = int32(m)
	return
}

// 反转切片
func reverse(a []byte) []byte {
	n := len(a)
	for i := 0; i < n/2; i++ {
		a[i], a[n-1-i] = a[n-1-i], a[i]
	}
	return a
}

// Exec 执行指令
func (t *Task) Exec(cmd string, testData string) (*Output, error) {
	zap.S().Infof("执行命令：%s", cmd)
	if cmd == "" {
		return &Output{}, nil
	}
	// 创建执行命令的对象
	execResp, err := t.dockerClient.ContainerExecCreate(context.Background(), t.containerID, types.ExecConfig{
		User:         "root",
		Tty:          false,
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  true,
		Env:          nil,
		WorkingDir:   "/judge",
		Cmd:          []string{"sh", "-c", cmd},
	})
	if err != nil {
		return nil, err
	}

	// 创建连接对象
	attachResp, err := t.dockerClient.ContainerExecAttach(context.Background(), execResp.ID, types.ExecStartCheck{
		Detach: false, // 按下什么键可以分离命令
		Tty:    false, // 不使用交互式终端
	})
	if err != nil {
		return nil, err
	}
	defer attachResp.Close()

	// 往输入流中写入测试数据
	if testData != "" {
		//写数据
		if _, err := io.WriteString(attachResp.Conn, testData); err != nil {
			zap.S().Info("写入数据失败")
			return nil, err
		}
	}

	// 启动协程读取结果
	done := make(chan struct{})
	//读取输出
	var msg []byte
	go func() {
		msg, err = io.ReadAll(attachResp.Reader)
		close(done)
	}()

	select {
	case <-time.After(global.ExecTimeOut):
		return &Output{
			msg:      []byte("程序运行超时"),
			exitCode: 2001,
		}, nil
	case <-done:
		if err != nil {
			return nil, err
		}
		msg = parseDockerLog(msg)

		//检查执行结果
		inspectResp, err := t.dockerClient.ContainerExecInspect(context.Background(), execResp.ID)
		if err != nil {
			return nil, err
		}
		zap.S().Infof("退出状态码:%d", inspectResp.ExitCode)
		return &Output{
			msg:      msg,
			exitCode: inspectResp.ExitCode,
		}, nil
	}
}

// docker ps -a | awk '$2 == "endmax/go:latest" {print $1}' | xargs docker rm -f
// rm -rf code/
