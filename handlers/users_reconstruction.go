package handlers

import (
	"context"
	"github.com/Godzizizilla/Management-System/db"
	"github.com/Godzizizilla/Management-System/middlewares"
	"github.com/Godzizizilla/Management-System/models"
	"github.com/Godzizizilla/Management-System/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

func Login(c *gin.Context) {
	/*
		所有可能的Response
		1. 400 为提供用户名及密码
		2. 401 用户名或密码错误
		3. 500 数据库查询错误
		4. 200 登录成功
	*/

	// 获取json request
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		sendErrorResponse(c, 400, "获取json失败")
		return
	}

	// 查询users表-用户存在
	if user, err := db.FindUserByStudentID(request.Username); err == nil {
		// 验证密码
		if !verifyPassword(c, request.Password, user.Password) {
			return
		}
		// 返回token
		token := generateToken(request.Username, "student")
		c.JSON(http.StatusOK, models.TokenResponse{
			Response: models.Response{Success: true},
			Role:     "student",
			Token:    token,
		})
		return
	}

	// 查询admins表-用户存在
	if admin, err := db.FindAdminByName(request.Username); err == nil {
		// 验证密码
		if !verifyPassword(c, request.Password, admin.Password) {
			return
		}
		// 返回token
		token := generateToken(request.Username, "admin")
		c.JSON(http.StatusOK, models.TokenResponse{
			Response: models.Response{Success: true},
			Role:     "admin",
			Token:    token,
		})
		return
	}

	// 用户不存在
	sendErrorResponse(c, 404, "用户不存在")
}

func Register(c *gin.Context) {
	/*
		所有可能的Response
		1. 400 未提供完整注册信息
		2. 400 学生ID重复
		3. 500 数据库添加失败
		4. 200 注册成功
	*/
	log := utils.NewHandleLog("POST", "Register")

	// 获取json request
	var request models.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		sendErrorResponse(c, 400, "获取json失败")
		return
	}

	// 简单验证数据是否合法
	if request.StudentID == 0 || request.Username == "" || request.Password == "" || request.Grade == "" {
		sendErrorResponse(c, 400, "请提供完整信息")
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
			sendErrorResponse(c, 400, "student id重复")
			return
		}
		log.Error("添加用户失败, 数据库错误")
		sendErrorResponse(c, 500, "添加用户失败, 服务器错误")
		return
	}

	// 注册成功
	token := generateToken(strconv.Itoa(int(request.StudentID)), "student")
	c.JSON(http.StatusOK, models.TokenResponse{
		Response: models.Response{Success: true, Message: "注册成功"},
		Role:     "student",
		Token:    token,
	})
}

func Update(c *gin.Context) {
	/*
		所有可能的Response
		1. 400 未提供完整信息, 旧密码等
		2. 400 旧密码错误
		3. 500 数据库出错导致更新失败
		4. 200 更新成功
	*/
	log := utils.NewHandleLog("PUT", "UpdateUser")

	// 获取token payload信息
	id := c.MustGet("id").(string)
	role := c.MustGet("role").(string)

	// 获取json request
	var request models.UpdateInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		sendErrorResponse(c, 400, "获取json失败")
		return
	}

	// 成功修改密码标志
	changedPasswordFlag := false

	// 判断用户类型
	if role == "student" {
		if user, err := db.FindUserByStudentID(id); err != nil {
			log.Errorf("获取学生信息失败, 不存在该学生, student_id = %s", id)
			sendErrorResponse(c, 404, "不存在该学生")
			return
		} else {
			// 修改名字
			if request.NewName != "" {
				user.Name = request.NewName
			}
			// 修改密码
			if request.NewPassword != "" && request.OldPassword != "" {
				if hashedPassword, success := verifyAndGenerateHashedPassword(c, log, request.OldPassword, user.Password, request.NewPassword); !success {
					return
				} else {
					user.Password = hashedPassword
					changedPasswordFlag = true
				}
			}
			// 更新数据库
			if err := db.UpdateUser(user); err != nil {
				changedPasswordFlag = false
				log.Error("修改信息失败, 数据库出错")
				sendErrorResponse(c, 500, "修改信息失败, 服务器错误")
				return
			}
		}
	} else if role == "admin" {
		if admin, err := db.FindAdminByName(id); err != nil {
			log.Errorf("获取管理员信息失败, 不存在该管理员, name = %s", id)
			sendErrorResponse(c, 404, "不存在该管理员")
			return
		} else {
			// 修改电话
			if request.Phone != "" {
				admin.Phone = request.Phone
			}
			// 修改邮箱
			if request.Email != "" {
				admin.Email = request.Email
			}
			// 修改密码
			if request.NewPassword != "" && request.OldPassword != "" {
				if hashedPassword, success := verifyAndGenerateHashedPassword(c, log, request.OldPassword, admin.Password, request.NewPassword); !success {
					return
				} else {
					admin.Password = hashedPassword
					changedPasswordFlag = true
				}
			}
			// 更新数据库
			if err := db.UpdateAdmin(admin); err != nil {
				changedPasswordFlag = false
				log.Error("修改信息失败, 数据库出错")
				sendErrorResponse(c, 500, "修改信息失败, 服务器错误")
				return
			}
		}
	}
	// 是否更新了密码
	if !changedPasswordFlag {
		sendSuccessResponse(c, "修改信息成功")
	} else {
		// 返回新token
		token := generateToken(id, role)
		c.JSON(http.StatusOK, models.TokenResponse{
			Response: models.Response{Success: true, Message: "修改信息成功"},
			Role:     role,
			Token:    token,
		})
		c.Set("removeJTI", true)
	}
}

