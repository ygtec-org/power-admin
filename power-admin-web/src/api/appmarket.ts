import axios from 'axios'
import request from './request'

// 创建一个用于应用市场API的axios实例
const apiInstance = axios.create({
  baseURL: '/api/admin',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 为应用市场API实例添加请求拦截器
apiInstance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 为应用市场API实例添加响应拦截器
apiInstance.interceptors.response.use(
  (response) => {
    // 应用市场API直接返回数据中的data字段
    if (response.status === 200 && response.data) {
      return Promise.resolve(response.data)
    }
    return Promise.reject(new Error('请求失败'))
  },
  (error) => {
    return Promise.reject(error)
  }
)

export const getApps = (params: any) => apiInstance.get('/app-market/list', { params })
export const getApp = (id: number) => apiInstance.get(`/app-market/${id}`)
export const getAppsByCategory = (category: string, params: any) => apiInstance.get(`/app-market/category/${category}`, { params })
export const searchApps = (keyword: string, params: any) => apiInstance.get('/app-market/search', { params: { keyword, ...params } })
export const installApp = (appId: number, appKey: string) => apiInstance.post('/app-market/install', null, { params: { appId, appKey } })
export const uninstallApp = (appKey: string) => apiInstance.post('/app-market/uninstall', null, { params: { appKey } })
export const checkInstallStatus = (appKey: string) => apiInstance.get('/app-market/install-status', { params: { appKey } })
