package controllers

import (
	"github.com/gin-gonic/gin"
	"go-tasks/boot/global"
	"go-tasks/master/model"
	"go-tasks/master/moduls/request"
	"go-tasks/utils"
	"go-tasks/utils/response"
)

func ProjectList(c *gin.Context) {
	page := utils.StrToInt(c.DefaultQuery("page", "1"))
	pageSize := utils.StrToInt(c.DefaultQuery("pageSize", "20"))
	name := c.DefaultQuery("name", "")

	projectModel := model.NewProjectModel()
	data := projectModel.GetListByPage(c, page, pageSize, name)
	response.OkWithData("ok", data, c)
}

func ProjectSave(c *gin.Context) {
	requestData := request.ProjectSave{}

	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}

	if requestData.Name == "" {
		response.Fail("名称不能为空", c)
		return
	}
	if requestData.Describe == "" {
		response.Fail("描述信息不能为空", c)
		return
	}
	projectModel := model.NewProjectModel()

	var project *model.ProjectModel
	if requestData.Id <= 0 {
		project = &model.ProjectModel{}
	} else {
		project, err = projectModel.FindById(requestData.Id, c)
		if err != nil {
			response.Fail("未查询到项目信息", c)
			return
		}

	}

	project.Name = requestData.Name
	project.Describe = requestData.Describe

	err = global.DB.WithContext(c).Save(project).Error
	if err != nil {
		response.Fail("保存失败", c)
		utils.Logger().Error(err.Error())
		return
	}
	response.Ok("保存成功", c)
}

func ProjectDelete(c *gin.Context) {
	requestData := request.ProjectSave{}

	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}

	projectModel := model.NewProjectModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	project, err := projectModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到项目信息", c)
		return
	}
	err = global.DB.WithContext(c).Delete(project).Error
	if err != nil {
		response.Fail("删除失败", c)
		utils.Logger().Error(err.Error())
		return
	}
	response.Ok("删除成功", c)
}

func ProjectSelector(c *gin.Context) {
	project := model.NewProjectModel()
	data := project.GetSelector(c)
	response.OkWithData("ok", data, c)
}
