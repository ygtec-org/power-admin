<template>
  <div class="page">
    <div class="page-header">
      <h1>角色管理</h1>
      <button @click="showAddDialog = true" class="btn-primary">+ 新增角色</button>
    </div>

    <div class="table-box">
      <table class="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>角色名称</th>
            <th>描述</th>
            <th>状态</th>
            <th width="250">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="roles.length === 0">
            <td colspan="5" class="empty">暂无数据</td>
          </tr>
          <tr v-for="role in roles" :key="role.id">
            <td>{{ role.id }}</td>
            <td>{{ role.name }}</td>
            <td>{{ role.description || '-' }}</td>
            <td>
              <span class="badge" :class="role.status === 1 ? 'success' : 'danger'">
                {{ role.status === 1 ? '启用' : '禁用' }}
              </span>
            </td>
            <td>
              <button @click="handleEdit(role)" class="btn-sm">编辑</button>
              <button @click="handleViewMenuAndApi(role)" class="btn-sm">权限</button>
              <button @click="handleDelete(role)" class="btn-sm danger">删除</button>
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
          <h2>{{ isEdit ? '编辑角色' : '新增角色' }}</h2>
          <button @click="closeDialog" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>角色名称 *</label>
            <input v-model="form.name" type="text" placeholder="请输入角色名称" />
          </div>
          <div class="form-group">
            <label>角色描述</label>
            <textarea v-model="form.description" placeholder="请输入角色描述"></textarea>
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

    <!-- 权限分配对话框 -->
    <div v-if="showPermDialog" class="modal">
      <div class="modal-content" style="width: 500px">
        <div class="modal-header">
          <h2>为「{{ selectedRole?.name }}」分配权限</h2>
          <button @click="showPermDialog = false" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="permission-list">
            <div v-for="perm in permissions" :key="perm.id" class="permission-item">
              <input
                type="checkbox"
                :id="'perm-' + perm.id"
                v-model="selectedPermissions"
                :value="perm.id"
              />
              <label :for="'perm-' + perm.id">{{ perm.name }} ({{ perm.resource }}/{{ perm.action }})</label>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showPermDialog = false" class="btn-cancel">取消</button>
          <button @click="handleSavePermissions" class="btn-primary">保存</button>
        </div>
      </div>
    </div>

    <!-- 右侧权限配置侧边栏 -->
    <div v-if="showMenuAndApiDialog" class="permission-sidebar-overlay" @click="showMenuAndApiDialog = false">
      <div class="permission-sidebar" @click.stop>
        <div class="sidebar-header">
          <h2>为「{{ selectedRole?.name }}」分配权限</h2>
          <button @click="showMenuAndApiDialog = false" class="close-btn">×</button>
        </div>
        
        <div class="sidebar-tabs">
          <button 
            :class="['tab-btn', { active: activePermTab === 'menu' }]"
            @click="activePermTab = 'menu'"
          >
            角色菜单
          </button>
          <button 
            :class="['tab-btn', { active: activePermTab === 'api' }]"
            @click="activePermTab = 'api'"
          >
            角色API
          </button>
        </div>

        <div class="sidebar-body">
          <!-- 菜单树形展示 -->
          <div v-if="activePermTab === 'menu'" class="menu-tree-wrapper">
            <div class="search-box">
              <input v-model="menuSearchText" type="text" placeholder="筛选菜单..." />
            </div>
            <div class="menu-tree">
              <MenuTreeNode 
                v-for="menu in rootMenus" 
                :key="menu.id"
                :menu="menu" 
                :selected-menus="selectedMenus" 
                :expanded-menus="expandedMenus"
                @toggle-expand="toggleMenuExpand"
                @update-selection="updateMenuSelection"
              />
            </div>
          </div>

          <!-- API列表 -->
          <div v-if="activePermTab === 'api'" class="api-list-wrapper">
            <div class="search-box">
              <input v-model="apiSearchText" type="text" placeholder="筛选API..." />
            </div>
            <ApiTreeGroup
              :apis="apis"
              :selected-apis="selectedApis"
              @update-selection="updateApiSelection"
            />
          </div>
        </div>

        <div class="sidebar-footer">
          <button @click="showMenuAndApiDialog = false" class="btn-cancel">取消</button>
          <button @click="handleSaveMenuAndApi" class="btn-primary">确定</button>
        </div>
      </div>
    </div> 
    <ConfirmDialog
      title="删除确认"
      :message="`确定要删除角色「${deleteTarget?.name}」吗？`"
      @confirm="handleConfirmDelete"
      @cancel="handleCancelDelete"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getRoles, createRole, updateRole, deleteRole, assignPermissions, getRolePermissions } from '../../../api/role'
