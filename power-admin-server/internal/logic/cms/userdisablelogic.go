// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDisableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 禁用用户
func NewUserDisableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDisableLogic {
	return &UserDisableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDisableLogic) UserDisable(req *types.UserDisableReq) error {
	user := &models.CmsUser{Status: 0}
	err := l.svcCtx.CmsUserRepo.Update(l.ctx, req.Id, user)
	if err != nil {
		l.Logger.Errorf("禁用用户失败: %v", err)
		return err
	}
	l.Logger.Infof("成功禁用用户: id=%d", req.Id)
	return nil
}
