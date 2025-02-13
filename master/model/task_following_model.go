package model

import (
	"context"
	"go-tasks/boot/global"
	"sync"
	"time"
)

type TaskFollowingModel struct {
	Id        int       `gorm:"column:id;primaryKey"`
	UserId    int       `gorm:"column:user_id"`
	TaskId    int       `gorm:"column:task_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

var taskFollowingModelOnce sync.Once
var taskFollowingModel *TaskFollowingModel

func (TaskFollowingModel) TableName() string {
	return "task_following"
}

func NewTaskFollowingModel() *TaskFollowingModel {
	taskFollowingModelOnce.Do(func() {
		taskFollowingModel = &TaskFollowingModel{}
	})
	return taskFollowingModel
}

func (n *TaskFollowingModel) FindByKey(userId, taskId int, c context.Context) (*TaskFollowingModel, error) {
	model := &TaskFollowingModel{}
	err := global.DB.WithContext(c).Where("user_id = ? and task_id = ?", userId, taskId).First(model).Error

	return model, err
}

func (n *TaskFollowingModel) GetTaskIdsByUserId(userId int, c context.Context) ([]int, error) {
	var taskIds []int
	err := global.DB.WithContext(c).Model(&TaskFollowingModel{}).Where("user_id = ?", userId).Pluck("task_id", &taskIds).Error

	return taskIds, err
}
