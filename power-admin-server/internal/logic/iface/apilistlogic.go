// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iface

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiListLogic {
	return &ApiListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiListLogic) ApiList(req *types.ApiListReq) (resp *types.ApiListResp, err error) {
	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	apis, total, err := l.svcCtx.APIRepo.List(offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	data := make([]types.ApiInfo, 0, len(apis))
	for _, a := range apis {
		data = append(data, types.ApiInfo{
			Id:          a.ID,
			ApiName:     a.APIName,
			ApiPath:     a.APIPath,
			ApiMethod:   a.APIMethod,
			Group:       a.Group,
			Description: a.Description,
			Status:      a.Status,
			CreatedAt:   a.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp = &types.ApiListResp{
		Total:    total,
		ApiInfos: data,
	}
	return
}
