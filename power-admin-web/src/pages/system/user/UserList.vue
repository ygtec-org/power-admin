<template>
  <div class="page">
    <div class="page-header">
      <h1>用户管理</h1>
      <button @click="showAddDialog = true" class="btn-primary">+ 新增用户</button>
    </div>

    <div class="table-box">
      <table class="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>用户名</th>
            <th>手机号</th>
            <th>昵称</th>
            <th>邮箱</th>
            <th>状态</th>
            <th>角色</th>
            <th width="200">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="users.length === 0">
            <td colspan="7" class="empty">暂无数据</td>
          </tr>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.username || '-' }}</td>
            <td>{{ user.phone }}</td>
            <td>{{ user.nickname || '-' }}</td>
            <td>{{ user.email || '-' }}</td>
            <td><span class="badge" :class="user.status === 1 ? 'success' : 'danger'">{{ user.status === 1 ? '启用' : '禁用' }}</span></td>
            <td>
              <div class="role-tags">
                <div v-for="role in (userRoles[user.id] || [])" :key="role.id" class="role-tag">
                  <span>{{ role.name }}</span>
                  <button @click="handleRemoveRole(user.id, role.id)" class="remove-btn">×</button>
                </div>
              </div>
            </td>
            <td>
              <button @click="handleEdit(user)" class="btn-sm">编辑</button>
              <button @click="handleDelete(user)" class="btn-sm danger">删除</button>
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

    <!-- 新增/编辑用户对话框 -->
    <div v-if="showAddDialog" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ isEdit ? '编辑用户' : '新增用户' }}</h2>
          <button @click="closeDialog" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>用户名 *</label>
            <input v-model="form.username" type="text" placeholder="请输入用户名" :disabled="isEdit" />
          </div>
          <div class="form-group">
            <label>手机号 *</label>
            <input v-model="form.phone" type="text" placeholder="请输入手机号" />
          </div>
          <div class="form-group">
            <label>昵称</label>
            <input v-model="form.nickname" type="text" placeholder="请输入昵称" />
          </div>
          <div class="form-group">
            <label>邮箱</label>
            <input v-model="form.email" type="email" placeholder="请输入邮箱" />
          </div>
          <div v-if="!isEdit" class="form-group">
            <label>密码 *</label>
            <input v-model="form.password" type="password" placeholder="请输入密码" />
          </div>
          <div class="form-group">
            <label>状态</label>
            <select v-model="form.status">
              <option value="1">启用</option>
              <option value="0">禁用</option>
            </select>
          </div>
          <div class="form-group">
            <label>绑定角色</label>
            <div class="role-list">
              <div v-for="role in roles" :key="role.id" class="role-item">
                <input
                  type="checkbox"
                  :id="'role-' + role.id"
                  v-model="selectedRoles"
                  :value="role.id"
                />
                <label :for="'role-' + role.id">{{ role.name }}</label>
              </div>
            </div>
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
      :message="`确定要删除用户「${deleteTarget?.username}」吗？`"
      @confirm="handleConfirmDelete"
      @cancel="handleCancelDelete"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getUsers, createUser, updateUser, deleteUser, assignRolesToUser, getUserRoles } from '../../../api/user'
import { getRoles } from '../../../api/role'
import { notify } from '../../../components/Notification'
import ConfirmDialog from '../../../components/ConfirmDialog.vue'
import Pagination from '../../../components/Pagination.vue'

const users = ref([])
const roles = ref([])
const userRoles = ref({})
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const showAddDialog = ref(false)
const isEdit = ref(false)
const selectedRoles = ref([])
const form = ref({
  username: '',
  phone: '',
  nickname: '',
  email: '',
  password: '',
  status: 1,
})
const selectedUser = ref(null)

// 确认对话框相关
const showConfirm = ref(false)
const deleteTarget = ref(null)

// 加载用户列表
const loadUsers = async (pageNum = page.value, pageSizeNum = pageSize.value) => {
  try {
    const res = await getUsers({ page: pageNum, pageSize: pageSizeNum })
    users.value = res.data.list || []
    total.value = res.data.total || 0
    page.value = pageNum
    pageSize.value = pageSizeNum
    
    // 加载每个用户的角色
    for (const user of users.value) {
      try {
        const roleRes = await getUserRoles(user.id)
        userRoles.value[user.id] = roleRes.data || []
      } catch (error) {
        userRoles.value[user.id] = []
      }
    }
  } catch (error) {
    notify.error(error.message || '获取用户列表失败')
  }
}

