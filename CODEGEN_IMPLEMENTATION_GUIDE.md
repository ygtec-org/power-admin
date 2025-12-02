# 代码生成器 Logic 层完整实现指南

## ✅ 已完成的工作

### 1. 数据库设计
- ✅ `db/codegen.sql` - 包含3张表和菜单数据

### 2. Model层
- ✅ `pkg/models/codegen_models.go` - GenConfig、GenTableColumn、GenHistory

### 3. Repository层
- ✅ `pkg/repository/codegen_repository.go` - 完整的数据访问层实现

### 4. 模板引擎
- ✅ `pkg/codegen/template_engine.go` - 代码生成模板引擎

### 5. ServiceContext
- ✅ 已注入 CodegenRepo
- ✅ 已添加模型自动迁移

### 6. Handler层
- ✅ goctl已生成12个Handler文件

### 7. Logic层骨架
- ✅ goctl已生成12个Logic文件(需要实现)

## 🔨 待实现：Logic 层业务逻辑

根据开发规范，需要删除Logic文件中的 `// todo:` 注释并实现完整的业务逻辑。

### Logic文件列表

1. `createconfiglogic.go` - 创建配置
2. `updateconfiglogic.go` - 更新配置
3. `deleteconfiglogic.go` - 删除配置
4. `getconfiglogic.go` - 获取配置详情
5. `listconfiglogic.go` - 配置列表
6. `generatecodelogic.go` - **核心**：生成代码
7. `previewcodelogic.go` - 预览代码
8. `listhistorylogic.go` - 历史列表
9. `gethistorylogic.go` - 历史详情
10. `deletehistorylogic.go` - 删除历史
11. `getdatabasetableslogic.go` - 获取数据库表列表
12. `importtablelogic.go` - 导入表结构

## 📝 Logic实现模板

### 1. createconfiglogic.go

```go
package codegen

import (
	"context"
	"power-admin/internal/svc"
	"power-admin/internal/types"
	"power-admin/pkg/models"
	"power-admin/pkg/repository"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateConfigLogic {
	return &CreateConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateConfigLogic) CreateConfig(req *types.GenConfigReq) (resp *types.GenConfigResp, err error) {
	// 1. 检查表名是否已存在
	existing, _ := l.svcCtx.CodegenRepo.GetConfigByTableName(l.ctx, req.TableName)
	if existing != nil {
		return nil, fmt.Errorf("表 %s 的配置已存在", req.TableName)
	}

	// 2. 创建配置
	config := &models.GenConfig{
		TableName:    req.TableName,
		TablePrefix:  req.TablePrefix,
		BusinessName: req.BusinessName,
		ModuleName:   req.ModuleName,
		PackageName:  req.PackageName,
		Author:       req.Author,
		Remark:       req.Remark,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := l.svcCtx.CodegenRepo.CreateConfig(l.ctx, config); err != nil {
		logx.Errorf("创建代码生成配置失败: %v", err)
		return nil, err
	}

	// 3. 创建字段配置
	columns := make([]*models.GenTableColumn, 0, len(req.Columns))
	for _, col := range req.Columns {
		column := &models.GenTableColumn{
			GenConfigID:   config.ID,
			ColumnName:    col.ColumnName,
			ColumnComment: col.ColumnComment,
			ColumnType:    col.ColumnType,
			GoType:        col.GoType,
			GoField:       col.GoField,
			IsPk:          col.IsPk,
			IsIncrement:   col.IsIncrement,
			IsRequired:    col.IsRequired,
			IsInsert:      col.IsInsert,
			IsEdit:        col.IsEdit,
			IsList:        col.IsList,
			IsQuery:       col.IsQuery,
			QueryType:     col.QueryType,
			HtmlType:      col.HtmlType,
			DictType:      col.DictType,
			Sort:          col.Sort,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		columns = append(columns, column)
	}

	if err := l.svcCtx.CodegenRepo.CreateColumns(l.ctx, columns); err != nil {
		logx.Errorf("创建字段配置失败: %v", err)
		return nil, err
	}

	// 4. 返回响应
	return l.configToResp(config, columns), nil
}

func (l *CreateConfigLogic) configToResp(config *models.GenConfig, columns []*models.GenTableColumn) *types.GenConfigResp {
	resp := &types.GenConfigResp{
		ID:           config.ID,
		TableName:    config.TableName,
		TablePrefix:  config.TablePrefix,
		BusinessName: config.BusinessName,
		ModuleName:   config.ModuleName,
		PackageName:  config.PackageName,
		Author:       config.Author,
		Remark:       config.Remark,
		CreatedAt:    config.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    config.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if len(columns) > 0 {
		resp.Columns = make([]types.GenTableColumnResp, 0, len(columns))
		for _, col := range columns {
			resp.Columns = append(resp.Columns, types.GenTableColumnResp{
				ID:            col.ID,
				GenConfigID:   col.GenConfigID,
				ColumnName:    col.ColumnName,
				ColumnComment: col.ColumnComment,
				ColumnType:    col.ColumnType,
				GoType:        col.GoType,
				GoField:       col.GoField,
				IsPk:          col.IsPk,
				IsIncrement:   col.IsIncrement,
				IsRequired:    col.IsRequired,
				IsInsert:      col.IsInsert,
				IsEdit:        col.IsEdit,
				IsList:        col.IsList,
				IsQuery:       col.IsQuery,
				QueryType:     col.QueryType,
				HtmlType:      col.HtmlType,
				DictType:      col.DictType,
				Sort:          col.Sort,
			})
		}
	}

	return resp
}
```

### 2. generatecodelogic.go (核心)

