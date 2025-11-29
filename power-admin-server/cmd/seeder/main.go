package main

import (
	"flag"
	"fmt"
	"log"

	"power-admin-server/internal/config"
	"power-admin-server/pkg/db"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/conf"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	// 删除旧表并重新创建（带正确的中文注释）
	if err := recreateTables(database); err != nil {
		log.Fatalf("Failed to recreate tables: %v", err)
	}

	// 插入数据
	if err := seedDatabase(database); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("✓ 所有数据初始化成功！")
}

func recreateTables(database *gorm.DB) error {
	log.Println("正在重新创建表...")

	// 重新创建所有表（GORM会根据模型的comment标签自动创建中文注释）
	// 也会自动处理多对多的关联表
	if err := database.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Menu{},
		&models.Dictionary{},
		&models.API{},
		&models.Plugin{},
		&models.Log{},
		&models.App{},
		&models.Review{},
	); err != nil {
		return fmt.Errorf("failed to recreate tables: %w", err)
	}

	log.Println("✓ 表创建成功（含中文注释）")
	return nil
}

func seedDatabase(database *gorm.DB) error {
	// 生成密码哈希 (密码: 123456)
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	passwordHashStr := string(passwordHash)

	// 插入用户
	users := []models.User{
		{
			Username: "admin",
			Password: passwordHashStr,
			Nickname: "管理员",
			Email:    "admin@example.com",
			Phone:    "13800000000",
			Status:   1,
		},
		{
			Username: "editor",
			Password: passwordHashStr,
			Nickname: "编辑",
			Email:    "editor@example.com",
			Phone:    "13800000001",
			Status:   1,
		},
		{
			Username: "user",
			Password: passwordHashStr,
			Nickname: "普通用户",
			Email:    "user@example.com",
			Phone:    "13800000002",
			Status:   1,
		},
	}

	for _, user := range users {
		// 先查找是否存在
		var existingUser models.User
		if err := database.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
			// 存在则更新密码
			database.Model(&existingUser).Update("password", user.Password)
		} else {
			// 不存在则创建
			if err := database.Create(&user).Error; err != nil {
				return fmt.Errorf("failed to seed users: %w", err)
			}
		}
	}
	log.Println("✓ 用户数据插入成功 (默认密码: 123456)")

	// 插入角色
	roles := []models.Role{
		{Name: "admin", Description: "管理员", Status: 1},
		{Name: "editor", Description: "编辑", Status: 1},
		{Name: "user", Description: "普通用户", Status: 1},
	}

	for _, role := range roles {
		if err := database.FirstOrCreate(&role, models.Role{Name: role.Name}).Error; err != nil {
			return fmt.Errorf("failed to seed roles: %w", err)
		}
	}
	log.Println("✓ 角色数据插入成功")

	// 插入菜单
	menus := []models.Menu{
		{ParentID: 0, MenuName: "系统管理", MenuPath: "/system", Icon: "setting", Sort: 10, Status: 1, MenuType: 1},
		{ParentID: 1, MenuName: "用户管理", MenuPath: "/system/users", Component: "system/user/UserList", Icon: "user", Sort: 1, Status: 1, MenuType: 1},
		{ParentID: 1, MenuName: "角色管理", MenuPath: "/system/roles", Component: "system/role/RoleList", Icon: "admin", Sort: 2, Status: 1, MenuType: 1},
		{ParentID: 1, MenuName: "菜单管理", MenuPath: "/system/menus", Component: "system/menu/MenuList", Icon: "menu", Sort: 3, Status: 1, MenuType: 1},
		{ParentID: 1, MenuName: "权限管理", MenuPath: "/system/permissions", Component: "system/permission/PermissionList", Icon: "lock", Sort: 4, Status: 1, MenuType: 1},
		{ParentID: 1, MenuName: "API管理", MenuPath: "/system/apis", Component: "system/api/ApiList", Icon: "link", Sort: 5, Status: 1, MenuType: 1},
		{ParentID: 0, MenuName: "内容管理", MenuPath: "/content", Icon: "document", Sort: 20, Status: 1, MenuType: 1},
		{ParentID: 7, MenuName: "字典管理", MenuPath: "/content/dicts", Component: "content/dict/DictList", Icon: "list", Sort: 1, Status: 1, MenuType: 1},
		{ParentID: 0, MenuName: "应用中心", MenuPath: "/market", Icon: "shopping", Sort: 30, Status: 1, MenuType: 1},
		{ParentID: 9, MenuName: "应用市场", MenuPath: "/market/apps", Component: "market/AppMarket", Icon: "shop", Sort: 1, Status: 1, MenuType: 1},
		{ParentID: 0, MenuName: "系统设置", MenuPath: "/system-config", Icon: "setting", Sort: 40, Status: 1, MenuType: 1},
		{ParentID: 11, MenuName: "日志管理", MenuPath: "/logs", Component: "logs/LogList", Icon: "monitor", Sort: 1, Status: 1, MenuType: 1},
	}

	for _, menu := range menus {
		if err := database.FirstOrCreate(&menu, models.Menu{MenuPath: menu.MenuPath}).Error; err != nil {
			return fmt.Errorf("failed to seed menus: %w", err)
		}
	}
	log.Println("✓ 菜单数据插入成功")

	// 插入字典
	dictionaries := []models.Dictionary{
		{DictType: "gender", DictLabel: "男", DictValue: "1", Sort: 1, Status: 1},
		{DictType: "gender", DictLabel: "女", DictValue: "2", Sort: 2, Status: 1},
		{DictType: "gender", DictLabel: "未知", DictValue: "0", Sort: 3, Status: 1},
		{DictType: "status", DictLabel: "启用", DictValue: "1", Sort: 1, Status: 1},
		{DictType: "status", DictLabel: "禁用", DictValue: "0", Sort: 2, Status: 1},
		{DictType: "menu_type", DictLabel: "菜单", DictValue: "1", Sort: 1, Status: 1},
		{DictType: "menu_type", DictLabel: "按钮", DictValue: "2", Sort: 2, Status: 1},
		{DictType: "user_status", DictLabel: "正常", DictValue: "1", Sort: 1, Status: 1},
		{DictType: "user_status", DictLabel: "禁用", DictValue: "0", Sort: 2, Status: 1},
		{DictType: "user_status", DictLabel: "锁定", DictValue: "2", Sort: 3, Status: 1},
	}

	for _, dict := range dictionaries {
		if err := database.FirstOrCreate(&dict, models.Dictionary{DictType: dict.DictType, DictValue: dict.DictValue}).Error; err != nil {
			return fmt.Errorf("failed to seed dictionaries: %w", err)
		}
	}
	log.Println("✓ 字典数据插入成功")

	return nil
}
