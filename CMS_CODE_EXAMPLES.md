# CMS插件 - 核心代码示例集合

> 本文档包含所有关键组件的代码示例，开发者可以复制使用或作为参考。

---

## 📝 目录

1. [后端代码示例](#后端代码示例)
2. [前端代码示例](#前端代码示例)
3. [数据库查询](#数据库查询)
4. [权限配置](#权限配置)

---

## 后端代码示例

### 1. API定义文件 (api/cms.api)

```api
syntax = "v1"

info(
    title: "CMS 内容管理系统"
    desc: "可插拔的内容管理系统API"
    author: "Your Team"
    version: "1.0.0"
)

// ===== 内容管理 =====
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
    
    // ===== 分类管理 =====
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
}

// ===== 请求类型 =====
type ListContentsReq {
    Page     int    `form:"page,default=1"`
    PageSize int    `form:"pageSize,default=10"`
    Category int64  `form:"category_id,optional"`
    Status   int    `form:"status,optional"`
    Search   string `form:"search,optional"`
}

type ListContentsResp {
    Code  int           `json:"code"`
    Msg   string        `json:"msg"`
    Data  ListContentsData `json:"data"`
}

type ListContentsData {
    List  []CmsContent `json:"list"`
    Total int64        `json:"total"`
}

type CmsContent {
    Id          int64  `json:"id"`
    Title       string `json:"title"`
    Content     string `json:"content"`
    Excerpt     string `json:"excerpt"`
    CategoryId  int64  `json:"category_id"`
    Status      int    `json:"status"`
    ViewCount   int    `json:"view_count"`
    CreatedAt   string `json:"created_at"`
    PublishedAt string `json:"published_at"`
}

type DetailReq {
    Id int64 `path:"id"`
}

type CreateContentReq {
    Title      string `json:"title"`
    Content    string `json:"content"`
    Excerpt    string `json:"excerpt,optional"`
    CategoryId int64  `json:"category_id,optional"`
}

type UpdateContentReq {
    Title      string `json:"title"`
    Content    string `json:"content"`
    Excerpt    string `json:"excerpt,optional"`
    CategoryId int64  `json:"category_id,optional"`
}

type CommonResp {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}
```

### 2. Logic层实现 (internal/logic/cms/cmscontentlogic.go)

```go
package cms

import (
    "context"
    "fmt"
    "log"
    "power-admin/internal/svc"
    "power-admin/internal/types"
)

type CmsContentLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewCmsContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CmsContentLogic {
    return &CmsContentLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

// List 获取内容列表
func (l *CmsContentLogic) List(req *types.ListContentsReq) (*types.ListContentsResp, error) {
    // 1. 参数验证
    if req.Page <= 0 {
        req.Page = 1
    }
    if req.PageSize <= 0 || req.PageSize > 100 {
        req.PageSize = 10
    }
    
    // 2. 构建查询条件
    query := "SELECT id, title, content, excerpt, category_id, status, view_count, created_at, published_at FROM cms_content WHERE 1=1"
    var args []interface{}
    
    if req.Status > 0 {
        query += " AND status = ?"
        args = append(args, req.Status)
    } else {
        // 默认只显示已发布的内容（管理后台显示草稿）
        query += " AND status != 3"
    }
    
    if req.Category > 0 {
        query += " AND category_id = ?"
        args = append(args, req.Category)
    }
    
    if req.Search != "" {
        query += " AND (title LIKE ? OR content LIKE ?)"
        searchTerm := "%" + req.Search + "%"
        args = append(args, searchTerm, searchTerm)
    }
    
    // 3. 统计总数
    countQuery := "SELECT COUNT(*) FROM cms_content WHERE 1=1"
    countQuery = query[:len(query)-len(fmt.Sprintf("SELECT id, title, content, excerpt, category_id, status, view_count, created_at, published_at FROM cms_content"))] 
    // (简化，实际应该单独统计)
    
    var total int64
    l.svcCtx.DB.QueryRowContext(l.ctx, countQuery, args...).Scan(&total)
    
    // 4. 分页和排序
    offset := (req.Page - 1) * req.PageSize
    query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
    args = append(args, req.PageSize, offset)
    
    // 5. 查询数据
    rows, err := l.svcCtx.DB.QueryContext(l.ctx, query, args...)
    if err != nil {
        log.Printf("query error: %v", err)
        return nil, fmt.Errorf("查询失败")
    }
    defer rows.Close()
    
    var list []types.CmsContent
    for rows.Next() {
        var content types.CmsContent
        err := rows.Scan(
            &content.Id,
            &content.Title,
            &content.Content,
            &content.Excerpt,
            &content.CategoryId,
            &content.Status,
            &content.ViewCount,
            &content.CreatedAt,
            &content.PublishedAt,
        )
        if err != nil {
            log.Printf("scan error: %v", err)
            continue
        }
        list = append(list, content)
    }
    
    return &types.ListContentsResp{
        Code: 0,
        Msg:  "success",
        Data: types.ListContentsData{
            List:  list,
            Total: total,
        },
    }, nil
}

// Create 创建内容
func (l *CmsContentLogic) Create(req *types.CreateContentReq, authorId int64) error {
    // 1. 验证输入
    if err := l.validateCreateRequest(req); err != nil {
        return err
    }
    
    // 2. 插入数据
    query := `INSERT INTO cms_content (title, content, excerpt, category_id, author_id, status, created_at) 
              VALUES (?, ?, ?, ?, ?, 1, NOW())`
    
    result, err := l.svcCtx.DB.ExecContext(l.ctx,
        query,
        req.Title,
        req.Content,
        req.Excerpt,
        req.CategoryId,
        authorId,
    )
    
    if err != nil {
        log.Printf("insert error: %v", err)
        return fmt.Errorf("创建失败")
    }
    
    id, err := result.LastInsertId()
    log.Printf("Created content with ID: %d", id)
    
    return nil
}

// Update 更新内容
func (l *CmsContentLogic) Update(id int64, req *types.UpdateContentReq) error {
    // 1. 验证输入
    if err := l.validateCreateRequest(&types.CreateContentReq{
        Title:      req.Title,
        Content:    req.Content,
        Excerpt:    req.Excerpt,
        CategoryId: req.CategoryId,
    }); err != nil {
        return err
    }
    
    // 2. 检查内容是否存在
    var count int
    l.svcCtx.DB.QueryRowContext(l.ctx, "SELECT COUNT(*) FROM cms_content WHERE id = ?", id).Scan(&count)
    if count == 0 {
        return fmt.Errorf("内容不存在")
    }
    
    // 3. 更新数据
    query := `UPDATE cms_content SET title = ?, content = ?, excerpt = ?, category_id = ?, updated_at = NOW() WHERE id = ?`
    
    _, err := l.svcCtx.DB.ExecContext(l.ctx,
        query,
        req.Title,
        req.Content,
        req.Excerpt,
        req.CategoryId,
        id,
    )
    
    if err != nil {
        log.Printf("update error: %v", err)
        return fmt.Errorf("更新失败")
    }
    
    return nil
}

// Delete 删除内容（逻辑删除）
func (l *CmsContentLogic) Delete(id int64) error {
    // 1. 检查内容是否存在
    var count int
    l.svcCtx.DB.QueryRowContext(l.ctx, "SELECT COUNT(*) FROM cms_content WHERE id = ?", id).Scan(&count)
    if count == 0 {
        return fmt.Errorf("内容不存在")
    }
    
    // 2. 逻辑删除（将status设为3）
    query := "UPDATE cms_content SET status = 3, updated_at = NOW() WHERE id = ?"
    
    _, err := l.svcCtx.DB.ExecContext(l.ctx, query, id)
    if err != nil {
        log.Printf("delete error: %v", err)
        return fmt.Errorf("删除失败")
    }
    
    return nil
}

// Publish 发布内容
func (l *CmsContentLogic) Publish(id int64) error {
    query := "UPDATE cms_content SET status = 2, published_at = NOW(), updated_at = NOW() WHERE id = ?"
    
    _, err := l.svcCtx.DB.ExecContext(l.ctx, query, id)
    if err != nil {
        log.Printf("publish error: %v", err)
        return fmt.Errorf("发布失败")
    }
    
    return nil
}

// validateCreateRequest 验证创建请求
func (l *CmsContentLogic) validateCreateRequest(req *types.CreateContentReq) error {
    if req.Title == "" {
        return fmt.Errorf("标题不能为空")
    }
    if len(req.Title) > 255 {
        return fmt.Errorf("标题长度不能超过255个字符")
    }
    if req.Content == "" {
        return fmt.Errorf("内容不能为空")
    }
    return nil
}
```

### 3. Handler层实现 (internal/handler/cms/cmscontenthandler.go)

```go
package cms

import (
    "encoding/json"
    "net/http"
    "strconv"
    "power-admin/internal/logic/cms"
    "power-admin/internal/svc"
    "power-admin/internal/types"
)

// CmsContentListHandler 获取内容列表
func CmsContentListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. 权限检查
        userId := r.Header.Get("X-User-ID")
        if userId == "" {
            writeError(w, http.StatusUnauthorized, "未授权")
            return
        }
        
        // 检查Casbin权限
        allowed, err := svcCtx.Enforcer.Enforce(userId, "/api/cms/admin/contents", "GET")
        if err != nil || !allowed {
            writeError(w, http.StatusForbidden, "无权限访问")
            return
        }
        
        // 2. 解析参数
        var req types.ListContentsReq
        page := r.URL.Query().Get("page")
        pageSize := r.URL.Query().Get("pageSize")
        
        if page != "" {
            req.Page, _ = strconv.Atoi(page)
        } else {
            req.Page = 1
        }
        
        if pageSize != "" {
            req.PageSize, _ = strconv.Atoi(pageSize)
        } else {
            req.PageSize = 10
        }
        
        req.Category, _ = strconv.ParseInt(r.URL.Query().Get("category_id"), 10, 64)
        req.Search = r.URL.Query().Get("search")
        
        // 3. 调用Logic
        l := cms.NewCmsContentLogic(r.Context(), svcCtx)
        resp, err := l.List(&req)
        
        if err != nil {
            writeError(w, http.StatusInternalServerError, err.Error())
            return
        }
        
        // 4. 返回响应
        writeJSON(w, http.StatusOK, resp)
    }
}

// CmsContentCreateHandler 创建内容
func CmsContentCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. 权限检查
        userId := r.Header.Get("X-User-ID")
        allowed, _ := svcCtx.Enforcer.Enforce(userId, "/api/cms/admin/contents", "POST")
        if !allowed {
            writeError(w, http.StatusForbidden, "无权限访问")
            return
        }
        
        // 2. 解析请求体
        var req types.CreateContentReq
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeError(w, http.StatusBadRequest, "请求格式错误")
            return
        }
        
        // 3. 调用Logic
        l := cms.NewCmsContentLogic(r.Context(), svcCtx)
        authorId, _ := strconv.ParseInt(userId, 10, 64)
        err := l.Create(&req, authorId)
        
        if err != nil {
            writeError(w, http.StatusInternalServerError, err.Error())
            return
        }
        
        // 4. 返回成功
        writeJSON(w, http.StatusOK, types.CommonResp{
            Code: 0,
            Msg:  "创建成功",
        })
    }
}

// 辅助函数
func writeJSON(w http.ResponseWriter, code int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, code int, msg string) {
    writeJSON(w, code, map[string]interface{}{
        "code": code,
        "msg":  msg,
    })
}
```

---

## 前端代码示例

### 1. API调用 (src/api/cms.ts)

```typescript
import request from './request'

// ===== 内容管理 =====
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

export const publishContent = (id: number) =>
  request.post(`/api/cms/admin/publish`, { id })

// ===== 分类管理 =====
export const getCategoryTree = () =>
  request.get('/api/cms/admin/categories/tree')

export const getCategoryList = (params: any) =>
  request.get('/api/cms/admin/categories', { params })

export const createCategory = (data: any) =>
  request.post('/api/cms/admin/categories', data)

export const updateCategory = (id: number, data: any) =>
  request.put(`/api/cms/admin/categories/${id}`, data)

export const deleteCategory = (id: number) =>
  request.delete(`/api/cms/admin/categories/${id}`)
```

### 2. 状态管理 (src/stores/cms.ts)

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as cmsApi from '@/api/cms'
import { ElMessage } from 'element-plus'

export const useCmsStore = defineStore('cms', () => {
  // ===== 状态 =====
  const contentList = ref<any[]>([])
  const categoryList = ref<any[]>([])
  const categoryTree = ref<any[]>([])
  const currentPage = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const loading = ref(false)
  
  // ===== 计算属性 =====
  const contentCount = computed(() => contentList.value.length)
  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))
  
  // ===== 内容管理 =====
  async function fetchContentList(params?: any) {
    loading.value = true
    try {
      const res = await cmsApi.getContentList({
        page: currentPage.value,
        pageSize: pageSize.value,
        ...params
      })
      
      if (res.data.code === 0) {
        contentList.value = res.data.data.list || []
        total.value = res.data.data.total || 0
      }
      return res
    } catch (error) {
      console.error('Failed to fetch contents:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  async function createContent(data: any) {
    try {
      await cmsApi.createContent(data)
      await fetchContentList()
      return true
    } catch (error) {
      console.error('Failed to create content:', error)
      throw error
    }
  }
  
  async function updateContent(id: number, data: any) {
    try {
      await cmsApi.updateContent(id, data)
      await fetchContentList()
      return true
    } catch (error) {
      console.error('Failed to update content:', error)
      throw error
    }
  }
  
  async function deleteContent(id: number) {
    try {
      await cmsApi.deleteContent(id)
      await fetchContentList()
      return true
    } catch (error) {
      console.error('Failed to delete content:', error)
      throw error
    }
  }
  
  async function publishContent(id: number) {
    try {
      await cmsApi.publishContent(id)
      await fetchContentList()
      return true
    } catch (error) {
      console.error('Failed to publish content:', error)
      throw error
    }
  }
  
  // ===== 分类管理 =====
  async function fetchCategoryTree() {
    try {
      const res = await cmsApi.getCategoryTree()
      if (res.data.code === 0) {
        categoryTree.value = res.data.data || []
      }
      return res
    } catch (error) {
      console.error('Failed to fetch category tree:', error)
      throw error
    }
  }
  
  async function fetchCategoryList(params?: any) {
    try {
      const res = await cmsApi.getCategoryList({
        page: 1,
        pageSize: 100,
        ...params
      })
      
      if (res.data.code === 0) {
        categoryList.value = res.data.data.list || []
      }
      return res
    } catch (error) {
      console.error('Failed to fetch categories:', error)
      throw error
    }
  }
  
  async function createCategory(data: any) {
    try {
      await cmsApi.createCategory(data)
      await fetchCategoryTree()
      return true
    } catch (error) {
      console.error('Failed to create category:', error)
      throw error
    }
  }
  
  // ===== 分页 =====
  function goPage(page: number) {
    currentPage.value = page
    return fetchContentList()
  }
  
  return {
    // 状态
    contentList,
    categoryList,
    categoryTree,
    currentPage,
    pageSize,
    total,
    loading,
    
    // 计算属性
    contentCount,
    totalPages,
    
    // 方法
    fetchContentList,
    createContent,
    updateContent,
    deleteContent,
    publishContent,
    fetchCategoryTree,
    fetchCategoryList,
    createCategory,
    goPage
  }
})
```

### 3. 内容列表页面 (src/pages/cms/content/ContentList.vue)

```vue
<template>
  <div class="cms-content">
    <div class="page-header">
      <h1>内容管理</h1>
      <button @click="showDialog = true" class="btn-primary">
        + 新增文章
      </button>
    </div>

    <div class="search-box">
      <input v-model="searchText" placeholder="搜索文章..." @keyup.enter="handleSearch" />
      <button @click="handleSearch">搜索</button>
    </div>

    <!-- 内容列表 -->
    <div class="table-box" v-if="!loading">
      <table class="table">
        <thead>
          <tr>
            <th>标题</th>
            <th>分类</th>
            <th>状态</th>
            <th>浏览数</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in cmsStore.contentList" :key="item.id">
            <td>{{ item.title }}</td>
            <td>{{ getCategoryName(item.category_id) }}</td>
            <td>
              <span :class="['badge', `status-${item.status}`]">
                {{ getStatusLabel(item.status) }}
              </span>
            </td>
            <td>{{ item.view_count }}</td>
            <td>{{ formatDate(item.created_at) }}</td>
            <td class="actions">
              <button @click="handleEdit(item)" class="btn-sm">编辑</button>
              <button 
                @click="handlePublish(item)" 
                v-if="item.status === 1" 
                class="btn-sm publish"
              >
                发布
              </button>
              <button @click="handleDelete(item.id)" class="btn-sm danger">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 加载状态 -->
    <div v-else class="loading">
      <span>加载中...</span>
    </div>

    <!-- 分页 -->
    <div class="pagination">
      <button @click="prevPage" :disabled="cmsStore.currentPage === 1">上一页</button>
      <span>第 {{ cmsStore.currentPage }} 页 / 共 {{ cmsStore.totalPages }} 页</span>
      <button @click="nextPage" :disabled="cmsStore.currentPage >= cmsStore.totalPages">下一页</button>
    </div>

    <!-- 编辑对话框 -->
    <ContentForm 
      v-if="showDialog"
      :isEdit="isEdit"
      :categories="cmsStore.categoryList"
      @save="handleSave"
      @close="showDialog = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useCmsStore } from '@/stores/cms'
import ContentForm from './ContentForm.vue'

const cmsStore = useCmsStore()
const showDialog = ref(false)
const isEdit = ref(false)
const searchText = ref('')
const loading = ref(false)

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    await Promise.all([
      cmsStore.fetchContentList(),
      cmsStore.fetchCategoryList()
    ])
  } finally {
    loading.value = false
  }
}

function getCategoryName(categoryId: number) {
  const category = cmsStore.categoryList.find(c => c.id === categoryId)
  return category?.name || '-'
}

function getStatusLabel(status: number) {
  const labels: { [key: number]: string } = {
    1: '草稿',
    2: '已发布',
    3: '已删除'
  }
  return labels[status] || '未知'
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString()
}

function handleSearch() {
  cmsStore.currentPage = 1
  cmsStore.fetchContentList({ search: searchText.value })
}

function handleEdit(item: any) {
  isEdit.value = true
  showDialog.value = true
  // 传递item到编辑表单
}

async function handleDelete(id: number) {
  if (confirm('确定要删除吗？')) {
    try {
      await cmsStore.deleteContent(id)
    } catch (error) {
      console.error('删除失败:', error)
    }
  }
}

async function handlePublish(item: any) {
  try {
    await cmsStore.publishContent(item.id)
  } catch (error) {
    console.error('发布失败:', error)
  }
}

async function handleSave(data: any) {
  try {
    if (isEdit.value) {
      await cmsStore.updateContent(data.id, data)
    } else {
      await cmsStore.createContent(data)
    }
    showDialog.value = false
    isEdit.value = false
  } catch (error) {
    console.error('保存失败:', error)
  }
}

function prevPage() {
  if (cmsStore.currentPage > 1) {
    cmsStore.currentPage--
    cmsStore.fetchContentList()
  }
}

function nextPage() {
  if (cmsStore.currentPage < cmsStore.totalPages) {
    cmsStore.currentPage++
    cmsStore.fetchContentList()
  }
}
</script>

<style scoped>
.cms-content {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
}

.btn-primary {
  padding: 8px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.search-box {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.search-box input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.search-box button {
  padding: 8px 16px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.table-box {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.table {
  width: 100%;
  border-collapse: collapse;
}

.table th,
.table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #e6e9f0;
}

.table th {
  background: #f5f7fa;
  font-weight: 600;
}

.badge {
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 12px;
  font-weight: 500;
}

.badge.status-1 {
  background: #e6f7ff;
  color: #0050b3;
}

.badge.status-2 {
  background: #d4edda;
  color: #155724;
}

.actions {
  display: flex;
  gap: 8px;
}

.btn-sm {
  padding: 4px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
}

.btn-sm:hover {
  border-color: #667eea;
  color: #667eea;
}

.btn-sm.danger {
  color: #f56c6c;
}

.btn-sm.danger:hover {
  border-color: #f56c6c;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #999;
}

.pagination {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 20px;
}

.pagination button {
  padding: 8px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
```

---

## 数据库查询

### 常用查询

```sql
-- 获取发布的内容列表（按创建时间倒序）
SELECT id, title, excerpt, category_id, view_count, created_at 
FROM cms_content 
WHERE status = 2 
ORDER BY created_at DESC 
LIMIT 10;

-- 获取分类的树形结构
SELECT id, name, parent_id, sort 
FROM cms_category 
WHERE status = 1 
ORDER BY parent_id, sort;

-- 统计各分类的内容数
SELECT c.id, c.name, COUNT(ct.id) as count
FROM cms_category c
LEFT JOIN cms_content ct ON c.id = ct.category_id AND ct.status != 3
GROUP BY c.id, c.name;

-- 获取热门内容
SELECT id, title, view_count, created_at
FROM cms_content
WHERE status = 2
ORDER BY view_count DESC
LIMIT 10;

-- 检查用户权限
SELECT p.* FROM casbin_rule p
WHERE p.v0 = 'cms_admin' AND p.v1 = '/api/cms/admin/contents';
```

---

## 权限配置

### Casbin规则示例

```sql
-- CMS管理员（所有权限）
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'cms_admin', '/api/cms/admin/contents', 'GET'),
('p', 'cms_admin', '/api/cms/admin/contents', 'POST'),
('p', 'cms_admin', '/api/cms/admin/contents', 'PUT'),
('p', 'cms_admin', '/api/cms/admin/contents', 'DELETE'),
('p', 'cms_admin', '/api/cms/admin/categories', 'GET'),
('p', 'cms_admin', '/api/cms/admin/categories', 'POST'),
('p', 'cms_admin', '/api/cms/admin/categories', 'PUT'),
('p', 'cms_admin', '/api/cms/admin/categories', 'DELETE');

-- CMS编辑（可创建和编辑，不能删除）
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'cms_editor', '/api/cms/admin/contents', 'GET'),
('p', 'cms_editor', '/api/cms/admin/contents', 'POST'),
('p', 'cms_editor', '/api/cms/admin/contents', 'PUT'),
('p', 'cms_editor', '/api/cms/admin/categories', 'GET');

-- CMS查看者（仅查看）
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'cms_viewer', '/api/cms/admin/contents', 'GET'),
('p', 'cms_viewer', '/api/cms/admin/categories', 'GET');

-- 将用户分配到角色
INSERT INTO casbin_rule (ptype, v0, v1) VALUES
('g', 'user_id_1', 'cms_admin'),    -- 用户1是CMS管理员
('g', 'user_id_2', 'cms_editor'),   -- 用户2是CMS编辑
('g', 'user_id_3', 'cms_viewer');   -- 用户3是CMS查看者
```

---

**所有代码都可以直接复制使用，或作为开发的参考基础。**

