package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
)

func ServeSwaggerUI(router *gin.Engine) {
	router.StaticFile("/swagger.json", "./docs/swagger.json")
	opts := middleware.SwaggerUIOpts{}
	handler := middleware.SwaggerUI(opts, nil)
	router.GET("/docs", func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	})
}
