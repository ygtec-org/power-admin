// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除评论
func NewCommentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentDeleteLogic {
	return &CommentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentDeleteLogic) CommentDelete(req *types.CommentDeleteReq) error {
	err := l.svcCtx.CmsCommentRepo.Delete(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("删除评论失败: %v", err)
		return err
	}
	l.Logger.Infof("成功删除评论: id=%d", req.Id)
	return nil
}
