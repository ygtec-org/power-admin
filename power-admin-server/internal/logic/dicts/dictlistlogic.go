// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dicts

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictListLogic {
	return &DictListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictListLogic) DictList(req *types.DictListReq) (resp *types.DictListResp, err error) {
	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	dicts, total, err := l.svcCtx.DictRepo.List(offset, req.PageSize, req.Dict)
	if err != nil {
		return nil, err
	}

	data := make([]types.DictInfo, 0, len(dicts))
	for _, d := range dicts {
		data = append(data, types.DictInfo{
			Id:          d.ID,
			DictType:    d.DictType,
			DictLabel:   d.DictLabel,
			DictValue:   d.DictValue,
			Description: d.Description,
			Sort:        d.Sort,
			Status:      d.Status,
		})
	}

	resp = &types.DictListResp{
		Total: total,
		Data:  data,
	}
	return
}
