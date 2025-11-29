package model

import (
	"time"
)

// Activity 活动表
type Activity struct {
	Id              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string    `gorm:"type:varchar(100);not null;comment:活动标题" json:"title"`
	Description     string    `gorm:"type:text;comment:活动描述" json:"description"`
	CoverImage      string    `gorm:"type:varchar(255);comment:封面图片" json:"coverImage"`
	City            string    `gorm:"type:varchar(50);comment:城市" json:"city"`
	Location        string    `gorm:"type:varchar(100);comment:详细地址" json:"location"`
	StartTime       time.Time `gorm:"comment:开始时间" json:"startTime"`
	EndTime         time.Time `gorm:"comment:结束时间" json:"endTime"`
	MaxParticipants int       `gorm:"type:int;default:0;comment:最大参与人数,0为无限制" json:"maxParticipants"`
	SignupCount     int       `gorm:"type:int;default:0;comment:已报名人数" json:"signupCount"`
	Status          int       `gorm:"type:tinyint;default:1;comment:1草稿2已发布3已结束4已取消" json:"status"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Activity) TableName() string {
	return "activity"
}

// ActivitySignup 活动报名表
type ActivitySignup struct {
	Id         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityId int64     `gorm:"index;not null;comment:活动ID" json:"activityId"`
	UserId     int64     `gorm:"index;not null;comment:用户ID" json:"userId"`
	Name       string    `gorm:"type:varchar(50);not null;comment:姓名" json:"name"`
	Phone      string    `gorm:"type:varchar(20);not null;comment:手机号" json:"phone"`
	Gender     int       `gorm:"type:tinyint;comment:性别:1男2女" json:"gender"`
	Age        int       `gorm:"type:int;comment:年龄" json:"age"`
	Job        string    `gorm:"type:varchar(50);comment:职业" json:"job"`
	Remark     string    `gorm:"type:text;comment:备注" json:"remark"`
	Status     int       `gorm:"type:tinyint;default:1;comment:1待审核2已通过3已拒绝" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (ActivitySignup) TableName() string {
	return "activity_signup"
}

// Topic 话题表
type Topic struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"type:varchar(100);not null;uniqueIndex;comment:话题标题" json:"title"`
	Description string    `gorm:"type:text;comment:话题描述" json:"description"`
	PostCount   int       `gorm:"type:int;default:0;comment:动态数" json:"postCount"`
	FollowCount int       `gorm:"type:int;default:0;comment:关注数" json:"followCount"`
	Sort        int       `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status      int       `gorm:"type:tinyint;default:1;comment:1启用2禁用" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Topic) TableName() string {
	return "topic"
}
