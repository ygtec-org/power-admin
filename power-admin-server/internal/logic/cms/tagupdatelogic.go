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

type TagUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新标签
func NewTagUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagUpdateLogic {
	return &TagUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagUpdateLogic) TagUpdate(req *types.TagUpdateReq) error {
	tag := &models.CmsTag{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Color:       req.Color,
		UpdatedAt:   time.Now(),
	}
	err := l.svcCtx.CmsTagRepo.Update(l.ctx, req.Id, tag)
	if err != nil {
		l.Logger.Errorf("更新标签失败: %v", err)
		return err
	}
	l.Logger.Infof("成功更新标签: id=%d", req.Id)
	return nil
}
