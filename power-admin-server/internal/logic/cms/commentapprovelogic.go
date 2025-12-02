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

type CommentApproveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 审核评论
func NewCommentApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentApproveLogic {
	return &CommentApproveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentApproveLogic) CommentApprove(req *types.CommentApproveReq) error {
	comment := &models.CmsComment{Status: 1}
	err := l.svcCtx.CmsCommentRepo.Update(l.ctx, req.Id, comment)
	if err != nil {
		l.Logger.Errorf("批准评论失败: %v", err)
		return err
	}
	l.Logger.Infof("成功批准评论: id=%d", req.Id)
	return nil
}
