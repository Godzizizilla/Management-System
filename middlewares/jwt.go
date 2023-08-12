package routes

import (
	"context"
	"github.com/Godzizizilla/Management-System/cache"
	"github.com/Godzizizilla/Management-System/models"
	"github.com/Godzizizilla/Management-System/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func JWTMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "Authorization header not found",
		})
		c.Abort()
		return
	}

	// Split the token from the 'Bearer' keyword.
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "Authorization header not found",
		})
		c.Abort()
		return
	}

	tokenString := parts[1]

	// TODO 判断token是否在黑名单
	isBad, err := cache.RC.Get(context.TODO(), tokenString).Result()
	if isBad == "invalid" {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "Authorization header not found",
		})
		c.Abort()
		return
	}

	// 鉴权
	id, role, _, err := utils.AuthenticateToken(tokenString)

	// 鉴权失败
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
		})
		c.Abort()
		return
	}

	// 鉴权成功
	c.Set("id", id)
	c.Set("role", role)

	c.Next()

	// 修改了密码或删除了账户
	if _, exists := c.Get("changed"); exists {
		// TODO 将token加入黑名单
		cache.RC.Set(context.TODO(), tokenString, "invalid", 7*24*time.Hour).Result()
	}
}
