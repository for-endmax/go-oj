package main

import (
	"encoding/json"
	"fmt"
	"question_srv/global"
	"question_srv/initialize"
	"question_srv/model"
	"strconv"
)

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
			QID:    int32(i),
			Input:  string(jsonString),
			Output: string(jsonString),
		}

		global.DB.Create(&test)
		test = model.Test{
			QID:    int32(i),
			Input:  string(jsonString) + "123",
			Output: string(jsonString) + "123",
		}
		global.DB.Create(&test)
	}
}

func main() {

	initialize.InitLogger() //初始化日志
	initialize.InitConfig() //读取配置信息
	initialize.InitDB()     //初始化MySQL

	insertData() // 建表并插入数据

}
