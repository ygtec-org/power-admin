# CMS插件实施完整检查清单

## 📋 项目概览

- **方案选择**: 集成式架构
- **预计周期**: 2周
- **技术栈**: Go + Vue3 + MySQL + Casbin
- **团队规模**: 2-3人

---

## ✅ 第一周：后端开发

### Day 1: 项目准备和数据库设计

- [ ] **1.1 创建项目文件夹**
  - [ ] `mkdir internal/handler/cms`
  - [ ] `mkdir internal/logic/cms`
  - [ ] `mkdir db/migrations`
  - [ ] `touch internal/types/cms.go`
  - [ ] `touch api/cms.api`

- [ ] **1.2 编写数据库迁移脚本**
  - [ ] `cms_content` 表 - 文章内容表
  - [ ] `cms_category` 表 - 分类表
  - [ ] `cms_users` 表 - CMS访客表
  - [ ] `cms_permissions` 表 - 权限表
  - [ ] `cms_admin_roles` 表 - 管理员-角色映射
  - [ ] `plugin_status` 表 - 插件状态表
  - [ ] 添加适当的索引和外键约束

- [ ] **1.3 执行数据库初始化**
  ```bash
  mysql -u root -p power_admin < db/migrations/001_init_cms_schema.sql
  ```

- [ ] **1.4 验证数据库**
  ```bash
  # 连接MySQL并验证表是否创建成功
  mysql> SHOW TABLES LIKE 'cms%';
  ```

### Day 2: API定义和类型定义

- [ ] **2.1 编写API定义文件** (api/cms.api)
  - [ ] 内容管理接口 (List/Get/Create/Update/Delete)
  - [ ] 分类管理接口 (List/Tree/Create/Update/Delete)
  - [ ] 发布管理接口 (Publish/Unpublish)
  - [ ] CMS用户管理接口 (List/Create/Update/Delete)
  - [ ] 前台公开接口 (PublicList/PublicDetail)

- [ ] **2.2 定义Types** (internal/types/cms.go)
  - [ ] 请求类型 (ListContentsReq, CreateContentReq等)
  - [ ] 响应类型 (ListContentsResp, CmsContentDetailResp等)
  - [ ] 数据模型 (CmsContent, CmsCategory等)

- [ ] **2.3 生成Handler和Logic骨架**
  ```bash
  goctl api go -api api/cms.api -dir . --style=go-zero
  ```
  或手动创建：
  - [ ] `internal/handler/cms/cmscontenthandler.go`
  - [ ] `internal/handler/cms/cmscategoryhandler.go`
  - [ ] `internal/handler/cms/cmspublishhandler.go`
  - [ ] `internal/handler/cms/cmsusermanagementhandler.go`

### Day 3: Handler和Logic实现

- [ ] **3.1 实现内容管理Logic**
  ```go
  internal/logic/cms/cmscontentlogic.go
  ```
  - [ ] `List()` - 列表查询
  - [ ] `Detail()` - 详情查询
  - [ ] `Create()` - 创建内容
  - [ ] `Update()` - 更新内容
  - [ ] `Delete()` - 删除内容
  - [ ] 输入验证
  - [ ] 权限检查

- [ ] **3.2 实现分类管理Logic**
  ```go
  internal/logic/cms/cmscategorylogic.go
  ```
  - [ ] `List()` - 分类列表
  - [ ] `Tree()` - 树形结构
  - [ ] `Create()` - 创建分类
  - [ ] `Update()` - 更新分类
  - [ ] `Delete()` - 删除分类
  - [ ] 验证父分类存在性

- [ ] **3.3 实现发布管理Logic**
  ```go
  internal/logic/cms/cmspublishlogic.go
  ```
  - [ ] `Publish()` - 发布内容
  - [ ] `Unpublish()` - 取消发布
  - [ ] 设置发布时间

- [ ] **3.4 实现CMS用户管理Logic**
  ```go
  internal/logic/cms/cmsusermanagementlogic.go
  ```
  - [ ] `List()` - 用户列表
  - [ ] `Create()` - 创建用户
  - [ ] `Update()` - 更新用户
  - [ ] `Delete()` - 删除用户

- [ ] **3.5 实现Handler层**
  - [ ] 参数绑定
  - [ ] 权限检查中间件
  - [ ] 调用Logic层
  - [ ] 错误处理
  - [ ] 返回JSON响应

### Day 4: 插件框架和权限集成

- [ ] **4.1 创建插件框架**
  ```go
  pkg/plugins/cms-plugin.go
  ```
  - [ ] 定义 `PluginInterface`
  - [ ] 实现 `CMSPlugin` 结构体
  - [ ] 实现 `GetInfo()` - 获取插件信息
  - [ ] 实现 `Init()` - 初始化
  - [ ] 实现 `Enable()` - 启用插件
  - [ ] 实现 `Disable()` - 禁用插件
  - [ ] 实现 `GetMenuItems()` - 返回菜单项
  - [ ] 实现 `GetPermissionRules()` - 返回权限规则

