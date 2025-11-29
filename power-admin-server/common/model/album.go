package model

import (
	"time"
)

// AlbumCategory 相册分类表
type AlbumCategory struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"type:bigint;not null;index:idx_user_id" json:"userId"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	Sort      int       `gorm:"type:int;default:0" json:"sort"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (AlbumCategory) TableName() string {
	return "album_categories"
}

// Album 相册照片表
type Album struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      int64     `gorm:"type:bigint;not null;index:idx_user_id,idx_user_category" json:"userId"`
	CategoryId  int64     `gorm:"type:bigint;index:idx_category_id,idx_user_category" json:"categoryId"`
	PhotoUrl    string    `gorm:"type:varchar(500);not null" json:"photoUrl"`
	Description string    `gorm:"type:varchar(200)" json:"description"`
	Sort        int       `gorm:"type:int;default:0" json:"sort"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Album) TableName() string {
	return "album"
}
