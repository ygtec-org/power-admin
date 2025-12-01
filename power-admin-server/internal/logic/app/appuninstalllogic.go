// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppUninstallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppUninstallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppUninstallLogic {
	return &AppUninstallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppUninstallLogic) AppUninstall(req *types.AppUninstallReq) (resp *types.AppUninstallResp, err error) {
	// 检查插件是否已安装
	if !l.svcCtx.PluginService.IsPluginInstalled(req.AppKey) {
		resp = &types.AppUninstallResp{
			Success: false,
			Message: "插件未安装",
		}
		return
	}

	// 删除插件目录
	if err := l.svcCtx.PluginService.UninstallPlugin(req.AppKey); err != nil {
		l.Errorf("Failed to uninstall plugin: %v", err)
		return nil, err
	}

	resp = &types.AppUninstallResp{
		Success: true,
		Message: "插件卸载成功",
	}
	return
}
