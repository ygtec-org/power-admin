<template>
  <div class="page">
    <div class="page-header">
      <h1>权限管理</h1>
      <button @click="showAddDialog = true" class="btn-primary">+ 新增权限</button>
    </div>

    <div class="table-box">
      <table class="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>权限名称</th>
            <th>资源</th>
            <th>操作</th>
            <th>描述</th>
            <th>状态</th>
            <th width="200">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="permissions.length === 0">
            <td colspan="7" class="empty">暂无数据</td>
          </tr>
          <tr v-for="perm in permissions" :key="perm.id">
            <td>{{ perm.id }}</td>
            <td>{{ perm.name }}</td>
            <td><code>{{ perm.resource }}</code></td>
            <td><code>{{ perm.action }}</code></td>
            <td>{{ perm.description || '-' }}</td>
            <td>
              <span class="badge" :class="perm.status === 1 ? 'success' : 'danger'">
                {{ perm.status === 1 ? '启用' : '禁用' }}
              </span>
            </td>
            <td>
              <button @click="handleEdit(perm)" class="btn-sm">编辑</button>
              <button @click="handleDelete(perm)" class="btn-sm danger">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 使用通用分页组件 -->
    <Pagination 
      v-model:currentPage="page"
      v-model:pageSize="pageSize"
      :total="total"
      @page-change="handlePageChange"
      @size-change="handleSizeChange"
    />

    <!-- 新增/编辑对话框 -->
    <div v-if="showAddDialog" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ isEdit ? '编辑权限' : '新增权限' }}</h2>
          <button @click="closeDialog" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>权限名称 *</label>
            <input v-model="form.name" type="text" placeholder="请输入权限名称" />
          </div>
          <div class="form-row">
            <div class="form-group" style="flex: 1; margin-right: 12px">
              <label>资源 *</label>
              <input v-model="form.resource" type="text" placeholder="如: user, role" />
            </div>
            <div class="form-group" style="flex: 1">
              <label>操作 *</label>
              <input v-model="form.action" type="text" placeholder="如: list, create, edit, delete" />
            </div>
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="form.description" placeholder="请输入权限描述"></textarea>
          </div>
          <div class="form-group">
            <label>状态</label>
            <select v-model="form.status">
              <option value="1">启用</option>
              <option value="0">禁用</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeDialog" class="btn-cancel">取消</button>
          <button @click="handleSave" class="btn-primary">保存</button>
        </div>
      </div>
    </div>

    <ConfirmDialog 
      v-model:visible="showConfirm"
      title="删除确认"
      :message="`确定要删除权限「${deleteTarget?.name}」吗？`"
      @confirm="handleConfirmDelete"
      @cancel="handleCancelDelete"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPermissions, createPermission, updatePermission, deletePermission } from '../../../api/permission'
import { notify } from '../../../components/Notification'
import ConfirmDialog from '../../../components/ConfirmDialog.vue'
import Pagination from '../../../components/Pagination.vue'

const permissions = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const showAddDialog = ref(false)
const isEdit = ref(false)
const form = ref({
  name: '',
  resource: '',
  action: '',
  description: '',
  status: 1,
})
const selectedPerm = ref(null)

// 获取权限列表
const loadPermissions = async (pageNum = page.value, pageSizeNum = pageSize.value) => {
  try {
    const res = await getPermissions({ page: pageNum, pageSize: pageSizeNum })
    permissions.value = res.data.list || []
    total.value = res.data.total || 0
    page.value = pageNum
    pageSize.value = pageSizeNum
  } catch (error) {
    notify.error(error.message || '获取权限列表失败')
  }
}

const handlePageChange = (pageNum) => {
  loadPermissions(pageNum, pageSize.value)
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadPermissions(1, size) // 页码重置为第一页
}

// 编辑权限
const handleEdit = (perm) => {
  isEdit.value = true
  selectedPerm.value = perm
  form.value = {
    name: perm.name,
    resource: perm.resource,
    action: perm.action,
    description: perm.description,
    status: perm.status,
  }
  showAddDialog.value = true
}

// 删除权限
const handleDelete = async (perm) => {
  deleteTarget.value = perm
  showConfirm.value = true
}

const handleConfirmDelete = async () => {
  try {
    await deletePermission(deleteTarget.value.id)
    notify.success('删除成功')
    loadPermissions()
  } catch (error) {
    notify.error(error.message || '删除失败')
  } finally {
    deleteTarget.value = null
  }
}

const handleCancelDelete = () => {
  deleteTarget.value = null
}

// 保存权限
const handleSave = async () => {
  if (!form.value.name || !form.value.resource || !form.value.action) {
    notify.warning('请填写必填项')
    return
  }

  try {
    if (isEdit.value) {
      await updatePermission(selectedPerm.value.id, form.value)
      notify.success('编辑成功')
    } else {
      await createPermission(form.value)
      notify.success('创建成功')
    }
    closeDialog()
    loadPermissions()
  } catch (error) {
    notify.error(error.message || '操作失败')
  }
}

// 关闭对话框
const closeDialog = () => {
  showAddDialog.value = false
  isEdit.value = false
  selectedPerm.value = null
  form.value = {
    name: '',
    resource: '',
    action: '',
    description: '',
    status: 1,
  }
}

onMounted(() => {
  loadPermissions()
})

// 确认对话框相关
const showConfirm = ref(false)
const deleteTarget = ref(null)

</script>

<style scoped>
.page {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
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
  color: #333;
}

.btn-primary {
  padding: 8px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-primary:hover {
  opacity: 0.9;
  transform: translateY(-2px);
}

.table-box {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  overflow: hidden;
  margin-bottom: 20px;
}

.table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

thead {
  background: #f5f7fa;
}

th {
  padding: 12px;
  text-align: left;
  font-weight: 600;
  color: #666;
  border-bottom: 1px solid #e6e9f0;
}

td {
  padding: 12px;
  border-bottom: 1px solid #e6e9f0;
  color: #333;
}

code {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 12px;
  color: #d63384;
}

.empty {
  text-align: center;
  color: #999;
  padding: 40px 12px !important;
}

.badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 12px;
  font-weight: 500;
}

.badge.success {
  background: #d4edda;
  color: #155724;
}

.badge.danger {
  background: #f8d7da;
  color: #721c24;
}

.btn-sm {
  padding: 4px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.3s;
  margin-right: 8px;
}

.btn-sm:hover {
  border-color: #667eea;
  color: #667eea;
}

.btn-sm.danger:hover {
  border-color: #f56c6c;
  color: #f56c6c;
}

.modal {
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
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  width: 90%;
  max-width: 600px;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e6e9f0;
}

.modal-header h2 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
}

.close-btn:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-row {
  display: flex;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #333;
  font-weight: 500;
  font-size: 14px;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid #e6e9f0;
  background: #f9f9f9;
}

.btn-cancel {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: white;
  color: #666;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-cancel:hover {
  border-color: #999;
}
</style>