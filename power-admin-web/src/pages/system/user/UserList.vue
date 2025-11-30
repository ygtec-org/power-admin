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
              <div class="role-cell">
                <div class="role-selector-wrapper" :class="{ active: activeRoleUserId === user.id }">
                  <div class="role-input-box" @click.stop="toggleRoleSelector(user.id)">
                    <div class="role-input-content">
                      <div v-if="(Array.isArray(userRoles[user.id]) ? userRoles[user.id] : []).length === 0" class="placeholder">选择角色</div>
                      <div v-else class="role-tags-inline">
                        <span v-for="(role, index) in Array.isArray(userRoles[user.id]) ? userRoles[user.id].slice(0, 1) : []" :key="role.id" class="role-tag-item">
                          {{ role.name }}
                          <button @click.stop="handleRemoveRole(user.id, role.id)" class="tag-remove">×</button>
                        </span>
                        <span v-if="Array.isArray(userRoles[user.id]) && userRoles[user.id].length > 1" class="role-count">+{{ userRoles[user.id].length - 1 }}</span>
                      </div>
                    </div>
                    <span class="dropdown-arrow" :class="{ active: activeRoleUserId === user.id }">▼</span>
                  </div>
                  
                  <div v-if="activeRoleUserId === user.id" class="role-dropdown-menu" @click.stop="">
                    <div v-if="roles.length === 0" class="empty-message">没有可选择的角色</div>
                    <div v-else>
                      <div
                        v-for="role in roles"
                        :key="role.id"
                        class="dropdown-item"
                        :class="{ selected: isRoleSelected(user.id, role.id) }"
                        @click.stop="handleAddRole(user.id, role.id)"
                      >
                        {{ role.name }}
                        <span v-if="isRoleSelected(user.id, role.id)" class="checkmark">✓</span>
                      </div>
                    </div>
                  </div>
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
import { ref, onMounted, computed } from 'vue'
import { getUsers, createUser, updateUser, deleteUser, assignRolesToUser, getUserRoles } from '../../../api/user'
import { getRoles } from '../../../api/role'
import { ElMessage } from 'element-plus'
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

// 角色下拉框相关
const activeRoleUserId = ref(null) // 追踪当前打开的下拉框下止
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
        console.log('getUserRoles response:', roleRes)
        // API返回可能是 {data: [...]} 或 {data: {data: [...]}} 或 {data: {list: [...]}}
        let userRolesList = []
        if (Array.isArray(roleRes.data)) {
          userRolesList = roleRes.data
        } else if (roleRes.data && Array.isArray(roleRes.data.data)) {
          userRolesList = roleRes.data.data
        } else if (roleRes.data && Array.isArray(roleRes.data.list)) {
          userRolesList = roleRes.data.list
        }
        userRoles.value[user.id] = userRolesList
      } catch (error) {
        userRoles.value[user.id] = []
      }
    }
  } catch (error) {
    ElMessage.error(error.message || '获取用户列表失败')
  }
}

// 加载角色列表
const loadRoles = async () => {
  try {
    const res = await getRoles({ pageSize: 1000 })
    console.log('getRoles response:', res)
    // 处理不同的API返回格式
    let roleList = []
    if (Array.isArray(res.data)) {
      roleList = res.data
    } else if (res.data && Array.isArray(res.data.list)) {
      roleList = res.data.list
    } else if (res.data && res.data.data && Array.isArray(res.data.data)) {
      roleList = res.data.data
    }
    roles.value = roleList
    console.log('roles.value after loading:', roles.value)
  } catch (error) {
    console.log('获取角色列表失败:', error.message)
    ElMessage.error('获取角色列表失败')
  }
}

// 切换用户的下拉框可见性
const toggleRoleSelector = (userId) => {
  activeRoleUserId.value = activeRoleUserId.value === userId ? null : userId
}

// 判断角色是否已选中
const isRoleSelected = (userId, roleId) => {
  const userRolesList = Array.isArray(userRoles.value[userId]) ? userRoles.value[userId] : []
  return userRolesList.some(r => r.id === roleId)
}

// 添加角色
const handleAddRole = async (userId, roleId) => {
  try {
    const currentRoles = Array.isArray(userRoles.value[userId]) ? userRoles.value[userId] : []
    if (currentRoles.some(r => r.id === roleId)) {
      ElMessage.warning('该角色已选中')
      return
    }
    
    const roleToAdd = roles.value.find(r => r.id === roleId)
    if (roleToAdd) {
      userRoles.value[userId] = [...currentRoles, roleToAdd]
    }
    
    const allRoleIds = userRoles.value[userId].map(r => r.id)
    await assignRolesToUser(userId, allRoleIds)
    ElMessage.success('角色添加成功')
    activeRoleUserId.value = null
  } catch (error) {
    ElMessage.error(error.message || '角色添加失败')
  }
}

