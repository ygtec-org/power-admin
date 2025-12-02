// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package codegen

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/codegen"
	"power-admin-server/pkg/models"
	"runtime"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateCodeLogic {
	return &GenerateCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateCodeLogic) GenerateCode(req *types.CodeGenerateReq) (resp *types.CodeGenerateResp, err error) {
	config, err := l.svcCtx.CodegenRepo.GetConfig(l.ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("配置不存在")
	}

	// 获取项目根目录
	projectRoot, err := getProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("获取项目根目录失败: %v", err)
	}

	engine := codegen.NewTemplateEngine()
	files := make([]types.GeneratedFile, 0)
	histories := make([]*models.GenHistory, 0)
	var writeErrors []string

	// 生成API文件
	if apiContent, err := engine.RenderAPI(config, config.Columns); err == nil {
		relativePath := fmt.Sprintf("api/%s.api", config.ModuleName)
		absolutePath := filepath.Join(projectRoot, "power-admin-server", relativePath)

		// 写入文件
		if writeErr := writeFile(absolutePath, apiContent); writeErr != nil {
			writeErrors = append(writeErrors, fmt.Sprintf("写入API文件失败: %v", writeErr))
		}

		apiFile := types.GeneratedFile{
			FilePath: relativePath,
			FileType: "api",
			Content:  apiContent,
		}
		files = append(files, apiFile)
		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			Table:       config.Table,
			FilePath:    relativePath,
			FileType:    "api",
			Content:     apiContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	// 生成Model文件
	if modelContent, err := engine.RenderModel(config, config.Columns); err == nil {
		relativePath := fmt.Sprintf("pkg/models/%s_models.go", config.ModuleName)
		absolutePath := filepath.Join(projectRoot, "power-admin-server", relativePath)

		// 写入文件
		if writeErr := writeFile(absolutePath, modelContent); writeErr != nil {
			writeErrors = append(writeErrors, fmt.Sprintf("写入Model文件失败: %v", writeErr))
		}

		modelFile := types.GeneratedFile{
			FilePath: relativePath,
			FileType: "model",
			Content:  modelContent,
		}
		files = append(files, modelFile)
		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			Table:       config.Table,
			FilePath:    relativePath,
			FileType:    "model",
			Content:     modelContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	// 生成Repository文件
	if repoContent, err := engine.RenderRepository(config, config.Columns); err == nil {
		relativePath := fmt.Sprintf("pkg/repository/%s_repository.go", config.ModuleName)
		absolutePath := filepath.Join(projectRoot, "power-admin-server", relativePath)

		// 写入文件
		if writeErr := writeFile(absolutePath, repoContent); writeErr != nil {
			writeErrors = append(writeErrors, fmt.Sprintf("写入Repository文件失败: %v", writeErr))
		}

		repoFile := types.GeneratedFile{
			FilePath: relativePath,
			FileType: "repository",
			Content:  repoContent,
		}
		files = append(files, repoFile)
		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			Table:       config.Table,
			FilePath:    relativePath,
			FileType:    "repository",
			Content:     repoContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	// 生成前端Vue页面
	if vueContent, err := engine.RenderVue(config, config.Columns); err == nil {
		relativePath := fmt.Sprintf("pages/%s/%sList.vue", config.ModuleName, strings.Title(config.ModuleName))
		absolutePath := filepath.Join(projectRoot, "power-admin-web", "src", relativePath)

		// 写入文件
		if writeErr := writeFile(absolutePath, vueContent); writeErr != nil {
			writeErrors = append(writeErrors, fmt.Sprintf("写入Vue文件失败: %v", writeErr))
		}

		vueFile := types.GeneratedFile{
			FilePath: "web/" + relativePath,
			FileType: "vue",
			Content:  vueContent,
		}
		files = append(files, vueFile)
		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			Table:       config.Table,
			FilePath:    "web/" + relativePath,
			FileType:    "vue",
			Content:     vueContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	// 保存历史记录
	if len(histories) > 0 {
		if err := l.svcCtx.CodegenRepo.CreateHistories(l.ctx, histories); err != nil {
			logx.Errorf("保存生成历史失败: %v", err)
		}
	}

	// 执行 make gen 生成 handler 和 router
	if err := runMakeGen(projectRoot); err != nil {
		writeErrors = append(writeErrors, fmt.Sprintf("执行make gen失败: %v", err))
	}

	message := fmt.Sprintf("代码生成成功！共生成 %d 个文件", len(files))
	if len(writeErrors) > 0 {
		message += "\n\n⚠️ 部分文件写入失败:\n" + strings.Join(writeErrors, "\n")
	}

	return &types.CodeGenerateResp{
		Success: len(writeErrors) == 0,
		Message: message,
		Files:   files,
	}, nil
}

// getProjectRoot 获取项目根目录
func getProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("无法获取当前文件路径")
	}
	// 从当前文件路径向上查找到 power-admin 目录
	dir := filepath.Dir(filename)
	for {
		if filepath.Base(dir) == "power-admin" {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("未找到 power-admin 根目录")
		}
		dir = parent
	}
}

// writeFile 写入文件，自动创建目录
func writeFile(filePath, content string) error {
	// 创建目录
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// runMakeGen 执行 make gen 命令
func runMakeGen(projectRoot string) error {
	serverPath := filepath.Join(projectRoot, "power-admin-server")
	cmd := exec.Command("make", "gen")
	cmd.Dir = serverPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("执行make gen失败: %v", err)
	}

	return nil
}
