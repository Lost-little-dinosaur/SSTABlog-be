package Mysql

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Draft struct {
	Base
	Title       string `gorm:"type:varchar(90);not null;" json:"title"`
	Cover       string `gorm:"type:varchar(90);" json:"cover"`
	CreateBy    string `gorm:"type:varchar(90);not null;"` //创建者，和User表关联ID
	Description string `gorm:"type:varchar(255);" json:"description"`
	Content     string `gorm:"type:longtext;" json:"content"`
}

func (c *Draft) BeforeCreate(tx *gorm.DB) (err error) { //使用钩子函数
	c.ID = uuid.NewV4().String()
	return
}
