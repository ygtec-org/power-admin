// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package codegen

import (
	"context"
	"fmt"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/codegen"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreviewCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPreviewCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewCodeLogic {
	return &PreviewCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreviewCodeLogic) PreviewCode(req *types.CodePreviewReq) (resp *types.CodePreviewResp, err error) {
	config, err := l.svcCtx.CodegenRepo.GetConfig(l.ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("配置不存在")
	}

	engine := codegen.NewTemplateEngine()
	files := make([]types.GeneratedFile, 0)

	if apiContent, err := engine.RenderAPI(config, config.Columns); err == nil {
		files = append(files, types.GeneratedFile{
			FilePath: fmt.Sprintf("api/%s.api", config.ModuleName),
			FileType: "api",
			Content:  apiContent,
		})
	}

	if modelContent, err := engine.RenderModel(config, config.Columns); err == nil {
		files = append(files, types.GeneratedFile{
			FilePath: fmt.Sprintf("pkg/models/%s_models.go", config.ModuleName),
			FileType: "model",
			Content:  modelContent,
		})
	}

	if logicContent, err := engine.RenderLogic(config, config.Columns); err == nil {
		files = append(files, types.GeneratedFile{
			FilePath: fmt.Sprintf("internal/logic/%s/", config.ModuleName),
			FileType: "logic",
			Content:  logicContent,
		})
	}

	if repoContent, err := engine.RenderRepository(config, config.Columns); err == nil {
		files = append(files, types.GeneratedFile{
			FilePath: fmt.Sprintf("pkg/repository/%s_repository.go", config.ModuleName),
			FileType: "repository",
			Content:  repoContent,
		})
	}

	return &types.CodePreviewResp{
		Files: files,
	}, nil
}
