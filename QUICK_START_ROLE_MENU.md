# 角色菜单权限快速开始

## 5分钟快速指南

### 1. 编译并启动服务器（30秒）

```bash
cd power-admin-server
go build -o power-admin.exe ./power.go
./power-admin.exe
```

服务器运行在 `http://localhost:8888`

### 2. 测试菜单权限（2分钟）

#### 步骤1: 登录获取Token
```bash
curl -X POST http://localhost:8888/api/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'
```

返回示例：
```json
{
  "code": 200,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### 步骤2: 请求菜单列表
```bash
curl -X GET http://localhost:8888/api/admin/system/menus \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

返回示例：
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
      "icon": "setting",
      "sort": 1,
      "status": 1,
      "menuType": 1,
      "createdAt": "2025-11-30 11:00:00",
      "children": [...]
    }
  ],
  "total": 5
}
```

### 3. 前端集成（3分钟）

#### 步骤1: 在API文件中添加菜单获取函数

创建 `src/api/menu.ts`:
```typescript
import axios from 'axios'

export async function getMenuList() {
  const response = await axios.get('/api/admin/system/menus', {
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  })
  return response.data.data || []
}
```

#### 步骤2: 在应用启动时加载菜单

修改 `src/main.ts`:
```typescript
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { getMenuList } from '@/api/menu'

const app = createApp(App)
const token = localStorage.getItem('token')

if (token) {
  // 获取菜单并创建路由
  getMenuList().then(menus => {
    const routes = menus.map(menu => ({
      path: menu.menuPath,
      component: () => import(`@/pages/${menu.component}.vue`)
    }))
    
    const router = createRouter({
      history: createWebHistory(),
      routes
    })
    
    app.use(router).mount('#app')
  })
} else {
  app.mount('#app')
}
```

## 完整示例

### 创建测试数据

```sql
-- 1. 创建用户
INSERT INTO users (username, password, nickname) 
VALUES ('testuser', 'hashed_password', 'Test User');

-- 2. 创建菜单
INSERT INTO menus (parent_id, menu_name, menu_path, component, icon, sort, status) 
VALUES 
  (0, '系统管理', '/system', 'Layout', 'setting', 1, 1),
  (1, '用户管理', '/system/users', 'system/User', 'user', 1, 1),
  (1, '菜单管理', '/system/menus', 'system/Menu', 'menu', 2, 1),
  (1, '角色管理', '/system/roles', 'system/Role', 'role', 3, 1);

-- 3. 创建角色
INSERT INTO roles (name, description) 
VALUES 
  ('Admin', '管理员'),
  ('User', '普通用户');

-- 4. 分配用户角色
INSERT INTO user_roles (user_id, role_id) 
VALUES (1, 1); -- 给testuser分配Admin角色

-- 5. 绑定角色菜单
INSERT INTO role_menus (role_id, menu_id) 
VALUES 
  (1, 1), -- Admin可以看到系统管理
  (1, 2), -- Admin可以看到用户管理
  (1, 3), -- Admin可以看到菜单管理
  (1, 4), -- Admin可以看到角色管理
  (2, 1), -- 普通用户只能看到系统管理菜单项
  (2, 2); -- 普通用户可以看到用户管理
```

### 验证权限

```bash
# 用Admin账户登录，应该看到4个菜单
curl -X POST http://localhost:8888/api/admin/auth/login \
  -d '{"username":"admin","password":"admin"}'

# 用testuser账户登录，应该看到2个菜单（1和2）
curl -X POST http://localhost:8888/api/admin/auth/login \
  -d '{"username":"testuser","password":"password"}'
```

## 主要修改文件

| 文件 | 修改内容 |
|------|--------|
| `internal/logic/menu/menulistlogic.go` | 添加基于角色的菜单过滤逻辑 |
| `pkg/repository/menu.go` | 新增`GetMenusByRoleIDs()`方法 |
| `pkg/repository/role.go` | 新增`GetRolesByUserID()`方法 |
| `internal/middleware/adminauthmiddleware.go` | 将用户ID存储到Context |

## 工作流程

```
用户请求菜单列表
    ↓
Middleware验证JWT Token
    ↓
提取用户ID放入Context
    ↓
MenuListLogic处理请求
    ├─ 从Context获取用户ID
    ├─ 查询用户的所有角色
    ├─ 查询这些角色绑定的菜单
    └─ 返回菜单树形结构
    ↓
前端接收菜单数据
    ├─ 动态生成路由
    ├─ 构建导航菜单
    └─ 用户可访问授权的菜单
```

## 常见问题

**Q: 为什么看不到任何菜单？**
A: 检查user_roles和role_menus表是否有对应的绑定记录

**Q: 如何禁用某个菜单？**
A: 更新menus表的status字段为0：`UPDATE menus SET status=0 WHERE id=xxx`

**Q: 如何给用户分配菜单？**
A: 1) 先给用户分配角色 2) 再给角色绑定菜单

**Q: API权限如何控制？**
A: API权限通过Casbin RBAC独立管理，与菜单权限分开

## 下一步

1. 阅读详细文档
   - `ROLE_BASED_MENU_IMPLEMENTATION.md` - 实现细节
   - `FRONTEND_MENU_INTEGRATION.md` - 前端集成

2. 运行测试
   - 参考 `MENU_PERMISSION_TESTING.md`

3. 部署上线
   - 确保数据库schema已更新
   - 创建必要的菜单和角色数据
   - 配置权限绑定

## 支持

需要帮助？查看以下文档：
- 错误排查：见MENU_PERMISSION_TESTING.md中的"常见问题排查"
- API文档：见MENU_PERMISSION_TESTING.md中的"API调用流程"
- 数据库：见ROLE_BASED_MENU_IMPLEMENTATION.md中的"数据库表关系"

---

**提示：** 保持菜单数据和权限配置的一致性。前端菜单受本实现控制，而API调用权限受Casbin RBAC控制。
