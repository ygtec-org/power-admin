# 代码生成器Logic层完整实现代码

## 说明
将以下代码复制到对应的Logic文件中，替换todo部分。

## 1. updateconfiglogic.go
```go
package codegen

import (
	"context"
	"fmt"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/models"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConfigLogic {
	return &UpdateConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateConfigLogic) UpdateConfig(req *types.GenConfigReq) (resp *types.GenConfigResp, err error) {
	// 获取ID从路径参数
	id := l.ctx.Value("id").(int64)
	
	// 获取原配置
	config, err := l.svcCtx.CodegenRepo.GetConfig(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("配置不存在")
	}

	// 更新配置
	config.TableName = req.TableName
	config.TablePrefix = req.TablePrefix
	config.BusinessName = req.BusinessName
	config.ModuleName = req.ModuleName
	config.PackageName = req.PackageName
	config.Author = req.Author
	config.Remark = req.Remark
	config.UpdatedAt = time.Now()

	if err := l.svcCtx.CodegenRepo.UpdateConfig(l.ctx, config); err != nil {
		logx.Errorf("更新配置失败: %v", err)
		return nil, err
	}

	// 删除旧字段
	if err := l.svcCtx.CodegenRepo.DeleteColumnsByConfigID(l.ctx, id); err != nil {
		logx.Errorf("删除旧字段失败: %v", err)
	}

	// 创建新字段
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

	if len(columns) > 0 {
		if err := l.svcCtx.CodegenRepo.CreateColumns(l.ctx, columns); err != nil {
			logx.Errorf("创建字段配置失败: %v", err)
			return nil, err
		}
	}

	return NewCreateConfigLogic(l.ctx, l.svcCtx).configToResp(config, columns), nil
}
```

## 2. deleteconfiglogic.go
```go
package codegen

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"power-admin-server/internal/svc"
)

type DeleteConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConfigLogic {
	return &DeleteConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteConfigLogic) DeleteConfig() error {
	id := l.ctx.Value("id").(int64)
	return l.svcCtx.CodegenRepo.DeleteConfig(l.ctx, id)
}
```

## 3. getconfiglogic.go
```go
package codegen

import (
	"context"
	"fmt"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigLogic {
	return &GetConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConfigLogic) GetConfig() (resp *types.GenConfigResp, err error) {
	id := l.ctx.Value("id").(int64)
	
	config, err := l.svcCtx.CodegenRepo.GetConfig(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("配置不存在")
	}

	return NewCreateConfigLogic(l.ctx, l.svcCtx).configToResp(config, config.Columns), nil
}
```

## 4. listconfiglogic.go
```go
package codegen

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListConfigLogic {
	return &ListConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListConfigLogic) ListConfig(req *types.GenConfigListReq) (resp *types.GenConfigListResp, err error) {
	configs, total, err := l.svcCtx.CodegenRepo.ListConfig(l.ctx, req.Page, req.PageSize, req.TableName)
	if err != nil {
		return nil, err
	}

	resp = &types.GenConfigListResp{
		Total: total,
		Data:  make([]types.GenConfigResp, 0, len(configs)),
	}

	for _, config := range configs {
		resp.Data = append(resp.Data, types.GenConfigResp{
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
		})
	}

	return resp, nil
}
```

## 5. generatecodelogic.go (核心)
```go
package codegen

import (
	"context"
	"fmt"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/codegen"
	"power-admin-server/pkg/models"
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
	// 获取配置
	config, err := l.svcCtx.CodegenRepo.GetConfig(l.ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("配置不存在")
	}

	// 初始化模板引擎
	engine := codegen.NewTemplateEngine()

	// 生成代码文件
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

	// 保存历史记录
	if len(histories) > 0 {
		if err := l.svcCtx.CodegenRepo.CreateHistories(l.ctx, histories); err != nil {
			logx.Errorf("保存生成历史失败: %v", err)
		}
	}

	resp = &types.CodeGenerateResp{
		Success: true,
		Message: "代码生成成功",
		Files:   files,
	}

	return resp, nil
}
```

## 6. previewcodelogic.go
```go
package codegen

import (
	"context"
	"fmt"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/codegen"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreviewCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPreviewCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewCodeLogic {
	return &PreviewCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreviewCodeLogic) PreviewCode(req *types.CodePreviewReq) (resp *types.CodePreviewResp, err error) {
	config, err := l.svcCtx.CodegenRepo.GetConfig(l.ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("配置不存在")
	}

	engine := codegen.NewTemplateEngine()

	files := make([]types.GeneratedFile, 0)

	// 预览API文件
	if apiContent, err := engine.RenderAPI(config, config.Columns); err == nil {
		files = append(files, types.GeneratedFile{
			FilePath: fmt.Sprintf("api/%s.api", config.ModuleName),
			FileType: "api",
			Content:  apiContent,
		})
	}

	// 预览Model文件
	if modelContent, err := engine.RenderModel(config, config.Columns); err == nil {
		files = append(files, types.GeneratedFile{
			FilePath: fmt.Sprintf("pkg/models/%s_models.go", config.ModuleName),
			FileType: "model",
			Content:  modelContent,
		})
	}

	// 预览Logic文件
	if logicContent, err := engine.RenderLogic(config, config.Columns); err == nil {
		files = append(files, types.GeneratedFile{
			FilePath: fmt.Sprintf("internal/logic/%s/", config.ModuleName),
			FileType: "logic",
			Content:  logicContent,
		})
	}

	// 预览Repository文件
	if repoContent, err := engine.RenderRepository(config, config.Columns); err == nil {
		files = append(files, types.GeneratedFile{
			FilePath: fmt.Sprintf("pkg/repository/%s_repository.go", config.ModuleName),
			FileType: "repository",
			Content:  repoContent,
		})
	}

	return &types.CodePreviewResp{
		Files: files,
	}, nil
}
```

