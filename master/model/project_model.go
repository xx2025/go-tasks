package model

import (
	"context"
	"go-tasks/boot/global"
	"go-tasks/master/model/base"
	"sync"
	"time"
)

type ProjectModel struct {
	Id        int       `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Describe  string    `gorm:"column:describe"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

var projectModelOnce sync.Once
var projectModel *ProjectModel

func (ProjectModel) TableName() string {
	return "project"
}

func NewProjectModel() *ProjectModel {
	projectModelOnce.Do(func() {
		projectModel = &ProjectModel{}
	})
	return projectModel
}

func (n *ProjectModel) FindById(id int, c context.Context) (*ProjectModel, error) {
	model := &ProjectModel{}
	err := global.DB.WithContext(c).Where("id = ?", id).First(model).Error

	return model, err
}

func (n *ProjectModel) GetListByPage(c context.Context, page, pageSize int, name string) *base.Pagination {
	query := global.DB.Model(&ProjectModel{})

	if name != "" {
		name = "%" + name + "%"
		query = query.Where("name LIKE ?", name)
	}

	offset := (page - 1) * pageSize
	list := []ProjectModel{}

	var total int64
	_ = query.Count(&total)
	res := query.Limit(pageSize).Offset(offset).Find(&list)

	if res.Error != nil {

	}

	taskNumModel := NewTaskModel()
	taskNumMap := taskNumModel.GetTaskNumByProject()

	processNumModel := NewProcessModel()
	processNumMap := processNumModel.GetProcessNumByProject()

	data := &base.Pagination{
		Page:  page,
		Total: total,
		Limit: pageSize,
	}
	for _, v := range list {

		item := map[string]interface{}{
			"id":         v.Id,
			"name":       v.Name,
			"describe":   v.Describe,
			"taskNum":    taskNumMap[v.Id],
			"processNum": processNumMap[v.Id],
			"createdAt":  v.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		data.Items = append(data.Items, item)
	}

	return data

}

func (n *ProjectModel) GetSelector(c context.Context) []interface{} {
	query := global.DB.Model(&ProjectModel{})

	list := []ProjectModel{}

	var total int64
	_ = query.Count(&total)
	res := query.Find(&list)

	if res.Error != nil {

	}

	data := make([]interface{}, 0, 10)

	for _, v := range list {
		item := map[string]interface{}{
			"id":   v.Id,
			"name": v.Name,
		}
		data = append(data, item)
	}

	return data

}

func (ProjectModel) GetProjectCount() int64 {
	var count int64
	_ = global.DB.Model(&ProjectModel{}).Count(&count).Error
	return count
}
