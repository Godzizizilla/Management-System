package main

import (
	"github.com/Godzizizilla/Management-System/cache"
	"github.com/Godzizizilla/Management-System/config"
	_ "github.com/Godzizizilla/Management-System/docs"
	"github.com/Godzizizilla/Management-System/routes"
)

func main() {
	config.Init()
	defer cache.RC.Close()

	r := routes.InitRouters() // 初始化并获取路由
	r.Run(":7890")
}
