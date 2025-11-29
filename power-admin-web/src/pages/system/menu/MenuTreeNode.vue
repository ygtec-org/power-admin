<template>
  <tr :class="['menu-row', 'level-'+level]">
    <td>
      <span v-for="i in level" :key="i" class="tree-indent"></span>
      <span v-if="hasChildren" class="expander" @click="expanded = !expanded">{{ expanded ? '▼' : '▶' }} &emsp;</span>
      <span v-else class="tree-indent"></span>
      <span>{{ menu.id }}</span>
    </td>
    <td class="menu-name">
      <span class="icon">{{ getMenuIcon(menu.icon) }}</span>
      <span class="name-text">{{ menu.menuName }}</span>
    </td>
    <td><code>{{ menu.menuPath }}</code></td>
    <td><code>{{ menu.component || '-' }}</code></td>
    <td>{{ menu.icon || '-' }}</td>
    <td class="center">{{ menu.sort || 0 }}</td>
    <td>
      <span class="badge" :class="menu.status === 1 ? 'success' : 'danger'">
        {{ menu.status === 1 ? '显示' : '隐藏' }}
      </span>
    </td>
    <td>
      <button @click="handleAddChild" class="btn-sm">新增子菜单</button>
      <button @click="$emit('edit', menu)" class="btn-sm">编辑</button>
      <button @click="$emit('delete', menu)" class="btn-sm danger">删除</button>
    </td>
  </tr>
  <template v-if="expanded && hasChildren">
    <MenuTreeNode
      v-for="child in menu.children"
      :key="child.id"
      :menu="child"
      :level="level + 1"
      @edit="$emit('edit', $event)"
      @delete="$emit('delete', $event)"
      @add-child="$emit('add-child', $event)"
    />
  </template>

</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  menu: {
    type: Object,
    required: true
  },
  level: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['edit', 'delete', 'add-child'])

const expanded = ref(false)
const hasChildren = computed(() => props.menu.children && props.menu.children.length > 0)

// 图标映射与显示
const iconMap = {
  setting: '⚙️',
  user: '👤',
  admin: '🎯',
  menu: '📋',
  lock: '🔐',
  link: '🔗',
  document: '📄',
  list: '📚',
  shopping: '🛍️',
  shop: '🛑',
  monitor: '📊'
}
const getMenuIcon = (icon) => iconMap[icon] || '📌'

const handleAddChild = () => {
  emit('add-child', props.menu.id)
}
</script>

<style scoped>
.tree-indent {
  display: inline-block;
  width: 20px;
}

.menu-row { transition: background 0.2s; }
.menu-row:hover { background: #fafbff; }

.menu-row td { padding: 12px 10px; }

.menu-row:not(:last-child) td { border-bottom: 1px solid #f0f0f0; }

.menu-name {  align-items: center; gap: 8px; }
.expander { cursor: pointer; color: #667eea; font-size: 12px; }
.icon { width: 18px; text-align: center; }
.name-text { font-weight: 600; color: #333; }
.subtext { margin-left: 38px; font-size: 12px; color: #999; }
.center { text-align: center; }

/* 垂直对齐所有单元格 */
td { vertical-align: middle; }

.btn-sm { padding: 4px 12px; border: 1px solid #ddd; background: white; border-radius: 3px; cursor: pointer; font-size: 12px; margin-right: 4px; }
.btn-sm:hover { border-color: #667eea; color: #667eea; }
.btn-sm.danger:hover { border-color: #f56c6c; color: #f56c6c; }

.badge { display: inline-block; padding: 2px 8px; border-radius: 12px; font-size: 12px; }
.badge.success { background: #e9f6ea; color: #2f8a3b; }
.badge.danger { background: #fdecec; color: #d23a3a; }

code { background: #f5f7fa; padding: 2px 6px; border-radius: 3px; font-size: 12px; color: #d63384; }
</style>
