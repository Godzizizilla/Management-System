package handlers

/*
import (
	"context"
	"errors"
	"github.com/Godzizizilla/Management-System/db"
	"github.com/Godzizizilla/Management-System/middlewares"
	"github.com/Godzizizilla/Management-System/models"
	"github.com/Godzizizilla/Management-System/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func Login2(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "请提供用户名及密码"})
		return
	}

	var role string
	var comparativePassword string
	var err error
	var user *models.User

	// 首先查询users表
	studentID, _ := strconv.Atoi(request.Username)
	if studentID == 0 {
		err = gorm.ErrRecordNotFound
	} else {
		user, err = db.FindUserByStudentID(uint(studentID))
	}

	// 存在该用户
	if err == nil {
		role = "student"
		comparativePassword = user.Password
	} else {
		// 数据库错误
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, models.Response{Success: false})
			return
		}

		// 不存在该用户, 再查询admins表
		admin, err := db.FindAdminByName(request.Username)
		// 存在该管理员
		if err == nil {
			role = "admin"
			comparativePassword = admin.Password
		} else {
			// 数据库错误
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusInternalServerError, models.Response{Success: false})
				return
			}
			// 不存在该管理员
			c.JSON(http.StatusUnauthorized, models.Response{Success: false, Message: "用户名或密码错误"})
			return
		}
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(comparativePassword), []byte(request.Password)); err != nil {
		// 密码错误
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "用户名或密码错误",
		})
		return
	}
	// 密码正确, 生成token
	token, jti, _ := utils.GenerateToken(request.Username, role)
	// 确保jti不重复
	for {
		// 没有记录, 返回err, 说明不重复
		if _, err := middlewares.RC.Get(context.TODO(), jti).Result(); err != nil {
			break
		}
		token, jti, _ = utils.GenerateToken(request.Username, role)
	}

	c.JSON(http.StatusOK, models.TokenResponse{
		Response: models.Response{Success: true, Message: "登录成功"},
		Token:    token,
		Role:     role,
	})

	// 添加{key: jti, val: role}到redis
	middlewares.RC.Set(context.TODO(), jti, role, 7*24*time.Hour).Result()
}

func CreateUser(c *gin.Context) {
	var request models.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Success: false})
		return
	}

	// 简单验证数据是否合法
	if request.StudentID == 0 || request.Username == "" || request.Password == "" || request.Grade == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
		})
		return
	}

	var user = models.User{
		Name:      request.Username,
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

	// 注册成功, 生成token
	token, jti, _ := utils.GenerateToken(request.Username, "student")
	// 确保jti不重复
	for {
		if _, err := middlewares.RC.Get(context.TODO(), jti).Result(); err != nil {
			break
		}
		token, jti, _ = utils.GenerateToken(request.Username, "student")
	}

	c.JSON(http.StatusOK, models.TokenResponse{
		Response: models.Response{Success: true, Message: "注册成功"},
		Token:    token,
		Role:     "student",
	})

	// 添加{key: jti, val: role}到redis
	middlewares.RC.Set(context.TODO(), jti, "student", 7*24*time.Hour).Result()
}

func UpdateUser(c *gin.Context) {
	// 修改学生信息, 管理员跳到另一个handle
	role := c.MustGet("role")
	if role == "admin" {
		return
	}

	var info models.UpdateInfoRequest
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Success: false})
		return
	}
	id, _ := strconv.Atoi(c.MustGet("id").(string))
	studentID := uint(id)

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

	// TODO 需要判断具体的错误
	if err := db.UpdateUser(user); err != nil {
		hasChangedPassword = false
		c.JSON(http.StatusInternalServerError, models.Response{Success: false, Message: "更新失败"})
		return
	}

	// 修改了密码
	if hasChangedPassword {
		token, jti, _ := utils.GenerateToken(string(studentID), "student")
		// 确保jti不重复
		for {
			if _, err := middlewares.RC.Get(context.TODO(), jti).Result(); err != nil {
				break
			}
			token, jti, _ = utils.GenerateToken(string(studentID), "student")
		}

		c.JSON(http.StatusOK, models.TokenResponse{
			Response: models.Response{Success: true, Message: "更新成功"},
			Token:    token,
			Role:     "student",
		})

		// 添加{key: jti, val: role}到redis
		middlewares.RC.Set(context.TODO(), jti, "student", 7*24*time.Hour).Result()
		c.Set("changed", true)
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "更新成功"})
}

func DeleteUser(c *gin.Context) {
	var studentID uint

	role := c.MustGet("role")
	if role == "student" {
		id, _ := strconv.Atoi(c.MustGet("id").(string))
		studentID := uint(id)

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
		if c.Param("id") == "me" {
			c.JSON(http.StatusForbidden, models.Response{
				Success: false,
				Message: "不允许删除管理员账户",
			})
			return
		}
		idParam, _ := strconv.Atoi(c.Param("id"))
		studentID = uint(idParam)
	}

	if err := db.DeleteUserByStudentId(studentID); err != nil {
		c.JSON(http.StatusNotFound, models.Response{Success: false, Message: "用户不存在"})
		return
	}

	// 管理员删除账户不会拉黑token
	if role != "admin" {
		c.Set("changed", true)
	}

	c.JSON(http.StatusOK, models.Response{Success: true})
}

func GetUser(c *gin.Context) {
	var studentID uint
	role := c.MustGet("role")
	if role == "student" {
		if c.Param("id") != "me" {
			c.JSON(http.StatusNotFound, models.Response{Success: false, Message: "请指定查询对象, me 或者 学生ID"})
			return
		}
		id, _ := strconv.Atoi(c.MustGet("id").(string))
		studentID = uint(id)
		user, err := db.FindUserByStudentID(studentID)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Success: false, Message: "不存在该学生"})
			return
		}

		c.JSON(http.StatusOK, models.GetUserResponse{
			Response: models.Response{Success: true},
			User:     *user,
		})
		return
	} else if role == "admin" {
		if c.Param("id") == "me" {
			// 查询管理员个人信息
			admin, err := db.FindAdminByName(c.MustGet("id").(string))
			if err != nil {
				c.JSON(http.StatusNotFound, models.Response{Success: false, Message: "不存在该管理员"})
				return
			}
			c.JSON(http.StatusOK, models.GetAdminResponse{
				Response: models.Response{Success: true},
				Role:     "admin",
				Admin:    *admin,
			})
			return
		} else if c.Param("id") == "all" {
			// 查询用户列表
			users := db.FindAllUser()
			if len(*users) == 0 {
				// 未找到学生
				c.JSON(http.StatusOK, models.GetUserListResponse{
					Response: models.Response{Success: true, Message: "学生名单为空"},
					UserList: nil,
				})
				return

			}
			c.JSON(http.StatusOK, models.GetUserListResponse{
				Response: models.Response{Success: true},
				UserList: *users,
			})
			return
		} else {
			// 查询指定学生
			idParam, _ := strconv.Atoi(c.Param("id"))
			studentID = uint(idParam)
			user, err := db.FindUserByStudentID(studentID)
			if err != nil {
				c.JSON(http.StatusNotFound, models.Response{Success: false, Message: "不存在该学生"})
				return
			}

			c.JSON(http.StatusOK, models.GetUserResponse{
				Response: models.Response{Success: true},
				User:     *user,
			})
			return
		}
	}
}
*/
