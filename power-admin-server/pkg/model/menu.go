package model

import (
	"time"
)

// Menu 菜单表
type Menu struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(50);not null;comment:菜单名称" json:"name"`
	Path      string    `gorm:"type:varchar(100);comment:路由路径" json:"path"`
	Component string    `gorm:"type:varchar(100);comment:组件路径" json:"component"`
	Icon      string    `gorm:"type:varchar(50);comment:图标" json:"icon"`
	ParentId  int64     `gorm:"default:0;comment:父级ID,0为顶级菜单" json:"parentId"`
	Sort      int       `gorm:"default:0;comment:排序" json:"sort"`
	Status    int       `gorm:"type:tinyint;default:1;comment:1启用2禁用" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Menu) TableName() string {
	return "menu"
}
