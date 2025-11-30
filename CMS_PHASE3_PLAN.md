# CMS 第三阶段开发计划 - Handler和API层实现

## 概述
第三阶段专注于实现Handler层和完整的API接口。我们将基于第二阶段完成的Logic层，为所有CMS模块创建Handler，完成API的全部实现。

---

## 开发阶段详解

### 阶段1：API定义文件 ✅ 已完成
**文件**：`api/cms.api`（594行）

**内容**：
- 所有Request和Response类型定义（34个类型）
- 所有API端点定义（31个接口）
- API路径和HTTP方法规范定义

**特点**：
- 遵循go-zero标准api语法
- 包含完整的字段验证和默认值
- 支持path参数、query参数、body参数

---

### 阶段2：Handler实现

#### 2.1 ContentHandler（内容管理）
**目录**：`internal/handler/cms/content.go`

**需要实现的方法**：
1. `ContentListHandler` - 列表查询（分页、排序、搜索、筛选）
2. `CreateContentHandler` - 创建内容（参数验证、分类验证）
3. `GetContentHandler` - 获取详情（自动增加浏览数）
4. `UpdateContentHandler` - 更新内容（原数据验证）
5. `DeleteContentHandler` - 删除内容（软删除）
6. `PublishContentHandler` - 发布内容
7. `UnpublishContentHandler` - 取消发布
8. `BatchUpdateStatusHandler` - 批量更新状态

**Handler实现模式**：
```go
func ContentListHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.ContentListReq
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            // 错误处理
        }
        
        // 调用Logic层
        result, err := serverCtx.CmsContentLogic.ListContent(r.Context(), &cms.ListContentRequest{
            Page: req.Page,
            PageSize: req.PageSize,
            // ... 其他参数映射
        })
        
        if err != nil {
            // 错误处理和日志
        }
        
        // 响应序列化
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    }
}
```

#### 2.2 CategoryHandler（分类管理）
**目录**：`internal/handler/cms/category.go`

**需要实现的方法**：
1. `CategoryListHandler` - 列表查询
2. `CategoryTreeHandler` - 树形结构查询
3. `CreateCategoryHandler` - 创建分类
4. `GetCategoryHandler` - 获取详情
5. `UpdateCategoryHandler` - 更新分类
6. `DeleteCategoryHandler` - 删除分类

#### 2.3 TagHandler（标签管理）
**目录**：`internal/handler/cms/tag.go`

**需要实现的方法**：
1. `TagListHandler` - 列表查询
2. `CreateTagHandler` - 创建标签
3. `GetTagHandler` - 获取详情
4. `UpdateTagHandler` - 更新标签
5. `DeleteTagHandler` - 删除标签

#### 2.4 CommentHandler（评论管理）
**目录**：`internal/handler/cms/comment.go`

**需要实现的方法**：
1. `CommentListHandler` - 列表查询
2. `CreateCommentHandler` - 创建评论
3. `GetCommentHandler` - 获取详情
4. `UpdateCommentHandler` - 更新评论
5. `DeleteCommentHandler` - 删除评论
6. `ApproveCommentHandler` - 审核通过
7. `RejectCommentHandler` - 拒绝评论
8. `LikeCommentHandler` - 点赞评论

#### 2.5 CmsUserHandler（访客用户管理）
**目录**：`internal/handler/cms/user.go`

**需要实现的方法**：
1. `CmsUserListHandler` - 列表查询
2. `CmsRegisterHandler` - 用户注册
3. `CmsLoginHandler` - 用户登录
4. `GetCmsUserHandler` - 获取用户详情
5. `UpdateCmsUserHandler` - 更新用户信息
6. `ChangePasswordHandler` - 修改密码
7. `DisableCmsUserHandler` - 禁用用户
8. `DeleteCmsUserHandler` - 删除用户

#### 2.6 PublishHandler（发布和工作流）
**目录**：`internal/handler/cms/publish.go`

**需要实现的方法**：
1. `PublishImmediateHandler` - 立即发布
2. `PublishScheduledHandler` - 定时发布
3. `GetPublishStatusHandler` - 获取发布状态
4. `CancelScheduledPublishHandler` - 取消定时发布
5. `BatchPublishHandler` - 批量发布

---

### 阶段3：类型定义文件

