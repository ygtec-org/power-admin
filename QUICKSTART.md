# Power Admin 快速启动指南

## 前置要求

- Go 1.21+ 
- MySQL 8.0+
- Redis 6.0+
- Node.js 16+

## 一、快速设置

### 1. 初始化数据库

```bash
# 使用MySQL客户端登录
mysql -u root -p

# 执行初始化脚本
source power-admin-server/db/init.sql
```

### 2. 启动Redis

```bash
# 本地启动
redis-server

# 或使用Docker
docker run -d -p 6379:6379 --name redis redis:latest
```

### 3. 编译后端

```bash
cd power-admin-server

# 下载依赖
go mod tidy

# 编译
go build -o power-admin.exe

# 或直接运行（不编译）
go run power.go -f etc/power-api.yaml
```

### 4. 启动前端

```bash
cd power-admin-web

# 安装依赖
npm install

# 开发服务器
npm run dev

# 生产构建
npm run build
```

## 二、测试登录

### 初始管理员账户

```
手机号: 13800138000
密码: admin123
```

### API调用示例

#### 1. 登录

```bash
curl -X POST http://localhost:8888/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"admin123"}'
```

**返回示例:**

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "userId": 1,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "nickname": "Admin",
    "avatar": ""
  }
}
```

#### 2. 获取用户列表

```bash
TOKEN="your_token_here"

curl -X GET "http://localhost:8888/api/v1/system/users?page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"
```

#### 3. 创建用户

```bash
curl -X POST http://localhost:8888/api/v1/system/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "testuser",
    "phone": "13900000000",
    "nickname": "Test User",
    "password": "password123",
    "gender": 1
  }'
```

## 三、项目结构速览

```
power-admin/
├── power-admin-server/          # Go后端
│   ├── api/                     # API定义（goctl生成）
│   ├── pkg/
│   │   ├── auth/               # JWT认证
│   │   ├── db/                 # 数据库
│   │   ├── cache/              # Redis缓存
│   │   ├── models/             # 数据模型
│   │   ├── permission/         # Casbin权限
│   │   └── repository/         # 数据访问层
│   ├── internal/
│   │   ├── config/             # 配置
│   │   ├── handler/            # HTTP处理器
│   │   ├── logic/              # 业务逻辑
│   │   ├── middleware/         # 中间件
│   │   └── svc/                # 服务上下文
│   ├── db/init.sql             # 数据库初始化脚本
│   ├── etc/power-api.yaml      # 配置文件
│   └── power.go                # 主入口
│
└── power-admin-web/             # Vue3前端
    ├── src/
    │   ├── components/          # Vue组件
    │   ├── pages/              # 页面
    │   ├── api/                # API调用
    │   └── stores/             # Pinia状态管理
    └── vite.config.ts
```

## 四、常用开发命令

### 生成API代码

编辑 `api/admin.api` 后执行:

```bash
cd power-admin-server
make gen
```

### 运行后端

```bash
# 开发模式（支持热重载建议使用air）
air

# 或直接运行
go run power.go -f etc/power-api.yaml

# 指定配置文件
./power-admin.exe -f etc/power-api.yaml
```

### 运行前端

```bash
cd power-admin-web
npm run dev     # 开发
npm run build   # 构建
npm run preview # 预览
```

## 五、配置文件说明

### `power-admin-server/etc/power-api.yaml`

```yaml
Name: power-api              # 服务名称
Host: 0.0.0.0              # 监听地址
Port: 8888                 # 监听端口

Mysql:
  DataSource: "root:root@tcp(127.0.0.1:3306)/power_admin?charset=utf8mb4&parseTime=True&loc=Local"

Redis:
  Host: 127.0.0.1
  Port: 6379
  Pass: ""
  Db: 0
```

## 六、主要功能

### ✅ 已实现

- [x] 用户管理（登录、注册、个人信息）
- [x] JWT认证
- [x] RBAC权限管理（Casbin）
- [x] 用户、角色、权限、菜单、字典数据库表
- [x] 基础仓储层
- [x] 错误处理和日志

### 🚀 开发中

- [ ] 用户、角色、权限、菜单、字典 CRUD API实现
- [ ] 权限验证中间件完整集成
- [ ] 前端管理台
- [ ] 代码生成器
- [ ] 插件系统

### 📋 计划中

- [ ] 文件上传
- [ ] 消息队列
- [ ] 定时任务
- [ ] 系统日志记录
- [ ] 操作审计
- [ ] 数据备份恢复

## 七、常见问题

### Q: 编译失败

**A:** 确保已运行 `go mod tidy` 下载所有依赖

```bash
go mod tidy
go build -o power-admin.exe
```

### Q: 连接数据库失败

**A:** 检查MySQL是否运行，配置文件中的DSN是否正确

```bash
# 测试MySQL连接
mysql -u root -p -h 127.0.0.1 -P 3306
```

### Q: 连接Redis失败

**A:** 确保Redis服务正常运行

```bash
# 测试Redis连接
redis-cli ping
# 应该返回 PONG
```

### Q: 前端无法连接后端API

**A:** 检查：
1. 后端服务是否运行在 `http://localhost:8888`
2. 前端的 API 基础URL配置是否正确
3. 浏览器控制台是否有CORS错误

## 八、性能优化建议

### 数据库
- 启用查询缓存
- 创建必要的索引
- 使用连接池

### Redis
- 合理设置过期时间
- 监控内存使用
- 定期清理过期数据

### 应用
- 启用Gzip压缩
- 实现API限流
- 使用CDN加速静态资源

## 九、部署

### Docker部署

```bash
# 构建镜像
docker build -t power-admin:latest .

# 运行容器
docker run -d --name power-admin \
  -p 8888:8888 \
  -e MYSQL_DSN="root:password@tcp(mysql:3306)/power_admin" \
  -e REDIS_HOST="redis" \
  power-admin:latest
```

### Kubernetes部署

参考 `k8s/` 目录下的配置文件

## 十、获取帮助

- 📖 [完整文档](./DEVELOPMENT_GUIDE.md)
- 🐛 [报告Bug](https://github.com/your-repo/issues)
- 💬 [讨论](https://github.com/your-repo/discussions)

## 十一、许可证

MIT License

---

祝您开发愉快！如有问题，欢迎提出Issue或Pull Request。
