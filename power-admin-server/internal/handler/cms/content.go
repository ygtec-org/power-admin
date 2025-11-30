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

// ContentListHandler 内容列表
func ContentListHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ContentListReq

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

		// 获取过滤参数
		if categoryIdStr := r.URL.Query().Get("categoryId"); categoryIdStr != "" {
			if cid, err := strconv.ParseInt(categoryIdStr, 10, 64); err == nil && cid > 0 {
				req.CategoryId = cid
			}
		}

		if statusStr := r.URL.Query().Get("status"); statusStr != "" {
			if status, err := strconv.Atoi(statusStr); err == nil {
				req.Status = status
			}
		}

		req.Search = r.URL.Query().Get("search")
		req.SortBy = r.URL.Query().Get("sortBy")
		req.SortOrder = r.URL.Query().Get("sortOrder")

		// 调用Logic层
		logicReq := &cms.ListContentRequest{
			Page:       req.Page,
			PageSize:   req.PageSize,
			CategoryID: req.CategoryId,
			Status:     int8(req.Status),
			Search:     req.Search,
			SortBy:     req.SortBy,
			SortOrder:  req.SortOrder,
		}

		result, err := serverCtx.CmsContentLogic.ListContent(r.Context(), logicReq)
		if err != nil {
			logx.Errorf("list content failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取内容列表失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    result,
		})
	}
}

