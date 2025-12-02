<template>
  <div class="cms-tag-list">
    <div class="header">
      <h1>标签管理</h1>
      <button class="btn-primary" @click="openCreateDialog">新建标签</button>
    </div>

    <!-- 标签卡片网格 -->
    <div class="tag-grid">
      <div v-if="tags.length === 0" class="empty-state">
        暂无标签
      </div>
      <div v-for="tag in tags" :key="tag.id" class="tag-card">
        <div class="tag-header" :style="{ borderLeftColor: tag.color || '#007bff' }">
          <h3 class="tag-name">{{ tag.name }}</h3>
          <div class="tag-actions">
            <button class="btn-edit" @click="editTag(tag)">编辑</button>
            <button class="btn-delete" @click="deleteTag(tag.id)">删除</button>
          </div>
        </div>
        <div class="tag-content">
          <p class="tag-desc">{{ tag.description }}</p>
          <div class="tag-meta">
            <span class="tag-slug">Slug: {{ tag.slug }}</span>
            <span class="tag-count">使用: {{ tag.usageCount }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑对话框 -->
    <div v-if="showDialog" class="dialog-overlay">
      <div class="dialog">
        <div class="dialog-header">
          <h2>{{ editingId ? '编辑标签' : '新建标签' }}</h2>
          <button class="btn-close" @click="closeDialog">×</button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label>标签名称</label>
            <input v-model="formData.name" type="text" class="form-input" />
          </div>
          <div class="form-group">
            <label>标签描述</label>
            <textarea v-model="formData.description" class="form-textarea"></textarea>
          </div>
          <div class="form-group">
            <label>Slug</label>
            <input v-model="formData.slug" type="text" class="form-input" />
          </div>
          <div class="form-group">
            <label>颜色</label>
            <input v-model="formData.color" type="color" class="form-color" />
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn-cancel" @click="closeDialog">取消</button>
          <button class="btn-primary" @click="saveTag">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTagList, createTag, updateTag, deleteTag as apiDeleteTag } from '@/api/cms'

const tags = ref([])
const showDialog = ref(false)
const editingId = ref(null)

const formData = ref({
  name: '',
  description: '',
  slug: '',
  color: '#007bff',
})

onMounted(() => {
  loadTags()
})

const loadTags = async () => {
  try {
    const res = await getTagList()
    tags.value = res.data || []
  } catch (error) {
    console.error('加载标签失败:', error)
  }
}

const openCreateDialog = () => {
  editingId.value = null
  formData.value = {
    name: '',
    description: '',
    slug: '',
    color: '#007bff',
  }
  showDialog.value = true
}

const editTag = (tag) => {
  editingId.value = tag.id
  formData.value = {
    name: tag.name,
    description: tag.description,
    slug: tag.slug,
    color: tag.color || '#007bff',
  }
  showDialog.value = true
}

const saveTag = async () => {
  try {
    if (editingId.value) {
      await updateTag(editingId.value, formData.value)
      ElMessage({
        message: '更新成功',
        type: 'success',
        duration: 2000,
        offset: 20,
      })
    } else {
      await createTag(formData.value)
      ElMessage({
        message: '创建成功',
        type: 'success',
        duration: 2000,
        offset: 20,
      })
    }
    closeDialog()
    loadTags()
  } catch (error) {
    console.error('保存失败:', error)
  }
}

const closeDialog = () => {
  showDialog.value = false
}

const deleteTag = async (id) => {
  ElMessageBox.confirm('确定删除此标签吗？', '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      await apiDeleteTag(id)
      ElMessage({
        message: '删除成功',
        type: 'success',
        duration: 2000,
        offset: 20,
      })
      loadTags()
    } catch (error) {
      console.error('删除失败:', error)
    }
  }).catch(() => {
    // 用户取消删除
  })
}
</script>

<style scoped>
.cms-tag-list {
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

.tag-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  color: #999;
  padding: 40px;
  background: white;
  border-radius: 4px;
}

.tag-card {
  background: white;
  border-radius: 4px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.tag-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.tag-header {
  padding: 15px;
  border-left: 4px solid #007bff;
  background: #f8f9fa;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.tag-name {
  margin: 0;
  font-size: 16px;
  color: #333;
  flex: 1;
}

.tag-actions {
  display: flex;
  gap: 5px;
  margin-left: 10px;
}

.btn-edit,
.btn-delete {
  padding: 4px 8px;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  background: #e9ecef;
  color: #333;
  white-space: nowrap;
}

.btn-edit:hover {
  background: #17a2b8;
  color: white;
}

.btn-delete:hover {
  background: #dc3545;
  color: white;
}

.tag-content {
  padding: 15px;
}

.tag-desc {
  margin: 0 0 10px 0;
  color: #666;
  font-size: 13px;
  line-height: 1.5;
}

.tag-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
}

.tag-slug,
.tag-count {
  white-space: nowrap;
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
  max-width: 500px;
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
  box-sizing: border-box;
}

.form-color {
  width: 50px;
  height: 40px;
  padding: 2px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
}

.form-textarea {
  min-height: 80px;
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
