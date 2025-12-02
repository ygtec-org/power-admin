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
	id := l.ctx.Value("id").(int64)

	config, err := l.svcCtx.CodegenRepo.GetConfig(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("配置不存在")
	}

	config.Table = req.TableName
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

	if err := l.svcCtx.CodegenRepo.DeleteColumnsByConfigID(l.ctx, id); err != nil {
		logx.Errorf("删除旧字段失败: %v", err)
	}

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
