# 前端菜单权限集成指南

## 概述

前端需要调用菜单列表API获取当前用户可访问的菜单列表，并根据返回结果构建导航菜单。

## API接口

### 获取菜单列表

**端点：**
```
GET /api/admin/system/menus
```

**认证：**
- 必须提供有效的JWT Token
- 格式：`Authorization: Bearer <TOKEN>`

**查询参数：**
```typescript
interface MenuListRequest {
  page?: number;        // 分页页码（可选，默认1）
  pageSize?: number;    // 每页数量（可选，默认10）
  parentId?: number;    // 父菜单ID（可选，不传则返回所有）
}
```

**响应数据结构：**
```typescript
interface MenuListResponse {
  code: number;         // 状态码
  data: MenuItem[];     // 菜单列表
  total: number;        // 菜单总数
}

interface MenuItem {
  id: number;           // 菜单ID
  parentId: number;     // 父菜单ID（0为顶级菜单）
  menuName: string;     // 菜单名称
  menuPath: string;     // 菜单路由路径
  component: string;    // 组件路径
  icon: string;         // 菜单图标
  sort: number;         // 排序号（升序）
  status: number;       // 状态（1显示 0隐藏）
  menuType: number;     // 菜单类型（1菜单 2按钮）
  createdAt: string;    // 创建时间
  children?: MenuItem[]; // 子菜单列表
}
```

**示例请求：**
```typescript
// 使用 axios 或 fetch
const response = await axios.get('/api/admin/system/menus', {
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('token')}`
  }
});

console.log(response.data);
// 输出：
// {
//   "code": 200,
//   "data": [
//     {
//       "id": 1,
//       "parentId": 0,
//       "menuName": "系统管理",
//       "menuPath": "/system",
//       "component": "Layout",
//       "icon": "setting",
//       "sort": 1,
//       "status": 1,
//       "menuType": 1,
//       "createdAt": "2025-11-30 11:00:00",
//       "children": [...]
//     }
//   ],
//   "total": 5
// }
```

## 前端实现示例

### 1. Vue 3 + TypeScript 菜单获取

**创建服务文件：** `src/api/menu.ts`

```typescript
import axios from 'axios'

export interface MenuItem {
  id: number
  parentId: number
  menuName: string
  menuPath: string
  component: string
  icon: string
  sort: number
  status: number
  menuType: number
  createdAt: string
  children?: MenuItem[]
}

export interface MenuListResponse {
  code: number
  data: MenuItem[]
  total: number
}

/**
 * 获取菜单列表
 */
export async function getMenuList(): Promise<MenuItem[]> {
  try {
    const response = await axios.get<MenuListResponse>(
      '/api/admin/system/menus',
      {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      }
    )
    
    if (response.data.code === 200) {
      return response.data.data || []
    } else {
      throw new Error(`Failed to fetch menus: ${response.data.code}`)
    }
  } catch (error) {
    console.error('Error fetching menus:', error)
    return []
  }
}

/**
 * 构建路由数组
 */
export function buildRoutes(menus: MenuItem[]): any[] {
  return menus
    .filter(menu => menu.status === 1)
    .sort((a, b) => a.sort - b.sort)
    .map(menu => ({
      path: menu.menuPath,
      component: menu.component,
      name: menu.menuName,
      meta: {
        title: menu.menuName,
        icon: menu.icon,
        menuType: menu.menuType
      },
      children: menu.children && menu.children.length > 0
        ? buildRoutes(menu.children)
        : []
    }))
}
```

### 2. Vue Router 动态路由集成

**修改文件：** `src/router/index.ts`

```typescript
import { createRouter, createWebHistory, Router, RouteRecordRaw } from 'vue-router'
import { getMenuList, buildRoutes, MenuItem } from '@/api/menu'
import Layout from '@/layout/Layout.vue'

const baseRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    component: () => import('@/pages/Login.vue'),
    meta: { requiresAuth: false }
  }
]

let router: Router

export async function setupRouter() {
  // 获取菜单列表
  const menus = await getMenuList()
  
  if (menus.length === 0) {
    console.warn('No menus available for user')
    // 可以返回一个空仪表板或提示用户权限不足
    return
  }
  
  // 将菜单转换为路由
  const dynamicRoutes = buildRoutes(menus)
  
  // 为动态路由添加Layout包装
  const mainRoute: RouteRecordRaw = {
    path: '/',
    component: Layout,
    children: dynamicRoutes
  }
  
  // 初始化路由器
  router = createRouter({
    history: createWebHistory(),
    routes: [
      ...baseRoutes,
      mainRoute,
      {
        path: '/:pathMatch(.*)*',
        redirect: '/404'
      }
    ]
  })
  
  return router
}

