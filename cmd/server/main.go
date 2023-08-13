package main

import (
	"fmt"
	"github.com/Godzizizilla/Management-System/config"
	"github.com/Godzizizilla/Management-System/db"
	_ "github.com/Godzizizilla/Management-System/docs"
	"github.com/Godzizizilla/Management-System/middlewares"
	"github.com/Godzizizilla/Management-System/routes"
	"github.com/Godzizizilla/Management-System/utils"
)

func main() {
	config.Load() // 加载./config/config.yml文件

	db.SetupDB()   // 连接数据库
	db.InitAdmin() // 初始化管理员账户, 需要在连接数据库之后

	utils.InitLog() // 初始化log配置

	r := routes.InitRouters()      // 初始化并获取路由
	defer middlewares.CloseRedis() // 使用了Redis实现的中间件, 确保main函数结束时关闭Redis客户端

	r.Run(fmt.Sprintf(":%d", config.C.Server.Port))
}
