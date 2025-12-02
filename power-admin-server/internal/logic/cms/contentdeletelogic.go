// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ContentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除内容
func NewContentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentDeleteLogic {
	return &ContentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentDeleteLogic) ContentDelete(req *types.IdReq) error {
	err := l.svcCtx.CmsContentRepo.Delete(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("删除内容失败: %v", err)
		return err
	}
	l.Logger.Infof("成功删除内容: id=%d", req.Id)
	return nil
}
