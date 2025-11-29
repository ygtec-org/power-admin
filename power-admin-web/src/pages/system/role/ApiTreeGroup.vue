<template>
  <div class="api-tree">
    <div v-for="(group, groupName) in groupedApis" :key="groupName">
      <div class="api-group">
        <div class="api-group-header">
          <span class="expand-icon" :class="{ expanded: expandedGroups[groupName] }" @click="toggleGroupExpand(groupName)">▶</span>
          <input
            type="checkbox"
            :id="'group-' + groupName"
            :checked="isGroupAllSelected(groupName)"
            :indeterminate="isGroupPartialSelected(groupName)"
            @change="handleGroupCheckChange(groupName, $event)"
          />
          <label :for="'group-' + groupName" class="group-label" @click="toggleGroupExpand(groupName)">{{ groupName || '未分组' }}</label>
        </div>
        <div v-if="expandedGroups[groupName]" class="api-group-items">
          <div v-for="api in group" :key="api.id" class="api-item">
            <input
              type="checkbox"
              :id="'api-' + api.id"
              :checked="isApiSelected(api.id)"
              @change="handleApiCheckChange(api.id, $event)"
            />
            <label :for="'api-' + api.id">
              <span class="api-name">{{ api.apiName }}</span>
              <span class="api-method" :class="'method-' + api.apiMethod.toLowerCase()">{{ api.apiMethod }}</span>
              <span class="api-path">{{ api.apiPath }}</span>
            </label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  apis: {
    type: Array,
    required: true
  },
  selectedApis: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['update-selection'])

const expandedGroups = ref({})
const groupCheckboxRefs = ref({})

const groupedApis = computed(() => {
  const grouped = {}
  props.apis.forEach(api => {
    const group = api.group || '未分组'
    if (!grouped[group]) {
      grouped[group] = []
      expandedGroups.value[group] = true
    }
    grouped[group].push(api)
  })
  return grouped
})

const toggleGroupExpand = (groupName) => {
  expandedGroups.value[groupName] = !expandedGroups.value[groupName]
}

const isApiSelected = (apiId) => {
  return props.selectedApis.includes(apiId)
}

// 检查一个API分组是否全部被选中
const isGroupAllSelected = (groupName) => {
  const group = groupedApis.value[groupName] || []
  return group.length > 0 && group.every(api => props.selectedApis.includes(api.id))
}

// 检查一个API分组是否部分被选中
const isGroupPartialSelected = (groupName) => {
  const group = groupedApis.value[groupName] || []
  if (group.length === 0) return false
  const selectedCount = group.filter(api => props.selectedApis.includes(api.id)).length
  return selectedCount > 0 && selectedCount < group.length
}

// 处理API分组的勾选变化
const handleGroupCheckChange = (groupName, event) => {
  const group = groupedApis.value[groupName] || []
  let newSelection = [...props.selectedApis]
  
  if (event.target.checked) {
    // 勾选：将整个API分组的所有API添加
    group.forEach(api => {
      if (!newSelection.includes(api.id)) {
        newSelection.push(api.id)
      }
    })
  } else {
    // 取消勾选：移除整个API分组的所有API
    const groupApiIds = group.map(api => api.id)
    newSelection = newSelection.filter(id => !groupApiIds.includes(id))
  }
  
  emit('update-selection', newSelection)
}

const handleApiCheckChange = (apiId, event) => {
  const newSelection = event.target.checked
    ? [...props.selectedApis, apiId]
    : props.selectedApis.filter(id => id !== apiId)
  emit('update-selection', newSelection)
}

// 监听selectedApis变化，更新checkbox的indeterminate状态
watch(() => props.selectedApis, () => {
  Object.keys(groupedApis.value).forEach(groupName => {
    const checkbox = document.getElementById('group-' + groupName)
    if (checkbox) {
      checkbox.indeterminate = isGroupPartialSelected(groupName)
    }
  })
}, { deep: true })
</script>

<style scoped>
.api-tree {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.api-group {
  border: 1px solid #e6e9f0;
  border-radius: 4px;
  overflow: hidden;
  background: white;
}

.api-group-header {
  display: flex;
  align-items: center;
  padding: 12px;
  background: #f5f7fa;
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
  gap: 8px;
}

.api-group-header:hover {
  background: #e6e9f0;
}

.expand-icon {
  display: inline-block;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 12px;
  transition: transform 0.2s;
  cursor: pointer;
  flex-shrink: 0;
}

.expand-icon.expanded {
  transform: rotate(90deg);
}

.group-name {
  font-size: 13px;
  font-weight: 500;
  color: #333;
}

.api-group-header input {
  width: auto;
  cursor: pointer;
  flex-shrink: 0;
}

.api-group-header .group-label {
  font-size: 13px;
  font-weight: 500;
  color: #333;
  margin: 0;
  cursor: pointer;
  flex: 1;
}

.api-group-items {
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.api-item {
  display: flex;
  align-items: flex-start;
  padding: 8px;
  border-radius: 3px;
  cursor: pointer;
  transition: background 0.2s;
}

.api-item:hover {
  background: #f5f7fa;
}

.api-item input {
  width: auto;
  margin-right: 8px;
  margin-top: 2px;
  cursor: pointer;
}

.api-item label {
  margin: 0;
  flex: 1;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.api-name {
  font-size: 13px;
  color: #333;
  font-weight: 500;
}

.api-method {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 2px;
  width: fit-content;
  color: white;
}

.api-method.method-get {
  background: #61affe;
}

.api-method.method-post {
  background: #49cc90;
}

.api-method.method-put {
  background: #fca130;
}

.api-method.method-delete {
  background: #f93e3e;
}

.api-path {
  font-size: 12px;
  color: #999;
  word-break: break-all;
}
</style>
