<template>
  <div class="tree-item">
    <div class="tree-node" @click.stop="toggleExpand">
      <span 
        v-if="hasChildren" 
        class="expand-icon" 
        :class="{ expanded: isExpanded }"
      >▶</span>
      <span v-else class="expand-icon placeholder"></span>
      <input
        type="checkbox"
        :id="'menu-' + menu.id"
        :checked="isSelected"
        @change="handleCheckChange"
      />
      <label :for="'menu-' + menu.id">{{ menu.menuName }}</label>
    </div>
    <div v-if="isExpanded && hasChildren" class="tree-children">
      <MenuTreeNode 
        v-for="child in menu.children" 
        :key="child.id"
        :menu="child" 
        :selected-menus="selectedMenus" 
        :expanded-menus="expandedMenus"
        @toggle-expand="$emit('toggle-expand', $event)"
        @update-selection="$emit('update-selection', $event)"
      />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  menu: {
    type: Object,
    required: true
  },
  selectedMenus: {
    type: Array,
    required: true
  },
  expandedMenus: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['toggle-expand', 'update-selection'])

const hasChildren = computed(() => props.menu.children && props.menu.children.length > 0)

const isExpanded = computed(() => props.expandedMenus[props.menu.id] || false)

const isSelected = computed(() => props.selectedMenus.includes(props.menu.id))

const toggleExpand = () => {
  emit('toggle-expand', props.menu.id)
}

const handleCheckChange = (event) => {
  // 递归获取所有子菜单ID
  const getAllChildIds = (node) => {
    let ids = [node.id]
    if (node.children && node.children.length > 0) {
      node.children.forEach(child => {
        ids = ids.concat(getAllChildIds(child))
      })
    }
    return ids
  }
  
  let newSelection = [...props.selectedMenus]
  
  if (event.target.checked) {
    // 勾选：添加自己和所有子菜单
    const allIds = getAllChildIds(props.menu)
    allIds.forEach(id => {
      if (!newSelection.includes(id)) {
        newSelection.push(id)
      }
    })
  } else {
    // 取消勾选：移除自己和所有子菜单
    const allIds = getAllChildIds(props.menu)
    newSelection = newSelection.filter(id => !allIds.includes(id))
  }
  
  emit('update-selection', newSelection)
}
</script>

<style scoped>
.tree-item {
  margin: 0;
}

.tree-node {
  display: flex;
  align-items: center;
  padding: 8px 0;
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
}

.tree-node:hover {
  background: #f5f7fa;
}

.expand-icon {
  display: inline-block;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 12px;
  transition: transform 0.2s;
  cursor: pointer;
}

.expand-icon.expanded {
  transform: rotate(90deg);
}

.expand-icon.placeholder {
  cursor: default;
}

.tree-node input {
  width: auto;
  margin: 0 8px 0 4px;
  cursor: pointer;
}

.tree-node label {
  margin: 0;
  flex: 1;
  cursor: pointer;
  font-size: 13px;
  color: #333;
}

.tree-children {
  padding-left: 24px;
  display: flex;
  flex-direction: column;
  gap: 0;
  background: #f9f9f9;
  margin-left: 12px;
  border-left: 2px solid #e6e9f0;
  padding-left: 12px;
  margin-left: 12px;
}
</style>
