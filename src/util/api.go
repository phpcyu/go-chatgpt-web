package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
	strBase "qiming-server/src/struct"
)

func ApiRet(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, strBase.ApiResult{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func ApiSuc(c *gin.Context, data interface{}) {
	ApiRet(c, 20000, "", data)
}

func ApiFail(c *gin.Context, code int, msg string) {
	ApiRet(c, code, msg, nil)
}