export default router
```

### 3. 登录后初始化菜单

**修改文件：** `src/main.ts`

```typescript
import { createApp } from 'vue'
import App from './App.vue'
import { setupRouter } from '@/router'

async function bootstrap() {
  const app = createApp(App)
  
  // 检查是否已登录
  const token = localStorage.getItem('token')
  if (token) {
    // 初始化路由（包括获取菜单）
    const router = await setupRouter()
    if (router) {
      app.use(router)
    } else {
      console.error('Failed to setup router')
      // 清除token并重定向到登录
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
  } else {
    // 使用基础路由
    const { createRouter, createWebHistory } = await import('vue-router')
    const router = createRouter({
      history: createWebHistory(),
      routes: [
        {
          path: '/login',
          component: () => import('@/pages/Login.vue')
        }
      ]
    })
    app.use(router)
  }
  
  app.mount('#app')
}

bootstrap()
```

### 4. 菜单组件显示

**创建组件：** `src/components/Sidebar.vue`

```vue
<template>
  <aside class="sidebar">
    <div class="menu-list">
      <template v-for="menu in menus" :key="menu.id">
        <div v-if="menu.menuType === 1" class="menu-item">
          <router-link
            :to="menu.menuPath"
            class="menu-link"
            :class="{ active: isActive(menu.menuPath) }"
          >
            <i :class="menu.icon" class="menu-icon"></i>
            <span class="menu-name">{{ menu.menuName }}</span>
            <i v-if="menu.children?.length" class="expand-icon"></i>
          </router-link>
          
          <!-- 子菜单 -->
          <div v-if="menu.children?.length" class="submenu">
            <template v-for="child in menu.children" :key="child.id">
              <router-link
                :to="child.menuPath"
                class="submenu-link"
                :class="{ active: isActive(child.menuPath) }"
              >
                <i :class="child.icon" class="submenu-icon"></i>
                <span>{{ child.menuName }}</span>
              </router-link>
            </template>
          </div>
        </div>
        
        <!-- 按钮菜单（如果需要） -->
        <button
          v-else-if="menu.menuType === 2"
          class="menu-button"
          @click="handleButtonClick(menu)"
        >
          <i :class="menu.icon"></i>
          {{ menu.menuName }}
        </button>
      </template>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getMenuList, MenuItem } from '@/api/menu'

const route = useRoute()
const menus = ref<MenuItem[]>([])

onMounted(async () => {
  menus.value = await getMenuList()
})

const isActive = (path: string) => {
  return route.path.startsWith(path)
}

const handleButtonClick = (menu: MenuItem) => {
  console.log('Button clicked:', menu)
  // 可以发送事件或调用API
}
</script>

<style scoped>
.sidebar {
  width: 250px;
  background-color: #f5f5f5;
  border-right: 1px solid #e0e0e0;
  overflow-y: auto;
}

.menu-list {
  padding: 10px 0;
}

.menu-item {
  margin-bottom: 5px;
}

.menu-link,
.submenu-link {
  display: flex;
  align-items: center;
  padding: 10px 15px;
  color: #333;
  text-decoration: none;
  transition: all 0.3s;
}

.menu-link:hover,
.submenu-link:hover {
  background-color: #e0e0e0;
}

.menu-link.active,
.submenu-link.active {
  background-color: #1890ff;
  color: white;
}

.menu-icon {
  margin-right: 10px;
}

.submenu-icon {
  margin-right: 20px;
  margin-left: 20px;
}

.submenu {
  background-color: #fafafa;
}

.expand-icon {
  margin-left: auto;
}

