<template>
  <div class="app-market">
    <div class="search-bar">
      <input 
        v-model="searchKeyword" 
        type="text" 
        placeholder="搜索应用..."
        @keyup.enter="handleSearch"
        class="search-input"
      />
      <button @click="handleSearch" class="search-btn">搜索</button>
      <select v-model="selectedCategory" @change="handleCategoryChange" class="category-select">
        <option value="">全部分类</option>
        <option value="productivity">生产力工具</option>
        <option value="development">开发工具</option>
        <option value="design">设计工具</option>
        <option value="communication">沟通协作</option>
        <option value="business">商业应用</option>
      </select>
    </div>

    <!-- 已安装应用 -->
    <div class="app-section" v-if="installedApps.length > 0">
      <h2 class="section-title">已安装应用</h2>
      <div class="app-grid">
        <div class="app-card" v-for="app in installedApps" :key="app.id" @click="viewAppDetail(app.id)">
          <div class="app-icon">
            <img :src="app.icon || '/default-app-icon.png'" :alt="app.appName" />
          </div>
          <div class="app-info">
            <h3>{{ app.appName }}</h3>
            <p class="version">v{{ app.version }}</p>
            <p class="description">{{ app.description }}</p>
            <div class="app-meta">
              <span class="rating">⭐ {{ app.rating || 0 }}</span>
              <span class="downloads">📥 {{ app.downloads || 0 }}</span>
            </div>
            <div class="app-actions">
              <button @click.stop="uninstallApp(app.appKey)" class="uninstall-btn">卸载</button>
              <a v-if="app.demoUrl" :href="app.demoUrl" target="_blank" class="demo-btn">演示</a>
            </div>
            <div class="status-badge installed">已安装</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 未安装应用 -->
    <div class="app-section" v-if="uninstalledApps.length > 0">
      <h2 class="section-title">未安装应用</h2>
      <div class="app-grid">
        <div class="app-card" v-for="app in uninstalledApps" :key="app.id" @click="viewAppDetail(app.id)">
          <div class="app-icon">
            <img :src="app.icon || '/default-app-icon.png'" :alt="app.appName" />
          </div>
          <div class="app-info">
            <h3>{{ app.appName }}</h3>
            <p class="version">v{{ app.version }}</p>
            <p class="description">{{ app.description }}</p>
            <div class="app-meta">
              <span class="rating">⭐ {{ app.rating || 0 }}</span>
              <span class="downloads">📥 {{ app.downloads || 0 }}</span>
            </div>
            <div class="app-actions">
              <button @click.stop="installApp(app.id, app.appKey)" class="install-btn">安装</button>
              <a v-if="app.demoUrl" :href="app.demoUrl" target="_blank" class="demo-btn">演示</a>
            </div>
            <div class="status-badge uninstalled">未安装</div>
          </div>
        </div>
      </div>
    </div>

    <div class="no-data" v-if="apps.length === 0">
      <p>未找到相关应用</p>
    </div>

    <div class="pagination" v-if="total > 0">
      <button @click="previousPage" :disabled="page <= 1">上一页</button>
      <span>第 {{ page }} 页 / 共 {{ Math.ceil(total / pageSize) }} 页</span>
      <button @click="nextPage" :disabled="page >= Math.ceil(total / pageSize)">下一页</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getApps, getAppsByCategory, searchApps, installApp as apiInstallApp, uninstallApp as apiUninstallApp } from '@/api/appmarket'

const apps = ref([])
const installedApps = computed(() => apps.value.filter(app => app.installed))
const uninstalledApps = computed(() => apps.value.filter(app => !app.installed))
const total = ref(0)
const page = ref(1)
const pageSize = ref(12)
const searchKeyword = ref('')
const selectedCategory = ref('')

