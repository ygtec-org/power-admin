import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Login from '../pages/Login.vue'
import Layout from '../layout/Layout.vue'
import Dashboard from '../pages/Dashboard.vue'

// 系统管理
import UserList from '../pages/system/user/UserList.vue'
import RoleList from '../pages/system/role/RoleList.vue'
import MenuList from '../pages/system/menu/MenuList.vue'
import PermissionList from '../pages/system/permission/PermissionList.vue'
import ApiList from '../pages/system/api/ApiList.vue'

// 内容管理
import DictList from '../pages/content/dict/DictList.vue'

// CMS内容系统
import CMSContentList from '../pages/cms/content/ContentList.vue'
import CMSCategoryList from '../pages/cms/category/CategoryList.vue'
import CMSTagList from '../pages/cms/tag/TagList.vue'
import CMSCommentList from '../pages/cms/comment/CommentList.vue'
import CMSUserList from '../pages/cms/user/UserList.vue'
import CMSPublishList from '../pages/cms/publish/PublishList.vue'

// 应用市场
import AppMarket from '../pages/market/AppMarket.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: Dashboard,
        meta: { title: '仪表板' },
      },
      // 系统管理
      {
        path: '/system/users',
        name: 'UserList',
        component: UserList,
        meta: { title: '用户管理' },
      },
      {
        path: '/system/roles',
        name: 'RoleList',
        component: RoleList,
        meta: { title: '角色管理' },
      },
      {
        path: '/system/menus',
        name: 'MenuList',
        component: MenuList,
        meta: { title: '菜单管理' },
      },
      {
        path: '/system/permissions',
        name: 'PermissionList',
        component: PermissionList,
        meta: { title: '权限管理' },
      },
      {
        path: '/system/apis',
        name: 'ApiList',
        component: ApiList,
        meta: { title: 'API管理' },
      },
      // 内容管理
      {
        path: '/content/dicts',
        name: 'DictList',
        component: DictList,
        meta: { title: '字典管理' },
      },
      // CMS内容系统
      {
        path: '/cms/content',
        name: 'CMSContentList',
        component: CMSContentList,
        meta: { title: '内容管理' },
      },
      {
        path: '/cms/category',
        name: 'CMSCategoryList',
        component: CMSCategoryList,
        meta: { title: '分类管理' },
      },
      {
        path: '/cms/tag',
        name: 'CMSTagList',
        component: CMSTagList,
        meta: { title: '标签管理' },
      },
      {
        path: '/cms/comment',
        name: 'CMSCommentList',
        component: CMSCommentList,
        meta: { title: '评论管理' },
      },
      {
        path: '/cms/user',
        name: 'CMSUserList',
        component: CMSUserList,
        meta: { title: '用户管理' },
      },
      {
        path: '/cms/publish',
        name: 'CMSPublishList',
        component: CMSPublishList,
        meta: { title: '发布管理' },
      },
      // 应用市场
      {
        path: '/market/apps',
        name: 'AppMarket',
        component: AppMarket,
        meta: { title: '应用市场' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/login',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const requiresAuth = to.meta.requiresAuth !== false

  if (requiresAuth && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
