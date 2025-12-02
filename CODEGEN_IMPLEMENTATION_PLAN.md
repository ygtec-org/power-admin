# 代码生成器功能实现方案

## 📋 功能概述

代码生成器是一个完整的开发工具模块,用于根据数据库表结构自动生成 CRUD 代码,包括:
- API 定义文件 (.api)
- Model 模型文件
- Handler 处理层
- Logic 业务逻辑层  
- Repository 数据访问层

## 🎯 核心功能

### 1. 代码生成配置管理
- 创建/编辑/删除代码生成配置
- 配置表名、业务名、模块名等基本信息
- 配置字段映射关系(MySQL类型 → Go类型)
- 设置字段是否用于查询、列表、编辑等

### 2. 智能代码生成
- 根据配置自动生成完整的CRUD代码
- 支持代码预览
- 生成API文件(go-zero格式)
- 生成Model/Handler/Logic/Repository文件
- 自动处理表前缀

### 3. 生成历史记录
- 记录每次代码生成的详细信息
- 可查看历史生成的代码
- 支持重新生成

### 4. 数据库表导入
- 读取现有数据库表结构
- 自动解析字段信息
- 一键导入生成配置

## 🗂️ 数据库设计

已创建3张表:

### gen_config - 代码生成配置表
```sql
- id: 配置ID
- table_name: 表名称(唯一)
- table_prefix: 表前缀
- business_name: 业务名称
- module_name: 模块名称  
- package_name: 包路径
- author: 作者
- remark: 备注
```

### gen_table_column - 表字段信息表
```sql
- id: 字段ID
- gen_config_id: 配置ID(外键)
- column_name: 字段名称
- column_type: MySQL字段类型
- go_type: Go类型
- go_field: Go字段名
- is_pk/is_increment/is_required: 字段属性
- is_insert/is_edit/is_list/is_query: 使用场景
- query_type: 查询方式(=,LIKE等)
- html_type: 前端显示类型
- dict_type: 字典类型
- sort: 排序
```

### gen_history - 代码生成历史表
```sql
- id: 历史ID
- gen_config_id: 配置ID(外键)
- table_name: 表名称
- file_path: 生成的文件路径
- file_type: 文件类型(api/model/handler/logic/repository)
- content: 生成的文件内容
- status: 状态(1成功 0失败)
- error_msg: 错误信息
- operator: 操作人
- created_at: 生成时间
```

## 📁 文件结构

```
power-admin-server/
├── api/
│   └── codegen.api          # 代码生成器API定义(已创建)
├── db/
│   └── codegen.sql          # 数据库初始化脚本(已创建)
├── pkg/
│   ├── models/
│   │   └── codegen_models.go  # 代码生成模型(已创建)
│   └── repository/
│       └── codegen_repository.go  # 代码生成仓储(待创建)
├── internal/
│   ├── handler/
│   │   └── codegen/         # Handler层(goctl生成)
│   ├── logic/
│   │   └── codegen/         # Logic层(goctl生成后实现)
│   └── svc/
│       └── servicecontext.go  # 注入代码生成仓储
└── templates/               # 代码模板目录(待创建)
    ├── api.tpl
    ├── model.tpl
    ├── handler.tpl
    ├── logic.tpl
    └── repository.tpl
```

## 🔧 实现步骤

### 步骤1: 数据库初始化 ✅
```bash
# 执行SQL脚本
mysql -u root -p power_admin < db/codegen.sql
```

### 步骤2: 使用goctl生成基础代码
```bash
cd power-admin-server

# 生成Handler和Logic
goctl api go -api api/codegen.api -dir .
```

这将自动生成:
- `internal/handler/codegen/*.go` - 所有Handler
- `internal/logic/codegen/*.go` - 所有Logic(需要实现)
- `internal/types/types.go` - Request/Response类型定义

### 步骤3: 创建Repository层
创建 `pkg/repository/codegen_repository.go`:

