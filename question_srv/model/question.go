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

// Question 定义用户信息
type Question struct {
	BaseModel
	Seq     int32  `gorm:"type: int comment '序号';index: idx_seq,unique"`
	Name    string `gorm:"type: varchar(20) comment '名称';index: idx_name,unique"`
	Content string `gorm:"type: text comment '名称'"`
}

// Test 测试信息表
type Test struct {
	BaseModel
	Question  Question `gorm:"foreignKey:q_id"`
	QID       int32    `gorm:"type: int comment '对应的题目id';column:q_id"` // 对应的问题ID
	TimeLimit int32    `gorm:"type: int comment '时间上限/毫秒'"`
	MemLimit  int32    `gorm:"type: int comment '内存上限/KB'"`
	Input     string   `gorm:"type: text comment '测试JSON数据输入'"`
	Output    string   `gorm:"type: text comment '测试JSON数据输出'"`
}
