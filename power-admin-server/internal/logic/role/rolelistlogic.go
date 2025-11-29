// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleListLogic {
	return &RoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleListLogic) RoleList(req *types.RoleListReq) (resp *types.RoleListResp, err error) {
	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	roles, total, err := l.svcCtx.RoleRepo.List(offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	data := make([]types.RoleInfo, 0, len(roles))
	for _, r := range roles {
		data = append(data, types.RoleInfo{
			Id:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Status:      r.Status,
			CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp = &types.RoleListResp{
		Total: total,
		Data:  data,
	}
	return
}
