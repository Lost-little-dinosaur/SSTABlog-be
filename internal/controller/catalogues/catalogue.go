package catalogues

import (
	"github.com/wujunyi792/crispy-waffle-be/internal/controller/articles"
	"github.com/wujunyi792/crispy-waffle-be/internal/db"
	"github.com/wujunyi792/crispy-waffle-be/internal/logger"
	"github.com/wujunyi792/crispy-waffle-be/internal/model/Mysql"
	"gorm.io/gorm"
	"sync"
)

var dbManage *CatalogueDBManage = nil

func init() {
	logger.Info.Println("[ catalogues ]start init Table ...")
	dbManage = GetManage()
}

type CatalogueDBManage struct {
	mDB     *db.MainGORM
	sDBLock sync.RWMutex
}

func (m *CatalogueDBManage) getGOrmDB() *gorm.DB {
	return m.mDB.GetDB()
}

func (m *CatalogueDBManage) atomicDBOperation(op func()) {
	m.sDBLock.Lock()
	op()
	m.sDBLock.Unlock()
}

func GetManage() *CatalogueDBManage {
	if dbManage == nil {
		var catalogueDb = db.MustCreateGorm()
		err := catalogueDb.GetDB().AutoMigrate(&Mysql.Catalogue{}) //自动创建表
		if err != nil {
			logger.Error.Fatalln(err)
			return nil
		}
		dbManage = &CatalogueDBManage{mDB: catalogueDb}
	}
	return dbManage
}

//以上代码是初始化数据库表以及自动创建表所需的代码，下面是查询数据库表的代码

func AddCatalogue(catalogue *Mysql.Catalogue) (error, string) {
	return GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Create(catalogue).Error, catalogue.ID
}

