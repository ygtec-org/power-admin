package repository

import (
	"context"
	"errors"
	"power-admin-server/pkg/models"
	"time"

	"gorm.io/gorm"
)

// =============================================
// CMS 仓储接口定义
// =============================================

type ContentRepository interface {
	Create(ctx context.Context, content *models.CmsContent) error
	Get(ctx context.Context, id int64) (*models.CmsContent, error)
	Update(ctx context.Context, id int64, content *models.CmsContent) error
	Delete(ctx context.Context, id int64) error
	SoftDelete(ctx context.Context, id int64) error
	List(ctx context.Context, req *ContentListRequest) (*PagedResult, error)
	GetAll(ctx context.Context) ([]*models.CmsContent, error)
	GetBySlug(ctx context.Context, slug string) (*models.CmsContent, error)
	Publish(ctx context.Context, id int64) error
	Unpublish(ctx context.Context, id int64) error
	IncrementViewCount(ctx context.Context, id int64) error
	BatchUpdateStatus(ctx context.Context, ids []int64, status int8) error
}

type CategoryRepository interface {
	Create(ctx context.Context, category *models.CmsCategory) error
	Get(ctx context.Context, id int64) (*models.CmsCategory, error)
	Update(ctx context.Context, id int64, category *models.CmsCategory) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, parentID *int64) ([]*models.CmsCategory, error)
	GetAll(ctx context.Context) ([]*models.CmsCategory, error)
	GetTree(ctx context.Context) ([]*models.CmsCategory, error)
	GetBySlug(ctx context.Context, slug string) (*models.CmsCategory, error)
	UpdateContentCount(ctx context.Context, categoryID int64, count int) error
}

type TagRepository interface {
	Create(ctx context.Context, tag *models.CmsTag) error
	Get(ctx context.Context, id int64) (*models.CmsTag, error)
	Update(ctx context.Context, id int64, tag *models.CmsTag) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, pageSize int) ([]*models.CmsTag, int64, error)
	GetAll(ctx context.Context) ([]*models.CmsTag, error)
	GetByName(ctx context.Context, name string) (*models.CmsTag, error)
	GetByIDs(ctx context.Context, ids []int64) ([]*models.CmsTag, error)
	UpdateUsageCount(ctx context.Context, tagID int64, count int) error
}

type CmsUserRepository interface {
	Create(ctx context.Context, user *models.CmsUser) error
	Get(ctx context.Context, id int64) (*models.CmsUser, error)
	Update(ctx context.Context, id int64, user *models.CmsUser) error
	Delete(ctx context.Context, id int64) error
	SoftDelete(ctx context.Context, id int64) error
	GetByUsername(ctx context.Context, username string) (*models.CmsUser, error)
	GetByEmail(ctx context.Context, email string) (*models.CmsUser, error)
	List(ctx context.Context) ([]*models.CmsUser, error)
	GetAll(ctx context.Context) ([]*models.CmsUser, error)
	UpdateLastLogin(ctx context.Context, id int64, ip string) error
}

type CommentRepository interface {
	Create(ctx context.Context, comment *models.CmsComment) error
	Get(ctx context.Context, id int64) (*models.CmsComment, error)
	Update(ctx context.Context, id int64, comment *models.CmsComment) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, contentID, page, pageSize int64) ([]*models.CmsComment, int64, error)
	GetAll(ctx context.Context) ([]*models.CmsComment, error)
	Approve(ctx context.Context, id int64) error
	UpdateLikeCount(ctx context.Context, id int64, count int) error
}

// =============================================
// CMS 内容仓储实现
// =============================================
type contentRepositoryImpl struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepositoryImpl{db: db}
}

func (r *contentRepositoryImpl) Create(ctx context.Context, content *models.CmsContent) error {
	return r.db.WithContext(ctx).Create(content).Error
}

func (r *contentRepositoryImpl) Get(ctx context.Context, id int64) (*models.CmsContent, error) {
	var content models.CmsContent
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&content).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &content, nil
}

func (r *contentRepositoryImpl) Update(ctx context.Context, id int64, content *models.CmsContent) error {
	return r.db.WithContext(ctx).Model(&models.CmsContent{}).Where("id = ?", id).Updates(content).Error
}

func (r *contentRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.CmsContent{}, id).Error
}

func (r *contentRepositoryImpl) SoftDelete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&models.CmsContent{}).Where("id = ?", id).
		Update("status", 3).Error
}

type ContentListRequest struct {
	Page       int
	PageSize   int
	CategoryID *int64
	Status     *int8
	Visibility *int8
	Search     string
	SortBy     string
	SortOrder  string
	IsSticky   *int8
	IsFeatured *int8
}

