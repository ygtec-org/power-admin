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

type CategoryCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建分类
func NewCategoryCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryCreateLogic {
	return &CategoryCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryCreateLogic) CategoryCreate(req *types.CategoryCreateReq) error {
	category := &models.CmsCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if req.ParentId > 0 {
		category.ParentID = &req.ParentId
	}
	err := l.svcCtx.CmsCategoryRepo.Create(l.ctx, category)
	if err != nil {
		l.Logger.Errorf("创建分类失败: %v", err)
		return err
	}
	l.Logger.Infof("成功创建分类: id=%d", category.ID)
	return nil
}
