package draftBoxs

import (
	"github.com/wujunyi792/crispy-waffle-be/internal/db"
	"github.com/wujunyi792/crispy-waffle-be/internal/logger"
	"github.com/wujunyi792/crispy-waffle-be/internal/model/Mysql"
	"gorm.io/gorm"
	"sync"
)

var dbManage *DraftDBManage = nil

func init() {
	logger.Info.Println("[ catalogues ]start init Table ...")
	dbManage = GetManage()
}

type DraftDBManage struct {
	mDB     *db.MainGORM
	sDBLock sync.RWMutex
}

func (m *DraftDBManage) getGOrmDB() *gorm.DB {
	return m.mDB.GetDB()
}

func (m *DraftDBManage) atomicDBOperation(op func()) {
	m.sDBLock.Lock()
	op()
	m.sDBLock.Unlock()
}

func GetManage() *DraftDBManage {
	if dbManage == nil {
		var catalogueDb = db.MustCreateGorm()
		err := catalogueDb.GetDB().AutoMigrate(&Mysql.Draft{}) //自动创建表
		if err != nil {
			logger.Error.Fatalln(err)
			return nil
		}
		dbManage = &DraftDBManage{mDB: catalogueDb}
	}
	return dbManage
}

//以上代码是初始化数据库表以及自动创建表所需的代码，下面是查询数据库表的代码

func AddDraft(draft *Mysql.Draft) error {
	return GetManage().getGOrmDB().Model(&Mysql.Draft{}).Create(draft).Error
}

func SearchDraftsTitle(title string) ([]Mysql.Draft, error) {
	drafts := make([]Mysql.Draft, 0)
	if err := GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("title like ?", "%"+title+"%").Find(&drafts).Error; err != nil {
		return nil, err
	}
	return drafts, nil
}

func SearchDraftsContent(content string) ([]Mysql.Draft, error) {
	drafts := make([]Mysql.Draft, 0)
	if err := GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("content like ?", "%"+content+"%").Find(&drafts).Error; err != nil {
		return nil, err
	}
	return drafts, nil
}

func SearchDraftsDescription(description string) ([]Mysql.Draft, error) {
	drafts := make([]Mysql.Draft, 0)
	if err := GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("description like ?", "%"+description+"%").Find(&drafts).Error; err != nil {
		return nil, err
	}
	return drafts, nil
}

func UpdateDraft(draft *Mysql.Draft, id string) error {
	return GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("id = ?", id).
		Updates(draft). //根据结构体批量更新数据
		Error
}

func DeleteDraft(id string, uid string) error {
	return GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("id = ?", id).
		Update("last_modifier", uid). //更新最后修改人
		Delete(&Mysql.Draft{}).Error
}

func DeleteDraftsByCatalogueID(catalogueID string, uid string) error {
	return GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("catalogue_id = ?", catalogueID).
		Update("last_modifier", uid). //更新最后修改人
		Delete(&Mysql.Draft{}).Error
}

func CheckDraftExistByCatalogueIDAndTitle(catalogueID string, title string) bool {
	tempDraft := &Mysql.Draft{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("catalogue_id = ? and title = ?", catalogueID, title).First(tempDraft).Error; err != nil {
		return false
	}
	return true
}
func CheckDraftExistByID(id string) (bool, string) {
	tempDraft := &Mysql.Draft{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("id = ?", id).First(tempDraft).Error; err != nil {
		return false, ""
	}
	return true, tempDraft.CreateBy
}

func GetDraftByID(id string) (*Mysql.Draft, error) {
	tempDraft := &Mysql.Draft{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("id = ?", id).First(tempDraft).Error; err != nil {
		return nil, err
	}
	return tempDraft, nil
}

func GetDraftByUID(uid string) ([]Mysql.Draft, error) {
	drafts := make([]Mysql.Draft, 0)
	if err := GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("create_by = ?", uid).Find(&drafts).Error; err != nil {
		return nil, err
	}
	return drafts, nil
}

//func GetDraftCatalogueIDAndTitleByID(id string) (string, string, error) {
//	tempDraft := &Mysql.Draft{}
//	if err := GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("id = ?", id).First(tempDraft).Error; err != nil {
//		return "", "", err
//	}
//	return tempDraft.CatalogueID, tempDraft.Title, nil
//}

func GetDraftsByCatalogueID(catalogueID string) ([]Mysql.Draft, error) {
	drafts := make([]Mysql.Draft, 0)
	if err := GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("catalogue_id = ?", catalogueID).Find(&drafts).Error; err != nil {
		return nil, err
	}
	return drafts, nil
}

func GetDeletedDraftsByCatalogueID(catalogueID string) ([]Mysql.Draft, error) {
	drafts := make([]Mysql.Draft, 0)
	if err := GetManage().getGOrmDB().Unscoped().Model(&Mysql.Draft{}).Where("catalogue_id = ? AND deleted_at IS NOT NULL", catalogueID).Find(&drafts).Error; err != nil {
		return nil, err
	}
	return drafts, nil
}

func GetAllArticleInfo(uid string) ([]Mysql.Draft, error) {
	drafts := make([]Mysql.Draft, 0)
	if err := GetManage().getGOrmDB().Model(&Mysql.Draft{}).Where("create_by = ?", uid).Find(&drafts).Error; err != nil {
		return nil, err
	}
	return drafts, nil
}

func DeleteDraftForever(id string) error {
	return GetManage().getGOrmDB().Model(&Mysql.Draft{}).Unscoped().Where("id = ?", id).Delete(&Mysql.Draft{}).Error
}
