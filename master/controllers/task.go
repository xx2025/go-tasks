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

func TaskList(c *gin.Context) {
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

	taskModel := model.NewTaskModel()
	data := taskModel.GetListByPage(c, page, pageSize, name, status, nodeId, projectId, following, adminId)
	response.OkWithData("ok", data, c)
}

func TaskSave(c *gin.Context) {
	requestData := request.TaskSave{}

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
	if requestData.Spec == "" {
		response.Fail("频率不能为空", c)
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

	if !scheduler.IsCronFormat(requestData.Spec) {
		response.Fail("执行频率必须是合法的cron表达式", c)
		return
	}

	nodeModel := model.NewNodeModel()
	node, _ := nodeModel.FindById(requestData.NodeId, c)
	if node.Id <= 0 {
		response.Fail("节点不存在或者状态有误", c)
		return
	}

	if requestData.ProjectId <= 0 {
		response.Fail("项目选择有误", c)
		return
	}
	projectModel := model.NewProjectModel()
	project, _ := projectModel.FindById(requestData.ProjectId, c)
	if project.Id <= 0 {
		response.Fail("项目不存在", c)
		return
	}

	if requestData.Describe == "" {
		response.Fail("任务描述不能为空", c)
		return
	}

	task := &model.TaskModel{}
	if requestData.Id > 0 {
		taskModel := model.NewTaskModel()
		task, _ = taskModel.FindById(requestData.Id, c)
		if task.Id <= 0 {
			response.Fail("未查询到相关任务", c)
			return
		}
	}

	task.Name = requestData.Name
	task.Spec = requestData.Spec
	task.Status = requestData.Status
	task.ProjectId = requestData.ProjectId
	task.NodeId = requestData.NodeId
	task.Describe = requestData.Describe

	err = global.DB.WithContext(c).Save(task).Error
	if err != nil {
		response.Fail("保存失败", c)
		utils.Logger().Error(err.Error())
		return
	}

	if task.Status == 1 {
		_ = scheduler.TaskSchedule.UpdateTask(task.Id, task.Spec, task.Name, node.Url, int64(task.IsSingle))
	}

	if task.Status != 1 {
		scheduler.TaskSchedule.RemoveTask(task.Id)
	}

	response.Ok("保存成功", c)
}

func TaskFollowing(c *gin.Context) {

	requestData := struct {
		Id int
	}{}
	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}
	taskModel := model.NewTaskModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	task, err := taskModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到任务信息", c)
		return
	}
	adminId := utils.GetAdminUserId(c)

	taskFollowingModel := model.NewTaskFollowingModel()

	followingTask, _ := taskFollowingModel.FindByKey(adminId, task.Id, c)
	if followingTask.Id > 0 {
		response.Fail("您已关注了该任务", c)
		return
	}
	followingTask.UserId = adminId
	followingTask.TaskId = task.Id

	err = global.DB.WithContext(c).Save(&followingTask).Error
	if err != nil {
		response.Fail("关注失败", c)
		utils.Logger().Error(err.Error())
		return
	}
	response.Ok("关注成功", c)

}

func TaskUnFollowing(c *gin.Context) {
	requestData := struct {
		Id int
	}{}
	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}
	taskModel := model.NewTaskModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	task, err := taskModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到任务信息", c)
		return
	}
	adminId := utils.GetAdminUserId(c)

	taskFollowingModel := model.NewTaskFollowingModel()
	taskFollowing, _ := taskFollowingModel.FindByKey(adminId, task.Id, c)
	if taskFollowing.Id <= 0 {
		response.Fail("您未关注该任务", c)
		return
	}

	err = global.DB.WithContext(c).Delete(taskFollowing).Error
	if err != nil {
		response.Fail("操作失败", c)
		utils.Logger().Error(err.Error())
		return
	}
	response.Ok("操作成功", c)

}

