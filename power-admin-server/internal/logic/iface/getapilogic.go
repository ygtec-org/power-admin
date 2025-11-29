// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iface

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApiLogic {
	return &GetApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetApiLogic) GetApi(req *types.GetApiReq) (resp *types.ApiInfo, err error) {
	if req.Id == 0 {
		return nil, nil
	}

	api, err := l.svcCtx.APIRepo.GetByID(req.Id)
	if err != nil || api == nil {
		return nil, err
	}

	resp = &types.ApiInfo{
		Id:          api.ID,
		ApiName:     api.APIName,
		ApiPath:     api.APIPath,
		ApiMethod:   api.APIMethod,
		Description: api.Description,
		Status:      api.Status,
		CreatedAt:   api.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return
}
