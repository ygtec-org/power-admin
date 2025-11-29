package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 数据库中的哈希值
	hash := "$2a$10$9L89bPEx.1S4DBsv0blEgu9rK3MmSWmqtd/LbOWBxPi3iuXG3UwxW"

	// 测试常见密码
	passwords := []string{
		"password",
		"123456",
		"admin123",
		"admin",
		"12345678",
		"111111",
		"000000",
	}

	fmt.Println("测试密码哈希...")
	for _, pwd := range passwords {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
		if err == nil {
			fmt.Printf("✓ 找到匹配的密码: %s\n", pwd)
			return
		}
	}

	fmt.Println("✗ 未找到匹配的密码")

	// 生成新的密码哈希
	fmt.Println("\n生成新密码哈希...")
	newPassword := "123456"
	newHash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	fmt.Printf("密码 '%s' 的哈希: %s\n", newPassword, string(newHash))
}