// 加载角色列表
const loadRoles = async () => {
  try {
    const res = await getRoles({ pageSize: 1000 })
    roles.value = res.data.list || []
  } catch (error) {
    console.log('获取角色列表失败:', error.message)
  }
}

const handleRemoveRole = async (userId, roleId) => {
  if (!confirm('\u786e定要取消该角色么?')) {
    return
  }
  
  try {
    // 获取当前用户的所有角色，然后删除指定的角色
    const currentRoles = userRoles.value[userId] || []
    const newRoleIds = currentRoles.filter(r => r.id !== roleId).map(r => r.id)
    
    // 重新分配角色（不含需要删除的那个）
    await assignRolesToUser(userId, newRoleIds)
    notify.success('角色删除成功')
    
    // 更新本地角色数据
    userRoles.value[userId] = currentRoles.filter(r => r.id !== roleId)
  } catch (error) {
    notify.error(error.message || '角色删除失败')
  }
}

const handlePageChange = (pageNum) => {
  loadUsers(pageNum, pageSize.value)
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadUsers(1, size) // 页码重置为第一页
}

const handleEdit = async (user) => {
  isEdit.value = true
  selectedUser.value = user
  form.value = {
    username: user.username,
    phone: user.phone,
    nickname: user.nickname || '',
    email: user.email || '',
    password: '',
    status: user.status,
  }
  // 获取用户当前的角色
  try {
    const res = await getUserRoles(user.id)
    if (res.data && res.data.data) {
      selectedRoles.value = res.data.data.map((r) => r.id)
    } else {
      selectedRoles.value = []
    }
  } catch (error) {
    console.log('获取用户角色失败:', error.message)
    selectedRoles.value = []
  }
  showAddDialog.value = true
}

const handleDelete = async (user) => {
  deleteTarget.value = user
  showConfirm.value = true
}

const handleConfirmDelete = async () => {
  try {
    await deleteUser(deleteTarget.value.id)
    notify.success('删除成功')
    loadUsers()
  } catch (error) {
    notify.error(error.message || '删除失败')
  } finally {
    deleteTarget.value = null
  }
}

const handleCancelDelete = () => {
  deleteTarget.value = null
}

// 保存用户
const handleSave = async () => {
  if (!form.value.username || !form.value.phone) {
    notify.warning('请填写必填项')
    return
  }

  try {
    if (isEdit.value) {
      await updateUser(selectedUser.value.id, form.value)
      notify.success('编辑成功')
    } else {
      if (!form.value.password) {
        notify.warning('新用户必须设置密码')
        return
      }
      await createUser(form.value)
      notify.success('创建成功')
    }
    
    // 为用户分配角色
    if (selectedRoles.value.length > 0) {
      await assignRolesToUser(selectedUser.value?.id || null, selectedRoles.value)
    }
    
    closeDialog()
    loadUsers()
  } catch (error) {
    notify.error(error.message || '操作失败')
  }
}

// 关闭对话框
const closeDialog = () => {
  showAddDialog.value = false
  isEdit.value = false
  selectedUser.value = null
  selectedRoles.value = []
  form.value = {
    username: '',
    phone: '',
    nickname: '',
    email: '',
    password: '',
    status: 1,
  }
}

onMounted(() => {
  loadUsers()
  loadRoles()
})
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
  max-height: 90vh;
  overflow-y: auto;
  min-width: 500px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.modal-header {
  padding: 20px;
  border-bottom: 1px solid #e6e9f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
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
  cursor: pointer;
  color: #999;
  transition: color 0.3s;
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

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
}

.form-group input:disabled {
  background: #f5f7fa;
  cursor: not-allowed;
}

.role-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.role-item {
  display: flex;
  align-items: center;
}

.role-item input[type='checkbox'] {
  width: auto;
  margin-right: 8px;
}

.role-item label {
  margin: 0;
  cursor: pointer;
  font-weight: normal;
}

.modal-footer {
  padding: 16px 20px;
  border-top: 1px solid #e6e9f0;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.btn-cancel {
  padding: 8px 16px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.btn-cancel:hover {
  border-color: #999;
  color: #999;
}

.role-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.role-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: #e8f4f8;
  color: #0066cc;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  border: 1px solid #b3d9e6;
}

.role-tag .remove-btn {
  background: none;
  border: none;
  color: #ff4444;
  cursor: pointer;
  font-size: 14px;
  font-weight: bold;
  padding: 0;
  margin-left: 2px;
  transition: color 0.3s;
}

.role-tag .remove-btn:hover {
  color: #cc0000;
}
</style>
