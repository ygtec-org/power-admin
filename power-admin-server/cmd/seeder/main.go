package main

import (
	"fmt"

	"power-admin-server/pkg/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// MySQL 连接字符串
	dsn := "root:root@tcp(127.0.0.1:3306)/power_admin?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	fmt.Println("✓ 数据库连接成功")

	// 创建表
	err = db.AutoMigrate(
		&models.App{},
		&models.AppInstallation{},
	)
	if err != nil {
		panic("Failed to auto migrate: " + err.Error())
	}

	fmt.Println("✓ 表结构已创建")

	// 清空旧数据
	db.Exec("DELETE FROM apps")
	db.Exec("DELETE FROM app_installations")

	// 插入应用数据（使用稳定的 placeholder 图片）
	apps := []models.App{
		{
			AppKey:      "power-cms",
			AppName:     "CMS内容系统",
			Version:     "1.0.0",
			Author:      "Power Admin Team",
			Description: "完整的内容管理系统，支持文章、分类、标签、评论管理，适合博客和内容运营",
			Icon:        "https://picsum.photos/300/200?random=1",
			Category:    "content",
			Tags:        "内容管理,博客",
			Rating:      4.8,
			Downloads:   1250,
			Status:      1,
			Published:   1,
		},
		{
			AppKey:      "shop",
			AppName:     "电商系统",
			Version:     "2.0.0",
			Author:      "Power Admin Team",
			Description: "功能完整的电商平台，包含商品管理、订单、支付、物流跟踪等功能",
			Icon:        "https://picsum.photos/300/200?random=2",
			Category:    "business",
			Tags:        "电商,商城",
			Rating:      4.6,
			Downloads:   980,
			Status:      1,
			Published:   1,
		},
		{
			AppKey:      "crm",
			AppName:     "CRM客户管理",
			Version:     "1.5.0",
			Author:      "Power Admin Team",
			Description: "企业客户关系管理系统，支持客户跟进、商机管理、合同管理等",
			Icon:        "https://picsum.photos/300/200?random=3",
			Category:    "business",
			Tags:        "客户管理,销售",
			Rating:      4.5,
			Downloads:   750,
			Status:      1,
			Published:   1,
		},
		{
			AppKey:      "finance",
			AppName:     "财务管理系统",
			Version:     "1.2.0",
			Author:      "Power Admin Team",
			Description: "企业财务管理工具，涵盖应收应付、报表、成本管理等功能",
			Icon:        "https://picsum.photos/300/200?random=4",
			Category:    "finance",
			Tags:        "财务,报表",
			Rating:      4.7,
			Downloads:   680,
			Status:      1,
			Published:   1,
		},
		{
			AppKey:      "hr",
			AppName:     "HR人力资源",
			Version:     "1.1.0",
			Author:      "Power Admin Team",
			Description: "人力资源管理平台，支持招聘、员工档案、考勤、薪资管理",
			Icon:        "https://picsum.photos/300/200?random=5",
			Category:    "business",
			Tags:        "人力资源,招聘",
			Rating:      4.4,
			Downloads:   520,
			Status:      1,
			Published:   1,
		},
		{
			AppKey:      "marketing",
			AppName:     "营销工具箱",
			Version:     "1.0.0",
			Author:      "Power Admin Team",
			Description: "集成营销工具，支持活动管理、积分管理、优惠券等营销功能",
			Icon:        "https://picsum.photos/300/200?random=6",
			Category:    "marketing",
			Tags:        "营销,促销",
			Rating:      4.3,
			Downloads:   400,
			Status:      1,
			Published:   1,
		},
	}

	if err := db.Create(&apps).Error; err != nil {
		panic("Failed to insert apps: " + err.Error())
	}

	fmt.Printf("✓ 应用数据插入成功 (%d 条)\n", len(apps))

	// 插入已安装的应用记录（仅 power-cms）
	cmsApp := &models.App{}
	db.Where("app_key = ?", "power-cms").First(cmsApp)

	if cmsApp.ID > 0 {
		installation := &models.AppInstallation{
			AppKey:  "power-cms",
			AppID:   cmsApp.ID,
			AppName: "CMS内容系统",
			Version: "1.0.0",
			Status:  1,
		}

		if err := db.Create(installation).Error; err != nil {
			panic("Failed to insert app installation: " + err.Error())
		}

		fmt.Println("✓ CMS应用标记为已安装")
	}

	fmt.Println("\n✓ 数据初始化完成！")
}
