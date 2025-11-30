# CMS 第三阶段开发启动说明

## 当前状态

### ✅ 已完成
1. **第二阶段**：后端Logic层完全实现（1700+行代码）
   - ContentLogic、CategoryLogic、TagLogic、CommentLogic、CmsUserLogic、PublishLogic
   - 所有单元测试通过（25/25）
   - 已集成到ServiceContext中

2. **第三阶段第一步**：API定义文件完成
   - `api/cms.api` 已创建（594行）
   - 包含31个API端点定义
   - 包含34个Request/Response类型定义

### ⏳ 待完成（第三阶段）
1. **Handler实现**（6个Handler，约30-40个方法）
2. **类型定义文件**（Request/Response类型）
3. **路由注册**（在routes.go中添加CMS路由）
4. **编译验证和集成测试**

---

## 下一步工作内容

### 1️⃣ 实现ContentHandler（内容管理API处理）
- 8个方法：列表、创建、获取、更新、删除、发布、取消发布、批量操作
- 调用ContentLogic进行业务逻辑处理
- 完整的参数验证和错误处理

### 2️⃣ 实现CategoryHandler（分类管理API处理）
- 6个方法：列表、树形结构、创建、获取、更新、删除
- 支持多级分类树形结构
- 分类存在性验证

### 3️⃣ 实现其他Handlers
- TagHandler（5个方法）
- CommentHandler（8个方法）
- CmsUserHandler（8个方法）
- PublishHandler（5个方法）

### 4️⃣ 注册路由
- 在routes.go中添加CMS路由组
- 配置JWT验证中间件
- 配置权限验证中间件

### 5️⃣ 集成测试
- API端点功能测试
- 参数验证测试
- 错误处理测试

---

## 技术细节

### Handler基本结构
```go
package cms

import (
    "net/http"
    "power-admin-server/internal/svc"
    "power-admin-server/internal/types"
    "encoding/json"
    "github.com/zeromicro/go-zero/core/logx"
)

// ContentListHandler 获取内容列表
func ContentListHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. 解析请求参数
        var req types.ContentListReq
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            logx.Errorf("decode request failed: %v", err)
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }
        
        // 2. 参数验证
        if req.Page < 1 {
            req.Page = 1
        }
        if req.PageSize < 1 || req.PageSize > 100 {
            req.PageSize = 10
        }
        
        // 3. 调用Logic层
        result, err := serverCtx.CmsContentLogic.ListContent(r.Context(), &cms.ListContentRequest{
            Page:       req.Page,
            PageSize:   req.PageSize,
            CategoryID: // 转换categoryId
            Status:     // 转换status
            // ... 其他参数
        })
        
        if err != nil {
            logx.Errorf("list content failed: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // 4. 返回响应
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    }
}
```

### 路由注册模式
```go
// 在routes.go的RegisterHandlers函数中添加
server.AddRoutes(
    rest.WithMiddlewares(
        []rest.Middleware{serverCtx.AdminAuthMiddleware},
        []rest.Route{
            {
                Method:  http.MethodGet,
                Path:    "/cms/content",
                Handler: cms.ContentListHandler(serverCtx),
            },
            {
                Method:  http.MethodPost,
                Path:    "/cms/content",
                Handler: cms.CreateContentHandler(serverCtx),
            },
            // ... 更多路由
        }...,
    ),
    rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
    rest.WithPrefix("/api"),
)
```

---

## 工作流程

### 每个Handler的开发流程
1. **分析Logic接口** - 理解Logic层提供的方法和参数
2. **定义Request/Response** - 在types中定义API的请求和响应格式
3. **实现Handler** - 实现HTTP handler函数
4. **参数映射** - 将HTTP请求参数映射到Logic的参数
5. **错误处理** - 完整的error handling和logging
6. **响应处理** - 序列化并返回JSON响应
7. **单元测试** - 编写handler的单元测试
8. **集成测试** - 进行API端到端测试

---

## 建议的实现优先级

### 推荐顺序（按优先级）
1. **ContentHandler** ⭐⭐⭐ - 核心功能，复杂度最高
2. **CategoryHandler** ⭐⭐ - 支撑模块，中等复杂度
3. **CmsUserHandler** ⭐⭐ - 用户认证，优先级高
4. **CommentHandler** ⭐ - 功能完整性，中等复杂度
5. **TagHandler** ⭐ - 简单实现，低复杂度
6. **PublishHandler** ⭐ - 工作流，中等复杂度

### 每天可完成
- 日均1-2个Handler的实现和测试
- 预计3-5工作天完成全部Handler

---

## 质量检查清单

### 代码质量
- [ ] 所有handler都有错误处理
- [ ] 所有参数都有验证
- [ ] 所有操作都有日志记录
- [ ] 遵循Go编码规范
- [ ] 没有重复代码

### API规范
- [ ] 请求和响应格式一致
- [ ] 使用合适的HTTP方法和状态码
- [ ] 错误信息清晰明确
- [ ] 支持Content negotiation

### 功能完整性
- [ ] 所有CRUD操作都有实现
- [ ] 特殊业务逻辑都有实现
- [ ] 边界情况都有处理
- [ ] 数据验证完善

### 安全性
- [ ] JWT验证到位
- [ ] 权限检查完善
- [ ] 输入参数验证
- [ ] 防止SQL注入（通过ORM）

---

## 参考资源

### 现有Handler参考
- `internal/handler/user/` - 用户管理Handler
- `internal/handler/role/` - 角色管理Handler
- `internal/handler/menu/` - 菜单管理Handler

### API定义参考
- `api/admin.api` - 现有API定义格式

### Logic层参考
- `CMS_LOGIC_API_REFERENCE.md` - Logic层API完整参考

---

## 预期成果

完成第三阶段后，系统将具备：
- ✅ 31个完整的API端点
- ✅ 6个完整的Handler模块
- ✅ 100%的Logic层集成
- ✅ JWT和权限验证
- ✅ 完整的错误处理
- ✅ 完整的API文档（通过cms.api）

---

## 接下来的步骤

1. **当前**：已完成第三阶段第一步（API定义）
2. **下一步**：开始实现Handler层
   - 推荐从ContentHandler开始
   - 按优先级逐个实现
   - 每个Handler完成后进行单元和集成测试

3. **后续**：
   - 第四阶段：前端开发（Vue3 CMS管理界面）
   - 第五阶段：性能优化和集成测试
   - 第六阶段：文档编写和部署

---

## 开发建议

### 代码组织
- 每个Handler类型创建一个文件
- 相关的helper函数放在同一文件
- 类型定义集中在types/cms_types.go

### 测试策略
- 先写Handler逻辑，再补充单元测试
- 每个Handler完成后立即进行集成测试
- 使用Postman或curl进行API测试

### 代码审查
- 遵循Go标准库约定
- 优先使用现有工具和包
- 保持代码风格一致

---

## 任何问题？

如果在开发过程中遇到问题，请参考：
1. `CMS_LOGIC_API_REFERENCE.md` - Logic层API参考
2. `CMS_PHASE3_PLAN.md` - 第三阶段详细计划
3. 现有Handler实现 - 学习go-zero pattern

---

**准备就绪！** 🚀 

现在可以开始实现Handler层了。建议从ContentHandler开始，然后按优先级逐个实现其他Handler。

**预计时间**：9-13.5小时完成所有Handler实现和测试

---

**文档版本**：1.0
**更新日期**：2025-11-30
**状态**：Ready to Start
