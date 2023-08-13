package handlers

import (
	"context"
	"github.com/Godzizizilla/Management-System/db"
	"github.com/Godzizizilla/Management-System/middlewares"
	"github.com/Godzizizilla/Management-System/models"
	"github.com/Godzizizilla/Management-System/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func UpdateAdmin(c *gin.Context) {
	// 学生无权修改
	role := c.MustGet("role")
	if role == "student" {
		return
	}

	var request models.UpdateInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Success: false})
		return
	}

	adminName := c.MustGet("id").(string)

	admin, err := db.FindAdminByName(adminName)
	// 数据库错误
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Success: false})
		return
	}

	// 修改名字
	if request.NewName != "" {
		admin.Name = request.NewName
	}

	// 修改密码
	hasChangedPassword := false
	if request.NewPassword != "" && request.OldPassword != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.OldPassword)); err != nil {
			// 密码错误
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "旧密码错误",
			})
			return
		}
		// 设置新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Success: false})
			return
		}
		admin.Password = string(hashedPassword)
		hasChangedPassword = true
	}

	// TODO 需要判断具体的错误
	if err := db.UpdateAdmin(admin); err != nil {
		hasChangedPassword = false
		c.JSON(http.StatusInternalServerError, models.Response{Success: false})
		return
	}

	// 修改了密码
	if hasChangedPassword {
		token, jti, _ := utils.GenerateToken(string(adminName), "admin")
		// 确保jti不重复
		for {
			if _, err := middlewares.RC.Get(context.TODO(), jti).Result(); err != nil {
				break
			}
			token, jti, _ = utils.GenerateToken(string(adminName), "admin")
		}

		c.JSON(http.StatusOK, models.TokenResponse{
			Response: models.Response{Success: true, Message: "修改成功"},
			Token:    token,
			Role:     "admin",
		})

		// 添加{key: jti, val: role}到redis
		middlewares.RC.Set(context.TODO(), jti, "admin", 7*24*time.Hour).Result()
		c.Set("changed", true)
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true})
}
