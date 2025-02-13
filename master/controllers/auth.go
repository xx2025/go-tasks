package controllers

import (
	"github.com/gin-gonic/gin"
	"go-tasks/boot/config"
	"go-tasks/boot/global"
	"go-tasks/master/logic"
	"go-tasks/master/model"
	"go-tasks/utils"
	"go-tasks/utils/response"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	_ = c.ShouldBindJSON(&data)

	token, err := logic.AdminLogin(c, data.Username, data.Password)
	if err != nil {
		response.Fail(err.Error(), c)
		return
	}
	res := map[string]string{
		"token": token,
	}
	c.SetCookie("token", token, 3600, "/", "", false, true)
	response.OkWithData("登录成功", res, c)
}

func LoginOut(c *gin.Context) {
	adminId := utils.GetAdminUserId(c)
	adminUserModel := model.NewAdminUserModel()
	user, _ := adminUserModel.FindById(adminId, c)

	err := global.DB.Model(user).Unscoped().Update("login_time", 0).Error
	if err != nil {
		response.Fail("退出失败："+err.Error(), c)
		return
	}

	response.OkWithData("操作成功", nil, c)
}

func MyInfo(c *gin.Context) {
	adminId := utils.GetAdminUserId(c)

	adminUserModel := model.NewAdminUserModel()
	user, _ := adminUserModel.FindById(adminId, c)
	data := map[string]any{
		"username":  user.Username,
		"nickname":  user.Nickname,
		"avatar":    config.MasterConfigure.Host + "/" + user.Avatar,
		"roleName":  user.GetRoleName(),
		"status":    user.Status,
		"createdAt": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"updatedAt": user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	response.OkWithData("ok", data, c)
}

func SaveMyPassword(c *gin.Context) {
	data := struct {
		OldPassword     string `json:"oldPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail("请求数据有误", c)
		return
	}
	if data.OldPassword == "" {
		response.Fail("原密码不能为空", c)
		return
	}
	if data.ConfirmPassword == "" {
		response.Fail("新密码不能为空", c)
		return
	}
	adminId := utils.GetAdminUserId(c)
	adminUserModel := model.NewAdminUserModel()
	user, _ := adminUserModel.FindById(adminId, c)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.OldPassword))
	if err != nil {
		response.Fail("原密码错误", c)
		return
	}
	newPassword, err := user.NewPassword(data.ConfirmPassword)
	if err != nil {
		response.Fail("新密码不合格", c)
		return
	}
	err = global.DB.Model(user).Unscoped().Update("password", newPassword).Error
	if err != nil {
		response.Fail("修改密码失败", c)
		return
	}
	response.Ok("修改密码成功", c)
}

func SaveMyNickname(c *gin.Context) {
	data := struct {
		Nickname string `json:"nickname"`
	}{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail("请求数据有误", c)
		return
	}
	if data.Nickname == "" {
		response.Fail("昵称不能为空", c)
		return
	}

	adminId := utils.GetAdminUserId(c)
	adminUserModel := model.NewAdminUserModel()
	user, _ := adminUserModel.FindById(adminId, c)

	err = global.DB.Model(user).Unscoped().Update("nickname", data.Nickname).Error
	if err != nil {
		response.Fail("修改昵称失败", c)
		return
	}
	response.Ok("修改昵称成功", c)
}

func SaveMyAvatar(c *gin.Context) {
	data := struct {
		Avatar string `json:"avatar"`
	}{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail("请求数据有误", c)
		return
	}
	if data.Avatar == "" {
		response.Fail("请上传头像", c)
		return
	}

	avatar := data.Avatar
	adminId := utils.GetAdminUserId(c)
	adminUserModel := model.NewAdminUserModel()
	user, _ := adminUserModel.FindById(adminId, c)

	err = global.DB.Model(user).Unscoped().Update("avatar", avatar).Error
	if err != nil {
		response.Fail("修改头像失败", c)
		return
	}
	response.OkWithData("修改头像成功", map[string]string{
		"avatar": data.Avatar,
	}, c)
}
