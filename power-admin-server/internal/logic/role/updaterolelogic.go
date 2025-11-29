// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleReq) (err error) {

	role, err := l.svcCtx.RoleRepo.GetByID(req.Id)
	if err != nil || role == nil {
		return err
	}

	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status > 0 {
		role.Status = req.Status
	}

	err = l.svcCtx.RoleRepo.Update(role)
	if err != nil {
		return err
	}

	return nil
}
