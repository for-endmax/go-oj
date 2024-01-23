package model

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel 公共字段
type BaseModel struct {
	ID        int32     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"` //column 定义别名
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
}

// Question 定义题目信息
type Question struct {
	BaseModel
	Seq     int32  `gorm:"type: int comment '序号'"`
	Name    string `gorm:"type: varchar(20) comment '名称'"`
	Content string `gorm:"type: text comment '名称'"`
	Tests   []Test `gorm:"foreignKey:q_id"`
}

// Test 测试信息表
type Test struct {
	BaseModel
	QID    int32  `gorm:"type: int comment '对应的题目id';column:q_id"` // 对应的问题ID
	Input  string `gorm:"type: text comment '测试JSON数据输入'"`
	Output string `gorm:"type: text comment '测试JSON数据输出'"`
}
