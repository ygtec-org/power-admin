# 基于角色的菜单列表实现

## 概述

已实现菜单列表接口基于角色权限的过滤功能。用户只能看到所属角色绑定的菜单，未绑定的菜单将不会显示。同时，API访问也受到权限控制。

## 实现原理

### 1. 权限检查流程

```
用户请求菜单列表
    ↓
AdminAuthMiddleware 验证JWT Token
    ↓
提取用户ID并存储到Context和Header中
    ↓
MenuListLogic 获取用户所有角色
    ↓
根据角色ID查询已绑定的菜单
    ↓
构建树形结构并返回给前端
```

### 2. 核心修改文件

#### a) MenuListLogic (`internal/logic/menu/menulistlogic.go`)

**修改内容：**
- 添加了用户认证检查，从context中获取用户ID
- 根据用户ID查询其所有角色
- 如果用户没有任何角色，返回空菜单列表
- 根据用户的所有角色ID查询已绑定的菜单
- 构建菜单树形结构并返回

**关键代码：**
```go
// 从上下文中获取用户ID
userIDStr := l.ctx.Value(constant.AdminUserKey)
userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)

// 获取用户的角色
userRoles, err := l.svcCtx.RoleRepo.GetRolesByUserID(userID)

// 获取用户角色绑定的菜单
menus, err := l.svcCtx.MenuRepo.GetMenusByRoleIDs(roleIDs)
```

#### b) MenuRepository (`pkg/repository/menu.go`)

**新增方法：**
```go
// GetMenusByRoleIDs 根据角色ID列表获取菜单（去重）
func (r *MenuRepository) GetMenusByRoleIDs(roleIDs []int64) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Distinct("menus.*").
		Joins("LEFT JOIN role_menus ON role_menus.menu_id = menus.id").
		Where("role_menus.role_id IN ? AND menus.status = 1", roleIDs).
		Order("menus.sort").
		Find(&menus).Error
	return menus, err
}
```

**特点：**
- 使用DISTINCT避免重复（多个角色可能绑定同一菜单）
- 只查询状态为启用的菜单（status = 1）
- 按排序号升序返回菜单

#### c) RoleRepository (`pkg/repository/role.go`)

**新增方法：**
```go
// GetRolesByUserID 根据用户ID获取所有角色
func (r *RoleRepository) GetRolesByUserID(userID int64) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Joins("LEFT JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.status = 1", userID).
		Find(&roles).Error
	return roles, err
}
```

**特点：**
- 只返回启用的角色（status = 1）
- 通过user_roles关联表进行连接查询

#### d) AdminAuthMiddleware (`internal/middleware/adminauthmiddleware.go`)

**修改内容：**
- 添加context import
- 在JWT验证成功后，将用户ID存储到Request Context中
- 供后续Logic层使用

**关键代码：**
```go
// 将用户ID存储到context中（供logic层使用）
ctx := context.WithValue(r.Context(), constant.AdminUserKey, fmt.Sprintf("%d", claims.ID))
ctx = context.WithValue(ctx, constant.AdminUserName, claims.Username)
r = r.WithContext(ctx)
```

## 数据库表关系

```
users (用户表)
  ↓
user_roles (用户角色关联表)
  ↓
roles (角色表)
  ↓
role_menus (角色菜单关联表)
  ↓
menus (菜单表)
```

## API调用流程

### 获取菜单列表

**请求：**
```
GET /api/admin/system/menus
Authorization: Bearer <JWT_TOKEN>
```

**流程：**
1. Middleware验证JWT Token
2. 提取用户ID（claims.ID）
3. 将用户ID存储到context
4. MenuListHandler调用MenuListLogic
5. MenuListLogic从context获取用户ID
6. 查询用户所有角色
7. 查询这些角色绑定的菜单
8. 构建树形结构返回

