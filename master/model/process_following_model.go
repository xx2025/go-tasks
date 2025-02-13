package model

import (
	"context"
	"go-tasks/boot/global"
	"sync"
	"time"
)

type ProcessFollowingModel struct {
	Id        int       `gorm:"column:id;primaryKey"`
	UserId    int       `gorm:"column:user_id"`
	ProcessId int       `gorm:"column:process_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

var processFollowingModelOnce sync.Once
var processFollowingModel *ProcessFollowingModel

func (ProcessFollowingModel) TableName() string {
	return "process_following"
}

func NewProcessFollowingModel() *ProcessFollowingModel {
	processFollowingModelOnce.Do(func() {
		processFollowingModel = &ProcessFollowingModel{}
	})
	return processFollowingModel
}

func (n *ProcessFollowingModel) FindByKey(userId, processId int, c context.Context) (*ProcessFollowingModel, error) {
	model := &ProcessFollowingModel{}
	err := global.DB.WithContext(c).Where("user_id = ? and process_id = ?", userId, processId).First(model).Error

	return model, err
}

func (n *ProcessFollowingModel) GetProcessIdsByUserId(userId int, c context.Context) ([]int, error) {
	var processIds []int
	err := global.DB.WithContext(c).Model(&ProcessFollowingModel{}).Where("user_id = ?", userId).Pluck("process_id", &processIds).Error

	return processIds, err
}
