// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

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

	if len(columns) > 0 {
		if err := l.svcCtx.CodegenRepo.CreateColumns(l.ctx, columns); err != nil {
			logx.Errorf("创建字段配置失败: %v", err)
			return nil, err
		}
	}

	// 4. 返回响应
	return l.configToResp(config, columns), nil
}

func (l *CreateConfigLogic) configToResp(config *models.GenConfig, columns []*models.GenTableColumn) *types.GenConfigResp {
	resp := &types.GenConfigResp{
		ID:           config.ID,
		TableName:    config.Table,
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
