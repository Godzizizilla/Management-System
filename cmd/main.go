package main

import (
	"github.com/Godzizizilla/Management-System/cache"
	"github.com/Godzizizilla/Management-System/config"
	_ "github.com/Godzizizilla/Management-System/docs"
	"github.com/Godzizizilla/Management-System/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	defer cache.RC.Close()

	r := gin.Default()
	routes.SetupRouters(r)
	r.Run(":7890")
}
