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
