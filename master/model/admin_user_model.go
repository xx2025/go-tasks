package model

import (
	"context"
	"errors"
	"go-tasks/boot/global"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"time"
)

type AdminUserModel struct {
	Id        int `gorm:"column:id;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Username  string `gorm:"column:username"`
	Nickname  string `gorm:"column:nickname"`
	Password  string `gorm:"column:password"`
	RoleId    int    `gorm:"column:role_id"`
	Status    int    `gorm:"column:status"`
	Avatar    string `gorm:"column:avatar"`
	LoginTime int64  `gorm:"column:login_time"`
}

// 用户等级
const (
	Role_Id_1 int = iota + 1
	Role_Id_2
	Role_Id_3
)

var adminUserModelOnce sync.Once
var adminUserModel *AdminUserModel

func (AdminUserModel) TableName() string {
	return "admin_user"
}

func NewAdminUserModel() *AdminUserModel {
	adminUserModelOnce.Do(func() {
		adminUserModel = &AdminUserModel{}
	})
	return adminUserModel
}

func (t *AdminUserModel) IsRootUser() bool {
	return t.Username == "root"
}

func (t *AdminUserModel) CheckUserRole(roleId int) error {
	if roleId != Role_Id_1 &&
		roleId != Role_Id_2 &&
		roleId != Role_Id_3 {
		return errors.New("用户等级有误")
	}
	return nil
}

func (t *AdminUserModel) GetRoleName() string {
	if t.RoleId == Role_Id_1 {
		return "超级管理员"
	} else if t.RoleId == Role_Id_2 {
		return "管理员"
	} else {
		return "普通用户"
	}
}
func (t *AdminUserModel) IsValid() bool {
	if t.Status == 1 {
		return true
	} else {
		return false
	}
}

func (t *AdminUserModel) FindById(id int, c context.Context) (*AdminUserModel, error) {
	adminUser := &AdminUserModel{}

	err := global.DB.WithContext(c).Where("id = ?", id).First(adminUser).Error

	return adminUser, err
}

func (t *AdminUserModel) FindByIdLocked(id int, tx *gorm.DB, c context.Context) (*AdminUserModel, error) {
	adminUser := &AdminUserModel{}

	err := tx.WithContext(c).Clauses(
		clause.Locking{Strength: "UPDATE"},
	).Where("id = ?", id).First(adminUser).Error

	return adminUser, err
}

func (t *AdminUserModel) NewPassword(password string) (string, error) {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	password = string(passwordByte)
	return password, nil
}

func (t *AdminUserModel) FindByUsername(username string) (*AdminUserModel, error) {
	adminUser := AdminUserModel{}
	err := global.DB.Where("username = ?", username).First(&adminUser).Error
	return &adminUser, err
}

func (t *AdminUserModel) UpdateLoginTime(loginTime int64) error {
	t.LoginTime = loginTime
	err := global.DB.Select("LoginTime").Updates(t).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *AdminUserModel) GetByIds(c context.Context, ids []int) map[int]AdminUserModel {
	if len(ids) <= 0 {
		return map[int]AdminUserModel{}
	}

	list := []AdminUserModel{}
	err := global.DB.WithContext(c).Model(t).Where("id in ?", ids).Find(&list).Error
	if err != nil {
		return map[int]AdminUserModel{}
	}
	userMap := make(map[int]AdminUserModel)
	for _, v := range list {
		userMap[v.Id] = v
	}
	return userMap
}
