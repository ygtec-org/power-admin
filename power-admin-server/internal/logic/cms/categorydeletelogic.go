// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除分类
func NewCategoryDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryDeleteLogic {
	return &CategoryDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryDeleteLogic) CategoryDelete(req *types.CategoryDeleteReq) error {
	err := l.svcCtx.CmsCategoryRepo.Delete(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("删除分类失败: %v", err)
		return err
	}
	l.Logger.Infof("成功删除分类: id=%d", req.Id)
	return nil
}
