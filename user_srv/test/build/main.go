package main

import (
	"fmt"
	"strconv"
	"user_srv/global"
	"user_srv/initialize"
	"user_srv/model"
	"user_srv/utils"
)

// insertData 建表并插入数据
func insertData() {
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

func main() {

	initialize.InitLogger() //初始化日志
	initialize.InitConfig() //读取配置信息
	initialize.InitDB()     //初始化MySQL

	insertData() // 建表并插入数据

}
