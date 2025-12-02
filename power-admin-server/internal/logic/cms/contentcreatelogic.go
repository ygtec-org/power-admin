// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ContentCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建内容
func NewContentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentCreateLogic {
	return &ContentCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentCreateLogic) ContentCreate(req *types.ContentCreateReq) error {
	// todo: add your logic here and delete this line

	return nil
}
