package middlewares

import (
	"context"
	"github.com/Godzizizilla/Management-System/models"
	"github.com/Godzizizilla/Management-System/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JWTMiddleware(c *gin.Context) {
	log := utils.NewFuncLog("JWTMiddleware")

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Error("无Authorization标头")
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
		log.Errorf("无法获取token: \nlen(parts): %d\nparts[0]: %s\nparts[1]: %s", len(parts), parts[0], parts[1])
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "无法获取token",
		})
		c.Abort()
		return
	}

	tokenString := parts[1]

	// 鉴权
	id, roleToken, jti, err := utils.AuthenticateToken(tokenString)

	// 鉴权失败
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "token错误",
		})
		c.Abort()
		return
	}

	// 判断jti是否在白名单内, 以及role是否正确
	roleRedis, err := RC.Get(context.TODO(), jti).Result()
	if err != nil || roleToken != roleRedis {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "token错误",
		})
		c.Abort()
		return
	}

	c.Set("id", id)
	c.Set("role", roleRedis)

	c.Next()

	// 修改了密码或删除了账户
	if _, exists := c.Get("removeJTI"); exists {
		// 从白名单中删除
		RC.Del(context.TODO(), jti).Result()
	}
}
