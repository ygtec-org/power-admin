# CMS插件集成方案文档

## 一、架构设计概览

### 1.1 系统整体架构

```
┌─────────────────────────────────────────────────────────┐
│                    Power Admin 主系统                      │
│  ┌───────────────────────────────────────────────────┐  │
│  │           Admin Dashboard (Web/管理后台)           │  │
│  │  - 用户管理 / 角色管理 / 权限管理 / 菜单管理        │  │
│  │  - CMS菜单项（动态注入）                            │  │
│  └───────────────────────────────────────────────────┘  │
│                          │                                │
│  ┌───────────┬───────────┼───────────┬──────────────┐  │
│  │           │           │           │              │  │
│  ▼           ▼           ▼           ▼              ▼  │
│ System    Content      Content    Permission      CMS   │
│ Module    Module       Publish    Module         Module │
│           ┌─────────────────────────────────────────┐  │
│           │                                         │  │
│           │    ┌──────────────────────────┐        │  │
│           │    │  CMS Plugin API Gateway  │        │  │
│           │    │  (内容管理 API 聚合层)    │        │  │
│           │    └──────────────────────────┘        │  │
│           │              │                          │  │
│           │    ┌─────────┴──────────┐             │  │
│           │    │                    │             │  │
│           │    ▼                    ▼             │  │
│           │  CMS Content DB    UniApp Frontend   │  │
│           │  (独立表结构)        (手机端)         │  │
│           └─────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### 1.2 核心设计原则

| 原则 | 说明 | 实现方式 |
|------|------|---------|
| **松耦合** | CMS通过API与主系统通信 | REST API + 微服务模式 |
| **高内聚** | CMS业务逻辑独立完整 | 独立的数据表、Logic层、Handler层 |
| **可插拔** | 支持启用/禁用 | 数据库标志位 + 路由动态注册 |
| **权限隔离** | CMS权限独立管理 | 专属Casbin规则集 |
| **数据隔离** | CMS数据与主系统分离 | 表名前缀cms_开头 |

---

## 二、后端实现方案

### 2.1 目录结构规划

```
power-admin-server/
├── api/
│   └── cms-module.api          # CMS模块API定义（新增）
├── internal/
│   ├── handler/
│   │   └── cms/                # CMS处理层（新增）
│   │       ├── cmscontenthandler.go
│   │       ├── cmscategoryhandler.go
│   │       ├── cmspublishhandler.go
│   │       └── cmsusermanagementhandler.go
│   ├── logic/
│   │   └── cms/                # CMS业务逻辑层（新增）
│   │       ├── cmscontentlogic.go
│   │       ├── cmscategorylogic.go
│   │       ├── cmspublishlogic.go
│   │       └── cmsusermanagementlogic.go
│   ├── svc/
│   │   └── servicecontext.go   # 添加CMS相关服务（修改）
│   └── types/
│       └── cms.go              # CMS类型定义（新增）
├── db/
│   └── migrations/
│       └── cms_schema.sql       # CMS数据库表（新增）
└── pkg/
    └── plugins/
        └── cms-plugin.go        # 插件管理核心（新增）
```

### 2.2 数据库设计

#### CMS核心表结构

```sql
-- CMS内容表
CREATE TABLE cms_content (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL COMMENT '文章标题',
  slug VARCHAR(255) UNIQUE COMMENT 'URL别名',
  content LONGTEXT COMMENT '文章内容',
  excerpt VARCHAR(500) COMMENT '文章摘要',
  category_id BIGINT COMMENT '分类ID',
  author_id BIGINT COMMENT '作者ID',
  status TINYINT DEFAULT 1 COMMENT '1:草稿 2:已发布 3:已删除',
  view_count INT DEFAULT 0 COMMENT '浏览次数',
  comment_count INT DEFAULT 0 COMMENT '评论数',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  published_at TIMESTAMP COMMENT '发布时间',
  KEY idx_category (category_id),
  KEY idx_author (author_id),
  KEY idx_status (status),
  KEY idx_created (created_at)
);

