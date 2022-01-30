package home

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Hello(context *gin.Context) {
	context.String(http.StatusOK, "hello")
}

func Ping(context *gin.Context) {
	context.String(http.StatusOK, "Pong")
}

func PanicTest(context *gin.Context) {
	panic("something")
}