- [ ] **4.2 集成权限系统**
  - [ ] 添加CMS权限规则到Casbin
  - [ ] 定义CMS角色 (cms_admin/cms_editor/cms_viewer)
  - [ ] 创建权限检查中间件
  - [ ] 验证权限逻辑

- [ ] **4.3 修改ServiceContext**
  ```go
  internal/svc/servicecontext.go
  ```
  - [ ] 添加CMSPlugin字段
  - [ ] 初始化CMS插件
  - [ ] 添加权限验证方法

### Day 5: 路由注册和测试

- [ ] **5.1 注册路由**
  - [ ] 修改main.go，添加CMS路由
  - [ ] 检查路由是否正确注册
  - [ ] 验证路由顺序（CMS路由应在权限检查后）

- [ ] **5.2 单元测试**
  - [ ] 编写CMS API单元测试
  - [ ] 测试内容CRUD操作
  - [ ] 测试权限检查
  - [ ] 测试分类树形结构

- [ ] **5.3 集成测试**
  - [ ] 测试整个请求流程
  - [ ] 验证权限验证
  - [ ] 测试错误处理

- [ ] **5.4 后端验证**
  ```bash
  # 启动后端服务
  go run ./cmd/api/main.go
  
  # 使用curl测试API
  curl -X GET http://localhost:8888/api/cms/admin/contents \
    -H "Authorization: Bearer YOUR_TOKEN"
  ```

---

## ✅ 第二周：前端开发

### Day 1: 页面框架和API接口

- [ ] **1.1 创建CMS页面目录**
  ```bash
  mkdir -p src/pages/cms/content
  mkdir -p src/pages/cms/category
  mkdir -p src/pages/cms/users
  ```

- [ ] **1.2 创建CMS API接口**
  ```typescript
  src/api/cms.ts
  ```
  - [ ] 内容API (getContentList/getContentDetail/createContent等)
  - [ ] 分类API (getCategoryTree/createCategory等)
  - [ ] 发布API (publishContent/unpublishContent)
  - [ ] 用户API (getCmsUserList/createCmsUser等)

- [ ] **1.3 创建状态管理**
  ```typescript
  src/stores/cms.ts (Pinia)
  ```
  - [ ] contentList 状态
  - [ ] categoryList 状态
  - [ ] cmsUsers 状态
  - [ ] fetchContentList() 方法
  - [ ] fetchCategoryList() 方法
  - [ ] fetchCmsUserList() 方法

### Day 2: 内容管理页面

- [ ] **2.1 创建内容列表页面**
  ```vue
  src/pages/cms/content/ContentList.vue
  ```
  - [ ] 表格展示内容列表
  - [ ] 分页控件
  - [ ] 搜索和筛选功能
  - [ ] 新增/编辑/删除按钮
  - [ ] 发布/取消发布按钮

- [ ] **2.2 创建内容详情/编辑页面**
  ```vue
  src/pages/cms/content/ContentForm.vue
  ```
  - [ ] 标题输入框
  - [ ] 内容编辑器 (使用富文本编辑器如Quill)
  - [ ] 分类选择
  - [ ] 摘要编辑
  - [ ] 保存/取消按钮

- [ ] **2.3 实现内容管理功能**
  - [ ] 加载内容列表
  - [ ] 创建新内容
  - [ ] 编辑内容
  - [ ] 删除内容
  - [ ] 发布内容
  - [ ] 错误处理和提示

### Day 3: 分类管理和用户管理

- [ ] **3.1 创建分类管理页面**
  ```vue
  src/pages/cms/category/CategoryList.vue
  ```
  - [ ] 树形结构展示分类
  - [ ] 新增/编辑/删除分类
  - [ ] 排序功能
  - [ ] 父分类级联

- [ ] **3.2 创建访客管理页面**
  ```vue
  src/pages/cms/users/UserList.vue
  ```
  - [ ] 用户表格
  - [ ] 用户搜索
  - [ ] 新增/编辑/禁用用户
  - [ ] 用户权限分配

- [ ] **3.3 创建CMS布局组件**
  ```vue
  src/pages/cms/CmsLayout.vue
  ```
  - [ ] 顶部导航
  - [ ] 左侧菜单 (内容/分类/用户)
  - [ ] 主内容区域

### Day 4: 菜单和路由集成

- [ ] **4.1 动态菜单加载**
  - [ ] 修改菜单加载逻辑，检查CMS是否启用
  - [ ] 动态添加CMS菜单项到左侧菜单
  - [ ] 验证菜单显示

