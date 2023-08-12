package routes

import (
	"github.com/Godzizizilla/Management-System/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	middlewares.ServeCors(r)
	middlewares.ServeSwaggerUI(r)
	SetupRouters(r)
	return r
}
