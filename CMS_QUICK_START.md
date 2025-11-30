# CMS 系统快速启动指南

## 系统架构概览

本项目包含三个独立的应用：

1. **power-admin-web**（管理后台）- 5173 端口
   - 包含系统管理、内容管理等功能
   - 集成了 CMS 内容管理模块（内容、分类、标签、评论、用户、发布管理）

2. **power-cms-server**（CMS 后端服务）- 8801 端口
   - 独立的微服务
   - 提供所有 CMS 相关的 API

3. **power-cms-web**（CMS 前台网站）- 5174 端口
   - 面向用户的展示平台
   - 展示已发布的内容、分类浏览、标签浏览等

## 启动步骤

### 第1步：启动 CMS 后端服务（8801 端口）

```bash
cd d:\Workspace\project\app\power-admin\plugins\power-cms\power-cms-server

# 编译（如果需要）
go build -o power-cms.exe

# 运行
.\power-cms.exe
```

或者直接运行编译好的可执行文件（如果存在）。

### 第2步：启动管理后台前端（5173 端口）

```bash
cd d:\Workspace\project\app\power-admin\power-admin-web

# 安装依赖（如果还没安装）
npm install

# 启动开发服务器
npm run dev
```

访问：http://localhost:5173

### 第3步：启动 CMS 前台网站（5174 端口）

```bash
cd d:\Workspace\project\app\power-cms-web

# 安装依赖（如果还没安装）
npm install

# 启动开发服务器
npm run dev
```

访问：http://localhost:5174

## CMS 管理后台功能

在 power-admin-web 的管理后台，进入左侧菜单 "CMS内容系统"，可以访问以下功能：

### 1. 内容管理（/cms/content）
- 创建、编辑、删除文章内容
- 搜索和过滤功能
- 状态管理（草稿、已发布、已删除）
- 分页支持

### 2. 分类管理（/cms/category）
- 创建、编辑、删除分类
- 支持树形结构
- 分类名称和 Slug 管理

### 3. 标签管理（/cms/tag）
- 创建、编辑、删除标签
- 颜色自定义
- 标签使用数统计

### 4. 评论管理（/cms/comment）
- 评论审核（待审核、已通过、已拒绝）
- 批准、拒绝、删除评论
- 评论过滤

### 5. 用户管理（/cms/user）
- 查看访客用户
- 启用、禁用用户账户
- 删除用户

### 6. 发布管理（/cms/publish）
- 立即发布内容
- 定时发布（选择发布时间）
- 取消定时发布
- 发布状态查询

## CMS 前台功能

在 power-cms-web 中，用户可以：

### 1. 首页（/）
- 浏览最新发布的文章
- 按时间倒序显示
- 支持分页

### 2. 内容详情页（/content/:id）
- 查看完整文章内容
- 查看文章元信息（作者、发布时间、浏览次数）
- 查看相关标签
- 查看所属分类
- 浏览已通过审核的评论

### 3. 分类浏览（/categories）
- 查看所有分类
- 点击分类查看该分类下的所有文章
- 显示每个分类的文章数

### 4. 标签浏览（/tags）
- 标签云展示（根据使用频率显示不同大小）
- 点击标签查看该标签下的所有文章

## API 代理配置

### power-admin-web 的代理配置（vite.config.ts）
```javascript
server: {
  port: 5173,
  proxy: {
    '/api/admin': {
      target: 'http://localhost:8888',  // 管理后台服务
      changeOrigin: true,
    },
    '/api/cms': {
      target: 'http://localhost:8801',   // CMS服务
      changeOrigin: true,
    },
  },
},
```

### power-cms-web 的代理配置（vite.config.js）
```javascript
server: {
  port: 5174,
  proxy: {
    '/api/cms': {
      target: 'http://localhost:8801',
      changeOrigin: true,
    },
  },
},
```

## 常见问题

### Q: CMS 后端服务启动后无法连接
A: 检查 CMS 后端服务配置文件（cms.yaml），确保：
- 数据库连接正确
- 服务监听的端口是 8801
- 防火墙没有阻止该端口

### Q: 管理后台无法调用 CMS API
A: 检查以下几点：
- CMS 后端服务是否正常运行
- vite.config.ts 中的 API 代理配置是否正确
- 前端开发服务器是否正常运行

### Q: 前台网站显示"暂无内容"
A: 可能原因：
- 需要在管理后台发布一些文章
- 确保发布时设置了发布状态为"已发布"
- 检查内容是否有分类和标签关联

## 项目文件结构

```
power-admin/
 power-admin-web/                 # 管理后台前端
    src/
       pages/
          cms/                # CMS 管理页面
             content/        # 内容管理
             category/       # 分类管理
             tag/            # 标签管理
             comment/        # 评论管理
             user/           # 用户管理
             publish/        # 发布管理
          system/             # 系统管理
       router/index.ts         # 路由配置
       layout/Layout.vue       # 布局组件（包含 CMS 菜单）
       ...
    ...

 plugins/
    power-cms/
        power-cms-server/       # CMS 后端服务
            internal/
               logic/          # 业务逻辑
               repository/     # 数据访问层
               handler/        # HTTP 处理器
               ...
            cms.yaml            # 配置文件
            ...

power-cms-web/                       # CMS 前台网站
 src/
    pages/
       Home.vue                # 首页
       ContentDetail.vue       # 内容详情页
       Categories.vue          # 分类浏览
       Tags.vue                # 标签浏览
    router/
       index.js                # 路由配置
    App.vue                     # 根组件
    main.js                     # 应用入口
 ...
```

## 开发建议

1. 修改管理后台的 CMS 功能时，只需修改 power-admin-web 中的相应文件
2. 需要修改后端 API 时，修改 power-cms-server 中的 Logic、Repository 或 Handler
3. 前台页面修改时，修改 power-cms-web 中的相应页面组件

## 部署提示

- 生产环境请使用 `npm run build` 构建优化版本
- 配置好后端服务的数据库连接
- 设置正确的跨域 (CORS) 配置
- 使用反向代理（如 Nginx）统一管理前后端服务的端口映射

