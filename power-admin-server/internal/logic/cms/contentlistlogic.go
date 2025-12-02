// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cms

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ContentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取内容列表
func NewContentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentListLogic {
	return &ContentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentListLogic) ContentList(req *types.ContentListReq) (resp *types.ContentListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
