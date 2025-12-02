<template>
  <div class="cms-category-list">
    <div class="header">
      <h1>分类管理</h1>
      <button class="btn-primary" @click="openCreateDialog">新建分类</button>
    </div>

    <!-- 分类树 -->
    <div class="tree-container">
      <div v-if="categories.length === 0" class="empty-state">
        暂无分类
      </div>
      <div v-for="category in categories" :key="category.id" class="tree-node">
        <div class="node-content">
          <span class="node-name">{{ category.name }}</span>
          <span class="node-desc">{{ category.description }}</span>
          <div class="node-actions">
            <button class="btn-edit" @click="editCategory(category)">编辑</button>
            <button class="btn-delete" @click="deleteCategory(category.id)">删除</button>
          </div>
        </div>
        <!-- 子分类 -->
        <div v-if="category.children && category.children.length > 0" class="children">
          <div v-for="child in category.children" :key="child.id" class="tree-node child">
            <div class="node-content">
              <span class="node-name">{{ child.name }}</span>
              <span class="node-desc">{{ child.description }}</span>
              <div class="node-actions">
                <button class="btn-edit" @click="editCategory(child)">编辑</button>
                <button class="btn-delete" @click="deleteCategory(child.id)">删除</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑对话框 -->
    <div v-if="showDialog" class="dialog-overlay">
      <div class="dialog">
        <div class="dialog-header">
          <h2>{{ editingId ? '编辑分类' : '新建分类' }}</h2>
          <button class="btn-close" @click="closeDialog">×</button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label>分类名称</label>
            <input v-model="formData.name" type="text" class="form-input" />
          </div>
          <div class="form-group">
            <label>分类描述</label>
            <textarea v-model="formData.description" class="form-textarea"></textarea>
          </div>
          <div class="form-group">
            <label>Slug</label>
            <input v-model="formData.slug" type="text" class="form-input" />
          </div>
          <div class="form-group">
            <label>父分类</label>
            <select v-model="formData.parentId" class="form-select">
              <option value="">无父分类</option>
              <option
                v-for="cat in categories"
                :key="cat.id"
                :value="cat.id"
                :disabled="cat.id === editingId"
              >
                {{ cat.name }}
              </option>
            </select>
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn-cancel" @click="closeDialog">取消</button>
          <button class="btn-primary" @click="saveCategory">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getCategoryTree,
  createCategory,
  updateCategory,
  deleteCategory as apiDeleteCategory,
} from '@/api/cms'

const categories = ref([])
const showDialog = ref(false)
const editingId = ref(null)

const formData = ref({
  name: '',
  description: '',
  slug: '',
  parentId: '',
})

onMounted(() => {
  loadCategories()
})

const loadCategories = async () => {
  try {
    const res = await getCategoryTree()
    categories.value = res.data || []
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

const openCreateDialog = () => {
  editingId.value = null
  formData.value = {
    name: '',
    description: '',
    slug: '',
    parentId: '',
  }
  showDialog.value = true
}

const editCategory = (category) => {
  editingId.value = category.id
  formData.value = {
    name: category.name,
    description: category.description,
    slug: category.slug,
    parentId: category.parentId || '',
  }
  showDialog.value = true
}

const saveCategory = async () => {
  try {
    const data = {
      name: formData.value.name,
      description: formData.value.description,
      slug: formData.value.slug,
      parentId: formData.value.parentId ? parseInt(formData.value.parentId) : 0,
    }

    if (editingId.value) {
      await updateCategory(editingId.value, data)
      ElMessage({
        message: '更新成功',
        type: 'success',
        duration: 2000,
        offset: 20,
      })
    } else {
      await createCategory(data)
      ElMessage({
        message: '创建成功',
        type: 'success',
        duration: 2000,
        offset: 20,
      })
    }
    closeDialog()
    loadCategories()
  } catch (error) {
    console.error('保存失败:', error)
  }
}

const closeDialog = () => {
  showDialog.value = false
}

const deleteCategory = async (id) => {
  ElMessageBox.confirm('确定删除此分类吗？', '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      await apiDeleteCategory(id)
      ElMessage({
        message: '删除成功',
        type: 'success',
        duration: 2000,
        offset: 20,
      })
      loadCategories()
    } catch (error) {
      console.error('删除失败:', error)
    }
  }).catch(() => {
    // 用户取消删除
  })
}
</script>

<style scoped>
.cms-category-list {
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

.tree-container {
  background: white;
  border-radius: 4px;
  padding: 20px;
}

.empty-state {
  text-align: center;
  color: #999;
  padding: 40px;
}

.tree-node {
  margin-bottom: 15px;
  border-left: 2px solid #007bff;
  padding-left: 15px;
}

.tree-node.child {
  margin-left: 20px;
  border-left-color: #17a2b8;
}

.node-content {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 4px;
}

.node-name {
  font-weight: 600;
  color: #333;
  min-width: 100px;
}

.node-desc {
  flex: 1;
  color: #666;
  font-size: 13px;
}

.node-actions {
  display: flex;
  gap: 5px;
}

.btn-edit,
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

.btn-delete:hover {
  background: #dc3545;
  color: white;
}

.children {
  margin-top: 10px;
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
