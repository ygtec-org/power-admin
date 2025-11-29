<template>
  <div class="page">
    <div class="page-header">
      <h1>字典管理</h1>
      <button @click="showDialog=true" class="btn-primary">+ 新增字典</button>
    </div>
    <!-- 搜索和筛选区域 -->
    <div class="filter-bar">
      <div class="filter-item">
        <input 
          v-model="searchKey" 
          placeholder="请输入字典名称" 
          class="search-input"
          @keyup.enter="handleSearch"
        />
        <button @click="handleSearch" class="btn-primary search-btn">搜索</button>
      </div>
    </div>
    <div class="table-box">
      <table class="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>key</th>
            <th>value</th>
            <th>类型</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="!dicts.length">
            <td colspan="5" class="empty">暂无数据</td>
          </tr>
          <tr v-for="d in dicts" :key="d.id">
            <td>{{d.id}}</td>
            <td><code>{{d.dictLabel}}</code></td>
            <td>{{d.dictValue}}</td>
            <td>{{d.dictType||'-'}}</td>
            <td>
              <button @click="edit(d)" class="btn-sm">编辑</button>
              <button @click="del(d)" class="btn-sm danger">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <!-- 使用通用分页组件 -->
    <Pagination 
      v-model:currentPage="currentPage"
      v-model:pageSize="pageSize"
      :total="total"
      @page-change="handlePageChange"
      @size-change="handleSizeChange"
    />
  </div>
  <div v-if="showDialog" class="modal">
    <div class="modal-content">
      <div class="modal-header">
        <h2>{{isEdit?'编辑字典':'新增字典'}}</h2>
        <button @click="showDialog=false" class="close-btn">×</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label>Key</label>
          <input v-model="form.key" />
        </div>
        <div class="form-group">
          <label>Value</label>
          <input v-model="form.value" />
        </div>
        <div class="form-group">
          <label>类型</label>
          <input v-model="form.type" placeholder="如: config, status" />
        </div>
      </div>
      <div class="modal-footer">
        <button @click="showDialog=false" class="btn-cancel">取消</button>
        <button @click="save" class="btn-primary">保存</button>
      </div>
    </div>
  </div>
  <ConfirmDialog 
    v-model:visible="showConfirm"
    title="删除确认"
    :message="`确定要删除「${deleteTarget?.key}」吗？`"
    @confirm="handleConfirmDelete"
    @cancel="handleCancelDelete"
  />
</template>

<script setup>
import { ref } from 'vue'
import { getDicts, createDict, updateDict, deleteDict } from '../../../api/menu'
import { notify } from '../../../components/Notification'
import ConfirmDialog from '../../../components/ConfirmDialog.vue'
import Pagination from '../../../components/Pagination.vue'

const dicts = ref([])
const showDialog = ref(false)
const isEdit = ref(false)
const form = ref({key:'',value:'',type:''})
const selected = ref(null)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 搜索相关
const searchKey = ref('')

// 确认对话框相关
const showConfirm = ref(false)
const deleteTarget = ref(null)

const load = async (page = currentPage.value, size = pageSize.value) => {
  try {
    const params = {
      page: page,
      pageSize: size
    }
    
    // 如果有搜索关键词，则添加到参数中
    if (searchKey.value.trim()) {
      params.dictLabel = searchKey.value.trim()
    }
    
    const res = await getDicts(params)
    dicts.value = res.data.list || []
    total.value = res.data.total || 0
    currentPage.value = page
    pageSize.value = size
  } catch (e) {
    console.error('获取字典列表失败:', e)
  }
}

const handlePageChange = (page) => {
  load(page, pageSize.value)
}

const handleSizeChange = (size) => {
  pageSize.value = size
  load(1, size) // 页码重置为第一页
}

const handleSearch = () => {
  currentPage.value = 1
  load(1, pageSize.value)
}

const save = async () => {
  if(!form.value.key||!form.value.value) {notify.warning('请填写必填项');return}
  try {
    if(isEdit.value) {
      await updateDict(selected.value.id, form.value)
      notify.success('编辑成功')
    } else {
      await createDict(form.value)
      notify.success('创建成功')
    }
    showDialog.value = false
    form.value = {key:'',value:'',type:''}
    load()
  } catch (e) {
    notify.error(e.message || '操作失败')
  }
}

const edit = (d) => {
  isEdit.value = true
  selected.value = d
  form.value = {key:d.dictLabel,value:d.dictValue,type:d.dictType}
  showDialog.value = true
}

const del = async (d) => {
  deleteTarget.value = d
  showConfirm.value = true
}

const handleConfirmDelete = async () => {
  try { 
    await deleteDict(deleteTarget.value.id)
    notify.success('删除成功')
    load() 
  } catch(e) {
    notify.error(e.message || '删除失败')
  } finally {
    deleteTarget.value = null
  }
}

const handleCancelDelete = () => {
  deleteTarget.value = null
}

load()
</script>

<style scoped>
.page { animation: fadeIn 0.3s; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h1 { margin: 0; font-size: 24px; }
.btn-primary { padding: 8px 16px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; border: none; border-radius: 4px; cursor: pointer; }
.table-box { background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,.05); }
.table { width: 100%; border-collapse: collapse; font-size: 14px; }
tr { border-bottom: 1px solid #e6e9f0; }
th { padding: 12px; text-align: left; font-weight: 600; color: #666; background: #f5f7fa; }
td { padding: 12px; }
code { background: #f5f7fa; padding: 2px 6px; border-radius: 3px; font-size: 12px; }
.empty { text-align: center; color: #999; padding: 40px !important; }
.btn-sm { padding: 4px 12px; border: 1px solid #ddd; background: white; border-radius: 3px; cursor: pointer; font-size: 12px; margin-right: 8px; }
.btn-sm:hover { border-color: #667eea; color: #667eea; }
.btn-sm.danger:hover { border-color: #f56c6c; color: #f56c6c; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-content { background: white; border-radius: 8px; width: 90%; max-width: 500px; }
.modal-header { display: flex; justify-content: space-between; padding: 20px; border-bottom: 1px solid #e6e9f0; }
.modal-header h2 { margin: 0; }
.close-btn { background: none; border: none; font-size: 24px; cursor: pointer; }
.modal-body { padding: 20px; }
.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 8px; font-weight: 500; }
.form-group input { width: 100%; padding: 8px 12px; border: 1px solid #ddd; border-radius: 4px; }
.modal-footer { display: flex; justify-content: flex-end; gap: 12px; padding: 20px; border-top: 1px solid #e6e9f0; }
.btn-cancel { padding: 8px 16px; border: 1px solid #ddd; background: white; border-radius: 4px; cursor: pointer; }
@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }

/* 新增样式 */
.filter-bar { 
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  margin-bottom: 20px; 
  background: white; 
  padding: 16px; 
  border-radius: 8px; 
  box-shadow: 0 2px 4px rgba(0,0,0,.05); 
}
.filter-item { display: flex; gap: 12px; }
.search-input { 
  padding: 8px 12px; 
  border: 1px solid #ddd; 
  border-radius: 4px; 
  width: 200px; 
}
.search-btn { padding: 8px 16px; }


</style>