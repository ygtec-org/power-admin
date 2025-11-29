import request from './request'

export const getApps = (params: any) => request.get('/app-market/list', { params })
export const getApp = (id: number) => request.get(`/app-market/${id}`)
export const getAppsByCategory = (category: string, params: any) => request.get(`/app-market/category/${category}`, { params })
export const searchApps = (keyword: string, params: any) => request.get('/app-market/search', { params: { keyword, ...params } })
