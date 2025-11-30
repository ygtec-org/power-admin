<template>
  <div class="publish-container">
    <div class="publish-header">
      <h2>发布管理</h2>
      <button @click="openAddModal" class="btn btn-primary">立即发布</button>
    </div>

    <div class="publish-controls">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索内容标题..."
        class="search-input"
      />
      <select v-model="filterStatus" class="filter-select">
        <option value="">所有状态</option>
        <option value="draft">草稿</option>
        <option value="scheduled">定时</option>
        <option value="published">已发布</option>
        <option value="archived">已归档</option>
      </select>
      <button @click="loadPublishList" class="btn btn-secondary">刷新</button>
    </div>

    <div class="publish-list">
      <table v-if="publishList.length > 0" class="table">
        <thead>
          <tr>
            <th>内容ID</th>
            <th>标题</th>
            <th>当前状态</th>
            <th>发布状态</th>
            <th>计划时间</th>
            <th>发布时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in publishList" :key="item.id">
            <td>{{ item.contentId }}</td>
            <td>{{ item.title }}</td>
            <td>
              <span class="status-badge" :class="'status-' + item.contentStatus">
                {{ getStatusText(item.contentStatus) }}
              </span>
            </td>
            <td>
              <span class="status-badge" :class="'status-' + item.publishStatus">
                {{ getPublishStatusText(item.publishStatus) }}
              </span>
            </td>
            <td>{{ formatDate(item.scheduledTime) }}</td>
            <td>{{ formatDate(item.publishedTime) }}</td>
            <td class="actions">
              <button
                v-if="item.publishStatus === 'draft' || item.publishStatus === 'scheduled'"
                @click="handlePublish(item)"
                class="btn btn-small btn-success"
              >
                发布
              </button>
              <button
                v-if="item.publishStatus === 'scheduled'"
                @click="handleCancel(item)"
                class="btn btn-small btn-danger"
              >
                取消定时
              </button>
              <button @click="handleDelete(item)" class="btn btn-small btn-danger">
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-else class="empty-state">
        <p>暂无发布任务</p>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="totalPages > 1" class="pagination">
      <button
        @click="currentPage > 1 && (currentPage -= 1) && loadPublishList()"
        :disabled="currentPage <= 1"
      >
        上一页
      </button>
      <span class="page-info">第 {{ currentPage }} / {{ totalPages }} 页</span>
      <button
        @click="currentPage < totalPages && (currentPage += 1) && loadPublishList()"
        :disabled="currentPage >= totalPages"
      >
        下一页
      </button>
    </div>

    <!-- 发布模态框 -->
    <div v-if="showPublishModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ editingItem ? '编辑发布' : '立即发布' }}</h3>
          <button @click="closeModal" class="close-btn">&times;</button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>选择内容</label>
            <select v-model="formData.contentId" class="form-input" required>
              <option value="">-- 请选择内容 --</option>
              <option v-for="content in availableContents" :key="content.id" :value="content.id">
                {{ content.title }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>发布方式</label>
            <div class="radio-group">
              <label>
                <input v-model="formData.publishType" type="radio" value="immediately" />
                立即发布
              </label>
              <label>
                <input v-model="formData.publishType" type="radio" value="scheduled" />
                定时发布
              </label>
            </div>
          </div>

          <div v-if="formData.publishType === 'scheduled'" class="form-group">
            <label>发布时间</label>
            <input v-model="formData.scheduledTime" type="datetime-local" class="form-input" />
          </div>

          <div class="form-group">
            <label>备注</label>
            <textarea v-model="formData.remark" class="form-input" rows="3"></textarea>
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeModal" class="btn btn-secondary">取消</button>
          <button @click="handleSavePublish" class="btn btn-primary">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const searchQuery = ref('')
const filterStatus = ref('')
const publishList = ref([])
const totalCount = ref(0)
const totalPages = ref(1)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)

const showPublishModal = ref(false)
const editingItem = ref(null)
const availableContents = ref([])

const formData = ref({
  contentId: '',
  publishType: 'immediately',
  scheduledTime: '',
  remark: '',
})

const getStatusText = (status) => {
  const map = {
    draft: '草稿',
    reviewing: '审核中',
    published: '已发布',
    archived: '已归档',
  }
  return map[status] || '未知'
}

const getPublishStatusText = (status) => {
  const map = {
    draft: '未发布',
    scheduled: '定时发布',
    published: '已发布',
  }
  return map[status] || '未知'
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

const loadPublishList = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: currentPage.value,
      pageSize: pageSize.value,
    })
    if (searchQuery.value) params.append('search', searchQuery.value)
    if (filterStatus.value) params.append('status', filterStatus.value)

    const response = await fetch(`/api/cms/publish/list?${params}`)
    const data = await response.json()

    if (data.code === 0 && data.data) {
      publishList.value = data.data.Data || []
      totalCount.value = data.data.Total || 0
      totalPages.value = Math.ceil(totalCount.value / pageSize.value) || 1
    }
  } catch (error) {
    console.error('加载失败:', error)
  } finally {
    loading.value = false
  }
}

