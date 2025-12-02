// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

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
