// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package codegen

import (
	"context"
	"fmt"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/codegen"
	"power-admin-server/pkg/models"
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

	engine := codegen.NewTemplateEngine()
	files := make([]types.GeneratedFile, 0)
	histories := make([]*models.GenHistory, 0)

	// 生成API文件
	if apiContent, err := engine.RenderAPI(config, config.Columns); err == nil {
		apiFile := types.GeneratedFile{
			FilePath: fmt.Sprintf("api/%s.api", config.ModuleName),
			FileType: "api",
			Content:  apiContent,
		}
		files = append(files, apiFile)
		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			Table:       config.Table,
			FilePath:    apiFile.FilePath,
			FileType:    "api",
			Content:     apiContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	// 生成Model文件
	if modelContent, err := engine.RenderModel(config, config.Columns); err == nil {
		modelFile := types.GeneratedFile{
			FilePath: fmt.Sprintf("pkg/models/%s_models.go", config.ModuleName),
			FileType: "model",
			Content:  modelContent,
		}
		files = append(files, modelFile)
		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			Table:       config.Table,
			FilePath:    modelFile.FilePath,
			FileType:    "model",
			Content:     modelContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	// 生成Logic文件
	if logicContent, err := engine.RenderLogic(config, config.Columns); err == nil {
		logicFile := types.GeneratedFile{
			FilePath: fmt.Sprintf("internal/logic/%s/", config.ModuleName),
			FileType: "logic",
			Content:  logicContent,
		}
		files = append(files, logicFile)
		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			Table:       config.Table,
			FilePath:    logicFile.FilePath,
			FileType:    "logic",
			Content:     logicContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	// 生成Repository文件
	if repoContent, err := engine.RenderRepository(config, config.Columns); err == nil {
		repoFile := types.GeneratedFile{
			FilePath: fmt.Sprintf("pkg/repository/%s_repository.go", config.ModuleName),
			FileType: "repository",
			Content:  repoContent,
		}
		files = append(files, repoFile)
		histories = append(histories, &models.GenHistory{
			GenConfigID: config.ID,
			Table:       config.Table,
			FilePath:    repoFile.FilePath,
			FileType:    "repository",
			Content:     repoContent,
			Status:      1,
			CreatedAt:   time.Now(),
		})
	}

	if len(histories) > 0 {
		if err := l.svcCtx.CodegenRepo.CreateHistories(l.ctx, histories); err != nil {
			logx.Errorf("保存生成历史失败: %v", err)
		}
	}

	return &types.CodeGenerateResp{
		Success: true,
		Message: "代码生成成功",
		Files:   files,
	}, nil
}
