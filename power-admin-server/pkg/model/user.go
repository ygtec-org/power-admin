package model

import (
	"time"
)

// User 用户表
type User struct {
	Id           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Phone        string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	Password     string    `gorm:"type:varchar(255);not null" json:"-"`
	Nickname     string    `gorm:"type:varchar(50)" json:"nickname"`
	Avatar       string    `gorm:"type:varchar(255)" json:"avatar"`
	Gender       int       `gorm:"type:tinyint;default:1;comment:1男2女" json:"gender"`
	Age          int       `gorm:"type:int" json:"age"`
	City         string    `gorm:"type:varchar(50)" json:"city"`
	Job          string    `gorm:"type:varchar(50)" json:"job"`
	Height       int       `gorm:"type:int;comment:身高cm" json:"height"`
	Weight       string    `gorm:"type:varchar(50);comment:体重kg" json:"weight"`
	Education    string    `gorm:"type:varchar(20)" json:"education"`
	Income       string    `gorm:"type:varchar(20)" json:"income"`
	Introduction string    `gorm:"type:text" json:"introduction"`
	Photos       string    `gorm:"type:text;comment:照片JSON数组" json:"photos"`
	Petal        int       `gorm:"type:int;default:0;comment:花瓣余额" json:"petal"`
	Status       int       `gorm:"type:tinyint;default:1;comment:1正常2审核中3已封禁" json:"status"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (User) TableName() string {
	return "user"
}

// Like 喜欢记录表
type Like struct {
	Id           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       int64     `gorm:"index;not null;comment:发起者ID" json:"userId"`
	TargetUserId int64     `gorm:"index;not null;comment:目标用户ID" json:"targetUserId"`
	Type         int       `gorm:"type:tinyint;default:1;comment:1普通喜欢2超级喜欢" json:"type"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (Like) TableName() string {
	return "like"
}

// Match 匹配表
type Match struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId1   int64     `gorm:"index;not null" json:"userId1"`
	UserId2   int64     `gorm:"index;not null" json:"userId2"`
	Status    int       `gorm:"type:tinyint;default:1;comment:1正常2已解除" json:"status"`
	MatchedAt time.Time `gorm:"autoCreateTime" json:"matchedAt"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (Match) TableName() string {
	return "match"
}

// Message 消息表
type Message struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FromUserId  int64     `gorm:"index;not null" json:"fromUserId"`
	ToUserId    int64     `gorm:"index;not null" json:"toUserId"`
	Content     string    `gorm:"type:text" json:"content"`
	MessageType int       `gorm:"type:tinyint;default:1;comment:1文本2图片3语音" json:"messageType"`
	IsRead      int       `gorm:"type:tinyint;default:0;comment:0未读1已读" json:"isRead"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (Message) TableName() string {
	return "message"
}

// Greet 打招呼记录表
type Greet struct {
	Id           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       int64     `gorm:"index;not null;comment:发起者ID" json:"userId"`
	TargetUserId int64     `gorm:"index;not null;comment:目标用户ID" json:"targetUserId"`
	Content      string    `gorm:"type:varchar(500);comment:打招呼内容" json:"content"`
	PetalCost    int       `gorm:"type:int;default:10;comment:花瓣消耗" json:"petalCost"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (Greet) TableName() string {
	return "greet"
}

// Admin 管理员表
type Admin struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Admin) TableName() string {
	return "admin"
}
