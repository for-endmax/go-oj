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
	IsDeleted bool
}

// User 定义用户信息
type User struct {
	BaseModel
	Password string `gorm:"type:varchar(200) comment '加密后的密码';not null"`
	Nickname string `gorm:"type: varchar(20) comment '昵称';index: idx_nickname,unique "`
	Gender   string `gorm:"column:gender;default:male;type:varchar(6) comment 'male表示男， female表示女'"`
	Role     int32  `gorm:"column:role;default:1;type:int comment '1表示用户， 2表示管理员'"`
}
