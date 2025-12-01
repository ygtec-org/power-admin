<template>
  <template v-for="menu in menus" :key="menu.id">
    <!-- 有子菜单的父级菜单 -->
    <div v-if="menu.children && menu.children.length > 0" class="menu-group">
      <div class="menu-parent" @click="$emit('toggle', menu.id)">
        <span class="menu-icon">{{ expandedMenus.has(menu.id) ? '▼' : '▶' }}</span>
        <span>{{ getMenuIcon(menu.icon) }} {{ menu.menu_name || menu.menuName }}</span>
      </div>
      <transition name="slide">
        <div v-show="expandedMenus.has(menu.id)" class="menu-children">
          <!-- 递归渲染子菜单 -->
          <MenuTree 
            :menus="menu.children" 
            :expanded-menus="expandedMenus"
            :level="level + 1"
            @toggle="$emit('toggle', $event)"
          />
        </div>
      </transition>
    </div>
    <!-- 没有子菜单的叶子菜单（允许 menu_path 或 menuPath） -->
    <RouterLink 
      v-else 
      :to="menu.menu_path || menu.menuPath || '#'" 
      class="menu-item"
      :class="{ 'menu-child': level > 0 }"
      :style="{ paddingLeft: `${24 + level * 24}px` }"
    >
      <span>{{ getMenuIcon(menu.icon) }} {{ menu.menu_name || menu.menuName }}</span>
    </RouterLink>
  </template>
</template>

<script setup>
import { RouterLink } from 'vue-router'

const props = defineProps({
  menus: {
    type: Array,
    default: () => []
  },
  expandedMenus: {
    type: Set,
    default: () => new Set()
  },
  level: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['toggle'])

// 根据图标名称返回 emoji
const getMenuIcon = (icon) => {
  const iconMap = {
    'setting': '⚙️',
    'user': '👤',
    'admin': '🎯',
    'menu': '📋',
    'lock': '🔐',
    'link': '🔗',
    'document': '📄',
    'list': '📚',
    'shopping': '🛍️',
    'shop': '🛑',
    'monitor': '📊',
    'content': '📰'
  }
  return iconMap[icon] || '📌'
}
</script>

<style scoped>
.menu-group {
  margin-top: 8px;
}

.menu-parent {
  display: flex;
  align-items: center;
  padding: 12px 24px;
  color: #333;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  user-select: none;
}

.menu-parent:hover {
  background: #f5f7fa;
  color: #667eea;
}

.menu-icon {
  display: inline-block;
  width: 16px;
  font-size: 10px;
  margin-right: 8px;
  transition: transform 0.3s;
}

.menu-children {
  overflow: hidden;
  background-color: #fafafa;
}

.menu-item {
  display: block;
  padding: 12px 24px;
  color: #666;
  text-decoration: none;
  transition: all 0.3s;
  border-left: 3px solid transparent;
}

.menu-item:hover {
  background: #f5f7fa;
  color: #667eea;
}

.menu-item.router-link-active {
  background: #f0f4ff;
  color: #667eea;
  border-left-color: #667eea;
  font-weight: 500;
}

.menu-child {
  font-size: 13px;
}

/* 展开/收起动画 */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
  max-height: 500px;
}

.slide-enter-from,
.slide-leave-to {
  max-height: 0;
  opacity: 0;
}
</style>
