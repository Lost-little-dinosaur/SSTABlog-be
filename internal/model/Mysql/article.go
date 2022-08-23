package Mysql

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Article struct {
	Base
	Title         string `gorm:"type:varchar(90);not null;" json:"title"`
	Cover         string `gorm:"type:varchar(90);" json:"cover"`
	CreateBy      string `gorm:"type:varchar(90);not null;"`                    //创建者，和User表关联ID
	LastModifier  string `gorm:"type:varchar(90);not null;"`                    //最后修改者，和User表关联ID
	CatalogueID   string `gorm:"type:varchar(90);not null;" json:"catalogueID"` //所属目录，和Catalogue表关联ID
	Description   string `gorm:"type:varchar(255);" json:"description"`
	Content       string `gorm:"type:longtext;" json:"content"`
	CommentNumber int    `gorm:"type:text" json:"commentNumber"` //评论数，作为拓展功能
	PraiseNumber  int    `gorm:"type:int" json:"praiseNumber"`   //点赞数，作为拓展功能
	WatchTimes    int    `gorm:"type:int" json:"watchTimes"`     //浏览次数，作为拓展功能
}

func (c *Article) BeforeCreate(tx *gorm.DB) (err error) { //使用钩子函数
	c.ID = uuid.NewV4().String()
	return
}
