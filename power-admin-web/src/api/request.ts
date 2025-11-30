import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

const instance: AxiosInstance = axios.create({
  baseURL: '/api/admin',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
instance.interceptors.request.use(
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

// 标志位：防止无限跳转
let isRedirecting = false

// 响应拦截器
instance.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, msg, data } = response.data
    if (code === 0) {
      return Promise.resolve({ data, msg } as any)
    } else if (code === 401) {
      // 401: 未授权，跳转登录
      if (!isRedirecting) {
        isRedirecting = true
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        ElMessage({
          message: msg || '登录已过期，请重新登录',
          type: 'error',
          duration: 4000,
          offset: 20,
        })
        // 延迟跳转，给用户时间看到提示
        setTimeout(() => {
          window.location.href = '/login'
          isRedirecting = false
        }, 1500)
      }
      return Promise.reject(new Error(msg || '登录已过期'))
    } else if (code === 403) {
      // 403: 无权限，显示顶部提示
      ElMessage({
        message: msg || '无权限,请联系管理员开通',
        type: 'error',
        duration: 4000,
        offset: 20,
      })
      return Promise.reject(new Error(msg || '无权限,请联系管理员开通'))
    } else {
      return Promise.reject(new Error(msg || '请求失败'))
    }
  },
  (error) => {
    let errorMsg = '请求失败'
    if (error.response) {
      if (error.response.status === 401) {
        // 401: 未授权，跳转登录
        if (!isRedirecting) {
          isRedirecting = true
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          ElMessage({
            message: '登录已过期，请重新登录',
            type: 'error',
            duration: 4000,
            offset: 20,
          })
          setTimeout(() => {
            window.location.href = '/login'
            isRedirecting = false
          }, 1500)
        }
        errorMsg = '登录已过期'
      } else if (error.response.status === 403) {
        // 403: 无权限，显示顶部提示
        const msg403 = error.response.data?.msg || '无权限,请联系管理员开通'
        ElMessage({
          message: msg403,
          type: 'error',
          duration: 4000,
          offset: 20,
        })
        return Promise.reject(new Error(msg403))
      } else if (error.response.status === 404) {
        errorMsg = '请求的资源不存在'
      } else if (error.response.status === 500) {
        errorMsg = '服务器错误'
      } else {
        errorMsg = error.response.data?.msg || `HTTP ${error.response.status}`
      }
    } else if (error.request) {
      errorMsg = '网络错误，请检查连接'
    }
    return Promise.reject(new Error(errorMsg))
  }
)

export default instance
