<template>
  <div class="cms-content-list">
    <div class="header">
      <h1>内容管理</h1>
      <button class="btn-primary" @click="openCreateDialog">新建内容</button>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <input
        v-model="searchForm.search"
        type="text"
        placeholder="搜索内容..."
        class="search-input"
      />
      <select v-model="searchForm.categoryId" class="filter-select">
        <option value="">全部分类</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </option>
      </select>
      <select v-model="searchForm.status" class="filter-select">
        <option value="">全部状态</option>
        <option value="1">草稿</option>
        <option value="2">已发布</option>
        <option value="3">已删除</option>
      </select>
      <button class="btn-search" @click="searchContent">搜索</button>
    </div>

    <!-- 内容表格 -->
    <div class="table-container">
      <table class="content-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>标题</th>
            <th>分类</th>
            <th>状态</th>
            <th>浏览量</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in contentList" :key="item.id">
            <td>{{ item.id }}</td>
            <td class="title">{{ item.title }}</td>
            <td>{{ item.categoryName }}</td>
            <td>
              <span :class="`status status-${item.status}`">
                {{ getStatusText(item.status) }}
              </span>
            </td>
            <td>{{ item.viewCount }}</td>
            <td>{{ formatDate(item.createdAt) }}</td>
            <td class="actions">
              <button class="btn-edit" @click="editContent(item)">编辑</button>
              <button
                v-if="item.status !== 2"
                class="btn-publish"
                @click="publishContent(item.id)"
              >
                发布
              </button>
              <button
                v-else
                class="btn-unpublish"
                @click="unpublishContent(item.id)"
              >
                取消发布
              </button>
              <button class="btn-delete" @click="deleteContent(item.id)">
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 分页 -->
    <div class="pagination">
      <button
        :disabled="pagination.page === 1"
        @click="pagination.page--"
        class="btn-page"
      >
        上一页
      </button>
      <span class="page-info">第 {{ pagination.page }} 页</span>
      <button
        :disabled="contentList.length < pagination.pageSize"
        @click="pagination.page++"
        class="btn-page"
      >
        下一页
      </button>
    </div>

    <!-- 编辑对话框 -->
    <div v-if="showDialog" class="dialog-overlay">
      <div class="dialog">
        <div class="dialog-header">
          <h2>{{ editingId ? '编辑内容' : '新建内容' }}</h2>
          <button class="btn-close" @click="closeDialog">×</button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label>标题</label>
            <input v-model="formData.title" type="text" class="form-input" />
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea
              v-model="formData.description"
              class="form-textarea"
            ></textarea>
          </div>
          <div class="form-group">
            <label>分类</label>
            <select v-model="formData.categoryId" class="form-select">
              <option value="">选择分类</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                {{ cat.name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>Slug</label>
            <input v-model="formData.slug" type="text" class="form-input" />
          </div>
          <div class="form-group">
            <label>内容</label>
            <textarea v-model="formData.content" class="form-textarea"></textarea>
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn-cancel" @click="closeDialog">取消</button>
          <button class="btn-primary" @click="saveContent">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getContentList, createContent, updateContent, deleteContent, publishContent as apiPublish, unpublishContent as apiUnpublish, getCategoryList } from '@/api/cms'

// 数据
const contentList = ref<any[]>([])
const categories = ref<any[]>([])
const showDialog = ref(false)
const editingId = ref<number | null>(null)

const searchForm = ref({
  search: '',
  categoryId: '',
  status: '',
})

const pagination = ref({
  page: 1,
  pageSize: 10,
})

const formData = ref({
  title: '',
  description: '',
  content: '',
  slug: '',
  categoryId: 0,
  status: 1,
})

// 生命周期
onMounted(() => {
  loadCategories()
  loadContent()
})

// 加载分类
const loadCategories = async () => {
  try {
    const res = await getCategoryList()
    categories.value = res.data || []
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

// 加载内容
const loadContent = async () => {
  try {
    const res = await getContentList({
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
      categoryId: searchForm.value.categoryId ? parseInt(searchForm.value.categoryId) : undefined,
      status: searchForm.value.status ? parseInt(searchForm.value.status) : undefined,
      search: searchForm.value.search,
    })
    contentList.value = res.data?.data || []
  } catch (error) {
    console.error('加载内容失败:', error)
  }
}

// 搜索
const searchContent = () => {
  pagination.value.page = 1
  loadContent()
}

// 打开创建对话框
const openCreateDialog = () => {
  editingId.value = null
  formData.value = {
    title: '',
    description: '',
    content: '',
    slug: '',
    categoryId: 0,
    status: 1,
  }
  showDialog.value = true
}

// 编辑内容
const editContent = (item: any) => {
  editingId.value = item.id
  formData.value = {
    title: item.title,
    description: item.description,
    content: item.content,
    slug: item.slug,
    categoryId: item.categoryId,
    status: item.status,
  }
  showDialog.value = true
}

// 保存内容
const saveContent = async () => {
  try {
    if (editingId.value) {
      await updateContent(editingId.value, formData.value)
      alert('更新成功')
    } else {
      await createContent(formData.value)
      alert('创建成功')
    }
    closeDialog()
    loadContent()
  } catch (error) {
    console.error('保存失败:', error)
    alert('保存失败')
  }
}

// 关闭对话框
const closeDialog = () => {
  showDialog.value = false
}

// 删除内容
const deleteContent = async (id: number) => {
  if (confirm('确定要删除吗?')) {
    try {
      await deleteContent(id)
      alert('删除成功')
      loadContent()
    } catch (error) {
      console.error('删除失败:', error)
      alert('删除失败')
    }
  }
}

// 发布内容
const publishContent = async (id: number) => {
  try {
    await apiPublish(id)
    alert('发布成功')
    loadContent()
  } catch (error) {
    console.error('发布失败:', error)
    alert('发布失败')
  }
}

// 取消发布
const unpublishContent = async (id: number) => {
  try {
    await apiUnpublish(id)
    alert('取消发布成功')
    loadContent()
  } catch (error) {
    console.error('取消发布失败:', error)
    alert('取消发布失败')
  }
}

// 辅助函数
const getStatusText = (status: number) => {
  const map: Record<number, string> = {
    1: '草稿',
    2: '已发布',
    3: '已删除',
  }
  return map[status] || '未知'
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.cms-content-list {
  padding: 20px;
  background: #f5f5f5;
  min-height: 100vh;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h1 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.btn-primary {
  padding: 10px 20px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-primary:hover {
  background: #0056b3;
}

.search-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  background: white;
  padding: 15px;
  border-radius: 4px;
}

.search-input,
.filter-select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.search-input {
  flex: 1;
  min-width: 200px;
}

.filter-select {
  min-width: 120px;
}

.btn-search {
  padding: 8px 16px;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-search:hover {
  background: #218838;
}

.table-container {
  background: white;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 20px;
}

.content-table {
  width: 100%;
  border-collapse: collapse;
}

.content-table thead {
  background: #f8f9fa;
  border-bottom: 2px solid #ddd;
}

.content-table th {
  padding: 12px;
  text-align: left;
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.content-table td {
  padding: 12px;
  border-bottom: 1px solid #eee;
  font-size: 14px;
}

.content-table tbody tr:hover {
  background: #f9f9f9;
}

.title {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 12px;
  font-weight: 600;
}

.status-1 {
  background: #ffc107;
  color: #333;
}

.status-2 {
  background: #28a745;
  color: white;
}

.status-3 {
  background: #dc3545;
  color: white;
}

.actions {
  display: flex;
  gap: 5px;
}

.btn-edit,
.btn-publish,
.btn-unpublish,
.btn-delete {
  padding: 4px 10px;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  background: #e9ecef;
  color: #333;
}

.btn-edit:hover {
  background: #17a2b8;
  color: white;
}

.btn-publish:hover {
  background: #28a745;
  color: white;
}

.btn-unpublish:hover {
  background: #ffc107;
  color: #333;
}

.btn-delete:hover {
  background: #dc3545;
  color: white;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  margin-top: 20px;
}

.btn-page {
  padding: 8px 16px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-page:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: #666;
}

/* 对话框样式 */
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.dialog {
  background: white;
  border-radius: 4px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.dialog-header h2 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
}

.btn-close:hover {
  color: #333;
}

.dialog-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  font-family: inherit;
}

.form-textarea {
  min-height: 100px;
  resize: vertical;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 20px;
  border-top: 1px solid #eee;
}

.btn-cancel {
  padding: 8px 16px;
  background: #e9ecef;
  color: #333;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-cancel:hover {
  background: #ddd;
}
</style>
