package db

import (
	"errors"
	"fmt"
	"github.com/Godzizizilla/Management-System/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

var DB *gorm.DB

func SetupDB() {
	dsn := "host=localhost user=postgres password=postgres2023 dbname=management_system port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connected to the database!")
	DB = db

	db.AutoMigrate(&models.User{}, &models.Admin{})
}

func FindUserByStudentID(studentID uint) (*models.User, error) {
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

func FindAdminByID(id uint) (*models.Admin, error) {
	var admin models.Admin
	if err := DB.Where("id = ?", id).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

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

func DeleteUserByStudentId(studentId uint) error {
	fmt.Println(studentId)
	if err := DB.Unscoped().Where("student_id = ?", studentId).Delete(&models.User{}).Error; err != nil {
		return err
	}
	return nil
}
