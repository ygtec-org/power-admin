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

// CategoryListHandler 分类列表
func CategoryListHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析查询参数
		parentID := int64(0)

		if parentIDStr := r.URL.Query().Get("parentId"); parentIDStr != "" {
			if pid, err := strconv.ParseInt(parentIDStr, 10, 64); err == nil && pid > 0 {
				parentID = pid
			}
		}

		// 调用Logic层
		var parentIDPtr *int64
		if parentID > 0 {
			parentIDPtr = &parentID
		}

		categories, err := serverCtx.CmsCategoryLogic.ListCategories(r.Context(), parentIDPtr)
		if err != nil {
			logx.Errorf("list categories failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取分类列表失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		// 构建响应数据
		categoryInfos := make([]map[string]interface{}, 0)
		for _, cat := range categories {
			categoryInfos = append(categoryInfos, map[string]interface{}{
				"id":             cat.ID,
				"name":           cat.Name,
				"slug":           cat.Slug,
				"description":    cat.Description,
				"thumbnail":      cat.Thumbnail,
				"parentId":       cat.ParentID,
				"sort":           cat.Sort,
				"status":         cat.Status,
				"contentCount":   cat.ContentCount,
				"seoKeywords":    cat.SeoKeywords,
				"seoDescription": cat.SeoDescription,
				"createdAt":      cat.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":      cat.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"total": len(categories),
				"list":  categoryInfos,
			},
		})
	}
}

// CategoryTreeHandler 分类树形结构
func CategoryTreeHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 调用Logic层获取树形结构
		tree, err := serverCtx.CmsCategoryLogic.GetCategoryTree(r.Context())
		if err != nil {
			logx.Errorf("get category tree failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取分类树失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    tree,
		})
	}
}

// CreateCategoryHandler 创建分类
func CreateCategoryHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateCategoryReq
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
				"message": "分类名称不能为空",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		logicReq := &cms.CreateCategoryRequest{
			Name:           req.Name,
			Slug:           req.Slug,
			Description:    req.Description,
			Thumbnail:      req.Thumbnail,
			ParentID:       req.ParentId,
			Sort:           req.Sort,
			Status:         int8(req.Status),
			SeoKeywords:    req.SeoKeywords,
			SeoDescription: req.SeoDescription,
		}

		category, err := serverCtx.CmsCategoryLogic.CreateCategory(r.Context(), logicReq)
		if err != nil {
			logx.Errorf("create category failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "创建分类失败: " + err.Error(),
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
				"id":             category.ID,
				"name":           category.Name,
				"slug":           category.Slug,
				"description":    category.Description,
				"thumbnail":      category.Thumbnail,
				"parentId":       category.ParentID,
				"sort":           category.Sort,
				"status":         category.Status,
				"contentCount":   category.ContentCount,
				"seoKeywords":    category.SeoKeywords,
				"seoDescription": category.SeoDescription,
				"createdAt":      category.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":      category.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// GetCategoryHandler 获取分类详情
func GetCategoryHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取路径参数
		categoryID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "分类ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		category, err := serverCtx.CmsCategoryLogic.GetCategory(r.Context(), categoryID)
		if err != nil {
			logx.Errorf("get category failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "获取分类失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":             category.ID,
				"name":           category.Name,
				"slug":           category.Slug,
				"description":    category.Description,
				"thumbnail":      category.Thumbnail,
				"parentId":       category.ParentID,
				"sort":           category.Sort,
				"status":         category.Status,
				"contentCount":   category.ContentCount,
				"seoKeywords":    category.SeoKeywords,
				"seoDescription": category.SeoDescription,
				"createdAt":      category.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":      category.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// UpdateCategoryHandler 更新分类
func UpdateCategoryHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取路径参数
		categoryID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "分类ID格式错误",
				"data":    nil,
			})
			return
		}

		var req types.UpdateCategoryReq
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
		logicReq := &cms.UpdateCategoryRequest{
			ID:             categoryID,
			Name:           req.Name,
			Slug:           req.Slug,
			Description:    req.Description,
			Thumbnail:      req.Thumbnail,
			ParentID:       req.ParentId,
			Sort:           req.Sort,
			Status:         int8(req.Status),
			SeoKeywords:    req.SeoKeywords,
			SeoDescription: req.SeoDescription,
		}

		category, err := serverCtx.CmsCategoryLogic.UpdateCategory(r.Context(), logicReq)
		if err != nil {
			logx.Errorf("update category failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "更新分类失败: " + err.Error(),
				"data":    nil,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":             category.ID,
				"name":           category.Name,
				"slug":           category.Slug,
				"description":    category.Description,
				"thumbnail":      category.Thumbnail,
				"parentId":       category.ParentID,
				"sort":           category.Sort,
				"status":         category.Status,
				"contentCount":   category.ContentCount,
				"seoKeywords":    category.SeoKeywords,
				"seoDescription": category.SeoDescription,
				"createdAt":      category.CreatedAt.Format("2006-01-02 15:04:05"),
				"updatedAt":      category.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// DeleteCategoryHandler 删除分类
func DeleteCategoryHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取路径参数
		categoryID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			logx.Errorf("parse path id failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    400,
				"message": "分类ID格式错误",
				"data":    nil,
			})
			return
		}

		// 调用Logic层
		if err := serverCtx.CmsCategoryLogic.DeleteCategory(r.Context(), categoryID); err != nil {
			logx.Errorf("delete category failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    500,
				"message": "删除分类失败: " + err.Error(),
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
