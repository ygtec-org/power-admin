# Power Admin - 完整项目结构

本文档详细说明 Power Admin 项目的完整文件结构和各文件的用途。

## 📁 根目录结构

```
power-admin/
├── power-admin-server/           # 后端 API 服务（Go-Zero）
├── power-admin-web/              # 前端管理台（Vue 3）
├── 🚀 init.bat                   # Windows 一键初始化脚本
├── 🚀 init.sh                    # Linux/Mac 初始化脚本
├── 📄 casbin_rule.sql            # Casbin 规则表 SQL（参考用）
│
├── 📖 README.md                  # 项目主文档
├── 📖 QUICK_START_CN.md          # 中文快速开始指南 ⭐
├── 📖 QUICKSTART.md              # 英文快速开始指南
├── 📖 RUN_GUIDE.md               # 详细运行指南
├── 📖 COMPLETION_REPORT.md       # 完成报告 ⭐
├── 📖 OPTIMIZATION_SUMMARY.md    # 优化总结
├── 📖 DEVELOPMENT_GUIDE.md       # 开发指南
├── 📖 ARCHITECTURE.md            # 架构设计文档
├── 📖 PROJECT_SUMMARY.md         # 项目概览
└── 📖 PROJECT_STRUCTURE.md       # 本文件
```

---

## 📁 后端项目结构 (power-admin-server/)

```
power-admin-server/
├── 🚀 go.mod                     # Go 模块定义
├── 🚀 go.sum                     # 依赖版本锁定
├── 🚀 power.go                   # 主程序入口
│
├── api/                          # API 定义（Protocol Buffer）
│   ├── power.proto               # API 协议定义
│   └── ...
│
├── db/                           # 数据库相关
│   └── init.sql                  # ⭐ 数据库初始化脚本
│       ├── 10 个表定义
│       ├── 3 个初始角色
│       ├── 20 个初始权限
│       └── Casbin 规则表
│
├── etc/                          # 配置文件目录
│   ├── power-api.yaml            # ⭐ 应用配置文件
│   │   ├── 服务端口（8888）
│   │   ├── MySQL 连接信息
│   │   ├── Redis 连接信息
│   │   └── JWT 密钥
│   │
│   └── rbac_model.conf           # ⭐ Casbin RBAC 模型配置
│       ├── request_definition
│       ├── policy_definition
│       ├── role_definition
│       ├── policy_effect
│       └── matchers
│
├── pkg/                          # 业务包
│   ├── auth/                     # 认证模块
│   │   ├── jwt.go                # JWT 令牌生成和验证
│   │   └── password.go           # 密码加密（bcrypt）
│   │
│   ├── db/                       # 数据库配置
│   │   └── db.go                 # GORM 初始化、数据源配置
│   │
│   ├── models/                   # 数据模型
│   │   ├── user.go               # 用户模型
│   │   ├── role.go               # 角色模型
│   │   ├── menu.go               # 菜单模型
│   │   ├── permission.go         # 权限模型
│   │   ├── api.go                # API 管理模型
│   │   └── dictionary.go         # 字典模型
│   │
│   ├── permission/               # 权限管理
│   │   └── rbac.go               # ⭐ Casbin RBAC 实现
│   │       ├── NewRBACEnforcer() - 初始化
│   │       ├── CheckPermission() - 权限检查
│   │       ├── AddRoleForUser() - 分配角色
│   │       └── ... 其他权限方法
│   │
│   └── repository/               # 数据仓储层
│       ├── user.go               # 用户仓储
│       ├── role.go               # 角色仓储
│       ├── menu.go               # 菜单仓储
│       ├── permission.go         # 权限仓储
│       ├── api.go                # API 仓储
│       └── dictionary.go         # 字典仓储
│
└── internal/                     # 内部实现
    ├── config/                   # 配置管理
    │   └── config.go             # 配置结构体
    │
    ├── handler/                  # HTTP 处理器
    │   ├── auth_handler.go        # 登录/注册处理
    │   ├── user_handler.go        # 用户管理处理
    │   ├── role_handler.go        # 角色管理处理
    │   ├── menu_handler.go        # 菜单管理处理
    │   ├── permission_handler.go  # 权限管理处理
    │   ├── api_handler.go         # API 管理处理
    │   └── dict_handler.go        # 字典管理处理
    │
    ├── logic/                    # 业务逻辑层
    │   ├── auth_logic.go          # 认证逻辑
    │   ├── user_logic.go          # 用户逻辑
    │   ├── role_logic.go          # 角色逻辑
    │   ├── menu_logic.go          # 菜单逻辑
    │   ├── permission_logic.go    # 权限逻辑
    │   ├── api_logic.go           # API 逻辑
    │   └── dict_logic.go          # 字典逻辑
    │
    ├── middleware/               # 中间件
    │   ├── auth.go                # ⭐ JWT 认证中间件
    │   ├── cors.go                # CORS 中间件
    │   └── logger.go              # 日志中间件
    │
    ├── router/                   # 路由定义
    │   └── routes.go              # API 路由注册
    │
    └── svc/                      # 服务上下文
        └── servicecontext.go      # ⭐ 依赖注入容器
            ├── DB 实例
            ├── Redis 实例
            ├── 各仓储实例
            └── 权限管理器实例
```