const loadApps = async () => {
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value,
    }

    let res
    if (searchKeyword.value) {
      res = await searchApps(searchKeyword.value, params)
    } else if (selectedCategory.value) {
      res = await getAppsByCategory(selectedCategory.value, params)
    } else {
      res = await getApps(params)
    }

    // API直接返回 data 对象（包含 list 和 total 字段）
    if (res && res.data.list) {
      apps.value = res.data.list
      total.value = res.total || 0
      console.log('加载应用成功:', apps.value)
    } else {
      console.warn('意料的响应格式:', res)
      apps.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('加载应用列表失败:', error)
    ElMessage.error('加载应用列表失败')
  }
}

const handleSearch = () => {
  page.value = 1
  loadApps()
}

const handleCategoryChange = () => {
  page.value = 1
  searchKeyword.value = ''
  loadApps()
}

const previousPage = () => {
  if (page.value > 1) {
    page.value--
    loadApps()
  }
}

const nextPage = () => {
  if (page.value < Math.ceil(total.value / pageSize.value)) {
    page.value++
    loadApps()
  }
}

const viewAppDetail = (appId) => {
  console.log('View app detail:', appId)
}

const installApp = async (appId, appKey) => {
  try {
    const res = await apiInstallApp(appId, appKey)
    if (res.data && res.data.success) {
      ElMessage.success(res.data.message || '应用安装成功')
      await loadApps() // 安装后刷新列表
    } else {
      ElMessage.error(res.data?.message || '应用安装失败')
    }
  } catch (error) {
    ElMessage.error('应用安装失败: ' + (error.message || ''))
  }
}

const uninstallApp = async (appKey) => {
  try {
    const res = await apiUninstallApp(appKey)
    if (res.data && res.data.success) {
      ElMessage.success(res.data.message || '应用卸载成功')
      await loadApps() // 卸载后刷新列表
    } else {
      ElMessage.error(res.data?.message || '应用卸载失败')
    }
  } catch (error) {
    ElMessage.error('应用卸载失败: ' + (error.message || ''))
  }
}

onMounted(() => {
  loadApps()
})
</script>

<style scoped>
.app-market {
  padding: 20px;
}

.search-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
  flex-wrap: wrap;
}

.search-input {
  flex: 1;
  min-width: 200px;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.search-input:focus {
  outline: none;
  border-color: #409eff;
}

.search-btn,
.install-btn,
.uninstall-btn,
.demo-btn {
  padding: 8px 16px;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.search-btn:hover,
.install-btn:hover,
.demo-btn:hover {
  background-color: #66b1ff;
}

.uninstall-btn {
  background-color: #f56c6c;
}

.uninstall-btn:hover {
  background-color: #f78989;
}

.category-select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: white;
  font-size: 14px;
  cursor: pointer;
}

/* 应用分类段 */
.app-section {
  margin-bottom: 40px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 2px solid #409eff;
}

.app-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.app-card {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: box-shadow 0.3s;
  position: relative;
}

.app-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.app-icon {
  width: 100%;
  height: 150px;
  background-color: #f5f5f5;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
  overflow: hidden;
}

.app-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.app-info h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #333;
}

.version {
  color: #999;
  font-size: 12px;
  margin: 4px 0;
}

.description {
  color: #666;
  font-size: 13px;
  margin: 8px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 32px;
}

.app-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
  margin: 10px 0;
}

.app-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.install-btn,
.uninstall-btn,
.demo-btn {
  flex: 1;
  padding: 6px 12px;
  font-size: 12px;
  border-radius: 3px;
}

.demo-btn {
  background-color: #67c23a;
  text-decoration: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

.demo-btn:hover {
  background-color: #85ce61;
}

/* 状态休闯 */
.status-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 11px;
  font-weight: 600;
  color: white;
}

.status-badge.installed {
  background-color: #67c23a;
}

.status-badge.uninstalled {
  background-color: #909399;
}

.no-data {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  padding: 20px;
}

.pagination button {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: white;
  cursor: pointer;
}

.pagination button:hover:not(:disabled) {
  background-color: #f5f5f5;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
