import request from './request'

// 菜单管理 API
export const getMenus = (params) => request.get('/system/menus', { params })
export const getMenuTree = () => request.get('/system/menus', { params: { page: 1, pageSize: 10000 } })
export const getMenu = (id) => request.get(`/system/menus/${id}`)
export const createMenu = (data) => {
  // 将前端字段名转换为后端字段名
  const payload = {
    parentId: data.parent_id || 0,
    menuName: data.menu_name,
    menuPath: data.menu_path,
    component: data.component || '',
    icon: data.icon || '',
    sort: data.sort || 0,
    menuType: data.menu_type || 1,
  }
  return request.post('/system/menus', payload)
}
export const updateMenu = (id, data) => {
  // 将前端字段名转换为后端字段名
  const payload = {
    id: id,
    menuName: data.menu_name,
    menuPath: data.menu_path,
    component: data.component,
    icon: data.icon,
    sort: data.sort,
    status: data.status,
    menuType: data.menu_type,
  }
  return request.put('/system/menus', payload)
}
export const deleteMenu = (id) => request.delete('/system/menus', { data: { id } })

export const getAPIs = (params) => request.get('/system/apis', { params })
export const getAPI = (id) => request.get(`/system/apis/${id}`)
export const createAPI = (data) => request.post('/system/apis', data)
export const updateAPI = (id, data) => request.put(`/system/apis/${id}`, { ...data, id })
export const deleteAPI = (id) => request.delete(`/system/apis/${id}`)

export const getDicts = (params) => request.get('/content/dicts', { params })
export const getDict = (id) => request.get(`/content/dicts/${id}`)
export const createDict = (data) => request.post('/content/dicts', data)
export const updateDict = (id, data) => request.put(`/content/dicts/${id}`, data)
export const deleteDict = (id) => request.delete(`/content/dicts/${id}`)
