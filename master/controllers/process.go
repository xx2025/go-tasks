package controllers

import (
	"github.com/gin-gonic/gin"
	"go-tasks/boot/global"
	"go-tasks/internal/scheduler"
	"go-tasks/master/model"
	"go-tasks/master/moduls/request"
	"go-tasks/utils"
	"go-tasks/utils/response"
	"gorm.io/gorm"
)

func ProcessList(c *gin.Context) {
	adminId := utils.GetAdminUserId(c)
	page := utils.StrToInt(c.DefaultQuery("page", "1"))
	pageSize := utils.StrToInt(c.DefaultQuery("pageSize", "20"))
	name := c.DefaultQuery("name", "")
	statusStr := c.DefaultQuery("status", "")
	nodeIdStr := c.DefaultQuery("nodeId", "")
	projectIdStr := c.DefaultQuery("projectId", "")
	following := c.DefaultQuery("following", "false")
	var status *int
	if statusStr != "" {
		statusInt := utils.StrToInt(statusStr)
		status = &statusInt
	}
	nodeId := utils.StrToInt(nodeIdStr)
	projectId := utils.StrToInt(projectIdStr)

	processModel := model.NewProcessModel()
	data := processModel.GetListByPage(c, page, pageSize, name, status, nodeId, projectId, following, adminId)
	response.OkWithData("ok", data, c)
}

func ProcessSave(c *gin.Context) {
	requestData := request.ProcessSave{}

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

	if requestData.Status != 1 && requestData.Status != 0 {
		response.Fail("状态有误", c)
		return
	}

	if requestData.NodeId <= 0 {
		response.Fail("节点选择有误", c)
		return
	}
	nodeModel := model.NewNodeModel()
	node, _ := nodeModel.FindById(requestData.NodeId, c)
	if node.Id <= 0 {
		response.Fail("节点不存", c)
		return
	}

	if requestData.ProjectId <= 0 {
		response.Fail("项目选择有误", c)
		return
	}
	if requestData.MaxRetries < 0 {
		response.Fail("重试次数设置有误", c)
		return
	}
	projectModel := model.NewProjectModel()
	project, _ := projectModel.FindById(requestData.ProjectId, c)
	if project.Id <= 0 {
		response.Fail("项目不存在", c)
		return
	}

	if requestData.Describe == "" {
		response.Fail("进程描述不能为空", c)
		return
	}

	process := &model.ProcessModel{}
	if requestData.Id > 0 {
		processModel := model.NewProcessModel()
		process, _ = processModel.FindById(requestData.Id, c)
		if process.Id <= 0 {
			response.Fail("未查询到相关进程", c)
			return
		}
	}

	process.Name = requestData.Name
	process.Status = requestData.Status
	process.ProjectId = requestData.ProjectId
	process.NodeId = requestData.NodeId
	process.MaxRetries = requestData.MaxRetries
	process.Describe = requestData.Describe

	err = global.DB.WithContext(c).Save(process).Error
	if err != nil {
		response.Fail("保存失败", c)
		utils.Logger().Error(err.Error())
		return
	}
	if process.Status == 1 {
		_ = scheduler.ProcessStart(process.Id, process.MaxRetries, process.Name, node.Url)
	} else {
		_ = scheduler.ProcessStop(process.Id, node.Url)
	}

	response.Ok("保存成功", c)
}

func ProcessFollowing(c *gin.Context) {

	requestData := struct {
		Id int
	}{}
	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}
	processModel := model.NewProcessModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	process, err := processModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到进程信息", c)
		return
	}
	adminId := utils.GetAdminUserId(c)

	processFollowingModel := model.NewProcessFollowingModel()

	processFollowing, _ := processFollowingModel.FindByKey(adminId, process.Id, c)
	if processFollowing.Id > 0 {
		response.Fail("您已关注了该进程", c)
		return
	}
	processFollowing.UserId = adminId
	processFollowing.ProcessId = process.Id

	err = global.DB.WithContext(c).Save(&processFollowing).Error
	if err != nil {
		response.Fail("关注失败", c)
		utils.Logger().Error(err.Error())
		return
	}
	response.Ok("关注成功", c)

}