func (r *contentRepositoryImpl) List(ctx context.Context, req *ContentListRequest) (*PagedResult, error) {
	query := r.db.WithContext(ctx)

	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.Visibility != nil {
		query = query.Where("visibility = ?", *req.Visibility)
	}
	if req.IsSticky != nil {
		query = query.Where("is_sticky = ?", *req.IsSticky)
	}
	if req.IsFeatured != nil {
		query = query.Where("is_featured = ?", *req.IsFeatured)
	}
	if req.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if req.SortBy != "" {
		order := "DESC"
		if req.SortOrder == "asc" {
			order = "ASC"
		}
		query = query.Order(req.SortBy + " " + order)
	} else {
		query = query.Order("created_at DESC")
	}

	var total int64
	if err := query.Model(&models.CmsContent{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	var contents []*models.CmsContent
	err := query.Offset(offset).Limit(req.PageSize).Find(&contents).Error
	if err != nil {
		return nil, err
	}

	return &PagedResult{
		Total:    total,
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		Items:    contents,
	}, nil
}

func (r *contentRepositoryImpl) GetAll(ctx context.Context) ([]*models.CmsContent, error) {
	var contents []*models.CmsContent
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&contents).Error
	return contents, err
}

func (r *contentRepositoryImpl) GetBySlug(ctx context.Context, slug string) (*models.CmsContent, error) {
	var content models.CmsContent
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&content).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &content, nil
}

func (r *contentRepositoryImpl) Publish(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.CmsContent{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       2,
			"published_at": now,
		}).Error
}

func (r *contentRepositoryImpl) Unpublish(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&models.CmsContent{}).Where("id = ?", id).
		Update("status", 1).Error
}

func (r *contentRepositoryImpl) IncrementViewCount(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&models.CmsContent{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *contentRepositoryImpl) BatchUpdateStatus(ctx context.Context, ids []int64, status int8) error {
	return r.db.WithContext(ctx).Model(&models.CmsContent{}).Where("id IN ?", ids).
		Update("status", status).Error
}

// =============================================
// CMS 分类仓储实现
// =============================================
type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}

func (r *categoryRepositoryImpl) Create(ctx context.Context, category *models.CmsCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepositoryImpl) Get(ctx context.Context, id int64) (*models.CmsCategory, error) {
	var category models.CmsCategory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepositoryImpl) Update(ctx context.Context, id int64, category *models.CmsCategory) error {
	return r.db.WithContext(ctx).Model(&models.CmsCategory{}).Where("id = ?", id).Updates(category).Error
}

func (r *categoryRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.CmsCategory{}, id).Error
}

func (r *categoryRepositoryImpl) List(ctx context.Context, parentID *int64) ([]*models.CmsCategory, error) {
	var categories []*models.CmsCategory
	query := r.db.WithContext(ctx)

	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	err := query.Order("sort DESC, created_at DESC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepositoryImpl) GetAll(ctx context.Context) ([]*models.CmsCategory, error) {
	var categories []*models.CmsCategory
	err := r.db.WithContext(ctx).Order("sort DESC, created_at DESC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepositoryImpl) GetTree(ctx context.Context) ([]*models.CmsCategory, error) {
	var categories []*models.CmsCategory
	err := r.db.WithContext(ctx).Where("parent_id IS NULL").Order("sort DESC").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	for _, cat := range categories {
		r.loadChildren(ctx, cat)
	}

	return categories, nil
}

func (r *categoryRepositoryImpl) loadChildren(ctx context.Context, category *models.CmsCategory) {
	var children []*models.CmsCategory
	r.db.WithContext(ctx).Where("parent_id = ?", category.ID).Order("sort DESC").Find(&children)
	category.Children = children

	for _, child := range children {
		r.loadChildren(ctx, child)
	}
}

func (r *categoryRepositoryImpl) GetBySlug(ctx context.Context, slug string) (*models.CmsCategory, error) {
	var category models.CmsCategory
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepositoryImpl) UpdateContentCount(ctx context.Context, categoryID int64, count int) error {
	return r.db.WithContext(ctx).Model(&models.CmsCategory{}).Where("id = ?", categoryID).
		Update("content_count", count).Error
}

// =============================================
// CMS 标签仓储实现
// =============================================
type tagRepositoryImpl struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepositoryImpl{db: db}
}

func (r *tagRepositoryImpl) Create(ctx context.Context, tag *models.CmsTag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepositoryImpl) Get(ctx context.Context, id int64) (*models.CmsTag, error) {
	var tag models.CmsTag
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepositoryImpl) Update(ctx context.Context, id int64, tag *models.CmsTag) error {
	return r.db.WithContext(ctx).Model(&models.CmsTag{}).Where("id = ?", id).Updates(tag).Error
}

func (r *tagRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.CmsTag{}, id).Error
}

