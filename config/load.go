package config

import (
	"github.com/Godzizizilla/Management-System/utils"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	Swagger struct {
		FilePath string `yaml:"filePath"`
		OpenUI   bool   `yaml:"openUI"`
	} `yaml:"swagger"`
}

var C Config

func Load() {
	// 加载配置文件
	f, err := os.Open("./config/config.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&C)
	if err != nil {
		panic(err)
	}

	// 检测Server字段
	if C.Server.Host == "" {
		panic("Server.Host is not defined in config.yml")
	}
	if C.Server.Port == 0 {
		panic("Server.Port is not defined in config.yml")
	}

	// 检查Database字段
	if C.Database.Host == "" {
		panic("Database.Host is not defined in config.yml")
	}
	if C.Database.Port == 0 {
		panic("Database.Port is not defined in config.yml")
	}
	if C.Database.User == "" {
		panic("Database.User is not defined in config.yml")
	}
	if C.Database.Password == "" {
		panic("Database.Password is not defined in config.yml")
	}
	if C.Database.DBName == "" {
		panic("Database.DBName is not defined in config.yml")
	}

	// 检查Redis字段
	if C.Redis.Address == "" {
		panic("Redis.Address is not defined in config.yml")
	}

	// 检查swagger字段
	if C.Swagger.FilePath == "" {
		panic("Swagger.FilePath is not defined in config.yml")
	}

	// 加载环境变量
	err = godotenv.Load()
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if len(secretKey) == 0 {
		panic("获取JWT_SECRET_KEY失败\n请将JWT_SECRET_KEY添加到环境变量, 或运行该命令生成:\ngo run cmd/generate_secret_key/main.go -dir .")
	}
	utils.SecretKey = []byte(secretKey)
}
