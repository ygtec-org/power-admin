# CMS Logic API 快速参考指南

本文档为Handler层开发提供快速参考，列举所有可用的Logic方法及其签名。

---

## 目录
1. [ContentLogic](#contentlogic)
2. [CategoryLogic](#categorylogic)
3. [TagLogic](#taglogic)
4. [CommentLogic](#commentlogic)
5. [CmsUserLogic](#cmsユlogic)
6. [PublishLogic](#publishlogic)
7. [使用示例](#使用示例)

---

## ContentLogic

### 创建内容
```go
func (l *ContentLogic) CreateContent(ctx context.Context, req *CreateContentRequest) (*models.CmsContent, error)
```
**请求字段**：Title(必填)、Content(必填)、AuthorID(必填)、CategoryID、Slug、Status、Visibility、CommentStatus、SEO字段等

### 更新内容
```go
func (l *ContentLogic) UpdateContent(ctx context.Context, req *UpdateContentRequest) (*models.CmsContent, error)
```

### 删除内容
```go
func (l *ContentLogic) DeleteContent(ctx context.Context, contentID int64) error  // 软删除
func (l *ContentLogic) HardDeleteContent(ctx context.Context, contentID int64) error  // 硬删除
```

### 查询内容
```go
func (l *ContentLogic) GetContent(ctx context.Context, contentID int64) (*models.CmsContent, error)
func (l *ContentLogic) GetContentBySlug(ctx context.Context, slug string) (*models.CmsContent, error)
func (l *ContentLogic) ListContent(ctx context.Context, req *ListContentRequest) (*repository.PagedResult, error)
```

### 发布操作
```go
func (l *ContentLogic) PublishContent(ctx context.Context, contentID int64) error
func (l *ContentLogic) UnpublishContent(ctx context.Context, contentID int64) error
func (l *ContentLogic) BatchUpdateContentStatus(ctx context.Context, ids []int64, status int8) error
```

---

## CategoryLogic

### 创建分类
```go
func (l *CategoryLogic) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*models.CmsCategory, error)
```

### 更新分类
```go
func (l *CategoryLogic) UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*models.CmsCategory, error)
```

### 删除分类
```go
func (l *CategoryLogic) DeleteCategory(ctx context.Context, categoryID int64) error
```

### 查询分类
```go
func (l *CategoryLogic) GetCategory(ctx context.Context, categoryID int64) (*models.CmsCategory, error)
func (l *CategoryLogic) GetCategoryBySlug(ctx context.Context, slug string) (*models.CmsCategory, error)
func (l *CategoryLogic) ListCategories(ctx context.Context, parentID *int64) ([]*models.CmsCategory, error)
func (l *CategoryLogic) GetCategoryTree(ctx context.Context) ([]*models.CmsCategory, error)  // 返回树形结构
```

---

## TagLogic

### CRUD操作
```go
func (l *TagLogic) CreateTag(ctx context.Context, req *CreateTagRequest) (*models.CmsTag, error)
func (l *TagLogic) UpdateTag(ctx context.Context, req *UpdateTagRequest) (*models.CmsTag, error)
func (l *TagLogic) DeleteTag(ctx context.Context, tagID int64) error
func (l *TagLogic) GetTag(ctx context.Context, tagID int64) (*models.CmsTag, error)
```

### 查询操作
```go
func (l *TagLogic) GetTagByName(ctx context.Context, name string) (*models.CmsTag, error)
func (l *TagLogic) ListTags(ctx context.Context) ([]*models.CmsTag, error)
func (l *TagLogic) GetTagsByIDs(ctx context.Context, tagIDs []int64) ([]*models.CmsTag, error)
```

### 使用数管理
```go
func (l *TagLogic) IncrementTagUsage(ctx context.Context, tagID int64) error
func (l *TagLogic) DecrementTagUsage(ctx context.Context, tagID int64) error
func (l *TagLogic) BatchGetOrCreateTags(ctx context.Context, tagNames []string) ([]*models.CmsTag, error)
```

---

## CommentLogic

### CRUD操作
```go
func (l *CommentLogic) CreateComment(ctx context.Context, req *CreateCommentRequest) (*models.CmsComment, error)
func (l *CommentLogic) UpdateComment(ctx context.Context, req *UpdateCommentRequest) (*models.CmsComment, error)
func (l *CommentLogic) DeleteComment(ctx context.Context, commentID int64) error
func (l *CommentLogic) GetComment(ctx context.Context, commentID int64) (*models.CmsComment, error)
```

### 查询操作
```go
func (l *CommentLogic) ListComments(ctx context.Context, contentID int64) ([]*models.CmsComment, error)
```

### 审核操作
```go
func (l *CommentLogic) ApproveComment(ctx context.Context, commentID int64) error  // 状态1
func (l *CommentLogic) RejectComment(ctx context.Context, commentID int64) error   // 状态2
func (l *CommentLogic) SpamComment(ctx context.Context, commentID int64) error     // 状态3
```

### 交互操作
```go
func (l *CommentLogic) LikeComment(ctx context.Context, commentID int64) error
func (l *CommentLogic) UnlikeComment(ctx context.Context, commentID int64) error
func (l *CommentLogic) ReplyComment(ctx context.Context, parentCommentID int64, req *CreateCommentRequest) (*models.CmsComment, error)
```

---

## CmsUserLogic

### 认证操作
```go
func (l *CmsUserLogic) Register(ctx context.Context, req *RegisterRequest) (*models.CmsUser, error)
func (l *CmsUserLogic) Login(ctx context.Context, req *LoginRequest) (*models.CmsUser, error)
func (l *CmsUserLogic) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error
```

### CRUD操作
```go
func (l *CmsUserLogic) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*models.CmsUser, error)
func (l *CmsUserLogic) GetUser(ctx context.Context, userID int64) (*models.CmsUser, error)
func (l *CmsUserLogic) ListUsers(ctx context.Context) ([]*models.CmsUser, error)
```

### 状态管理
```go
func (l *CmsUserLogic) DisableUser(ctx context.Context, userID int64) error
func (l *CmsUserLogic) EnableUser(ctx context.Context, userID int64) error
func (l *CmsUserLogic) DeleteUser(ctx context.Context, userID int64) error  // 软删除
func (l *CmsUserLogic) HardDeleteUser(ctx context.Context, userID int64) error  // 硬删除
```

### 验证操作
```go
func (l *CmsUserLogic) VerifyEmail(ctx context.Context, userID int64) error
func (l *CmsUserLogic) VerifyPhone(ctx context.Context, userID int64) error
```

---

## PublishLogic

### 发布操作
```go
func (l *PublishLogic) PublishImmediate(ctx context.Context, req *PublishImmediateRequest) error
func (l *PublishLogic) PublishScheduled(ctx context.Context, req *PublishScheduledRequest) error
func (l *PublishLogic) Unpublish(ctx context.Context, req *UnpublishRequest) error
func (l *PublishLogic) CancelScheduledPublish(ctx context.Context, req *CancelScheduledPublishRequest) error
```

### 批量操作
```go
func (l *PublishLogic) BatchPublish(ctx context.Context, req *BatchPublishRequest) error
func (l *PublishLogic) BatchUnpublish(ctx context.Context, req *BatchUnpublishRequest) error
```

### 查询操作
```go
func (l *PublishLogic) GetPublishStatus(ctx context.Context, contentID int64) (*PublishStatus, error)
func (l *PublishLogic) GetScheduledContents(ctx context.Context) ([]*models.CmsContent, error)
```

### 后台任务
```go
func (l *PublishLogic) ProcessScheduledPublish(ctx context.Context) error  // 供后台定时任务调用
```

---

## 使用示例

### 在Handler中使用Logic

#### 创建内容API
```go
func (h *ContentHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req cms.CreateContentRequest
    
    // 解析请求
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // 调用Logic层
    content, err := h.svc.CmsContentLogic.CreateContent(r.Context(), &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 返回响应
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(content)
}
```

#### 发布内容API
```go
func (h *ContentHandler) Publish(w http.ResponseWriter, r *http.Request) {
    contentID := parseID(r.PathValue("id"))
    
    // 调用Logic层
    err := h.svc.CmsContentLogic.PublishContent(r.Context(), contentID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 返回成功
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
```

#### 列表查询API
```go
func (h *ContentHandler) List(w http.ResponseWriter, r *http.Request) {
    req := &cms.ListContentRequest{
        Page:     parseInt(r.Query("page"), 1),
        PageSize: parseInt(r.Query("pageSize"), 10),
    }
    
    if categoryID := r.Query("categoryId"); categoryID != "" {
        id := parseInt(categoryID, 0)
        req.CategoryID = &id
    }
    
    // 调用Logic层
    result, err := h.svc.CmsContentLogic.ListContent(r.Context(), req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 返回结果
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

---

## 注意事项

### 错误处理
所有Logic方法都返回error，需要检查和处理：
```go
if err != nil {
    // 记录错误日志
    // 返回适当的HTTP状态码
    // 返回用户友好的错误信息
}
```

### 密码安全
- 用户密码使用bcrypt加密
- 不要返回密码字段给前端
- Login方法会自动清除返回的密码

### 状态值
- 内容状态：1=草稿, 2=已发布, 3=已删除
- 评论状态：0=待审核, 1=已批准, 2=已拒绝, 3=垃圾
- 用户状态：0=禁用, 1=正常

### 上下文处理
- 始终传递request的context给Logic方法
- 支持context超时控制

---

## 下一步

1. **创建API定义文件**（api/cms.api）
2. **实现Handler类**
3. **注册路由**
4. **编写集成测试**

---

**更新日期**：2025-11-30
**版本**：1.0
