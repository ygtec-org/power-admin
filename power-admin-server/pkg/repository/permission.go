package repository

import (
	"power-admin-server/pkg/models"

	"gorm.io/gorm"
)

// PermissionRepository 权限仓储
type PermissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓储
func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// Create 创建权限
func (r *PermissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

// Update 更新权限
func (r *PermissionRepository) Update(permission *models.Permission) error {
	return r.db.Model(permission).Updates(permission).Error
}

// Delete 删除权限
func (r *PermissionRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Permission{}).Error
}

// GetByID 根据ID获取权限
func (r *PermissionRepository) GetByID(id int64) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetByName 根据名称获取权限
func (r *PermissionRepository) GetByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// List 获取权限列表
func (r *PermissionRepository) List(offset, limit int) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	err := r.db.Model(&models.Permission{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// ListByResource 根据资源获取权限列表
func (r *PermissionRepository) ListByResource(resource string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("resource = ? AND status = 1", resource).Find(&permissions).Error
	return permissions, err
}

// DictionaryRepository 字典仓储
type DictionaryRepository struct {
	db *gorm.DB
}

// NewDictionaryRepository 创建字典仓储
func NewDictionaryRepository(db *gorm.DB) *DictionaryRepository {
	return &DictionaryRepository{db: db}
}

// Create 创建字典
func (r *DictionaryRepository) Create(dict *models.Dictionary) error {
	return r.db.Create(dict).Error
}

// Update 更新字典
func (r *DictionaryRepository) Update(dict *models.Dictionary) error {
	return r.db.Model(dict).Updates(dict).Error
}

// Delete 删除字典
func (r *DictionaryRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Dictionary{}).Error
}

// GetByID 根据ID获取字典
func (r *DictionaryRepository) GetByID(id int64) (*models.Dictionary, error) {
	var dict models.Dictionary
	err := r.db.Where("id = ?", id).First(&dict).Error
	if err != nil {
		return nil, err
	}
	return &dict, nil
}

// List 获取字典列表
func (r *DictionaryRepository) List(offset, limit int, dict string) ([]models.Dictionary, int64, error) {
	var dicts []models.Dictionary
	var total int64

	err := r.db.Model(&models.Dictionary{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	where := r.db.Offset(offset).Limit(limit)
	if dict != "" {
		where = where.Where("dict_type = ? OR dict_label LIKE ? OR dict_value LIKE ?", dict, "%"+dict+"%", "%"+dict+"%")
	}
	err = where.Order("sort").Find(&dicts).Error
	if err != nil {
		return nil, 0, err
	}

	return dicts, total, nil
}

// ListByType 根据类型获取字典列表
func (r *DictionaryRepository) ListByType(dictType string) ([]models.Dictionary, error) {
	var dicts []models.Dictionary
	err := r.db.Where("dict_type = ? AND status = 1", dictType).Order("sort").Find(&dicts).Error
	return dicts, err
}

// GetByTypeAndValue 根据类型和值获取字典项
func (r *DictionaryRepository) GetByTypeAndValue(dictType, dictValue string) (*models.Dictionary, error) {
	var dict models.Dictionary
	err := r.db.Where("dict_type = ? AND dict_value = ?", dictType, dictValue).First(&dict).Error
	if err != nil {
		return nil, err
	}
	return &dict, nil
}
