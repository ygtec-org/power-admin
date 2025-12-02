package service

import (
	"os"
	"path/filepath"

	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// PluginService 处理插件相关的业务逻辑
type PluginService struct {
	pluginsDir string
	db         *gorm.DB
}

// NewPluginService 创建插件服务
func NewPluginService(pluginsDir string, db *gorm.DB) *PluginService {
	return &PluginService{
		pluginsDir: pluginsDir,
		db:         db,
	}
}

// IsPluginInstalled 检查插件是否已安装（双重验证：检查目录 + 数据库记录）
func (s *PluginService) IsPluginInstalled(pluginKey string) bool {
	// 第一步：检查文件系统中的目录是否存在
	pluginPath := filepath.Join(s.pluginsDir, pluginKey)
	info, err := os.Stat(pluginPath)
	if err != nil {
		if os.IsNotExist(err) {
			logx.Infof("Plugin directory not found for appKey: %s", pluginKey)
			return false
		}
		logx.Errorf("Failed to check plugin path: %v", err)
		return false
	}

	// 检查是否为目录
	if !info.IsDir() {
		logx.Infof("Plugin path exists but is not a directory: %s", pluginPath)
		return false
	}

	// 第二步：验证数据库中是否有对应的安装记录
	if s.db == nil {
		logx.Infof("Database connection not initialized, skipping database verification")
		return true // 仅依靠目录检查
	}

	var installation models.AppInstallation
	result := s.db.Where("app_key = ? AND status = ?", pluginKey, 1).First(&installation)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logx.Infof("Plugin directory exists but no installation record in database for appKey: %s", pluginKey)
			return false // 目录存在但数据库无记录，认为未安装
		}
		logx.Errorf("Failed to query app_installations table: %v", result.Error)
		return false
	}

	logx.Infof("Plugin installed (verified in both filesystem and database): %s", pluginKey)
	return true
}

// GetInstalledPlugins 获取所有已安装的插件列表（仅从数据库获取已验证的安装记录）
func (s *PluginService) GetInstalledPlugins() ([]string, error) {
	if s.db == nil {
		logx.Infof("Database connection not initialized, falling back to filesystem check")
		return s.getInstalledPluginsFromFilesystem()
	}

	var installations []models.AppInstallation
	result := s.db.Where("status = ?", 1).Find(&installations)
	if result.Error != nil {
		logx.Errorf("Failed to query app_installations table: %v", result.Error)
		return s.getInstalledPluginsFromFilesystem()
	}

	var plugins []string
	for _, installation := range installations {
		// 验证对应的目录是否存在
		pluginPath := filepath.Join(s.pluginsDir, installation.AppKey)
		if info, err := os.Stat(pluginPath); err == nil && info.IsDir() {
			plugins = append(plugins, installation.AppKey)
		} else {
			logx.Infof("Database record exists but plugin directory missing: %s", installation.AppKey)
		}
	}
	return plugins, nil
}

// getInstalledPluginsFromFilesystem 从文件系统获取已安装的插件列表
func (s *PluginService) getInstalledPluginsFromFilesystem() ([]string, error) {
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

// IsPluginInstalledInDatabase 仅检查数据库中是否有对应的安装记录
func (s *PluginService) IsPluginInstalledInDatabase(pluginKey string) bool {
	if s.db == nil {
		logx.Infof("Database connection not initialized")
		return false
	}

	var installation models.AppInstallation
	result := s.db.Where("app_key = ? AND status = ?", pluginKey, 1).First(&installation)
	return result.Error == nil
}
