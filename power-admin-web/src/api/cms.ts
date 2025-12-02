import { cmsRequest as request } from '@/api/request'

// ======================== 内容管理 ========================

/**
 * 获取内容列表
 */
export const getContentList = (params: {
  page: number
  pageSize: number
  categoryId?: number
  status?: number
  search?: string
}) => {
  return request.get('/content/list', { params })
}

/**
 * 获取内容详情
 */
export const getContentDetail = (id: number) => {
  return request.get(`/content/${id}`)
}

/**
 * 创建内容
 */
export const createContent = (data: {
  title: string
  description: string
  content: string
  slug: string
  categoryId: number
  status?: number
}) => {
  return request.post('/content', data)
}

/**
 * 更新内容
 */
export const updateContent = (id: number, data: any) => {
  return request.put(`/content/${id}`, data)
}

/**
 * 删除内容
 */
export const deleteContent = (id: number) => {
  return request.delete(`/content/${id}`)
}

/**
 * 发布内容
 */
export const publishContent = (id: number) => {
  return request.post(`/content/${id}/publish`, {})
}

/**
 * 取消发布
 */
export const unpublishContent = (id: number) => {
  return request.post(`/content/${id}/unpublish`, {})
}

// ======================== 分类管理 ========================

/**
 * 获取分类树
 */
export const getCategoryTree = () => {
  return request.get('/category/tree')
}

/**
 * 获取分类列表
 */
export const getCategoryList = (params: { parentId?: number } = {}) => {
  return request.get('/category/list', { params })
}

/**
 * 创建分类
 */
export const createCategory = (data: {
  name: string
  slug: string
  description: string
  parentId?: number
}) => {
  return request.post('/category', data)
}

/**
 * 更新分类
 */
export const updateCategory = (id: number, data: any) => {
  return request.put(`/category/${id}`, data)
}

/**
 * 删除分类
 */
export const deleteCategory = (id: number) => {
  return request.delete(`/category/${id}`)
}

// ======================== 标签管理 ========================

/**
 * 获取标签列表
 */
export const getTagList = (params: { page?: number; pageSize?: number } = {}) => {
  return request.get('/tag/list', { params })
}

/**
 * 创建标签
 */
export const createTag = (data: {
  name: string
  slug: string
  description: string
  color?: string
}) => {
  return request.post('/tag', data)
}

/**
 * 更新标签
 */
export const updateTag = (id: number, data: any) => {
  return request.put(`/tag/${id}`, data)
}

/**
 * 删除标签
 */
export const deleteTag = (id: number) => {
  return request.delete(`/tag/${id}`)
}

// ======================== 评论管理 ========================

/**
 * 获取评论列表
 */
export const getCommentList = (params: {
  contentId?: number
  page?: number
  pageSize?: number
}) => {
  return request.get('/comment/list', { params })
}

/**
 * 审核评论
 */
export const approveComment = (id: number) => {
  return request.post(`/comment/${id}/approve`, {})
}

/**
 * 拒绝评论
 */
export const rejectComment = (id: number) => {
  return request.post(`/comment/${id}/reject`, {})
}

/**
 * 删除评论
 */
export const deleteComment = (id: number) => {
  return request.delete(`/comment/${id}`)
}

// ======================== 用户管理 ========================

/**
 * 获取用户列表
 */
export const getUserList = () => {
  return request.get('/user/list')
}

/**
 * 获取用户详情
 */
export const getUserDetail = (id: number) => {
  return request.get(`/user/${id}`)
}

/**
 * 禁用用户
 */
export const disableUser = (id: number) => {
  return request.post(`/user/${id}/disable`, {})
}

/**
 * 启用用户
 */
export const enableUser = (id: number) => {
  return request.post(`/user/${id}/enable`, {})
}

// ======================== 发布管理 ========================

/**
 * 立即发布
 */
export const publishImmediate = (data: { id: number }) => {
  return request.post('/publish/immediate', data)
}

/**
 * 定时发布
 */
export const publishSchedule = (data: { id: number; scheduleAt: string }) => {
  return request.post('/publish/schedule', data)
}

/**
 * 取消发布
 */
export const publishCancel = (id: number) => {
  return request.post(`/publish/${id}/cancel`, {})
}

/**
 * 批量发布
 */
export const publishBatch = (data: { ids: number[] }) => {
  return request.post('/publish/batch', data)
}
