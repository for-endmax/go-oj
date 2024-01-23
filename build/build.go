package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// writeConfig2Consul 向consul写入配置信息
func writeConfig2Consul(url string, data string) {
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
	// 向consul写入各个服务的具体配置
	writeConfig2Consul("http://127.0.0.1:8500/v1/kv/go-oj/user_srv", "./build/user_srv_content.yaml")
	writeConfig2Consul("http://127.0.0.1:8500/v1/kv/go-oj/user_web", "./build/user_web_content.yaml")
	writeConfig2Consul("http://127.0.0.1:8500/v1/kv/go-oj/question_srv", "./build/question_srv_content.yaml")
	writeConfig2Consul("http://127.0.0.1:8500/v1/kv/go-oj/question_web", "./build/question_web_content.yaml")
	writeConfig2Consul("http://127.0.0.1:8500/v1/kv/go-oj/record_srv", "./build/record_srv_content.yaml")
	writeConfig2Consul("http://127.0.0.1:8500/v1/kv/go-oj/submit_web", "./build/submit_web_content.yaml")
}
