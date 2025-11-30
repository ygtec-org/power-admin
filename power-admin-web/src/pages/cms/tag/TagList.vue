<template>
  <div class="page-container">
    <div class="page-header">
      <h2>标签管理</h2>
      <button @click="openCreateModal" class="btn-primary">+ 新增标签</button>
    </div>
    <div class="table-container">
      <table class="table">
        <thead><tr><th>ID</th><th>标签名称</th><th>Slug</th><th>颜色</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="item in list" :key="item.Id">
            <td>{{ item.Id }}</td>
            <td>{{ item.Name }}</td>
            <td>{{ item.Slug }}</td>
            <td><span :style="{ background: item.Color || '#ccc', padding: '4px 8px', color: 'white', borderRadius: '3px' }">{{ item.Color }}</span></td>
            <td class="action-cell">
              <button @click="edit(item)" class="btn-small">编辑</button>
              <button @click="delete_(item.Id)" class="btn-small delete">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="list.length === 0" class="empty-state">暂无数据</div>
    </div>
    <div v-if="showModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ editingId ? '编辑标签' : '新增标签' }}</h3>
          <button @click="closeModal" class="btn-close">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>标签名称</label>
            <input v-model="form.Name" type="text" class="form-input">
          </div>
          <div class="form-group">
            <label>Slug</label>
            <input v-model="form.Slug" type="text" class="form-input">
          </div>
          <div class="form-group">
            <label>颜色</label>
            <input v-model="form.Color" type="color" class="form-input">
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn-cancel">取消</button>
          <button @click="save" class="btn-primary">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
const list = ref([])
const showModal = ref(false)
const editingId = ref(null)
const form = ref({ Name: '', Slug: '', Color: '#667eea' })
const load = async () => {
  try {
    const response = await fetch('/api/cms/tag/list?page=1&pageSize=100')
    const data = await response.json()
    if (data.code === 0 && data.data) list.value = data.data.Data || []
  } catch (error) { console.error('加载失败:', error) }
}
const openCreateModal = () => { editingId.value = null; form.value = { Name: '', Slug: '', Color: '#667eea' }; showModal.value = true }
const edit = (item) => { editingId.value = item.Id; form.value = { ...item }; showModal.value = true }
const closeModal = () => { showModal.value = false }
const save = async () => {
  try {
    const url = editingId.value ? `/api/cms/tag/${editingId.value}` : '/api/cms/tag'
    const method = editingId.value ? 'PUT' : 'POST'
    const response = await fetch(url, { method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(form.value) })
    const data = await response.json()
    if (data.code === 0) { alert('保存成功'); closeModal(); load() }
  } catch (error) { console.error('保存失败:', error) }
}
const delete_ = async (id) => {
  if (!confirm('确定删除？')) return
  try {
    const response = await fetch(`/api/cms/tag/${id}`, { method: 'DELETE' })
    const data = await response.json()
    if (data.code === 0) { alert('删除成功'); load() }
  } catch (error) { console.error('删除失败:', error) }
}
onMounted(() => { load() })
</script>
<style scoped>
.page-container { padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 24px; }
.btn-primary { padding: 8px 16px; background: #667eea; color: white; border: none; border-radius: 4px; cursor: pointer; }
.table-container { background: white; border-radius: 4px; overflow: auto; box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1); }
.table { width: 100%; border-collapse: collapse; }
.table th { padding: 12px; text-align: left; background: #f5f7fa; border-bottom: 2px solid #e6e9f0; font-weight: 600; }
.table td { padding: 12px; border-bottom: 1px solid #e6e9f0; }
.action-cell { display: flex; gap: 8px; }
.btn-small { padding: 4px 12px; border: 1px solid #ddd; background: white; border-radius: 3px; cursor: pointer; font-size: 12px; color: #667eea; }
.btn-small.delete { color: #ff4d4f; border-color: #ff4d4f; }
.empty-state { padding: 40px; text-align: center; color: #999; }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-content { background: white; border-radius: 8px; max-width: 500px; width: 90%; }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 20px; border-bottom: 1px solid #e6e9f0; }
.modal-header h3 { margin: 0; }
.btn-close { background: none; border: none; font-size: 24px; color: #999; cursor: pointer; }
.modal-body { padding: 20px; }
.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 8px; font-weight: 500; }
.form-input { width: 100%; padding: 8px 12px; border: 1px solid #ddd; border-radius: 4px; }
.modal-footer { display: flex; justify-content: flex-end; gap: 12px; padding: 20px; border-top: 1px solid #e6e9f0; }
.btn-cancel { padding: 8px 16px; background: white; border: 1px solid #ddd; border-radius: 4px; cursor: pointer; }
</style>