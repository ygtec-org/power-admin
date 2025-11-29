import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { notify } from '../components/Notification'

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
        notify.error(msg || '登录已过期', '未授权')
        // 延迟跳转，给用户时间看到提示
        setTimeout(() => {
          window.location.href = '/login'
          isRedirecting = false
        }, 1500)
      }
      return Promise.reject(new Error(msg || '登录已过期'))
    } else if (code === 403) {
      // 403: 无权限，只提示错误，不跳转
      notify.error(msg || '无权限访问', '权限不足')
      return Promise.reject(new Error(msg || '无权限访问'))
    } else {
      return Promise.reject(new Error(msg || '请求失败'))
    }
  },
  (error) => {
    let message = '请求失败'
    if (error.response) {
      if (error.response.status === 401) {
        // 401: 未授权，跳转登录
        if (!isRedirecting) {
          isRedirecting = true
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          notify.error('登录已过期，请重新登录', '未授权')
          setTimeout(() => {
            window.location.href = '/login'
            isRedirecting = false
          }, 1500)
        }
        message = '登录已过期'
      } else if (error.response.status === 403) {
        // 403: 无权限，只提示错误，不跳转
        message = error.response.data?.msg || '无权限访问'
        notify.error(message, '权限不足')
      } else if (error.response.status === 404) {
        message = '请求的资源不存在'
      } else if (error.response.status === 500) {
        message = '服务器错误'
      } else {
        message = error.response.data?.msg || `HTTP ${error.response.status}`
      }
    } else if (error.request) {
      message = '网络错误，请检查连接'
    }
    return Promise.reject(new Error(message))
  }
)

export default instance
