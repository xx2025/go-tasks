package model

import (
	"context"
	"go-tasks/boot/global"
	"go-tasks/master/model/base"
	"sync"
	"time"
)

type ProcessLogsModel struct {
	Id        int       `gorm:"column:id;primaryKey"`
	ProcessId int       `gorm:"column:process_id"`
	Message   string    `gorm:"column:message"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

var processLogsModelOnce sync.Once
var processLogsModel *ProcessLogsModel

func (ProcessLogsModel) TableName() string {
	return "process_logs"
}

func NewProcessLogsModel() *ProcessLogsModel {
	processLogsModelOnce.Do(func() {
		processLogsModel = &ProcessLogsModel{}
	})
	return processLogsModel
}

func (ProcessLogsModel) NewLog(processId int, message string) (int, error) {
	log := &ProcessLogsModel{
		ProcessId: processId,
		Message:   message,
	}
	err := global.DB.Save(log).Error
	if err != nil {
		return 0, err
	}
	return log.Id, nil
}

func (ProcessLogsModel) GetListByPage(c context.Context, page, pageSize int, processName string, processId *int) *base.Pagination {
	query := global.DB.Model(&ProcessLogsModel{})

	if processName != "" {
		processName = "%" + processName + "%"
		var processIds []int
		_ = global.DB.Model(&ProcessModel{}).WithContext(c).Where("name LIKE ?", processName).Pluck("id", &processIds).Error
		if len(processIds) <= 0 {
			processIds = []int{-99}
		}
		query = query.Where("process_id in ?", processIds)
	}

	if processId != nil {
		query = query.Where("process_id =?", *processId)
	}

	offset := (page - 1) * pageSize
	list := []ProcessLogsModel{}

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

	var processIds []int
	for _, v := range list {
		processIds = append(processIds, v.ProcessId)
	}

	processModel := NewProcessModel()
	processList := processModel.GetListWithIds(processIds)
	var pMap = make(map[int]ProcessModel)
	if len(processList) > 0 {
		for _, process := range processList {
			pMap[process.Id] = process
		}
	}

	for _, v := range list {
		processName := ""
		if process, ok := pMap[v.ProcessId]; ok {
			processName = process.Name
		}

		item := map[string]interface{}{
			"id":          v.Id,
			"processId":   v.ProcessId,
			"processName": processName,
			"message":     v.Message,
			"createdAt":   v.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt":   v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		data.Items = append(data.Items, item)
	}

	return data

}
