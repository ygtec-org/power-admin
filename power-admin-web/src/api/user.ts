import request from './request'

export interface LoginRequest {
  phone: string
  password: string
}

export interface LoginResponse {
  token: string
  userId: number
  nickname: string
  avatar?: string
}

export interface Role {
  id: number
  name: string
  description: string
  status: number
}

export interface AssignRolesToUserRequest {
  roleIds: number[]
}

export const login = (data: LoginRequest) => {
  return request.post<LoginResponse>('/auth/login', data)
}

export const logout = () => {
  return request.post('/auth/logout')
}

export const getUserInfo = () => {
  return request.get('/auth/info')
}

export const getUsers = (params?: any) => {
  return request.get('/system/users', { params })
}

export const createUser = (data: any) => {
  return request.post('/system/users', data)
}

export const updateUser = (id: number, data: any) => {
  return request.put('/system/users', { id, ...data })
}

export const deleteUser = (id: number) => {
  return request.delete('/system/users', { data: { id } })
}

// 为用户分配角色
export const assignRolesToUser = (userId: number, roleIds: number[]) => {
  return request.post(`/system/users/${userId}/roles`, { roleIds })
}

// 获取用户的角色
export const getUserRoles = (userId: number) => {
  return request.get<{ data: Role[] }>(`/system/users/${userId}/roles`)
}
