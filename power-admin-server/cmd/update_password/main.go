package main

import (
	"flag"
	"log"

	"power-admin-server/internal/config"
	"power-admin-server/pkg/db"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/conf"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	configFile := flag.String("f", "etc/power-api.yaml", "the config file")
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化数据库
	database, err := db.InitDB(c.Mysql.DataSource)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// 生成新密码的哈希
	newPassword := "123456"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to generate password hash: %v", err)
	}

	// 更新所有用户的密码
	result := database.Model(&models.User{}).
		Where("username IN ?", []string{"admin", "editor", "user"}).
		Update("password", string(passwordHash))

	if result.Error != nil {
		log.Fatalf("Failed to update passwords: %v", result.Error)
	}

	log.Printf("✓ 成功更新 %d 个用户的密码为: %s\n", result.RowsAffected, newPassword)
	log.Println("用户列表:")
	log.Println("  - 手机号: 13800000000, 密码: 123456 (管理员)")
	log.Println("  - 手机号: 13800000001, 密码: 123456 (编辑)")
	log.Println("  - 手机号: 13800000002, 密码: 123456 (普通用户)")
}
