<template>
  <div class="layout">
    <header class="header">
      <div class="header-left">
        <h1>Power Admin</h1>
      </div>
      <div class="header-right">
        <span class="user-name">{{ userName }}</span>
        <button @click="logout" class="logout-btn">退出登录</button>
      </div>
    </header>

    <div class="main">
      <aside class="sidebar">
        <nav class="menu">
          <RouterLink to="/dashboard" class="menu-item">
            <span>📊 仪表板</span>
          </RouterLink>

          <!-- 动态可折叠树形菜单 -->
          <template v-for="menu in menus" :key="menu.id">
            <div class="menu-group" v-if="menu.children && menu.children.length > 0">
              <div class="menu-parent" @click="toggleMenu(menu.id)">
                <span class="menu-icon">{{ expandedMenus.has(menu.id) ? '▼' : '▶' }}</span>
                <span>{{ getMenuIcon(menu.icon) }} {{ menu.menuName }}</span>
              </div>
              <transition name="slide">
                <div v-show="expandedMenus.has(menu.id)" class="menu-children">
                  <template v-for="child in menu.children" :key="child.id">
                    <RouterLink :to="child.menuPath" class="menu-item menu-child">
                      <span>{{ getMenuIcon(child.icon) }} {{ child.menuName }}</span>
                    </RouterLink>
                  </template>
                </div>
              </transition>
            </div>
            <RouterLink v-else :to="menu.menuPath" class="menu-item">
              <span>{{ getMenuIcon(menu.icon) }} {{ menu.menuName }}</span>
            </RouterLink>
          </template>
        </nav>
      </aside>

      <main class="content">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { RouterLink, RouterView } from 'vue-router'
import { getMenuTree } from '../api/menu'

const router = useRouter()
const userName = ref('Admin')
const menus = ref([])
const expandedMenus = ref(new Set()) // 记录展开的菜单ID

// 切换菜单展开/收起
const toggleMenu = (menuId) => {
  if (expandedMenus.value.has(menuId)) {
    expandedMenus.value.delete(menuId)
  } else {
    expandedMenus.value.add(menuId)
  }
  // 触发响应式更新
  expandedMenus.value = new Set(expandedMenus.value)
}

// 获取菜单数据
const loadMenus = async () => {
  try {
    const res = await getMenuTree()
    // 后端返回格式: { data: { total, data: [...] } }
    // 响应拦截器提取后: { data: { total, data: [...] } }
    menus.value = res.data.list || res.data || []
    console.log('加载的菜单数据:', menus.value)
  } catch (error) {
    console.error('获取菜单失败:', error)
  }
}

// 根据图标名称返回 emoji
const getMenuIcon = (icon: string) => {
  const iconMap: Record<string, string> = {
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
    'monitor': '📊'
  }
  return iconMap[icon] || '📌'
}

onMounted(() => {
  const user = localStorage.getItem('user')
  if (user) {
    const userData = JSON.parse(user)
    userName.value = userData.nickname || 'Admin'
  }
  loadMenus()
})

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}
</script>

<style scoped>
.layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: #f5f7fa;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  height: 64px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-left h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-name {
  font-size: 14px;
}

.logout-btn {
  padding: 6px 16px;
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.4);
  color: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.main {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 256px;
  background: white;
  border-right: 1px solid #e6e9f0;
  overflow-y: auto;
}

.menu {
  padding: 16px 0;
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
}

.menu-child {
  padding-left: 48px !important;
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

.menu-title {
  margin: 8px 0;
  padding: 0 24px;
  font-size: 12px;
  color: #999;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: #d9d9d9;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #bbb;
}
</style>