**接口定义**:
```go
type CodegenRepository interface {
    // 配置管理
    CreateConfig(ctx context.Context, config *models.GenConfig) error
    UpdateConfig(ctx context.Context, config *models.GenConfig) error
    DeleteConfig(ctx context.Context, id int64) error
    GetConfig(ctx context.Context, id int64) (*models.GenConfig, error)
    ListConfig(ctx context.Context, page, pageSize int, tableName string) ([]*models.GenConfig, int64, error)
    
    // 历史管理
    CreateHistory(ctx context.Context, history *models.GenHistory) error
    ListHistory(ctx context.Context, page, pageSize int, tableName string) ([]*models.GenHistory, int64, error)
    GetHistory(ctx context.Context, id int64) (*models.GenHistory, error)
    DeleteHistory(ctx context.Context, id int64) error
    
    // 数据库表信息
    GetDatabaseTables(ctx context.Context, tableName string) ([]DatabaseTable, error)
    GetTableColumns(ctx context.Context, tableName string) ([]TableColumn, error)
}
```

### 步骤4: 实现Logic层

每个Logic需要实现的核心功能:

**CreateConfigLogic** - 创建配置
```go
func (l *CreateConfigLogic) CreateConfig(req *types.GenConfigReq) (*types.GenConfigResp, error) {
    // 1. 参数验证
    // 2. 转换为Model
    // 3. 调用Repository保存
    // 4. 转换为Response返回
}
```

**GenerateCodeLogic** - 生成代码(核心逻辑)
```go
func (l *GenerateCodeLogic) GenerateCode(req *types.CodeGenerateReq) (*types.CodeGenerateResp, error) {
    // 1. 获取配置信息
    // 2. 加载代码模板
    // 3. 渲染模板(替换表名、字段等)
    // 4. 生成文件
    // 5. 保存历史记录
    // 6. 返回生成结果
}
```

### 步骤5: 创建代码模板

在 `templates/` 目录创建模板文件:

**api.tpl** - API定义模板
```go
syntax = "v1"

type (
    {{.StructName}}Req {
        {{range .Fields}}
        {{.GoField}} {{.GoType}} `json:"{{.JsonField}}"`
        {{end}}
    }
    
    {{.StructName}}Resp {
        {{range .Fields}}
        {{.GoField}} {{.GoType}} `json:"{{.JsonField}}"`
        {{end}}
    }
)

@server(
    prefix: /api/admin/{{.ModuleName}}
)
service power-admin {
    @handler List{{.StructName}}
    get /{{.BusinessName}}/list returns ([]{{.StructName}}Resp)
    
    @handler Create{{.StructName}}
    post /{{.BusinessName}} ({{.StructName}}Req)
    
    // ... 其他CRUD操作
}
```

### 步骤6: 实现模板渲染引擎

创建 `pkg/codegen/template_engine.go`:

```go
type TemplateEngine struct {
    templates map[string]*template.Template
}

func (e *TemplateEngine) Render(tplName string, data interface{}) (string, error) {
    // 1. 加载模板
    // 2. 渲染数据
    // 3. 返回生成的代码
}
```

### 步骤7: 注册路由

在 `internal/handler/routes.go` 添加代码生成器路由:

```go
server.AddRoutes(
    rest.WithMiddlewares(
        []rest.Middleware{serverCtx.AdminAuthMiddleware},
        []rest.Route{
            {
                Method:  http.MethodPost,
                Path:    "/codegen/config",
                Handler: codegen.CreateConfigHandler(serverCtx),
            },
            // ... 其他路由
        }...,
    ),
    rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
    rest.WithPrefix("/api/admin"),
)
```

### 步骤8: 前端页面开发

创建Vue页面:

**power-admin-web/src/pages/devtools/codegen/CodeGen.vue**
- 配置表单(表名、字段信息)
- 字段配置表格
- 代码预览对话框
- 生成代码按钮

**power-admin-web/src/pages/devtools/codegen/GenHistory.vue**
- 历史记录列表
- 代码查看对话框
- 删除/重新生成操作

## 🎨 核心算法

### 1. MySQL类型 → Go类型映射

```go
func MySQLTypeToGoType(mysqlType string) string {
    typeMap := map[string]string{
        "int":        "int64",
        "bigint":     "int64",
        "tinyint":    "int",
        "varchar":    "string",
        "text":       "string",
        "datetime":   "time.Time",
        "timestamp":  "time.Time",
        "decimal":    "float64",
        // ...
    }
    return typeMap[mysqlType]
}
```

### 2. 表名 → 结构体名转换

```go
func TableNameToStructName(tableName, prefix string) string {
    // 1. 去除表前缀
    name := strings.TrimPrefix(tableName, prefix+"_")
    
    // 2. 下划线转驼峰
    parts := strings.Split(name, "_")
    for i, part := range parts {
        parts[i] = strings.Title(part)
    }
    
    return strings.Join(parts, "")
}

// 示例:
// admin_users + "admin" => Users
// cms_content + "cms" => Content
```

