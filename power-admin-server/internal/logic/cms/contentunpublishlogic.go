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

type ContentUnpublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消发布
func NewContentUnpublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentUnpublishLogic {
	return &ContentUnpublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentUnpublishLogic) ContentUnpublish(req *types.IdReq) error {
	content := &models.CmsContent{Status: 0}
	err := l.svcCtx.CmsContentRepo.Update(l.ctx, req.Id, content)
	if err != nil {
		l.Logger.Errorf("取消发布内容失败: %v", err)
		return err
	}
	l.Logger.Infof("成功取消发布内容: id=%d", req.Id)
	return nil
}
