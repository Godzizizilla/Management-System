package handlers

import (
	"github.com/Godzizizilla/Management-System/db"
	"github.com/Godzizizilla/Management-System/models"
	"github.com/Godzizizilla/Management-System/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func UpdateAdmin(c *gin.Context) {
	// 学生无权修改
	role := c.MustGet("role")
	if role == "student" {
		c.JSON(http.StatusForbidden, models.Response{Success: false})
		return
	}

	var request models.UpdateInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Success: false})
		return
	}

	adminID := c.MustGet("id").(uint)

	admin, err := db.FindAdminByID(adminID)
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
	}

	if err := db.UpdateAdmin(admin); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Success: false})
		return
	}

	// 生成新token
	token, issuedAt, _ := utils.GenerateToken(adminID, "student")
	c.Set("changed", true)
	c.Set("issuedAt", issuedAt)
	c.JSON(http.StatusOK, models.TokenResponse{
		Response: models.Response{Success: true},
		Token:    token,
	})
}
