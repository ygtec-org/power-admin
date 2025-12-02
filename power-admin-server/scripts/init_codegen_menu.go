package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Menu struct {
	ID        int64  `gorm:"primaryKey"`
	ParentID  int64  `gorm:"column:parent_id"`
	MenuName  string `gorm:"column:menu_name"`
	MenuPath  string `gorm:"column:menu_path"`
	Component string `gorm:"column:component"`
	Icon      string `gorm:"column:icon"`
	Sort      int    `gorm:"column:sort"`
	Status    int    `gorm:"column:status"`
	MenuType  int    `gorm:"column:menu_type"`
	Remark    string `gorm:"column:remark"`
}

func (Menu) TableName() string {
	return "admin_menus"
}

func main() {
	// 连接数据库
	dsn := "root:root@tcp(127.0.0.1:3306)/power_admin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 先删除旧的菜单（如果存在）
	db.Where("menu_name IN (?, ?, ?)", "开发工具", "代码生成", "生成历史").Delete(&Menu{})

	// 插入顶级菜单
	devtoolsMenu := &Menu{
		ParentID:  0,
		MenuName:  "开发工具",
		MenuPath:  "/devtools",
		Component: "",
		Icon:      "Tool",
		Sort:      20,
		Status:    1,
		MenuType:  1,
		Remark:    "开发工具菜单",
	}

	if err := db.Create(devtoolsMenu).Error; err != nil {
		log.Fatal("插入开发工具菜单失败:", err)
	}
	fmt.Printf("✅ 插入开发工具菜单成功, ID: %d\n", devtoolsMenu.ID)

	// 插入子菜单 - 代码生成
	codegenMenu := &Menu{
		ParentID:  devtoolsMenu.ID,
		MenuName:  "代码生成",
		MenuPath:  "/devtools/codegen",
		Component: "devtools/codegen/CodeGen",
		Icon:      "Code",
		Sort:      1,
		Status:    1,
		MenuType:  1,
		Remark:    "代码生成器",
	}

	if err := db.Create(codegenMenu).Error; err != nil {
		log.Fatal("插入代码生成菜单失败:", err)
	}
	fmt.Printf("✅ 插入代码生成菜单成功, ID: %d\n", codegenMenu.ID)

	// 插入子菜单 - 生成历史
	historyMenu := &Menu{
		ParentID:  devtoolsMenu.ID,
		MenuName:  "生成历史",
		MenuPath:  "/devtools/history",
		Component: "devtools/history/GenHistory",
		Icon:      "List",
		Sort:      2,
		Status:    1,
		MenuType:  1,
		Remark:    "代码生成历史",
	}

	if err := db.Create(historyMenu).Error; err != nil {
		log.Fatal("插入生成历史菜单失败:", err)
	}
	fmt.Printf("✅ 插入生成历史菜单成功, ID: %d\n", historyMenu.ID)

	fmt.Println("\n🎉 菜单初始化完成！")
}
