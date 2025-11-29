<template>
  <div class="dashboard">
    <h1>仪表板</h1>
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon users">👤</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.totalUsers }}</div>
          <div class="stat-label">总用户数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon roles">🎯</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.totalRoles }}</div>
          <div class="stat-label">角色数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon menus">📋</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.totalMenus }}</div>
          <div class="stat-label">菜单数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon apis">🔗</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.totalApis }}</div>
          <div class="stat-label">API数</div>
        </div>
      </div>
    </div>

    <div class="welcome-card">
      <h2>欢迎使用 Power Admin</h2>
      <p>这是一个基于 Go-Zero + Vue 3 的企业级管理后台系统</p>
      <div class="features">
        <div class="feature">
          <span class="check">✓</span> 完整的RBAC权限管理
        </div>
        <div class="feature">
          <span class="check">✓</span> 灵活的菜单和权限配置
        </div>
        <div class="feature">
          <span class="check">✓</span> 高效的API管理
        </div>
        <div class="feature">
          <span class="check">✓</span> 模块化架构设计
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const stats = ref({
  totalUsers: 0,
  totalRoles: 0,
  totalMenus: 0,
  totalApis: 0,
})

onMounted(() => {
  // 模拟数据，实际应从API获取
  stats.value = {
    totalUsers: 10,
    totalRoles: 3,
    totalMenus: 12,
    totalApis: 25,
  }
})
</script>

<style scoped>
.dashboard {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

h1 {
  margin: 0 0 24px 0;
  font-size: 28px;
  color: #333;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.stat-icon {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  border-radius: 8px;
}

.stat-icon.users {
  background: #e8f4f8;
}

.stat-icon.roles {
  background: #f0e8f8;
}

.stat-icon.menus {
  background: #f8f4e8;
}

.stat-icon.apis {
  background: #e8f8f0;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 4px;
}

.welcome-card {
  padding: 32px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.welcome-card h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  color: #333;
}

.welcome-card > p {
  color: #666;
  margin: 0 0 24px 0;
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
}

.feature {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
  color: #666;
}

.check {
  color: #667eea;
  font-weight: 600;
  font-size: 18px;
}
</style>
