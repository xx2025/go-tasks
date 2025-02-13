package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/mem"
	"go-tasks/master/model"
	"go-tasks/utils/response"
)

func Dashboard(c *gin.Context) {
	nodeModel := model.NewNodeModel()
	nodeCount := nodeModel.GetNodeCount()

	projectModel := model.NewProjectModel()
	projectCount := projectModel.GetProjectCount()

	taskModel := model.NewTaskModel()
	taskCount := taskModel.GetTaskNum()

	processModel := model.NewProcessModel()
	processCount := processModel.GetProcessNum()

	var totalMem, freeMem string = "", ""
	vmStat, err := mem.VirtualMemory()
	if err == nil {
		totalMem = fmt.Sprintf("%.2f", float64(vmStat.Total)/(1024*1024*1024))
		freeMem = fmt.Sprintf("%.2f", float64(vmStat.Free)/(1024*1024*1024))

	}

	data := make(map[string]interface{})

	data["nodeCount"] = nodeCount
	data["projectCount"] = projectCount
	data["taskCount"] = taskCount
	data["processCount"] = processCount
	data["cpu"] = ""
	data["totalMem"] = totalMem
	data["freeMem"] = freeMem

	response.OkWithData("ok", data, c)
}
