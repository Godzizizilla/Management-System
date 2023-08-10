package routes

import (
	"github.com/Godzizizilla/Management-System/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouters(router *gin.Engine) {
	cfg := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTION"},
		AllowHeaders:     []string{"Authorization"},
		AllowCredentials: true,
	}
	router.Use(cors.New(cfg))

	apiPublic := router.Group("/v1")
	apiProtected := router.Group("/v1")
	apiProtected.Use(JWTMiddleware)

	// public api
	apiPublic.POST("/login", handlers.Login)
	apiPublic.POST("/users", handlers.CreateUser)

	// users api
	apiProtected.PUT("/users", handlers.UpdateUser)
	apiProtected.DELETE("/users", handlers.DeleteUser)
	apiProtected.GET("/users", handlers.GetUser)
	apiProtected.DELETE("/users/:id", handlers.DeleteUser)
	apiProtected.GET("/users/:id", handlers.GetUser)

	// admin api
	apiProtected.PUT("/admin", handlers.UpdateAdmin)
}
