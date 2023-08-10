package config

import (
	"errors"
	"fmt"
	"github.com/Godzizizilla/Management-System/db"
	"github.com/Godzizizilla/Management-System/models"
	"golang.org/x/crypto/ssh/terminal"
	"gorm.io/gorm"
	"os"
)

func InitAdmin() {
	var admin models.Admin
	if err := db.DB.First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有找到任何管理员记录
			fmt.Println("[Init]: 未设置管理员账户!")
		} else {
			// 其他数据库错误
			fmt.Println("[Init]: 查询数据库错误: ", err)
		}
	} else {
		// 存在管理员记录
		fmt.Println("[Init]: 管理员账户存在, 账户名为: ", admin.Name)
		return
	}

	// 获取用户名
	fmt.Print("[Init]: 请设置管理员账户\n请输入用户名: ")
	fmt.Scan(&admin.Name)

	// 获取密码
	fmt.Print("请输入密码: ")
	bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("\n读取密码失败:", err)
		return
	}
	admin.Password = string(bytePassword)

	if err := db.DB.Create(&admin).Error; err != nil {
		fmt.Println("\n[Init]: 设置管理员账户失败")
	} else {
		fmt.Println("\n[Init]: 设置管理员账户成功")
	}
}