import { getPermissions } from '../../../api/permission'
import { getMenuTree, getAPIs } from '../../../api/menu'
import { ElMessage } from 'element-plus'
import ConfirmDialog from '../../../components/ConfirmDialog.vue'
import Pagination from '../../../components/Pagination.vue'
import MenuTreeNode from './MenuTreeNode.vue'
import ApiTreeGroup from './ApiTreeGroup.vue'

const roles = ref([])
const permissions = ref([])
const menus = ref([])
const rootMenus = ref([])
const apis = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const showAddDialog = ref(false)
const showMenuAndApiDialog = ref(false)
const activePermTab = ref('menu')
const menuSearchText = ref('')
const apiSearchText = ref('')
const expandedMenus = ref({})
const isEdit = ref(false)
const selectedRole = ref(null)
const selectedMenus = ref([])
const selectedApis = ref([])
const form = ref({
  name: '',
  description: '',
  status: 1,
})

// 确认对话框相关
const showConfirm = ref(false)
const deleteTarget = ref(null)

// 获取角色列表
const loadRoles = async (pageNum = page.value, pageSizeNum = pageSize.value) => {
  try {
    const res = await getRoles({ page: pageNum, pageSize: pageSizeNum })
    roles.value = res.data.list || []
    total.value = res.data.total || 0
    page.value = pageNum
    pageSize.value = pageSizeNum
  } catch (error) {
    notify.error(error.message || '获取角色列表失败')
  }
}

const handlePageChange = (pageNum) => {
  loadRoles(pageNum, pageSize.value)
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadRoles(1, size) // 页码重置为第一页
}

// 获取权限列表
const loadPermissions = async () => {
  try {
    const res = await getPermissions({ pageSize: 1000 })
    permissions.value = res.data.list || []
  } catch (error) {
    console.log('获取权限列表失败:', error.message)
  }
}

// 编辑角色
const handleEdit = (role) => {
  isEdit.value = true
  selectedRole.value = role
  form.value = {
    name: role.name,
    description: role.description,
    status: role.status,
  }
  showAddDialog.value = true
}

// 删除角色
const handleDelete = async (role) => {
  deleteTarget.value = role
  showConfirm.value = true
}

const handleConfirmDelete = async () => {
  try {
    await deleteRole(deleteTarget.value.id)
    notify.success('删除成功')
    loadRoles()
  } catch (error) {
    notify.error(error.message || '删除失败')
  } finally {
    deleteTarget.value = null
  }
}

const handleCancelDelete = () => {
  deleteTarget.value = null
}

// 查看权限
// 加载菜单列表
const loadMenus = async () => {
  try {
    const res = await getMenuTree()
    // getMenuTree接口直接返回已构建好的树形结构
    rootMenus.value = res.data.list || res.data.data || []
  } catch (error) {
    console.log('获取菜单列表失败:', error.message)
  }
}

// 切换菜单展开状态
const toggleMenuExpand = (menuId) => {
  if (!expandedMenus.value[menuId]) {
    expandedMenus.value[menuId] = true
  } else {
    expandedMenus.value[menuId] = false
  }
}

