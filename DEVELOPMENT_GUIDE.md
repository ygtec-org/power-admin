# Power Admin 开发指南

## 项目概述

Power Admin 是一个基于 Go-Zero 框架的通用管理后台系统，支持 RBAC（基于角色的访问控制）和插件扩展。

## 技术栈

- **后端框架**: go-zero
- **数据库**: MySQL 8.0+
- **缓存**: Redis
- **权限管理**: Casbin
- **认证**: JWT
- **密码加密**: bcrypt
- **前端**: Vue 3 + TypeScript

## 项目结构

```
power-admin-server/
├── api/                    # API定义文件（goctl生成）
├── common/                 # 通用模块
│   ├── constant/          # 常量定义
│   ├── model/             # 公共模型
│   └── response/          # 统一响应结构
├── pkg/                    # 业务包
│   ├── auth/              # 认证相关（JWT、密码加密）
│   ├── db/                # 数据库初始化
│   ├── cache/             # 缓存管理（Redis）
│   ├── models/            # 数据模型（GORM）
│   ├── permission/        # Casbin权限管理
│   └── repository/        # 数据访问层
├── internal/              # 内部模块
│   ├── config/            # 配置管理
│   ├── handler/           # HTTP处理器
│   ├── logic/             # 业务逻辑
│   ├── svc/               # 服务上下文
│   ├── types/             # 类型定义
│   └── middleware/        # 中间件
├── db/                    # 数据库脚本
├── etc/                   # 配置文件
├── Makefile               # 构建脚本
├── go.mod                 # Go模块定义
└── power.go               # 主入口

power-admin-web/
├── src/
│   ├── components/        # Vue组件
│   ├── pages/             # 页面
│   ├── api/               # API调用
│   ├── stores/            # Pinia状态管理
│   └── App.vue
├── package.json
└── vite.config.ts
```

## 快速开始

### 1. 环境准备

```bash
# 安装MySQL
# 创建数据库并导入初始化脚本
mysql -u root -p < power-admin-server/db/init.sql

# 启动Redis（需要Redis服务）
redis-server

# 或使用Docker启动Redis
docker run -d -p 6379:6379 redis:latest
```

### 2. 后端启动

```bash
cd power-admin-server

# 下载依赖
go mod tidy

# 编译
go build -o power-admin.exe

# 启动服务
./power-admin.exe -f etc/power-api.yaml
```

### 3. 前端启动

```bash
cd power-admin-web

# 安装依赖
npm install

# 开发模式启动
npm run dev

# 生产构建
npm run build
```

## 代码生成

### 使用 goctl 生成API代码

编辑 `api/admin.api` 定义API，然后执行：

```bash
make gen
```

这会自动生成：
- Handler（处理器）
- Router（路由）
- Types（类型）

## 数据库设置

### 初始化脚本

运行 `db/init.sql` 会创建以下表：
- users（用户表）
- roles（角色表）
- permissions（权限表）
- menus（菜单表）
- user_roles（用户-角色关联表）
- role_permissions（角色-权限关联表）
- dictionaries（字典表）
- apis（API管理表）
- plugins（插件表）
- logs（系统日志表）

## API使用示例

### 1. 用户注册

```bash
curl -X POST http://localhost:8888/api/v1/admin/register \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "password": "password123",
    "nickname": "User123",
    "gender": 1
  }'
```

### 2. 用户登录

```bash
curl -X POST http://localhost:8888/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "password": "password123"
  }'

# 响应
{
  "code": 0,
  "msg": "成功",
  "data": {
    "userId": 1,
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "nickname": "User123",
    "avatar": ""
  }
}
```

### 3. 获取用户列表（需要认证）

```bash
curl -X GET "http://localhost:8888/api/v1/system/users?page=1&pageSize=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## 权限管理（RBAC）

### 权限检查原理

系统使用Casbin实现RBAC权限控制：

```
用户 -> 角色 -> 权限 -> API
```

### 添加权限检查中间件

在API路由中添加权限检查：

```go
@handler GetUser
@permission get,users,view
get /system/users/:id returns (UserInfo)
```

## 配置文件

编辑 `etc/power-api.yaml`：

```yaml
Name: power-api
Host: 0.0.0.0
Port: 8888

Mysql:
  DataSource: "root:root@tcp(127.0.0.1:3306)/power_admin?charset=utf8mb4&parseTime=True&loc=Local"

Redis:
  Host: 127.0.0.1
  Port: 6379
  Pass: ""
  Db: 0
```

## 开发规范

### 1. 命名规范

- 包名：小写，不使用下划线
- 函数：驼峰式命名，首字母大写表示导出
- 变量：驼峰式命名
- 常量：大写加下划线（CONSTANT_NAME）

### 2. 项目结构约定

- `pkg/` 目录：可复用的业务包
- `internal/` 目录：项目内部使用的包
- `handler/` 目录：HTTP请求处理
- `logic/` 目录：业务逻辑实现
- `repository/` 目录：数据访问层

### 3. 错误处理

统一使用以下方式返回错误：

```go
// 业务错误
return nil, fmt.Errorf("error message")

// 系统错误
return nil, errors.New("error message")
```

## 扩展开发

### 1. 添加新的API

1. 编辑 `api/admin.api` 添加新的API定义
2. 执行 `make gen` 生成代码
3. 在 `logic/` 目录实现业务逻辑
4. 重新编译运行

### 2. 添加新的权限

1. 在数据库 `permissions` 表中添加新权限
2. 为角色分配权限
3. 在API处理器中验证权限

### 3. 创建插件

插件开发文档见 `PLUGIN_DEVELOPMENT.md`

## 调试

### 启用日志

在 `etc/power-api.yaml` 中配置日志级别：

```yaml
Logx:
  Level: debug
```

### 查看数据库日志

GORM会在控制台输出SQL日志（开发环境）

## 常见问题

### Q: 如何修改JWT Secret？
A: 编辑 `pkg/auth/jwt.go` 中的 `JwtSecret` 变量，但建议从环境变量或配置文件读取。

### Q: 如何添加新的数据表？
A: 编辑对应的SQL脚本，在 `pkg/models/` 创建GORM模型，在 `pkg/repository/` 创建仓储。

### Q: 如何集成其他认证方式？
A: 修改 `internal/middleware/auth.go` 中的认证逻辑。

## 部署

### Docker部署

构建镜像：

```bash
docker build -t power-admin:latest .
docker run -d --name power-admin \
  -p 8888:8888 \
  -e MYSQL_DSN="user:password@tcp(host:3306)/power_admin" \
  power-admin:latest
```

## 相关文档

- [API文档](./API_DOCUMENTATION.md)
- [数据库设计](./DATABASE_DESIGN.md)
- [插件开发指南](./PLUGIN_DEVELOPMENT.md)
- [前端开发指南](./power-admin-web/README.md)

## 贡献指南

欢迎提交Issue和Pull Request！

## License

MIT
