package repository

import (
	"power-admin-server/pkg/models"

	"gorm.io/gorm"
)

// MenuRepository 菜单仓储
type MenuRepository struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓储
func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// Create 创建菜单
func (r *MenuRepository) Create(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

// Update 更新菜单
func (r *MenuRepository) Update(menu *models.Menu) error {
	return r.db.Model(menu).Updates(menu).Error
}

// Delete 删除菜单
func (r *MenuRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Menu{}).Error
}

// GetByID 根据ID获取菜单
func (r *MenuRepository) GetByID(id int64) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetByParentID 根据父ID获取子菜单
func (r *MenuRepository) GetByParentID(parentID int64) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Where("parent_id = ? AND status = 1", parentID).Order("sort").Find(&menus).Error
	return menus, err
}

// List 获取菜单列表（分页）
func (r *MenuRepository) List(offset, limit, parentId int) ([]models.Menu, int64, error) {
	var menus []models.Menu
	var total int64

	err := r.db.Model(&models.Menu{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	where := r.db.Where("status = 1")
	if parentId > 0 {
		where = where.Where("parent_id = ?", parentId)
	}
	if offset > 0 {
		where = where.Offset(offset)
	}
	if limit > 0 {
		where = where.Limit(limit)
	}
	err = where.Order("sort").Find(&menus).Error
	if err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

// All 获取全部菜单
func (r *MenuRepository) All(parentId int64) ([]models.Menu, int64, error) {
	var menus []models.Menu
	var total int64

	err := r.db.Model(&models.Menu{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	where := r.db.Where("status = 1")
	if parentId > 0 {
		where = where.Where("parent_id = ?", parentId)
	}
	err = where.Order("sort").Find(&menus).Error
	if err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

// GetMenuTree 获取菜单树形结构
func (r *MenuRepository) GetMenuTree(parentID int64) ([]*models.Menu, error) {
	var menus []models.Menu
	err := r.db.Where("parent_id = ? AND status = 1", parentID).Order("sort").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 递归加载子菜单
	var result []*models.Menu
	for i := range menus {
		children, err := r.GetMenuTree(menus[i].ID)
		if err != nil {
			return nil, err
		}
		menus[i].Children = children
		result = append(result, &menus[i])
	}

	return result, nil
}

// GetMenusByRole 根据角色获取菜单
func (r *MenuRepository) GetMenusByRole(roleID int64) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Joins("LEFT JOIN role_menus ON role_menus.menu_id = menus.id").
		Where("role_menus.role_id = ? AND menus.status = 1", roleID).
		Order("menus.sort").
		Find(&menus).Error
	return menus, err
}

// GetMenusTreeByRole 根据角色获取菜单树形结构
func (r *MenuRepository) GetMenusTreeByRole(roleID int64, parentID int64) ([]*models.Menu, error) {
	var menus []models.Menu
	err := r.db.Joins("LEFT JOIN role_menus ON role_menus.menu_id = menus.id").
		Where("role_menus.role_id = ? AND menus.parent_id = ? AND menus.status = 1", roleID, parentID).
		Order("menus.sort").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 递归加载子菜单
	var result []*models.Menu
	for i := range menus {
		children, err := r.GetMenusTreeByRole(roleID, menus[i].ID)
		if err != nil {
			return nil, err
		}
		menus[i].Children = children
		result = append(result, &menus[i])
	}

	return result, nil
}

// GetMenusByRoleIDs 根据角色ID列表获取菜单（去重）
func (r *MenuRepository) GetMenusByRoleIDs(roleIDs []int64) ([]models.Menu, error) {
	var menus []models.Menu
	// 使用DISTINCT来避免重复，因为多个角色可能绑定了同一个菜单
	err := r.db.Distinct("menus.*").
		Joins("LEFT JOIN role_menus ON role_menus.menu_id = menus.id").
		Where("role_menus.role_id IN ? AND menus.status = 1", roleIDs).
		Order("menus.sort").
		Find(&menus).Error
	return menus, err
}