### 后端数据库表（10 个）

| 表名 | 说明 |
|-----|------|
| `users` | 用户表 |
| `roles` | 角色表 |
| `permissions` | 权限表 |
| `menus` | 菜单表（树形结构） |
| `user_roles` | 用户-角色关联表 |
| `role_permissions` | 角色-权限关联表 |
| `role_menus` | 角色-菜单关联表 |
| `dictionaries` | 字典数据表 |
| `apis` | API 管理表 |
| `casbin_rule` | Casbin 规则表 |

---

## 📁 前端项目结构 (power-admin-web/)

```
power-admin-web/
├── 🚀 package.json               # NPM 依赖定义
├── 🚀 package-lock.json          # 依赖版本锁定
├── 🚀 vite.config.ts             # Vite 构建配置
├── 🚀 tsconfig.json              # TypeScript 配置
├── 🚀 index.html                 # HTML 入口文件
│
├── src/                          # 源代码目录
│   ├── 🎯 main.ts                # ⭐ 程序入口
│   ├── 🎯 App.vue                # ⭐ 根组件
│   │
│   ├── api/                      # API 调用模块
│   │   ├── request.ts            # ⭐ Axios 实例 + 拦截器
│   │   │   ├── 请求拦截器（token）
│   │   │   ├── 响应拦截器（错误处理）
│   │   │   └── 自动登录重定向
│   │   │
│   │   ├── user.ts               # ⭐ 用户 API
│   │   │   ├── login() - 登录
│   │   │   ├── logout() - 登出
│   │   │   ├── getUsers() - 获取用户列表
│   │   │   ├── createUser() - 创建用户
│   │   │   ├── updateUser() - 更新用户
│   │   │   └── deleteUser() - 删除用户
│   │   │
│   │   ├── role.ts               # 角色 API (待实现)
│   │   ├── menu.ts               # 菜单 API (待实现)
│   │   ├── permission.ts         # 权限 API (待实现)
│   │   └── dict.ts               # 字典 API (待实现)
│   │
│   ├── router/                   # 路由配置
│   │   └── index.ts              # ⭐ 完整路由定义
│   │       ├── 登录路由 (/login)
│   │       ├── 仪表板路由 (/dashboard)
│   │       ├── 系统管理路由 (/system/*)
│   │       ├── 内容管理路由 (/content/*)
│   │       └── 路由守卫（权限检查）
│   │
│   ├── stores/                   # Pinia 状态管理 (预留)
│   │   └── user.ts               # 用户状态存储 (待实现)
│   │
│   ├── layout/                   # 布局组件
│   │   └── Layout.vue            # ⭐ 管理台主布局
│   │       ├── Header（导航栏）
│   │       ├── Sidebar（左侧菜单）
│   │       ├── Content（主内容区）
│   │       └── RouterView（页面切换）
│   │
│   ├── pages/                    # 页面组件
│   │   ├── Login.vue             # ⭐ 登录页面
│   │   │   ├── 手机号输入
│   │   │   ├── 密码输入
│   │   │   ├── 登录按钮
│   │   │   └── 错误提示
│   │   │
│   │   ├── Dashboard.vue         # ⭐ 仪表板
│   │   │   ├── 统计卡片（用户、角色、菜单、API）
│   │   │   ├── 功能介绍卡片
│   │   │   └── 动画效果
│   │   │
│   │   ├── system/               # 系统管理模块
│   │   │   ├── user/
│   │   │   │   └── UserList.vue  # ⭐ 用户管理列表
│   │   │   │       ├── 用户列表表格
│   │   │   │       ├── 编辑按钮
│   │   │   │       ├── 删除按钮
│   │   │   │       └── 新增按钮
│   │   │   │
│   │   │   ├── role/
│   │   │   │   └── RoleList.vue  # 角色管理 (占位)
│   │   │   ├── menu/
│   │   │   │   └── MenuList.vue  # 菜单管理 (占位)
│   │   │   ├── permission/
│   │   │   │   └── PermissionList.vue  # 权限管理 (占位)
│   │   │   └── api/
│   │   │       └── ApiList.vue   # API 管理 (占位)
│   │   │
│   │   └── content/              # 内容管理模块
│   │       └── dict/
│   │           └── DictList.vue  # 字典管理 (占位)
│   │
│   ├── components/               # 可复用组件
│   │   ├── Table.vue             # 通用表格组件
│   │   └── ...
│   │
│   └── styles/ (待补充)          # 全局样式
│       └── main.css              # 全局样式表
```

