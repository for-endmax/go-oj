package main

import (
	"question_srv/global"
	"question_srv/initialize"
	"question_srv/model"
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
	question := model.Question{
		Seq:     int32(1),
		Name:    "题目：" + "初始测试题目",
		Content: "将输入的每行字符串拼接起来,输入的最后一行为exit",
	}
	global.DB.Create(&question)

	testInput := "aa\nbb\ncc\nexit\n"
	testOutput := "aabbcc"
	test := model.Test{
		QID:    int32(1),
		Input:  testInput,
		Output: testOutput,
	}
	global.DB.Create(&test)

	testInput = "aaa\nbbb\nccc\nexit\n"

	testOutput = "aaabbbccc"
	test = model.Test{
		QID:    int32(1),
		Input:  testInput,
		Output: testOutput,
	}
	global.DB.Create(&test)
}

func main() {

	initialize.InitLogger() //初始化日志
	initialize.InitConfig() //读取配置信息
	initialize.InitDB()     //初始化MySQL

	insertData() // 建表并插入数据

}
