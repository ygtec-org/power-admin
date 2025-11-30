# CMS插件开发计划 - 稳定完整版

**目标**: 构建一个**稳定、完整、生产级别**的CMS内容管理系统  
**时间**: 充足时间（质量优先）  
**质量标准**: 企业级应用  

---

## 📋 开发阶段划分

### 第一阶段：基础设施和数据库设计（1-2周）

#### 1.1 数据库设计和初始化
- [ ] 创建CMS核心数据表（cms_content、cms_category等）
- [ ] 创建权限和角色管理表
- [ ] 创建评论和标签系统表
- [ ] 创建审计日志表
- [ ] 添加所有必要的索引和约束
- [ ] 编写数据库初始化脚本和迁移文件

**交付物**：
- `power-admin-server/db/migrations/003_cms_schema.sql` - 完整的CMS数据库定义
- `power-admin-server/db/migrations/004_cms_permissions.sql` - CMS权限初始化
- `power-admin-server/db/seeds/cms_seed.sql` - 初始测试数据

#### 1.2 定义数据模型和仓储层
- [ ] 创建CMS模型（Model）
- [ ] 创建CMS仓储层（Repository）
- [ ] 实现数据访问接口
- [ ] 添加完整的错误处理

**交付物**：
- `power-admin-server/pkg/models/cms_*.go` - 所有CMS模型
- `power-admin-server/pkg/repository/cms_repository.go` - CMS仓储实现

---

### 第二阶段：后端业务逻辑开发（2-3周）

#### 2.1 内容管理业务逻辑
- [ ] 创建ContentLogic
  - [ ] 创建内容（包含富文本验证、SEO字段）
  - [ ] 更新内容
  - [ ] 删除内容（硬删除和软删除）
  - [ ] 查询内容列表（支持分页、排序、筛选）
  - [ ] 获取内容详情（包含关联数据）
- [ ] 内容版本管理（可选）
- [ ] 草稿箱功能
- [ ] 内容搜索功能

**交付物**：
- `power-admin-server/internal/logic/cms/content_logic.go`

#### 2.2 分类管理业务逻辑
- [ ] 创建CategoryLogic
  - [ ] 创建分类（支持多级）
  - [ ] 更新分类
  - [ ] 删除分类（包含级联处理）
  - [ ] 获取分类树形结构
  - [ ] 获取分类列表
- [ ] 分类排序功能
- [ ] 分类状态管理

**交付物**：
- `power-admin-server/internal/logic/cms/category_logic.go`

#### 2.3 标签和关键词管理
- [ ] 标签CRUD操作
- [ ] 关键词管理
- [ ] 内容-标签关联

**交付物**：
- `power-admin-server/internal/logic/cms/tag_logic.go`

#### 2.4 发布和工作流
- [ ] 发布内容
- [ ] 取消发布
- [ ] 定时发布功能
- [ ] 工作流状态管理

**交付物**：
- `power-admin-server/internal/logic/cms/publish_logic.go`

#### 2.5 评论管理
- [ ] 评论CRUD
- [ ] 评论审核工作流
- [ ] 评论回复功能
- [ ] 评论垃圾检测

**交付物**：
- `power-admin-server/internal/logic/cms/comment_logic.go`

#### 2.6 CMS用户管理
- [ ] 创建用户（访客注册）
- [ ] 更新用户信息
- [ ] 用户状态管理
- [ ] 用户禁用/封禁

**交付物**：
- `power-admin-server/internal/logic/cms/user_logic.go`

#### 2.7 权限和访问控制
- [ ] CMS角色定义（管理员、编辑、查看者）
- [ ] 权限规则定义
- [ ] 权限验证逻辑
- [ ] 操作审计

**交付物**：
- `power-admin-server/internal/logic/cms/permission_logic.go`

---

### 第三阶段：后端API层开发（1-2周）

#### 3.1 API定义
- [ ] 编写完整的API定义文件（cms.api）
- [ ] 定义所有请求/响应类型
- [ ] 添加API文档注释

**交付物**：
- `power-admin-server/api/cms.api`
- `power-admin-server/internal/types/cms_types.go`

#### 3.2 Handler层实现
- [ ] 内容管理Handler
- [ ] 分类管理Handler
- [ ] 标签管理Handler
- [ ] 发布管理Handler
- [ ] 评论管理Handler
- [ ] 用户管理Handler
- [ ] 权限管理Handler

**要求**：
- 完整的参数验证
- 完整的错误处理
- 统一的响应格式
- 请求日志记录

**交付物**：
- `power-admin-server/internal/handler/cms/*.go`

#### 3.3 路由注册
- [ ] 在routes.go中注册所有CMS路由
- [ ] 添加路由中间件（权限验证、日志等）
- [ ] 支持路由版本管理

**交付物**：
- `power-admin-server/internal/handler/cms/routes.go`

---

### 第四阶段：插件管理和集成（1周）

#### 4.1 插件框架
- [ ] 创建PluginInterface接口
- [ ] 实现CMSPlugin
- [ ] 菜单注入逻辑
- [ ] 权限初始化逻辑

