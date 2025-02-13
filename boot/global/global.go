package global

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	Mod         string
	WorkerDir   string
	LogDir      string
	ResourceDir string
)

var DB *gorm.DB

var RouterMap map[string]Router

type Router struct {
	Method   string
	Path     string
	Handle   gin.HandlerFunc
	Desc     string
	RoleAuth int
}
