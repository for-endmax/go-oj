package main

import (
	"record_srv/global"
	"record_srv/initialize"
	"record_srv/model"
)

// insertData 建表并插入数据
func insertData() {
	// 建表
	err := global.DB.AutoMigrate(&model.RecordModel{})
	if err != nil {
		panic(err)
	}

	code := `package main
import "fmt"
func main(){
	fmt.Println("hello world")
}`
	// 插入数据
	for i := 0; i < 10; i++ {
		record := model.RecordModel{
			UID:        int32(i),
			QID:        int32(i),
			Lang:       "go",
			Status:     0,
			ErrCode:    0,
			ErrMsg:     "没有错误",
			TimeLimit:  1000,
			MemLimit:   1000,
			SubmitCode: code,
			MemUsage:   0,
			TimeUsage:  0,
		}
		global.DB.Create(&record)

	}
}

func main() {

	initialize.InitLogger() //初始化日志
	initialize.InitConfig() //读取配置信息
	initialize.InitDB()     //初始化MySQL

	insertData() // 建表并插入数据

}