func TaskDelete(c *gin.Context) {
	requestData := struct {
		Id int
	}{}

	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}

	taskModel := model.NewTaskModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	task, err := taskModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到任务信息", c)
		return
	}

	txErr := global.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.WithContext(c).Delete(task).Error
		if err != nil {
			return err
		}

		err = tx.WithContext(c).Where("task_id = ?", task.Id).Delete(&model.TaskFollowingModel{}).Error
		if err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		// 事务执行失败
		utils.Logger().Error("删除任务失败" + txErr.Error())
		response.Fail("删除失败", c)
		return
	}
	scheduler.TaskSchedule.RemoveTask(task.Id)
	response.Ok("删除成功", c)
}

func TaskExec(c *gin.Context) {
	requestData := struct {
		Id int
	}{}

	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}

	taskModel := model.NewTaskModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	task, err := taskModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到任务信息", c)
		return
	}

	nodeModel := model.NewNodeModel()
	node, err := nodeModel.FindById(task.NodeId, c)
	if err != nil {
		response.Fail("该任务节点异常", c)
		return
	}

	err = scheduler.TaskSchedule.ExecTask(task.Id, task.Name, node.Url, int64(task.IsSingle))
	if err != nil {
		response.Fail(err.Error(), c)
		return
	}
	response.Ok("操作成功", c)
}

func TaskStop(c *gin.Context) {
	requestData := struct {
		Id int
	}{}

	err := c.BindJSON(&requestData)
	if err != nil {
		response.Fail("请求数据有误", c)
		utils.Logger().Error(err.Error())
		return
	}

	taskModel := model.NewTaskModel()
	if requestData.Id <= 0 {
		response.Fail("id不能小于0", c)
		return
	}
	task, err := taskModel.FindById(requestData.Id, c)
	if err != nil {
		response.Fail("未查询到任务信息", c)
		return
	}

	nodeModel := model.NewNodeModel()
	node, err := nodeModel.FindById(task.NodeId, c)
	if err != nil {
		response.Fail("该任务节点异常", c)
		return
	}

	err = scheduler.TaskStop(task.Id, node.Url)
	if err != nil {
		response.Fail(err.Error(), c)
		return
	}
	response.Ok("操作成功", c)
}

func TaskDetail(c *gin.Context) {
	id := utils.StrToInt(c.DefaultQuery("id", ""))
	taskModel := model.NewTaskModel()
	task, err := taskModel.FindById(id, c)
	if err != nil {
		response.Fail(err.Error(), c)
		return
	}
	nodeModel := model.NewNodeModel()
	node, _ := nodeModel.FindById(task.NodeId, c)

	projectModel := model.NewProjectModel()
	project, _ := projectModel.FindById(task.ProjectId, c)

	scheduleState := 0
	_, err = scheduler.GetTaskCron(task.Id)
	if err == nil {
		scheduleState = 1 //调度中
	}

	pid, _ := scheduler.GetTaskPID(task.Id, node.Url)
	data := map[string]interface{}{
		"id":            task.Id,
		"name":          task.Name,
		"spec":          task.Spec,
		"describe":      task.Describe,
		"pid":           pid,
		"nodeName":      node.Name,
		"projectName":   project.Name,
		"scheduleState": scheduleState,
	}
	response.OkWithData("获取成功", data, c)
}

func TaskLogs(c *gin.Context) {
	page := utils.StrToInt(c.DefaultQuery("page", "1"))
	pageSize := utils.StrToInt(c.DefaultQuery("pageSize", "20"))
	taskName := c.DefaultQuery("taskName", "")
	statusStr := c.DefaultQuery("status", "")
	taskIdStr := c.DefaultQuery("taskId", "")
	var status *int
	if statusStr != "" {
		statusInt := utils.StrToInt(statusStr)
		status = &statusInt
	}

	var taskId *int
	if taskIdStr != "" {
		taskIdInt := utils.StrToInt(taskIdStr)
		taskId = &taskIdInt
	}

	taskLogsModel := model.NewTaskLogsModel()
	data := taskLogsModel.GetListByPage(c, page, pageSize, taskName, status, taskId)
	response.OkWithData("获取成功", data, c)
}