### 3. 字段名转换

```go
func ColumnNameToGoField(columnName string) string {
    // 下划线转大驼峰
    // user_id => UserID
    // created_at => CreatedAt
}

func ColumnNameToJsonField(columnName string) string {
    // 下划线转小驼峰
    // user_id => userId
    // created_at => createdAt
}
```

## 🔒 约束和规范

### 1. 结构体唯一性

**问题**: 同一个API文件中不能有重复的结构体名称

**解决方案**:
```go
// 为每个表生成独立的API文件
// 文件名: api/{module_name}.api

// 结构体命名规范:
// {TableName}Req
// {TableName}Resp
// {TableName}ListReq
// {TableName}ListResp

// 示例: users表生成
type (
    UserReq {
        Username string `json:"username"`
    }
    
    UserResp {
        ID int64 `json:"id"`
        Username string `json:"username"`
    }
    
    UserListReq {
        Page int `form:"page"`
    }
    
    UserListResp {
        Total int64 `json:"total"`
        Data []UserResp `json:"data"`
    }
)
```

### 2. 文件生成位置

```
api/
└── {module_name}.api        # 独立的API文件

pkg/models/
└── {module_name}_models.go  # Model文件

internal/
├── handler/{module_name}/   # Handler目录
├── logic/{module_name}/     # Logic目录(在此实现业务)
└── repository/
    └── {module_name}_repository.go
```

### 3. 命名规范

- **表名**: 小写+下划线 (admin_users)
- **Go结构体**: 大驼峰 (AdminUser)
- **Go字段**: 大驼峰 (UserName)
- **JSON字段**: 小驼峰 (userName)
- **API路径**: 小写+连字符 (/admin/users)

## 📊 使用示例

### 示例1: 生成用户表代码

**输入配置**:
```json
{
  "tableName": "admin_users",
  "tablePrefix": "admin",
  "businessName": "user",
  "moduleName": "admin",
  "packageName": "power-admin/internal",
  "author": "PowerAdmin",
  "columns": [
    {
      "columnName": "id",
      "columnType": "bigint",
      "goType": "int64",
      "goField": "ID",
      "isPk": 1,
      "isIncrement": 1
    },
    {
      "columnName": "username",
      "columnType": "varchar(100)",
      "goType": "string",
      "goField": "Username",
      "isRequired": 1
    }
  ]
}
```

**生成文件**:
1. `api/admin.api` - 包含User相关的API定义
2. `pkg/models/admin_models.go` - AdminUser结构体
3. `internal/handler/admin/user*.go` - Handler文件
4. `internal/logic/admin/user*.go` - Logic文件
5. `pkg/repository/admin_repository.go` - Repository接口

## ⚙️ 配置项

在 `etc/power-api.yaml` 添加代码生成配置:

```yaml
CodeGen:
  OutputPath: "."              # 代码输出根目录
  TemplatePath: "templates"    # 模板文件目录
  Author: "PowerAdmin"         # 默认作者
  PackageName: "power-admin"   # 默认包名
  EnableBackup: true           # 是否备份已存在文件
```

## 🚀 快速开始

### 1. 初始化数据库
```bash
mysql -u root -p power_admin < db/codegen.sql
```

### 2. 生成基础代码
```bash
goctl api go -api api/codegen.api -dir .
```

### 3. 实现Logic层
根据业务需求实现各个Logic方法

### 4. 创建前端页面
在Vue项目中创建代码生成器管理页面

### 5. 测试
- 导入数据库表
- 配置字段映射
- 预览代码
- 生成代码
- 查看历史

## 📝 待完成任务清单

- [ ] 创建Repository层实现
- [ ] 实现所有Logic业务逻辑
- [ ] 创建代码模板文件
- [ ] 实现模板渲染引擎
- [ ] 注册路由
- [ ] 创建前端页面
- [ ] 编写单元测试
- [ ] 编写使用文档

## 🎯 下一步工作

建议按以下顺序实现:

1. **Repository层** - 数据访问基础
2. **模板引擎** - 代码生成核心
3. **Logic层** - 业务逻辑实现
4. **前端页面** - 用户交互界面
5. **测试验证** - 确保功能正确

需要我继续实现哪个部分?
