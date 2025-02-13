package model

import (
	"context"
	"go-tasks/boot/global"
	"go-tasks/master/model/base"
	"go-tasks/utils"
	"sync"
	"time"
)

type ProcessModel struct {
	Id         int       `gorm:"column:id;primaryKey"`
	Name       string    `gorm:"column:name"`
	Status     int       `gorm:"column:status"`
	ProjectId  int       `gorm:"column:project_id"`
	NodeId     int       `gorm:"column:node_id"`
	MaxRetries int       `gorm:"column:max_retries"`
	Describe   string    `gorm:"column:describe"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

var processModelOnce sync.Once
var processModel *ProcessModel

func (ProcessModel) TableName() string {
	return "process"
}

func NewProcessModel() *ProcessModel {
	processModelOnce.Do(func() {
		processModel = &ProcessModel{}
	})
	return processModel
}

func (n *ProcessModel) FindById(id int, c context.Context) (*ProcessModel, error) {
	model := &ProcessModel{}
	err := global.DB.WithContext(c).Where("id = ?", id).First(model).Error

	return model, err
}

func (n *ProcessModel) GetListByPage(c context.Context, page, pageSize int, name string, status *int, nodeId, projectId int, following string, adminId int) *base.Pagination {
	query := global.DB.Model(&ProcessModel{})

	if name != "" {
		name = "%" + name + "%"
		query = query.Where("name LIKE ?", name)
	}
	if status != nil {
		query = query.Where("status =?", *status)
	}

	if nodeId > 0 {
		query = query.Where("node_id =?", nodeId)
	}
	if projectId > 0 {
		query = query.Where("project_id =?", projectId)
	}
	processFollowingModel := NewProcessFollowingModel()
	myFollowingProcessIds, _ := processFollowingModel.GetProcessIdsByUserId(adminId, c)
	if len(myFollowingProcessIds) <= 0 {
		myFollowingProcessIds = []int{-99}
	}

	if following == "true" {
		query = query.Where("id in ?", myFollowingProcessIds)
	}

	offset := (page - 1) * pageSize
	list := []ProcessModel{}

	var total int64
	_ = query.Count(&total)
	res := query.Limit(pageSize).Offset(offset).Find(&list)

	if res.Error != nil {

	}

	data := &base.Pagination{
		Page:  page,
		Total: total,
		Limit: pageSize,
	}
	for _, v := range list {
		item := map[string]interface{}{
			"id":          v.Id,
			"name":        v.Name,
			"projectId":   v.ProjectId,
			"nodeId":      v.NodeId,
			"status":      v.Status,
			"maxRetries":  v.MaxRetries,
			"describe":    v.Describe,
			"isFollowing": utils.ContainsInSlice(myFollowingProcessIds, v.Id),
			"createdAt":   v.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt":   v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		data.Items = append(data.Items, item)
	}

	return data

}

func (ProcessModel) GetListWithIds(ids []int) []ProcessModel {
	list := []ProcessModel{}
	_ = global.DB.Model(&ProcessModel{}).Where("id in ?", ids).Find(&list).Error

	return list
}

func (ProcessModel) GetAllProcess(nodeId int) []ProcessModel {
	query := global.DB.Model(&ProcessModel{}).Where("status = 1")

	if nodeId > 0 {
		query = query.Where("node_id = ?", nodeId)
	}

	list := []ProcessModel{}
	_ = query.Find(&list).Error

	return list
}

func (ProcessModel) GetProcessNum() int64 {
	var count int64
	_ = global.DB.Model(&ProcessModel{}).Count(&count).Error

	return count
}

func (ProcessModel) GetProcessNumByNode() map[int]int {
	var results []struct {
		NodeId     int `gorm:"column:node_id"`
		ProcessNum int `gorm:"column:process_num"`
	}
	_ = global.DB.Model(&ProcessModel{}).Select("node_id, COUNT(*) as process_num").Group("node_id").Find(&results).Error

	data := make(map[int]int)
	for _, v := range results {
		data[v.NodeId] = v.ProcessNum
	}
	return data
}

func (ProcessModel) GetProcessNumByProject() map[int]int {
	var results []struct {
		ProjectId  int `gorm:"column:project_id"`
		ProcessNum int `gorm:"column:process_num"`
	}
	_ = global.DB.Model(&ProcessModel{}).Select("project_id, COUNT(*) as process_num").Group("project_id").Find(&results).Error

	data := make(map[int]int)
	for _, v := range results {
		data[v.ProjectId] = v.ProcessNum
	}
	return data
}
