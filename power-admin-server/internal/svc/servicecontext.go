package svc

import (
	"power-admin-server/internal/config"
	"power-admin-server/internal/middleware"
	"power-admin-server/pkg/auth"
	"power-admin-server/pkg/models"
	"power-admin-server/pkg/permission"
	"power-admin-server/pkg/repository"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client

	// Repositories
	UserRepo       *repository.UserRepository
	RoleRepo       *repository.RoleRepository
	MenuRepo       *repository.MenuRepository
	PermissionRepo *repository.PermissionRepository
	DictRepo       *repository.DictionaryRepository
	APIRepo        *repository.APIRepository
	AppRepo        *repository.AppRepository
	ReviewRepo     *repository.ReviewRepository

	// Permission
	Permission          *permission.RBACEnforcer
	AdminAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {

	// 设置 JWT Secret
	auth.SetJwtSecret(c.Auth.AccessSecret)

	// 这里完成mysql和Redis的初始化
	rds := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
		DB:       c.Redis.Db,
	})
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	err = autoMigrate(db)
	if err != nil {
		panic("Failed to auto migrate database: " + err.Error())
	}
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	dictRepo := repository.NewDictionaryRepository(db)
	apiRepo := repository.NewAPIRepository(db)
	appRepo := repository.NewAppRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	// 使用配置文件初始化Casbin
	permissionManager, err := permission.NewRBACEnforcer(db, "etc/rbac_model.conf")
	if err != nil {
		panic("Failed to initialize permission manager: " + err.Error())
	}

	return &ServiceContext{
		Config:              c,
		DB:                  db,
		Redis:               rds,
		UserRepo:            userRepo,
		RoleRepo:            roleRepo,
		MenuRepo:            menuRepo,
		PermissionRepo:      permissionRepo,
		DictRepo:            dictRepo,
		APIRepo:             apiRepo,
		AppRepo:             appRepo,
		ReviewRepo:          reviewRepo,
		Permission:          permissionManager,
		AdminAuthMiddleware: middleware.NewAdminAuthMiddleware(&c, permissionManager).Handle,
	}
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 在这里添加需要自动迁移的模型
		&models.User{},
		&models.Role{},
		&models.Menu{},
		&models.Dictionary{},
		&models.API{},
		&models.App{},
		&models.Review{},
		&models.Log{},
		&models.Permission{},
		&models.Plugin{},
		&models.RoleMenu{},
		&models.CasbinRule{},
	)
}
