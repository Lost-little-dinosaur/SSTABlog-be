package Mysql

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Catalogue struct {
	Base
	CatalogueName string `gorm:"type:varchar(90);not null;primarykey;" json:"catalogueName"`
	Description   string `gorm:"type:varchar(255);" json:"description"`
	CreateBy      string `gorm:"type:varchar(90);not null;"`   //创建者，和User表关联ID
	LastModifier  string `gorm:"type:varchar(90);not null;"`   //最后修改者，和User表关联ID
	FatherID      string `gorm:"type:varchar(90);primarykey;"` //父级目录，为空则为根目录，同级目录下不能有相同名字的目录
	//Sons          []Element `json:"-"`  //一对多，只需要在子节点中声明父亲即可，不需要在父节点中声明子节点
}

func (c *Catalogue) BeforeCreate(tx *gorm.DB) (err error) { //使用钩子函数
	c.ID = uuid.NewV4().String()
	return
}