### 前端依赖清单

```json
{
  "vue": "^3.3.4",                 # Vue 3 框架
  "vue-router": "^4.2.5",          # 路由管理
  "axios": "^1.5.0",               # HTTP 客户端
  "pinia": "^2.1.6",               # 状态管理
  "element-plus": "^2.4.1",        # UI 组件库
  "dayjs": "^1.11.10"              # 日期时间库
}
```

---

## 🔗 关键文件和职责

### 后端关键文件

| 文件 | 职责 | 修改频率 |
|-----|------|--------|
| `etc/power-api.yaml` | ⭐⭐ 应用配置（端口、数据库） | 低 |
| `etc/rbac_model.conf` | ⭐⭐ Casbin 权限模型 | 低 |
| `pkg/permission/rbac.go` | ⭐⭐ RBAC 权限管理实现 | 低 |
| `internal/svc/servicecontext.go` | ⭐ 依赖注入容器 | 低 |
| `internal/handler/*.go` | ⭐⭐ HTTP 处理器 | 高 |
| `internal/logic/*.go` | ⭐⭐ 业务逻辑 | 高 |
| `pkg/repository/*.go` | ⭐ 数据访问层 | 中 |
| `db/init.sql` | ⭐ 数据库初始化 | 低 |

### 前端关键文件

| 文件 | 职责 | 修改频率 |
|-----|------|--------|
| `src/main.ts` | ⭐ 应用入口 | 低 |
| `src/router/index.ts` | ⭐⭐ 路由配置 | 中 |
| `src/api/request.ts` | ⭐ API 基础配置 | 低 |
| `src/layout/Layout.vue` | ⭐ 主布局组件 | 低 |
| `src/pages/**/*.vue` | ⭐⭐ 页面组件 | 高 |
| `src/api/*.ts` | ⭐ API 调用 | 高 |

---

## 🔄 数据流向示意

```
┌─────────────────────────────────────────────────┐
│                  浏览器                          │
│        ┌─────────────────────────────┐           │
│        │     Vue 3 + TypeScript       │           │
│        ├─────────────────────────────┤           │
│        │ • Pages (Dashboard, Users)   │           │
│        │ • Components (Layout, Table) │           │
│        │ • Router (权限守卫)          │           │
│        │ • API (request + 拦截器)    │           │
│        └────────────┬────────────────┘           │
│                     │ HTTP                       │
└─────────────────────┼───────────────────────────┘
                      │
┌─────────────────────┼───────────────────────────┐
│                     ▼                           │
│              localhost:8888                    │
│  ┌───────────────────────────────────────────┐ │
│  │         Go-Zero API Server                 │ │
│  ├───────────────────────────────────────────┤ │
│  │ • Handler (请求处理)                      │ │
│  │ • Middleware (JWT、CORS)                  │ │
│  │ • Logic (业务逻辑)                        │ │
│  │ • Repository (数据访问)                   │ │
│  │ • Permission (Casbin RBAC)                │ │
│  └────────┬────────────────────────┬─────────┘ │
│           │                        │           │
│           ▼ SQL                     ▼ Redis    │
│      ┌────────┐                ┌────────┐      │
│      │ MySQL  │                │ Redis  │      │
│      │Database│                │Cache  │      │
│      └────────┘                └────────┘      │
└────────────────────────────────────────────────┘
```

---

## 📦 构建和部署

### 后端构建

```bash
cd power-admin-server
go build -o power-admin.exe
```

输出：`power-admin.exe` （可执行文件）

### 前端构建

```bash
cd power-admin-web
npm run build
```

输出：`dist/` 目录（生产版本）

---

## 🎯 开发工作流

1. **修改页面** → `src/pages/` 修改 Vue 文件 → 自动热重载
2. **修改 API 调用** → `src/api/` 修改方法 → 测试 API 响应
3. **修改后端逻辑** → `internal/logic/` 修改逻辑 → 重新编译运行
4. **修改数据库** → `db/init.sql` 修改 SQL → 重新初始化数据库

---

## 📊 项目规模统计

| 指标 | 数值 |
|-----|------|
| **后端代码行数** | ~2500 行 |
| **前端代码行数** | ~1500 行 |
| **数据库表数** | 10 个 |
| **API 端点数** | 30+ 个 |
| **前端页面数** | 8+ 个 |
| **配置文件** | 2 个（YAML + CONF） |

---

## 🚀 快速导航

- **快速开始** → [QUICK_START_CN.md](./QUICK_START_CN.md)
- **详细运行** → [RUN_GUIDE.md](./RUN_GUIDE.md)
- **开发指南** → [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md)
- **架构设计** → [ARCHITECTURE.md](./ARCHITECTURE.md)
- **完成报告** → [COMPLETION_REPORT.md](./COMPLETION_REPORT.md)

---

**此文档最后更新**: 2025-11-29
