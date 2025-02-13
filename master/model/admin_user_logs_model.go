package model

import (
	"github.com/gin-gonic/gin"
	"go-tasks/boot/global"
	"go-tasks/master/model/base"
	"sync"
	"time"
)

type AdminUserLogsModel struct {
	Id        int       `gorm:"column:id;primaryKey"`
	UserId    int       `gorm:"column:user_id"`
	Uri       string    `gorm:"column:uri"`
	ClientIp  string    `gorm:"column:client_ip"`
	Data      *string   `gorm:"column:data"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:created_at"`
}

var adminUserLogsModelOnce sync.Once
var adminUserLogsModel *AdminUserLogsModel

func (AdminUserLogsModel) TableName() string {
	return "admin_user_logs"
}

func NewAdminUserLogsModel() *AdminUserLogsModel {
	adminUserLogsModelOnce.Do(func() {
		adminUserLogsModel = &AdminUserLogsModel{}
	})
	return adminUserLogsModel
}

func (t *AdminUserLogsModel) NewLog(adminId int, requestBody *string, path string, c *gin.Context) {
	t.UserId = adminId
	t.Uri = path
	t.ClientIp = c.ClientIP()
	t.Data = requestBody
	err := global.DB.WithContext(c).Model(t).Save(t).Error
	if err != nil {

	}
}

func (t *AdminUserLogsModel) GetLogs(c *gin.Context, page, pageSize, userId int) interface{} {
	query := global.DB.Model(&AdminUserLogsModel{})

	if userId > 0 {
		query = query.Where("user_id =?", userId)
	}

	offset := (page - 1) * pageSize
	list := []AdminUserLogsModel{}

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
	userIdsMap := make(map[int]int)
	for _, v := range list {
		userIdsMap[v.UserId] = v.UserId
	}
	userIds := make([]int, 0, 0)
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}
	adminUserModel := NewAdminUserModel()
	userMap := adminUserModel.GetByIds(c, userIds)

	for _, v := range list {
		username := ""
		if user, ok := userMap[v.UserId]; ok {
			username = user.Username
		}

		item := map[string]interface{}{
			"id":        v.Id,
			"userId":    v.UserId,
			"username":  username,
			"data":      v.Data,
			"uri":       v.Uri,
			"ip":        v.ClientIp,
			"updatedAt": v.UpdatedAt.Format("2006-01-02 15:04:05"),
			"createdAt": v.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		data.Items = append(data.Items, item)
	}

	return data
}
