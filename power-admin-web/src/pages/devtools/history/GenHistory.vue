<template>
  <div class="history-container">
    <el-card>
      <template #header>
        <span>代码生成历史</span>
      </template>

      <!-- 搜索区域 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="表名">
          <el-input v-model="searchForm.tableName" placeholder="请输入表名" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <el-table :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="tableName" label="表名" />
        <el-table-column prop="filePath" label="文件路径" show-overflow-tooltip />
        <el-table-column prop="fileType" label="文件类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getFileTypeTag(row.fileType)">{{ row.fileType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="operator" label="操作人" width="100" />
        <el-table-column prop="createdAt" label="生成时间" width="170" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewContent(row)">查看</el-button>
            <el-button link type="danger" @click="deleteHistory(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="searchForm.page"
          v-model:page-size="searchForm.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>

    <!-- 内容查看对话框 -->
    <el-dialog v-model="showContentDialog" :title="currentFile.filePath" width="80%" top="5vh">
      <el-input
        v-model="currentFile.content"
        type="textarea"
        :rows="30"
        readonly
        style="font-family: 'Courier New', monospace;"
      />
      <template v-if="currentFile.status === 0">
        <el-divider />
        <el-alert
          title="生成失败"
          type="error"
          :description="currentFile.errorMsg"
          show-icon
          :closable="false"
        />
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/api/request'

const searchForm = reactive({
  tableName: '',
  page: 1,
  pageSize: 10
})

const tableData = ref([])
const total = ref(0)
const showContentDialog = ref(false)
const currentFile = ref({
  filePath: '',
  content: '',
  status: 1,
  errorMsg: ''
})

onMounted(() => {
  loadList()
})

const loadList = async () => {
  try {
    const res = await request.get('/codegen/history/list', { params: searchForm })
    // 后端返回 { total, data }，响应拦截器包装后是 res.data = { total, data }
    tableData.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (err) {
    ElMessage.error('加载列表失败')
  }
}

const resetSearch = () => {
  searchForm.tableName = ''
  searchForm.page = 1
  loadList()
}

const getFileTypeTag = (type) => {
  const typeMap = {
    api: 'primary',
    model: 'success',
    logic: 'warning',
    repository: 'info',
    handler: 'danger'
  }
  return typeMap[type] || 'info'
}

const viewContent = async (row) => {
  try {
    const res = await request.get(`/codegen/history/${row.id}`)
    currentFile.value = res.data
    showContentDialog.value = true
  } catch (err) {
    ElMessage.error('加载内容失败')
  }
}

const deleteHistory = (row) => {
  ElMessageBox.confirm('确定删除此历史记录吗？', '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.delete(`/codegen/history/${row.id}`)
      ElMessage.success('删除成功')
      loadList()
    } catch (err) {
      ElMessage.error('删除失败')
    }
  })
}

// 组件卸载前清理
onBeforeUnmount(() => {
  showContentDialog.value = false
})
</script>

<style scoped>
.history-container {
  padding: 20px;
}

.search-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
