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

type CategoryUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新分类
func NewCategoryUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryUpdateLogic {
	return &CategoryUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryUpdateLogic) CategoryUpdate(req *types.CategoryUpdateReq) error {
	category := &models.CmsCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		UpdatedAt:   time.Now(),
	}
	if req.ParentId > 0 {
		category.ParentID = &req.ParentId
	}
	err := l.svcCtx.CmsCategoryRepo.Update(l.ctx, req.Id, category)
	if err != nil {
		l.Logger.Errorf("更新分类失败: %v", err)
		return err
	}
	l.Logger.Infof("成功更新分类: id=%d", req.Id)
	return nil
}
