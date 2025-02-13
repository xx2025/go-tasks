package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func OkWithData(msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		200,
		data,
		msg,
	})
}

func Ok(msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		200,
		nil,
		msg,
	})
}

func Fail(msg string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		400,
		nil,
		msg,
	})
}

func NotFound(c *gin.Context) {

	c.JSON(http.StatusBadRequest, Response{
		http.StatusNotFound,
		nil,
		"Not Found",
	})
}

func Unauthorized(c *gin.Context) {

	c.JSON(http.StatusBadRequest, Response{
		http.StatusUnauthorized,
		nil,
		"请登录",
	})
}

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		http.StatusForbidden,
		nil,
		"无权操作",
	})
}