const handleRemoveRole = async (userId, roleId) => {
  if (!confirm('确定要取消该角色么?')) {
    return
  }
  
  try {
    // 获取当前用户的所有角色，然后删除指定的角色
    const currentRoles = Array.isArray(userRoles.value[userId]) ? userRoles.value[userId] : []
    const newRoleIds = currentRoles.filter(r => r.id !== roleId).map(r => r.id)
    
    // 重新分配角色（不含需要删除的那个）
    await assignRolesToUser(userId, newRoleIds)
    ElMessage.success('角色删除成功')
    
    // 更新本地角色数据
    userRoles.value[userId] = currentRoles.filter(r => r.id !== roleId)
  } catch (error) {
    ElMessage.error(error.message || '角色删除失败')
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
    ElMessage.success('删除成功')
    loadUsers()
  } catch (error) {
    ElMessage.error(error.message || '删除失败')
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
    ElMessage.warning('请填写必填项')
    return
  }

  try {
    let userId = null
    if (isEdit.value) {
      // 编辑用户
      await updateUser(selectedUser.value.id, form.value)
      ElMessage.success('编辑成功')
      userId = selectedUser.value.id
    } else {
      // 新建用户
      if (!form.value.password) {
        ElMessage.warning('新用户必须设置密码')
        return
      }
      const res = await createUser(form.value)
      ElMessage.success('创建成功')
      // 获取新用户的ID (API返回中找用户ID)
      if (res.data && res.data.id) {
        userId = res.data.id
      } else if (res.data && res.data.data && res.data.data.id) {
        userId = res.data.data.id
      }
    }
    
    // 为用户分配角色 - 确保userId有效
    if (selectedRoles.value.length > 0 && userId) {
      await assignRolesToUser(userId, selectedRoles.value)
    }
    
    closeDialog()
    loadUsers()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
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

// 关闭下拉框
const closeRoleSelector = () => {
  activeRoleUserId.value = null
}

onMounted(() => {
  loadUsers()
  loadRoles()
  // 添加文档点击事件，关闭下拉框
  document.addEventListener('click', closeRoleSelector)
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
  overflow: visible;
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
  margin-bottom: 8px;
}

.empty-roles {
  color: #999;
  font-size: 12px;
  padding: 4px 8px;
  margin-bottom: 8px;
}

.role-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.role-tag .remove-btn {
  background: none;
  border: none;
  color: white;
  cursor: pointer;
  font-size: 16px;
  font-weight: bold;
  padding: 0;
  margin-left: 4px;
  transition: transform 0.2s;
  line-height: 1;
}

.role-tag .remove-btn:hover {
  transform: scale(1.2);
}

.role-display {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.btn-add-role {
  padding: 6px 12px;
  background: white;
  border: 1px dashed #667eea;
  color: #667eea;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.3s;
  align-self: fit-content;
}

.btn-add-role:hover {
  background: #f0f5ff;
  border-color: #667eea;
  color: #667eea;
}

/* 角色选择下拉框样式 */
.role-cell {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.role-selector-wrapper {
  position: relative;
  width: 100%;
  min-width: 250px;
  z-index: 10;
}

.role-input-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
  cursor: pointer;
  transition: all 0.3s;
  min-height: 36px;
  box-sizing: border-box;
}

.role-selector-wrapper.active .role-input-box {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
  background: #fafbff;
}

.role-input-box:hover {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.05);
}

.role-input-content {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
  min-width: 0;
}

.placeholder {
  color: #999;
  font-size: 13px;
}

.role-tags-inline {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: wrap;
  max-width: 100%;
}

.role-tag-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: #f0f5ff;
  color: #667eea;
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.tag-remove {
  background: none;
  border: none;
  color: #999;
  cursor: pointer;
  font-size: 12px;
  font-weight: bold;
  padding: 0;
  margin: 0;
  transition: color 0.2s;
  line-height: 1;
  display: flex;
  align-items: center;
}

.tag-remove:hover {
  color: #667eea;
}

.role-count {
  color: #667eea;
  font-size: 12px;
  font-weight: 500;
}

.dropdown-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  color: #999;
  font-size: 10px;
  transition: transform 0.3s, color 0.3s;
  flex-shrink: 0;
}

.dropdown-arrow.active {
  transform: rotate(180deg);
  color: #667eea;
}

.role-dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  max-height: 280px;
  overflow-y: auto;
  z-index: 100;
  animation: slideDown 0.2s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.empty-message {
  padding: 16px 12px;
  text-align: center;
  color: #999;
  font-size: 12px;
}

.dropdown-item {
  padding: 10px 12px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 13px;
  color: #333;
  border-bottom: 1px solid #f5f5f5;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dropdown-item:last-child {
  border-bottom: none;
}

.dropdown-item:hover {
  background: #f5f7fa;
  color: #667eea;
  padding-left: 16px;
}

.dropdown-item.selected {
  background: #f0f5ff;
  color: #667eea;
  font-weight: 500;
}

.checkmark {
  color: #667eea;
  font-weight: bold;
  font-size: 14px;
}
</style>