package codegen

// Vue模板
const vueTemplate = `<template>
  <div class="{{.ModuleName}}-container">
    <el-card>
      <template #header>
        <div class="header-actions">
          <span>{{.BusinessName}}</span>
          <el-button type="primary" @click="openCreateDialog">新建</el-button>
        </div>
      </template>

      <!-- 搜索区域 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        {{range .Fields}}{{if .IsQuery}}<el-form-item label="{{.ColumnComment}}">
          <el-input v-model="searchForm.{{.JsonField}}" placeholder="请输入{{.ColumnComment}}" clearable />
        </el-form-item>
        {{end}}{{end}}<el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <el-table :data="tableData" border stripe>
        {{range .Fields}}{{if .IsList}}<el-table-column prop="{{.JsonField}}" label="{{.ColumnComment}}" />
        {{end}}{{end}}<el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editItem(row)">编辑</el-button>
            <el-button link type="danger" @click="deleteItem(row)">删除</el-button>
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

    <!-- 新建/编辑对话框 -->
    <el-dialog 
      v-model="showEditDialog" 
      :title="editForm.id ? '编辑' : '新建'" 
      width="600px"
    >
      <el-form :model="editForm" :rules="formRules" ref="formRef" label-width="120px">
        {{range .Fields}}{{if or .IsInsert .IsEdit}}{{if not .IsIncrement}}<el-form-item label="{{.ColumnComment}}" prop="{{.JsonField}}">
          <el-input v-model="editForm.{{.JsonField}}" placeholder="请输入{{.ColumnComment}}" />
        </el-form-item>
        {{end}}{{end}}{{end}}
      </el-form>

      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="saveItem">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/api/request'

const searchForm = reactive({
  {{range .Fields}}{{if .IsQuery}}{{.JsonField}}: '',
  {{end}}{{end}}page: 1,
  pageSize: 10
})

const tableData = ref([])
const total = ref(0)
const showEditDialog = ref(false)
const formRef = ref(null)

const editForm = ref({
  id: 0,
  {{range .Fields}}{{if or .IsInsert .IsEdit}}{{.JsonField}}: '',
  {{end}}{{end}}
})

const formRules = {
  {{range .Fields}}{{if and .IsRequired (or .IsInsert .IsEdit)}}{{.JsonField}}: [
    { required: true, message: '请输入{{.ColumnComment}}', trigger: 'blur' }
  ],
  {{end}}{{end}}
}

// 加载列表
const loadList = async () => {
  try {
    const res = await request.get('/{{.ModuleName}}/{{.BusinessName}}/list', { params: searchForm })
    tableData.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (err) {
    console.error('加载列表失败:', err)
  }
}

// 重置搜索
const resetSearch = () => {
  {{range .Fields}}{{if .IsQuery}}searchForm.{{.JsonField}} = ''
  {{end}}{{end}}searchForm.page = 1
  loadList()
}

// 打开新建对话框
const openCreateDialog = () => {
  editForm.value = {
    id: 0,
    {{range .Fields}}{{if or .IsInsert .IsEdit}}{{.JsonField}}: '',
    {{end}}{{end}}
  }
  showEditDialog.value = true
}

// 编辑
const editItem = async (row) => {
  try {
    const res = await request.get(` + "`" + `/{{.ModuleName}}/{{.BusinessName}}/${row.id}` + "`" + `)
    editForm.value = res.data
    showEditDialog.value = true
  } catch (err) {
    console.error('加载数据失败:', err)
  }
}

// 保存
const saveItem = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    try {
      if (editForm.value.id) {
        await request.put(` + "`" + `/{{.ModuleName}}/{{.BusinessName}}/${editForm.value.id}` + "`" + `, editForm.value)
      } else {
        await request.post('/{{.ModuleName}}/{{.BusinessName}}', editForm.value)
      }
      showEditDialog.value = false
      loadList()
    } catch (err) {
      console.error('保存失败:', err)
    }
  })
}

// 删除
const deleteItem = (row) => {
  ElMessageBox.confirm('确定删除吗？', '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.delete(` + "`" + `/{{.ModuleName}}/{{.BusinessName}}/${row.id}` + "`" + `)
      loadList()
    } catch (err) {
      console.error('删除失败:', err)
    }
  })
}

// 初始化
loadList()

// 组件卸载前清理
onBeforeUnmount(() => {
  showEditDialog.value = false
})
</script>

<style scoped>
.{{.ModuleName}}-container {
  padding: 20px;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
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
`
