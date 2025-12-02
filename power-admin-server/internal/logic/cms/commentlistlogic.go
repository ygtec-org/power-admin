// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取评论列表
func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListResp, err error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, total, err := l.svcCtx.CmsCommentRepo.List(l.ctx, req.ContentId, int64(req.Page), int64(req.PageSize))
	if err != nil {
		l.Logger.Errorf("查询评论列表失败: %v", err)
		return nil, err
	}

	var data []types.CommentData
	for _, comment := range result {
		data = append(data, types.CommentData{
			Id:      comment.ID,
			Content: comment.Content,
			Status:  int(comment.Status),
		})
	}

	return &types.CommentListResp{
		Total: total,
		Data:  data,
	}, nil
}
