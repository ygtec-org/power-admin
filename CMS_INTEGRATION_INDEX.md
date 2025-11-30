# CMS插件集成方案 - 完整索引

> 本文档汇总所有CMS集成相关的资源和指南，帮助你快速定位需要的内容。

---

## 📚 文档导航

### 1️⃣ 快速开始（10分钟）

| 文档 | 内容 | 用途 |
|------|------|------|
| **CMS_QUICK_START.md** | 6步快速实施指南 | 👨‍💻 **新手必读** - 从0到1的快速路径 |
| **CMS_IMPLEMENTATION_CHECKLIST.md** | 逐天的任务清单 | ✅ 项目管理 - 追踪开发进度 |

### 2️⃣ 架构设计（30分钟）

| 文档 | 内容 | 用途 |
|------|------|------|
| **CMS_PLUGIN_INTEGRATION_PLAN.md** | 完整架构设计方案 | 🏗️ **深度理解** - 系统设计和实现细节 |
| **CMS_ARCHITECTURE_COMPARISON.md** | 三种方案对比分析 | 🔍 **方案选择** - 了解优劣并做出决策 |

### 3️⃣ 数据库和数据（5分钟）

| 文件 | 内容 | 用途 |
|------|------|------|
| **power-admin-server/db/migrations/002_init_cms_schema.sql** | CMS数据库初始化脚本 | 🗄️ **开箱即用** - 一键创建所有表 |

---

## 🎯 场景导航

### 场景1：我想快速了解CMS插件能做什么

**推荐路径：**
```
CMS_QUICK_START.md
    ↓
CMS_ARCHITECTURE_COMPARISON.md (看"最终建议"部分)
```

**预计阅读时间**: 15分钟

---

### 场景2：我需要开始开发

**推荐路径：**
```
CMS_QUICK_START.md (第一步和第二步)
    ↓
power-admin-server/db/migrations/002_init_cms_schema.sql (执行SQL)
    ↓
CMS_PLUGIN_INTEGRATION_PLAN.md (参考具体实现)
    ↓
CMS_IMPLEMENTATION_CHECKLIST.md (逐个完成任务)
```

**预计时间**: 2周

---

### 场景3：我需要理解整个系统架构

**推荐路径：**
```
CMS_ARCHITECTURE_COMPARISON.md
    ↓
CMS_PLUGIN_INTEGRATION_PLAN.md
    ↓
实际代码实现
```

**预计阅读时间**: 1小时

---

### 场景4：我被分配到某个具体任务

**推荐路径：**
```
CMS_IMPLEMENTATION_CHECKLIST.md (定位你的任务)
    ↓
CMS_QUICK_START.md (对应的章节)
    ↓
CMS_PLUGIN_INTEGRATION_PLAN.md (技术细节)
```

**预计时间**: 取决于任务复杂度

---

## 📋 核心概念速览

### 关键术语

| 术语 | 解释 | 举例 |
|------|------|------|
| **插件** | 可独立启用/禁用的功能模块 | CMS作为插件集成到主系统 |
| **可插拔** | 支持安装/卸载而不影响主系统 | CMS禁用后主系统正常运行 |
| **菜单注入** | 动态将菜单项添加到左侧菜单栏 | CMS启用后"CMS管理"菜单出现 |
| **权限隔离** | CMS权限独立管理，不与系统权限混合 | CMS角色(cms_admin/cms_editor) |
| **数据隔离** | CMS数据表与系统表分离 | cms_content, cms_category等 |
| **API网关** | CMS通过API与主系统通信 | /api/cms/admin/contents |

### 核心流程

```
1️⃣ 用户启用CMS插件
       ↓
2️⃣ 系统初始化CMS表、权限、菜单
       ↓
3️⃣ CMS菜单出现在左侧菜单栏
       ↓
4️⃣ 用户访问CMS功能(内容、分类、用户管理)
       ↓
5️⃣ 系统检查用户CMS权限
       ↓
6️⃣ 允许访问(如果有权限) 或 返回403(如果无权限)
```

---

## 🏗️ 架构快览

### 选择的方案：集成式

```
系统登录
    ↓
认证(JWT)
    ↓
加载菜单(包括CMS菜单)
    ↓
用户可访问:
├─ 系统管理(权限/角色/用户/菜单)
└─ CMS管理(内容/分类/访客) ← 新增
```

### 为什么选择集成式

