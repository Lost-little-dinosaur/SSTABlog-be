package catalogues

import (
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

func DeleteCatalogue(id string) error { //todo 回收站功能
	err := GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("id = ?", id).Delete(&Mysql.Catalogue{}).Error
	if err != nil {
		return err
	}
	//删除子目录
	var catalogues []Mysql.Catalogue
	if err = GetManage().getGOrmDB().Model(&Mysql.Catalogue{}).Where("father_id = ?", id).Find(&catalogues).Error; err != nil {
		return err
	}
	//todo 删除子目录下的文章

	//递归删除子目录的子目录
	for _, catalogue := range catalogues {
		if err = DeleteCatalogue(catalogue.ID); err != nil {
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
