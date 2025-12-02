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

type PublishBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量发布
func NewPublishBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishBatchLogic {
	return &PublishBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishBatchLogic) PublishBatch(req *types.PublishBatchReq) error {
	if len(req.Ids) == 0 {
		l.Logger.Infof("批量发布ID列表为空")
		return nil
	}

	for _, id := range req.Ids {
		content := &models.CmsContent{Status: 1}
		err := l.svcCtx.CmsContentRepo.Update(l.ctx, id, content)
		if err != nil {
			l.Logger.Errorf("批量发布内容失败: id=%d, err=%v", id, err)
		}
	}
	l.Logger.Infof("批量发布完成: %d条", len(req.Ids))
	return nil
}
