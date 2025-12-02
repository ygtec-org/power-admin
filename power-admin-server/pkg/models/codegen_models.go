package models

import "time"

// GenConfig 代码生成配置表
type GenConfig struct {
	ID           int64     `gorm:"primaryKey;comment:配置ID" json:"id"`
	Table        string    `gorm:"column:table_name;size:100;uniqueIndex;not null;comment:表名称" json:"tableName"`
	TablePrefix  string    `gorm:"size:50;comment:表前缀" json:"tablePrefix"`
	BusinessName string    `gorm:"size:100;comment:业务名称" json:"businessName"`
	ModuleName   string    `gorm:"size:100;comment:模块名称" json:"moduleName"`
	PackageName  string    `gorm:"size:200;comment:包路径" json:"packageName"`
	Author       string    `gorm:"size:50;comment:作者" json:"author"`
	Remark       string    `gorm:"size:500;comment:备注" json:"remark"`
	CreatedAt    time.Time `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`

	// 关联字段
	Columns []*GenTableColumn `gorm:"foreignKey:GenConfigID" json:"columns,omitempty"`
}

// TableName 指定表名
func (GenConfig) TableName() string {
	return "gen_config"
}

// GenTableColumn 表字段信息表
type GenTableColumn struct {
	ID            int64     `gorm:"primaryKey;comment:字段ID" json:"id"`
	GenConfigID   int64     `gorm:"not null;index;comment:配置ID" json:"genConfigId"`
	ColumnName    string    `gorm:"size:100;not null;comment:字段名称" json:"columnName"`
	ColumnComment string    `gorm:"size:500;comment:字段注释" json:"columnComment"`
	ColumnType    string    `gorm:"size:100;comment:字段类型(MySQL类型)" json:"columnType"`
	GoType        string    `gorm:"size:100;comment:Go类型" json:"goType"`
	GoField       string    `gorm:"size:100;comment:Go字段名" json:"goField"`
	IsPk          int       `gorm:"default:0;comment:是否主键(1是)" json:"isPk"`
	IsIncrement   int       `gorm:"default:0;comment:是否自增(1是)" json:"isIncrement"`
	IsRequired    int       `gorm:"default:0;comment:是否必填(1是)" json:"isRequired"`
	IsInsert      int       `gorm:"default:1;comment:是否为插入字段(1是)" json:"isInsert"`
	IsEdit        int       `gorm:"default:1;comment:是否编辑字段(1是)" json:"isEdit"`
	IsList        int       `gorm:"default:1;comment:是否列表字段(1是)" json:"isList"`
	IsQuery       int       `gorm:"default:1;comment:是否查询字段(1是)" json:"isQuery"`
	QueryType     string    `gorm:"size:50;default:=;comment:查询方式(=,!=,>,<,LIKE等)" json:"queryType"`
	HtmlType      string    `gorm:"size:50;comment:显示类型(input,textarea,select等)" json:"htmlType"`
	DictType      string    `gorm:"size:100;comment:字典类型" json:"dictType"`
	Sort          int       `gorm:"default:0;index;comment:排序" json:"sort"`
	CreatedAt     time.Time `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
}

// TableName 指定表名
func (GenTableColumn) TableName() string {
	return "gen_table_column"
}

// GenHistory 代码生成历史表
type GenHistory struct {
	ID          int64     `gorm:"primaryKey;comment:历史ID" json:"id"`
	GenConfigID int64     `gorm:"not null;index;comment:配置ID" json:"genConfigId"`
	Table       string    `gorm:"column:table_name;size:100;not null;index;comment:表名称" json:"tableName"`
	FilePath    string    `gorm:"size:500;comment:生成的文件路径" json:"filePath"`
	FileType    string    `gorm:"size:50;comment:文件类型(api,model,handler,logic,repository)" json:"fileType"`
	Content     string    `gorm:"type:longtext;comment:生成的文件内容" json:"content"`
	Status      int       `gorm:"default:1;comment:状态(1成功 0失败)" json:"status"`
	ErrorMsg    string    `gorm:"type:text;comment:错误信息" json:"errorMsg"`
	Operator    string    `gorm:"size:50;comment:操作人" json:"operator"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index;comment:生成时间" json:"createdAt"`
}

// TableName 指定表名
func (GenHistory) TableName() string {
	return "gen_history"
}
