package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"question_srv/global"
	"question_srv/initialize"
	"question_srv/model"
	"strconv"
)

// writeConfig2Consul 向consul写入配置信息
func writeConfig2Consul() {
	url := "http://127.0.0.1:8500/v1/kv/go-oj/question_srv"
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

// insertData 建表并插入数据
func insertData() {
	// 建表
	err := global.DB.AutoMigrate(&model.Question{})
	if err != nil {
		panic(err)
	}

	err = global.DB.AutoMigrate(&model.Test{})
	if err != nil {
		panic(err)
	}
	// 插入数据
	for i := 0; i < 10; i++ {
		question := model.Question{
			Seq:     int32(i),
			Name:    "题目：" + strconv.Itoa(i),
			Content: "这是一道题目",
		}
		global.DB.Create(&question)

		myArray := []interface{}{1, 2, 3, 4, 5}

		// 将数组转换为 JSON 字符串
		jsonString, err := json.Marshal(myArray)
		if err != nil {
			fmt.Println("转换为 JSON 字符串时发生错误:", err)
			return
		}

		test := model.Test{
			QID:       int32(i),
			TimeLimit: 500,
			MemLimit:  3000,
			Input:     string(jsonString),
			Output:    string(jsonString),
		}
		global.DB.Create(&test)
	}
}

func main() {
	writeConfig2Consul() //写入配置信息

	initialize.InitLogger() //初始化日志
	initialize.InitConfig() //读取配置信息
	initialize.InitDB()     //初始化MySQL

	//insertData() // 建表并插入数据

}
