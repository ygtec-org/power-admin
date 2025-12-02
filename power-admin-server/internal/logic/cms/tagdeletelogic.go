// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除标签
func NewTagDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagDeleteLogic {
	return &TagDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagDeleteLogic) TagDelete(req *types.TagDeleteReq) error {
	err := l.svcCtx.CmsTagRepo.Delete(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("删除标签失败: %v", err)
		return err
	}
	l.Logger.Infof("成功删除标签: id=%d", req.Id)
	return nil
}
