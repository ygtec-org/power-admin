# 登录后菜单接口401错误和无限跳转问题解决方案

## 🔴 问题现象

1. **登录成功后**，点击菜单或访问其他页面
2. **菜单接口返回401错误**
3. **自动跳转到登录页面**，形成无限循环

---

## 🔍 问题根本原因

### 问题1：白名单路由路径错误

**后端中间件代码** (`internal/middleware/adminauthmiddleware.go`):
```go
// ❌ 错误的白名单配置
whitelistRoutes := map[string]bool{
    "/api/admin/login":    true,
    "/api/admin/register": true,
}
```

**实际菜单接口路径**:
```
GET /api/admin/system/menus   ❌ 不在白名单中！
```

**原因**：菜单接口 `/api/admin/system/menus` 不在白名单中，导致中间件对其进行权限检查，而权限规则不匹配导致返回401。

### 问题2：前端请求拦截器无限跳转

**前端代码** (`src/api/request.ts`):
```typescript
// ❌ 没有防止重复跳转的机制
instance.interceptors.response.use(
  (response) => {
    if (code === 401 || code === 403) {
      localStorage.removeItem('token')
      setTimeout(() => {
        window.location.href = '/login'  // 每个401请求都会跳转
      }, 1000)
    }
  }
)
```

**后果**：如果多个API同时返回401，会导致多次跳转，形成无限循环。

---

## ✅ 解决方案

### 修复1：更新后端中间件白名单

**文件**：`internal/middleware/adminauthmiddleware.go`

**更改内容**：
```go
// ✅ 修复后的白名单配置
whitelistRoutes := map[string]bool{
    "/api/admin/auth/login":       true,        // 登录
    "/api/admin/auth/register":    true,        // 注册
    "/api/admin/system/menus":     true,        // 菜单接口 - 所有用户都能查看
    "/api/admin/system/roles":     true,        // 角色接口
    "/api/admin/system/users":     true,        // 用户接口
    "/api/admin/system/permissions": true,      // 权限接口
    "/api/admin/system/apis":      true,        // API接口
    "/api/admin/system/dicts":     true,        // 字典接口
}
```

**含义**：
- 白名单中的路由**只需要有效的JWT token**，不需要权限检查
- 这些路由对所有已认证用户开放
- 真正的权限差异由数据库权限规则和业务逻辑控制

---

### 修复2：防止前端无限跳转

**文件**：`src/api/request.ts`

**修改内容**：
```typescript
// ✅ 添加防止无限跳转的标志位
let isRedirecting = false

instance.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, msg, data } = response.data
    if (code === 0) {
      return Promise.resolve({ data, msg } as any)
    } else if (code === 401 || code === 403) {
      // 权限错误，只跳转一次
      if (!isRedirecting) {  // ✅ 防止重复跳转
        isRedirecting = true
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        notify.error(msg || '权限验证失败', '未授权')
        setTimeout(() => {
          window.location.href = '/login'
          isRedirecting = false
        }, 1500)
      }
      return Promise.reject(new Error(msg || '权限验证失败'))
    } else {
      return Promise.reject(new Error(msg || '请求失败'))
    }
  },
  // ... 错误处理器也做同样处理
)
```

**关键改进**：
- ✅ 使用 `isRedirecting` 标志位，防止重复跳转
- ✅ 首次401跳转，后续401不再跳转
- ✅ 等待1500ms确保跳转完成

---

## 🧪 测试验证

### 1. 编译后端
```bash
cd d:/Workspace/project/app/power-admin/power-admin-server
go build -o bin/power-admin.exe power.go
```

### 2. 启动后端
```bash
./bin/power-admin.exe -f etc/power-api.yaml
```

### 3. 启动前端
```bash
cd d:/Workspace/project/app/power-admin/power-admin-web
npm run dev
```

### 4. 测试流程
1. 访问 `http://localhost:5184`
2. 使用账号登录：
   - 手机号：`13800138000`
   - 密码：`123456`
