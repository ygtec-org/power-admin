import request from './request'

export interface Permission {
  id: number
  name: string
  resource: string
  action: string
  description?: string
  status: number
  created_at?: string
}

export interface PermissionListResponse {
  list: Permission[]
  total: number
}

export interface CreatePermissionRequest {
  name: string
  resource: string
  action: string
  description?: string
  status?: number
}

export interface UpdatePermissionRequest {
  name?: string
  resource?: string
  action?: string
  description?: string
  status?: number
}

// 获取权限列表
export const getPermissions = (params?: { page?: number; pageSize?: number }) => {
  return request.get<PermissionListResponse>('/system/permissions', { params })
}

// 获取单个权限
export const getPermission = (id: number) => {
  return request.get<Permission>(`/system/permissions/${id}`)
}

// 创建权限
export const createPermission = (data: CreatePermissionRequest) => {
  return request.post<Permission>('/system/permissions', data)
}

// 更新权限
export const updatePermission = (id: number, data: UpdatePermissionRequest) => {
  return request.put<Permission>('/system/permissions', { id, ...data })
}

// 删除权限
export const deletePermission = (id: number) => {
  return request.delete('/system/permissions', { data: { id } })
}
