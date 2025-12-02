// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"
	"time"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishScheduleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 定时发布
func NewPublishScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishScheduleLogic {
	return &PublishScheduleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishScheduleLogic) PublishSchedule(req *types.PublishScheduleReq) error {
	scheduleAt, err := time.Parse("2006-01-02 15:04:05", req.ScheduleAt)
	if err != nil {
		l.Logger.Errorf("解析定时发布时间失败: %v", err)
		return err
	}
	content := &models.CmsContent{
		Status:      1,
		ScheduledAt: &scheduleAt,
	}
	err = l.svcCtx.CmsContentRepo.Update(l.ctx, req.Id, content)
	if err != nil {
		l.Logger.Errorf("定时发布失败: %v", err)
		return err
	}
	l.Logger.Infof("成功定时发布: id=%d", req.Id)
	return nil
}
