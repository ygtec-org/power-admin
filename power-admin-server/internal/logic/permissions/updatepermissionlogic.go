// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package permissions

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePermissionLogic {
	return &UpdatePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePermissionLogic) UpdatePermission(req *types.UpdatePermissionReq) (err error) {
	if req.Id == 0 {
		return nil
	}

	perm, err := l.svcCtx.PermissionRepo.GetByID(req.Id)
	if err != nil || perm == nil {
		return err
	}

	if req.Status > 0 {
		perm.Status = req.Status
	}

	err = l.svcCtx.PermissionRepo.Update(perm)
	if err != nil {
		return nil
	}

	return nil
}
