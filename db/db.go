package db

import (
	"errors"
	"fmt"
	"github.com/Godzizizilla/Management-System/config"
	"github.com/Godzizizilla/Management-System/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var DB *gorm.DB

func SetupDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.C.Database.Host,
		config.C.Database.User,
		config.C.Database.Password,
		config.C.Database.DBName,
		config.C.Database.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connected to the database!")
	DB = db

	db.AutoMigrate(&models.User{}, &models.Admin{})
}

func FindUserByStudentID(idStr string) (*models.User, error) {
	studentID, _ := strconv.Atoi(idStr)
	var user models.User
	if err := DB.Where("student_id = ?", studentID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindAllUser() *[]models.User {
	var users []models.User
	DB.Find(&users)
	return &users
}

func FindAdminByName(name string) (*models.Admin, error) {
	var admin models.Admin
	if err := DB.Where("name = ?", name).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

/*func FindAdminByID(id uint) (*models.Admin, error) {
	var admin models.Admin
	if err := DB.Where("id = ?", id).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}*/

func AddUser(user *models.User) error {
	if err := DB.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("student id重复")
		}
		return errors.New("添加用户失败")
	}
	return nil
}

func UpdateUser(user *models.User) error {
	if err := DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func UpdateAdmin(admin *models.Admin) error {
	if err := DB.Save(admin).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUserByStudentId(idStr string) error {
	studentID, _ := strconv.Atoi(idStr)
	result := DB.Where("student_id = ?", studentID).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("不存在该用户")
	}
	return nil
}
