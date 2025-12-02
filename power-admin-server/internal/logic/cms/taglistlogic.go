// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取标签列表
func NewTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagListLogic {
	return &TagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagListLogic) TagList(req *types.TagListReq) (resp *types.TagListResp, err error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, total, err := l.svcCtx.CmsTagRepo.List(l.ctx, req.Page, req.PageSize)
	if err != nil {
		l.Logger.Errorf("查询标签列表失败: %v", err)
		return nil, err
	}

	var data []types.TagData
	for _, tag := range result {
		data = append(data, types.TagData{
			Id:          tag.ID,
			Name:        tag.Name,
			Slug:        tag.Slug,
			Description: tag.Description,
			Color:       tag.Color,
		})
	}

	return &types.TagListResp{
		Total: total,
		Data:  data,
	}, nil
}
