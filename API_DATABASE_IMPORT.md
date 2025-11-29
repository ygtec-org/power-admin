# API接口数据导入指南

## 📋 概述

本文档说明如何将系统所有的API接口数据导入到数据库的`apis`表中，实现API的动态管理。

---

## ✅ 已插入的API接口统计

| 分类 | 数量 | 接口类型 |
|------|------|--------|
| 认证API | 3个 | 登录、登出、获取用户信息 |
| 用户管理 | 7个 | 列表、创建、编辑、删除、详情、分配角色、获取角色 |
| 角色管理 | 7个 | 列表、创建、编辑、删除、详情、分配权限、获取权限 |
| 菜单管理 | 5个 | 列表、创建、编辑、删除、详情 |
| 权限管理 | 5个 | 列表、创建、编辑、删除、详情 |
| 字典管理 | 5个 | 列表、创建、编辑、删除、详情 |
| API管理 | 5个 | 列表、创建、编辑、删除、详情 |
| **总计** | **37个** | **完整的管理系统API** |

---

## 🚀 导入方式

### 方式1：使用初始化脚本（推荐）

系统已将所有API数据集成到主初始化脚本中。

**步骤1：检查数据库是否存在**
```bash
mysql -u root -p -e "SHOW DATABASES LIKE 'power_admin';"
```

**步骤2：导入初始化脚本**
```bash
# 方式1：使用MySQL客户端
mysql -u root -p < d:/Workspace/project/app/power-admin/power-admin-server/db/init.sql

# 方式2：在数据库管理工具中执行
# 打开数据库工具（如Navicat、DBeaver等）
# 选择 power_admin 数据库
# 打开并执行 init.sql 文件
```

**步骤3：验证导入结果**
```sql
-- 查看导入的API数量
SELECT COUNT(*) as api_count FROM apis;

-- 查看所有API接口
SELECT id, api_name, api_path, api_method, description FROM apis ORDER BY id;

-- 按分类统计API
SELECT 
  CASE 
    WHEN api_path LIKE '%/auth/%' THEN '认证API'
    WHEN api_path LIKE '%/system/users%' THEN '用户管理'
    WHEN api_path LIKE '%/system/roles%' THEN '角色管理'
    WHEN api_path LIKE '%/system/menus%' THEN '菜单管理'
    WHEN api_path LIKE '%/system/permissions%' THEN '权限管理'
    WHEN api_path LIKE '%/content/dicts%' THEN '字典管理'
    WHEN api_path LIKE '%/system/apis%' THEN 'API管理'
    ELSE '其他'
  END as category,
  COUNT(*) as count
FROM apis
GROUP BY category
ORDER BY count DESC;
```

---

### 方式2：单独导入API数据

如果已有数据库，只想导入API数据，可使用单独的脚本。

```bash
mysql -u root -p power_admin < d:/Workspace/project/app/power-admin/power-admin-server/db/insert_apis.sql
```

---

## 📊 API接口详细清单

### 1. 认证API（3个）

| 接口名称 | 路径 | 方法 | 描述 |
|--------|------|------|------|
| 用户登录 | `/api/admin/auth/login` | POST | 用户登录接口 |
| 用户登出 | `/api/admin/auth/logout` | POST | 用户登出接口 |
| 获取用户信息 | `/api/admin/auth/info` | GET | 获取当前登录用户信息 |

### 2. 用户管理API（7个）

| 接口名称 | 路径 | 方法 | 描述 |
|--------|------|------|------|
| 获取用户列表 | `/api/admin/system/users` | GET | 分页获取用户列表 |
| 创建用户 | `/api/admin/system/users` | POST | 创建新用户 |
| 编辑用户 | `/api/admin/system/users` | PUT | 编辑用户信息 |
| 删除用户 | `/api/admin/system/users` | DELETE | 删除用户 |
| 获取用户详情 | `/api/admin/system/users/:id` | GET | 根据ID获取用户详情 |
| 为用户分配角色 | `/api/admin/system/users/:id/roles` | POST | 为用户分配角色 |
| 获取用户角色 | `/api/admin/system/users/:id/roles` | GET | 获取用户已分配的角色 |

### 3. 角色管理API（7个）

| 接口名称 | 路径 | 方法 | 描述 |
|--------|------|------|------|
| 获取角色列表 | `/api/admin/system/roles` | GET | 分页获取角色列表 |
| 创建角色 | `/api/admin/system/roles` | POST | 创建新角色 |
| 编辑角色 | `/api/admin/system/roles` | PUT | 编辑角色信息 |
| 删除角色 | `/api/admin/system/roles` | DELETE | 删除角色 |
| 获取角色详情 | `/api/admin/system/roles/:id` | GET | 根据ID获取角色详情 |
| 为角色分配权限 | `/api/admin/system/roles/:id/permissions` | POST | 为角色分配权限 |
| 获取角色权限 | `/api/admin/system/roles/:id/permissions` | GET | 获取角色已分配的权限 |

### 4. 菜单管理API（5个）

| 接口名称 | 路径 | 方法 | 描述 |
|--------|------|------|------|
| 获取菜单列表 | `/api/admin/system/menus` | GET | 分页获取菜单列表 |
| 创建菜单 | `/api/admin/system/menus` | POST | 创建新菜单 |
| 编辑菜单 | `/api/admin/system/menus` | PUT | 编辑菜单信息 |
| 删除菜单 | `/api/admin/system/menus` | DELETE | 删除菜单 |
| 获取菜单详情 | `/api/admin/system/menus/:id` | GET | 根据ID获取菜单详情 |

### 5. 权限管理API（5个）