func ProcessUnFollowing(c *gin.Context) {
	requestData := struct {
		Id int
	}{}
	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}
	processModel := model.NewProcessModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	process, err := processModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到任务信息", c)
		return
	}
	adminId := utils.GetAdminUserId(c)

	processFollowingModel := model.NewProcessFollowingModel()
	processFollowing, _ := processFollowingModel.FindByKey(adminId, process.Id, c)
	if processFollowing.Id <= 0 {
		response.Fail("您未关注该进程", c)
		return
	}

	err = global.DB.WithContext(c).Delete(&processFollowing).Error
	if err != nil {
		response.Fail("操作失败", c)
		utils.Logger().Error(err.Error())
		return
	}
	response.Ok("操作成功", c)

}

func ProcessDelete(c *gin.Context) {
	requestData := struct {
		Id int
	}{}

	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}

	processModel := model.NewProcessModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	process, err := processModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到进程信息", c)
		return
	}

	txErr := global.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.WithContext(c).Delete(process).Error
		if err != nil {
			return err
		}

		err = tx.WithContext(c).Where("process_id = ?", process.Id).Delete(&model.ProcessFollowingModel{}).Error
		if err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		utils.Logger().Error("删除进程失败" + txErr.Error())
		response.Fail("删除进程失败", c)
		return
	}
	response.Ok("删除成功", c)

}

func ProcessStart(c *gin.Context) {
	requestData := struct {
		Id int
	}{}

	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}

	processModel := model.NewProcessModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	process, err := processModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到进程信息", c)
		return
	}

	//if process.Status == 1 {
	//	response.Fail("进程已启动", c)
	//	return
	//}

	if process.Status != 1 {
		process.Status = 1
		err = global.DB.WithContext(c).Save(process).Error
		if err != nil {
			response.Fail(err.Error(), c)
			return
		}
	}

	nodeModel := model.NewNodeModel()
	node, _ := nodeModel.FindById(process.NodeId, c)

	_ = scheduler.ProcessStart(process.Id, process.MaxRetries, process.Name, node.Url)
	response.Ok("启动成功", c)

}

func ProcessStop(c *gin.Context) {
	requestData := struct {
		Id int
	}{}

	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}

	processModel := model.NewProcessModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	process, err := processModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到进程信息", c)
		return
	}

	if process.Status == 1 {
		process.Status = 0
		err = global.DB.WithContext(c).Save(process).Error
		if err != nil {
			response.Fail(err.Error(), c)
			return
		}
	}

	//if process.Status != 1 {
	//	response.Fail("进程已暂停", c)
	//	return
	//}
	//
	//process.Status = 0
	//err = global.DB.WithContext(c).Save(process).Error
	//if err != nil {
	//	response.Fail(err.Error(), c)
	//	return
	//}

	nodeModel := model.NewNodeModel()
	node, _ := nodeModel.FindById(process.NodeId, c)

	_ = scheduler.ProcessStop(process.Id, node.Url)
	response.Ok("操作成功", c)

}

func ProcessDetail(c *gin.Context) {
	id := utils.StrToInt(c.DefaultQuery("id", ""))
	processModel := model.NewProcessModel()
	process, err := processModel.FindById(id, c)
	if err != nil {
		response.Fail(err.Error(), c)
		return
	}
	nodeModel := model.NewNodeModel()
	node, _ := nodeModel.FindById(process.NodeId, c)

	projectModel := model.NewProjectModel()
	project, _ := projectModel.FindById(process.ProjectId, c)

	pid, _ := scheduler.GetProcessPID(process.Id, node.Url)

	data := map[string]interface{}{
		"id":            process.Id,
		"name":          process.Name,
		"describe":      process.Describe,
		"status":        process.Status,
		"nodeName":      node.Name,
		"projectName":   project.Name,
		"runningStatus": pid > 0, //运行状态
		"pid":           pid,
	}
	response.OkWithData("获取成功", data, c)
}

func ProcessLogs(c *gin.Context) {
	page := utils.StrToInt(c.DefaultQuery("page", "1"))
	pageSize := utils.StrToInt(c.DefaultQuery("pageSize", "20"))
	processName := c.DefaultQuery("processName", "")
	processIdStr := c.DefaultQuery("processId", "")

	var processId *int
	if processIdStr != "" {
		taskIdInt := utils.StrToInt(processIdStr)
		processId = &taskIdInt
	}

	processLogsModel := model.NewProcessLogsModel()
	data := processLogsModel.GetListByPage(c, page, pageSize, processName, processId)
	response.OkWithData("获取成功", data, c)
}