## 7. listhistorylogic.go
```go
package codegen

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListHistoryLogic {
	return &ListHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListHistoryLogic) ListHistory(req *types.GenHistoryListReq) (resp *types.GenHistoryListResp, err error) {
	histories, total, err := l.svcCtx.CodegenRepo.ListHistory(l.ctx, req.Page, req.PageSize, req.TableName)
	if err != nil {
		return nil, err
	}

	resp = &types.GenHistoryListResp{
		Total: total,
		Data:  make([]types.GenHistoryResp, 0, len(histories)),
	}

	for _, history := range histories {
		resp.Data = append(resp.Data, types.GenHistoryResp{
			ID:          history.ID,
			GenConfigID: history.GenConfigID,
			TableName:   history.TableName,
			FilePath:    history.FilePath,
			FileType:    history.FileType,
			Content:     history.Content,
			Status:      history.Status,
			ErrorMsg:    history.ErrorMsg,
			Operator:    history.Operator,
			CreatedAt:   history.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}
```

## 8. gethistorylogic.go
```go
package codegen

import (
	"context"
	"fmt"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHistoryLogic {
	return &GetHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHistoryLogic) GetHistory() (resp *types.GenHistoryResp, err error) {
	id := l.ctx.Value("id").(int64)
	
	history, err := l.svcCtx.CodegenRepo.GetHistory(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("历史记录不存在")
	}

	return &types.GenHistoryResp{
		ID:          history.ID,
		GenConfigID: history.GenConfigID,
		TableName:   history.TableName,
		FilePath:    history.FilePath,
		FileType:    history.FileType,
		Content:     history.Content,
		Status:      history.Status,
		ErrorMsg:    history.ErrorMsg,
		Operator:    history.Operator,
		CreatedAt:   history.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
```

## 9. deletehistorylogic.go
```go
package codegen

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"power-admin-server/internal/svc"
)

type DeleteHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteHistoryLogic {
	return &DeleteHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteHistoryLogic) DeleteHistory() error {
	id := l.ctx.Value("id").(int64)
	return l.svcCtx.CodegenRepo.DeleteHistory(l.ctx, id)
}
```

## 10. getdatabasetableslogic.go
```go
package codegen

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

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
	tables, err := l.svcCtx.CodegenRepo.GetDatabaseTables(l.ctx, req.TableName)
	if err != nil {
		return nil, err
	}

	resp = &types.GetDatabaseTablesResp{
		Tables: make([]types.DatabaseTableInfo, 0, len(tables)),
	}

	for _, table := range tables {
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

## 11. importtablelogic.go
```go
package codegen

import (
	"context"
	"fmt"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/models"
	"power-admin-server/pkg/repository"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportTableLogic {
	return &ImportTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportTableLogic) ImportTable(req *types.GenConfigReq) (resp *types.GenConfigResp, err error) {
	// 检查表是否存在
	columns, err := l.svcCtx.CodegenRepo.GetTableColumns(l.ctx, req.TableName)
	if err != nil || len(columns) == 0 {
		return nil, fmt.Errorf("表 %s 不存在或没有字段", req.TableName)
	}

	// 检查配置是否已存在
	existing, _ := l.svcCtx.CodegenRepo.GetConfigByTableName(l.ctx, req.TableName)
	if existing != nil {
		return nil, fmt.Errorf("表 %s 的配置已存在", req.TableName)
	}

	// 创建配置
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
		return nil, err
	}

	// 导入字段配置
	genColumns := make([]*models.GenTableColumn, 0, len(columns))
	for i, col := range columns {
		goType := repository.MySQLTypeToGoType(col.DataType)
		goField := repository.ColumnNameToGoField(col.ColumnName)
		
		genColumn := &models.GenTableColumn{
			GenConfigID:   config.ID,
			ColumnName:    col.ColumnName,
			ColumnComment: col.ColumnComment,
			ColumnType:    col.ColumnType,
			GoType:        goType,
			GoField:       goField,
			IsPk:          boolToInt(col.ColumnKey == "PRI"),
			IsIncrement:   boolToInt(col.Extra == "auto_increment"),
			IsRequired:    boolToInt(col.IsNullable == "NO"),
			IsInsert:      1,
			IsEdit:        1,
			IsList:        1,
			IsQuery:       1,
			QueryType:     "=",
			HtmlType:      repository.GetHtmlType(goType),
			Sort:          i + 1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		genColumns = append(genColumns, genColumn)
	}

	if err := l.svcCtx.CodegenRepo.CreateColumns(l.ctx, genColumns); err != nil {
		return nil, err
	}

	return NewCreateConfigLogic(l.ctx, l.svcCtx).configToResp(config, genColumns), nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
```
