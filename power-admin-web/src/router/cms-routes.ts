import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/cms',
    component: () => import('@/layout/MainLayout.vue'),
    children: [
      {
        path: 'content',
        name: 'CMSContent',
        component: () => import('@/pages/cms/ContentList.vue'),
        meta: { title: '内容管理' },
      },
      {
        path: 'category',
        name: 'CMSCategory',
        component: () => import('@/pages/cms/CategoryList.vue'),
        meta: { title: '分类管理' },
      },
      {
        path: 'tag',
        name: 'CMSTag',
        component: () => import('@/pages/cms/TagList.vue'),
        meta: { title: '标签管理' },
      },
      {
        path: 'comment',
        name: 'CMSComment',
        component: () => import('@/pages/cms/CommentList.vue'),
        meta: { title: '评论管理' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
