<template>
  <div class="page">
    <div class="page-header">
      <h1>{{ title }}</h1>
      <button @click="showAddDialog = true" class="btn-primary">+ 新增</button>
    </div>

    <div class="table-box">
      <table class="table">
        <thead>
          <tr>
            <th v-for="col in columns" :key="col">{{ col }}</th>
            <th width="200">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="dataList.length === 0">
            <td :colspan="columns.length + 1" class="empty">暂无数据</td>
          </tr>
          <tr v-for="item in dataList" :key="item.id">
            <td v-for="col in columns" :key="col">{{ item[col.toLowerCase()] || '-' }}</td>
            <td>
              <button @click="handleEdit(item)" class="btn-sm">编辑</button>
              <button @click="handleDelete(item)" class="btn-sm danger">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  title: String,
  columns: Array,
  mockData: Array,
})

const dataList = ref(props.mockData || [])
const showAddDialog = ref(false)

const handleEdit = (item) => {
  console.log('Edit:', item)
}

const handleDelete = (item) => {
  if (confirm('确定要删除吗？')) {
    dataList.value = dataList.value.filter(i => i.id !== item.id)
  }
}
</script>

<style scoped>
.page {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
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
}

.btn-primary:hover {
  opacity: 0.9;
}

.table-box {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

thead {
  background: #f5f7fa;
}

th, td {
  padding: 12px;
  text-align: left;
}

th {
  font-weight: 600;
  color: #666;
  border-bottom: 1px solid #e6e9f0;
}

td {
  border-bottom: 1px solid #e6e9f0;
  color: #333;
}

.empty {
  text-align: center;
  color: #999;
  padding: 40px 12px !important;
}

.btn-sm {
  padding: 4px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
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
</style>
