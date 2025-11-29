package main

import (
	"fmt"
	"log"
	"power-admin-server/pkg/models"

	"gorm.io/gorm"
)

// SeedDatabase 初始化数据库数据
func SeedDatabase(db *gorm.DB) error {
	// 清空现有数据（仅在开发环境使用）
	// 生产环境不建议使用

	// 插入用户
	users := []models.User{
		{
			Username: "admin",
			Password: "$2a$10$9L89bPEx.1S4DBsv0blEgu9rK3MmSWmqtd/LbOWBxPi3iuXG3UwxW",
			Nickname: "管理员",
			Email:    "admin@example.com",
			Phone:    "13800000000",
			Status:   1,
		},
		{
			Username: "editor",
			Password: "$2a$10$9L89bPEx.1S4DBsv0blEgu9rK3MmSWmqtd/LbOWBxPi3iuXG3UwxW",
			Nickname: "编辑",
			Email:    "editor@example.com",
			Phone:    "13800000001",
			Status:   1,
		},
		{
			Username: "user",
			Password: "$2a$10$9L89bPEx.1S4DBsv0blEgu9rK3MmSWmqtd/LbOWBxPi3iuXG3UwxW",
			Nickname: "普通用户",
			Email:    "user@example.com",
			Phone:    "13800000002",
			Status:   1,
		},
	}

	for _, user := range users {
		if err := db.FirstOrCreate(&user, models.User{Username: user.Username}).Error; err != nil {
			return fmt.Errorf("failed to seed users: %w", err)
		}
	}
	log.Println("✓ 用户数据插入成功")

	// 插入角色
	roles := []models.Role{
		{
			Name:        "admin",
			Description: "管理员",
			Status:      1,
		},
		{
			Name:        "editor",
			Description: "编辑",
			Status:      1,
		},
		{
			Name:        "user",
			Description: "普通用户",
			Status:      1,
		},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, models.Role{Name: role.Name}).Error; err != nil {
			return fmt.Errorf("failed to seed roles: %w", err)
		}
	}
	log.Println("✓ 角色数据插入成功")

	// 关联用户到角色
	userRoles := []map[string]int64{
		{"user_id": 1, "role_id": 1},
		{"user_id": 2, "role_id": 2},
		{"user_id": 3, "role_id": 3},
	}

	for _, ur := range userRoles {
		if err := db.Create(map[string]interface{}{
			"user_id": ur["user_id"],
			"role_id": ur["role_id"],
		}).Error; err != nil && !isUniqueConstraintError(err) {
			return fmt.Errorf("failed to seed user_roles: %w", err)
		}
	}
	log.Println("✓ 用户角色关联插入成功")

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
		if err := db.FirstOrCreate(&menu, models.Menu{MenuPath: menu.MenuPath}).Error; err != nil {
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
		if err := db.FirstOrCreate(&dict, models.Dictionary{DictType: dict.DictType, DictValue: dict.DictValue}).Error; err != nil {
			return fmt.Errorf("failed to seed dictionaries: %w", err)
		}
	}
	log.Println("✓ 字典数据插入成功")

	return nil
}

// isUniqueConstraintError 检查是否是唯一约束错误
func isUniqueConstraintError(err error) bool {
	return err != nil && (err.Error() == "UNIQUE constraint failed" ||
		err.Error() == "Error 1062: Duplicate entry" ||
		err.Error() == "duplicate key value violates unique constraint")
}
