package repository

import (
	"context"
	"power-admin-server/pkg/models"
	"strings"

	"gorm.io/gorm"
)

// CodegenRepository 代码生成器仓储接口
type CodegenRepository interface {
	// 配置管理
	CreateConfig(ctx context.Context, config *models.GenConfig) error
	UpdateConfig(ctx context.Context, config *models.GenConfig) error
	DeleteConfig(ctx context.Context, id int64) error
	GetConfig(ctx context.Context, id int64) (*models.GenConfig, error)
	ListConfig(ctx context.Context, page, pageSize int, tableName string) ([]*models.GenConfig, int64, error)
	GetConfigByTableName(ctx context.Context, tableName string) (*models.GenConfig, error)

	// 字段管理
	CreateColumns(ctx context.Context, columns []*models.GenTableColumn) error
	UpdateColumns(ctx context.Context, columns []*models.GenTableColumn) error
	DeleteColumnsByConfigID(ctx context.Context, configID int64) error
	GetColumnsByConfigID(ctx context.Context, configID int64) ([]*models.GenTableColumn, error)

	// 历史管理
	CreateHistory(ctx context.Context, history *models.GenHistory) error
	CreateHistories(ctx context.Context, histories []*models.GenHistory) error
	ListHistory(ctx context.Context, page, pageSize int, tableName string) ([]*models.GenHistory, int64, error)
	GetHistory(ctx context.Context, id int64) (*models.GenHistory, error)
	DeleteHistory(ctx context.Context, id int64) error

	// 数据库表信息查询
	GetDatabaseTables(ctx context.Context, tableName string) ([]*DatabaseTableInfo, error)
	GetTableColumns(ctx context.Context, tableName string) ([]*TableColumnInfo, error)
}

// DatabaseTableInfo 数据库表信息
type DatabaseTableInfo struct {
	TableName      string
	TableComment   string
	Engine         string
	TableCollation string
}

// TableColumnInfo 表字段信息
type TableColumnInfo struct {
	ColumnName      string
	DataType        string
	ColumnType      string
	ColumnComment   string
	IsNullable      string
	ColumnKey       string
	Extra           string
	OrdinalPosition int
}

// codegenRepositoryImpl 代码生成器仓储实现
type codegenRepositoryImpl struct {
	db *gorm.DB
}

// NewCodegenRepository 创建代码生成器仓储实例
func NewCodegenRepository(db *gorm.DB) CodegenRepository {
	return &codegenRepositoryImpl{db: db}
}

// ==================== 配置管理 ====================

func (r *codegenRepositoryImpl) CreateConfig(ctx context.Context, config *models.GenConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

func (r *codegenRepositoryImpl) UpdateConfig(ctx context.Context, config *models.GenConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

func (r *codegenRepositoryImpl) DeleteConfig(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除配置
		if err := tx.Delete(&models.GenConfig{}, id).Error; err != nil {
			return err
		}
		// 删除关联的字段（级联删除会自动处理）
		// 删除关联的历史（级联删除会自动处理）
		return nil
	})
}

func (r *codegenRepositoryImpl) GetConfig(ctx context.Context, id int64) (*models.GenConfig, error) {
	var config models.GenConfig
	err := r.db.WithContext(ctx).
		Preload("Columns", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort ASC")
		}).
		First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *codegenRepositoryImpl) ListConfig(ctx context.Context, page, pageSize int, tableName string) ([]*models.GenConfig, int64, error) {
	var configs []*models.GenConfig
	var total int64

	query := r.db.WithContext(ctx).Model(&models.GenConfig{})

	// 条件查询
	if tableName != "" {
		query = query.Where("table_name LIKE ?", "%"+tableName+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&configs).Error

	return configs, total, err
}

func (r *codegenRepositoryImpl) GetConfigByTableName(ctx context.Context, tableName string) (*models.GenConfig, error) {
	var config models.GenConfig
	err := r.db.WithContext(ctx).
		Where("table_name = ?", tableName).
		Preload("Columns", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort ASC")
		}).
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// ==================== 字段管理 ====================

func (r *codegenRepositoryImpl) CreateColumns(ctx context.Context, columns []*models.GenTableColumn) error {
	if len(columns) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&columns).Error
}

func (r *codegenRepositoryImpl) UpdateColumns(ctx context.Context, columns []*models.GenTableColumn) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, col := range columns {
			if err := tx.Save(col).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *codegenRepositoryImpl) DeleteColumnsByConfigID(ctx context.Context, configID int64) error {
	return r.db.WithContext(ctx).
		Where("gen_config_id = ?", configID).
		Delete(&models.GenTableColumn{}).Error
}

