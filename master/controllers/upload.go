package controllers

import (
	"github.com/gin-gonic/gin"
	"go-tasks/boot/global"
	"go-tasks/utils/response"
)

func UploadImg(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		response.Fail("获取图片失败", c)
		return
	}

	// 构建保存图片的目标路径，可以根据实际需求调整路径
	filePath := global.ResourceDir + "img/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response.Fail("保存图片失败", c)
		return
	}
	data := map[string]string{
		"url": "img/" + file.Filename,
	}
	response.OkWithData("上传成功", data, c)
}