// CreateContentHandler 创建内容
func CreateContentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateContentReq
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
		if req.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容标题不能为空",
				"data":    nil,
			})
			return
		}

		if req.Content == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容不能为空",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		logicReq := &cms.CreateContentRequest{
			Title:            req.Title,
			Slug:             req.Slug,
			Description:      req.Description,
			Content:          req.Content,
			FeaturedImage:    req.FeaturedImage,
			FeaturedImageAlt: req.FeaturedImageAlt,
			CategoryID:       req.CategoryId,
			AuthorID:         req.AuthorId,
			Status:           int8(req.Status),
			Visibility:       int8(req.Visibility),
			CommentStatus:    int8(req.CommentStatus),
			SeoTitle:         req.SeoTitle,
			SeoKeywords:      req.SeoKeywords,
			SeoDescription:   req.SeoDescription,
			IsFeatured:       int8(req.IsFeatured),
			IsSticky:         int8(req.IsSticky),
			ScheduledAt:      req.ScheduledAt,
		}

		content, err := serverCtx.CmsContentLogic.CreateContent(r.Context(), logicReq)
		if err != nil {
			logx.Errorf("create content failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "创建内容失败: " + err.Error(),
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
				"id":               content.ID,
				"title":            content.Title,
				"slug":             content.Slug,
				"description":      content.Description,
				"content":          content.Content,
				"featuredImage":    content.FeaturedImage,
				"featuredImageAlt": content.FeaturedImageAlt,
				"categoryId":       content.CategoryID,
				"authorId":         content.AuthorID,
				"status":           content.Status,
				"visibility":       content.Visibility,
				"commentStatus":    content.CommentStatus,
				"viewCount":        content.ViewCount,
				"commentCount":     content.CommentCount,
				"likeCount":        content.LikeCount,
				"seoTitle":         content.SeoTitle,
				"seoKeywords":      content.SeoKeywords,
				"seoDescription":   content.SeoDescription,
				"publishedAt":      content.PublishedAt,
				"scheduledAt":      content.ScheduledAt,
				"isFeatured":       content.IsFeatured,
				"isSticky":         content.IsSticky,
				"createdAt":        content.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":        content.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// GetContentHandler 获取内容详情
func GetContentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		content, err := serverCtx.CmsContentLogic.GetContent(r.Context(), contentID)
		if err != nil {
			logx.Errorf("get content failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取内容失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":               content.ID,
				"title":            content.Title,
				"slug":             content.Slug,
				"description":      content.Description,
				"content":          content.Content,
				"featuredImage":    content.FeaturedImage,
				"featuredImageAlt": content.FeaturedImageAlt,
				"categoryId":       content.CategoryID,
				"authorId":         content.AuthorID,
				"status":           content.Status,
				"visibility":       content.Visibility,
				"commentStatus":    content.CommentStatus,
				"viewCount":        content.ViewCount,
				"commentCount":     content.CommentCount,
				"likeCount":        content.LikeCount,
				"seoTitle":         content.SeoTitle,
				"seoKeywords":      content.SeoKeywords,
				"seoDescription":   content.SeoDescription,
				"publishedAt":      content.PublishedAt,
				"scheduledAt":      content.ScheduledAt,
				"isFeatured":       content.IsFeatured,
				"isSticky":         content.IsSticky,
				"createdAt":        content.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":        content.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// UpdateContentHandler 更新内容
func UpdateContentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容ID格式错误",
				"data":    nil,
			})
			return
		}

		var req types.UpdateContentReq
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
		logicReq := &cms.UpdateContentRequest{
			ID:               contentID,
			Title:            req.Title,
			Slug:             req.Slug,
			Description:      req.Description,
			Content:          req.Content,
			FeaturedImage:    req.FeaturedImage,
			FeaturedImageAlt: req.FeaturedImageAlt,
			CategoryID:       req.CategoryId,
			Status:           int8(req.Status),
			Visibility:       int8(req.Visibility),
			CommentStatus:    int8(req.CommentStatus),
			SeoTitle:         req.SeoTitle,
			SeoKeywords:      req.SeoKeywords,
			SeoDescription:   req.SeoDescription,
			IsFeatured:       int8(req.IsFeatured),
			IsSticky:         int8(req.IsSticky),
			ScheduledAt:      req.ScheduledAt,
		}

		content, err := serverCtx.CmsContentLogic.UpdateContent(r.Context(), logicReq)
		if err != nil {
			logx.Errorf("update content failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "更新内容失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":               content.ID,
				"title":            content.Title,
				"slug":             content.Slug,
				"description":      content.Description,
				"content":          content.Content,
				"featuredImage":    content.FeaturedImage,
				"featuredImageAlt": content.FeaturedImageAlt,
				"categoryId":       content.CategoryID,
				"authorId":         content.AuthorID,
				"status":           content.Status,
				"visibility":       content.Visibility,
				"commentStatus":    content.CommentStatus,
				"viewCount":        content.ViewCount,
				"commentCount":     content.CommentCount,
				"likeCount":        content.LikeCount,
				"seoTitle":         content.SeoTitle,
				"seoKeywords":      content.SeoKeywords,
				"seoDescription":   content.SeoDescription,
				"publishedAt":      content.PublishedAt,
				"scheduledAt":      content.ScheduledAt,
				"isFeatured":       content.IsFeatured,
				"isSticky":         content.IsSticky,
				"createdAt":        content.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":        content.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// DeleteContentHandler 删除内容
func DeleteContentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsContentLogic.DeleteContent(r.Context(), contentID); err != nil {
			logx.Errorf("delete content failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "删除内容失败: " + err.Error(),
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

// PublishContentHandler 发布内容
func PublishContentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsContentLogic.PublishContent(r.Context(), contentID); err != nil {
			logx.Errorf("publish content failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "发布内容失败: " + err.Error(),
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

// UnpublishContentHandler 取消发布内容
func UnpublishContentHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsContentLogic.UnpublishContent(r.Context(), contentID); err != nil {
			logx.Errorf("unpublish content failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "取消发布内容失败: " + err.Error(),
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

// BatchUpdateStatusHandler 批量更新内容状态
func BatchUpdateStatusHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchUpdateStatusReq
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
		if len(req.Ids) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "内容ID列表不能为空",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsContentLogic.BatchUpdateContentStatus(r.Context(), req.Ids, int8(req.Status)); err != nil {
			logx.Errorf("batch update status failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "批量更新状态失败: " + err.Error(),
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
