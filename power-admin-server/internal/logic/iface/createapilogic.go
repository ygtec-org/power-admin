// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iface

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApiLogic {
	return &CreateApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateApiLogic) CreateApi(req *types.CreateApiReq) (err error) {

	api := &models.API{
		APIName:     req.ApiName,
		APIPath:     req.ApiPath,
		APIMethod:   req.ApiMethod,
		Description: req.Description,
		Status:      1,
	}

	err = l.svcCtx.APIRepo.Create(api)
	if err != nil {
		return nil
	}

	return nil
}
