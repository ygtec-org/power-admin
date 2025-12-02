// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分类列表
func NewCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryListLogic {
	return &CategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryListLogic) CategoryList(req *types.CategoryListReq) (resp *types.CategoryListResp, err error) {
	categories, err := l.svcCtx.CmsCategoryRepo.List(l.ctx, &req.ParentId)
	if err != nil {
		l.Logger.Errorf("查询分类列表失败: %v", err)
		return nil, err
	}

	var data []types.CategoryData
	for _, category := range categories {
		item := types.CategoryData{
			Id:          category.ID,
			Name:        category.Name,
			Slug:        category.Slug,
			Description: category.Description,
		}
		if category.ParentID != nil {
			item.ParentId = *category.ParentID
		}
		data = append(data, item)
	}

	return &types.CategoryListResp{
		Total: int64(len(data)),
		Data:  data,
	}, nil
}
