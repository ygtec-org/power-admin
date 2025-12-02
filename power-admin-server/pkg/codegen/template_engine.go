package codegen

import (
	"bytes"
	"fmt"
	"power-admin-server/pkg/models"
	"power-admin-server/pkg/repository"
	"strings"
	"text/template"
	"time"
)

// TemplateData 模板数据
type TemplateData struct {
	TableName    string
	TablePrefix  string
	BusinessName string
	ModuleName   string
	PackageName  string
	StructName   string
	Author       string
	Date         string
	Fields       []FieldData
}

// FieldData 字段数据
type FieldData struct {
	ColumnName    string
	ColumnComment string
	ColumnType    string
	GoType        string
	GoField       string
	JsonField     string
	IsPk          bool
	IsIncrement   bool
	IsRequired    bool
	IsInsert      bool
	IsEdit        bool
	IsList        bool
	IsQuery       bool
	QueryType     string
	HtmlType      string
}

// TemplateEngine 模板引擎
type TemplateEngine struct {
	templates map[string]string
}

// NewTemplateEngine 创建模板引擎
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{
		templates: GetDefaultTemplates(),
	}
}

// RenderAPI 渲染API文件
func (e *TemplateEngine) RenderAPI(config *models.GenConfig, columns []*models.GenTableColumn) (string, error) {
	data := e.prepareTemplateData(config, columns)
	return e.render("api", data)
}

// RenderModel 渲染Model文件
func (e *TemplateEngine) RenderModel(config *models.GenConfig, columns []*models.GenTableColumn) (string, error) {
	data := e.prepareTemplateData(config, columns)
	return e.render("model", data)
}

// RenderLogic 渲染Logic文件
func (e *TemplateEngine) RenderLogic(config *models.GenConfig, columns []*models.GenTableColumn) (string, error) {
	data := e.prepareTemplateData(config, columns)
	return e.render("logic", data)
}

// RenderRepository 渲染Repository文件
func (e *TemplateEngine) RenderRepository(config *models.GenConfig, columns []*models.GenTableColumn) (string, error) {
	data := e.prepareTemplateData(config, columns)
	return e.render("repository", data)
}

// RenderVue 渲染Vue文件
func (e *TemplateEngine) RenderVue(config *models.GenConfig, columns []*models.GenTableColumn) (string, error) {
	data := e.prepareTemplateData(config, columns)
	return e.render("vue", data)
}

// prepareTemplateData 准备模板数据
func (e *TemplateEngine) prepareTemplateData(config *models.GenConfig, columns []*models.GenTableColumn) *TemplateData {
	structName := repository.TableNameToStructName(config.Table, config.TablePrefix)

	fields := make([]FieldData, 0, len(columns))
	for _, col := range columns {
		fields = append(fields, FieldData{
			ColumnName:    col.ColumnName,
			ColumnComment: col.ColumnComment,
			ColumnType:    col.ColumnType,
			GoType:        col.GoType,
			GoField:       col.GoField,
			JsonField:     toJsonField(col.GoField),
			IsPk:          col.IsPk == 1,
			IsIncrement:   col.IsIncrement == 1,
			IsRequired:    col.IsRequired == 1,
			IsInsert:      col.IsInsert == 1,
			IsEdit:        col.IsEdit == 1,
			IsList:        col.IsList == 1,
			IsQuery:       col.IsQuery == 1,
			QueryType:     col.QueryType,
			HtmlType:      col.HtmlType,
		})
	}

	return &TemplateData{
		TableName:    config.Table,
		TablePrefix:  config.TablePrefix,
		BusinessName: config.BusinessName,
		ModuleName:   config.ModuleName,
		PackageName:  config.PackageName,
		StructName:   structName,
		Author:       config.Author,
		Date:         time.Now().Format("2006-01-02"),
		Fields:       fields,
	}
}

