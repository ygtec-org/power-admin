package repository

import (
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type AppRepository interface {
	// 获取应用列表
	GetAppList(page, pageSize int, category string) ([]models.App, int64, error)
	// 按应用标识获取应用
	GetAppByKey(appKey string) (*models.App, error)
	// 按ID获取应用
	GetAppByID(id int64) (*models.App, error)
	// 创建应用
	CreateApp(app *models.App) error
	// 更新应用
	UpdateApp(app *models.App) error
	// 删除应用
	DeleteApp(id int64) error
	// 搜索应用
	SearchApps(keyword string, page, pageSize int) ([]models.App, int64, error)

	// 应用安装相关
	// 检查应用是否已安装
	IsAppInstalled(appKey string) (bool, error)
	// 获取已安装应用列表
	GetInstalledApps() ([]models.AppInstallation, error)
	// 安装应用
	InstallApp(installation *models.AppInstallation) error
	// 卸载应用
	UninstallApp(appKey string) error
	// 获取应用安装记录
	GetAppInstallation(appKey string) (*models.AppInstallation, error)
}

type AppRepositoryImpl struct {
	conn sqlx.SqlConn
	db   *gorm.DB
}

func NewAppRepository(db *gorm.DB) AppRepository {
	return &AppRepositoryImpl{
		db: db,
	}
}

// GetAppList 获取应用列表
func (r *AppRepositoryImpl) GetAppList(page, pageSize int, category string) ([]models.App, int64, error) {
	var apps []models.App
	var total int64

	query := r.db.Where("status = ?", 1)
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 获取总数
	if err := query.Model(&models.App{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&apps).Error; err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// GetAppByKey 按应用标识获取应用
func (r *AppRepositoryImpl) GetAppByKey(appKey string) (*models.App, error) {
	var app models.App
	if err := r.db.Where("app_key = ?", appKey).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// GetAppByID 按ID获取应用
func (r *AppRepositoryImpl) GetAppByID(id int64) (*models.App, error) {
	var app models.App
	if err := r.db.Where("id = ?", id).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// CreateApp 创建应用
func (r *AppRepositoryImpl) CreateApp(app *models.App) error {
	return r.db.Create(app).Error
}

// UpdateApp 更新应用
func (r *AppRepositoryImpl) UpdateApp(app *models.App) error {
	return r.db.Save(app).Error
}

// DeleteApp 删除应用
func (r *AppRepositoryImpl) DeleteApp(id int64) error {
	return r.db.Delete(&models.App{}, id).Error
}

// SearchApps 搜索应用
func (r *AppRepositoryImpl) SearchApps(keyword string, page, pageSize int) ([]models.App, int64, error) {
	var apps []models.App
	var total int64

	query := r.db.Where("status = ? AND (app_name LIKE ? OR description LIKE ?)", 1, "%"+keyword+"%", "%"+keyword+"%")

	// 获取总数
	if err := query.Model(&models.App{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&apps).Error; err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// IsAppInstalled 检查应用是否已安装
func (r *AppRepositoryImpl) IsAppInstalled(appKey string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.AppInstallation{}).Where("app_key = ? AND status = ?", appKey, 1).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetInstalledApps 获取已安装应用列表
func (r *AppRepositoryImpl) GetInstalledApps() ([]models.AppInstallation, error) {
	var installations []models.AppInstallation
	if err := r.db.Where("status = ?", 1).Find(&installations).Error; err != nil {
		return nil, err
	}
	return installations, nil
}

// InstallApp 安装应用
func (r *AppRepositoryImpl) InstallApp(installation *models.AppInstallation) error {
	installation.Status = 1
	return r.db.Create(installation).Error
}

// UninstallApp 卸载应用
func (r *AppRepositoryImpl) UninstallApp(appKey string) error {
	return r.db.Model(&models.AppInstallation{}).Where("app_key = ?", appKey).Update("status", 0).Error
}

// GetAppInstallation 获取应用安装记录
func (r *AppRepositoryImpl) GetAppInstallation(appKey string) (*models.AppInstallation, error) {
	var installation models.AppInstallation
	if err := r.db.Where("app_key = ?", appKey).First(&installation).Error; err != nil {
		return nil, err
	}
	return &installation, nil
}