-- CMS分类表
CREATE TABLE cms_category (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) UNIQUE,
  description TEXT,
  parent_id BIGINT COMMENT '父分类ID',
  sort INT DEFAULT 0,
  status TINYINT DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_parent (parent_id),
  KEY idx_status (status)
);

-- CMS用户表（访问者）
CREATE TABLE cms_users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(100) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE,
  password VARCHAR(255) NOT NULL,
  nickname VARCHAR(100),
  avatar VARCHAR(500),
  bio TEXT,
  status TINYINT DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_email (email)
);

-- CMS权限配置表
CREATE TABLE cms_permissions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) UNIQUE NOT NULL,
  description TEXT,
  resource VARCHAR(100),
  action VARCHAR(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  KEY idx_resource_action (resource, action)
);

-- CMS-系统用户角色映射表
CREATE TABLE cms_admin_roles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  admin_id BIGINT NOT NULL COMMENT '主系统管理员ID',
  role_name VARCHAR(50) NOT NULL COMMENT 'cms_admin/cms_editor/cms_viewer',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_admin_role (admin_id, role_name),
  KEY idx_admin (admin_id)
);

-- 插件启用状态表
CREATE TABLE plugin_status (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  plugin_name VARCHAR(100) UNIQUE NOT NULL COMMENT 'cms',
  enabled TINYINT DEFAULT 0,
  config JSON,
  installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 2.3 API定义（cms-module.api）

```api
syntax = "v1"

info(
    title: "CMS内容管理系统API"
    desc: "可插拔CMS模块，提供内容、分类、用户、权限管理"
    author: "Your Team"
    email: "support@example.com"
    version: "1.0.0"
)

// 内容管理相关
service cms {
    @handler CmsContentList
    get /api/cms/admin/contents (ListContentsReq) returns (ListContentsResp)
    
    @handler CmsContentDetail
    get /api/cms/admin/contents/:id (DetailReq) returns (CmsContentDetailResp)
    
    @handler CmsContentCreate
    post /api/cms/admin/contents (CreateContentReq) returns (CommonResp)
    
    @handler CmsContentUpdate
    put /api/cms/admin/contents/:id (UpdateContentReq) returns (CommonResp)
    
    @handler CmsContentDelete
    delete /api/cms/admin/contents/:id (DetailReq) returns (CommonResp)
    
    // 分类管理
    @handler CmsCategoryList
    get /api/cms/admin/categories (ListCategoriesReq) returns (ListCategoriesResp)
    
    @handler CmsCategoryTree
    get /api/cms/admin/categories/tree returns (CategoryTreeResp)
    
    @handler CmsCategoryCreate
    post /api/cms/admin/categories (CreateCategoryReq) returns (CommonResp)
    
    @handler CmsCategoryUpdate
    put /api/cms/admin/categories/:id (UpdateCategoryReq) returns (CommonResp)
    
    @handler CmsCategoryDelete
    delete /api/cms/admin/categories/:id (DetailReq) returns (CommonResp)
    
    // 发布管理
    @handler CmsPublish
    post /api/cms/admin/publish (PublishReq) returns (CommonResp)
    
    @handler CmsUnpublish
    post /api/cms/admin/unpublish/:id (DetailReq) returns (CommonResp)
    
    // CMS用户管理
    @handler CmsUserList
    get /api/cms/admin/users (ListCmsUsersReq) returns (ListCmsUsersResp)
    
    @handler CmsUserCreate
    post /api/cms/admin/users (CreateCmsUserReq) returns (CommonResp)
    
    @handler CmsUserUpdate
    put /api/cms/admin/users/:id (UpdateCmsUserReq) returns (CommonResp)
    
    @handler CmsUserDelete
    delete /api/cms/admin/users/:id (DetailReq) returns (CommonResp)
}

// 前台API（UniApp访问）
service cms-public {
    @handler PublicContentList
    get /api/cms/public/contents (ListContentsReq) returns (ListContentsResp)
    
    @handler PublicContentDetail
    get /api/cms/public/contents/:id (DetailReq) returns (CmsContentDetailResp)
    
    @handler PublicCategoryList
    get /api/cms/public/categories returns (ListCategoriesResp)
}

// 类型定义
type ListContentsReq {
    Page     int    `form:"page,default=1"`
    PageSize int    `form:"pageSize,default=10"`
    Category int64  `form:"category_id,optional"`
    Status   int    `form:"status,optional"`
    Search   string `form:"search,optional"`
}

type CmsContentDetailResp {
    Code int       `json:"code"`
    Msg  string    `json:"msg"`
    Data CmsContent `json:"data"`
}

type CmsContent {
    Id          int64     `json:"id"`
    Title       string    `json:"title"`
    Slug        string    `json:"slug"`
    Content     string    `json:"content"`
    Excerpt     string    `json:"excerpt"`
    CategoryId  int64     `json:"category_id"`
    AuthorId    int64     `json:"author_id"`
    Status      int       `json:"status"`
    ViewCount   int       `json:"view_count"`
    CreatedAt   string    `json:"created_at"`
    PublishedAt string    `json:"published_at"`
}

// ... 其他类型定义
```

### 2.4 插件管理核心（pkg/plugins/cms-plugin.go）

```go
package plugins

import (
    "context"
    "database/sql"
    "encoding/json"
)

type PluginInterface interface {
    // 获取插件信息
    GetInfo() PluginInfo
    // 初始化插件
    Init(ctx context.Context, config map[string]interface{}) error
    // 启用插件
    Enable(ctx context.Context) error
    // 禁用插件
    Disable(ctx context.Context) error
    // 卸载插件
    Uninstall(ctx context.Context) error
    // 获取插件菜单
    GetMenuItems(adminId int64) ([]MenuItem, error)
    // 获取权限规则
    GetPermissionRules() ([]PermissionRule, error)
}

type PluginInfo struct {
    Name        string   `json:"name"`
    Version     string   `json:"version"`
    Description string   `json:"description"`
    Author      string   `json:"author"`
    Permissions []string `json:"permissions"`
}

type MenuItem struct {
    Id       int64       `json:"id"`
    Name     string      `json:"name"`
    Path     string      `json:"path"`
    Icon     string      `json:"icon"`
    Component string     `json:"component"`
    Order    int         `json:"order"`
    Children []MenuItem  `json:"children"`
}

type PermissionRule struct {
    PType string // p, g, g2, g3
    V1    string // role, user, etc
    V2    string // resource
    V3    string // action
}

// CMS插件实现
type CMSPlugin struct {
    db *sql.DB
    enabled bool
}

func (p *CMSPlugin) GetInfo() PluginInfo {
    return PluginInfo{
        Name:        "CMS内容管理系统",
        Version:     "1.0.0",
        Description: "可插拔的内容管理系统，支持文章、分类、发布、权限管理",
        Author:      "Your Team",
        Permissions: []string{
            "cms_admin",
            "cms_editor",
            "cms_viewer",
        },
    }
}

func (p *CMSPlugin) Init(ctx context.Context, config map[string]interface{}) error {
    // 初始化CMS表、权限等
    return nil
}

func (p *CMSPlugin) Enable(ctx context.Context) error {
    _, err := p.db.ExecContext(ctx,
        "UPDATE plugin_status SET enabled = 1 WHERE plugin_name = ?",
        "cms")
    return err
}

func (p *CMSPlugin) Disable(ctx context.Context) error {
    _, err := p.db.ExecContext(ctx,
        "UPDATE plugin_status SET enabled = 0 WHERE plugin_name = ?",
        "cms")
    return err
}

func (p *CMSPlugin) GetMenuItems(adminId int64) ([]MenuItem, error) {
    // 返回CMS菜单项，管理员将其注入到系统菜单中
    return []MenuItem{
        {
            Id:        9999,
            Name:      "CMS管理",
            Path:      "/cms",
            Icon:      "mdi:file-document-multiple",
            Component: "CmsLayout",
            Order:     10,
            Children: []MenuItem{
                {
                    Id:        10001,
                    Name:      "内容管理",
                    Path:      "/cms/content",
                    Icon:      "mdi:file-document",
                    Component: "CmsContentList",
                },
                {
                    Id:        10002,
                    Name:      "分类管理",
                    Path:      "/cms/category",
                    Icon:      "mdi:folder-multiple",
                    Component: "CmsCategoryList",
                },
                {
                    Id:        10003,
                    Name:      "访客管理",
                    Path:      "/cms/users",
                    Icon:      "mdi:account-multiple",
                    Component: "CmsUserList",
                },
            },
        },
    }
}

func (p *CMSPlugin) GetPermissionRules() ([]PermissionRule, error) {
    rules := []PermissionRule{
        // CMS管理员权限
        {PType: "p", V1: "cms_admin", V2: "cms", V3: "*"},
        // CMS编辑权限
        {PType: "p", V1: "cms_editor", V2: "cms_content", V3: "create"},
        {PType: "p", V1: "cms_editor", V2: "cms_content", V3: "edit"},
        {PType: "p", V1: "cms_editor", V2: "cms_content", V3: "delete"},
        {PType: "p", V1: "cms_editor", V2: "cms_content", V3: "publish"},
        // CMS查看权限
        {PType: "p", V1: "cms_viewer", V2: "cms_content", V3: "read"},
        {PType: "p", V1: "cms_viewer", V2: "cms_category", V3: "read"},
    }
    return rules, nil
}
```

### 2.5 Logic层示例（cms/content）

```go
package cms

import (
    "context"
    "fmt"
    "power-admin/internal/svc"
    "power-admin/internal/types"
)

type CmsContentLogic struct {
    ctx context.Context
    svcCtx *svc.ServiceContext
}

func NewCmsContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CmsContentLogic {
    return &CmsContentLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

// 列表
func (l *CmsContentLogic) List(req *types.ListContentsReq) (*types.ListContentsResp, error) {
    // 实现内容列表查询
    // 权限检查：当前用户是否有cms:content:read权限
    return &types.ListContentsResp{
        Code: 0,
        Msg:  "success",
    }, nil
}

// 创建
func (l *CmsContentLogic) Create(req *types.CreateContentReq) error {
    // 权限检查：当前用户是否有cms:content:create权限
    // 验证输入
    if err := l.validateCreateRequest(req); err != nil {
        return err
    }
    
    // 插入到cms_content表
    _, err := l.svcCtx.DB.ExecContext(l.ctx,
        "INSERT INTO cms_content (title, slug, content, category_id) VALUES (?, ?, ?, ?)",
        req.Title, req.Slug, req.Content, req.CategoryId)
    
    return err
}

func (l *CmsContentLogic) validateCreateRequest(req *types.CreateContentReq) error {
    if req.Title == "" {
        return fmt.Errorf("title is required")
    }
    return nil
}
```

---

## 三、前端实现方案

### 3.1 前端架构决策

**推荐方案：集成到现有Admin Dashboard中**

**原因：**
- 避免重复开发认证、授权系统
- 复用现有的UI框架和样式
- 统一用户体验
- 便于权限管理集成

### 3.2 前端目录结构

```
power-admin-web/
├── src/
│   ├── pages/
│   │   ├── system/           # 现有系统管理页面
│   │   └── cms/              # CMS管理页面（新增）
│   │       ├── content/
│   │       │   ├── ContentList.vue
│   │       │   ├── ContentDetail.vue
│   │       │   └── ContentForm.vue
│   │       ├── category/
│   │       │   ├── CategoryList.vue
│   │       │   └── CategoryTree.vue
│   │       ├── users/
│   │       │   ├── UserList.vue
│   │       │   └── UserForm.vue
│   │       └── CmsLayout.vue
│   ├── api/
│   │   ├── system.ts         # 现有系统API
│   │   └── cms.ts            # CMS API（新增）
│   ├── stores/
│   │   └── cms.ts            # CMS状态管理（新增）
│   └── router/
│       └── index.ts          # 动态添加CMS路由
```

### 3.3 前端API接口（api/cms.ts）

```typescript
import request from './request'

// 内容管理
export const getContentList = (params: any) =>
  request.get('/api/cms/admin/contents', { params })

export const getContentDetail = (id: number) =>
  request.get(`/api/cms/admin/contents/${id}`)

export const createContent = (data: any) =>
  request.post('/api/cms/admin/contents', data)

export const updateContent = (id: number, data: any) =>
  request.put(`/api/cms/admin/contents/${id}`, data)

export const deleteContent = (id: number) =>
  request.delete(`/api/cms/admin/contents/${id}`)

// 分类管理
export const getCategoryList = (params: any) =>
  request.get('/api/cms/admin/categories', { params })

export const getCategoryTree = () =>
  request.get('/api/cms/admin/categories/tree')

export const createCategory = (data: any) =>
  request.post('/api/cms/admin/categories', data)

export const updateCategory = (id: number, data: any) =>
  request.put(`/api/cms/admin/categories/${id}`, data)

export const deleteCategory = (id: number) =>
  request.delete(`/api/cms/admin/categories/${id}`)

// 发布管理
export const publishContent = (id: number) =>
  request.post(`/api/cms/admin/publish`, { id })

export const unpublishContent = (id: number) =>
  request.post(`/api/cms/admin/unpublish/${id}`)

// CMS用户管理
export const getCmsUserList = (params: any) =>
  request.get('/api/cms/admin/users', { params })

export const createCmsUser = (data: any) =>
  request.post('/api/cms/admin/users', data)

export const updateCmsUser = (id: number, data: any) =>
  request.put(`/api/cms/admin/users/${id}`, data)

export const deleteCmsUser = (id: number) =>
  request.delete(`/api/cms/admin/users/${id}`)
```

### 3.4 菜单动态注入（router/index.ts修改）

```typescript
import { createRouter, createWebHistory } from 'vue-router'
import { useMenuStore } from '@/stores/menu'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // ... 现有路由
  ]
})

// 菜单更新钩子
router.beforeEach(async (to, from, next) => {
  const menuStore = useMenuStore()
  
  if (!menuStore.menuLoaded) {
    await menuStore.loadMenus()
    
    // 检查CMS插件是否启用
    const cmsEnabled = await checkCmsPluginStatus()
    
    if (cmsEnabled) {
      // 动态添加CMS路由
      addCmsRoutes(router, menuStore)
    }
  }
  
  next()
})

function addCmsRoutes(router, menuStore) {
  const cmsRoutes = [
    {
      path: '/cms',
      component: () => import('@/pages/cms/CmsLayout.vue'),
      meta: { title: 'CMS管理' },
      children: [
        {
          path: 'content',
          component: () => import('@/pages/cms/content/ContentList.vue'),
          meta: { title: '内容管理' }
        },
        {
          path: 'category',
          component: () => import('@/pages/cms/category/CategoryList.vue'),
          meta: { title: '分类管理' }
        },
        {
          path: 'users',
          component: () => import('@/pages/cms/users/UserList.vue'),
          meta: { title: '访客管理' }
        }
      ]
    }
  ]
  
  cmsRoutes.forEach(route => router.addRoute(route))
}
```

### 3.5 CMS状态管理（stores/cms.ts）

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as cmsApi from '@/api/cms'

export const useCmsStore = defineStore('cms', () => {
  // 状态
  const contentList = ref([])
  const categoryList = ref([])
  const cmsUsers = ref([])
  const loading = ref(false)
  
  // 计算属性
  const contentCount = computed(() => contentList.value.length)
  
  // 方法
  async function fetchContentList(params) {
    loading.value = true
    try {
      const res = await cmsApi.getContentList(params)
      contentList.value = res.data.list || []
      return res
    } finally {
      loading.value = false
    }
  }
  
  async function createContent(data) {
    await cmsApi.createContent(data)
    // 刷新列表
    await fetchContentList({})
  }
  
  async function fetchCategoryList(params) {
    const res = await cmsApi.getCategoryList(params)
    categoryList.value = res.data.list || []
    return res
  }
  
  async function fetchCmsUserList(params) {
    const res = await cmsApi.getCmsUserList(params)
    cmsUsers.value = res.data.list || []
    return res
  }
  
  return {
    // 状态
    contentList,
    categoryList,
    cmsUsers,
    loading,
    
    // 计算属性
    contentCount,
    
    // 方法
    fetchContentList,
    createContent,
    fetchCategoryList,
    fetchCmsUserList
  }
})
```

---

## 四、UniApp 手机端适配

### 4.1 独立UniApp项目

```
cms-uniapp/
├── src/
│   ├── pages/
│   │   ├── index/
│   │   │   └── index.vue      # 首页-文章列表
│   │   ├── article/
│   │   │   └── detail.vue     # 文章详情
│   │   ├── category/
│   │   │   └── list.vue       # 分类浏览
│   │   └── user/
│   │       ├── login.vue
│   │       ├── register.vue
│   │       └── profile.vue
│   ├── api/
│   │   └── cms.ts             # CMS API调用
│   ├── components/
│   │   ├── ArticleCard.vue
│   │   ├── CategoryTag.vue
│   │   └── CommentList.vue
│   ├── stores/
│   │   └── cms.ts
│   └── utils/
│       └── request.ts          # HTTP请求工具
└── pages.json                  # 页面路由配置
```

### 4.2 UniApp API调用示例

```typescript
// api/cms.ts
import { http } from '@/utils/request'

// 获取文章列表
export const getPublicContentList = (params: any) =>
  http.request({
    method: 'GET',
    url: '/api/cms/public/contents',
    params
  })

// 获取文章详情
export const getPublicContentDetail = (id: number) =>
  http.request({
    method: 'GET',
    url: `/api/cms/public/contents/${id}`
  })

// 获取分类列表
export const getPublicCategoryList = () =>
  http.request({
    method: 'GET',
    url: '/api/cms/public/categories'
  })

// 用户登录
export const cmsUserLogin = (data: any) =>
  http.request({
    method: 'POST',
    url: '/api/cms/public/auth/login',
    data
  })
```

---

## 五、权限管理集成

### 5.1 Casbin规则集成

**现有系统Casbin规则添加：**

```sql
-- CMS管理员权限（继承系统管理员的所有CMS权限）
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES 
('p', 'cms_admin', '/api/cms/admin/contents', 'GET', ''),
('p', 'cms_admin', '/api/cms/admin/contents', 'POST', ''),
('p', 'cms_admin', '/api/cms/admin/contents', 'PUT', ''),
('p', 'cms_admin', '/api/cms/admin/contents', 'DELETE', ''),
('p', 'cms_editor', '/api/cms/admin/contents', 'GET', ''),
('p', 'cms_editor', '/api/cms/admin/contents', 'POST', ''),
('p', 'cms_editor', '/api/cms/admin/contents', 'PUT', ''),
('p', 'cms_viewer', '/api/cms/admin/contents', 'GET', ''),
-- 将系统角色映射到CMS角色
('g', '1', 'cms_admin');  -- 用户ID 1 是 CMS管理员
```

### 5.2 权限检查中间件

```go
// internal/middleware/cms-permission.go
package middleware

import (
    "context"
    "net/http"
    "power-admin/internal/svc"
)

func CmsPermissionMiddleware(svcCtx *svc.ServiceContext) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1. 检查CMS插件是否启用
            enabled, err := svcCtx.CheckPluginEnabled("cms")
            if err != nil || !enabled {
                http.Error(w, "CMS plugin not enabled", http.StatusForbidden)
                return
            }
            
            // 2. 检查用户是否有CMS权限
            userId := r.Header.Get("X-User-ID")
            hasPermission, err := svcCtx.Enforcer.Enforce(
                userId, 
                r.URL.Path, 
                r.Method,
            )
            
            if err != nil || !hasPermission {
                http.Error(w, "Permission denied", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 六、部署和启用流程

### 6.1 CMS插件启用流程图

```
┌─────────────────────┐
│  管理员购买CMS插件   │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────────────┐
│  插件管理界面 → 点击"启用"     │
└──────────┬──────────────────┘
           │
           ▼
┌──────────────────────────────────┐
│  1. 检查依赖和兼容性              │
│  2. 创建CMS数据表                 │
│  3. 初始化权限规则到Casbin        │
│  4. 更新plugin_status表为enabled  │
│  5. 加载CMS菜单到系统菜单        │
│  6. 生成API路由                   │
└──────────┬───────────────────────┘
           │
           ▼
┌──────────────────────────┐
│  CMS管理菜单出现在左侧菜单栏  │
│  用户可以开始使用CMS功能      │
└──────────────────────────┘
```

### 6.2 启用API（可选）

```go
// 插件启用接口
POST /api/admin/plugins/enable
{
  "plugin_name": "cms",
  "config": {
    "enable_comments": true,
    "enable_ratings": true,
    "posts_per_page": 10
  }
}

// 响应
{
  "code": 0,
  "msg": "CMS plugin enabled successfully",
  "data": {
    "plugin_name": "cms",
    "version": "1.0.0",
    "enabled": true,
    "menu_items": [...]
  }
}
```

---

## 七、数据隔离和共享策略

### 7.1 数据隔离

| 类型 | 隔离策略 | 实现方式 |
|------|---------|---------|
| **业务数据** | 完全隔离 | CMS表名统一cms_前缀 |
| **权限数据** | 部分共享 | 使用系统Casbin表，添加CMS规则 |
| **用户数据** | 部分共享 | 系统用户和CMS访客分离 |
| **配置数据** | 共享存储 | plugin_status表统一存储 |

### 7.2 数据流向

```
系统用户(system_user)
    │
    ├─→ 系统权限检查 (admin dashboard)
    │
    └─→ CMS权限检查 (CMS operations)
            │
            ├─→ CMS内容管理 (cms_content)
            ├─→ CMS分类管理 (cms_category)
            └─→ CMS用户管理 (cms_users)
                    │
                    └─→ 前台访客 (UniApp)
```

---

## 八、实施时间表

| 阶段 | 任务 | 时间 | 负责人 |
|------|------|------|--------|
| **第一阶段** | 后端API开发 + 数据库设计 | 2周 | Backend |
| **第二阶段** | 插件框架实现 + 权限集成 | 1周 | Backend |
| **第三阶段** | 前端页面开发 + 菜单集成 | 2周 | Frontend |
| **第四阶段** | UniApp手机端开发 | 2周 | Mobile |
| **第五阶段** | 集成测试 + 优化调整 | 1周 | QA |
| **第六阶段** | 文档编写 + 部署上线 | 1周 | DevOps |

---

## 九、风险评估和解决方案

| 风险 | 等级 | 解决方案 |
|------|------|---------|
| **权限冲突** | 高 | 使用独立Casbin规则集，定期检查 |
| **性能下降** | 中 | 优化SQL查询，添加适当索引 |
| **数据不一致** | 中 | 使用事务，确保原子性操作 |
| **插件卸载** | 中 | 保留数据但禁用功能，支持重新安装 |
| **跨域问题** | 低 | 配置CORS中间件 |

---

## 十、关键注意事项

1. **API版本管理**: CMS API使用统一的版本前缀 `/api/cms/v1/`
2. **安全认证**: 所有CMS API都需要系统JWT令牌认证
3. **审计日志**: 记录所有CMS操作到审计表
4. **数据备份**: CMS数据定期备份到独立位置
5. **性能监控**: 监控CMS API响应时间，设置告警阈值
6. **灰度发布**: 新版本CMS通过灰度发布逐步上线

---

## 总结

该方案通过以下方式实现了真正的**可插拔架构**：

✅ **松耦合**: CMS通过独立的API与主系统通信  
✅ **独立部署**: CMS可独立开发、测试、更新  
✅ **灵活启用**: 管理员可随时启用/禁用CMS  
✅ **权限隔离**: CMS权限完全独立管理  
✅ **数据安全**: CMS数据与系统数据隔离  
✅ **用户体验**: 无缝集成到现有admin dashboard  

