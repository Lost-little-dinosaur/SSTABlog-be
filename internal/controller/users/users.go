package users

import (
	"SSTABlog-be/internal/db"
	"SSTABlog-be/internal/logger"
	"SSTABlog-be/internal/model/Mysql"
	"errors"
	"gorm.io/gorm"
	"sync"
	"time"
)

var dbManage *UserDBManage = nil

func init() {
	logger.Info.Println("[ USER ]start init Table ...")
	dbManage = GetManage()
}

type UserDBManage struct {
	mDB     *db.MainGORM
	sDBLock sync.RWMutex
}

func (m *UserDBManage) getGOrmDB() *gorm.DB {
	return m.mDB.GetDB()
}

func (m *UserDBManage) atomicDBOperation(op func()) {
	m.sDBLock.Lock()
	op()
	m.sDBLock.Unlock()
}

func GetManage() *UserDBManage {
	if dbManage == nil {
		var userDb = db.MustCreateGorm()
		err := userDb.GetDB().AutoMigrate(&Mysql.Permission{}, &Mysql.User{}, &Mysql.Oauth{})
		if err != nil {
			logger.Error.Fatalln(err)
			return nil
		}
		dbManage = &UserDBManage{mDB: userDb}
	}
	return dbManage
}

func SetLoginLog(id string, token string) {
	GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{"last_login_token": token, "last_login_time": time.Now()})
}

func CheckPhoneExist(phone string) bool {
	res := false
	GetManage().atomicDBOperation(func() {
		res = GetManage().getGOrmDB().Model(&Mysql.User{}).Where("phone = ?", phone).Find(&Mysql.User{}).RowsAffected > 0
	})
	return res
}

func CheckUserNameExist(username string) bool {
	res := false
	GetManage().atomicDBOperation(func() {
		res = GetManage().getGOrmDB().Model(&Mysql.User{}).Where("nick_name = ?", username).Find(&Mysql.User{}).RowsAffected > 0
	})
	return res
}

func GetEntityByGithubId(githubID int64) (exist bool, entity *Mysql.User) {
	entity = &Mysql.User{}
	GetManage().atomicDBOperation(func() {
		exist = GetManage().getGOrmDB().Debug().Model(entity).Where("id = (SELECT user_id from `oauths` WHERE `oauths`.`oauth_id` = ?)", githubID).Find(entity).RowsAffected > 0
	})
	return
}

func GetEntity(entity *Mysql.User) *Mysql.User {
	GetManage().atomicDBOperation(func() {
		GetManage().getGOrmDB().Where(entity).Find(entity)
	})
	return entity
}

func GetUserAllAuthApp(uid string) (res []Mysql.Oauth) {
	GetManage().getGOrmDB().Model(&Mysql.Oauth{}).Where("user_id = ?", uid).Find(&res)
	return
}

func AddGithubOauth(entity *Mysql.Oauth) error {
	return GetManage().getGOrmDB().Create(entity).Error
}

func GetUserByID(id string) (entity *Mysql.User) {
	entity = &Mysql.User{}
	GetManage().atomicDBOperation(func() {
		GetManage().getGOrmDB().Where("id = ?", entity.ID).Find(entity)
	})
	return entity
}

func RegisterUser(user *Mysql.User) (err error) {
	GetManage().atomicDBOperation(func() {
		err = GetManage().getGOrmDB().Create(user).Error
	})
	return
}

func UpdateAvatar(id string, avatar string) (err error) {
	tx := GetManage().getGOrmDB().Begin()
	defer tx.Commit()

	entity := Mysql.User{
		Base: Mysql.Base{
			ID: id,
		},
	}

	if tx.Where("id = ?", entity.ID).Find(&entity).RowsAffected == 0 {
		tx.Rollback()
		return errors.New("???????????????")
	}

	entity.Avatar = avatar

	err = tx.Model(&entity).Update("avatar", avatar).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("avatar", avatar)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("???????????????")
	}
	return nil
}

func UpdatePhone(id string, phone string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("phone", phone)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("???????????????")
	}
	return nil
}

