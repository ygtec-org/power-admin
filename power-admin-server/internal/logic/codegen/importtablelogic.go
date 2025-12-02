// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

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
	columns, err := l.svcCtx.CodegenRepo.GetTableColumns(l.ctx, req.TableName)
	if err != nil || len(columns) == 0 {
		return nil, fmt.Errorf("表 %s 不存在或没有字段", req.TableName)
	}

	existing, _ := l.svcCtx.CodegenRepo.GetConfigByTableName(l.ctx, req.TableName)
	if existing != nil {
		return nil, fmt.Errorf("表 %s 的配置已存在", req.TableName)
	}

	config := &models.GenConfig{
		Table:        req.TableName,
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
