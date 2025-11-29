# Power Admin - 通用管理后台系统

> 一个基于 Go-Zero 框架的企业级、可扩展、高性能的通用管理后台系统

![Go Version](https://img.shields.io/badge/Go-1.21+-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Active-brightgreen)

## ✨ 核心特性

### 🔐 完整的RBAC权限管理
- **Casbin权限控制**：基于角色的访问控制（RBAC）
- **无限层级权限**：支持无限父子级权限分组
- **灵活分配**：一个用户可同时属于多个角色
- **权限聚合**：角色权限自动聚合

### 🚀 高效的代码生成工具
- **一键生成CRUD**：自动生成增删改查接口
- **自动生成菜单**：生成对应的菜单结构
- **API文档自动生成**：生成OpenAPI规范文档
- **前端代码生成**：生成Vue 3前端组件和API调用

### 🧩 灵活的插件扩展系统
- **插件隔离**：各插件独立，互不影响
- **动态加载**：无需重新编译即可加载新插件
- **统一接口**：所有插件实现统一的生命周期接口
- **应用市场**：内置插件市场和应用广场

### 📊 通用的会员和API模块
- **共用账号体系**：Web端和API端共用同一账号
- **统一认证**：JWT + Casbin双重认证
- **权限同步**：前后端权限实时同步
- **完整的会员系统**：支持会员中心和API权限

### 🎯 模块化架构
- **微服务友好**：易于扩展为微服务架构
- **清晰的分层**：Handler → Logic → Repository → Database
- **依赖注入**：使用ServiceContext进行依赖注入
- **高内聚低耦合**：各模块职责清晰

## 📦 技术栈

| 技术 | 用途 | 版本 |
|-----|------|------|
| Go-Zero | 后端框架 | v1.9.3+ |
| GORM | ORM框架 | v1.25+ |
| MySQL | 数据库 | 8.0+ |
| Redis | 缓存 | 6.0+ |
| Casbin | 权限管理 | v2.83+ |
| JWT | 认证 | go-jwt/v4 |
| Vue 3 | 前端框架 | 3.3+ |
| TypeScript | 类型检查 | 5.2+ |
| Element Plus | UI组件库 | 2.4+ |
| Pinia | 状态管理 | 2.1+ |

## 🚀 快速开始

### 前置要求
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- Node.js 16+

### 方式一：自动化初始化（推荐）

**Linux/Mac：**
```bash
bash init.sh
```

**Windows：**
```cmd
init.bat
```

### 方式二：手动初始化

**1. 初始化数据库**
```bash
mysql -u root -p < power-admin-server/db/init.sql
```

**2. 启动后端**
```bash
cd power-admin-server
go mod tidy
go build -o power-admin.exe
./power-admin.exe -f etc/power-api.yaml
```

**3. 启动前端**
```bash
cd power-admin-web
npm install
npm run dev
```

**4. 访问应用**
- 前端地址：http://localhost:5173
- 后端API：http://localhost:8888/api/v1
- 默认账户：13800138000 / admin123

## 📖 文档

| 文档 | 描述 |
|-----|------|
| [快速启动指南](./QUICKSTART.md) | 5分钟快速上手 |
| [开发指南](./DEVELOPMENT_GUIDE.md) | 完整的开发规范和API文档 |
| [架构设计](./ARCHITECTURE.md) | 系统架构、设计模式和扩展机制 |

## 🏗️ 项目结构

```
power-admin/
├── power-admin-server/          # Go后端
│   ├── api/                     # API定义
│   ├── pkg/                     # 可复用业务包
│   │   ├── auth/               # 认证（JWT、密码）
│   │   ├── db/                 # 数据库
│   │   ├── models/             # 数据模型
│   │   ├── permission/         # Casbin权限
│   │   └── repository/         # 数据访问层
│   ├── internal/                # 内部模块
│   │   ├── handler/            # HTTP处理器
│   │   ├── logic/              # 业务逻辑
│   │   ├── middleware/         # 中间件
│   │   └── svc/                # 服务容器
│   ├── db/                      # 数据库脚本
│   ├── etc/                     # 配置文件
│   └── power.go                 # 主程序
│
├── power-admin-web/             # Vue 3前端
│   ├── src/
│   │   ├── api/                # API调用
│   │   ├── pages/              # 页面
│   │   ├── components/         # 组件
│   │   ├── stores/             # 状态管理
│   │   └── router/             # 路由
│   └── vite.config.ts          # Vite配置
│
├── QUICKSTART.md               # 快速开始
├── DEVELOPMENT_GUIDE.md        # 开发指南
├── ARCHITECTURE.md             # 架构设计
└── README.md                   # 本文件
```

## 💡 核心功能

### ✅ 已完成

- [x] 用户管理（注册、登录、个人中心）
- [x] 角色管理（创建、编辑、删除、权限分配）
- [x] 权限管理（权限定义、角色授权）
- [x] 菜单管理（无限层级、权限关联）
- [x] 字典管理（数据字典管理）
- [x] API管理（API定义和权限关联）
- [x] JWT认证机制
- [x] 密码加密和验证
- [x] Casbin RBAC权限控制
- [x] 数据库设计和初始化脚本

### 🚀 开发中

- [ ] 所有CRUD API实现
- [ ] 权限验证中间件完整集成
- [ ] 前端管理台UI开发
- [ ] 代码生成器完善
- [ ] 系统日志记录

### 📋 计划中

- [ ] 文件上传下载
- [ ] 消息通知系统
- [ ] 定时任务管理
- [ ] 数据备份恢复
- [ ] 操作审计日志
- [ ] 插件市场
- [ ] 应用广场
- [ ] B2C商城
- [ ] B2B2C商城
- [ ] CRM系统
- [ ] 在线考试
- [ ] 在线投票

## 🔒 安全特性

- ✅ JWT Token认证
- ✅ bcrypt密码加密
- ✅ Casbin权限控制
- ✅ 请求签名验证
- ✅ SQL注入防护（使用ORM）
- ✅ XSS防护（前端框架内置）
- ✅ CORS配置
- ✅ 速率限制（可选）

## 📊 API示例

### 登录
```bash
curl -X POST http://localhost:8888/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"admin123"}'

# 响应
{
  "code": 0,
  "msg": "成功",
  "data": {
    "userId": 1,
    "token": "eyJhbGc...",
    "nickname": "Admin",
    "avatar": ""
  }
}
```

### 获取用户列表（需要认证）
```bash
curl -X GET "http://localhost:8888/api/v1/system/users?page=1&pageSize=10" \
  -H "Authorization: Bearer your_token_here"
```

## 🛠️ 常用开发命令

### 后端

```bash
# 编译
go build -o power-admin.exe

# 运行
go run power.go -f etc/power-api.yaml

# 生成API代码
make gen

# 格式化
go fmt ./...

# 测试
go test ./...
```

### 前端

```bash
# 安装依赖
npm install

# 开发
npm run dev

# 构建
npm run build

# 预览
npm run preview

# 代码检查
npm run lint
```

## 📈 性能指标

| 指标 | 目标 | 备注 |
|-----|------|------|
| 单API响应时间 | < 100ms | 简单查询 |
| 并发连接数 | > 10,000 | 使用连接池 |
| 缓存命中率 | > 80% | 热数据缓存 |
| 数据库查询时间 | < 50ms | 使用索引优化 |

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📝 开发规范

### 命名约定
- **包名**：小写，无下划线
- **函数**：大驼峰式（导出），小驼峰式（私有）
- **变量**：驼峰式
- **常量**：大写加下划线
- **数据库表**：蛇形小写
- **API路径**：蛇形小写

### 代码风格
- 使用 `gofmt` 格式化Go代码
- 使用 `prettier` 格式化前端代码
- 遵循Go官方编码规范
- 添加适当的注释和文档

## 📄 许可证

本项目采用 **MIT License** 开源许可，可自由用于商业项目。

详见 [LICENSE](./LICENSE) 文件。

## 🙏 致谢

感谢以下开源项目和社区：
- [Go-Zero](https://go-zero.dev) - 高性能Go框架
- [GORM](https://gorm.io) - Go ORM库
- [Casbin](https://casbin.org) - 权限管理库
- [Vue.js](https://vuejs.org) - 前端框架
- [Element Plus](https://element-plus.org) - UI组件库

## 📞 联系方式

- 📧 Email: admin@example.com
- 💬 讨论: [GitHub Discussions](https://github.com/your-repo/discussions)
- 🐛 Bug报告: [GitHub Issues](https://github.com/your-repo/issues)

## 📌 更新日志

### v1.0.0 (2024-11-29)
- ✨ 初始版本发布
- ✅ 完成核心功能开发
- 📚 添加完整的文档
- 🚀 支持基础CRUD操作
- 🔐 实现JWT认证和Casbin权限管理

---

**最后更新**：2024年11月29日  
**版本**：1.0.0  
**维护者**：Power Admin Team

如果觉得有帮助，请给个 ⭐️ Star！
