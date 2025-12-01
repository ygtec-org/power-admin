package service

import (
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"
)

// PluginService 处理插件相关的业务逻辑
type PluginService struct {
	pluginsDir string
}

// NewPluginService 创建插件服务
func NewPluginService(pluginsDir string) *PluginService {
	return &PluginService{
		pluginsDir: pluginsDir,
	}
}

// IsPluginInstalled 检查插件是否已安装（通过检测目录是否存在）
func (s *PluginService) IsPluginInstalled(pluginKey string) bool {
	pluginPath := filepath.Join(s.pluginsDir, pluginKey)
	info, err := os.Stat(pluginPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		logx.Errorf("Failed to check plugin path: %v", err)
		return false
	}
	return info.IsDir()
}

// GetInstalledPlugins 获取所有已安装的插件列表
func (s *PluginService) GetInstalledPlugins() ([]string, error) {
	entries, err := os.ReadDir(s.pluginsDir)
	if err != nil {
		return nil, err
	}

	var plugins []string
	for _, entry := range entries {
		// 跳过 common 目录，只统计目录
		if entry.IsDir() && entry.Name() != "common" {
			plugins = append(plugins, entry.Name())
		}
	}
	return plugins, nil
}

// InstallPlugin 安装插件（通常是将插件目录放到 plugins 目录下）
// 这里仅为占位，实际安装逻辑由前端或其他服务处理
func (s *PluginService) InstallPlugin(pluginKey string) error {
	// 实际安装逻辑需要由部署系统处理
	return nil
}

// UninstallPlugin 卸载插件（删除插件目录）
func (s *PluginService) UninstallPlugin(pluginKey string) error {
	pluginPath := filepath.Join(s.pluginsDir, pluginKey)
	return os.RemoveAll(pluginPath)
}

// GetPluginPath 获取插件目录路径
func (s *PluginService) GetPluginPath(pluginKey string) string {
	return filepath.Join(s.pluginsDir, pluginKey)
}