- [ ] **4.2 动态路由注册**
  ```typescript
  src/router/index.ts
  ```
  - [ ] 添加 `registerCmsRoutes()` 函数
  - [ ] 在路由守卫中注册CMS路由
  - [ ] 验证路由可访问

- [ ] **4.3 权限检查**
  - [ ] 检查用户是否有CMS权限
  - [ ] 只有授权用户才能访问CMS菜单
  - [ ] 检查各功能的细粒度权限

### Day 5: 测试和优化

- [ ] **5.1 功能测试**
  - [ ] [ ] 测试登录系统
  - [ ] [ ] 验证CMS菜单显示
  - [ ] [ ] 测试创建文章
  - [ ] [ ] 测试编辑文章
  - [ ] [ ] 测试删除文章
  - [ ] [ ] 测试发布/取消发布
  - [ ] [ ] 测试分类管理
  - [ ] [ ] 测试用户管理
  - [ ] [ ] 测试权限验证（无权限用户无法访问）

- [ ] **5.2 性能优化**
  - [ ] 优化表格加载速度
  - [ ] 添加加载状态指示器
  - [ ] 优化分类树形渲染
  - [ ] 缓存分类列表

- [ ] **5.3 UI/UX 改进**
  - [ ] 添加确认对话框（删除操作）
  - [ ] 添加成功/失败提示
  - [ ] 优化表单验证提示
  - [ ] 响应式设计（移动端适配）

- [ ] **5.4 前端验证**
  ```bash
  # 启动前端开发服务器
  npm run dev
  
  # 访问系统
  http://localhost:5173
  ```

---

## ✅ 权限和菜单配置

### 权限设置

- [ ] **添加CMS权限到系统**
  ```sql
  -- 权限表
  INSERT INTO permission (name, resource, action, description) VALUES
  ('CMS内容查看', 'cms_content', 'read', '查看CMS内容'),
  ('CMS内容管理', 'cms_content', 'manage', '创建/编辑/删除CMS内容'),
  ('CMS分类管理', 'cms_category', 'manage', '管理CMS分类'),
  ('CMS用户管理', 'cms_users', 'manage', '管理CMS访客');
  ```

- [ ] **添加CMS角色**
  ```sql
  INSERT INTO role (name, description, status) VALUES
  ('cms_admin', 'CMS管理员', 1),
  ('cms_editor', 'CMS编辑', 1),
  ('cms_viewer', 'CMS查看者', 1);
  ```

- [ ] **添加Casbin规则**
  ```sql
  INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
  ('p', 'cms_admin', '/api/cms/admin/contents', 'GET', ''),
  ('p', 'cms_admin', '/api/cms/admin/contents', 'POST', ''),
  ('p', 'cms_admin', '/api/cms/admin/contents', 'PUT', ''),
  ('p', 'cms_admin', '/api/cms/admin/contents', 'DELETE', ''),
  ...
  ```

### 菜单配置

- [ ] **在数据库中添加CMS菜单项**
  ```sql
  INSERT INTO menu (menu_name, menu_path, component, icon, parent_id, status) VALUES
  ('CMS管理', '/cms', 'CmsLayout', 'mdi:file-document-multiple', 0, 1),
  ('内容管理', '/cms/content', 'CmsContentList', 'mdi:file-document', LAST_INSERT_ID(), 1),
  ('分类管理', '/cms/category', 'CmsCategoryList', 'mdi:folder-multiple', LAST_INSERT_ID(), 1),
  ('访客管理', '/cms/users', 'CmsUserList', 'mdi:account-multiple', LAST_INSERT_ID(), 1);
  ```

- [ ] **为角色分配CMS菜单权限**
  - [ ] cms_admin 有所有CMS菜单权限
  - [ ] cms_editor 有内容和分类权限
  - [ ] cms_viewer 只有查看权限

---

## ✅ UniApp 手机端开发（可选）

### 项目设置

- [ ] **创建UniApp项目**
  ```bash
  mkdir cms-uniapp
  cd cms-uniapp
  # 使用 HBuilderX 或 vue-cli 创建UniApp项目
  ```

- [ ] **创建页面结构**
  - [ ] `pages/index/index.vue` - 首页文章列表
  - [ ] `pages/article/detail.vue` - 文章详情
  - [ ] `pages/category/list.vue` - 分类浏览
  - [ ] `pages/user/login.vue` - 用户登录
  - [ ] `pages/user/register.vue` - 用户注册
  - [ ] `pages/user/profile.vue` - 用户资料

- [ ] **创建API调用模块**
  ```typescript
  api/cms.ts
  ```
  - [ ] getPublicContentList() - 获取文章列表
  - [ ] getPublicContentDetail() - 获取文章详情
  - [ ] getPublicCategoryList() - 获取分类列表
  - [ ] cmsUserLogin() - 用户登录
  - [ ] cmsUserRegister() - 用户注册

