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

type ContentUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新内容
func NewContentUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentUpdateLogic {
	return &ContentUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentUpdateLogic) ContentUpdate(req *types.ContentUpdateReq) error {
	content := &models.CmsContent{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Slug:        req.Slug,
		Status:      int8(req.Status),
		UpdatedAt:   time.Now(),
	}
	if req.CategoryId > 0 {
		content.CategoryID = &req.CategoryId
	}
	err := l.svcCtx.CmsContentRepo.Update(l.ctx, req.Id, content)
	if err != nil {
		l.Logger.Errorf("更新内容失败: %v", err)
		return err
	}
	l.Logger.Infof("成功更新内容: id=%d", req.Id)
	return nil
}