func UpdatePasswordAndSalt(phone string, passwordHashed string, salt string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("phone = ?", phone).Updates(map[string]interface{}{"password": passwordHashed, "salt": salt})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("???????????????")
	}
	return nil
}

func UpdateUserName(id string, username string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("nick_name", username)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("???????????????")
	}
	return nil
}

func UpdateSex(id string, sex int) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("sex", sex)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("???????????????")
	}
	return nil
}

func UpdateSignature(id string, signature string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("signature", signature)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("???????????????")
	}
	return nil
}

func UpdateEmail(id string, email string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("email", email)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("???????????????")
	}
	return nil
}

func UpdateStatus(id string, status string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("status", status)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("???????????????")
	}
	return nil
}

func UpdateStudentID(id string, studentID string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("student_id", studentID)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("???????????????")
	}
	return nil
}

func PermissionAdd(id string, permissionName string) (err error) {
	var u Mysql.User
	u.ID = id
	_ = GetEntity(&u)
	if u.ID == "" {
		return errors.New("???????????????")
	}
	GetManage().atomicDBOperation(func() {
		err = GetManage().getGOrmDB().Model(&Mysql.User{
			Base: Mysql.Base{
				ID: u.ID,
			},
		}).
			Association("Permission").Append(&Mysql.Permission{
			PermissionName: permissionName,
		})
	})
	return err
}

func PermissionDel(id string, permissionName string) (err error) {
	var u Mysql.User
	u.ID = id
	_ = GetEntity(&u)
	if u.ID == "" {
		return errors.New("???????????????")
	}
	GetManage().atomicDBOperation(func() {
		err = GetManage().getGOrmDB().Model(&Mysql.User{
			Base: Mysql.Base{
				ID: u.ID,
			},
		}).
			Association("Permission").Delete(&Mysql.Permission{
			PermissionName: permissionName,
		})
	})
	return err
}

func PermissionClear(id string) (err error) {
	var u Mysql.User
	u.ID = id
	_ = GetEntity(&u)
	if u.ID == "" {
		return errors.New("???????????????")
	}
	GetManage().atomicDBOperation(func() {
		err = GetManage().getGOrmDB().Model(&Mysql.User{
			Base: Mysql.Base{
				ID: u.ID,
			},
		}).
			Association("Permission").Clear()
	})
	return err
}

//func PermissionCheck(id string, permission string) (exist bool) { //
//	var permissionEntity Mysql.Permission
//	err := GetManage().getGOrmDB().Model(&Mysql.User{
//		Base: Mysql.Base{
//			ID: id,
//		},
//	}).Where("permission_name = ?", permission).Association("Permission").Find(&permissionEntity)
//	if err != nil {
//		logger.Error.Println(err)
//		return false
//	}
//	if permissionEntity.PermissionName != "" {
//		return true
//	}
//	return
//}

func GetPermission(id string) []string {
	var permissionEntity []Mysql.Permission
	err := GetManage().getGOrmDB().Model(&Mysql.User{
		Base: Mysql.Base{
			ID: id,
		},
	}).Association("Permission").Find(&permissionEntity)
	if err != nil {
		logger.Error.Println(err)
		return nil
	}
	returnTemp := make([]string, len(permissionEntity))
	for i, v := range permissionEntity {
		returnTemp[i] = v.PermissionName
	}
	return returnTemp
}

func PermissionCheck(id string, permission string) (exist bool) {
	var permissionEntity Mysql.Permission
	err := GetManage().getGOrmDB().Model(&Mysql.User{
		Base: Mysql.Base{
			ID: id,
		},
	}).Association("Permission").Find(&permissionEntity)
	if err != nil {
		logger.Error.Println(err)
		return false
	}
	if permissionEntity.PermissionName != "" && permissionEntity.PermissionName <= permission { //?????????????????????????????????????????????????????????true
		return true
	} else {
		return false
	}
}

func GetUserNameByID(id string) (error, string) {
	var user Mysql.User
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).First(&user)
	return res.Error, user.NickName
}
