package handlers

import (
	"errors"
	"github.com/Godzizizilla/Management-System/db"
	"github.com/Godzizizilla/Management-System/models"
	"github.com/Godzizizilla/Management-System/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func Login(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Success: false})
		return
	}

	var id uint
	var comparativePassword string
	var err error

	if request.Role == "student" {
		studentID, _ := strconv.Atoi(request.Username)
		var user *models.User
		user, err = db.FindUserByStudentID(uint(studentID))
		id = user.StudentID
		comparativePassword = user.Password
	} else if request.Role == "admin" {
		var admin *models.Admin
		admin, err = db.FindAdminByName(request.Username)
		id = admin.ID
		comparativePassword = admin.Password
	} else {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "",
		})
		return
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在
			c.JSON(http.StatusNotFound, models.Response{Success: false})
			return
		} else {
			// 数据库错误
			c.JSON(http.StatusInternalServerError, models.Response{Success: false})
			return
		}
	}
	// 用户存在
	if err := bcrypt.CompareHashAndPassword([]byte(comparativePassword), []byte(request.Password)); err != nil {
		// 密码错误
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "密码错误",
		})
		return
	}
	// 登录成功
	token, _, _ := utils.GenerateToken(id, request.Role)
	c.JSON(http.StatusOK, models.TokenResponse{
		Response: models.Response{Success: true},
		Token:    token,
	})
}

func CreateUser(c *gin.Context) {
	var request models.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Success: false})
		return
	}

	// 简单验证数据是否合法
	if request.StudentID == 0 || request.UserName == "" || request.Password == "" || request.Grade == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
		})
		return
	}

	var user = models.User{
		Name:      request.UserName,
		StudentID: request.StudentID,
		Grade:     request.Grade,
		Password:  request.Password,
	}

	if err := db.AddUser(&user); err != nil {
		if err.Error() == "student id重复" {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "student id重复",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "注册失败",
		})
		return
	}

	// 注册成功
	token, _, _ := utils.GenerateToken(request.StudentID, "student")
	c.JSON(http.StatusOK, models.TokenResponse{
		Response: models.Response{Success: true},
		Token:    token,
	})
}

func UpdateUser(c *gin.Context) {
	// 管理员无权修改
	role := c.MustGet("role")
	if role == "admin" {
		c.JSON(http.StatusForbidden, models.Response{Success: false})
		return
	}

	var info models.UpdateInfoRequest
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Success: false})
		return
	}
	studentID := c.MustGet("id").(uint)

	user, err := db.FindUserByStudentID(studentID)
	// 数据库错误
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Success: false})
		return
	}

	// 修改名字
	if info.NewName != "" {
		user.Name = info.NewName
	}

	// 修改密码
	hasChangedPassword := false
	if info.NewPassword != "" && info.OldPassword != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.OldPassword)); err != nil {
			// 密码错误
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "旧密码错误",
			})
			return
		}
		// 设置新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(info.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Success: false})
			return
		}
		user.Password = string(hashedPassword)
		hasChangedPassword = true
	}

	if err := db.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Success: false})
		return
	}
	if hasChangedPassword {
		// 生成新token
		token, issuedAt, _ := utils.GenerateToken(studentID, "student")
		c.Set("changed", true)
		c.Set("issuedAt", issuedAt)
		c.JSON(http.StatusOK, models.TokenResponse{
			Response: models.Response{Success: true},
			Token:    token,
		})
	}
	c.JSON(http.StatusOK, models.Response{Success: true})
}

func DeleteUser(c *gin.Context) {
	var studentID uint

	role := c.MustGet("role")
	if role == "student" {
		studentID = c.MustGet("id").(uint)

		// 需要输入密码删除
		var request models.DeleteRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Success: false})
			return
		}
		if request.Password == "" {
			c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "请输入密码"})
			return
		}

		// 验证密码是否正确
		user, err := db.FindUserByStudentID(studentID)
		// 数据库错误
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Success: false})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			// 密码错误
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "旧密码错误",
			})
			return
		}

	} else if role == "admin" {
		idParam, _ := strconv.Atoi(c.Param("id"))
		studentID = uint(idParam)
	}

	if err := db.DeleteUserByStudentId(studentID); err != nil {
		c.JSON(http.StatusNotFound, models.Response{Success: false, Message: "用户不存在"})
		return
	}

	c.Set("changed", true)
	c.JSON(http.StatusOK, models.Response{Success: true})
}

func GetUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Header("Access-Control-Expose-Headers", "Content-Length")
	c.Header("Access-Control-Allow-Credentials", "true")

	var studentID uint
	role := c.MustGet("role")
	if role == "student" {
		studentID = c.MustGet("id").(uint)
	} else if role == "admin" {
		idParam, _ := strconv.Atoi(c.Param("id"))
		studentID = uint(idParam)
		// 查询用户列表
		if idParam == 0 {
			users := db.FindAllUser()
			c.JSON(http.StatusOK, models.GetUserListResponse{
				Response: models.Response{Success: true},
				UserList: *users,
			})
			return
		}
	}

	// 查询单个用户
	user, err := db.FindUserByStudentID(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Success: false})
		return
	}

	c.JSON(http.StatusOK, models.GetUserResponse{
		Response: models.Response{Success: true},
		User:     *user,
	})
}
