package cms

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"power-admin-server/internal/logic/cms"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
)

// TagListHandler 标签列表
func TagListHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagListReq

		// 获取分页参数
		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				req.Page = p
			}
		}
		if req.Page == 0 {
			req.Page = 1
		}

		if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
				req.PageSize = ps
			}
		}
		if req.PageSize == 0 {
			req.PageSize = 10
		}

		// 调用Logic层
		tags, err := serverCtx.CmsTagLogic.ListTags(r.Context())
		if err != nil {
			logx.Errorf("list tags failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取标签列表失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		// 构建响应数据
		tagInfos := make([]map[string]interface{}, 0)
		for _, tag := range tags {
			tagInfos = append(tagInfos, map[string]interface{}{
				"id":          tag.ID,
				"name":        tag.Name,
				"slug":        tag.Slug,
				"description": tag.Description,
				"color":       tag.Color,
				"usageCount":  tag.UsageCount,
				"status":      tag.Status,
				"createdAt":   tag.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":   tag.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"total": len(tags),
				"list":  tagInfos,
			},
		})
	}
}

// CreateTagHandler 创建标签
func CreateTagHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateTagReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logx.Errorf("decode request failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "请求参数格式错误",
				"data":    nil,
			})
			return
		}

		// 参数验证
		if req.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "标签名称不能为空",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		logicReq := &cms.CreateTagRequest{
			Name:        req.Name,
			Slug:        req.Slug,
			Description: req.Description,
			Color:       req.Color,
			Status:      int8(req.Status),
		}

		tag, err := serverCtx.CmsTagLogic.CreateTag(r.Context(), logicReq)
		if err != nil {
			logx.Errorf("create tag failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "创建标签失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":          tag.ID,
				"name":        tag.Name,
				"slug":        tag.Slug,
				"description": tag.Description,
				"color":       tag.Color,
				"usageCount":  tag.UsageCount,
				"status":      tag.Status,
				"createdAt":   tag.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":   tag.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// GetTagHandler 获取标签详情
func GetTagHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "标签ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		tag, err := serverCtx.CmsTagLogic.GetTag(r.Context(), tagID)
		if err != nil {
			logx.Errorf("get tag failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取标签失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":          tag.ID,
				"name":        tag.Name,
				"slug":        tag.Slug,
				"description": tag.Description,
				"color":       tag.Color,
				"usageCount":  tag.UsageCount,
				"status":      tag.Status,
				"createdAt":   tag.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":   tag.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// UpdateTagHandler 更新标签
func UpdateTagHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "标签ID格式错误",
				"data":    nil,
			})
			return
		}

		var req types.UpdateTagReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logx.Errorf("decode request failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "请求参数格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		logicReq := &cms.UpdateTagRequest{
			ID:          tagID,
			Name:        req.Name,
			Slug:        req.Slug,
			Description: req.Description,
			Color:       req.Color,
			Status:      int8(req.Status),
		}

		tag, err := serverCtx.CmsTagLogic.UpdateTag(r.Context(), logicReq)
		if err != nil {
			logx.Errorf("update tag failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "更新标签失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":          tag.ID,
				"name":        tag.Name,
				"slug":        tag.Slug,
				"description": tag.Description,
				"color":       tag.Color,
				"usageCount":  tag.UsageCount,
				"status":      tag.Status,
				"createdAt":   tag.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":   tag.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// DeleteTagHandler 删除标签
func DeleteTagHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "标签ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsTagLogic.DeleteTag(r.Context(), tagID); err != nil {
			logx.Errorf("delete tag failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "删除标签失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    nil,
		})
	}
}
