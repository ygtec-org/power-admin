<template>
  <div class="cms-comment-list">
    <div class="header">
      <h1>评论管理</h1>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <input
        v-model="filters.contentId"
        type="number"
        placeholder="按内容ID搜索..."
        class="search-input"
      />
      <button class="btn-search" @click="loadComments">搜索</button>
    </div>

    <!-- 评论表格 -->
    <div class="table-container">
      <table class="comment-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>内容</th>
            <th>作者</th>
            <th>评论内容</th>
            <th>状态</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="comment in comments" :key="comment.id">
            <td>{{ comment.id }}</td>
            <td>内容 #{{ comment.contentId }}</td>
            <td>{{ comment.authorName }}</td>
            <td class="content-cell">{{ comment.content }}</td>
            <td>
              <span :class="`status status-${comment.status}`">
                {{ getStatusText(comment.status) }}
              </span>
            </td>
            <td>{{ formatDate(comment.createdAt) }}</td>
            <td class="actions">
              <button
                v-if="comment.status === 0"
                class="btn-approve"
                @click="approveComment(comment.id)"
              >
                通过
              </button>
              <button
                v-if="comment.status === 0"
                class="btn-reject"
                @click="rejectComment(comment.id)"
              >
                拒绝
              </button>
              <button class="btn-delete" @click="deleteComment(comment.id)">
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
        :disabled="comments.length < pagination.pageSize"
        @click="pagination.page++"
        class="btn-page"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getCommentList,
  approveComment as apiApprove,
  rejectComment as apiReject,
  deleteComment as apiDelete,
} from '@/api/cms'

const comments = ref([])
const filters = ref({
  contentId: '',
})
const pagination = ref({
  page: 1,
  pageSize: 10,
})

onMounted(() => {
  loadComments()
})

const loadComments = async () => {
  try {
    const res = await getCommentList({
      contentId: filters.value.contentId ? parseInt(filters.value.contentId) : undefined,
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
    })
    comments.value = res.data?.data || []
  } catch (error) {
    console.error('加载评论失败:', error)
  }
}

const approveComment = async (id) => {
  try {
    await apiApprove(id)
    ElMessage({
      message: '通过评论成功',
      type: 'success',
      duration: 2000,
      offset: 20,
    })
    loadComments()
  } catch (error) {
    console.error('操作失败:', error)
  }
}

const rejectComment = async (id) => {
  try {
    await apiReject(id)
    ElMessage({
      message: '拒绝评论成功',
      type: 'success',
      duration: 2000,
      offset: 20,
    })
    loadComments()
  } catch (error) {
    console.error('操作失败:', error)
  }
}

const deleteComment = async (id) => {
  ElMessageBox.confirm('确定删除此评论吗？', '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      await apiDelete(id)
      ElMessage({
        message: '删除成功',
        type: 'success',
        duration: 2000,
        offset: 20,
      })
      loadComments()
    } catch (error) {
      console.error('删除失败:', error)
    }
  }).catch(() => {
    // 用户取消删除
  })
}

const getStatusText = (status) => {
  const map = {
    0: '待审核',
    1: '已通过',
    2: '已拒绝',
  }
  return map[status] || '未知'
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.cms-comment-list {
  padding: 20px;
  background: #f5f5f5;
  min-height: 100vh;
}

.header {
  margin-bottom: 20px;
}

.header h1 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.search-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  background: white;
  padding: 15px;
  border-radius: 4px;
}

.search-input {
  flex: 1;
  min-width: 200px;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
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

.comment-table {
  width: 100%;
  border-collapse: collapse;
}

.comment-table thead {
  background: #f8f9fa;
  border-bottom: 2px solid #ddd;
}

.comment-table th {
  padding: 12px;
  text-align: left;
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.comment-table td {
  padding: 12px;
  border-bottom: 1px solid #eee;
  font-size: 14px;
}

.comment-table tbody tr:hover {
  background: #f9f9f9;
}

.content-cell {
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

.status-0 {
  background: #ffc107;
  color: #333;
}

.status-1 {
  background: #28a745;
  color: white;
}

.status-2 {
  background: #dc3545;
  color: white;
}

.actions {
  display: flex;
  gap: 5px;
}

.btn-approve,
.btn-reject,
.btn-delete {
  padding: 4px 10px;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
}

.btn-approve {
  background: #28a745;
  color: white;
}

.btn-approve:hover {
  background: #218838;
}

.btn-reject {
  background: #ffc107;
  color: #333;
}

.btn-reject:hover {
  background: #e0a800;
}

.btn-delete {
  background: #dc3545;
  color: white;
}

.btn-delete:hover {
  background: #c82333;
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
</style>
