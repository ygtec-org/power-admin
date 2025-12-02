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

type PublishCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消定时发布
func NewPublishCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishCancelLogic {
	return &PublishCancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishCancelLogic) PublishCancel(req *types.PublishCancelReq) error {
	content := &models.CmsContent{Status: 0, ScheduledAt: nil}
	err := l.svcCtx.CmsContentRepo.Update(l.ctx, req.Id, content)
	if err != nil {
		l.Logger.Errorf("取消定时发布失败: %v", err)
		return err
	}
	l.Logger.Infof("成功取消定时发布: id=%d", req.Id)
	return nil
}
