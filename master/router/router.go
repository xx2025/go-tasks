package router

import (
	"github.com/gin-gonic/gin"
	"go-tasks/boot/config"
	"go-tasks/boot/global"
	"go-tasks/master/controllers"
	"go-tasks/master/middleware"
	"log"
	"strings"
)

func InitRouter() *gin.Engine {
	global.RouterMap = make(map[string]global.Router)

	routerList := getRouterList()
	r := gin.New()
	gin.SetMode(config.MasterConfigure.Mode)

	//r.Static("/img", global.ResourceDir)

	r.Use(middleware.Cors())
	r.Use(middleware.DefaultMiddleware())

	for _, route := range routerList {
		switch route.Method {
		case "GET":
			r.GET(route.Path, route.Handle)
		case "POST":
			r.POST(route.Path, route.Handle)
		default:
			log.Fatal("Routing can only be either get or post")
		}
		_, ok := global.RouterMap[route.Path]
		if ok {
			// 键存在，此时value为键对应的值
			log.Fatal("Route already registered on " + route.Path)
		} else {
			global.RouterMap[strings.Trim(route.Path, "/")] = route
		}
	}

	return r
}

func getRouterList() []global.Router {
	// 定义路由配置数组
	routes := []global.Router{
		{
			Method:   "POST",
			Path:     "/login",
			Handle:   controllers.Login,
			Desc:     "登录",
			RoleAuth: 0, // 0:不鉴权
		},
		{
			Method:   "POST",
			Path:     "/login/out",
			Handle:   controllers.LoginOut,
			Desc:     "退出登录",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/upload/img",
			Handle:   controllers.UploadImg,
			Desc:     "上传图片",
			RoleAuth: 0,
		},
		{
			Method:   "GET",
			Path:     "/my/info",
			Handle:   controllers.MyInfo,
			Desc:     "我的信息",
			RoleAuth: 0,
		},

		{
			Method:   "POST",
			Path:     "/save/my/password",
			Handle:   controllers.SaveMyPassword,
			Desc:     "更新密码",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/save/my/nickname",
			Handle:   controllers.SaveMyNickname,
			Desc:     "更新昵称",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/save/my/avatar",
			Handle:   controllers.SaveMyAvatar,
			Desc:     "更新头像",
			RoleAuth: 0,
		},

		{
			Method:   "GET",
			Path:     "/user/selector",
			Handle:   controllers.UserSelector,
			Desc:     "用户选择",
			RoleAuth: 0,
		},

		{
			Method:   "GET",
			Path:     "/user/list",
			Handle:   controllers.UserList,
			Desc:     "用户列表",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/user/add",
			Handle:   controllers.UserAdd,
			Desc:     "新增用户",
			RoleAuth: 1,
		},
		{
			Method:   "POST",
			Path:     "/user/save",
			Handle:   controllers.UserSave,
			Desc:     "编辑用户",
			RoleAuth: 1,
		},
		{
			Method:   "POST",
			Path:     "/user/delete",
			Handle:   controllers.UserDelete,
			Desc:     "删除用户",
			RoleAuth: 1,
		},
		{
			Method:   "POST",
			Path:     "/user/password/reset",
			Handle:   controllers.UserPasswordReset,
			Desc:     "重置密码",
			RoleAuth: 1,
		},

		{
			Method:   "GET",
			Path:     "/user/logs",
			Handle:   controllers.UserLogs,
			Desc:     "用户日志",
			RoleAuth: 0,
		},

		{
			Method:   "POST",
			Path:     "/node/add",
			Handle:   controllers.NodeSave,
			Desc:     "添加节点",
			RoleAuth: 1,
		},
		{
			Method:   "GET",
			Path:     "/node/list",
			Handle:   controllers.NodeList,
			Desc:     "节点列表",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/node/save",
			Handle:   controllers.NodeSave,
			Desc:     "节点编辑",
			RoleAuth: 1,
		},
		{
			Method:   "POST",
			Path:     "/node/delete",
			Handle:   controllers.NodeDelete,
			Desc:     "节点删除",
			RoleAuth: 1,
		},
		{
			Method:   "POST",
			Path:     "/node/health",
			Handle:   controllers.NodeHealth,
			Desc:     "健康检查",
			RoleAuth: 0,
		},

		{
			Method:   "GET",
			Path:     "/node/selector",
			Handle:   controllers.NodeSelector,
			Desc:     "节点选择",
			RoleAuth: 0,
		},

		{
			Method:   "GET",
			Path:     "/project/list",
			Handle:   controllers.ProjectList,
			Desc:     "项目列表",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/project/save",
			Handle:   controllers.ProjectSave,
			Desc:     "项目保存",
			RoleAuth: 2,
		},
		{
			Method:   "POST",
			Path:     "/project/delete",
			Handle:   controllers.ProjectDelete,
			Desc:     "项目删除",
			RoleAuth: 2,
		},

		{
			Method:   "GET",
			Path:     "/project/selector",
			Handle:   controllers.ProjectSelector,
			Desc:     "项目选择",
			RoleAuth: 0,
		},

		{
			Method:   "GET",
			Path:     "/task/list",
			Handle:   controllers.TaskList,
			Desc:     "任务列表",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/task/save",
			Handle:   controllers.TaskSave,
			Desc:     "新增任务",
			RoleAuth: 2,
		},
		{
			Method:   "POST",
			Path:     "/task/following",
			Handle:   controllers.TaskFollowing,
			Desc:     "关注任务",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/task/un/following",
			Handle:   controllers.TaskUnFollowing,
			Desc:     "取消关注任务",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/task/delete",
			Handle:   controllers.TaskDelete,
			Desc:     "删除任务",
			RoleAuth: 2,
		},
		{
			Method:   "POST",
			Path:     "/task/exec",
			Handle:   controllers.TaskExec,
			Desc:     "执行任务",
			RoleAuth: 2,
		},
		{
			Method:   "POST",
			Path:     "/task/stop",
			Handle:   controllers.TaskStop,
			Desc:     "停止任务",
			RoleAuth: 2,
		},
		{
			Method:   "GET",
			Path:     "/task/detail",
			Handle:   controllers.TaskDetail,
			Desc:     "任务详情",
			RoleAuth: 0,
		},

		{
			Method:   "GET",
			Path:     "/task/logs",
			Handle:   controllers.TaskLogs,
			Desc:     "任务日志",
			RoleAuth: 0,
		},

		{
			Method:   "GET",
			Path:     "/process/list",
			Handle:   controllers.ProcessList,
			Desc:     "进程列表",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/process/save",
			Handle:   controllers.ProcessSave,
			Desc:     "新增进程",
			RoleAuth: 2,
		},
		{
			Method:   "POST",
			Path:     "/process/following",
			Handle:   controllers.ProcessFollowing,
			Desc:     "关注进程",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/process/un/following",
			Handle:   controllers.ProcessUnFollowing,
			Desc:     "取消关注进程",
			RoleAuth: 0,
		},
		{
			Method:   "POST",
			Path:     "/process/delete",
			Handle:   controllers.ProcessDelete,
			Desc:     "删除进程",
			RoleAuth: 2,
		},
		{
			Method:   "POST",
			Path:     "/process/start",
			Handle:   controllers.ProcessStart,
			Desc:     "启动进程",
			RoleAuth: 2,
		},
		{
			Method:   "POST",
			Path:     "/process/stop",
			Handle:   controllers.ProcessStop,
			Desc:     "暂停进程",
			RoleAuth: 2,
		},
		{
			Method:   "GET",
			Path:     "/process/detail",
			Handle:   controllers.ProcessDetail,
			Desc:     "进程详情",
			RoleAuth: 0,
		},
		{
			Method:   "GET",
			Path:     "/process/logs",
			Handle:   controllers.ProcessLogs,
			Desc:     "进程日志",
			RoleAuth: 0,
		},

		{
			Method:   "GET",
			Path:     "/dashboard",
			Handle:   controllers.Dashboard,
			Desc:     "首页信息",
			RoleAuth: 0,
		},
	}

	return routes
}
