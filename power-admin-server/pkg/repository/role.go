package repository

import (
	"power-admin-server/pkg/models"

	"gorm.io/gorm"
)

// RoleRepository 角色仓储
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓储
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create 创建角色
func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// Update 更新角色
func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Model(role).Updates(role).Error
}

// Delete 删除角色
func (r *RoleRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Role{}).Error
}

// GetByID 根据ID获取角色
func (r *RoleRepository) GetByID(id int64) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").Preload("Menus").Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByName 根据名称获取角色
func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").Preload("Menus").Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List 获取角色列表
func (r *RoleRepository) List(offset, limit int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	err := r.db.Model(&models.Role{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Permissions").Preload("Menus").Offset(offset).Limit(limit).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// IsNameExist 检查角色名是否存在
func (r *RoleRepository) IsNameExist(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Role{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// ExistsByName 检查角色名是否存在（别名）
func (r *RoleRepository) ExistsByName(name string) (bool, error) {
	return r.IsNameExist(name)
}

// GetRoles 获取角色列表（与List相同）
func (r *RoleRepository) GetRoles(offset, limit int) ([]models.Role, int64, error) {
	return r.List(offset, limit)
}

// GetPermissions 获取角色的权限
func (r *RoleRepository) GetPermissions(roleID int64) ([]*models.Permission, error) {
	role, err := r.GetByID(roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	return role.Permissions, nil
}

// AddPermission 为角色添加权限
func (r *RoleRepository) AddPermission(roleID, permissionID int64) error {
	return r.db.Model(&models.Role{}).Where("id = ?", roleID).Association("Permissions").Append(&models.Permission{ID: permissionID})
}

// RemovePermission 移除角色权限
func (r *RoleRepository) RemovePermission(roleID, permissionID int64) error {
	return r.db.Model(&models.Role{}).Where("id = ?", roleID).Association("Permissions").Delete(&models.Permission{ID: permissionID})
}

// RemoveAllPermissions 移除角色的所有权限
func (r *RoleRepository) RemoveAllPermissions(roleID int64) error {
	return r.db.Model(&models.Role{}).Where("id = ?", roleID).Association("Permissions").Clear()
}

// AddMenu 为角色添加菜单
func (r *RoleRepository) AddMenu(roleID, menuID int64) error {
	return r.db.Model(&models.Role{}).Where("id = ?", roleID).Association("Menus").Append(&models.Menu{ID: menuID})
}

// RemoveMenu 移除角色菜单
func (r *RoleRepository) RemoveMenu(roleID, menuID int64) error {
	return r.db.Model(&models.Role{}).Where("id = ?", roleID).Association("Menus").Delete(&models.Menu{ID: menuID})
}

// RemoveAllMenus 移除角色的所有菜单
func (r *RoleRepository) RemoveAllMenus(roleID int64) error {
	return r.db.Model(&models.Role{}).Where("id = ?", roleID).Association("Menus").Clear()
}

// GetRolesByUserID 根据用户ID获取上所有角色
func (r *RoleRepository) GetRolesByUserID(userID int64) ([]models.Role, error) {
	var roles []models.Role
	// 使用 GORM 的 Joins 是会自动处理表名前缀
	// 不需要手动指定表名，自定义命名策略会自动处理
	err := r.db.Joins("LEFT JOIN admin_user_roles ON admin_user_roles.role_id = admin_roles.id").
		Where("admin_user_roles.user_id = ? AND admin_roles.status = 1", userID).
		Find(&roles).Error
	return roles, err
}
