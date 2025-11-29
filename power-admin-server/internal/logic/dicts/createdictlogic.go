// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dicts

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictLogic {
	return &CreateDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDictLogic) CreateDict(req *types.CreateDictReq) (err error) {

	dict := &models.Dictionary{
		DictType:    req.DictType,
		DictLabel:   req.DictLabel,
		DictValue:   req.DictValue,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      1,
	}

	err = l.svcCtx.DictRepo.Create(dict)
	if err != nil {
		return nil
	}

	return nil
}
