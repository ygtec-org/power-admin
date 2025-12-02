<template>
  <div class="cms-content-container">
    <div class="page-header">
      <h2>内容管理</h2>
      <button @click="openCreateModal" class="btn-primary">+ 新增内容</button>
    </div>

    <div class="filters">
      <input v-model="searchQuery" type="text" placeholder="搜索文章标题..." @keyup.enter="loadContent" class="search-input">
      <select v-model="filterStatus" @change="loadContent" class="filter-select">
        <option value="">全部状态</option>
        <option value="1">草稿</option>
        <option value="2">已发布</option>
      </select>
      <button @click="loadContent" class="btn-search">搜索</button>
    </div>

    <div class="table-container">
      <table class="content-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>标题</th>
            <th>状态</th>
            <th>浏览次数</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in contentList" :key="item.id">
            <td>{{ item.id }}</td>
            <td>{{ item.title }}</td>
            <td><span :class="['badge', getStatusClass(item.status)]">{{ getStatusText(item.status) }}</span></td>
            <td>{{ item.viewCount || 0 }}</td>
            <td>{{ formatDate(item.createdAt) }}</td>
            <td class="action-cell">
              <button @click="editContent(item)" class="btn-small">编辑</button>
              <button @click="deleteItem(item.id)" class="btn-small delete">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="!loading && contentList.length === 0" class="empty-state">暂无数据</div>
    </div>

    <div class="pagination">
      <button @click="prevPage" :disabled="currentPage === 1" class="btn-pagination">上一页</button>
      <span class="page-info">第 {{ currentPage }} 页，共 {{ totalPages }} 页</span>
      <button @click="nextPage" :disabled="currentPage >= totalPages" class="btn-pagination">下一页</button>
    </div>

    <!-- 模态框 -->
    <div v-if="showModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ editingId ? '编辑内容' : '新增内容' }}</h3>
          <button @click="closeModal" class="btn-close">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>标题</label>
            <input v-model="formData.title" type="text" class="form-input">
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="formData.description" class="form-textarea" rows="3"></textarea>
          </div>
          <div class="form-group">
            <label>内容</label>
            <textarea v-model="formData.content" class="form-textarea" rows="8"></textarea>
          </div>
          <div class="form-group">
            <label>状态</label>
            <select v-model="formData.status" class="form-input">
              <option :value="1">草稿</option>
              <option :value="2">已发布</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn-cancel">取消</button>
          <button @click="saveContent" class="btn-primary">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getContentList, createContent, updateContent, deleteContent } from '@/api/cms'

const contentList = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const totalCount = ref(0)
const totalPages = ref(0)
const searchQuery = ref('')
const filterStatus = ref('')
const showModal = ref(false)
const editingId = ref(null)
const loading = ref(false)

const formData = ref({
  title: '',
  description: '',
  content: '',
  slug: '',
  categoryId: 0,
  status: 1,
})

const loadContent = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
    }
    if (searchQuery.value) params.search = searchQuery.value
    if (filterStatus.value) params.status = parseInt(filterStatus.value)

    const res = await getContentList(params)
    console.log('API 返回数据:', res)
    
    if (res.data) {
      contentList.value = res.data.data || []
      totalCount.value = res.data.total || 0
      totalPages.value = Math.ceil(totalCount.value / pageSize.value) || 1
    }
  } catch (error) {
    console.error('加载失败:', error)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  editingId.value = null
  formData.value = { 
    title: '', 
    description: '', 
    content: '', 
    slug: '',
    categoryId: 0,
    status: 1 
  }
  showModal.value = true
}

const editContent = (item) => {
  editingId.value = item.id
  formData.value = { 
    title: item.title,
    description: item.description,
    content: item.content,
    slug: item.slug,
    categoryId: item.categoryId || 0,
    status: item.status
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveContent = async () => {
  try {
    if (editingId.value) {
      await updateContent(editingId.value, formData.value)
      ElMessage({
        message: '更新成功',
        type: 'success',
        duration: 2000,
        offset: 20,
      })
    } else {
      await createContent(formData.value)
      ElMessage({
        message: '创建成功',
        type: 'success',
        duration: 2000,
        offset: 20,
      })
    }
    closeModal()
    currentPage.value = 1
    loadContent()
  } catch (error) {
    console.error('保存失败:', error)
  }
}

const deleteItem = async (id) => {
  ElMessageBox.confirm('确定删除此内容吗？', '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      await deleteContent(id)
      ElMessage({
        message: '删除成功',
        type: 'success',
        duration: 2000,
        offset: 20,
      })
      loadContent()
    } catch (error) {
      console.error('删除失败:', error)
    }
  }).catch(() => {
    // 用户取消删除
  })
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadContent()
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadContent()
  }
}

const getStatusText = (status) => {
  return status === 1 ? '草稿' : '已发布'
}

const getStatusClass = (status) => {
  return status === 1 ? 'badge-draft' : 'badge-published'
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return dateStr.substring(0, 19).replace('T', ' ')
}

onMounted(() => {
  loadContent()
})
</script>

<style scoped>
.cms-content-container {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.btn-primary {
  padding: 8px 16px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-primary:hover {
  background: #5568d3;
}

.filters {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.search-input, .filter-select, .btn-search {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.search-input {
  flex: 1;
}

.btn-search {
  background: #667eea;
  color: white;
  border: none;
  cursor: pointer;
}

.table-container {
  background: white;
  border-radius: 4px;
  overflow: auto;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.content-table {
  width: 100%;
  border-collapse: collapse;
}

.content-table th {
  padding: 12px;
  text-align: left;
  background: #f5f7fa;
  border-bottom: 2px solid #e6e9f0;
  font-weight: 600;
  color: #333;
}

.content-table td {
  padding: 12px;
  border-bottom: 1px solid #e6e9f0;
}

.badge {
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 12px;
}

.badge.draft {
  background: #f0f0f0;
  color: #666;
}

.badge.published {
  background: #e6f7ff;
  color: #1890ff;
}

.action-cell {
  display: flex;
  gap: 8px;
}

.btn-small {
  padding: 4px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  color: #667eea;
}

.btn-small.delete {
  color: #ff4d4f;
  border-color: #ff4d4f;
}

.empty-state {
  padding: 40px;
  text-align: center;
  color: #999;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 12px;
}

.btn-pagination {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-pagination:disabled {
  cursor: not-allowed;
  color: #ccc;
}

.page-info {
  font-size: 14px;
  color: #666;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  max-width: 600px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e6e9f0;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  color: #999;
  cursor: pointer;
  padding: 0;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
  font-size: 14px;
}

.form-input, .form-textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-textarea {
  resize: vertical;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid #e6e9f0;
}

.btn-cancel {
  padding: 8px 16px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}
</style>