| 原因 | 说明 |
|------|------|
| ⚡ **快** | 2周快速上线，不需微服务复杂性 |
| 🔐 **安全** | 权限管理无缝集成Casbin |
| 📦 **简单** | 单一二进制部署 |
| 🚀 **性能** | 无网络延迟，内存中权限验证 |

---

## 📂 文件结构

### 后端

```
power-admin-server/
├── api/
│   └── cms.api                    ← API定义(新建)
├── internal/
│   ├── handler/cms/               ← Handler层(新建)
│   │   ├── cmscontenthandler.go
│   │   ├── cmscategoryhandler.go
│   │   ├── cmspublishhandler.go
│   │   └── cmsusermanagementhandler.go
│   ├── logic/cms/                 ← Logic层(新建)
│   │   ├── cmscontentlogic.go
│   │   ├── cmscategorylogic.go
│   │   ├── cmspublishlogic.go
│   │   └── cmsusermanagementlogic.go
│   ├── types/
│   │   └── cms.go                 ← 类型定义(新建)
│   └── svc/
│       └── servicecontext.go      ← 修改:添加CMSPlugin
├── db/
│   └── migrations/
│       └── 002_init_cms_schema.sql ← 数据库脚本(新建)
└── pkg/
    └── plugins/
        └── cms-plugin.go          ← 插件框架(新建)
```

### 前端

```
power-admin-web/
├── src/
│   ├── pages/cms/                 ← CMS页面(新建)
│   │   ├── content/
│   │   │   ├── ContentList.vue
│   │   │   ├── ContentDetail.vue
│   │   │   └── ContentForm.vue
│   │   ├── category/
│   │   │   ├── CategoryList.vue
│   │   │   └── CategoryTree.vue
│   │   ├── users/
│   │   │   ├── UserList.vue
│   │   │   └── UserForm.vue
│   │   └── CmsLayout.vue
│   ├── api/
│   │   └── cms.ts                 ← API调用(新建)
│   ├── stores/
│   │   └── cms.ts                 ← 状态管理(新建)
│   └── router/
│       └── index.ts               ← 修改:动态注册路由
```

---

## 🔧 核心实现步骤

### 后端（5天）

```
Day 1: 数据库设计
  └─ 执行 002_init_cms_schema.sql
  
Day 2: API定义和类型定义
  └─ 创建 api/cms.api
  └─ 创建 internal/types/cms.go

Day 3: Handler和Logic实现
  └─ 实现 CRUD 操作

Day 4: 插件框架和权限集成
  └─ 实现 PluginInterface
  └─ Casbin规则集成

Day 5: 路由注册和测试
  └─ 在main.go中注册CMS路由
  └─ 单元测试
```

### 前端（5天）

```
Day 1: API接口和状态管理
  └─ 创建 src/api/cms.ts
  └─ 创建 src/stores/cms.ts

Day 2: 内容管理页面
  └─ ContentList.vue
  └─ ContentForm.vue

Day 3: 分类和用户管理
  └─ CategoryList.vue
  └─ UserList.vue

Day 4: 菜单和路由集成
  └─ 动态加载CMS菜单
  └─ 注册CMS路由

Day 5: 测试和优化
  └─ 功能测试
  └─ 权限验证
  └─ 性能优化
```

---

## 📊 关键数据结构

### CMS内容 (cms_content)

```typescript
interface CmsContent {
  id: number;
  title: string;              // 文章标题
  content: string;            // 文章内容(HTML)
  excerpt: string;            // 文章摘要
  categoryId: number;         // 分类ID
  authorId: number;           // 作者ID(系统用户)
  status: 1|2|3;              // 1:草稿 2:已发布 3:已删除
  viewCount: number;          // 浏览次数
  commentCount: number;       // 评论数
  createdAt: string;
  publishedAt?: string;
}
```

### CMS分类 (cms_category)

```typescript
interface CmsCategory {
  id: number;
  name: string;
  description: string;
  parentId?: number;          // 支持多级分类
  sort: number;               // 排序号
  status: 0|1;                // 0:禁用 1:启用
}
```

### CMS前台用户 (cms_users)

```typescript
interface CmsUser {
  id: number;
  username: string;
  email: string;
  password: string;           // 加密存储
  nickname: string;
  avatar?: string;
  status: 0|1|2;              // 1:正常 0:禁用 2:封禁
  createdAt: string;
}
```

---

## 🔒 权限管理

### CMS角色

