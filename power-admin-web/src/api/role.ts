import request from './request'

export interface Role {
  id: number
  name: string
  description: string
  status: number
  created_at?: string
  updated_at?: string
}

export interface RoleListResponse {
  list: Role[]
  total: number
}

export interface Permission {
  id: number
  name: string
  resource: string
  action: string
  description?: string
  status: number
}

export interface CreateRoleRequest {
  name: string
  description?: string
  status?: number
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
  status?: number
}

export interface AssignPermissionsRequest {
  permissionIds: number[]
}

// 获取角色列表
export const getRoles = (params?: { page?: number; pageSize?: number }) => {
  return request.get<RoleListResponse>('/system/roles', { params })
}

// 获取单个角色
export const getRole = (id: number) => {
  return request.get<Role>(`/system/roles/${id}`)
}

// 创建角色
export const createRole = (data: CreateRoleRequest) => {
  return request.post<Role>('/system/roles', data)
}

// 更新角色
export const updateRole = (id: number, data: UpdateRoleRequest) => {
  return request.put<Role>('/system/roles', { id, ...data })
}

// 删除角色
export const deleteRole = (id: number) => {
  return request.delete('/system/roles', { data: { id } })
}

// 为角色分配权限
export const assignPermissions = (roleId: number, permissionIds: number[]) => {
  return request.post(`/system/roles/${roleId}/permissions`, { permissionIds })
}

// 获取角色的权限
export const getRolePermissions = (roleId: number) => {
  return request.get<{ data: Permission[] }>(`/system/roles/${roleId}/permissions`)
}