func (r *tagRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*models.CmsTag, int64, error) {
	var tags []*models.CmsTag
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Where("status = ?", 1).Order("usage_count DESC").Offset(offset).Limit(pageSize).Find(&tags).Error
	var total int64
	err = r.db.WithContext(ctx).Model(&models.CmsTag{}).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return tags, total, err
}

func (r *tagRepositoryImpl) GetAll(ctx context.Context) ([]*models.CmsTag, error) {
	var tags []*models.CmsTag
	err := r.db.WithContext(ctx).Order("usage_count DESC").Find(&tags).Error
	return tags, err
}

func (r *tagRepositoryImpl) GetByName(ctx context.Context, name string) (*models.CmsTag, error) {
	var tag models.CmsTag
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepositoryImpl) GetByIDs(ctx context.Context, ids []int64) ([]*models.CmsTag, error) {
	var tags []*models.CmsTag
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

func (r *tagRepositoryImpl) UpdateUsageCount(ctx context.Context, tagID int64, count int) error {
	return r.db.WithContext(ctx).Model(&models.CmsTag{}).Where("id = ?", tagID).
		Update("usage_count", count).Error
}

// =============================================
// CMS 用户仓储实现
// =============================================
type cmsUserRepositoryImpl struct {
	db *gorm.DB
}

func NewCmsUserRepository(db *gorm.DB) CmsUserRepository {
	return &cmsUserRepositoryImpl{db: db}
}

func (r *cmsUserRepositoryImpl) Create(ctx context.Context, user *models.CmsUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *cmsUserRepositoryImpl) Get(ctx context.Context, id int64) (*models.CmsUser, error) {
	var user models.CmsUser
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *cmsUserRepositoryImpl) Update(ctx context.Context, id int64, user *models.CmsUser) error {
	return r.db.WithContext(ctx).Model(&models.CmsUser{}).Where("id = ?", id).Updates(user).Error
}

func (r *cmsUserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.CmsUser{}, id).Error
}

func (r *cmsUserRepositoryImpl) SoftDelete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&models.CmsUser{}).Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *cmsUserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*models.CmsUser, error) {
	var user models.CmsUser
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *cmsUserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.CmsUser, error) {
	var user models.CmsUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *cmsUserRepositoryImpl) List(ctx context.Context) ([]*models.CmsUser, error) {
	var users []*models.CmsUser
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *cmsUserRepositoryImpl) GetAll(ctx context.Context) ([]*models.CmsUser, error) {
	var users []*models.CmsUser
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *cmsUserRepositoryImpl) UpdateLastLogin(ctx context.Context, id int64, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.CmsUser{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at": now,
			"last_login_ip": ip,
			"login_count":   gorm.Expr("login_count + ?", 1),
		}).Error
}

// =============================================
// CMS 评论仓储实现
// =============================================
type commentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepositoryImpl{db: db}
}

func (r *commentRepositoryImpl) Create(ctx context.Context, comment *models.CmsComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepositoryImpl) Get(ctx context.Context, id int64) (*models.CmsComment, error) {
	var comment models.CmsComment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepositoryImpl) Update(ctx context.Context, id int64, comment *models.CmsComment) error {
	return r.db.WithContext(ctx).Model(&models.CmsComment{}).Where("id = ?", id).Updates(comment).Error
}

func (r *commentRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.CmsComment{}, id).Error
}

func (r *commentRepositoryImpl) List(ctx context.Context, contentID, page, pageSize int64) ([]*models.CmsComment, int64, error) {
	var comments []*models.CmsComment
	var total int64
	err := r.db.WithContext(ctx).Model(&models.CmsComment{}).Where("content_id = ? AND status = ?", contentID, 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err = r.db.WithContext(ctx).Where("content_id = ? AND status = ?", contentID, 1).Offset(int(offset)).Limit(int(pageSize)).
		Order("created_at DESC").Find(&comments).Error
	return comments, total, err
}

func (r *commentRepositoryImpl) GetAll(ctx context.Context) ([]*models.CmsComment, error) {
	var comments []*models.CmsComment
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&comments).Error
	return comments, err
}

func (r *commentRepositoryImpl) Approve(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.CmsComment{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      1,
			"approved_at": now,
		}).Error
}

func (r *commentRepositoryImpl) UpdateLikeCount(ctx context.Context, id int64, count int) error {
	return r.db.WithContext(ctx).Model(&models.CmsComment{}).Where("id = ?", id).
		Update("like_count", count).Error
}

// =============================================
// PagedResult 分页结果
// =============================================
type PagedResult struct {
	Total    int64       `json:"total"`
	Page     int64       `json:"page"`
	PageSize int64       `json:"page_size"`
	Items    interface{} `json:"items"`
}
