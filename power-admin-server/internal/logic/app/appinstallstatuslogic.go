// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppInstallStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppInstallStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppInstallStatusLogic {
	return &AppInstallStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppInstallStatusLogic) AppInstallStatus(req *types.AppInstallStatusReq) (resp *types.AppInstallStatusResp, err error) {
	// 检查插件是否已安装（通过检测 plugins 目录）
	installed := l.svcCtx.PluginService.IsPluginInstalled(req.AppKey)

	resp = &types.AppInstallStatusResp{
		AppKey:    req.AppKey,
		Installed: installed,
	}
	return
}
