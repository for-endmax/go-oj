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

// RecordModel 记录表
type RecordModel struct {
	BaseModel
	UID        int32  `gorm:"column:u_id"`
	QID        int32  `gorm:"column:q_id"`
	Lang       string `gorm:"column:lang"`
	Status     int32  `gorm:"column:status"`
	ErrCode    int32  `gorm:"column:err_code"`
	ErrMsg     string `gorm:"column:err_msg"`
	TimeLimit  int32  `gorm:"column:time_limit"`
	MemLimit   int32  `gorm:"column:mem_limit"`
	SubmitCode string `gorm:"column:submit_code"`
}
