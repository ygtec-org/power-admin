// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"
	"os"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppInstallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppInstallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppInstallLogic {
	return &AppInstallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppInstallLogic) AppInstall(req *types.AppInstallReq) (resp *types.AppInstallResp, err error) {
	// 检查插件是否已经安装
	if l.svcCtx.PluginService.IsPluginInstalled(req.AppKey) {
		resp = &types.AppInstallResp{
			Success: false,
			Message: "插件已安装",
		}
		return
	}

	// 获取plugins目录路径（从PluginService中获取）
	// 创建插件目录
	pluginPath := l.svcCtx.PluginService.GetPluginPath(req.AppKey)
	if installErr := os.MkdirAll(pluginPath, 0755); installErr != nil {
		l.Errorf("Failed to create plugin directory: %v", installErr)
		resp = &types.AppInstallResp{
			Success: false,
			Message: "插件安装失败: 无法创建目录",
		}
		return
	}

	resp = &types.AppInstallResp{
		Success: true,
		Message: "插件安装成功",
	}
	return
}