**文件**：`internal/types/cms_types.go`

**内容**：
- 从cms.api自动生成的所有Request和Response类型
- 或手动创建与api一致的类型定义

---

### 阶段4：路由注册

**修改文件**：`internal/handler/routes.go`

**需要添加的代码**：
```go
// CMS路由注册
server.AddRoutes(
    rest.WithMiddlewares(
        []rest.Middleware{serverCtx.AdminAuthMiddleware},
        []rest.Route{
            // 内容路由
            { Method: http.MethodGet, Path: "/cms/content", Handler: cms.ContentListHandler(serverCtx) },
            { Method: http.MethodPost, Path: "/cms/content", Handler: cms.CreateContentHandler(serverCtx) },
            // ... 其他路由
        }...,
    ),
    rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
    rest.WithPrefix("/api"),
)
```

---

## 实现顺序建议

### 第一批（优先级高）
1. **ContentHandler** - 核心内容管理功能
2. **CategoryHandler** - 支撑内容分类
3. **CmsUserHandler** - 用户认证和管理

### 第二批（优先级中）
4. **TagHandler** - 内容标签管理
5. **CommentHandler** - 评论管理
6. **PublishHandler** - 发布工作流

### 第三批（优先级低但必要）
7. **类型定义** - 完整的Request/Response类型
8. **路由注册** - 所有路由的注册和中间件
9. **编译验证** - 全项目编译检验

---

## 技术要点

### 错误处理模式
```go
if err != nil {
    logx.Errorf("operation failed: %v", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
```

### 响应格式
```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

### 参数验证
- 必填字段检查
- 数据类型验证
- 业务逻辑验证（如分类存在、用户存在等）

### 日志记录
- 所有操作都应记录日志
- 错误日志记录详细信息
- 使用logx.Logger

### 中间件使用
- JWT验证：`rest.WithJwt()`
- 自定义中间件：`rest.WithMiddlewares()`
- 权限验证中间件

---

## 文件结构

```
internal/handler/
└── cms/
    ├── content.go      # ContentHandler实现
    ├── category.go     # CategoryHandler实现
    ├── tag.go          # TagHandler实现
    ├── comment.go      # CommentHandler实现
    ├── user.go         # CmsUserHandler实现
    ├── publish.go      # PublishHandler实现
    └── types.go        # CMS相关类型定义（可选）

internal/types/
├── cms_types.go        # CMS API的Request/Response类型

api/
└── cms.api            # CMS API定义（已完成）
```

---

## 测试计划

### 单元测试
- Handler参数验证测试
- Handler错误处理测试
- Handler正常流程测试

### 集成测试
- 完整的API调用流程
- 数据库持久化验证
- 事务回滚验证

### API测试
- 使用Postman或curl测试所有API端点
- 验证响应格式
- 验证错误处理

---

## 质量标准

### 代码质量
- ✅ 所有handler都应该有错误处理
- ✅ 所有参数都应该有验证
- ✅ 所有操作都应该有日志
- ✅ 遵循Go编码规范

### API响应
- ✅ 统一的响应格式
- ✅ 恰当的HTTP状态码
- ✅ 清晰的错误消息
- ✅ 正确的Content-Type

### 安全性
- ✅ JWT验证
- ✅ 权限检查
- ✅ 输入验证
- ✅ SQL注入防护（通过ORM）

---

## 时间估计

| 模块 | 时间估计 | 难度 |
|------|---------|------|
| ContentHandler | 2-3小时 | 中等 |
| CategoryHandler | 1.5-2小时 | 中等 |
| TagHandler | 1-1.5小时 | 简单 |
| CommentHandler | 1.5-2小时 | 中等 |
| CmsUserHandler | 1.5-2小时 | 中等 |
| PublishHandler | 1-1.5小时 | 简单 |
| 类型定义 | 0.5小时 | 简单 |
| 路由注册 | 0.5-1小时 | 简单 |
| **总计** | **9-13.5小时** | - |

---

## 下一阶段展望

完成第三阶段后，将进入：
- **第四阶段**：前端开发（Vue3管理界面）
- **第五阶段**：集成测试和性能优化
- **第六阶段**：文档和部署

---

**开始日期**：2025-11-30
**计划完成日期**：待更新
**版本**：1.0
