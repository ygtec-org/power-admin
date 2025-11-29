# API快速参考手册

## 📌 快速导航

### 认证 (Authentication)
```
POST   /api/admin/auth/login              用户登录
POST   /api/admin/auth/logout             用户登出
GET    /api/admin/auth/info               获取用户信息
```

### 用户 (User)
```
GET    /api/admin/system/users             获取用户列表
POST   /api/admin/system/users             创建用户
PUT    /api/admin/system/users             编辑用户
DELETE /api/admin/system/users             删除用户
GET    /api/admin/system/users/:id         获取用户详情
POST   /api/admin/system/users/:id/roles   分配角色
GET    /api/admin/system/users/:id/roles   获取用户角色
```

### 角色 (Role)
```
GET    /api/admin/system/roles                    获取角色列表
POST   /api/admin/system/roles                    创建角色
PUT    /api/admin/system/roles                    编辑角色
DELETE /api/admin/system/roles                    删除角色
GET    /api/admin/system/roles/:id                获取角色详情
POST   /api/admin/system/roles/:id/permissions    分配权限
GET    /api/admin/system/roles/:id/permissions    获取角色权限
```

### 菜单 (Menu)
```
GET    /api/admin/system/menus             获取菜单列表
POST   /api/admin/system/menus             创建菜单
PUT    /api/admin/system/menus             编辑菜单
DELETE /api/admin/system/menus             删除菜单
GET    /api/admin/system/menus/:id         获取菜单详情
```

### 权限 (Permission)
```
GET    /api/admin/system/permissions       获取权限列表
POST   /api/admin/system/permissions       创建权限
PUT    /api/admin/system/permissions       编辑权限
DELETE /api/admin/system/permissions       删除权限
GET    /api/admin/system/permissions/:id   获取权限详情
```

### 字典 (Dictionary)
```
GET    /api/admin/content/dicts            获取字典列表
POST   /api/admin/content/dicts            创建字典
PUT    /api/admin/content/dicts            编辑字典
DELETE /api/admin/content/dicts            删除字典
GET    /api/admin/content/dicts/:id        获取字典详情
```

### API管理 (API Management)
```
GET    /api/admin/system/apis              获取API列表
POST   /api/admin/system/apis              创建API记录
PUT    /api/admin/system/apis              编辑API
DELETE /api/admin/system/apis              删除API
GET    /api/admin/system/apis/:id          获取API详情
```

---

## 📊 API数据统计

```
认证相关:      3个 API
用户管理:      7个 API
角色管理:      7个 API
菜单管理:      5个 API
权限管理:      5个 API
字典管理:      5个 API
API管理:       5个 API
─────────────────
总计:         37个 API
```

---

## 🔐 认证方式

所有API（除登录/登出外）需要提供JWT Token：

```
Authorization: Bearer <your_token_here>
Content-Type: application/json
```

---

## 📋 常用查询

### 查看所有API
```sql
SELECT id, api_name, api_path, api_method, status FROM apis ORDER BY id;
```

### 查看启用的API
```sql
SELECT * FROM apis WHERE status = 1 ORDER BY api_path;
```

### 按分类统计
```sql
SELECT 
  CASE 
    WHEN api_path LIKE '%/auth/%' THEN '认证'
    WHEN api_path LIKE '%/system/users%' THEN '用户'
    WHEN api_path LIKE '%/system/roles%' THEN '角色'
    WHEN api_path LIKE '%/system/menus%' THEN '菜单'
    WHEN api_path LIKE '%/system/permissions%' THEN '权限'
    WHEN api_path LIKE '%/content/dicts%' THEN '字典'
    WHEN api_path LIKE '%/system/apis%' THEN 'API管理'
  END as category,
  COUNT(*) as count
FROM apis WHERE deleted_at IS NULL
GROUP BY category;
```

### 查看权限关联
```sql
SELECT a.api_name, a.api_path, a.api_method, p.name 
FROM apis a 
LEFT JOIN permissions p ON a.permission_id = p.id 
WHERE a.status = 1;
```

---

## 💾 导入命令

### 完整初始化
```bash
mysql -u root -p power_admin < db/init.sql
```

### 仅导入API
```bash
mysql -u root -p power_admin < db/insert_apis.sql
```

---

## ⚙️ 配置说明

| 字段 | 说明 | 示例 |
|------|------|------|
| api_name | API名称 | 用户登录 |
| api_path | API路径 | /api/admin/auth/login |
| api_method | HTTP方法 | POST |
| description | API描述 | 用户登录接口 |
| status | 启用状态 | 1(启用) / 0(禁用) |
| permission_id | 关联权限 | NULL(无权限控制) |

---

## 🎯 使用场景

### 场景1：获取用户列表
```
请求: GET /api/admin/system/users?page=1&pageSize=10
头部: Authorization: Bearer <token>
响应: { code: 0, msg: "success", data: { total: 100, list: [...] } }
```

### 场景2：创建用户
```
请求: POST /api/admin/system/users
头部: Authorization: Bearer <token>
体部: { username: "newuser", phone: "13800000000", password: "123456" }
响应: { code: 0, msg: "创建成功", data: { id: 10, username: "newuser" } }
```

### 场景3：分配角色
```
请求: POST /api/admin/system/users/10/roles
头部: Authorization: Bearer <token>
体部: { roleIds: [1, 2, 3] }
响应: { code: 0, msg: "分配成功" }
```

---

## 📱 前端集成示例

### 使用TypeScript
```typescript
import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8888',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
})

// 获取用户列表
const { data } = await api.get('/api/admin/system/users', {
  params: { page: 1, pageSize: 10 }
})

// 创建用户
await api.post('/api/admin/system/users', {
  username: 'newuser',
  phone: '13800000000',
  password: '123456'
})
```

---

## 🔔 常见问题

**Q: API数据在哪里管理?**  
A: 在系统的"API管理"页面或直接操作apis表

**Q: 如何关联权限?**  
A: 在API编辑页面选择permission_id，或直接更新apis表

**Q: 新增API如何注册?**  
A: 在admin.api定义接口，执行make gen生成代码，API数据自动插入或通过API管理页面手动添加

**Q: API路径有什么命名规范?**  
A: 遵循RESTful规范，格式为 `/api/admin/{module}/{resource}` 或 `/api/admin/{module}/{resource}/:id/{action}`

---

*最后更新: 2025-11-29*
