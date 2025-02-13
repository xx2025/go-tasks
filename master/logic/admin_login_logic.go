package logic

import (
	"context"
	"errors"
	"go-tasks/master/model"
	"go-tasks/utils"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func AdminLogin(ctx context.Context, username, password string) (string, error) {
	adminUserModel := model.NewAdminUserModel()

	username = strings.Trim(username, " ")
	password = strings.Trim(password, " ")
	if len(username) <= 0 {
		return "", errors.New("用户名不能为空")
	}
	if len(password) <= 0 {
		return "", errors.New("密码不能为空")
	}
	adminUser, err := adminUserModel.FindByUsername(username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(password))
	if err != nil {
		return "", errors.New("密码错误")
	}

	if !adminUser.IsValid() {
		return "", errors.New("该账号禁止登录")
	}

	userToken, loginTime, err := utils.NewToken(adminUser.Id, adminUser.RoleId)
	if err != nil {
		return "", err
	}

	err = adminUser.UpdateLoginTime(loginTime)
	if err != nil {
		return "", err
	}

	return userToken, nil
}
