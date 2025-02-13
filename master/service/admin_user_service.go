package service

import (
	"context"
	"errors"
	"go-tasks/boot/global"
	"go-tasks/master/model"
	"go-tasks/master/moduls/request"
	"go-tasks/utils"
	"gorm.io/gorm"
)

type AdminUserService struct {
}

func NewAdminUserService() *AdminUserService {
	return &AdminUserService{}
}

func (this AdminUserService) GetList(c context.Context) []map[string]interface{} {
	//adminUser := &model.AdminUserModel{}

	query := global.DB.Model(&model.AdminUserModel{})

	list := []model.AdminUserModel{}

	var total int64
	_ = query.Count(&total)
	res := query.Find(&list)

	if res.Error != nil {

	}

	data := make([]map[string]interface{}, 0, 0)
	for _, v := range list {
		item := map[string]interface{}{
			"id":       v.Id,
			"username": v.Username,
			"nickname": v.Nickname,
		}
		data = append(data, item)
	}

	return data

}

func (this AdminUserService) GetListByPage(c context.Context, page, pageSize int, username, nickname string, roleId, status int) *Pagination {
	//adminUser := &model.AdminUserModel{}

	query := global.DB.Model(&model.AdminUserModel{})

	if username != "" {
		query = query.Where("username =?", username)
	}

	if nickname != "" {
		nickname = "%" + nickname + "%"
		query = query.Where("nickname LIKE?", nickname)
	}

	if status > 0 {
		query = query.Where("status =?", status)
	}
	if roleId > 0 {
		query = query.Where("role_id =?", roleId)
	}

	offset := (page - 1) * pageSize
	list := []model.AdminUserModel{}

	var total int64
	_ = query.Count(&total)
	res := query.Limit(pageSize).Offset(offset).Find(&list)

	if res.Error != nil {

	}

	data := &Pagination{
		Page:  page,
		Total: total,
		Limit: pageSize,
	}
	for _, v := range list {
		item := map[string]interface{}{
			"id":        v.Id,
			"username":  v.Username,
			"nickname":  v.Nickname,
			"status":    v.Status,
			"roleId":    v.RoleId,
			"roleName":  v.GetRoleName(),
			"avatar":    utils.GetImgHost() + v.Avatar,
			"updatedAt": v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		data.Items = append(data.Items, item)
	}

	return data

}

func (this AdminUserService) NewUser(username, nickname, password string, roleId, status int) error {

	adminUser := &model.AdminUserModel{}

	user, err := adminUser.FindByUsername(username)

	if user.Id > 0 {
		return errors.New("用户名已存在")
	}

	password, err = adminUser.NewPassword(password)
	if err != nil {
		return err
	}

	//user := &AdminUserModel{
	//	Username:  username,
	//	Nickname:  nickname,
	//	Password:  password,
	//	Status:    status,
	//	UserLevel: userLevel,
	//}
	user.Username = username
	user.Nickname = nickname
	user.Password = password
	user.Status = status
	user.RoleId = roleId

	err = global.DB.Save(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (this AdminUserService) SaveUser(c context.Context, data *request.UserSaveRequest) error {

	adminUserModel := model.NewAdminUserModel()
	err := adminUserModel.CheckUserRole(data.RoleId)
	if err != nil {
		return err
	}
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user, err := adminUserModel.FindByIdLocked(data.Id, tx, c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未查询到该用户")
		} else {
			return err
		}
	}

	user.Status = data.Status
	user.Nickname = data.Nickname
	user.RoleId = data.RoleId

	err = tx.Select("Status", "Nickname", "RoleId").Updates(user).Error
	if err != nil {
		return err
	}
	tx.Commit()
	return nil

}

func (this AdminUserService) DeleteUser(c context.Context, id int) error {

	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user, err := adminUserModel.FindByIdLocked(id, tx, c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未查询到该用户")
		} else {
			return err
		}
	}

	if user.IsRootUser() {
		return errors.New("该用户不能删除")
	}

	err = tx.Where("id = ?", user.Id).Delete(user).Error
	if err != nil {
		return err
	}
	tx.Commit()
	return nil

}

func (this AdminUserService) ResetUserPassword(c context.Context, id int, password string) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user, err := adminUserModel.FindByIdLocked(id, tx, c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未查询到该用户")
		} else {
			return err
		}
	}

	password, err = user.NewPassword(password)
	if err != nil {
		return err
	}

	user.Password = password

	err = tx.Select("Password").Updates(user).Error
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
