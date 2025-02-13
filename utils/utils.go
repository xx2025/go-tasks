package utils

import (
	"github.com/gin-gonic/gin"
	"go-tasks/boot/config"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func GetImgHost() string {
	return config.MasterConfigure.Host + "/img/"
}

func Logger() *zap.Logger {
	return config.Logger
}

func GetAdminUserId(ctx *gin.Context) int {
	adminId, ok := ctx.Get("adminId")
	if ok {
		return adminId.(int)
	}
	return 0
}

func GetFormatDateTime() string {
	var datetime string = time.Now().Format("2006-01-02 15:04:05")
	return datetime
}

func StrToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}

func IntToStr(num int) string {
	str := strconv.Itoa(num)
	return str
}

func ContainsInSlice[T comparable](slice []T, target T) bool {
	for _, elem := range slice {
		if elem == target {
			return true
		}
	}
	return false
}