| 角色 | 权限 | 说明 |
|------|------|------|
| **cms_admin** | 所有操作 | CMS管理员，全权限 |
| **cms_editor** | C/R/U | CMS编辑，可创建/编辑/发布 |
| **cms_viewer** | R | CMS查看者，仅查看 |

### 权限规则示例

```sql
-- CMS管理员可执行所有操作
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'cms_admin', '/api/cms/admin/contents', 'GET'),
('p', 'cms_admin', '/api/cms/admin/contents', 'POST'),
('p', 'cms_admin', '/api/cms/admin/contents', 'PUT'),
('p', 'cms_admin', '/api/cms/admin/contents', 'DELETE');

-- CMS编辑可创建和编辑，但不能删除
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'cms_editor', '/api/cms/admin/contents', 'GET'),
('p', 'cms_editor', '/api/cms/admin/contents', 'POST'),
('p', 'cms_editor', '/api/cms/admin/contents', 'PUT');
```

---

## 🚀 快速部署

### 1️⃣ 初始化数据库

```bash
# 进入数据库客户端
mysql -u root -p power_admin

# 执行初始化脚本
source power-admin-server/db/migrations/002_init_cms_schema.sql;
```

### 2️⃣ 编译后端

```bash
cd power-admin-server
go build -o power-admin.exe ./cmd/api/main.go
./power-admin.exe
```

### 3️⃣ 启动前端

```bash
cd power-admin-web
npm install
npm run dev
```

### 4️⃣ 测试

访问 http://localhost:5173，登录后应该能看到左侧菜单中的"CMS管理"项。

---

## ⚠️ 常见问题

### Q1: CMS和系统用户是同一个吗？

**A**: 不是。系统有两种用户：
- **系统用户** (admin_users): 管理后台的管理员
- **CMS用户** (cms_users): CMS前台的访客

系统用户通过CMS角色(cms_admin/cms_editor)来访问CMS功能。

### Q2: CMS可以禁用吗？

**A**: 可以。更新 `plugin_status` 表中的 `enabled` 字段为0，然后重启后端服务。所有CMS菜单和功能将自动隐藏。

### Q3: 如何实现用户对CMS内容的权限控制？

**A**: 使用Casbin规则：
```go
// 检查用户是否可以访问CMS
hasPermission, _ := enforcer.Enforce(
    userId,                        // 用户ID
    "/api/cms/admin/contents",    // 资源
    "GET",                         // 操作
)
```

### Q4: 怎样支持多个插件？

**A**: 每个插件都实现相同的 `PluginInterface` 接口，然后通过插件管理器（PluginManager）统一管理。

---

## 📖 深入学习

### 推荐阅读顺序

```
1. CMS_QUICK_START.md (快速了解)
   ↓
2. CMS_ARCHITECTURE_COMPARISON.md (了解设计决策)
   ↓
3. CMS_PLUGIN_INTEGRATION_PLAN.md (深入技术细节)
   ↓
4. CMS_IMPLEMENTATION_CHECKLIST.md (按步骤实施)
   ↓
5. 开始编码!
```

### 相关技术文档

- **Go-Zero框架**: [官方文档](https://go-zero.dev/)
- **Casbin RBAC**: [官方文档](https://casbin.org/docs/get-started)
- **Vue3 + TypeScript**: [官方文档](https://vuejs.org/)
- **MySQL设计最佳实践**: [MySQL官方文档](https://dev.mysql.com/)

---

## 🤝 获取帮助

| 问题类型 | 查看文档 |
|---------|--------|
| 如何快速开始 | CMS_QUICK_START.md |
| 架构设计问题 | CMS_PLUGIN_INTEGRATION_PLAN.md |
| 方案对比 | CMS_ARCHITECTURE_COMPARISON.md |
| 任务管理 | CMS_IMPLEMENTATION_CHECKLIST.md |
| 权限配置 | CMS_PLUGIN_INTEGRATION_PLAN.md (第五部分) |
| 数据库 | power-admin-server/db/migrations/002_init_cms_schema.sql |

---

## 📞 核对清单

在开始开发之前，确保你：

- [ ] 已阅读 CMS_QUICK_START.md
- [ ] 理解了集成式架构的优势
- [ ] 了解CMS的核心功能(内容/分类/用户)
- [ ] 知道如何执行数据库脚本
- [ ] 准备好后端和前端开发环境
- [ ] 有任何问题都查阅了相应的文档

---

**最后更新**: 2024年  
**版本**: 1.0  
**状态**: ✅ 完成