**交付物**：
- `power-admin-server/pkg/plugins/plugin.go`
- `power-admin-server/pkg/plugins/cms_plugin.go`

#### 4.2 主系统集成
- [ ] 修改ServiceContext添加CMS服务
- [ ] 修改数据库迁移添加CMS表
- [ ] 修改权限配置添加CMS规则
- [ ] 修改菜单系统支持CMS菜单动态注入

**交付物**：
- 修改 `internal/svc/servicecontext.go`
- 修改 `internal/handler/routes.go`

---

### 第五阶段：前端开发（3-4周）

#### 5.1 创建CMS页面结构
- [ ] 创建CMS布局组件（CmsLayout）
- [ ] 创建顶部导航和面包屑
- [ ] 创建侧边栏菜单

**交付物**：
- `power-admin-web/src/pages/cms/CmsLayout.vue`
- `power-admin-web/src/pages/cms/components/*`

#### 5.2 内容管理模块
- [ ] 内容列表页面
  - [ ] 列表展示
  - [ ] 搜索和筛选
  - [ ] 分页
  - [ ] 批量操作
- [ ] 内容编辑页面
  - [ ] 富文本编辑器集成
  - [ ] 分类选择
  - [ ] SEO字段编辑
  - [ ] 缩略图上传
  - [ ] 预览功能
- [ ] 内容详情页面

**交付物**：
- `power-admin-web/src/pages/cms/content/*`

#### 5.3 分类管理模块
- [ ] 分类列表页面
- [ ] 分类树形展示
- [ ] 分类编辑页面
- [ ] 分类拖拽排序（可选）

**交付物**：
- `power-admin-web/src/pages/cms/category/*`

#### 5.4 评论管理模块
- [ ] 评论列表页面
- [ ] 评论审核页面
- [ ] 批量操作

**交付物**：
- `power-admin-web/src/pages/cms/comment/*`

#### 5.5 访客用户管理模块
- [ ] 用户列表页面
- [ ] 用户详情页面
- [ ] 用户禁用/启用

**交付物**：
- `power-admin-web/src/pages/cms/users/*`

#### 5.6 CMS统计面板
- [ ] 内容统计
- [ ] 评论统计
- [ ] 用户统计
- [ ] 流量统计（可选）

**交付物**：
- `power-admin-web/src/pages/cms/dashboard/*`

#### 5.7 API调用和状态管理
- [ ] 编写所有CMS API调用函数
- [ ] 创建CMS Pinia store
- [ ] 编写通用工具函数

**交付物**：
- `power-admin-web/src/api/cms.ts`
- `power-admin-web/src/stores/cms.ts`

#### 5.8 菜单和路由集成
- [ ] 动态菜单加载CMS项
- [ ] 动态路由注册
- [ ] 权限检查

**交付物**：
- 修改 `power-admin-web/src/router/index.ts`
- 修改 `power-admin-web/src/layout/Layout.vue`

---

### 第六阶段：高级功能（可选，2-3周）

#### 6.1 高级搜索和过滤
- [ ] 全文搜索
- [ ] 高级筛选条件
- [ ] 保存搜索条件

#### 6.2 内容版本管理
- [ ] 版本历史记录
- [ ] 版本对比
- [ ] 版本回滚

#### 6.3 定时发布
- [ ] 定时发布队列
- [ ] 定时任务管理

#### 6.4 分析和报告
- [ ] 内容热度分析
- [ ] 评论统计报告
- [ ] 导出功能

---

### 第七阶段：测试和优化（2-3周）

#### 7.1 单元测试
- [ ] Logic层单元测试
- [ ] Handler层单元测试
- [ ] Repository层单元测试
- [ ] 目标覆盖率: >80%

**交付物**：
- `*_test.go` 文件

#### 7.2 集成测试
- [ ] API集成测试
- [ ] 数据库集成测试
- [ ] 权限验证测试

**交付物**：
- `tests/integration/*_test.go`

#### 7.3 性能测试
- [ ] 压力测试
- [ ] 响应时间测试
- [ ] 数据库查询优化

#### 7.4 安全测试
- [ ] SQL注入测试
- [ ] XSS防护测试
- [ ] CSRF防护测试
- [ ] 权限穿透测试

#### 7.5 前端优化
- [ ] 性能优化（懒加载、虚拟滚动）
- [ ] 打包体积优化
- [ ] 缓存策略优化

#### 7.6 代码审查和重构
- [ ] 代码质量审查
- [ ] 性能瓶颈分析
- [ ] 可维护性改进

---

### 第八阶段：文档和部署（1-2周）

#### 8.1 代码文档
- [ ] 关键函数文档
- [ ] API文档（Swagger/OpenAPI）
- [ ] 架构文档
- [ ] 代码规范文档

#### 8.2 用户文档
- [ ] 管理员使用指南
- [ ] 编辑指南
- [ ] FAQ

#### 8.3 部署和上线
- [ ] 数据库脚本编写
- [ ] 数据迁移脚本
- [ ] 部署清单
- [ ] 上线测试

