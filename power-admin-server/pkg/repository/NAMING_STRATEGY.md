# GORM 全局表名前缀配置方案

## 问题背景

在项目开发过程中，经常会在 Repository 的查询语句中硬编码表名，导致：
1. **代码重复** - 每个查询都要写表名
2. **维护困难** - 修改表名需要改多个地方
3. **容易出错** - 手写表名容易打错或遗漏前缀

例如：
```go
// ❌ 不推荐：硬编码表名
err := r.db.Joins("LEFT JOIN user_roles ON user_roles.role_id = roles.id").
    Where("user_roles.user_id = ? AND roles.status = 1", userID).
    Find(&roles).Error
```

## 解决方案

使用 GORM 的 `NamingStrategy` 自定义命名策略，全局配置表名前缀。

### 1. 自定义命名策略 (naming.go)

```go
type AdminNamingStrategy struct {
    schema.NamingStrategy
    TablePrefix string
}

// 在表名前添加前缀
func (n AdminNamingStrategy) TableName(table string) string {
    // casbin_rule 表特殊处理，不添加前缀
    if table == "casbin_rule" {
        return table
    }
    if strings.HasPrefix(table, n.TablePrefix) {
        return table
    }
    return n.TablePrefix + table
}

// 处理多对多关联表名
func (n AdminNamingStrategy) JoinTableName(table string) string {
    // 特殊处理 casbin 相关表
    if strings.Contains(table, "casbin") {
        return table
    }
    if strings.HasPrefix(table, n.TablePrefix) {
        return table
    }
    return n.TablePrefix + table
}
```

### 2. 在 GORM 初始化时应用策略

**文件**: `internal/svc/servicecontext.go`

```go
// 创建自定义命名策略，全局添加 admin_ 前缀
naming := repository.AdminNamingStrategy{
    TablePrefix: "admin_",
}

db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{
    NamingStrategy: naming,
})
```

### 3. 使用示例

#### ✅ 推荐方式：使用 admin_ 前缀的表名

```go
// 使用带前缀的表名
err := r.db.Joins("LEFT JOIN admin_user_roles ON admin_user_roles.role_id = admin_roles.id").
    Where("admin_user_roles.user_id = ? AND admin_roles.status = 1", userID).
    Find(&roles).Error
```

#### ✅ 最佳实践：使用 GORM 关联

对于关联关系，尽量使用 GORM 的 Preload 或 Association 方法：

```go
// 最佳实践：使用 GORM 关联
var user models.User
r.db.Preload("Roles").First(&user, userID)
// Roles 会自动通过 many2many 关联加载
```

## 优势

| 特点 | 说明 |
|------|------|
| **全局配置** | 一处配置，全局生效，无需在每个查询中重复 |
| **自动处理** | GORM 会自动在 SQL 中使用正确的表名 |
| **灵活支持** | 支持特殊表（如 casbin_rule）的排除处理 |
| **易于维护** | 修改前缀只需改一处，所有查询自动更新 |
| **向后兼容** | 已有的 TableName() 方法仍然有效 |

## 数据库表对照表

| 逻辑名 | 数据库表名 | 前缀处理 |
|--------|----------|---------|
| users | admin_users | ✅ 自动添加前缀 |
| roles | admin_roles | ✅ 自动添加前缀 |
| permissions | admin_permissions | ✅ 自动添加前缀 |
| menus | admin_menus | ✅ 自动添加前缀 |
| user_roles | admin_user_roles | ✅ 自动添加前缀 |
| role_permissions | admin_role_permissions | ✅ 自动添加前缀 |
| role_menus | admin_role_menus | ✅ 自动添加前缀 |
| casbin_rule | casbin_rule | ❌ 不添加前缀（特殊处理） |

## 迁移指南

如果项目中已有硬编码表名的查询，可以按以下步骤迁移：

1. **第一步** - 应用命名策略（已完成）
2. **第二步** - 逐步更新查询语句，使用带前缀的表名
3. **第三步** - 优先使用 GORM 关联功能（Preload, Association）

## 注意事项

1. **特殊表处理** - casbin_rule 已在命名策略中特殊处理，无需添加前缀
2. **自定义表** - 如果添加新的业务表，模型中的 `TableName()` 方法仍然有效，会覆盖命名策略
3. **性能** - 命名策略处理是在 SQL 生成阶段进行的，不会有性能影响
