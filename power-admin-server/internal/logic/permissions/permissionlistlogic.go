// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package permissions

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionListLogic {
	return &PermissionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionListLogic) PermissionList(req *types.PermissionListReq) (resp *types.PermissionListResp, err error) {
	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	perms, total, err := l.svcCtx.PermissionRepo.List(offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	data := make([]types.PermissionInfo, 0, len(perms))
	for _, p := range perms {
		data = append(data, types.PermissionInfo{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Resource:    p.Resource,
			Action:      p.Action,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp = &types.PermissionListResp{
		Total: total,
		Data:  data,
	}
	return
}