// 加载API列表
const loadApis = async () => {
  try {
    const res = await getAPIs({ pageSize: 1000 })
    apis.value = res.data.list || []
  } catch (error) {
    console.log('获取API列表失败:', error.message)
  }
}

// 更新菜单选择状态
const updateMenuSelection = (newSelection) => {
  // 递归查找父菜单
const findParents = (menuId, allMenus) => {
    const parents = []
    const walk = (menus) => {
      for (let menu of menus) {
        if (menu.id === menuId) {
          return true
        }
        if (menu.children && menu.children.length > 0) {
          if (walk(menu.children)) {
            parents.push(menu.id)
            return true
          }
        }
      }
      return false
    }
    walk(allMenus)
    return parents
  }
  
  // 检查是否需要自动选中父菜单
  let finalSelection = [...newSelection]
  
  // 对于每个已选中的菜单，扮零其所有父菜单
  newSelection.forEach(menuId => {
    const parents = findParents(menuId, rootMenus.value)
    parents.forEach(parentId => {
      if (!finalSelection.includes(parentId)) {
        finalSelection.push(parentId)
      }
    })
  })
  
  selectedMenus.value = finalSelection
}

// 更新API选择状态
const updateApiSelection = (newSelection) => {
  selectedApis.value = newSelection
}

// 查看菜单和API
const handleViewMenuAndApi = async (role) => {
  selectedRole.value = role
  selectedMenus.value = []
  selectedApis.value = []
  activePermTab.value = 'menu'
  expandedMenus.value = {}
  menuSearchText.value = ''
  apiSearchText.value = ''
  showMenuAndApiDialog.value = true

  // 延迟一下，确保dialog已经打开，然后加载权限数据
  setTimeout(async () => {
    try {
      const res = await fetch(`/api/admin/system/roles/${role.id}/menu-and-api`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })
      const data = await res.json()
      console.log('角色权限回显数据:', data) // 添加日志便于调试
      if (data.code === 0 && data.data) {
        // 初始化已选的菜单和API
        selectedMenus.value = data.data.selectedMenuIds || []
        selectedApis.value = data.data.selectedApiIds || []
        console.log('已选菜单:', selectedMenus.value)
        console.log('已选API:', selectedApis.value)
      }
    } catch (error) {
      console.log('获取角色权限失败:', error.message)
    }
  }, 200)
}

// 保存菜单咊API权限
const handleSaveMenuAndApi = async () => {
  if (!selectedRole.value) return
  try {
    // 过滤出有效的ID，仇止字符串司数字串转整数
const menuIds = selectedMenus.value
      .map(id => typeof id === 'string' ? parseInt(id, 10) : id)
      .filter(id => !isNaN(id) && id > 0)
    const apiIds = selectedApis.value
      .map(id => typeof id === 'string' ? parseInt(id, 10) : id)
      .filter(id => !isNaN(id) && id > 0)
    
    // 调用后端 API 保存权限
    const res = await fetch(`/api/admin/system/roles/${selectedRole.value.id}/menu-and-api`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        roleId: selectedRole.value.id,
        menuIds: menuIds,
        apiIds: apiIds
      })
    })
    const data = await res.json()
    if (data.code === 0) {
      notify.success('菜单咊API权限分配成功')
      showMenuAndApiDialog.value = false
    } else {
      notify.error(data.msg || '权限分配失败')
    }
  } catch (error) {
    notify.error(error.message || '权限分配失败')
  }
}

// 保存角色
const handleSave = async () => {
  if (!form.value.name) {
    notify.warning('请输入角色名称')
    return
  }

  try {
    if (isEdit.value) {
      await updateRole(selectedRole.value.id, form.value)
      notify.success('编辑成功')
    } else {
      await createRole(form.value)
      notify.success('创建成功')
    }
    closeDialog()
    loadRoles()
  } catch (error) {
    notify.error(error.message || '操作失败')
  }
}

