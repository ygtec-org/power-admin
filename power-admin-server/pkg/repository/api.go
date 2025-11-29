package repository

import (
	"power-admin-server/pkg/models"

	"gorm.io/gorm"
)

// APIRepository API管理仓储
type APIRepository struct {
	db *gorm.DB
}

// NewAPIRepository 创建API管理仓储
func NewAPIRepository(db *gorm.DB) *APIRepository {
	return &APIRepository{db: db}
}

// Create 创建API
func (r *APIRepository) Create(api *models.API) error {
	return r.db.Create(api).Error
}

// Update 更新API
func (r *APIRepository) Update(api *models.API) error {
	return r.db.Model(api).Updates(api).Error
}

// Delete 删除API
func (r *APIRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.API{}).Error
}

// GetByID 根据ID获取API
func (r *APIRepository) GetByID(id int64) (*models.API, error) {
	var api models.API
	err := r.db.Where("id = ?", id).First(&api).Error
	if err != nil {
		return nil, err
	}
	return &api, nil
}

// GetByPathAndMethod 根据路径和方法获取API
func (r *APIRepository) GetByPathAndMethod(path, method string) (*models.API, error) {
	var api models.API
	err := r.db.Where("api_path = ? AND api_method = ?", path, method).First(&api).Error
	if err != nil {
		return nil, err
	}
	return &api, nil
}

// List 获取API列表
func (r *APIRepository) List(offset, limit int) ([]models.API, int64, error) {
	var apis []models.API
	var total int64

	err := r.db.Model(&models.API{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&apis).Error
	if err != nil {
		return nil, 0, err
	}

	return apis, total, nil
}

// GetAPIs 获取API列表（与List相同）
func (r *APIRepository) GetAPIs(offset, limit int) ([]models.API, int64, error) {
	return r.List(offset, limit)
}

// ListByMethod 根据方法获取API列表
func (r *APIRepository) ListByMethod(method string) ([]models.API, error) {
	var apis []models.API
	err := r.db.Where("api_method = ? AND status = 1", method).Find(&apis).Error
	return apis, err
}

// ListAll 获取所有API
func (r *APIRepository) ListAll() ([]models.API, error) {
	var apis []models.API
	err := r.db.Where("status = 1").Find(&apis).Error
	return apis, err
}
