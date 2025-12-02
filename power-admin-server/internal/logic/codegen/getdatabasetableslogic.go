// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

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
