<template>
  <div class="page">
    <div class="page-header">
      <h1>菜单管理</h1>
      <button @click="openAddDialog(0)" class="btn-primary">+ 新增菜单</button>
    </div>
    <div class="table-box">
      <table class="table">
        <thead>
          <tr>
            <th>Id</th>
            <th>菜单名称</th>
            <th>路径</th>
            <th>组件</th>
            <th>图标</th>
            <th>排序</th>
            <th>状态</th>
            <th width="250">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="menus.length === 0">
            <td colspan="7" class="empty">暂无数据</td>
          </tr>
          <template v-for="menu in menus" :key="menu.id">
            <MenuTreeNode :menu="menu" :level="0" @edit="handleEdit" @delete="handleDelete" @add-child="openAddDialog" />
          </template>
        </tbody>
      </table>
    </div>
    <div v-if="showAddDialog" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ isEdit ? '编辑菜单' : '新增菜单' }}</h2>
          <button @click="closeDialog" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="form-row">
            <div class="form-group" style="flex: 1; margin-right: 16px">
              <label>菜单名称 <span style="color: red">*</span></label>
              <input v-model="form.menu_name" type="text" placeholder="请输入菜单名称" />
            </div>
            <div class="form-group" style="flex: 1">
              <label>菜单路径 <span style="color: red">*</span></label>
              <input v-model="form.menu_path" type="text" placeholder="请输入菜单路径" />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group" style="flex: 1; margin-right: 16px">
              <label>组件路径</label>
              <input v-model="form.component" type="text" placeholder="可选：组件路径" />
            </div>
            <div class="form-group" style="flex: 1">
              <label>父级菜单</label>
              <select v-model.number="form.parent_id">
                <option v-for="opt in parentOptions" :key="opt.id" :value="opt.id">{{ opt.label }}</option>
              </select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group" style="flex: 1; margin-right: 16px">
              <label>菜单图标</label>
              <div style="display: flex; gap: 8px">
                <input v-model="form.icon" type="text" placeholder="可选：菜单图标，如: mdi:home" style="flex: 1" />
                <button type="button" @click="showIconPicker = true" class="btn-primary" style="padding: 8px 12px">🔍</button>
              </div>
            </div>
            <div class="form-group" style="flex: 1">
              <label>排序</label>
              <input v-model.number="form.sort" type="number" placeholder="排序号" />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group" style="flex: 1; margin-right: 16px">
              <label>菜单类型</label>
              <select v-model.number="form.menu_type">
                <option :value="1">菜单</option>
                <option :value="2">按钮</option>
              </select>
            </div>
            <div class="form-group" style="flex: 1">
              <label>状态</label>
              <select v-model.number="form.status">
                <option :value="1">显示</option>
                <option :value="0">隐藏</option>
              </select>
            </div>
          </div>
          <div class="form-group">
            <label>备注</label>
            <textarea v-model="form.remark" placeholder="可选：菜单备注" rows="3" />
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeDialog" class="btn-cancel">取消</button>
          <button @click="handleSave" class="btn-primary">保存</button>
        </div>
      </div>
      <div
        v-if="showIconPicker"
        class="icon-picker-modal"
        @click="showIconPicker = false"
      >
        <div class="icon-picker-content" @click.stop>
          <div class="icon-picker-header">
            <h3>选择图标</h3>
            <button class="close-btn" @click="showIconPicker = false">×</button>
          </div>
          <input 
            v-model="iconSearch" 
            placeholder="搜索图标，如：user, home" 
            class="icon-search-input"
          />
          <div class="icon-grid">
            <div 
              v-for="icon in filteredIcons" 
              :key="icon" 
              class="icon-item"
              @click="selectIcon(icon)"
            >
              <Icon :icon="icon" width="24" height="24" />
              <span class="icon-name">{{ icon.split(':')[1] }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    <ConfirmDialog 
      v-model:visible="showConfirm"
      title="删除确认"
      :message="`确定要删除菜单「${deleteTarget?.menuName}」吗？`"
      @confirm="handleConfirmDelete"
      @cancel="handleCancelDelete"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { getMenuTree, createMenu, updateMenu, deleteMenu } from '../../../api/menu'
import { notify } from '../../../components/Notification'
import MenuTreeNode from './MenuTreeNode.vue'
import { Icon } from '@iconify/vue'
import ConfirmDialog from '../../../components/ConfirmDialog.vue'

const showIconPicker = ref(false)
const iconSearch = ref('')

const popularIcons = [
  'mdi:home', 'mdi:account', 'mdi:cog', 'mdi:lock', 'mdi:menu',
  'mdi:view-list', 'mdi:plus', 'mdi:pencil', 'mdi:delete', 'mdi:magnify',
  'mdi:bell', 'mdi:email', 'mdi:calendar', 'mdi:chart-bar', 'mdi:file-document',
  'mdi:folder', 'mdi:image', 'mdi:video', 'mdi:music', 'mdi:book',
  'mdi:map-marker', 'mdi:heart', 'mdi:star', 'mdi:share', 'mdi:download',
  'mdi:upload', 'mdi:print', 'mdi:cloud', 'mdi:database', 'mdi:server'
]

const filteredIcons = computed(() => {
  const search = iconSearch.value.toLowerCase().trim()
  if (!search) return popularIcons
  return popularIcons.filter(icon => icon.includes(search))
})
const selectIcon = (icon) => { 
  form.value.icon = icon
  showIconPicker.value = false
}

const menus = ref([])
const parentOptions = ref([])

const buildParentOptions = (list) => {
  const result = []
  const walk = (items, level = 0) => {
    items.forEach(it => {
      result.push({ id: it.id, label: `${'  '.repeat(level)}${it.menuName}` })
      if (it.children && it.children.length) walk(it.children, level + 1)
    })
  }
  walk(list || [])
  // 顶级选项
  parentOptions.value = [{ id: 0, label: '无（顶级）' }, ...result]
}
const showAddDialog = ref(false)
const isEdit = ref(false)
const selectedMenu = ref(null)
const form = ref({
  id:0,
  menu_name: '',
  menu_path: '',
  component: '',
  icon: '',
  sort: 0,
  menu_type: 1,
  status: 1,
  remark: '',
  parent_id: 0
})

const loadMenus = async () => {
  try {
    const res = await getMenuTree()
    // 统一解析：data.list 与 data.total
    menus.value = res.data.list || res.data.data || []
    buildParentOptions(menus.value)
  } catch (error) {
    console.error('获取菜单失败:', error)
  }
}

const openAddDialog = (parentId) => {
  isEdit.value = false
  selectedMenu.value = null
  form.value = {
    id:0,
    menu_name: '',
    menu_path: '',
    component: '',
    icon: '',
    sort: 0,
    menu_type: 1,
    status: 1,
    remark: '',
    parent_id: parentId
  }
  showAddDialog.value = true
}

const handleEdit = (menu) => {
  isEdit.value = true
  selectedMenu.value = menu
  form.value = {
    id:menu.id,
    menu_name: menu.menuName,
    menu_path: menu.menuPath,
    component: menu.component || '',
    icon: menu.icon || '',
    sort: menu.sort || 0,
    menu_type: menu.menuType || 1,
    status: menu.status || 1,
    remark: menu.remark || '',
    parent_id: menu.parentId || 0
  }
  showAddDialog.value = true
}

const handleDelete = async (menu) => {
  deleteTarget.value = menu
  showConfirm.value = true
}

const handleConfirmDelete = async () => {
  try {
    await deleteMenu(deleteTarget.value.id)
    notify.success('删除成功')
    loadMenus()
  } catch (error) {
    notify.error(error.message || '删除失败')
  } finally {
    deleteTarget.value = null
  }
}

const handleCancelDelete = () => {
  deleteTarget.value = null
}

const handleSave = async () => {
  if (!form.value.menu_name || !form.value.menu_path) {
    notify.warning('请填写菜单名称和菜单路径')
    return
  }
  try {
    if (isEdit.value) {
      await updateMenu(selectedMenu.value.id, form.value)
      notify.success('编辑成功')
    } else {
      await createMenu(form.value)
      notify.success('创建成功')
    }
    closeDialog()
    loadMenus()
  } catch (error) {
    notify.error(error.message || '操作失败')
  }
}

const closeDialog = () => {
  showAddDialog.value = false
  isEdit.value = false
  selectedMenu.value = null
  form.value = {
    id: 0,
    menu_name: '',
    menu_path: '',
    component: '',
    icon: '',
    sort: 0,
    menu_type: 1,
    status: 1,
    remark: '',
    parent_id: 0
  }
}

onMounted(loadMenus)

// 确认对话框相关
const showConfirm = ref(false)
const deleteTarget = ref(null)

</script>
<style scoped>
.form-row { display: flex; gap: 16px; margin-bottom: 16px; }
.form-row .form-group { flex: 1; margin-bottom: 0; }

.page { animation: fadeIn 0.3s ease-in; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h1 { margin: 0; font-size: 24px; color: #333; }
.btn-primary { padding: 8px 16px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; border: none; border-radius: 4px; cursor: pointer; }

.table-box { background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); overflow: auto; margin-bottom: 20px; }
.table { width: 100%; border-collapse: collapse; font-size: 15px; padding-left:15px}
thead { background: #f5f7fa; position: sticky; top: 0; z-index: 1; }
th { padding: 14px; text-align: left; font-weight: 600; color: #666; border-bottom: 1px solid #e6e9f0; }
td { padding: 16px; border-bottom: 1px solid #e6e9f0; color: #333; }
.table tbody tr:nth-child(odd) { background: #fcfcff; }
.table tbody tr:hover { background: #f9fbff; }

code { background: #f5f7fa; padding: 2px 6px; border-radius: 3px; font-size: 12px; color: #d63384; }
.empty { text-align: center; color: #999; padding: 40px 12px !important; }

.btn-sm { padding: 4px 12px; border: 1px solid #ddd; background: white; border-radius: 3px; cursor: pointer; font-size: 12px; margin-right: 4px; }
.btn-sm:hover { border-color: #667eea; color: #667eea; }
.btn-sm.danger:hover { border-color: #f56c6c; color: #f56c6c; }

.pagination { display: flex; justify-content: space-between; padding: 16px 20px; background: white; border-radius: 8px; }
.pager { display: flex; gap: 12px; align-items: center; }
.pager button { padding: 6px 12px; border: 1px solid #ddd; background: white; border-radius: 3px; cursor: pointer; }
.pager button:disabled { opacity: 0.5; cursor: not-allowed; }

.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-content { background: white; border-radius: 8px; width: 90%; max-width: 700px; }
.modal-header { display: flex; justify-content: space-between; padding: 20px; border-bottom: 1px solid #e6e9f0; }
.modal-header h2 { margin: 0; font-size: 18px; }
.close-btn { background: none; border: none; font-size: 24px; cursor: pointer; }
.modal-body { padding: 20px; }
.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 8px; font-weight: 500; }
.form-group input,
.form-group select,
.form-group textarea { width: 100%; padding: 8px 12px; border: 1px solid #ddd; border-radius: 4px; font-family: inherit; }
.modal-footer { display: flex; justify-content: flex-end; gap: 12px; padding: 20px; border-top: 1px solid #e6e9f0; }
.btn-cancel { padding: 8px 16px; border: 1px solid #ddd; background: white; border-radius: 4px; cursor: pointer; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
.tree-indent { display: inline-block; width: 20px; }
.menu-name { display: flex; align-items: center; gap: 8px; }
</style>
<style scoped>
.icon-picker-modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1001; }
.icon-picker-content { background: white; border-radius: 8px; width: 90%; max-width: 600px; max-height: 80vh; overflow: hidden; display: flex; flex-direction: column; }
.icon-picker-header { display: flex; justify-content: space-between; padding: 20px; border-bottom: 1px solid #e6e9f0; }
.icon-picker-header h3 { margin: 0; font-size: 18px; }
.icon-search-input { width: calc(100% - 40px); margin: 20px auto; padding: 10px 15px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; }
.icon-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(80px, 1fr)); gap: 15px; padding: 20px; overflow-y: auto; flex: 1; }
.icon-item { display: flex; flex-direction: column; align-items: center; gap: 8px; padding: 12px 8px; border: 1px solid #eee; border-radius: 6px; cursor: pointer; background: #fff; transition: all 0.2s; }
.icon-item:hover { border-color: #667eea; background: #f0f4ff; transform: translateY(-2px); }
.icon-name { font-size: 11px; color: #666; text-align: center; overflow: hidden; text-overflow: ellipsis; width: 100%; }
</style>
