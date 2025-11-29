<template>
  <div class="pagination" v-if="total > 0">
    <div class="pagination-info">
      共 {{ total }} 条记录
      <select v-model="currentPageSize" @change="handleSizeChange" class="page-size-select">
        <option v-for="size in pageSizes" :key="size" :value="size">{{ size }} 条/页</option>
      </select>
    </div>
    <div class="pager">
      <button 
        :disabled="currentPage <= 1" 
        @click="handlePageChange(currentPage - 1)"
        class="page-btn"
      >
        上一页
      </button>
      <span class="page-number">第 {{ currentPage }} 页</span>
      <button 
        :disabled="currentPage * pageSize >= total" 
        @click="handlePageChange(currentPage + 1)"
        class="page-btn"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  total: {
    type: Number,
    default: 0
  },
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  pageSizes: {
    type: Array,
    default: () => [10, 20, 50, 100]
  }
})

const emit = defineEmits(['update:currentPage', 'update:pageSize', 'page-change', 'size-change'])

// 内部响应式数据
const currentPageSize = ref(props.pageSize)

// 监听外部pageSize变化
watch(() => props.pageSize, (newVal) => {
  currentPageSize.value = newVal
})

const handlePageChange = (page) => {
  emit('update:currentPage', page)
  emit('page-change', page)
}

const handleSizeChange = () => {
  emit('update:pageSize', currentPageSize.value)
  emit('size-change', currentPageSize.value)
  // 页码重置为第一页
  emit('update:currentPage', 1)
  emit('page-change', 1)
}
</script>

<style scoped>
.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  font-size: 14px;
  margin-top: 20px;
}

.pagination-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-size-select {
  padding: 4px 8px;
  border: 1px solid #ddd;
  border-radius: 3px;
  font-size: 12px;
}

.pager {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.3s;
}

.page-btn:hover:not(:disabled) {
  border-color: #667eea;
  color: #667eea;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-number {
  font-size: 12px;
  color: #666;
}
</style>