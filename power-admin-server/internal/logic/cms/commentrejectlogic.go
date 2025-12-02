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

type CommentRejectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 拒绝评论
func NewCommentRejectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentRejectLogic {
	return &CommentRejectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentRejectLogic) CommentReject(req *types.CommentRejectReq) error {
	comment := &models.CmsComment{Status: 2}
	err := l.svcCtx.CmsCommentRepo.Update(l.ctx, req.Id, comment)
	if err != nil {
		l.Logger.Errorf("拒绝评论失败: %v", err)
		return err
	}
	l.Logger.Infof("成功拒绝评论: id=%d", req.Id)
	return nil
}