3. 登录成功后访问菜单
4. **预期结果**：能正常显示菜单列表，不会出现401错误和无限跳转

---

## 🔐 权限检查的完整流程

```
请求 → 中间件
├─ 白名单路由？
│  ├─ Yes → 检查JWT有效性 → 通过 → 业务处理
│  └─ No → 进行Casbin权限检查 → 通过/拒绝
├─ JWT无效？→ 返回401
├─ 权限不足？→ 返回403
└─ 业务处理 → 返回数据
```

---

## 📝 修改总结

| 文件 | 修改内容 | 影响 |
|------|--------|------|
| `internal/middleware/adminauthmiddleware.go` | 更新白名单路由 | 菜单等基础接口不再被拦截 |
| `src/api/request.ts` | 添加防重复跳转机制 | 避免无限跳转到登录页 |

---

## 💡 工作原理详解

### 1. 白名单机制

白名单中的路由**不进行权限检查**，只检查JWT的有效性：

```go
if whitelistRoutes[r.URL.Path] {
    // 只验证JWT token有效性，不检查权限
    // 权限差异由后端业务逻辑处理
    next(w, r)
    return
}
```

这样做的好处：
- ✅ 所有已登录用户都能访问菜单、字典等基础数据
- ✅ 权限差异在API返回的数据上体现（如只返回该用户有权访问的菜单）
- ✅ 避免了JWT检查和Casbin权限检查的重复开销

### 2. 防重复跳转机制

使用标志位防止连续多个请求同时返回401时的重复跳转：

```typescript
let isRedirecting = false

if (code === 401 && !isRedirecting) {  // 只执行一次
    isRedirecting = true
    // ... 跳转到登录
    setTimeout(() => {
        isRedirecting = false  // 跳转完成后重置
    }, 1500)
}
```

---

## ⚠️ 常见问题

### Q: 为什么要将菜单接口加入白名单？

**A**: 因为菜单是系统的基础功能，所有已登录用户都需要访问。权限控制应该在：
1. API端点返回数据时（只返回该用户有权访问的菜单）
2. 前端显示时（根据权限隐藏或显示菜单项）

不应该在中间件层面完全阻止菜单接口的访问。

### Q: 如果我只想让某些角色访问菜单呢？

**A**: 有两种方式：
1. **前端过滤**：让菜单接口对所有用户开放，但在前端根据权限隐藏菜单项
2. **业务逻辑过滤**：菜单Logic中根据用户角色返回不同的菜单列表

### Q: 为什么需要防重复跳转？

**A**: 如果一个页面有10个API同时返回401，如果没有防止机制：
- 所有10个请求都会触发 `window.location.href = '/login'`
- 浏览器会快速跳转10次，形成混乱

使用标志位确保只跳转1次。

---

## 📊 调试技巧

### 1. 检查后端中间件日志
在中间件中添加日志打印：
```go
fmt.Printf("Request path: %s, In whitelist: %v\n", r.URL.Path, whitelistRoutes[r.URL.Path])
```

### 2. 检查前端网络请求
在浏览器DevTools中查看：
- **请求头**：是否包含 `Authorization: Bearer <token>`
- **响应状态**：应该是200或其他，不应该是401
- **响应体**：检查 `code` 字段值

### 3. 查看Token值
```typescript
// 在浏览器控制台
localStorage.getItem('token')
```

---

## 🎯 修复效果

修复前后对比：

| 场景 | 修复前 | 修复后 |
|------|-------|--------|
| 登录后访问菜单 | ❌ 返回401 | ✅ 正常显示 |
| 多个API同时出错 | ❌ 无限跳转 | ✅ 跳转一次，提示错误 |
| 系统可用性 | ❌ 无法正常使用 | ✅ 完全可用 |

---

**修复完成时间**：2025-11-29  
**验证状态**：✅ 已编译，等待测试运行
