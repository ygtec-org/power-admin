package svc

import (
	"os"
	"path/filepath"

	"power-admin-server/internal/config"
	"power-admin-server/internal/middleware"
	"power-admin-server/internal/service"
	"power-admin-server/pkg/auth"
	"power-admin-server/pkg/models"
	"power-admin-server/pkg/permission"
	"power-admin-server/pkg/repository"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
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
	AppRepository  repository.AppRepository

	// CMS Repositories
	CmsContentRepo  repository.ContentRepository
	CmsCategoryRepo repository.CategoryRepository
	CmsTagRepo      repository.TagRepository
	CmsUserRepo     repository.CmsUserRepository
	CmsCommentRepo  repository.CommentRepository

	// Services
	PluginService *service.PluginService

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
	appRepository := repository.NewAppRepository(db)

	// 使用配置文件初始化Casbin
	permissionManager, err := permission.NewRBACEnforcer(db, "etc/rbac_model.conf")
	if err != nil {
		panic("Failed to initialize permission manager: " + err.Error())
	}

	// 初始化CMS repositories
	cmsContentRepo := repository.NewContentRepository(db)
	cmsCategoryRepo := repository.NewCategoryRepository(db)
	cmsTagRepo := repository.NewTagRepository(db)
	cmsUserRepo := repository.NewCmsUserRepository(db)
	cmsCommentRepo := repository.NewCommentRepository(db)

	// 初始化PluginService - 使用绝对路径
	pluginService := service.NewPluginService(getPluginsDir())

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
		AppRepository:       appRepository,
		CmsContentRepo:      cmsContentRepo,
		CmsCategoryRepo:     cmsCategoryRepo,
		CmsTagRepo:          cmsTagRepo,
		CmsUserRepo:         cmsUserRepo,
		CmsCommentRepo:      cmsCommentRepo,
		PluginService:       pluginService,
		Permission:          permissionManager,
		AdminAuthMiddleware: middleware.NewAdminAuthMiddleware(&c, permissionManager).Handle,
	}
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 系统模型
		&models.User{},
		&models.Role{},
		&models.Menu{},
		&models.Dictionary{},
		&models.API{},
		&models.App{},
		&models.AppInstallation{},
		&models.Review{},
		&models.Log{},
		&models.Permission{},
		&models.Plugin{},
		&models.RoleMenu{},
		&models.CasbinRule{},

		// CMS模型
		//&models.CmsContent{},
		//&models.CmsCategory{},
		//&models.CmsTag{},
		//&models.CmsContentTag{},
		//&models.CmsContentRevision{},
		//&models.CmsUser{},
		//&models.CmsComment{},
		//&models.CmsPermission{},
		//&models.CmsAdminRole{},
		//&models.CmsPluginStatus{},
		//&models.CmsAuditLog{},
		//&models.CmsLike{},
		//&models.CmsDraft{},
	)
}

// getPluginsDir 获取插件目录的绝对路径
func getPluginsDir() string {
	// 获取可执行文件所在目录
	execPath, err := os.Executable()
	if err != nil {
		logx.Errorf("Failed to get executable path: %v", err)
		return "plugins" // 降级到相对路径
	}

	// 获取可执行文件的目录
	execDir := filepath.Dir(execPath)

	// 如果在 bin 目录中，往上两级到 power-admin 项目根目录
	// 例如: d:\Workspace\project\app\power-admin\power-admin-server\bin\power-admin.exe
	// 我们需要: d:\Workspace\project\app\power-admin\plugins
	// bin -> power-admin-server -> power-admin
	projectRoot := filepath.Join(execDir, "..", "..")
	pluginDir := filepath.Join(projectRoot, "plugins")

	// 规范化路径
	pluginDir, err = filepath.Abs(pluginDir)
	if err != nil {
		logx.Errorf("Failed to get absolute plugins path: %v", err)
		return "plugins"
	}

	logx.Infof("Plugin directory: %s", pluginDir)
	return pluginDir
}
