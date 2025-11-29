<template>
  <div v-if="visible" class="confirm-dialog-overlay" @click="handleClose">
    <div class="confirm-dialog" @click.stop>
      <div class="confirm-dialog-header">
        <h3>{{ title }}</h3>
        <button class="close-btn" @click="handleClose">×</button>
      </div>
      <div class="confirm-dialog-body">
        <p>{{ message }}</p>
      </div>
      <div class="confirm-dialog-footer">
        <button class="btn-cancel" @click="handleCancel">取消</button>
        <button class="btn-confirm" @click="handleConfirm">确定</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '确认'
  },
  message: {
    type: String,
    default: '确定要执行此操作吗？'
  }
})

const emit = defineEmits(['update:visible', 'confirm', 'cancel'])

const handleClose = () => {
  emit('update:visible', false)
  emit('cancel')
}

const handleConfirm = () => {
  emit('update:visible', false)
  emit('confirm')
}

const handleCancel = () => {
  emit('update:visible', false)
  emit('cancel')
}

// 监听 visible 变化，如果变为 true 则聚焦到确认按钮
watch(() => props.visible, (newVal) => {
  if (newVal) {
    setTimeout(() => {
      const confirmBtn = document.querySelector('.btn-confirm')
      if (confirmBtn) confirmBtn.focus()
    }, 100)
  }
})
</script>

<style scoped>
.confirm-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.confirm-dialog {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 400px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.confirm-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 20px 10px;
  border-bottom: 1px solid #e6e9f0;
}

.confirm-dialog-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
}

.close-btn:hover {
  color: #666;
}

.confirm-dialog-body {
  padding: 20px;
}

.confirm-dialog-body p {
  margin: 0;
  font-size: 14px;
  color: #666;
  line-height: 1.5;
}

.confirm-dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid #e6e9f0;
}

.btn-cancel {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  color: #333;
}

.btn-cancel:hover {
  border-color: #999;
}

.btn-confirm {
  padding: 8px 16px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.btn-confirm:hover {
  background: #5a6fd8;
}
</style>