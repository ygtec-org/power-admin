# Power Admin - RBAC权限管理系统完整实现指南

## 📋 目录

1. [系统架构](#系统架构)
2. [已实现功能](#已实现功能)
3. [核心组件](#核心组件)
4. [集成步骤](#集成步骤)
5. [API接口](#api接口)
6. [前端权限控制](#前端权限控制)
7. [测试验证](#测试验证)

---

## 系统架构

### RBAC（Role-Based Access Control）权限模型

```
用户 → 用户-角色关联 → 角色 → 角色-权限关联 → 权限 → API端点
```

### 三层权限验证

```
1. JWT认证层（检查用户身份是否有效）
   ↓
2. Casbin权限验证层（检查用户是否有该API的访问权限）
   ↓
3. 业务逻辑层（检查数据权限）
```

---

## 已实现功能

### 1. ✅ 后端权限验证中间件增强

**文件**: `internal/middleware/adminauthmiddleware.go`

**功能**:
- JWT Token 验证
- Casbin 权限检查
- 白名单路由支持（登录、注册等）
- 权限不足返回 403 Forbidden

**关键代码**:
```go
// 检查权限（subject, object, action）
if !m.RBAC.CheckPermission(subject, resource, action) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusForbidden)
    fmt.Fprintf(w, `{"code":403,"msg":"permission denied"}`)
    return
}
```

### 2. ✅ 角色权限分配功能

**后端Logic**:
- `internal/logic/role/assignpermissionslogic.go` - 为角色分配权限
- `internal/logic/role/getrolepermissionslogic.go` - 获取角色权限列表

**前端API**:
- `src/api/role.ts`:
  - `assignPermissions(roleId, permissionIds)` - 分配权限
  - `getRolePermissions(roleId)` - 获取权限

**调用方式**:
```typescript
// 获取角色的权限
const res = await getRolePermissions(roleId)
const selectedPermissions = res.data.data.map((p) => p.id)

// 为角色分配权限
await assignPermissions(roleId, selectedPermissions)
```

### 3. ✅ 用户角色分配功能

**后端Logic**:
- `internal/logic/user/assignrolestouserlogic.go` - 为用户分配角色
- `internal/logic/user/getuserroleslogic.go` - 获取用户角色列表

**前端API**:
- `src/api/user.ts`:
  - `assignRolesToUser(userId, roleIds)` - 分配角色
  - `getUserRoles(userId)` - 获取角色

**调用方式**:
```typescript
// 获取用户的角色
const res = await getUserRoles(userId)
const selectedRoles = res.data.data.map((r) => r.id)

// 为用户分配角色
await assignRolesToUser(userId, selectedRoles)
```

### 4. ✅ Casbin RBAC权限引擎

**文件**: `pkg/permission/rbac.go`

**核心方法**:
```go
// 检查权限
CheckPermission(subject, object, action string) bool

// 为用户添加角色
AddRoleForUser(user, role string) error

// 为角色添加权限
AddPermissionForRole(role, object, action string) error

// 获取用户所有角色
GetRolesForUser(user string) ([]string, error)

// 获取角色所有权限
GetPermissionsForRole(role string) ([][]string, error)
```

---

## 核心组件

### 1. 数据模型更新

**RoleRepository** - `pkg/repository/role.go`
```go
// 新增方法
RemoveAllPermissions(roleID int64) error  // 移除角色的所有权限
```

**UserRepository** - `pkg/repository/user.go`
```go
// 新增方法
GetRoles(userID int64) ([]*models.Role, error)        // 获取用户角色
RemoveAllRoles(userID int64) error                      // 移除用户的所有角色
```

### 2. 请求/响应类型

**internal/types/types.go**:
```go
// 角色权限分配
type AssignPermissionsReq struct {
    RoleID        int64   `json:"roleId" path:"id"`
    PermissionIds []int64 `json:"permissionIds"`
}

type AssignPermissionsResp struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

// 用户角色分配
type AssignRolesToUserReq struct {
    UserID  int64   `json:"userId" path:"id"`
    RoleIds []int64 `json:"roleIds"`
}

type AssignRolesToUserResp struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}
```

---

## 集成步骤

### Step 1: 编译后端

```bash
cd d:/Workspace/project/app/power-admin/power-admin-server
go build -o bin/power-admin.exe ./cmd/admin
```

### Step 2: 启动后端服务

```bash
# Windows
.\bin\power-admin.exe -f etc/power-api.yaml

# Linux/Mac
./bin/power-admin ./etc/power-api.yaml
```

### Step 3: 启动前端开发服务

```bash
cd d:/Workspace/project/app/power-admin/power-admin-web
npm install  # 首次需要
npm run dev
```

### Step 4: 初始化数据库权限规则

使用 Seeder 工具初始化 Casbin 规则:

```bash
cd d:/Workspace/project/app/power-admin/power-admin-server
go build -o bin/seeder.exe ./cmd/seeder
./bin/seeder.exe -f etc/power-api.yaml
```

---

## API接口

### 角色权限管理

| 方法 | 端点 | 描述 |
|------|------|------|
| POST | `/api/admin/system/roles/{id}/permissions` | 为角色分配权限 |
| GET | `/api/admin/system/roles/{id}/permissions` | 获取角色权限列表 |

### 用户角色管理

| 方法 | 端点 | 描述 |
|------|------|------|
| POST | `/api/admin/system/users/{id}/roles` | 为用户分配角色 |
| GET | `/api/admin/system/users/{id}/roles` | 获取用户角色列表 |

### 权限验证流程

1. **请求** → 包含 `Authorization: Bearer <token>` 头
2. **JWT验证** → 验证Token有效性，提取用户ID
3. **Casbin检查** → 检查 `(userId, apiPath, httpMethod)` 三元组
4. **允许/拒绝** → 200 OK 或 403 Forbidden

---

## 前端权限控制

### 1. 角色管理页面 - 权限分配

**文件**: `src/pages/system/role/RoleList.vue`

```vue
<!-- 权限分配对话框 -->
<div v-if="showPermDialog" class="modal">
  <div class="modal-content" style="width: 500px">
    <div class="modal-header">
      <h2>为「{{ selectedRole?.name }}」分配权限</h2>
    </div>
    <div class="modal-body">
      <div class="permission-list">
        <div v-for="perm in permissions" :key="perm.id" class="permission-item">
          <input
            type="checkbox"
            :id="'perm-' + perm.id"
            v-model="selectedPermissions"
            :value="perm.id"
          />
          <label :for="'perm-' + perm.id">
            {{ perm.name }} ({{ perm.resource }}/{{ perm.action }})
          </label>
        </div>
      </div>
    </div>
    <div class="modal-footer">
      <button @click="showPermDialog = false" class="btn-cancel">取消</button>
      <button @click="handleSavePermissions" class="btn-primary">保存</button>
    </div>
  </div>
</div>
```

**关键函数**:
```typescript
// 查看权限
const handleViewPermissions = async (role) => {
  selectedRole.value = role
  selectedPermissions.value = []
  try {
    const res = await getRolePermissions(role.id)
    if (res.data && res.data.data) {
      selectedPermissions.value = res.data.data.map((p) => p.id)
    }
  } catch (error) {
    console.log('获取权限列表失败:', error.message)
  }
  showPermDialog.value = true
}

// 保存权限
const handleSavePermissions = async () => {
  if (!selectedRole.value) return
  try {
    await assignPermissions(selectedRole.value.id, selectedPermissions.value)
    notify.success('权限分配成功')
    showPermDialog.value = false
  } catch (error) {
    notify.error(error.message || '权限分配失败')
  }
}
```

### 2. 用户管理页面 - 角色绑定（TODO）

**待实现功能**:
- 在编辑用户时添加角色选择
- 支持多角色绑定
- 角色变更自动刷新权限

---

## 测试验证

### 1. 测试权限分配流程

```bash
# 登录系统
POST /api/admin/auth/login
{
  "phone": "13800138000",
  "password": "123456"
}

# 获取角色权限
GET /api/admin/system/roles/1/permissions

# 分配权限给角色
POST /api/admin/system/roles/1/permissions
{
  "permissionIds": [1, 2, 3, 4, 5]
}

# 为用户分配角色
POST /api/admin/system/users/1/roles
{
  "roleIds": [1, 2]
}

# 获取用户角色
GET /api/admin/system/users/1/roles

# 验证权限 - 应该返回 200 OK
GET /api/admin/system/menus

# 验证权限不足 - 应该返回 403 Forbidden
DELETE /api/admin/system/users/2  # 假设该用户无此权限
```

### 2. 权限矩阵

| 用户ID | 角色 | 权限范围 | 示例API |
|-------|------|--------|--------|
| 1 | admin | 所有权限 | GET /api/admin/system/* |
| 2 | editor | 内容管理 | GET /api/admin/system/menus |
| 3 | user | 基础权限 | GET /api/admin/user/info |

### 3. Casbin规则示例

```sql
-- 管理员角色(ID=1)拥有所有权限
('p', '1', '/api/admin/system/*', 'GET'),
('p', '1', '/api/admin/system/*', 'POST'),
('p', '1', '/api/admin/system/*', 'PUT'),
('p', '1', '/api/admin/system/*', 'DELETE'),

-- 编辑角色(ID=2)权限
('p', '2', '/api/admin/system/menus', 'GET'),
('p', '2', '/api/admin/system/dicts', 'GET'),

-- 用户角色(ID=3)基础权限
('p', '3', '/api/admin/user/info', 'GET'),

-- 用户角色关联
('g', 'admin', '1'),
('g', 'editor', '2'),
('g', 'user', '3'),
```

---

## 常见问题

### Q1: 权限验证返回403，但用户确实有权限？

**A**: 检查以下几点:
1. 确认Casbin规则表(casbin_rule)中有该权限定义
2. 检查用户是否被分配了包含该权限的角色
3. 查看AdminAuthMiddleware中白名单设置是否正确

### Q2: 如何修改权限规则？

**A**: 有两种方式:
1. **动态修改**: 在业务逻辑中调用 `svcCtx.Permission.AddPermissionForRole()`
2. **静态修改**: 修改 `db/init.sql` 中的 `casbin_rule` 表数据，重新初始化数据库

### Q3: 如何实现前端权限控制？

**A**: 可以实现权限指令 `v-permission`:

```vue
<!-- 只有拥有 'user:delete' 权限的用户才能看到删除按钮 -->
<button v-permission="'user:delete'" @click="deleteUser">删除</button>
```

---

## 下一步优化

1. **权限缓存** - 使用Redis缓存Casbin规则，提高性能
2. **权限指令** - 实现 `v-permission` 前端权限控制指令
3. **权限管理界面** - 完善权限管理和权限分配的可视化界面
4. **审计日志** - 记录所有权限变更操作
5. **权限继承** - 支持权限继承和权限模板

---

## 参考资源

- [Casbin官方文档](https://casbin.org/zh/docs/rbac/)
- [Go-Zero框架文档](https://go-zero.dev/)
- [RBAC权限模型详解](https://en.wikipedia.org/wiki/Role-based_access_control)

---

**最后更新**: 2025-11-29
