// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package permissions

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePermissionLogic {
	return &CreatePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePermissionLogic) CreatePermission(req *types.CreatePermissionReq) (err error) {

	perm := &models.Permission{
		Name:        req.Name,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		Status:      1,
	}

	err = l.svcCtx.PermissionRepo.Create(perm)
	if err != nil {
		return err
	}

	return nil
}
