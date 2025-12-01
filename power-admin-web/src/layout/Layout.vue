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

          <!-- 递归树形菜单 -->
          <MenuTree :menus="menus.length > 0 ? menus : defaultMenus" :expanded-menus="expandedMenus" @toggle="toggleMenu" />
        </nav>
      </aside>

      <main class="content">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { RouterLink, RouterView } from 'vue-router'
import { getMenuTree } from '../api/menu'
import MenuTree from './MenuTree.vue'

const router = useRouter()
const userName = ref('Admin')
const menus = ref([])
const defaultMenus = ref([])
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

// 递归构建树形菜单（兼容 parent_id 和 parentId 两种字段名）
const buildMenuTree = (menuList, parentId = 0) => {
  const result = []
  menuList.forEach(menu => {
    const pId = menu.parent_id !== undefined ? menu.parent_id : menu.parentId
    if (pId === parentId) {
      const children = buildMenuTree(menuList, menu.id)
      const menuItem = {
        ...menu,
        parent_id: pId,
        parentId: pId,
        children: children.length > 0 ? children : []
      }
      result.push(menuItem)
    }
  })
  return result
}

// 获取菜单数据
const loadMenus = async () => {
  try {
    const res = await getMenuTree()
    console.log('原始菜单响应:', res)
    
    // 处理多种可能的响应格式
    let menuList = []
    if (res.data) {
      if (Array.isArray(res.data)) {
        menuList = res.data
      } else if (res.data.list) {
        menuList = res.data.list
      } else if (res.data.data) {
        menuList = res.data.data
      }
    } else if (Array.isArray(res)) {
      menuList = res
    }
    
    // 过滤出真正的菜单项（排除应用市场等可能混入的非菜单数据）
    menuList = menuList.filter(item => 
      item.menu_name || item.menuName || item.menu_path || item.menuPath
    )
    
    console.log('过滤后的菜单列表:', menuList)
    
    if (menuList.length === 0) {
      loadDefaultMenus()
      return
    }
    
    // 直接使用 API 返回的菜单（后端已经构建好树形结构）
    menus.value = menuList
    console.log('最终菜单树:', menus.value)
  } catch (error) {
    console.error('获取菜单失败:', error)
    // 加载失败时使用默认菜单
    loadDefaultMenus()
  }
}

// 默认菜单（当API失败时使用）
const loadDefaultMenus = () => {
  const defaultMenuList = [
    {
      id: 1,
      menu_name: '系统管理',
      menuName: '系统管理',
      menu_path: '/system',
      menuPath: '/system',
      icon: 'admin',
      parent_id: 0,
    },
    { id: 2, menu_name: '用户管理', menuName: '用户管理', menu_path: '/system/users', menuPath: '/system/users', icon: 'user', parent_id: 1 },
    { id: 3, menu_name: '角色管理', menuName: '角色管理', menu_path: '/system/roles', menuPath: '/system/roles', icon: 'admin', parent_id: 1 },
    {
      id: 9,
      menu_name: '应用中心',
      menuName: '应用中心',
      menu_path: '/market',
      menuPath: '/market',
      icon: 'shopping',
      parent_id: 0,
    },
    { id: 10, menu_name: '应用市场', menuName: '应用市场', menu_path: '/market/apps', menuPath: '/market/apps', icon: 'shop', parent_id: 9 }
  ]
  // 使用树形构建函数统一处理
  const treeMenus = buildMenuTree(defaultMenuList)
  defaultMenus.value = treeMenus
  menus.value = treeMenus
  console.log('使用默认菜单，树形结构:', menus.value)
}

onMounted(() => {
  // 先加载默认菜单，确保菜单一定会显示
  loadDefaultMenus()
  
  const user = localStorage.getItem('user')
  if (user) {
    const userData = JSON.parse(user)
    userName.value = userData.nickname || 'Admin'
  }
  // 然后尝试从 API 加载菜单
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
