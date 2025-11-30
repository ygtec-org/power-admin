<template>
  <div class="page-container">
    <h2>用户管理</h2>
    <div class="table-container">
      <table class="table">
        <thead><tr><th>ID</th><th>用户名</th><th>邮箱</th><th>状态</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="item in list" :key="item.Id">
            <td>{{ item.Id }}</td>
            <td>{{ item.Username }}</td>
            <td>{{ item.Email }}</td>
            <td><span :class="['badge', item.Status === 1 ? 'active' : 'disabled']">{{ item.Status === 1 ? '正常' : '禁用' }}</span></td>
            <td class="action-cell">
              <button v-if="item.Status === 1" @click="disable(item.Id)" class="btn-small delete">禁用</button>
              <button v-else @click="enable(item.Id)" class="btn-small">启用</button>
              <button @click="delete_(item.Id)" class="btn-small delete">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="list.length === 0" class="empty-state">暂无数据</div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
const list = ref([])
const load = async () => {
  try {
    const response = await fetch('/api/cms/user/list?page=1&pageSize=100')
    const data = await response.json()
    if (data.code === 0 && data.data) list.value = data.data.Data || []
  } catch (error) { console.error('加载失败:', error) }
}
const disable = async (id) => {
  try {
    const response = await fetch(`/api/cms/user/${id}/disable`, { method: 'POST' })
    const data = await response.json()
    if (data.code === 0) { alert('禁用成功'); load() }
  } catch (error) { console.error('禁用失败:', error) }
}
const enable = async (id) => {
  try {
    const response = await fetch(`/api/cms/user/${id}`, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ Status: 1 }) })
    const data = await response.json()
    if (data.code === 0) { alert('启用成功'); load() }
  } catch (error) { console.error('启用失败:', error) }
}
const delete_ = async (id) => {
  if (!confirm('确定删除？')) return
  try {
    const response = await fetch(`/api/cms/user/${id}`, { method: 'DELETE' })
    const data = await response.json()
    if (data.code === 0) { alert('删除成功'); load() }
  } catch (error) { console.error('删除失败:', error) }
}
onMounted(() => { load() })
</script>
<style scoped>
.page-container { padding: 24px; }
.page-container h2 { margin: 0 0 24px 0; font-size: 24px; }
.table-container { background: white; border-radius: 4px; overflow: auto; box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1); }
.table { width: 100%; border-collapse: collapse; }
.table th { padding: 12px; text-align: left; background: #f5f7fa; border-bottom: 2px solid #e6e9f0; font-weight: 600; }
.table td { padding: 12px; border-bottom: 1px solid #e6e9f0; }
.badge { padding: 4px 8px; border-radius: 3px; font-size: 12px; }
.badge.active { background: #74b9ff; color: #fff; }
.badge.disabled { background: #fab1a0; color: #d63031; }
.action-cell { display: flex; gap: 8px; }
.btn-small { padding: 4px 12px; border: 1px solid #ddd; background: white; border-radius: 3px; cursor: pointer; font-size: 12px; color: #667eea; }
.btn-small.delete { color: #ff4d4f; border-color: #ff4d4f; }
.empty-state { padding: 40px; text-align: center; color: #999; }
</style>
