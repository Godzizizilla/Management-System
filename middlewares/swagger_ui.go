package middlewares

import (
	"fmt"
	"github.com/Godzizizilla/Management-System/config"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"os/exec"
	"runtime"
)

func ServeSwaggerUI(router *gin.Engine) {
	router.StaticFile("/swagger.json", config.C.Swagger.FilePath)
	opts := middleware.SwaggerUIOpts{}
	handler := middleware.SwaggerUI(opts, nil)
	router.GET("/docs", func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	})

	if config.C.Swagger.OpenUI == true {
		url := fmt.Sprintf("http://%s:%d/docs", config.C.Server.Host, config.C.Server.Port)
		fmt.Println("Swagger UI 地址: ", url)
		openURL(url)
	}
}

func openURL(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		// Unsupported OS, you can add more switch cases if needed
		return
	}

	cmd.Run()
}
