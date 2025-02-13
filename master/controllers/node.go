package controllers

import (
	"github.com/gin-gonic/gin"
	"go-tasks/boot/global"
	"go-tasks/internal/scheduler"
	"go-tasks/master/model"
	"go-tasks/master/moduls/request"
	"go-tasks/utils"
	"go-tasks/utils/response"
)

func NodeList(c *gin.Context) {
	page := utils.StrToInt(c.DefaultQuery("page", "1"))
	pageSize := utils.StrToInt(c.DefaultQuery("pageSize", "20"))
	name := c.DefaultQuery("name", "")
	url := c.DefaultQuery("url", "")
	statusStr := c.DefaultQuery("status", "")

	var status *int
	if statusStr != "" {
		statusInt := utils.StrToInt(statusStr)
		status = &statusInt
	}
	nodeModel := model.NewNodeModel()
	data := nodeModel.GetListByPage(c, page, pageSize, name, url, status)
	response.OkWithData("ok", data, c)
}

func NodeSave(c *gin.Context) {
	requestData := request.NodeSave{}

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
	if requestData.Url == "" {
		response.Fail("URL不能为空", c)
		return
	}
	nodeModel := model.NewNodeModel()

	var node *model.NodeModel
	if requestData.Id <= 0 {
		node = &model.NodeModel{}
	} else {
		node, err = nodeModel.FindById(requestData.Id, c)
		if err != nil {
			response.Fail("未查询到节点信息", c)
			return
		}

	}

	node.Name = requestData.Name
	node.Url = requestData.Url

	err = scheduler.NodePing(node.Url)
	if err != nil {
		response.Fail(err.Error(), c)
		return
	}

	err = global.DB.WithContext(c).Save(node).Error
	if err != nil {
		response.Fail("保存失败", c)
		utils.Logger().Error(err.Error())
		return
	}
	response.Ok("保存成功", c)
}

func NodeDelete(c *gin.Context) {
	requestData := model.NodeModel{}

	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}

	nodeModel := model.NewNodeModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	node, err := nodeModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到节点信息", c)
		return
	}
	taskModel := model.NewTaskModel()
	taskCount := taskModel.GetTaskCountByNode(node.Id)
	if taskCount > 0 {
		response.Fail("该节点有任务， 请先删除任务", c)
		return
	}

	err = global.DB.WithContext(c).Delete(node).Error
	if err != nil {
		response.Fail("删除失败", c)
		utils.Logger().Error(err.Error())
		return
	}
	response.Ok("删除成功", c)
}

func NodeHealth(c *gin.Context) {
	requestData := request.NodeSave{}

	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}

	nodeModel := model.NewNodeModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	node, err := nodeModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到节点信息", c)
		return
	}
	err = scheduler.NodePing(node.Url)
	if err != nil {
		response.Fail(err.Error(), c)
	} else {
		response.Ok("节点状态正常", c)
	}

}

func NodeSelector(c *gin.Context) {
	node := model.NewNodeModel()
	data := node.GetSelector(c)
	response.OkWithData("ok", data, c)
}