// 关闭对话框
const closeDialog = () => {
  showAddDialog.value = false
  isEdit.value = false
  selectedRole.value = null
  form.value = {
    name: '',
    description: '',
    status: 1,
  }
}

onMounted(() => {
  loadRoles()
  loadPermissions()
  loadMenus()
  loadApis()
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

/* 模态框样式 */
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
  display: flex;
  align-items: center;
  justify-content: center;
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

.permission-list {
  max-height: 400px;
  overflow-y: auto;
}

.permission-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #e6e9f0;
}

.permission-item input {
  width: auto;
  margin-right: 12px;
  cursor: pointer;
}

.permission-item label {
  margin: 0;
  cursor: pointer;
  flex: 1;
  font-size: 13px;
  color: #666;
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
  transition: all 0.3s;
}

.btn-cancel:hover {
  border-color: #999;
}

/* 浒权限配置侧边栏 */
.permission-sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 2000;
  animation: fadeIn 0.3s ease-in;
}

.permission-sidebar {
  position: fixed;
  right: 0;
  top: 0;
  bottom: 0;
  width: 500px;
  background: white;
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  animation: slideInRight 0.3s ease-out;
  z-index: 2001;
}

@keyframes slideInRight {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e6e9f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sidebar-header h2 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.sidebar-tabs {
  display: flex;
  border-bottom: 1px solid #e6e9f0;
  background: #f5f7fa;
}

.tab-btn {
  flex: 1;
  padding: 12px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  border-bottom: 3px solid transparent;
  transition: all 0.3s;
}

.tab-btn.active {
  color: #667eea;
  border-bottom-color: #667eea;
  background: white;
}

.tab-btn:hover {
  color: #667eea;
}

.sidebar-body {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.search-box {
  margin-bottom: 16px;
}

.search-box input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.search-box input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

/* 菜单树形 */
.menu-tree-wrapper,
.api-list-wrapper {
  height: 100%;
}

.menu-tree {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.tree-item {
  border-left: 1px solid transparent;
  padding-left: 0;
}

.tree-node {
  display: flex;
  align-items: center;
  padding: 8px 0;
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
}

.tree-node:hover {
  background: #f5f7fa;
}

.expand-icon {
  display: inline-block;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 12px;
  transition: transform 0.2s;
  cursor: pointer;
}

.expand-icon.expanded {
  transform: rotate(90deg);
}

.expand-icon.placeholder {
  cursor: default;
}

.tree-node input {
  width: auto;
  margin: 0 8px 0 4px;
  cursor: pointer;
}

.tree-node label {
  margin: 0;
  flex: 1;
  cursor: pointer;
  font-size: 13px;
  color: #333;
}

.tree-children {
  padding-left: 24px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  background: #f9f9f9;
  margin-left: 12px;
  border-left: 2px solid #e6e9f0;
  padding-left: 12px;
  margin-left: 12px;
}

.child-node {
  padding: 4px 0;
}

/* API列表 */
.api-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.api-item {
  display: flex;
  align-items: flex-start;
  padding: 12px;
  border: 1px solid #e6e9f0;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.api-item:hover {
  border-color: #667eea;
  background: #f5f7fa;
}

.api-item input {
  width: auto;
  margin-right: 12px;
  margin-top: 2px;
  cursor: pointer;
}

.api-item label {
  margin: 0;
  flex: 1;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.api-name {
  font-size: 13px;
  color: #333;
  font-weight: 500;
}

.api-method {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 2px;
  width: fit-content;
  color: white;
}

.api-method.method-get {
  background: #61affe;
}

.api-method.method-post {
  background: #49cc90;
}

.api-method.method-put {
  background: #fca130;
}

.api-method.method-delete {
  background: #f93e3e;
}

.api-path {
  font-size: 12px;
  color: #999;
  word-break: break-all;
}

.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid #e6e9f0;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  background: #f9f9f9;
}
</style>
