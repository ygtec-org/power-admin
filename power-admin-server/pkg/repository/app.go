package repository

import (
	"power-admin-server/pkg/models"

	"gorm.io/gorm"
)

// AppRepository 应用仓储
type AppRepository struct {
	db *gorm.DB
}

// NewAppRepository 创建应用仓储
func NewAppRepository(db *gorm.DB) *AppRepository {
	return &AppRepository{db: db}
}

// Create 创建应用
func (r *AppRepository) Create(app *models.App) error {
	return r.db.Create(app).Error
}

// Update 更新应用
func (r *AppRepository) Update(app *models.App) error {
	return r.db.Model(app).Updates(app).Error
}

// Delete 删除应用
func (r *AppRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.App{}).Error
}

// GetByID 根据ID获取应用
func (r *AppRepository) GetByID(id int64) (*models.App, error) {
	var app models.App
	err := r.db.Where("id = ?", id).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetByAppKey 根据AppKey获取应用
func (r *AppRepository) GetByAppKey(appKey string) (*models.App, error) {
	var app models.App
	err := r.db.Where("app_key = ?", appKey).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// List 获取应用列表
func (r *AppRepository) List(offset, limit int) ([]models.App, int64, error) {
	var apps []models.App
	var total int64

	err := r.db.Model(&models.App{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Where("published = 1").Order("created_at DESC").Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// ListByCategory 根据分类获取应用列表
func (r *AppRepository) ListByCategory(category string, offset, limit int) ([]models.App, int64, error) {
	var apps []models.App
	var total int64

	err := r.db.Model(&models.App{}).Where("category = ? AND published = 1", category).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("category = ? AND published = 1", category).Offset(offset).Limit(limit).Order("created_at DESC").Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// Search 搜索应用
func (r *AppRepository) Search(keyword string, offset, limit int) ([]models.App, int64, error) {
	var apps []models.App
	var total int64

	query := r.db.Where("published = 1 AND (app_name LIKE ? OR description LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")

	err := query.Model(&models.App{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("downloads DESC").Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// ReviewRepository 评价仓储
type ReviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository 创建评价仓储
func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// Create 创建评价
func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

// Delete 删除评价
func (r *ReviewRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Review{}).Error
}

// GetByID 根据ID获取评价
func (r *ReviewRepository) GetByID(id int64) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("id = ?", id).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// ListByAppID 根据AppID获取评价列表
func (r *ReviewRepository) ListByAppID(appID int64, offset, limit int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	err := r.db.Model(&models.Review{}).Where("app_id = ?", appID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("app_id = ?", appID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}