func GetCatalogue(id string) (*Mysql.Catalogue, error) {
	catalogue := &Mysql.Catalogue{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("id = ?", id).First(catalogue).Error; err != nil {
		return nil, err
	}
	return catalogue, nil
}

//func GetCatalogueByName(name string) (*Mysql.Catalogue, error) { //目录名字也是不能重复的
//	catalogue := &Mysql.Catalogue{}
//	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("catalogue_name = ?", name).First(catalogue).Error; err != nil {
//		return nil, err
//	}
//	return catalogue, nil
//}

func CheckCatalogueExist(id string) *Mysql.Catalogue { //返回true表示存在
	catalogue := &Mysql.Catalogue{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("id = ?", id).First(catalogue).Error; err != nil {
		return nil
	}
	return catalogue
}

func GetCatalogueSons(id string) ([]Mysql.Catalogue, error) {
	var catalogues []Mysql.Catalogue
	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("father_id = ?", id).Order("catalogue_name").Find(&catalogues).Error; err != nil { //默认对目录按照目录名称排序
		return nil, err
	}
	return catalogues, nil
}

func CheckCatalogueExistByName(catalogueName string, id string) bool { //目录名字也是不能重复的,返回true表示存在
	catalogue := &Mysql.Catalogue{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("catalogue_name = ? AND father_id = ?", catalogueName, id).First(catalogue).Error; err != nil {
		return false
	}
	return true
}

func RenameCatalogue(id string, newCatalogueName string, uid string) error { //重命名前先检查是否存在以及是否同级目录下是否重名
	return GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("id = ?", id).Update("catalogue_name", newCatalogueName).Update("last_modifier", uid).Error
}

func DeleteCatalogue(id string, uid string) error { //todo 回收站功能
	err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("id = ?", id).
		Update("last_modifier", uid). //更新最后修改人
		Delete(&Mysql.Catalogue{}).Error
	if err != nil {
		return err
	}
	//删除目录下的文章
	err = articles.DeleteArticlesByCatalogueID(id, uid)
	if err != nil {
		return err
	}
	//删除子目录
	var catalogues []Mysql.Catalogue
	if err = GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("father_id = ?", id).Find(&catalogues).Error; err != nil {
		return err
	}

	//递归删除子目录的子目录
	for _, catalogue := range catalogues {
		if err = DeleteCatalogue(catalogue.ID, uid); err != nil {
			return err
		}
		err = articles.DeleteArticlesByCatalogueID(catalogue.ID, uid)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateCatalogueDescription(id string, newDescription string, uid string) error {
	return GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("id = ?", id).Update("description", newDescription).Update("last_modifier", uid).Error
}

func UpdateCatalogueFather(id string, newFatherID string, uid string) error {
	return GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("id = ?", id).Update("father_id", newFatherID).Update("last_modifier", uid).Error
}
func GetCatalogueFatherID(id string) (string, error) {
	catalogue := &Mysql.Catalogue{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("id = ?", id).First(catalogue).Error; err != nil {
		return "", err
	}
	return catalogue.FatherID, nil
}

func SearchCatalogueByName(keyword string) ([]Mysql.Catalogue, error) {
	var catalogues []Mysql.Catalogue
	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("catalogue_name LIKE ?", "%"+keyword+"%").Order("catalogue_name").Find(&catalogues).Error; err != nil { //默认对目录按照目录名称排序
		return nil, err
	}
	return catalogues, nil
}

func SearchCatalogueByDescription(keyword string) ([]Mysql.Catalogue, error) {
	var catalogues []Mysql.Catalogue
	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("description LIKE ?", "%"+keyword+"%").Order("catalogue_name").Find(&catalogues).Error; err != nil { //默认对目录按照目录名称排序
		return nil, err
	}
	return catalogues, nil
}

func GetCatalogueRoute(id string) ([]string, error) { //获取目录路径
	var tempCatalogue = &Mysql.Catalogue{}

	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("id = ?", id).First(tempCatalogue).Error; err != nil {
		return nil, err
	}
	route := []string{tempCatalogue.CatalogueName}
	for tempCatalogue.FatherID != "" {
		if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("id = ?", tempCatalogue.FatherID).First(tempCatalogue).Error; err != nil {
			return nil, err
		}
		route = append(route, tempCatalogue.CatalogueName)
	}
	return route, nil
}

func GetDeletedCatalogue(uid string) ([]Mysql.Catalogue, error) {
	var catalogues []Mysql.Catalogue
	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Unscoped().Where("last_modifier = ? AND deleted_at IS NOT NULL", uid).Order("catalogue_name").Find(&catalogues).Error; err != nil { //默认对目录按照目录名称排序
		return nil, err
	}
	return catalogues, nil
}

func CheckIfCatalogueDeleted(uid string, id string) error {
	catalogue := &Mysql.Catalogue{}
	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Unscoped().Where("last_modifier = ? AND deleted_at IS NOT NULL AND id = ?", uid, id).First(catalogue).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCatalogueForever(id string) error { //永久删除目录
	var err error
	err = GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Unscoped().Where("id = ?", id).Delete(&Mysql.Catalogue{}).Error
	if err != nil {
		return err
	}
	var tempArticleArr []Mysql.Article
	//删除当前目录下的文章
	tempArticleArr, err = articles.GetDeletedArticlesByCatalogueID(id)
	for _, tempArticle := range tempArticleArr {
		err = articles.DeleteArticleForever(tempArticle.ID)
		if err != nil {
			return err
		}
	}
	//删除当前目录下的子目录
	var tempCatalogueArr []Mysql.Catalogue
	err = GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Unscoped().Where("father_id = ?", id).Find(&tempCatalogueArr).Error
	if err != nil {
		return err
	}
	for _, tempCatalogue := range tempCatalogueArr {
		err = DeleteCatalogueForever(tempCatalogue.ID) //递归调用
		if err != nil {
			return err
		}
	}
	return nil
}

func RestoreCatalogue(id string) error { //恢复目录
	var err error
	err = GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error
	if err != nil {
		return err
	}
	//恢复当前目录下的文章
	var tempArticleArr []Mysql.Article
	tempArticleArr, err = articles.GetDeletedArticlesByCatalogueID(id)
	for _, tempArticle := range tempArticleArr {
		err = articles.RestoreArticle(tempArticle.ID)
		if err != nil {
			return err
		}
	}
	//恢复当前目录下的子目录
	var tempCatalogueArr []Mysql.Catalogue
	err = GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Unscoped().Where("father_id = ?", id).Find(&tempCatalogueArr).Error
	if err != nil {
		return err
	}
	for _, tempCatalogue := range tempCatalogueArr {
		err = RestoreCatalogue(tempCatalogue.ID) //递归调用
		if err != nil {
			return err
		}
	}
	return nil
}

//func CheckCatalogueValidForUpdateName(id string, newName string) bool { //返回true表示存在
//	catalogue := &Mysql.Catalogue{}
//	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("id != ? AND catalogue_name = ?", id, newName).First(catalogue).Error; err != nil {
//		return false
//	}
//	return true
//}

//func GetCatalogueSon(parentID string, catalogueName string) ([]Mysql.Catalogue, error) {
//	catalogue := []Mysql.Catalogue{}
//	if err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("parent_id = ? AND catalogue_name = ?", parentID, catalogueName).Find(catalogue).Error; err != nil {
//		return nil, err
//	}
//	return catalogue, nil
//}
