package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go-tasks/boot/global"
	"go-tasks/master/model"
	"go-tasks/utils"
	"go-tasks/utils/response"
	"io"
	"strings"
)

func DefaultMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := strings.Trim(c.Request.URL.Path, "/")

		if strings.HasSuffix(strings.ToLower(path), ".png") ||
			strings.HasSuffix(strings.ToLower(path), ".jpg") ||
			strings.HasSuffix(strings.ToLower(path), ".jpeg") {
			c.File(global.ResourceDir + path)
			return
		}

		if path == "login" {
			c.Next()
			return
		}

		r, ok := global.RouterMap[path]
		if !ok {
			response.NotFound(c)
			c.Abort() //
			return
		}

		if strings.ToUpper(r.Method) == "POST" && strings.ToUpper(c.Request.Method) != "POST" {
			response.Fail("Allow: POST", c)
			c.Abort()
			return
		}

		if strings.ToUpper(r.Method) == "GET" && strings.ToUpper(c.Request.Method) != "GET" {
			response.Fail("Allow: GET", c)
			c.Abort() //
			return
		}

		token := c.GetHeader("token")
		adminUser, loginOk := checkLogin(token, c)
		if !loginOk {
			response.Unauthorized(c)
			c.Abort() //
			return
		}

		c.Set("adminId", adminUser.Id)
		if r.RoleAuth != 0 {
			if adminUser.RoleId > r.RoleAuth {
				response.Forbidden(c)
				c.Abort()
				return
			}
		}
		adminId := adminUser.Id

		var requestBodyPointer *string
		if strings.ToUpper(r.Method) == "POST" {
			bodyByte, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyByte))
			requestBodyStr := string(bodyByte)
			if requestBodyStr != "" {
				requestBodyPointer = &requestBodyStr
			}
		}

		c.Next()

		if strings.ToUpper(r.Method) == "POST" {
			adminUserLogsModel := model.NewAdminUserLogsModel()
			adminUserLogsModel.NewLog(adminId, requestBodyPointer, path, c)
		}

	}
}

func checkLogin(token string, c *gin.Context) (*model.AdminUserModel, bool) {
	if token == "" {
		return nil, false
	}
	userToken, err := utils.ParseToken(token)

	if err != nil {
		return nil, false
	}
	adminUserModel := model.NewAdminUserModel()
	user, err := adminUserModel.FindById(userToken.UserId, c)
	if err != nil {
		return nil, false
	}

	if user.LoginTime != userToken.LoginTime {
		return nil, false
	}

	return user, true
}