**响应示例：**
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "parentId": 0,
      "menuName": "系统管理",
      "menuPath": "/system",
      "component": "Layout",
      "icon": "system",
      "sort": 1,
      "status": 1,
      "menuType": 1,
      "createdAt": "2025-11-30 11:00:00",
      "children": [
        {
          "id": 2,
          "parentId": 1,
          "menuName": "用户管理",
          "menuPath": "/system/users",
          "component": "system/User",
          "icon": "user",
          "sort": 1,
          "status": 1,
          "menuType": 1,
          "createdAt": "2025-11-30 11:00:00",
          "children": []
        }
      ]
    }
  ],
  "total": 1
}
```

## API权限控制

### 白名单路由

菜单列表接口在中间件的白名单中，所以已认证的用户都能访问（无需额外的API权限）：
```go
whitelistRoutes := map[string]bool{
    "/api/admin/system/menus": true, // 菜单接口
    // ...其他白名单路由
}
```

### 其他API权限

其他API的调用受到Casbin RBAC权限控制：
- 用户必须拥有访问该API对应的权限
- 通过角色绑定权限，再将角色分配给用户
- Casbin会检查user→role→permission的权限链

## 错误处理

### 可能的错误场景

1. **用户未认证**
   - 返回401 Unauthorized
   - 原因：缺少Authorization header或token无效

2. **用户没有任何角色**
   - 返回成功，但菜单列表为空
   - 原因：该用户未被分配任何角色

3. **用户角色没有绑定菜单**
   - 返回成功，但菜单列表为空
   - 原因：用户的角色未绑定菜单

4. **数据库查询失败**
   - 返回500 Internal Server Error
   - 原因：数据库连接问题

## 测试建议

### 1. 创建测试数据

```sql
-- 创建用户
INSERT INTO users (username, phone, email, password, nickname) 
VALUES ('testuser', '13800138000', 'test@example.com', 'encrypted_password', 'Test User');

-- 创建角色
INSERT INTO roles (name, description) 
VALUES ('Admin', 'Administrator Role');

-- 创建菜单
INSERT INTO menus (parent_id, menu_name, menu_path, component, icon, sort, status, menu_type)
VALUES (0, 'System', '/system', 'Layout', 'system', 1, 1, 1);

-- 用户分配角色
INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);

-- 角色绑定菜单
INSERT INTO role_menus (role_id, menu_id) VALUES (1, 1);
```

### 2. 测试场景

**场景1：用户有角色，角色绑定了菜单**
- 预期结果：返回该用户角色绑定的菜单

**场景2：用户有多个角色，每个角色绑定不同菜单**
- 预期结果：返回所有角色绑定菜单的并集（去重）

**场景3：用户没有任何角色**
- 预期结果：返回空菜单列表（code:200, data:[], total:0）

**场景4：用户角色没有绑定菜单**
- 预期结果：返回空菜单列表

**场景5：菜单状态为禁用（status=0）**
- 预期结果：即使角色绑定了该菜单，也不会显示

## 性能优化

### 使用的优化措施

1. **DISTINCT查询**
   - 避免多角色绑定同一菜单导致的重复数据

2. **状态过滤**
   - 在数据库层面过滤禁用的菜单和角色，减少应用层处理

3. **联接查询**
   - 使用JOIN而不是逐个查询，减少数据库往返次数

4. **排序**
   - 在数据库层面进行排序，避免应用层排序

## 扩展建议

1. **菜单权限缓存**
   - 可以使用Redis缓存用户的菜单列表，减少数据库查询

2. **权限预加载**
   - 登录时预加载用户的所有菜单权限到Session/JWT中

3. **菜单搜索**
   - 在菜单列表基础上添加搜索功能

4. **菜单分页**
   - 对于菜单数量很多的情况，可以考虑分页

## 总结

该实现完全基于已有的role_menus关联表和用户角色关系，通过在MenuListLogic中增加权限过滤逻辑，确保用户只能看到其角色绑定的菜单。整个流程安全、高效，符合RBAC权限管理的最佳实践。
