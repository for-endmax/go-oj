package judge

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-uuid"
	"io"
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

// Option 各个语言的配置
type Option struct {
	fileName  string
	ImageName string
	buildCmd  string
	RunCmd    string
}

var langMap = map[string]Option{
	"go": {
		fileName:  "main.go",
		ImageName: "endmax/go:latest",
		buildCmd:  "go mod init endmax && go build main.go",
		RunCmd:    "time -f \"-%K:%e\" ./main",
	},
	"python": {
		fileName:  "main.py",
		ImageName: "endmax/python:latest",
		buildCmd:  "",
		RunCmd:    "time -f \"-%K:%e\" python3 main.py",
	},
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
		fmt.Println("创建客户端失败")
		return nil, err
	}

	task.absPathDir = path.Join("/home/endmax/code", task.uuid)
	if err := os.MkdirAll(task.absPathDir, 0755); err != nil {
		fmt.Println("创建目录失败")
		return nil, err
	}

	if err := os.WriteFile(path.Join(task.absPathDir, task.option.fileName), []byte(task.msgSend.SubmitCode), 0755); err != nil {
		fmt.Println("创建文件失败")
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
			//	NanoCPUs: 1000000000 * 0.9,  // 总共是10^9
			//	Memory:   1024 * 1024 * 100, // 100MB
			//},
		}, nil, nil, task.uuid) // uuid作为container的名称
	if err != nil {
		fmt.Println("创建容器失败")
		return nil, err
	}
	task.containerID = createContainerResp.ID

	// 启动容器
	if err := task.dockerClient.ContainerStart(context.Background(), task.containerID, container.StartOptions{}); err != nil {
		fmt.Println("启动容器失败")
		return nil, err
	}
	return task, nil
}

// Clean 清理
func (t *Task) Clean() {
	fmt.Println("\n删除容器和目录")
	// 清除
	if err := t.dockerClient.ContainerKill(context.Background(), t.containerID, "9"); err != nil {
		fmt.Printf("Failed to kill container: %v", err)
	}

	if err := t.dockerClient.ContainerRemove(context.Background(), t.containerID, container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		fmt.Printf("Failed to remove container: %v", err)
	}
	err := os.RemoveAll(t.absPathDir)
	if err != nil {
		fmt.Printf("Failed to remove volume folder: %v", err)
	}
}

// Run 运行判题机进行判题(一次)
func (t *Task) Run(testInput string, testOutput string) (error, *Result) {

	// 编译
	fmt.Println("编译...")
	buildOutput, err := t.Exec(t.option.buildCmd, "")
	if err != nil {
		fmt.Println("\n编译服务出错")
		return err, nil
	}
	if buildOutput.exitCode != 0 {
		fmt.Println("\n编译错误")
		return nil, &Result{
			ErrCode: 1,
			ErrMsg:  string(buildOutput.msg),
		}
	}
	// 执行
	fmt.Println("\n执行...")
	fmt.Printf("测试输入数据：\n%s测试输出数据:\n%s\n", testInput, testOutput)
	runOutput, err := t.Exec(t.option.RunCmd, testInput)
	//fmt.Println(string(runOutput.msg))
	if err != nil {
		fmt.Println("\n执行服务出错")
		return err, nil
	}

	if runOutput.exitCode != 0 {
		if runOutput.exitCode == 137 {
			fmt.Println("内存超限")
			return nil, &Result{
				ErrCode: 4,
				ErrMsg:  string(runOutput.msg),
			}
		}
		if runOutput.exitCode == 2001 {
			fmt.Println("超时")
			return nil, &Result{
				ErrCode: 2,
				ErrMsg:  string(runOutput.msg),
			}
		}
		fmt.Println("\n运行时错误")
		return nil, &Result{
			ErrCode: 1,
			ErrMsg:  string(runOutput.msg),
		}
	}

	//获取时间 0.00
	runTime, runMem, realMsg := GetResource(runOutput.msg)
	fmt.Printf("运行时间 %d ms,运行内存 %d KB", runTime, runMem)
	// 测试用例输出比较
	if testOutput != realMsg {
		fmt.Println("\n测试用例不通过")
		return nil, &Result{
			ErrCode: 3,
			ErrMsg:  realMsg,
		}
	}
	fmt.Println("\n测试用例通过")

	result := Result{
		runTime: runTime,
		runMem:  runMem,
	}
	//判断是否超时或超内存
	if runMem > t.msgSend.MemLimit {
		//超内存
		result.ErrCode = 4
		result.ErrMsg = "超内存"
	}
	if runTime > t.msgSend.TimeLimit {
		//超时
		result.ErrCode = 2
		result.ErrMsg = "超时"
	}
	return nil, &result
}

// GetResource 获取时间(ms),内存KB  -1000:0.00
func GetResource(msg []byte) (time int32, mem int32, realMsg string) {
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
		fmt.Println("转换失败:", err)
		panic(err)
	}
	time = int32(f * 1000)

	m, err := strconv.Atoi(memStr)
	if err != nil {
		fmt.Println("转换失败:", err)
		panic(err)
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
	fmt.Printf("执行命令：%s\n", cmd)
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
			fmt.Println("写入数据失败")
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
	case <-time.After(time.Second * 60):
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
		fmt.Printf("退出状态码:%d\n", inspectResp.ExitCode)
		return &Output{
			msg:      msg,
			exitCode: inspectResp.ExitCode,
		}, nil
	}
}