func Delete(c *gin.Context) {
	log := utils.NewHandleLog("DELETE", "DeleteUser")

	// 获取路径参数
	pathParam := c.Param("id")

	// 获取token payload信息
	id := c.MustGet("id").(string)
	role := c.MustGet("role").(string)

	// 判断用户类型
	if role == "student" {
		if pathParam != "me" {
			log.Errorf("错误的路径参数: %s", pathParam)
			sendErrorResponse(c, 400, "请提供正确的路径参数")
			return
		}

		// 学生需要提供密码, 获取json字段
		var request models.DeleteRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			sendErrorResponse(c, 400, "获取json失败")
			return
		}

		// 获取学生密码, 并验证
		if user, err := db.FindUserByStudentID(id); err != nil {
			log.Errorf("获取学生信息密码失败, 不存在该学生, student_id = %s", id)
			sendErrorResponse(c, 404, "不存在该学生")
			return
		} else {
			// 验证密码
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
				sendErrorResponse(c, 400, "密码错误")
				return
			}
		}

		// 删除学生
		if err := db.DeleteUserByStudentId(id); err != nil {
			log.Errorf("删除失败, 数据库出错, student_id = %s", id)
			sendErrorResponse(c, 500, "删除学生失败")
		} else {
			sendSuccessResponse(c, "删除成功")
			c.Set("removeJTI", true)
		}
		return
	} else if role == "admin" {
		if pathParam == "me" {
			sendErrorResponse(c, 403, "不允许删除管理员账户")
			return
		}
		if err := db.DeleteUserByStudentId(id); err != nil {
			log.Errorf("删除失败, 不存在该学生, student_id = %s", id)
			sendErrorResponse(c, 404, "不存在该学生或学生ID错误")
		} else {
			sendSuccessResponse(c, "删除成功")
			c.Set("removeJTI", true)
		}
		return
	}
}

func Get(c *gin.Context) {
	log := utils.NewHandleLog("GET", "GetUser")

	// 获取路径参数
	pathParam := c.Param("id")

	// 获取token payload信息
	id := c.MustGet("id").(string)
	role := c.MustGet("role").(string)

	// 判断用户类型
	if role == "student" {
		if pathParam != "me" {
			log.Errorf("错误的路径参数: %s", pathParam)
			sendErrorResponse(c, 400, "请提供正确的路径参数")
			return
		}
		if user, err := db.FindUserByStudentID(id); err != nil {
			log.Errorf("获取学生信息失败, 不存在该学生, student_id = %s", id)
			sendErrorResponse(c, 404, "不存在该学生")
		} else {
			c.JSON(http.StatusOK, models.GetUserResponse{Response: models.Response{Success: true}, User: *user})
		}
		return
	} else if role == "admin" {
		if pathParam == "me" {
			if admin, err := db.FindAdminByName(id); err != nil {
				log.Errorf("获取管理员信息失败, 不存在该管理员, name = %s", id)
				sendErrorResponse(c, 404, "不存在该管理员")
			} else {
				c.JSON(http.StatusOK, models.GetAdminResponse{Response: models.Response{Success: true}, Admin: *admin})
			}
			return
		} else if pathParam == "all" {
			if students := db.FindAllUser(); len(*students) == 0 {
				log.Error("获取学生列表失败, 学生列表为空")
				sendErrorResponse(c, 404, "学生列表为空")
			} else {
				c.JSON(http.StatusOK, models.GetUserListResponse{Response: models.Response{Success: true}, UserList: *students})
			}
			return
		} else {
			if user, err := db.FindUserByStudentID(pathParam); err != nil {
				log.Errorf("获取学生信息失败, 不存在该学生, student_id = %s", id)
				sendErrorResponse(c, 404, "不存在该学生")
			} else {
				c.JSON(http.StatusOK, models.GetUserResponse{Response: models.Response{Success: true}, User: *user})
			}
			return
		}
	}
}

func sendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, models.Response{Success: false, Message: message})
}

func sendSuccessResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, models.Response{Success: true, Message: message})
}

func verifyAndGenerateHashedPassword(c *gin.Context, log *logrus.Entry, oldPassword string, comparativePassword string, newPassword string) (hashedPassword string, success bool) {
	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(comparativePassword), []byte(oldPassword)); err != nil {
		sendErrorResponse(c, 400, "旧密码错误, 请提供正确的旧密码")
		return "", false
	}
	// 设置新密码
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost); err != nil {
		log.Error("生成哈希密码出错")
		sendErrorResponse(c, 500, "修改密码失败, 服务器错误")
		return "", false
	} else {
		return string(hashedPassword), true
	}
}

func verifyPassword(c *gin.Context, password string, comparativePassword string) (correspond bool) {
	if err := bcrypt.CompareHashAndPassword([]byte(comparativePassword), []byte(password)); err != nil {
		// 密码错误
		sendErrorResponse(c, 401, "用户名或密码错误")
		return false
	}
	return true
}

func generateToken(id string, role string) (token string) {
	token, jti, err := utils.GenerateToken(id, role)
	if err != nil {
		logrus.Error(err.Error())
	}
	// 确保jti不重复
	for {
		if _, err := middlewares.RC.Get(context.TODO(), jti).Result(); err != nil {
			break
		}
		token, jti, _ = utils.GenerateToken(id, role)
	}
	// 添加{key: jti, val: role}到redis
	_, _ = middlewares.RC.Set(context.TODO(), jti, role, 7*24*time.Hour).Result()
	return token
}
