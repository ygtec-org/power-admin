// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAppLogic {
	return &SearchAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchAppLogic) SearchApp(req *types.AppListReq) (resp *types.AppListResp, err error) {
	// 搜索应用
	keyword := req.Category // 此处使用Category作为keyword（实际应用中需要优化）
	apps, total, err := l.svcCtx.AppRepository.SearchApps(keyword, req.Page, req.PageSize)
	if err != nil {
		l.Errorf("Failed to search apps: %v", err)
		return nil, err
	}

	// 获取已安装插件列表（通过检测 plugins 目录）
	installedPlugins, err := l.svcCtx.PluginService.GetInstalledPlugins()
	if err != nil {
		l.Errorf("Failed to get installed plugins: %v", err)
		// 非致命错误，继续处理
		installedPlugins = []string{}
	}

	// 创建已安装插件的Map
	installedMap := make(map[string]bool)
	for _, pluginKey := range installedPlugins {
		installedMap[pluginKey] = true
	}

	// 构造响应
	appInfos := make([]types.AppInfo, len(apps))
	for i, app := range apps {
		appInfos[i] = types.AppInfo{
			Id:          app.ID,
			AppKey:      app.AppKey,
			AppName:     app.AppName,
			Version:     app.Version,
			Author:      app.Author,
			Description: app.Description,
			Icon:        app.Icon,
			DownloadUrl: app.DownloadUrl,
			DemoUrl:     app.DemoUrl,
			Category:    app.Category,
			Tags:        app.Tags,
			Rating:      app.Rating,
			Downloads:   app.Downloads,
			Status:      int(app.Status),
			Published:   int(app.Published),
			Installed:   installedMap[app.AppKey],
			CreatedAt:   app.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	resp = &types.AppListResp{
		Total: total,
		List:  appInfos,
	}
	return
}