| 接口名称 | 路径 | 方法 | 描述 |
|--------|------|------|------|
| 获取权限列表 | `/api/admin/system/permissions` | GET | 分页获取权限列表 |
| 创建权限 | `/api/admin/system/permissions` | POST | 创建新权限 |
| 编辑权限 | `/api/admin/system/permissions` | PUT | 编辑权限信息 |
| 删除权限 | `/api/admin/system/permissions` | DELETE | 删除权限 |
| 获取权限详情 | `/api/admin/system/permissions/:id` | GET | 根据ID获取权限详情 |

### 6. 字典管理API（5个）

| 接口名称 | 路径 | 方法 | 描述 |
|--------|------|------|------|
| 获取字典列表 | `/api/admin/content/dicts` | GET | 分页获取字典列表 |
| 创建字典 | `/api/admin/content/dicts` | POST | 创建新字典项 |
| 编辑字典 | `/api/admin/content/dicts` | PUT | 编辑字典项 |
| 删除字典 | `/api/admin/content/dicts` | DELETE | 删除字典项 |
| 获取字典详情 | `/api/admin/content/dicts/:id` | GET | 根据ID获取字典项详情 |

### 7. API管理API（5个）

| 接口名称 | 路径 | 方法 | 描述 |
|--------|------|------|------|
| 获取API列表 | `/api/admin/system/apis` | GET | 分页获取API列表 |
| 创建API | `/api/admin/system/apis` | POST | 创建新API接口记录 |
| 编辑API | `/api/admin/system/apis` | PUT | 编辑API接口信息 |
| 删除API | `/api/admin/system/apis` | DELETE | 删除API接口记录 |
| 获取API详情 | `/api/admin/system/apis/:id` | GET | 根据ID获取API接口详情 |

---

## 🔧 API表结构说明

```sql
CREATE TABLE apis (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,           -- API ID
    api_name VARCHAR(100) NOT NULL,                 -- API名称
    api_path VARCHAR(255) NOT NULL,                 -- API路径
    api_method VARCHAR(10) NOT NULL,                -- HTTP方法 (GET, POST, PUT, DELETE)
    description VARCHAR(255),                       -- API描述
    permission_id BIGINT,                           -- 关联权限ID（可选）
    status INT DEFAULT 1,                           -- 状态 (1:启用, 0:禁用)
    remark VARCHAR(255),                            -- 备注
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP,                           -- 更新时间
    deleted_at TIMESTAMP NULL,                      -- 删除时间
    UNIQUE KEY uk_path_method (api_path, api_method),
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
);
```

---

## 💡 使用建议

### 1. API管理功能
已导入的API数据可通过前端"API管理"页面进行管理：
- ✅ 查看所有已注册的API
- ✅ 启用/禁用API
- ✅ 编辑API描述信息
- ✅ 关联权限

### 2. 权限绑定
为了实现完整的权限管理，可将API与权限关联：
```sql
-- 示例：将"创建用户"API与"user_create"权限关联
UPDATE apis 
SET permission_id = (SELECT id FROM permissions WHERE name = 'user_create')
WHERE api_name = '创建用户' AND api_path = '/api/admin/system/users' AND api_method = 'POST';
```

### 3. 动态权限控制
结合Casbin权限引擎，可实现：
- 基于API的权限验证
- 角色权限映射
- 动态权限授予

---

## ✨ 数据导入特性

✅ **完整性**：包含所有系统API接口  
✅ **规范性**：遵循RESTful命名规范  
✅ **可扩展**：支持新增API自动映射  
✅ **安全性**：支持权限绑定和访问控制  
✅ **易维护**：通过API管理页面动态调整  

---

## 🔍 验证导入

登录系统后，访问"系统管理 > API管理"页面：

1. **查看API列表**
   - 应该看到所有37个API接口
   - 按分类显示（认证、用户、角色、菜单、权限、字典、API管理）

2. **验证API信息**
   - API名称正确
   - API路径正确
   - HTTP方法正确
   - 状态为"启用"

3. **关联权限**（可选）
   - 点击API右侧的"编辑"按钮
   - 选择关联的权限
   - 保存更改

---

## 📝 文件清单

| 文件 | 说明 |
|------|------|
| `db/init.sql` | 主初始化脚本（包含API数据） |
| `db/insert_apis.sql` | 单独的API数据插入脚本 |
| `API_DATABASE_IMPORT.md` | 本文档 |

---

## 🚨 注意事项

1. **唯一约束**：API表对`(api_path, api_method)`有唯一约束，避免重复插入相同的路径和方法

2. **软删除**：API删除时使用软删除（设置deleted_at），不会物理删除数据

3. **权限关联**：permission_id字段为可选，不绑定权限时可为NULL

4. **数据库编码**：必须使用UTF-8编码，确保中文字符正确显示

---

## 📞 故障排除

### 问题1：导入失败 - 唯一约束违反

**原因**：同一个API已经存在  
**解决**：先清空apis表再导入，或检查是否重复导入
```sql
-- 清空apis表（谨慎操作）
TRUNCATE TABLE apis;
-- 或删除特定API
DELETE FROM apis WHERE api_path = '/api/admin/system/users' AND api_method = 'GET';
```

### 问题2：中文字符显示乱码

**原因**：数据库编码不是UTF-8  
**解决**：修改数据库编码
```sql
ALTER DATABASE power_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 问题3：权限关联出错

**原因**：permission_id指向的权限不存在  
**解决**：先检查权限表
```sql
-- 查看权限列表
SELECT id, name FROM permissions;
-- 清除无效的permission_id
UPDATE apis SET permission_id = NULL WHERE permission_id NOT IN (SELECT id FROM permissions);
```

---

**文档更新时间**：2025-11-29  
**API数据版本**：1.0.0  
**系统版本**：Power Admin 1.0