- [ ] **创建组件**
  - [ ] ArticleCard.vue - 文章卡片
  - [ ] CategoryTag.vue - 分类标签
  - [ ] CommentList.vue - 评论列表

- [ ] **状态管理**
  ```typescript
  stores/cms.ts (使用 uni-app store 或 Pinia)
  ```

---

## ✅ 部署和上线

### 部署前检查

- [ ] **代码质量检查**
  - [ ] 运行 linter (Go: golangci-lint, Vue: ESLint)
  - [ ] 运行单元测试
  - [ ] 代码审查

- [ ] **性能测试**
  - [ ] 后端API响应时间 (<100ms)
  - [ ] 前端页面加载时间 (<3s)
  - [ ] 数据库查询性能

- [ ] **安全检查**
  - [ ] SQL注入防护
  - [ ] XSS防护
  - [ ] CSRF防护
  - [ ] 权限验证完整性

### 部署步骤

- [ ] **后端部署**
  ```bash
  # 编译
  cd power-admin-server
  go build -o power-admin.exe ./cmd/api/main.go
  
  # 配置环境
  # 修改 etc/power-admin-api.yaml
  # 设置正确的数据库连接
  
  # 运行
  ./power-admin.exe
  ```

- [ ] **前端部署**
  ```bash
  # 编译
  cd power-admin-web
  npm run build
  
  # 部署到服务器
  # 将 dist 目录上传到Web服务器
  # 或配置反向代理指向构建后的文件
  ```

- [ ] **数据库备份**
  ```bash
  mysqldump -u root -p power_admin > power_admin_backup.sql
  ```

---

## ✅ 上线后验证

- [ ] **功能验证**
  - [ ] [ ] CMS菜单显示正常
  - [ ] [ ] 内容创建、编辑、删除正常
  - [ ] [ ] 分类管理正常
  - [ ] [ ] 用户管理正常
  - [ ] [ ] 发布功能正常

- [ ] **权限验证**
  - [ ] [ ] 超级管理员可访问所有功能
  - [ ] [ ] CMS编辑可创建/编辑内容
  - [ ] [ ] CMS查看者只能查看
  - [ ] [ ] 无权限用户无法访问

- [ ] **性能监控**
  - [ ] [ ] API响应时间正常
  - [ ] [ ] 数据库查询性能正常
  - [ ] [ ] 无内存泄漏

- [ ] **日志检查**
  - [ ] [ ] 无错误日志
  - [ ] [ ] 无警告日志
  - [ ] [ ] 操作日志完整

---

## 📊 进度追踪

### Week 1: 后端开发

| 日期 | 任务 | 状态 | 备注 |
|------|------|------|------|
| Day 1 | 项目准备和数据库 | ⬜ | |
| Day 2 | API定义和类型定义 | ⬜ | |
| Day 3 | Handler和Logic实现 | ⬜ | |
| Day 4 | 插件框架和权限集成 | ⬜ | |
| Day 5 | 路由注册和测试 | ⬜ | |

### Week 2: 前端开发

| 日期 | 任务 | 状态 | 备注 |
|------|------|------|------|
| Day 1 | 页面框架和API接口 | ⬜ | |
| Day 2 | 内容管理页面 | ⬜ | |
| Day 3 | 分类和用户管理 | ⬜ | |
| Day 4 | 菜单和路由集成 | ⬜ | |
| Day 5 | 测试和优化 | ⬜ | |

---

## 🎯 里程碑

- [ ] **Milestone 1**: 数据库和API框架 (Day 3)
- [ ] **Milestone 2**: 后端所有功能完成 (Day 5)
- [ ] **Milestone 3**: 前端页面开发完成 (Day 8)
- [ ] **Milestone 4**: 集成和权限测试 (Day 9)
- [ ] **Milestone 5**: 上线部署 (Day 10)

---

## 🚨 风险预警

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|---------|
| 权限验证复杂 | 中 | 高 | 提前理解Casbin规则 |
| 前后端接口不匹配 | 低 | 中 | 定期同步API定义 |
| 数据库迁移错误 | 低 | 高 | 先在测试环境验证 |
| 性能问题 | 低 | 中 | 提前优化查询和索引 |
| 权限规则遗漏 | 中 | 中 | 编写完整的权限测试用例 |

---

## 📞 支持和问题

- **技术问题**: 查看 CMS_PLUGIN_INTEGRATION_PLAN.md
- **快速开始**: 查看 CMS_QUICK_START.md
- **架构对比**: 查看 CMS_ARCHITECTURE_COMPARISON.md

---

**项目完成预计日期**: 2024年[DATE]

