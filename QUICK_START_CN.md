# Power Admin - 快速开始指南（5分钟上手）

> 本文档指导您在 5 分钟内快速启动 Power Admin 项目

## 前置条件

请确保您的系统已安装以下软件：

- **MySQL** 8.0+ （需要运行中）
- **Go** 1.19+
- **Node.js** 16+
- **npm** 8+

## 🚀 一键启动（最简单）

Windows 用户，直接运行：

```bash
init.bat
```

脚本会自动完成：
1. ✅ 检查数据库连接
2. ✅ 创建数据库
3. ✅ 导入 SQL 脚本
4. ✅ 编译后端程序
5. ✅ 打印启动指令

---

## 📝 手动启动步骤

### 第 1 步：初始化数据库（1 分钟）

打开 MySQL 客户端或 CMD：

```bash
# 创建数据库
mysql -u root -proot -e "CREATE DATABASE power_admin;"

# 导入初始化脚本
mysql -u root -proot power_admin < power-admin-server\db\init.sql
```

**验证成功**：没有报错信息

### 第 2 步：启动后端服务（1 分钟）

打开 **命令窗口 1**：

```bash
cd power-admin-server

# 编译
go build -o power-admin.exe

# 运行
.\power-admin.exe -f etc\power-api.yaml
```

**预期输出**：
```
Database connected successfully
Redis connected successfully
Starting server at 0.0.0.0:8888...
```

✅ 后端启动成功，不要关闭此窗口

### 第 3 步：启动前端应用（2 分钟）

打开 **新的命令窗口 2**：

```bash
cd power-admin-web

# 安装依赖（首次需要 1-2 分钟）
npm install

# 启动开发服务器
npm run dev
```

**预期输出**：
```
  VITE v5.0.2  ready in xxx ms

  ➜  Local:   http://localhost:5173/
  ➜  press h to show help
```

✅ 前端启动成功

### 第 4 步：打开浏览器访问

在浏览器中打开：

```
http://localhost:5173
```

您应该看到登录页面。

### 第 5 步：登录应用（1 分钟）

使用默认登录凭证：

- **手机号**：13800138000
- **密码**：admin123

点击"登录"按钮，进入管理台。

---

## ✨ 见证奇迹的时刻

登录成功后，您将看到：

1. 📊 **仪表板** - 显示系统统计数据
2. 📋 **左侧菜单** - 系统管理、内容管理等模块
3. 👤 **用户管理** - 查看和管理用户列表
4. 🎯 **其他模块** - 角色、菜单、权限、API、字典等

---

## 🎯 常见问题速查表

| 问题 | 解决方案 |
|-----|--------|
| **MySQL 连接失败** | 1. 确认 MySQL 已启动<br>2. 检查用户名密码（默认：root/root）<br>3. 修改 etc/power-api.yaml 中的数据库配置 |
| **后端编译失败** | 1. 运行 `go mod download`<br>2. 检查 Go 版本 `go version`<br>3. 检查网络连接 |
| **前端启动失败** | 1. 运行 `npm install` 重新安装<br>2. 删除 node_modules 文件夹后重试<br>3. 尝试 `npm ci` 代替 `npm install` |
| **登录显示 404** | 1. 检查后端是否运行在 :8888<br>2. 检查浏览器控制台查看具体错误<br>3. 尝试硬刷新（Ctrl+Shift+R） |
| **权限校验失败** | 1. 确保 rbac_model.conf 存在<br>2. 检查数据库中 casbin_rule 表有数据<br>3. 查看后端日志获取详细错误 |

---

## 📊 项目架构一览

```
Power Admin
├── 后端 (Go-Zero)
│   ├── 用户管理 API
│   ├── 权限管理 (Casbin RBAC)
│   ├── 角色和菜单管理
│   └── 字典和 API 管理
│
└── 前端 (Vue 3 + TypeScript)
    ├── 仪表板
    ├── 用户管理页面
    ├── 系统管理模块
    └── 内容管理模块
```

---

## 🔧 配置文件位置

如需修改配置，编辑这些文件：

| 文件 | 说明 |
|-----|------|
| `power-admin-server/etc/power-api.yaml` | 后端配置（端口、数据库、Redis） |
| `power-admin-server/etc/rbac_model.conf` | Casbin 权限模型配置 |
| `power-admin-web/vite.config.ts` | 前端打包配置和代理 |
| `power-admin-web/src/api/request.ts` | API 基础配置（超时、拦截器） |

---

## 🎓 下一步学习

快速启动成功后，可以继续阅读：

1. **[项目完整文档](./README.md)** - 了解项目全貌
2. **[优化总结](./OPTIMIZATION_SUMMARY.md)** - 了解此次优化内容
3. **[开发指南](./DEVELOPMENT_GUIDE.md)** - 学习如何扩展功能
4. **[架构设计](./ARCHITECTURE.md)** - 深入理解系统设计

---

## 💡 提示

- **热重载**：前端修改代码会自动刷新，无需手动重启
- **API 调用**：检查浏览器 F12 Network 标签查看 API 调用
- **权限调试**：查看浏览器 Console 查看权限检查结果
- **数据初始化**：初始用户存储在数据库中，可通过 API 新增

---

## ❓ 仍有问题？

1. 查看对应的详细文档文件
2. 检查后端和前端的日志输出
3. 参考项目中的代码注释
4. 查阅 [Casbin 官方文档](https://casbin.org/zh/)

---

**现在开始享受 Power Admin 吧！** 🎉

开发快乐！✨
