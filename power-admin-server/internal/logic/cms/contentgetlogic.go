// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ContentGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取内容详情
func NewContentGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentGetLogic {
	return &ContentGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentGetLogic) ContentGet(req *types.IdReq) (resp *types.ContentResp, err error) {
	content, err := l.svcCtx.CmsContentRepo.Get(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("获取内容详情失败: %v", err)
		return nil, err
	}
	if content == nil {
		l.Logger.Infof("内容不存在: id=%d", req.Id)
		return nil, nil
	}

	resp = &types.ContentResp{
		Data: types.ContentData{
			Id:          content.ID,
			Title:       content.Title,
			Description: content.Description,
			Content:     content.Content,
			Slug:        content.Slug,
			Status:      int(content.Status),
			ViewCount:   int64(content.ViewCount),
			CreatedAt:   content.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   content.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	// 获取分类名称
	if content.CategoryID != nil {
		resp.Data.CategoryId = *content.CategoryID
		category, err := l.svcCtx.CmsCategoryRepo.Get(l.ctx, *content.CategoryID)
		if err == nil && category != nil {
			resp.Data.CategoryName = category.Name
		}
	}

	return
}
