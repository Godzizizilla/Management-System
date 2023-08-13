package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func GenerateSecretKey(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(randomBytes)
}

func main() {
	// 使用flag包来接受命令行参数
	length := flag.Int("length", 32, "The length of the secret key in bytes")
	dir := flag.String("dir", ".", "The directory to save the .env file")
	flag.Parse()

	secretKey := GenerateSecretKey(*length)

	// 构建完整的文件路径
	envFilePath := filepath.Join(*dir, ".env")

	// 打开文件，如果文件不存在则创建它
	file, err := os.OpenFile(envFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 写入secret_key到.env文件
	_, err = file.WriteString(fmt.Sprintf("JWT_SECRET_KEY=%s\n", secretKey))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Secret key saved to %s\n", envFilePath)
}
