// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package codegen

import (
	"context"
	"fmt"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHistoryLogic {
	return &GetHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHistoryLogic) GetHistory() (resp *types.GenHistoryResp, err error) {
	id := l.ctx.Value("id").(int64)

	history, err := l.svcCtx.CodegenRepo.GetHistory(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("历史记录不存在")
	}

	return &types.GenHistoryResp{
		ID:          history.ID,
		GenConfigID: history.GenConfigID,
		TableName:   history.Table,
		FilePath:    history.FilePath,
		FileType:    history.FileType,
		Content:     history.Content,
		Status:      history.Status,
		ErrorMsg:    history.ErrorMsg,
		Operator:    history.Operator,
		CreatedAt:   history.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