---

## 📊 质量标准

### 代码质量
- ✅ 代码审查通过率 100%
- ✅ 单元测试覆盖率 >80%
- ✅ 集成测试覆盖率 >70%
- ✅ 代码规范检查 0 warning

### 性能指标
- ✅ API响应时间 <200ms (不含数据库查询)
- ✅ 列表查询 <500ms
- ✅ 前端首屏加载 <3s
- ✅ 数据库查询优化到位

### 安全标准
- ✅ 无SQL注入漏洞
- ✅ 无XSS漏洞
- ✅ 无权限穿透
- ✅ 完整的审计日志

### 功能完整性
- ✅ 所有需求功能实现
- ✅ 所有边界场景处理
- ✅ 完整的错误提示
- ✅ 用户友好的交互

---

## 🎯 关键里程碑

| 里程碑 | 时间 | 交付物 | 验收标准 |
|--------|------|--------|---------|
| **M1: 数据库设计完成** | 第2周末 | 数据库脚本 | 所有表已创建 |
| **M2: 后端Logic完成** | 第5周末 | Logic层代码 | 核心业务逻辑可用 |
| **M3: 后端API完成** | 第7周末 | Handler+Router | 所有接口可测试 |
| **M4: 前端页面完成** | 第12周末 | 前端代码 | 所有页面可交互 |
| **M5: 集成测试通过** | 第14周末 | 测试报告 | 所有测试用例通过 |
| **M6: 上线就绪** | 第15周末 | 部署文档 | 可正式上线 |

---

## 🛠️ 技术栈详细说明

### 后端

**核心框架**: Go-Zero v1.9.3
- HTTP框架
- 路由管理
- 中间件支持
- 日志管理

**数据库**:
- MySQL 8.0+ (GORM ORM)
- Redis 6.0+ (缓存)

**权限管理**: Casbin v2.83.0
- RBAC权限模型
- 动态规则加载

**其他库**:
- JWT认证: github.com/golang-jwt/jwt/v4
- 密码加密: golang.org/x/crypto

### 前端

**框架**: Vue3 + TypeScript + Vite

**状态管理**: Pinia

**HTTP客户端**: axios

**UI组件库**: Element Plus (已有)

**富文本编辑**: 待选择（Tiptap/Editor.js/Quill）

**其他库**:
- 日期处理: day.js
- 工具函数: lodash-es

---

## 📝 开发规范

### 后端代码规范

```go
// 文件命名: snake_case
// 例: content_logic.go, content_handler.go

// 函数命名: PascalCase
func (l *ContentLogic) CreateContent(ctx context.Context, req *CreateContentReq) error {}

// 变量命名: camelCase
var contentId int64

// 常量命名: UPPER_SNAKE_CASE
const MAX_CONTENT_SIZE = 1024 * 1024
```

### 前端代码规范

```typescript
// 文件命名: PascalCase (组件), camelCase (其他)
// ContentList.vue, contentApi.ts

// 函数命名: camelCase
const fetchContentList = async () => {}

// 常量命名: UPPER_SNAKE_CASE
const MAX_PAGE_SIZE = 100
```

### 注释规范

```go
// 公开函数必须有注释说明
// CreateContent 创建新的内容
// @param ctx 上下文
// @param req 创建请求
// @return error 错误信息
func (l *ContentLogic) CreateContent(ctx context.Context, req *CreateContentReq) error {}
```

---

## 🔄 开发流程

### 1. 需求分析和设计
- [ ] 明确功能需求
- [ ] 设计数据模型
- [ ] 设计API接口
- [ ] 设计UI/UX

### 2. 代码实现
- [ ] 实现后端逻辑
- [ ] 编写单元测试
- [ ] 实现前端页面
- [ ] 代码审查

### 3. 测试验证
- [ ] 集成测试
- [ ] 性能测试
- [ ] 安全测试
- [ ] 用户验收

### 4. 文档和发布
- [ ] 补充文档
- [ ] 准备部署
- [ ] 发布上线

---

## 📋 依赖关系

```
数据库设计
    ↓
Model & Repository
    ↓
Business Logic (后端)
    ↓
Handler & Router (后端)
    ↓
├─ API定义 ──┬──→ 单元测试
│           │
└──────────┼───→ 前端开发
           ↓
       集成测试
           ↓
       性能优化
           ↓
       上线部署
```

---

## ✅ 最终验收标准

### 功能验收
- [ ] 所有CRUD操作正常
- [ ] 权限控制正确
- [ ] 菜单动态注入成功
- [ ] 所有页面响应速度正常

### 质量验收
- [ ] 测试覆盖率达到目标
- [ ] 代码审查通过
- [ ] 性能指标达到目标
- [ ] 无严重缺陷

### 部署验收
- [ ] 数据库脚本可执行
- [ ] 前后端可正常启动
- [ ] 所有功能可正常使用
- [ ] 日志记录完整

---

这个计划强调**质量第一，充足的时间和资源来确保系统的稳定性和完整性**。

在每个阶段完成后进行充分的测试和审查，确保代码质量。