func (r *codegenRepositoryImpl) GetColumnsByConfigID(ctx context.Context, configID int64) ([]*models.GenTableColumn, error) {
	var columns []*models.GenTableColumn
	err := r.db.WithContext(ctx).
		Where("gen_config_id = ?", configID).
		Order("sort ASC").
		Find(&columns).Error
	return columns, err
}

// ==================== 历史管理 ====================

func (r *codegenRepositoryImpl) CreateHistory(ctx context.Context, history *models.GenHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *codegenRepositoryImpl) CreateHistories(ctx context.Context, histories []*models.GenHistory) error {
	if len(histories) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&histories).Error
}

func (r *codegenRepositoryImpl) ListHistory(ctx context.Context, page, pageSize int, tableName string) ([]*models.GenHistory, int64, error) {
	var histories []*models.GenHistory
	var total int64

	query := r.db.WithContext(ctx).Model(&models.GenHistory{})

	// 条件查询
	if tableName != "" {
		query = query.Where("table_name LIKE ?", "%"+tableName+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&histories).Error

	return histories, total, err
}

func (r *codegenRepositoryImpl) GetHistory(ctx context.Context, id int64) (*models.GenHistory, error) {
	var history models.GenHistory
	err := r.db.WithContext(ctx).First(&history, id).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *codegenRepositoryImpl) DeleteHistory(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.GenHistory{}, id).Error
}

// ==================== 数据库表信息查询 ====================

func (r *codegenRepositoryImpl) GetDatabaseTables(ctx context.Context, tableName string) ([]*DatabaseTableInfo, error) {
	var tables []*DatabaseTableInfo

	query := `
		SELECT 
			table_name AS TableName,
			table_comment AS TableComment,
			engine AS Engine,
			table_collation AS TableCollation
		FROM information_schema.tables 
		WHERE table_schema = DATABASE()
	`

	if tableName != "" {
		query += " AND table_name LIKE ?"
		err := r.db.WithContext(ctx).Raw(query, "%"+tableName+"%").Scan(&tables).Error
		return tables, err
	}

	err := r.db.WithContext(ctx).Raw(query).Scan(&tables).Error
	return tables, err
}

func (r *codegenRepositoryImpl) GetTableColumns(ctx context.Context, tableName string) ([]*TableColumnInfo, error) {
	var columns []*TableColumnInfo

	query := `
		SELECT 
			column_name AS ColumnName,
			data_type AS DataType,
			column_type AS ColumnType,
			column_comment AS ColumnComment,
			is_nullable AS IsNullable,
			column_key AS ColumnKey,
			extra AS Extra,
			ordinal_position AS OrdinalPosition
		FROM information_schema.columns 
		WHERE table_schema = DATABASE() 
		AND table_name = ?
		ORDER BY ordinal_position
	`

	err := r.db.WithContext(ctx).Raw(query, tableName).Scan(&columns).Error
	return columns, err
}

// ==================== 辅助方法 ====================

// ParseTablePrefix 解析表前缀
func ParseTablePrefix(tableName, prefix string) string {
	if prefix == "" {
		return tableName
	}
	return strings.TrimPrefix(tableName, prefix+"_")
}

// TableNameToStructName 表名转结构体名
func TableNameToStructName(tableName, prefix string) string {
	name := ParseTablePrefix(tableName, prefix)
	parts := strings.Split(name, "_")
	for i, part := range parts {
		if part != "" {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// ColumnNameToGoField 字段名转Go字段名
func ColumnNameToGoField(columnName string) string {
	parts := strings.Split(columnName, "_")
	for i, part := range parts {
		if part != "" {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// MySQLTypeToGoType MySQL类型转Go类型
func MySQLTypeToGoType(mysqlType string) string {
	mysqlType = strings.ToLower(mysqlType)

	typeMap := map[string]string{
		"tinyint":    "int",
		"smallint":   "int",
		"mediumint":  "int",
		"int":        "int64",
		"integer":    "int64",
		"bigint":     "int64",
		"float":      "float32",
		"double":     "float64",
		"decimal":    "float64",
		"char":       "string",
		"varchar":    "string",
		"tinytext":   "string",
		"text":       "string",
		"mediumtext": "string",
		"longtext":   "string",
		"date":       "time.Time",
		"datetime":   "time.Time",
		"timestamp":  "time.Time",
		"time":       "string",
		"year":       "int",
		"blob":       "[]byte",
		"tinyblob":   "[]byte",
		"mediumblob": "[]byte",
		"longblob":   "[]byte",
		"json":       "string",
	}

	for key, value := range typeMap {
		if strings.Contains(mysqlType, key) {
			return value
		}
	}

	return "string"
}

// GetHtmlType 根据Go类型获取HTML表单类型
func GetHtmlType(goType string) string {
	switch goType {
	case "int", "int64", "float32", "float64":
		return "input"
	case "string":
		return "input"
	case "time.Time":
		return "datetime"
	case "bool":
		return "radio"
	default:
		return "input"
	}
}
