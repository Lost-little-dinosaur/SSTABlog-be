package articles

import (
	"github.com/wujunyi792/crispy-waffle-be/internal/db"
	"github.com/wujunyi792/crispy-waffle-be/internal/dto/article"
	"github.com/wujunyi792/crispy-waffle-be/internal/logger"
	"github.com/wujunyi792/crispy-waffle-be/internal/model/Mysql"
	"gorm.io/gorm"
	"sync"
)

var dbManage *ArticleDBManage = nil

func init() {
	logger.Info.Println("[ catalogues ]start init Table ...")
	dbManage = GetManage()
}

type ArticleDBManage struct {
	mDB     *db.MainGORM
	sDBLock sync.RWMutex
}

func (m *ArticleDBManage) getGOrmDB() *gorm.DB {
	return m.mDB.GetDB()
}

func (m *ArticleDBManage) atomicDBOperation(op func()) {
	m.sDBLock.Lock()
	op()
	m.sDBLock.Unlock()
}

func GetManage() *ArticleDBManage {
	if dbManage == nil {
		var catalogueDb = db.MustCreateGorm()
		err := catalogueDb.GetDB().AutoMigrate(&Mysql.Article{}) //自动创建表
		if err != nil {
			logger.Error.Fatalln(err)
			return nil
		}
		dbManage = &ArticleDBManage{mDB: catalogueDb}
	}
	return dbManage
}

//以上代码是初始化数据库表以及自动创建表所需的代码，下面是查询数据库表的代码

func AddArticle(article *Mysql.Article) (error, string) {
	return GetManage().getGOrmDB().Model(&Mysql.Article{}).Create(article).Error, article.ID
}

func SearchArticlesTitle(title string) ([]Mysql.Article, error) {
	articles := make([]Mysql.Article, 0)
	if err := GetManage().getGOrmDB().Model(&Mysql.Article{}).Where("title like ?", "%"+title+"%").Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func SearchArticlesContent(content string) ([]Mysql.Article, error) {
	articles := make([]Mysql.Article, 0)
	if err := GetManage().getGOrmDB().Model(&Mysql.Article{}).Where("content like ?", "%"+content+"%").Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func SearchArticlesDescription(description string) ([]Mysql.Article, error) {
	articles := make([]Mysql.Article, 0)
	if err := GetManage().getGOrmDB().Model(&Mysql.Article{}).Where("description like ?", "%"+description+"%").Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func UpdateArticle(article *Mysql.Article, id string) error {
	return GetManage().getGOrmDB().Model(&Mysql.Article{}).Where("id = ?", id).
		Updates(article). //根据结构体批量更新数据
		Error
}

func DeleteArticle(id string, uid string) error {
	return GetManage().getGOrmDB().Model(&Mysql.Article{}).Where("id = ?", id).
		Update("last_modifier", uid). //更新最后修改人
		Delete(&Mysql.Article{}).Error
}

func DeleteArticlesByCatalogueID(catalogueID string, uid string) error {
	return GetManage().getGOrmDB().Model(&Mysql.Article{}).Where("catalogue_id = ?", catalogueID).
		Update("last_modifier", uid). //更新最后修改人
		Delete(&Mysql.Article{}).Error
}

func CheckArticleExistByCatalogueIDAndTitle(catalogueID string, title string) bool {
	tempArticle := &Mysql.Article{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Article{}).Where("catalogue_id = ? and title = ?", catalogueID, title).First(tempArticle).Error; err != nil {
		return false
	}
	return true
}
func CheckArticleExistByID(id string) (bool, string) {
	tempArticle := &Mysql.Article{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Article{}).Where("id = ?", id).First(tempArticle).Error; err != nil {
		return false, ""
	}
	return true, tempArticle.CreateBy
}

func GetArticleByID(id string) (*Mysql.Article, error) {
	tempArticle := &Mysql.Article{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Article{}).Where("id = ?", id).First(tempArticle).Error; err != nil {
		return nil, err
	}
	return tempArticle, nil
}

func GetArticleCatalogueIDAndTitleByID(id string) (string, string, error) {
	tempArticle := &Mysql.Article{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Article{}).Where("id = ?", id).First(tempArticle).Error; err != nil {
		return "", "", err
	}
	return tempArticle.CatalogueID, tempArticle.Title, nil
}

func GetArticlesByCatalogueID(catalogueID string) ([]Mysql.Article, error) {
	articles := make([]Mysql.Article, 0)
	if err := GetManage().getGOrmDB().Model(&Mysql.Article{}).Where("catalogue_id = ?", catalogueID).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func GetDeletedArticlesByCatalogueID(catalogueID string) ([]Mysql.Article, error) {
	articles := make([]Mysql.Article, 0)
	if err := GetManage().getGOrmDB().Unscoped().Model(&Mysql.Article{}).Where("catalogue_id = ? AND deleted_at IS NOT NULL", catalogueID).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func GetDeletedArticleInfo(uid string) ([]article.GetArticleInfoResponse, error) { //获取已删除的文章，自己只能看到自己删除的文章
	articles := make([]Mysql.Article, 0)
	if err := GetManage().getGOrmDB().Unscoped().Model(&Mysql.Article{}).Where("last_modifier = ? AND deleted_at IS NOT NULL", uid).Find(&articles).Error; err != nil {
		return nil, err
	}
	var tempReturnArr []article.GetArticleInfoResponse
	for _, v := range articles {
		tempReturnArr = append(tempReturnArr, article.GetArticleInfoResponse{

			CatalogueID:   v.CatalogueID,
			Title:         v.Title,
			Description:   v.Description,
			CreateBy:      v.CreateBy,
			LastModifier:  v.LastModifier,
			Cover:         v.Cover,
			ID:            v.ID,
			CommentNumber: v.CommentNumber,
			PraiseNumber:  v.PraiseNumber,
		})
	}
	return tempReturnArr, nil
}

func CheckIfArticleDeleted(uid string, id string) error {
	tempArticle := &Mysql.Article{}
	if err := GetManage().getGOrmDB().Unscoped().Model(&Mysql.Article{}).Where("last_modifier = ? AND deleted_at IS NOT NULL  AND id = ?", uid, id).First(tempArticle).Error; err != nil {
		return err
	}
	return nil
}

func DeleteArticleForever(id string) error {
	return GetManage().getGOrmDB().Model(&Mysql.Article{}).Unscoped().Where("id = ?", id).Delete(&Mysql.Article{}).Error
}

func RestoreArticle(id string) error {
	return GetManage().getGOrmDB().Model(&Mysql.Article{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error
}
