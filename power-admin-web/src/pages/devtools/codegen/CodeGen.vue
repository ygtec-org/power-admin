<template>
  <div class="codegen-container">
    <el-card>
      <template #header>
        <div class="header-actions">
          <span>代码生成器</span>
          <el-button type="primary" @click="openCreateDialog">新建配置</el-button>
        </div>
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
        <el-table-column prop="businessName" label="业务名称" />
        <el-table-column prop="moduleName" label="模块名称" />
        <el-table-column prop="author" label="作者" />
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column prop="createdAt" label="创建时间" width="170" />
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="previewCode(row)">预览</el-button>
            <el-button link type="success" @click="generateCode(row)">生成</el-button>
            <el-button link type="primary" @click="editConfig(row)">编辑</el-button>
            <el-button link type="danger" @click="deleteConfig(row)">删除</el-button>
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

    <!-- 新建/编辑配置对话框 -->
    <el-dialog 
      v-model="showEditDialog" 
      :title="editForm.id ? '编辑配置' : '新建配置'" 
      width="90%" 
      top="5vh"
    >
      <el-form :model="editForm" :rules="formRules" ref="formRef" label-width="120px">
        <el-card shadow="never" class="form-section">
          <template #header>
            <span>基本信息</span>
          </template>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="表名" prop="tableName" required>
                <el-input 
                  v-model="editForm.tableName" 
                  placeholder="如: sys_user"
                  :disabled="!!editForm.id"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="表前缀">
                <el-input v-model="editForm.tablePrefix" placeholder="如: sys_" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="业务名称" prop="businessName" required>
                <el-input v-model="editForm.businessName" placeholder="如: 用户管理" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="模块名称" prop="moduleName" required>
                <el-input v-model="editForm.moduleName" placeholder="如: user" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="包路径">
                <el-input v-model="editForm.packageName" placeholder="如: com.admin.user" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="作者">
                <el-input v-model="editForm.author" placeholder="作者名称" />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="备注">
                <el-input v-model="editForm.remark" type="textarea" :rows="2" placeholder="配置备注" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>

        <el-card shadow="never" class="form-section">
          <template #header>
            <div class="section-header">
              <span>字段配置</span>
              <el-button type="primary" size="small" @click="addColumn">添加字段</el-button>
            </div>
          </template>
          
          <el-table :data="editForm.columns || []" border max-height="400">
            <el-table-column label="字段名" width="150">
              <template #default="{ row, $index }">
                <el-input v-model="row.columnName" placeholder="字段名" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="字段注释" width="150">
              <template #default="{ row }">
                <el-input v-model="row.columnComment" placeholder="注释" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="MySQL类型" width="120">
              <template #default="{ row }">
                <el-select v-model="row.columnType" placeholder="类型" size="small" @change="onColumnTypeChange(row)">
                  <el-option label="INT" value="int" />
                  <el-option label="BIGINT" value="bigint" />
                  <el-option label="VARCHAR" value="varchar(255)" />
                  <el-option label="TEXT" value="text" />
                  <el-option label="DATETIME" value="datetime" />
                  <el-option label="DECIMAL" value="decimal(10,2)" />
                  <el-option label="TINYINT" value="tinyint" />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="Go类型" width="100">
              <template #default="{ row }">
                <el-tag size="small">{{ row.goType || '' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="Go字段名" width="120">
              <template #default="{ row }">
                <el-tag size="small">{{ row.goField || '' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="主键" width="60" align="center">
              <template #default="{ row }">
                <el-checkbox v-model="row.isPk" :true-label="1" :false-label="0" />
              </template>
            </el-table-column>
            <el-table-column label="自增" width="60" align="center">
              <template #default="{ row }">
                <el-checkbox v-model="row.isIncrement" :true-label="1" :false-label="0" />
              </template>
            </el-table-column>
            <el-table-column label="必填" width="60" align="center">
              <template #default="{ row }">
                <el-checkbox v-model="row.isRequired" :true-label="1" :false-label="0" />
              </template>
            </el-table-column>
            <el-table-column label="插入" width="60" align="center">
              <template #default="{ row }">
                <el-checkbox v-model="row.isInsert" :true-label="1" :false-label="0" />
              </template>
            </el-table-column>
            <el-table-column label="编辑" width="60" align="center">
              <template #default="{ row }">
                <el-checkbox v-model="row.isEdit" :true-label="1" :false-label="0" />
              </template>
            </el-table-column>
            <el-table-column label="列表" width="60" align="center">
              <template #default="{ row }">
                <el-checkbox v-model="row.isList" :true-label="1" :false-label="0" />
              </template>
            </el-table-column>
            <el-table-column label="查询" width="60" align="center">
              <template #default="{ row }">
                <el-checkbox v-model="row.isQuery" :true-label="1" :false-label="0" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80" fixed="right">
              <template #default="{ $index }">
                <el-button 
                  link 
                  type="danger" 
                  size="small"
                  @click="removeColumn($index)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-form>

      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="saveConfig">保存配置</el-button>
      </template>
    </el-dialog>

    <!-- 代码预览对话框 -->
    <el-dialog v-model="showPreviewDialog" title="代码预览" width="90%" top="5vh">
      <el-tabs v-model="activeTab">
        <el-tab-pane 
          v-for="file in previewFiles" 
          :key="file.filePath" 
          :label="file.filePath" 
          :name="file.filePath"
        >
          <el-input
            v-model="file.content"
            type="textarea"
            :rows="25"
            readonly
            style="font-family: 'Courier New', monospace;"
          />
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/api/request'

const searchForm = reactive({
  tableName: '',
  page: 1,
  pageSize: 10
})

const tableData = ref([])
const total = ref(0)
const showEditDialog = ref(false)
const showPreviewDialog = ref(false)
const previewFiles = ref([])
const activeTab = ref('')
const formRef = ref(null)

const editForm = ref({
  id: 0,
  tableName: '',
  tablePrefix: '',
  businessName: '',
  moduleName: '',
  packageName: '',
  author: '',
  remark: '',
  columns: []
})

const formRules = {
  tableName: [
    { required: true, message: '请输入表名', trigger: 'blur' }
  ],
  businessName: [
    { required: true, message: '请输入业务名称', trigger: 'blur' }
  ],
  moduleName: [
    { required: true, message: '请输入模块名称', trigger: 'blur' }
  ]
}

// 加载列表
const loadList = async () => {
  try {
    const res = await request.get('/codegen/config/list', { params: searchForm })
    // 后端返回 { total, data }，响应拦截器包装后是 res.data = { total, data }
    tableData.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (err) {
    ElMessage.error('加载列表失败')
  }
}

// 重置搜索
const resetSearch = () => {
  searchForm.tableName = ''
  searchForm.page = 1
  loadList()
}

// 打开新建对话框
const openCreateDialog = () => {
  editForm.value = {
    id: 0,
    tableName: '',
    tablePrefix: '',
    businessName: '',
    moduleName: '',
    packageName: '',
    author: 'admin',
    remark: '',
    columns: [
      createDefaultColumn('id', 'ID', 'bigint', true, true)
    ]
  }
  showEditDialog.value = true
}

// 创建默认字段
const createDefaultColumn = (name, comment, type, isPk = false, isIncrement = false) => {
  const goType = mysqlTypeToGoType(type)
  const goField = columnNameToGoField(name)
  
  return {
    columnName: name,
    columnComment: comment,
    columnType: type,
    goType: goType,
    goField: goField,
    isPk: isPk ? 1 : 0,
    isIncrement: isIncrement ? 1 : 0,
    isRequired: 1,
    isInsert: isPk ? 0 : 1,
    isEdit: isPk ? 0 : 1,
    isList: 1,
    isQuery: isPk ? 0 : 1,
    queryType: '=',
    htmlType: getHtmlType(goType),
    dictType: '',
    sort: 0
  }
}

// 添加字段
const addColumn = () => {
  if (!editForm.value.columns) {
    editForm.value.columns = []
  }
  editForm.value.columns.push({
    columnName: '',
    columnComment: '',
    columnType: 'varchar(255)',
    goType: 'string',
    goField: '',
    isPk: 0,
    isIncrement: 0,
    isRequired: 1,
    isInsert: 1,
    isEdit: 1,
    isList: 1,
    isQuery: 1,
    queryType: '=',
    htmlType: 'input',
    dictType: '',
    sort: editForm.value.columns.length
  })
}

// 删除字段
const removeColumn = (index) => {
  editForm.value.columns.splice(index, 1)
}

// 字段类型改变
const onColumnTypeChange = (row) => {
  const baseType = row.columnType.split('(')[0]
  row.goType = mysqlTypeToGoType(baseType)
  row.goField = columnNameToGoField(row.columnName)
  row.htmlType = getHtmlType(row.goType)
}

// MySQL类型转Go类型
const mysqlTypeToGoType = (mysqlType) => {
  const typeMap = {
    'int': 'int64',
    'tinyint': 'int',
    'bigint': 'int64',
    'varchar': 'string',
    'text': 'string',
    'datetime': 'time.Time',
    'timestamp': 'time.Time',
    'decimal': 'float64',
    'double': 'float64',
    'float': 'float32'
  }
  return typeMap[mysqlType.toLowerCase()] || 'string'
}

// 字段名转Go字段名
const columnNameToGoField = (columnName) => {
  if (!columnName) return ''
  return columnName.split('_').map(part => {
    return part.charAt(0).toUpperCase() + part.slice(1).toLowerCase()
  }).join('')
}

// 获取HTML类型
const getHtmlType = (goType) => {
  if (goType === 'time.Time') return 'datetime'
  if (goType === 'int' || goType === 'int64' || goType === 'float64') return 'input'
  return 'input'
}

// 保存配置
const saveConfig = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    if (editForm.value.columns.length === 0) {
      ElMessage.warning('请至少添加一个字段')
      return
    }

    // 更新字段排序
    editForm.value.columns.forEach((col, index) => {
      col.sort = index + 1
      // 确保goField和goType已设置
      if (!col.goField) {
        col.goField = columnNameToGoField(col.columnName)
      }
      if (!col.goType) {
        const baseType = col.columnType.split('(')[0]
        col.goType = mysqlTypeToGoType(baseType)
      }
    })

    try {
      if (editForm.value.id) {
        await request.put(`/codegen/config/${editForm.value.id}`, editForm.value)
        ElMessage({
          message: '更新成功',
          type: 'success',
          duration: 2000,
          offset: 20
        })
      } else {
        await request.post('/codegen/config', editForm.value)
        ElMessage({
          message: '创建成功',
          type: 'success',
          duration: 2000,
          offset: 20
        })
      }
      showEditDialog.value = false
      loadList()
    } catch (err) {
      ElMessage.error('保存失败：' + (err.message || '未知错误'))
    }
  })
}

// 编辑配置
const editConfig = async (row) => {
  try {
    const res = await request.get(`/codegen/config/${row.id}`)
    // 确保columns字段存在
    editForm.value = {
      ...res.data,
      columns: res.data.columns || []
    }
    showEditDialog.value = true
  } catch (err) {
    ElMessage.error('加载配置失败')
  }
}

// 删除配置
const deleteConfig = (row) => {
  ElMessageBox.confirm('确定删除此配置吗？', '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.delete(`/codegen/config/${row.id}`)
      ElMessage({
        message: '删除成功',
        type: 'success',
        duration: 2000,
        offset: 20
      })
      loadList()
    } catch (err) {
      ElMessage.error('删除失败')
    }
  })
}

// 预览代码
const previewCode = async (row) => {
  try {
    const res = await request.post('/codegen/preview', { id: row.id })
    previewFiles.value = res.data.files || []
    if (previewFiles.value.length > 0) {
      activeTab.value = previewFiles.value[0].filePath
    }
    showPreviewDialog.value = true
  } catch (err) {
    ElMessage.error('预览失败')
  }
}

// 生成代码
const generateCode = async (row) => {
  try {
    const res = await request.post('/codegen/generate', { id: row.id })
    ElMessageBox.alert(
      `成功生成 ${res.data.files?.length || 0} 个文件`,
      '代码生成成功',
      {
        confirmButtonText: '确定',
        type: 'success'
      }
    )
  } catch (err) {
    ElMessage.error('生成失败')
  }
}

// 初始化
loadList()

// 组件卸载前清理
onBeforeUnmount(() => {
  // 关闭所有对话框
  showEditDialog.value = false
  showPreviewDialog.value = false
})
</script>

<style scoped>
.codegen-container {
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

.form-section {
  margin-bottom: 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.el-table .el-input__inner) {
  padding: 5px 10px;
}

:deep(.el-table .el-select) {
  width: 100%;
}
</style>
