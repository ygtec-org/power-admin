// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package permissions

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePermissionLogic {
	return &DeletePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePermissionLogic) DeletePermission(req *types.DeletePermissionReq) (err error) {
	if req.Id == 0 {
		return nil
	}

	err = l.svcCtx.PermissionRepo.Delete(req.Id)
	if err != nil {
		return nil
	}

	return nil
}
