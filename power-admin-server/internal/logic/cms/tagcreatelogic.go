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

type TagCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建标签
func NewTagCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagCreateLogic {
	return &TagCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagCreateLogic) TagCreate(req *types.TagCreateReq) error {
	tag := &models.CmsTag{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Color:       req.Color,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := l.svcCtx.CmsTagRepo.Create(l.ctx, tag)
	if err != nil {
		l.Logger.Errorf("创建标签失败: %v", err)
		return err
	}
	l.Logger.Infof("成功创建标签: id=%d", tag.ID)
	return nil
}
