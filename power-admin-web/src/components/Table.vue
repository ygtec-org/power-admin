<template>
  <div class="table-container">
    <div class="table-header">
      <h2>{{ title }}</h2>
      <button @click="handleAdd" class="btn-primary">+ 新增</button>
    </div>

    <div class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th v-for="col in columns" :key="col.prop" :width="col.width">{{ col.label }}</th>
            <th width="200">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in data" :key="row.id">
            <td v-for="col in columns" :key="col.prop">
              {{ row[col.prop] }}
            </td>
            <td class="actions">
              <button @click="handleEdit(row)" class="btn-edit">编辑</button>
              <button @click="handleDelete(row)" class="btn-delete">删除</button>
            </td>
          </tr>
          <tr v-if="data.length === 0">
            <td :colspan="columns.length + 1" class="no-data">暂无数据</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="pagination">
      <span>共 {{ total }} 条</span>
      <div class="pager">
        <button :disabled="page === 1" @click="page--">上一页</button>
        <span>第 {{ page }} 页 / 共 {{ Math.ceil(total / pageSize) }} 页</span>
        <button :disabled="page >= Math.ceil(total / pageSize)" @click="page++">下一页</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Column {
  prop: string
  label: string
  width?: string
}

interface Props {
  title: string
  columns: Column[]
  data: any[]
  total: number
  pageSize?: number
}

interface Emits {
  (e: 'add'): void
  (e: 'edit', row: any): void
  (e: 'delete', row: any): void
}

withDefaults(defineProps<Props>(), {
  pageSize: 10,
})

defineEmits<Emits>()

const page = ref(1)
const pageSize = ref(10)

const handleAdd = () => {
  // 触发 add 事件
}

const handleEdit = (row: any) => {
  // 触发 edit 事件
}

const handleDelete = (row: any) => {
  if (confirm('确定要删除吗？')) {
    // 触发 delete 事件
  }
}
</script>

<style scoped>
.table-container {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e6e9f0;
}

.table-header h2 {
  margin: 0;
  font-size: 18px;
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
  transition: all 0.3s;
}

.btn-primary:hover {
  opacity: 0.9;
  transform: translateY(-2px);
}

.table-wrapper {
  overflow-x: auto;
}

table {
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

.no-data {
  text-align: center;
  color: #999;
}

.actions {
  display: flex;
  gap: 8px;
}

.btn-edit,
.btn-delete {
  padding: 4px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.3s;
}

.btn-edit:hover {
  color: #667eea;
  border-color: #667eea;
}

.btn-delete:hover {
  color: #f56c6c;
  border-color: #f56c6c;
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-top: 1px solid #e6e9f0;
  background: #f9f9f9;
  font-size: 14px;
}

.pager {
  display: flex;
  align-items: center;
  gap: 12px;
}

.pager button {
  padding: 6px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.3s;
}

.pager button:hover:not(:disabled) {
  border-color: #667eea;
  color: #667eea;
}

.pager button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
