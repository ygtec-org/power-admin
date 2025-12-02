<template>
  <div class="cms-menu">
    <div class="menu-header">
      <h3>CMS 管理</h3>
    </div>
    <ul class="menu-list">
      <li>
        <router-link to="/cms/content" class="menu-item" active-class="active">
          <span class="icon">📄</span>
          <span class="text">内容管理</span>
        </router-link>
      </li>
      <li>
        <router-link to="/cms/category" class="menu-item" active-class="active">
          <span class="icon">📁</span>
          <span class="text">分类管理</span>
        </router-link>
      </li>
      <li>
        <router-link to="/cms/tag" class="menu-item" active-class="active">
          <span class="icon">🏷️</span>
          <span class="text">标签管理</span>
        </router-link>
      </li>
      <li>
        <router-link to="/cms/comment" class="menu-item" active-class="active">
          <span class="icon">💬</span>
          <span class="text">评论管理</span>
        </router-link>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
// CMS菜单组件
</script>

<style scoped>
.cms-menu {
  background: #f8f9fa;
  border-radius: 4px;
  overflow: hidden;
}

.menu-header {
  padding: 15px;
  border-bottom: 1px solid #ddd;
}

.menu-header h3 {
  margin: 0;
  font-size: 16px;
  color: #333;
  font-weight: 600;
}

.menu-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 15px;
  color: #666;
  text-decoration: none;
  border-left: 3px solid transparent;
  transition: all 0.2s;
}

.menu-item:hover {
  background: #e9ecef;
  color: #333;
}

.menu-item.active {
  background: #e7f3ff;
  border-left-color: #007bff;
  color: #007bff;
  font-weight: 600;
}

.icon {
  font-size: 18px;
  display: inline-block;
  width: 24px;
  text-align: center;
}

.text {
  flex: 1;
}
</style>
