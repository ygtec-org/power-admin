package permission

import (
	"fmt"
	"log"

	"power-admin-server/pkg/models"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

// RBACEnforcer Casbin RBAC权限管理
type RBACEnforcer struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
}

// NewRBACEnforcer 初始化Casbin权限管理
// 由于MySQL索引长度限制，我们使用内存模式并手动同步到数据库
func NewRBACEnforcer(db *gorm.DB, modelPath string) (*RBACEnforcer, error) {
	// 删除旧表重建
	db.Exec("DROP TABLE IF EXISTS `casbin_rule`")
	db.Exec("DROP TABLE IF EXISTS `casbin_rules`")

	// 手动创建表，不含任何UNIQUE索引
	db.Exec(`
		CREATE TABLE IF NOT EXISTS casbin_rule (
			id INT AUTO_INCREMENT PRIMARY KEY,
			ptype VARCHAR(10) NOT NULL,
			v0 VARCHAR(50),
			v1 VARCHAR(50),
			v2 VARCHAR(50),
			v3 VARCHAR(50),
			v4 VARCHAR(50),
			v5 VARCHAR(50),
			INDEX idx_ptype (ptype),
			INDEX idx_v0 (v0)
		) DEFAULT CHARSET=utf8mb4
	`)

	// 不使用adapter，直接创建内存Enforcer
	enforcer, err := casbin.NewEnforcer(modelPath)
	if err != nil {
		log.Printf("Failed to create enforcer: %v", err)
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	// 从数据库加载已保存的策略到Enforcer
	var rules []models.CasbinRule
	if err := db.Find(&rules).Error; err != nil {
		log.Printf("Failed to load rules from database: %v", err)
	} else {
		for _, rule := range rules {
			tokens := []interface{}{rule.V0, rule.V1, rule.V2, rule.V3, rule.V4, rule.V5}
			// 移除末尾空字符串
			for len(tokens) > 0 && tokens[len(tokens)-1] == "" {
				tokens = tokens[:len(tokens)-1]
			}
			enforcer.AddPolicy(tokens...)
		}
	}

	return &RBACEnforcer{enforcer: enforcer, db: db}, nil
}

// CheckPermission 检查权限
func (r *RBACEnforcer) CheckPermission(subject, object, action string) bool {
	ok, err := r.enforcer.Enforce(subject, object, action)
	if err != nil {
		log.Printf("Failed to check permission: %v", err)
		return false
	}
	return ok
}

// AddRoleForUser 为用户添加角色
func (r *RBACEnforcer) AddRoleForUser(user, role string) error {
	ok, err := r.enforcer.AddGroupingPolicy(user, role)
	if err != nil {
		return fmt.Errorf("failed to add role: %w", err)
	}
	if !ok {
		return fmt.Errorf("failed to add role: policy already exists")
	}
	return nil
}

// RemoveRoleForUser 移除用户角色
func (r *RBACEnforcer) RemoveRoleForUser(user, role string) error {
	ok, err := r.enforcer.RemoveGroupingPolicy(user, role)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}
	if !ok {
		return fmt.Errorf("failed to remove role: policy not found")
	}
	return nil
}

// AddPermissionForRole 为角色添加权限
func (r *RBACEnforcer) AddPermissionForRole(role, object, action string) error {
	ok, err := r.enforcer.AddPolicy(role, object, action)
	if err != nil {
		return fmt.Errorf("failed to add permission: %w", err)
	}
	if !ok {
		return fmt.Errorf("failed to add permission: policy already exists")
	}

	// 同时保存到数据库
	casbinRule := models.CasbinRule{
		PType: "p",
		V0:    role,
		V1:    object,
		V2:    action,
	}
	return r.db.Create(&casbinRule).Error
}

// RemovePermissionForRole 移除角色权限
func (r *RBACEnforcer) RemovePermissionForRole(role, object, action string) error {
	ok, err := r.enforcer.RemovePolicy(role, object, action)
	if err != nil {
		return fmt.Errorf("failed to remove permission: %w", err)
	}
	if !ok {
		return fmt.Errorf("failed to remove permission: policy not found")
	}

	// 同时从数据库删除
	return r.db.Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", "p", role, object, action).
		Delete(&models.CasbinRule{}).Error
}

// GetRolesForUser 获取用户所有角色
func (r *RBACEnforcer) GetRolesForUser(user string) ([]string, error) {
	roles, err := r.enforcer.GetRolesForUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	return roles, nil
}

// GetPermissionsForRole 获取角色所有权限
func (r *RBACEnforcer) GetPermissionsForRole(role string) ([][]string, error) {
	permissions, err := r.enforcer.GetPermissionsForUser(role)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}
	return permissions, nil
}
