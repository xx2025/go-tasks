package controllers

import (
	"github.com/gin-gonic/gin"
	"go-tasks/master/model"
	"go-tasks/master/moduls/request"
	"go-tasks/utils"
	"go-tasks/utils/response"
	"strings"
)

func UserAdd(c *gin.Context) {
	user := &request.UserSave{}
	_ = c.ShouldBindJSON(&user)
	user.Username = strings.Trim(user.Username, "")
	user.Password = strings.Trim(user.Password, "")
	if user.Username == "" {
		response.Fail("用户名不能为空", c)
		return
	}

	if user.Password == "" {
		response.Fail("密码不能为空", c)
		return
	}

	if user.RoleId != model.Role_Id_1 &&
		user.RoleId != model.Role_Id_2 &&
		user.RoleId != model.Role_Id_3 {
		response.Fail("用户等级有误", c)
		return
	}

	err := adminUserService.NewUser(user.Username, user.Nickname, user.Password, user.RoleId, user.Status)

	if err != nil {
		response.Fail(err.Error(), c)
		return
	}

	response.Ok("新增成功", c)
}

func UserSelector(c *gin.Context) {
	data := adminUserService.GetList(c)
	response.OkWithData("ok", data, c)
}

func UserList(c *gin.Context) {

	page := utils.StrToInt(c.DefaultQuery("page", "1"))
	pageSize := utils.StrToInt(c.DefaultQuery("pageSize", "20"))
	username := c.DefaultQuery("username", "")
	nickname := c.DefaultQuery("nickname", "")
	roleId := utils.StrToInt(c.DefaultQuery("roleId", "0"))
	status := utils.StrToInt(c.DefaultQuery("status", "0"))

	data := adminUserService.GetListByPage(c, page, pageSize, username, nickname, roleId, status)
	response.OkWithData("ok", data, c)
}

func UserSave(c *gin.Context) {
	user := &request.UserSaveRequest{}
	_ = c.ShouldBindJSON(&user)
	if user.Id <= 0 {
		response.Fail("请求有误", c)
		return
	}

	err := adminUserService.SaveUser(c, user)
	if err != nil {
		response.Fail(err.Error(), c)
		return
	}

	response.Ok("操作成功", c)
}

func UserDelete(c *gin.Context) {
	user := &request.UserSaveRequest{}
	_ = c.ShouldBindJSON(&user)
	if user.Id <= 0 {
		response.Fail("请求有误", c)
		return
	}

	err := adminUserService.DeleteUser(c, user.Id)
	if err != nil {
		response.Fail(err.Error(), c)
		return
	}

	response.Ok("删除成功", c)
}

func UserPasswordReset(c *gin.Context) {
	reqData := &request.UserSaveRequest{}
	_ = c.ShouldBindJSON(&reqData)
	if reqData.Id <= 0 {
		response.Fail("请求有误", c)
		return
	}

	if reqData.Password == "" {
		response.Fail("密码不能为空", c)
		return
	}

	err := adminUserService.ResetUserPassword(c, reqData.Id, reqData.Password)
	if err != nil {
		response.Fail(err.Error(), c)
		return
	}

	response.Ok("删除成功", c)
}

func UserLogs(c *gin.Context) {
	page := utils.StrToInt(c.DefaultQuery("page", "1"))
	pageSize := utils.StrToInt(c.DefaultQuery("pageSize", "20"))
	userId := utils.StrToInt(c.DefaultQuery("userId", "0"))

	adminUserLogsModel := model.NewAdminUserLogsModel()

	data := adminUserLogsModel.GetLogs(c, page, pageSize, userId)
	response.OkWithData("ok", data, c)
}