```go
package codegen

import (
	"context"
	"fmt"
	"power-admin/internal/svc"
	"power-admin/internal/types"
	"power-admin/pkg/codegen"
	"power-admin/pkg/models"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateCodeLogic {
	return &GenerateCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateCodeLogic) GenerateCode(req *types.CodeGenerateReq) (resp *types.CodeGenerateResp, err error) {
	// 1. 获取配置
	config, err := l.svcCtx.CodegenRepo.GetConfig(l.ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("配置不存在")
	}

	// 2. 获取字段
	columns, err := l.svcCtx.CodegenRepo.GetColumnsByConfigID(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if len(config.Columns) == 0 {
		config.Columns = columns
	}

	// 3. 初始化模板引擎
	engine := codegen.NewTemplateEngine()

	// 4. 生成代码文件
	files := make([]types.GeneratedFile, 0)
	histories := make([]*models.GenHistory, 0)

	// 生成API文件
	apiContent, err := engine.RenderAPI(config, config.Columns)
	if err == nil {
		apiFile := types.GeneratedFile{
			FilePath: fmt.Sprintf("api/%s.api", config.ModuleName),
			FileType: "api",
			Content:  apiContent,
		}
		files = append(files, apiFile)

		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			TableName:   config.TableName,
			FilePath:    apiFile.FilePath,
			FileType:    "api",
			Content:     apiContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	// 生成Model文件
	modelContent, err := engine.RenderModel(config, config.Columns)
	if err == nil {
		modelFile := types.GeneratedFile{
			FilePath: fmt.Sprintf("pkg/models/%s_models.go", config.ModuleName),
			FileType: "model",
			Content:  modelContent,
		}
		files = append(files, modelFile)

		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			TableName:   config.TableName,
			FilePath:    modelFile.FilePath,
			FileType:    "model",
			Content:     modelContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	// 生成Logic文件
	logicContent, err := engine.RenderLogic(config, config.Columns)
	if err == nil {
		logicFile := types.GeneratedFile{
			FilePath: fmt.Sprintf("internal/logic/%s/", config.ModuleName),
			FileType: "logic",
			Content:  logicContent,
		}
		files = append(files, logicFile)

		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			TableName:   config.TableName,
			FilePath:    logicFile.FilePath,
			FileType:    "logic",
			Content:     logicContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	// 生成Repository文件
	repoContent, err := engine.RenderRepository(config, config.Columns)
	if err == nil {
		repoFile := types.GeneratedFile{
			FilePath: fmt.Sprintf("pkg/repository/%s_repository.go", config.ModuleName),
			FileType: "repository",
			Content:  repoContent,
		}
		files = append(files, repoFile)

		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			TableName:   config.TableName,
			FilePath:    repoFile.FilePath,
			FileType:    "repository",
			Content:     repoContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	// 5. 保存历史记录
	if err := l.svcCtx.CodegenRepo.CreateHistories(l.ctx, histories); err != nil {
		logx.Errorf("保存生成历史失败: %v", err)
	}

	// 6. 返回响应
	resp = &types.CodeGenerateResp{
		Success: true,
		Message: "代码生成成功",
		Files:   files,
	}

	return resp, nil
}
```

### 3. getdatabasetableslogic.go

```go
package codegen

import (
	"context"
	"power-admin/internal/svc"
	"power-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDatabaseTablesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDatabaseTablesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDatabaseTablesLogic {
	return &GetDatabaseTablesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDatabaseTablesLogic) GetDatabaseTables(req *types.GetDatabaseTablesReq) (resp *types.GetDatabaseTablesResp, err error) {
	// 获取数据库表列表
	tables, err := l.svcCtx.CodegenRepo.GetDatabaseTables(l.ctx, req.TableName)
	if err != nil {
		return nil, err
	}

	resp = &types.GetDatabaseTablesResp{
		Tables: make([]types.DatabaseTableInfo, 0, len(tables)),
	}

	for _, table := range tables {
		// 获取表字段信息
		columns, _ := l.svcCtx.CodegenRepo.GetTableColumns(l.ctx, table.TableName)

		columnInfos := make([]types.ColumnInfo, 0, len(columns))
		for _, col := range columns {
			columnInfos = append(columnInfos, types.ColumnInfo{
				ColumnName:    col.ColumnName,
				ColumnType:    col.ColumnType,
				DataType:      col.DataType,
				ColumnComment: col.ColumnComment,
				IsNullable:    col.IsNullable,
				ColumnKey:     col.ColumnKey,
				Extra:         col.Extra,
			})
		}

		resp.Tables = append(resp.Tables, types.DatabaseTableInfo{
			TableName:    table.TableName,
			TableComment: table.TableComment,
			Engine:       table.Engine,
			Charset:      table.TableCollation,
			Columns:      columnInfos,
		})
	}

	return resp, nil
}
```

## 📌 其他Logic实现要点

### UpdateConfigLogic
- 先删除旧的字段配置
- 再创建新的字段配置
- 更新config表

### DeleteConfigLogic
- 级联删除会自动处理字段和历史

### GetConfigLogic
- 使用Preload预加载Columns

### ListConfigLogic / ListHistoryLogic
- 分页查询
- 支持表名模糊搜索

### PreviewCodeLogic
- 与GenerateCode逻辑相同
- 但不保存历史记录

### ImportTableLogic
- 从数据库读取表结构
- 自动映射MySQL类型到Go类型
- 自动生成字段配置

## 🚀 下一步

1. **实现所有Logic** - 复制上述模板并根据具体需求调整
2. **注册路由** - 在routes.go中添加代码生成器路由
3. **初始化数据库** - 执行db/codegen.sql
4. **测试API** - 使用Postman测试各个接口
5. **开发前端页面** - Vue3页面开发

## 💡 提示

- 所有Logic都要删除`// todo:`注释
- 添加完整的错误处理
- 记录操作日志
- 返回友好的错误信息
