package model

import (
	"context"
	"go-tasks/boot/global"
	"go-tasks/master/model/base"
	"go-tasks/utils"
	"sync"
	"time"
)

type TaskModel struct {
	Id        int       `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Spec      string    `gorm:"column:spec"`
	Status    int       `gorm:"column:status"`
	IsSingle  int       `gorm:"column:is_single"`
	ProjectId int       `gorm:"column:project_id"`
	NodeId    int       `gorm:"column:node_id"`
	Describe  string    `gorm:"column:describe"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

var taskModelOnce sync.Once
var taskModel *TaskModel

func (TaskModel) TableName() string {
	return "task"
}

func NewTaskModel() *TaskModel {
	taskModelOnce.Do(func() {
		taskModel = &TaskModel{}
	})
	return taskModel
}

func (TaskModel) FindById(id int, c context.Context) (*TaskModel, error) {
	model := &TaskModel{}
	err := global.DB.WithContext(c).Where("id = ?", id).First(model).Error

	return model, err
}

func (n *TaskModel) GetListByPage(c context.Context, page, pageSize int, name string, status *int, nodeId, projectId int, following string, adminId int) *base.Pagination {
	query := global.DB.Model(&TaskModel{})

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
	taskFollowingModel := NewTaskFollowingModel()

	myFollowingTaskIds, _ := taskFollowingModel.GetTaskIdsByUserId(adminId, c)
	if len(myFollowingTaskIds) <= 0 {
		myFollowingTaskIds = []int{-99}
	}

	if following == "true" {
		query = query.Where("id in ?", myFollowingTaskIds)
	}

	offset := (page - 1) * pageSize
	list := []TaskModel{}

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
			"spec":        v.Spec,
			"projectId":   v.ProjectId,
			"nodeId":      v.NodeId,
			"status":      v.Status,
			"isSingle":    v.IsSingle,
			"describe":    v.Describe,
			"isFollowing": utils.ContainsInSlice(myFollowingTaskIds, v.Id),
			"createdAt":   v.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt":   v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		data.Items = append(data.Items, item)
	}

	return data

}

func (TaskModel) GetListWithSchedule() []TaskModel {
	list := []TaskModel{}
	_ = global.DB.Model(&TaskModel{}).Where("status = 1").Find(&list).Error
	return list
}

func (TaskModel) GetTaskCountByNode(nodeId int) int64 {
	var count int64
	_ = global.DB.Model(&TaskModel{}).Where("node_id = ?", nodeId).Count(&count).Error
	return count
}

func (TaskModel) GetListWithIds(ids []int) []TaskModel {
	list := []TaskModel{}
	_ = global.DB.Model(&TaskModel{}).Where("id in ?", ids).Find(&list).Error

	return list
}

func (TaskModel) GetTaskNum() int64 {
	var count int64
	_ = global.DB.Model(&TaskModel{}).Count(&count).Error

	return count
}

func (TaskModel) GetTaskNumByNode() map[int]int {
	var results []struct {
		NodeId  int `gorm:"column:node_id"`
		TaskNum int `gorm:"column:task_num"`
	}
	_ = global.DB.Model(&TaskModel{}).Select("node_id, COUNT(*) as task_num").Group("node_id").Find(&results).Error

	data := make(map[int]int)
	for _, v := range results {
		data[v.NodeId] = v.TaskNum
	}
	return data
}

func (TaskModel) GetTaskNumByProject() map[int]int {
	var results []struct {
		ProjectId int `gorm:"column:project_id"`
		TaskNum   int `gorm:"column:task_num"`
	}
	_ = global.DB.Model(&TaskModel{}).Select("project_id, COUNT(*) as task_num").Group("project_id").Find(&results).Error

	data := make(map[int]int)
	for _, v := range results {
		data[v.ProjectId] = v.TaskNum
	}
	return data
}
