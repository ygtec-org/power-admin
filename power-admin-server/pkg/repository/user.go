package repository

import (
	"power-admin-server/pkg/models"

	"gorm.io/gorm"
)

// UserRepository 用户仓储
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Model(user).Updates(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号获取用户
func (r *UserRepository) GetByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表
func (r *UserRepository) List(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	err := r.db.Model(&models.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Roles").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// IsPhoneExist 检查手机号是否存在
func (r *UserRepository) IsPhoneExist(phone string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error
	return count > 0, err
}

// IsUsernameExist 检查用户名是否存在
func (r *UserRepository) IsUsernameExist(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// IsEmailExist 检查邮箱是否存在
func (r *UserRepository) IsEmailExist(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// AddRole 为用户添加角色
func (r *UserRepository) AddRole(userID, roleID int64) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Association("Roles").Append(&models.Role{ID: roleID})
}

// RemoveRole 移除用户角色
func (r *UserRepository) RemoveRole(userID, roleID int64) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Association("Roles").Delete(&models.Role{ID: roleID})
}

// GetRoles 获取用户的所有角色
func (r *UserRepository) GetRoles(userID int64) ([]*models.Role, error) {
	user, err := r.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return user.Roles, nil
}

// RemoveAllRoles 移除用户的所有角色
func (r *UserRepository) RemoveAllRoles(userID int64) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Association("Roles").Clear()
}
