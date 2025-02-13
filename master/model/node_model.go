package model

import (
	"context"
	"go-tasks/boot/global"
	"go-tasks/master/model/base"
	"gorm.io/gorm"
	"sync"
	"time"
)

type NodeModel struct {
	Id        int       `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Url       string    `gorm:"column:url"`
	Status    int       `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

var nodeModelOnce sync.Once
var nodeModel *NodeModel

func (NodeModel) TableName() string {
	return "node"
}

func NewNodeModel() *NodeModel {
	nodeModelOnce.Do(func() {
		nodeModel = &NodeModel{}
	})
	return nodeModel
}

func (n *NodeModel) FindById(id int, c context.Context) (*NodeModel, error) {
	nodeModel := &NodeModel{}
	err := global.DB.WithContext(c).Where("id = ?", id).First(nodeModel).Error

	return nodeModel, err
}

func (n *NodeModel) FindByUrl(url string, c context.Context) (*NodeModel, error) {
	nodeModel := &NodeModel{}
	err := global.DB.WithContext(c).Where("url = ?", url).First(nodeModel).Error

	return nodeModel, err
}

func (n *NodeModel) GetListByPage(c context.Context, page, pageSize int, name, url string, status *int) *base.Pagination {
	query := global.DB.Model(&NodeModel{})

	if name != "" {
		name = "%" + name + "%"
		query = query.Where("name LIKE ?", name)
	}

	if url != "" {
		url = "%" + url + "%"
		query = query.Where("url LIKE ?", url)
	}

	if status != nil {
		query = query.Where("status =?", *status)
	}

	offset := (page - 1) * pageSize
	list := []NodeModel{}

	var total int64
	_ = query.Count(&total)
	res := query.Limit(pageSize).Offset(offset).Find(&list)

	if res.Error != nil {

	}

	taskNumModel := NewTaskModel()
	taskNumMap := taskNumModel.GetTaskNumByNode()

	processNumModel := NewProcessModel()
	processNumMap := processNumModel.GetProcessNumByNode()

	data := &base.Pagination{
		Page:  page,
		Total: total,
		Limit: pageSize,
	}
	for _, v := range list {
		item := map[string]interface{}{
			"id":         v.Id,
			"name":       v.Name,
			"url":        v.Url,
			"taskNum":    taskNumMap[v.Id],
			"processNum": processNumMap[v.Id],
			"updatedAt":  v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		data.Items = append(data.Items, item)
	}

	return data

}

func (n *NodeModel) GetSelector(c context.Context) []interface{} {
	query := global.DB.Model(&NodeModel{})

	list := []NodeModel{}

	var total int64
	_ = query.Count(&total)
	res := query.Find(&list)

	if res.Error != nil {

	}

	data := make([]interface{}, 0, 1)

	for _, v := range list {
		item := map[string]interface{}{
			"id":        v.Id,
			"name":      v.Name,
			"url":       v.Url,
			"updatedAt": v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		data = append(data, item)
	}

	return data

}

func (NodeModel) GetNodeWithHost() map[int]string {
	nodeMap := make(map[int]string)
	list := []NodeModel{}
	_ = global.DB.Model(&NodeModel{}).Find(&list).Error
	if len(list) <= 0 {
		return nodeMap
	}
	for _, v := range list {
		nodeMap[v.Id] = v.Url
	}
	return nodeMap
}

func (n *NodeModel) InsertIfNotExsit(c context.Context, name, host string) (*NodeModel, error) {
	node, err := n.FindByUrl(host, c)
	if err == gorm.ErrRecordNotFound {
		node = &NodeModel{
			Name:   name,
			Url:    host,
			Status: 1,
		}
		err = global.DB.WithContext(c).Save(&node).Error
		return node, err
	}

	return node, err
}

func (NodeModel) GetNodeCount() int64 {
	var count int64
	_ = global.DB.Model(&NodeModel{}).Count(&count).Error
	return count
}
