// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

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
			TableName:    config.Table,
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