.menu-button {
  width: 100%;
  padding: 10px 15px;
  border: none;
  background: none;
  cursor: pointer;
  text-align: left;
}
</style>
```

## 错误处理

### 常见错误及处理

```typescript
async function getMenuListWithErrorHandling(): Promise<MenuItem[]> {
  try {
    const response = await axios.get('/api/admin/system/menus', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    switch (response.data.code) {
      case 200:
        return response.data.data || []
      
      case 401:
        // Token无效或过期
        console.error('Token expired, please login again')
        localStorage.removeItem('token')
        window.location.href = '/login'
        return []
      
      case 403:
        // 权限不足
        console.error('Permission denied')
        return []
      
      default:
        console.error(`Unexpected response code: ${response.data.code}`)
        return []
    }
  } catch (error) {
    if (axios.isAxiosError(error)) {
      switch (error.response?.status) {
        case 401:
          console.error('Unauthorized')
          localStorage.removeItem('token')
          window.location.href = '/login'
          break
        
        case 403:
          console.error('Forbidden')
          break
        
        case 500:
          console.error('Server error')
          break
        
        default:
          console.error('Request failed:', error.message)
      }
    } else {
      console.error('Error:', error)
    }
    return []
  }
}
```

## 缓存策略

### 本地缓存菜单列表

```typescript
const MENU_CACHE_KEY = 'user_menus'
const MENU_CACHE_TTL = 1000 * 60 * 30 // 30分钟

export async function getMenuListWithCache(): Promise<MenuItem[]> {
  // 检查缓存
  const cached = sessionStorage.getItem(MENU_CACHE_KEY)
  if (cached) {
    try {
      const { menus, timestamp } = JSON.parse(cached)
      if (Date.now() - timestamp < MENU_CACHE_TTL) {
        console.log('Using cached menus')
        return menus
      }
    } catch (e) {
      console.error('Invalid cache data')
    }
  }
  
  // 获取菜单
  const menus = await getMenuList()
  
  // 保存到缓存
  sessionStorage.setItem(MENU_CACHE_KEY, JSON.stringify({
    menus,
    timestamp: Date.now()
  }))
  
  return menus
}

// 登出时清除缓存
export function clearMenuCache() {
  sessionStorage.removeItem(MENU_CACHE_KEY)
}
```

## 权限检查

### 前端权限守卫

```typescript
/**
 * 检查用户是否有权访问该页面
 */
function canAccessMenu(menuPath: string, menus: MenuItem[]): boolean {
  // 递归搜索菜单树
  function search(items: MenuItem[]): boolean {
    for (const item of items) {
      if (item.menuPath === menuPath && item.status === 1) {
        return true
      }
      if (item.children?.length) {
        if (search(item.children)) {
          return true
        }
      }
    }
    return false
  }
  
  return search(menus)
}

// 在路由导航守卫中使用
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (!token) {
    if (to.path === '/login') {
      next()
    } else {
      next('/login')
    }
  } else {
    if (to.path === '/login') {
      next('/')
    } else {
      // 检查权限
      if (canAccessMenu(to.path, menus.value)) {
        next()
      } else {
        console.warn(`No permission to access ${to.path}`)
        next('/403')
      }
    }
  }
})
```

## 性能优化

### 虚拟滚动（大菜单列表）

```vue
<template>
  <virtual-list
    :size="50"
    :remain="8"
    :bench="5"
  >
    <div v-for="menu in menus" :key="menu.id" class="menu-item">
      <!-- 菜单项内容 -->
    </div>
  </virtual-list>
</template>

<script setup>
import VirtualList from 'vue-virtual-scroller'
</script>
```

### 懒加载菜单组件

```typescript
const routes = menus.map(menu => ({
  path: menu.menuPath,
  component: () => import(`@/pages/${menu.component}.vue`),
  meta: {
    title: menu.menuName
  }
}))
```

## 最佳实践

1. **安全性**
   - 不要在前端完全依赖菜单权限，后端仍需检查API权限
   - 定期刷新Token和菜单列表

2. **性能**
   - 使用缓存减少API调用
   - 使用虚拟滚动处理大量菜单
   - 懒加载菜单对应的组件

3. **用户体验**
   - 提供加载状态反馈
   - 禁用菜单项而不是隐藏（保持布局稳定）
   - 面包屑导航显示当前位置

4. **可维护性**
   - 将菜单API逻辑独立到service层
   - 使用TypeScript确保类型安全
   - 编写单元测试

## 测试示例

```typescript
import { describe, it, expect, vi } from 'vitest'
import { getMenuList, buildRoutes } from '@/api/menu'

describe('Menu API', () => {
  it('should fetch menus successfully', async () => {
    const menus = await getMenuList()
    expect(Array.isArray(menus)).toBe(true)
  })
  
  it('should build routes from menus', () => {
    const menus = [
      {
        id: 1,
        parentId: 0,
        menuName: 'System',
        menuPath: '/system',
        component: 'Layout',
        icon: 'setting',
        sort: 1,
        status: 1,
        menuType: 1,
        createdAt: '2025-11-30'
      }
    ]
    
    const routes = buildRoutes(menus)
    expect(routes[0].path).toBe('/system')
    expect(routes[0].name).toBe('System')
  })
})
```
