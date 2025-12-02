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

type CategoryTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分类树
func NewCategoryTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryTreeLogic {
	return &CategoryTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryTreeLogic) CategoryTree() (resp *types.CategoryResp, err error) {
	categories, err := l.svcCtx.CmsCategoryRepo.GetAll(l.ctx)
	if err != nil {
		l.Logger.Errorf("获取分类树失败: %v", err)
		return nil, err
	}

	var rootCategories []*models.CmsCategory
	for i := range categories {
		if categories[i].ParentID == nil || *categories[i].ParentID == 0 {
			rootCategories = append(rootCategories, categories[i])
		}
	}

	categoryMap := make(map[int64]*models.CmsCategory)
	for i := range categories {
		categoryMap[categories[i].ID] = categories[i]
	}

	for _, category := range rootCategories {
		category.Children = buildCategoryTree(category, categoryMap)
	}

	var data []types.CategoryData
	for _, category := range rootCategories {
		data = append(data, l.convertToDTO(category))
	}

	return &types.CategoryResp{Data: data}, nil
}

func buildCategoryTree(parent *models.CmsCategory, categoryMap map[int64]*models.CmsCategory) []*models.CmsCategory {
	var children []*models.CmsCategory
	for _, category := range categoryMap {
		if category.ParentID != nil && *category.ParentID == parent.ID {
			category.Children = buildCategoryTree(category, categoryMap)
			children = append(children, category)
		}
	}
	return children
}

func (l *CategoryTreeLogic) convertToDTO(cat *models.CmsCategory) types.CategoryData {
	data := types.CategoryData{
		Id:          cat.ID,
		Name:        cat.Name,
		Slug:        cat.Slug,
		Description: cat.Description,
	}
	if cat.ParentID != nil {
		data.ParentId = *cat.ParentID
	}
	if len(cat.Children) > 0 {
		for _, child := range cat.Children {
			data.Children = append(data.Children, l.convertToDTO(child))
		}
	}
	return data
}
