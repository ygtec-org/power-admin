import { ElMessage, ElNotification } from 'element-plus'

export const notify = {
  success: (message: string, title: string = '成功') => {
    ElNotification({
      title,
      message,
      type: 'success',
      duration: 3000,
      position: 'top-right',
    })
  },
  error: (message: string, title: string = '错误') => {
    ElNotification({
      title,
      message,
      type: 'error',
      duration: 4000,
      position: 'top-right',
    })
  },
  warning: (message: string, title: string = '警告') => {
    ElNotification({
      title,
      message,
      type: 'warning',
      duration: 3000,
      position: 'top-right',
    })
  },
  info: (message: string, title: string = '提示') => {
    ElNotification({
      title,
      message,
      type: 'info',
      duration: 3000,
      position: 'top-right',
    })
  },
}

export const toast = {
  success: (message: string) => {
    ElMessage.success(message)
  },
  error: (message: string) => {
    ElMessage.error(message)
  },
  warning: (message: string) => {
    ElMessage.warning(message)
  },
  info: (message: string) => {
    ElMessage.info(message)
  },
}

// 顶部中心消息提示（推荐用于重要提示）
export const message = {
  success: (message: string) => {
    ElMessage({
      message,
      type: 'success',
      duration: 3000,
      offset: 20,
    })
  },
  error: (message: string) => {
    ElMessage({
      message,
      type: 'error',
      duration: 4000,
      offset: 20,
    })
  },
  warning: (message: string) => {
    ElMessage({
      message,
      type: 'warning',
      duration: 3000,
      offset: 20,
    })
  },
  info: (message: string) => {
    ElMessage({
      message,
      type: 'info',
      duration: 3000,
      offset: 20,
    })
  },
}