const loadAvailableContents = async () => {
  try {
    const response = await fetch('/api/cms/content/list?pageSize=1000')
    const data = await response.json()

    if (data.code === 0 && data.data) {
      availableContents.value = data.data.Data || []
    }
  } catch (error) {
    console.error('加载内容列表失败:', error)
  }
}

const openAddModal = () => {
  editingItem.value = null
  formData.value = {
    contentId: '',
    publishType: 'immediately',
    scheduledTime: '',
    remark: '',
  }
  showPublishModal.value = true
}

const closeModal = () => {
  showPublishModal.value = false
}

const handleSavePublish = async () => {
  if (!formData.value.contentId) {
    alert('请选择内容')
    return
  }

  try {
    const payload = {
      contentId: parseInt(formData.value.contentId),
      publishType: formData.value.publishType,
      scheduledTime: formData.value.scheduledTime || null,
      remark: formData.value.remark,
    }

    const response = await fetch('/api/cms/publish/create', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    const data = await response.json()
    if (data.code === 0) {
      alert('发布成功')
      closeModal()
      loadPublishList()
    } else {
      alert('发布失败: ' + data.msg)
    }
  } catch (error) {
    console.error('发布失败:', error)
    alert('发布失败')
  }
}

const handlePublish = async (item) => {
  if (!confirm('确定要发布此内容吗？')) return

  try {
    const response = await fetch(`/api/cms/publish/publish/${item.id}`, {
      method: 'POST',
    })

    const data = await response.json()
    if (data.code === 0) {
      alert('发布成功')
      loadPublishList()
    } else {
      alert('发布失败: ' + data.msg)
    }
  } catch (error) {
    console.error('发布失败:', error)
    alert('发布失败')
  }
}

const handleCancel = async (item) => {
  if (!confirm('确定要取消定时发布吗？')) return

  try {
    const response = await fetch(`/api/cms/publish/cancel/${item.id}`, {
      method: 'POST',
    })

    const data = await response.json()
    if (data.code === 0) {
      alert('取消成功')
      loadPublishList()
    } else {
      alert('取消失败: ' + data.msg)
    }
  } catch (error) {
    console.error('取消失败:', error)
    alert('取消失败')
  }
}

const handleDelete = async (item) => {
  if (!confirm('确定要删除此发布任务吗？')) return

  try {
    const response = await fetch(`/api/cms/publish/delete/${item.id}`, {
      method: 'DELETE',
    })

    const data = await response.json()
    if (data.code === 0) {
      alert('删除成功')
      loadPublishList()
    } else {
      alert('删除失败: ' + data.msg)
    }
  } catch (error) {
    console.error('删除失败:', error)
    alert('删除失败')
  }
}

onMounted(() => {
  loadPublishList()
  loadAvailableContents()
})
</script>

<style scoped>
.publish-container {
  padding: 20px;
  background: white;
  border-radius: 8px;
}

.publish-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.publish-header h2 {
  margin: 0;
  font-size: 20px;
  color: #333;
}

.publish-controls {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
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
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover {
  background: #5568d3;
}

.btn-secondary {
  background: #f0f0f0;
  color: #333;
}

.btn-secondary:hover {
  background: #e0e0e0;
}

.btn-success {
  background: #48bb78;
  color: white;
}

.btn-danger {
  background: #f56565;
  color: white;
}

.btn-small {
  padding: 4px 8px;
  font-size: 12px;
  margin-right: 4px;
}

.table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 20px;
}

.table thead {
  background: #f5f7fa;
}

.table th,
.table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.table th {
  font-weight: 600;
  color: #333;
}

.status-badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-draft {
  background: #fef3c7;
  color: #92400e;
}

.status-reviewing {
  background: #bfdbfe;
  color: #1e3a8a;
}

.status-published {
  background: #bbf7d0;
  color: #065f46;
}

.status-archived {
  background: #e5e7eb;
  color: #374151;
}

.status-scheduled {
  background: #c7d2fe;
  color: #312e81;
}

.actions {
  white-space: nowrap;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #999;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
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

.page-info {
  color: #666;
  font-size: 14px;
}

.modal-overlay {
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

.modal-content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  width: 90%;
  max-width: 500px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e6e9f0;
}

.modal-header h3 {
  margin: 0;
  font-size: 16px;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
}

.modal-body {
  padding: 16px;
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

.form-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  font-family: inherit;
}

.radio-group {
  display: flex;
  gap: 16px;
}

.radio-group label {
  display: flex;
  align-items: center;
  margin-bottom: 0;
}

.radio-group input {
  margin-right: 6px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px;
  border-top: 1px solid #e6e9f0;
}
</style>