// render 渲染模板
func (e *TemplateEngine) render(tplName string, data *TemplateData) (string, error) {
	tplStr, ok := e.templates[tplName]
	if !ok {
		return "", fmt.Errorf("template %s not found", tplName)
	}

	tmpl, err := template.New(tplName).Parse(tplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// toJsonField 转换为JSON字段名(小驼峰)
func toJsonField(goField string) string {
	if goField == "" {
		return ""
	}
	return strings.ToLower(goField[:1]) + goField[1:]
}

// GetDefaultTemplates 获取默认模板
func GetDefaultTemplates() map[string]string {
	return map[string]string{
		"api":        apiTemplate,
		"model":      modelTemplate,
		"logic":      logicTemplate,
		"repository": repositoryTemplate,
		"vue":        vueTemplate,
	}
}

// API模板
const apiTemplate = `syntax = "v1"

info(
	title: "{{.StructName}} API"
	desc: "{{.StructName}} 管理接口"
	author: "{{.Author}}"
	date: "{{.Date}}"
)

type (
	// 创建{{.StructName}}请求
	{{.StructName}}CreateReq {
		{{range .Fields}}{{if .IsInsert}}{{.GoField}} {{.GoType}} ` + "`json:\"{{.JsonField}}\"{{if .IsRequired}},validate:\"required\"{{end}}`" + ` {{if .ColumnComment}}// {{.ColumnComment}}{{end}}
		{{end}}{{end}}
	}

	// 更新{{.StructName}}请求
	{{.StructName}}UpdateReq {
		{{range .Fields}}{{if or .IsEdit .IsPk}}{{.GoField}} {{.GoType}} ` + "`json:\"{{.JsonField}}\"{{if .IsRequired}},validate:\"required\"{{end}}`" + ` {{if .ColumnComment}}// {{.ColumnComment}}{{end}}
		{{end}}{{end}}
	}

	// {{.StructName}}响应
	{{.StructName}}Resp {
		{{range .Fields}}{{.GoField}} {{.GoType}} ` + "`json:\"{{.JsonField}}\"`" + ` {{if .ColumnComment}}// {{.ColumnComment}}{{end}}
		{{end}}
	}

	// {{.StructName}}列表请求
	{{.StructName}}ListReq {
		Page     int    ` + "`form:\"page,default=1\"`" + `
		PageSize int    ` + "`form:\"pageSize,default=10\"`" + `
		{{range .Fields}}{{if .IsQuery}}{{.GoField}} {{.GoType}} ` + "`form:\"{{.JsonField}},optional\"`" + ` {{if .ColumnComment}}// {{.ColumnComment}}{{end}}
		{{end}}{{end}}
	}

	// {{.StructName}}列表响应
	{{.StructName}}ListResp {
		Total int64             ` + "`json:\"total\"`" + `
		Data  []{{.StructName}}Resp ` + "`json:\"data\"`" + `
	}
)

@server(
	prefix: /api/admin/{{.ModuleName}}
	group: {{.ModuleName}}
	middleware: AdminAuthMiddleware
)
service power-api {
	@handler Create{{.StructName}}
	post /{{.BusinessName}} ({{.StructName}}CreateReq) returns ({{.StructName}}Resp)

	@handler Update{{.StructName}}
	put /{{.BusinessName}}/:id ({{.StructName}}UpdateReq) returns ({{.StructName}}Resp)

	@handler Delete{{.StructName}}
	delete /{{.BusinessName}}/:id

	@handler Get{{.StructName}}
	get /{{.BusinessName}}/:id returns ({{.StructName}}Resp)

	@handler List{{.StructName}}
	get /{{.BusinessName}}/list ({{.StructName}}ListReq) returns ({{.StructName}}ListResp)
}
`

// Model模板
const modelTemplate = `package models

import "time"

// {{.StructName}} {{.StructName}}模型
type {{.StructName}} struct {
	{{range .Fields}}{{.GoField}} {{.GoType}} ` + "`gorm:\"{{if .IsPk}}primaryKey;{{end}}{{if .IsIncrement}}autoIncrement;{{end}}comment:{{.ColumnComment}}\" json:\"{{.JsonField}}\"`" + `
	{{end}}
}

// TableName 指定表名
func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`

// Logic模板
const logicTemplate = `// TODO: 此文件为生成的模板，请根据实际业务需求修改
package {{.ModuleName}}

import (
	"context"
	"power-admin/internal/svc"
	"power-admin/internal/types"
	"power-admin/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
)

// Create{{.StructName}}Logic 创建{{.StructName}}逻辑
type Create{{.StructName}}Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreate{{.StructName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Create{{.StructName}}Logic {
	return &Create{{.StructName}}Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Create{{.StructName}}Logic) Create{{.StructName}}(req *types.{{.StructName}}CreateReq) (*types.{{.StructName}}Resp, error) {
	// TODO: 实现创建{{.StructName}}逻辑
	return nil, nil
}
`

// Repository模板
const repositoryTemplate = `package repository

import (
	"context"
	"power-admin/pkg/models"

	"gorm.io/gorm"
)

// {{.StructName}}Repository {{.StructName}}仓储接口
type {{.StructName}}Repository interface {
	Create(ctx context.Context, m *models.{{.StructName}}) error
	Update(ctx context.Context, m *models.{{.StructName}}) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*models.{{.StructName}}, error)
	List(ctx context.Context, page, pageSize int) ([]*models.{{.StructName}}, int64, error)
}

type {{.StructName}}RepositoryImpl struct {
	db *gorm.DB
}

func New{{.StructName}}Repository(db *gorm.DB) {{.StructName}}Repository {
	return &{{.StructName}}RepositoryImpl{db: db}
}

func (r *{{.StructName}}RepositoryImpl) Create(ctx context.Context, m *models.{{.StructName}}) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *{{.StructName}}RepositoryImpl) Update(ctx context.Context, m *models.{{.StructName}}) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *{{.StructName}}RepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.{{.StructName}}{}, id).Error
}

func (r *{{.StructName}}RepositoryImpl) Get(ctx context.Context, id int64) (*models.{{.StructName}}, error) {
	var m models.{{.StructName}}
	err := r.db.WithContext(ctx).First(&m, id).Error
	return &m, err
}

func (r *{{.StructName}}RepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*models.{{.StructName}}, int64, error) {
	var list []*models.{{.StructName}}
	var total int64

	query := r.db.WithContext(ctx).Model(&models.{{.StructName}}{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&list).Error

	return list, total, err
}
`
