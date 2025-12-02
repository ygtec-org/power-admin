package repository

import (
	"strings"

	"gorm.io/gorm/schema"
)

// AdminNamingStrategy 自定义命名策略，用于添加表名前缀
type AdminNamingStrategy struct {
	schema.NamingStrategy
	TablePrefix string
}

// TableName 返回带前缀的表名
// 特殊处理 casbin_rule 表，不添加前缀
func (n AdminNamingStrategy) TableName(table string) string {
	// casbin_rule 表不添加前缀
	if table == "casbin_rule" {
		return table
	}

	// 检查是否已经有前缀
	if strings.HasPrefix(table, n.TablePrefix) {
		return table
	}

	// 添加前缀
	return n.TablePrefix + table
}

// JoinTableName 处理多对多关联表名
// 特殊处理带 casbin 的表名，不添加前缀
func (n AdminNamingStrategy) JoinTableName(table string) string {
	// casbin 相关的表不添加前缀
	if strings.Contains(table, "casbin") {
		return table
	}

	// 检查是否已经有前缀
	if strings.HasPrefix(table, n.TablePrefix) {
		return table
	}

	// 添加前缀
	return n.TablePrefix + table
}
