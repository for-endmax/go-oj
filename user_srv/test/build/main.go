package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"user_srv/global"
	"user_srv/initialize"
	"user_srv/model"
	"user_srv/utils"
)

// writeConfig2Consul 向consul写入配置信息
func writeConfig2Consul() {
	url := "http://127.0.0.1:8500/v1/kv/go-oj/user_srv"
	data := "./test/build/content.yaml"

	// 读取文件内容
	contents, err := os.ReadFile(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// 发送PUT请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(contents))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	// 打印响应
	fmt.Println("Response:", string(respBody))
}

func main() {
	writeConfig2Consul() //写入配置信息

	initialize.InitLogger() //初始化日志
	initialize.InitConfig() //读取配置信息
	initialize.InitDB()     //初始化MySQL

	// 建表
	err := global.DB.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}

	// 插入数据
	for i := 0; i < 5; i++ {
		salt, encodedPassword := utils.EncodePassword("admin123")
		user := model.User{
			BaseModel: model.BaseModel{},
			Password:  salt + ":" + encodedPassword,
			Nickname:  "somebody" + strconv.Itoa(i),
			Gender:    "male",
			Role:      1,
		}
		if i == 0 {
			user.Role = 2
			user.Nickname = "endmax"
		}
		fmt.Println(utils.VerifyPassword("admin123", salt+":"+encodedPassword))
		global.DB.Create(&user)
	}

}
