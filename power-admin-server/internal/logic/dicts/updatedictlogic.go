// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dicts

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictLogic {
	return &UpdateDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDictLogic) UpdateDict(req *types.UpdateDictReq) (err error) {
	if req.Id == 0 {
		return nil
	}

	dict, err := l.svcCtx.DictRepo.GetByID(req.Id)
	if err != nil || dict == nil {
		return err
	}

	if req.DictLabel != "" {
		dict.DictLabel = req.DictLabel
	}
	if req.Sort > 0 {
		dict.Sort = req.Sort
	}
	if req.Status > 0 {
		dict.Status = req.Status
	}

	err = l.svcCtx.DictRepo.Update(dict)
	if err != nil {
		return nil
	}

	return nil
}
