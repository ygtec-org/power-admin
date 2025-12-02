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

type UserEnableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 启用用户
func NewUserEnableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserEnableLogic {
	return &UserEnableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserEnableLogic) UserEnable(req *types.UserEnableReq) error {
	user := &models.CmsUser{Status: 1}
	err := l.svcCtx.CmsUserRepo.Update(l.ctx, req.Id, user)
	if err != nil {
		l.Logger.Errorf("启用用户失败: %v", err)
		return err
	}
	l.Logger.Infof("成功启用用户: id=%d", req.Id)
	return nil
}
