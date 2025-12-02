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

type PublishImmediateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 立即发布
func NewPublishImmediateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishImmediateLogic {
	return &PublishImmediateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishImmediateLogic) PublishImmediate(req *types.PublishImmediateReq) error {
	content := &models.CmsContent{Status: 1}
	err := l.svcCtx.CmsContentRepo.Update(l.ctx, req.Id, content)
	if err != nil {
		l.Logger.Errorf("立即发布内容失败: %v", err)
		return err
	}
	l.Logger.Infof("成功立即发布: id=%d", req.Id)
	return nil
}
