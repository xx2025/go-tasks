package model

import (
	"context"
	"go-tasks/boot/global"
	"go-tasks/master/model/base"
	"sync"
	"time"
)

type TaskLogsModel struct {
	Id        int       `gorm:"column:id;primaryKey"`
	TaskId    int       `gorm:"column:task_id"`
	Status    int       `gorm:"column:status"`
	Message   string    `gorm:"column:message"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

var taskLogsModelOnce sync.Once
var taskLogsModel *TaskLogsModel

func (TaskLogsModel) TableName() string {
	return "task_logs"
}

func NewTaskLogsModel() *TaskLogsModel {
	taskLogsModelOnce.Do(func() {
		taskLogsModel = &TaskLogsModel{}
	})
	return taskLogsModel
}

func (TaskLogsModel) NewLog(taskId int) (int, error) {
	log := &TaskLogsModel{
		TaskId: taskId,
	}
	err := global.DB.Save(log).Error
	if err != nil {
		return 0, err
	}
	return log.Id, nil
}

func (TaskLogsModel) FindById(id int) *TaskLogsModel {
	log := &TaskLogsModel{}
	err := global.DB.Where("id=?", id).First(log).Error
	if err != nil {
		return nil
	}
	return log
}

func (TaskLogsModel) GetListByPage(c context.Context, page, pageSize int, taskName string, status *int, taskId *int) *base.Pagination {
	query := global.DB.Model(&TaskLogsModel{})

	if taskName != "" {
		taskName = "%" + taskName + "%"
		var taskIds []int
		_ = global.DB.Model(&TaskModel{}).WithContext(c).Where("name LIKE ?", taskName).Pluck("id", &taskIds).Error
		if len(taskIds) <= 0 {
			taskIds = []int{-99}
		}
		query = query.Where("task_id in ?", taskIds)
	}
	if status != nil {
		query = query.Where("status =?", *status)
	}

	if taskId != nil {
		query = query.Where("task_id =?", *taskId)
	}

	offset := (page - 1) * pageSize
	list := []TaskLogsModel{}

	var total int64
	_ = query.Count(&total)
	res := query.Order("id desc").Limit(pageSize).Offset(offset).Find(&list)

	if res.Error != nil {

	}

	data := &base.Pagination{
		Page:  page,
		Total: total,
		Limit: pageSize,
	}

	var taskIds []int
	for _, v := range list {
		taskIds = append(taskIds, v.TaskId)
	}

	taskModel = NewTaskModel()
	taskList := taskModel.GetListWithIds(taskIds)
	var taskMap = make(map[int]TaskModel)
	if len(taskList) > 0 {
		for _, task := range taskList {
			taskMap[task.Id] = task
		}
	}

	for _, v := range list {
		taskName := ""
		if task, ok := taskMap[v.TaskId]; ok {
			taskName = task.Name
		}

		item := map[string]interface{}{
			"id":        v.Id,
			"taskId":    v.TaskId,
			"taskName":  taskName,
			"status":    v.Status,
			"message":   v.Message,
			"createdAt": v.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt": v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		data.Items = append(data.Items, item)
	}

	return data

}
