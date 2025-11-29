package model

import (
	"time"
)

// Post 动态表
type Post struct {
	Id           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       int64     `gorm:"index;not null;comment:用户ID" json:"userId"`
	Content      string    `gorm:"type:text;not null;comment:动态内容" json:"content"`
	Images       string    `gorm:"type:text;comment:图片URL JSON数组" json:"images"`
	Location     string    `gorm:"type:varchar(100);comment:位置信息" json:"location"`
	Status       int       `gorm:"type:tinyint;default:1;comment:1正常2已删除" json:"status"`
	LikeCount    int       `gorm:"type:int;default:0;comment:点赞数" json:"likeCount"`
	CommentCount int       `gorm:"type:int;default:0;comment:评论数" json:"commentCount"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Post) TableName() string {
	return "post"
}

// Comment 评论表
type Comment struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PostId    int64     `gorm:"index;not null;comment:动态ID" json:"postId"`
	UserId    int64     `gorm:"index;not null;comment:用户ID" json:"userId"`
	Content   string    `gorm:"type:text;not null;comment:评论内容" json:"content"`
	Status    int       `gorm:"type:tinyint;default:1;comment:1正常2已删除" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Comment) TableName() string {
	return "comment"
}

// LikePost 动态点赞表
type LikePost struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PostId    int64     `gorm:"index;not null;comment:动态ID" json:"postId"`
	UserId    int64     `gorm:"index;not null;comment:用户ID" json:"userId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (LikePost) TableName() string {
	return "like_post"
}

// Tag 标签表
type Tag struct {
	Id         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(50);not null;uniqueIndex;comment:标签名称" json:"name"`
	UsageCount int       `gorm:"type:int;default:0;comment:使用次数" json:"usageCount"`
	Sort       int       `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status     int       `gorm:"type:tinyint;default:1;comment:1启用0禁用" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Tag) TableName() string {
	return "tag"
}

// SensitiveWord 敏感词表
type SensitiveWord struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Word      string    `gorm:"type:varchar(100);not null;uniqueIndex;comment:敏感词" json:"word"`
	Level     int       `gorm:"type:tinyint;default:1;comment:级别:1警告2严重" json:"level"`
	HitCount  int       `gorm:"type:int;default:0;comment:命中次数" json:"hitCount"`
	Status    int       `gorm:"type:tinyint;default:1;comment:1启用0禁用" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (SensitiveWord) TableName() string {
	return "sensitive_word"
}

// Dict 字典表
type Dict struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DictKey     string    `gorm:"type:varchar(100);not null;comment:字典键" json:"dictKey"`
	DictValue   string    `gorm:"type:varchar(255);not null;comment:字典值" json:"dictValue"`
	DictType    string    `gorm:"type:varchar(50);not null;index;comment:字典类型" json:"dictType"`
	Description string    `gorm:"type:varchar(255);comment:描述" json:"description"`
	Sort        int       `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status      int       `gorm:"type:tinyint;default:1;comment:1启用0禁用" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Dict) TableName() string {
	return "dict"
}

// Help 帮助表
type Help struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"type:varchar(200);not null;comment:标题" json:"title"`
	Category  string    `gorm:"type:varchar(50);comment:分类" json:"category"`
	Content   string    `gorm:"type:text;not null;comment:内容" json:"content"`
	Sort      int       `gorm:"type:int;default:0;comment:排序" json:"sort"`
	ViewCount int       `gorm:"type:int;default:0;comment:浏览次数" json:"viewCount"`
	Status    int       `gorm:"type:tinyint;default:1;comment:1显示0隐藏" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Help) TableName() string {
	return "help"
}

// Protocol 协议表
type Protocol struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"type:varchar(200);not null;comment:协议标题" json:"title"`
	Type      int       `gorm:"type:tinyint;not null;comment:协议类型:1用户协议2隐私政策3服务条款4其他" json:"type"`
	Version   string    `gorm:"type:varchar(20);not null;comment:版本号" json:"version"`
	Content   string    `gorm:"type:longtext;not null;comment:协议内容" json:"content"`
	Status    int       `gorm:"type:tinyint;default:1;comment:1生效0停用" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Protocol) TableName() string {
	return "protocol"
}

// FrontendConfig 前台配置表
type FrontendConfig struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string    `gorm:"type:varchar(100);not null;uniqueIndex;comment:配置键" json:"configKey"`
	ConfigValue string    `gorm:"type:text;not null;comment:配置值" json:"configValue"`
	ConfigType  int       `gorm:"type:tinyint;not null;comment:配置类型:1文本2图片3JSON4数字" json:"configType"`
	Description string    `gorm:"type:varchar(255);comment:描述" json:"description"`
	Status      int       `gorm:"type:tinyint;default:1;comment:1启用0禁用" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (FrontendConfig) TableName() string {
	return "frontend_config"
}

// ActivityNotice 活动公告表
type ActivityNotice struct {
	Id         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"type:varchar(200);not null;comment:公告标题" json:"title"`
	Content    string    `gorm:"type:text;not null;comment:公告内容" json:"content"`
	CoverImage string    `gorm:"type:varchar(500);comment:封面图片" json:"coverImage"`
	StartTime  time.Time `gorm:"comment:开始时间" json:"startTime"`
	EndTime    time.Time `gorm:"comment:结束时间" json:"endTime"`
	Sort       int       `gorm:"type:int;default:0;comment:排序" json:"sort"`
	ViewCount  int       `gorm:"type:int;default:0;comment:浏览次数" json:"viewCount"`
	Status     int       `gorm:"type:tinyint;default:1;comment:1显示0隐藏" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (ActivityNotice) TableName() string {
	return "activity_notice"
}

// Partner 搭子表
type Partner struct {
	Id            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId        int64     `gorm:"index;not null;comment:用户ID" json:"userId"`
	Title         string    `gorm:"type:varchar(200);not null;comment:搭子标题" json:"title"`
	Description   string    `gorm:"type:text;comment:描述" json:"description"`
	ActivityType  string    `gorm:"type:varchar(50);comment:活动类型" json:"activityType"`
	CityId        int64     `gorm:"index;comment:城市ID" json:"cityId"`
	Location      string    `gorm:"type:varchar(200);comment:位置" json:"location"`
	ActivityTime  time.Time `gorm:"comment:活动时间" json:"activityTime"`
	MaxPeople     int       `gorm:"type:int;comment:最多人数" json:"maxPeople"`
	CurrentPeople int       `gorm:"type:int;default:0;comment:当前人数" json:"currentPeople"`
	Status        int       `gorm:"type:tinyint;default:1;comment:1招募中2已满3已结束" json:"status"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Partner) TableName() string {
	return "partner"
}

// PartnerEvent 搭子活动表
type PartnerEvent struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PartnerId int64     `gorm:"index;not null;comment:搭子ID" json:"partnerId"`
	UserId    int64     `gorm:"index;not null;comment:用户ID" json:"userId"`
	Status    int       `gorm:"type:tinyint;default:1;comment:1待确认2已确认3已拒绝" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (PartnerEvent) TableName() string {
	return "partner_event"
}

// City 城市表
type City struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(50);not null;comment:城市名称" json:"name"`
	Province  string    `gorm:"type:varchar(50);comment:省份" json:"province"`
	Code      string    `gorm:"type:varchar(20);comment:城市编码" json:"code"`
	Sort      int       `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status    int       `gorm:"type:tinyint;default:1;comment:1启用0禁用" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (City) TableName() string {
	return "city"
}

// RechargeConfig 充值配置表
type RechargeConfig struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string    `gorm:"type:varchar(100);not null;uniqueIndex;comment:配置键" json:"configKey"`
	ConfigValue string    `gorm:"type:text;not null;comment:配置值" json:"configValue"`
	Description string    `gorm:"type:varchar(255);comment:描述" json:"description"`
	Status      int       `gorm:"type:tinyint;default:1;comment:1启用0禁用" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (RechargeConfig) TableName() string {
	return "recharge_config"
}

// MemberPackage 会员套餐表
type MemberPackage struct {
	Id            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"type:varchar(100);not null;comment:套餐名称" json:"name"`
	Days          int       `gorm:"type:int;not null;comment:天数" json:"days"`
	Price         float64   `gorm:"type:decimal(10,2);not null;comment:价格" json:"price"`
	OriginalPrice float64   `gorm:"type:decimal(10,2);comment:原价" json:"originalPrice"`
	Description   string    `gorm:"type:text;comment:描述" json:"description"`
	Sort          int       `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status        int       `gorm:"type:tinyint;default:1;comment:1上架0下架" json:"status"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (MemberPackage) TableName() string {
	return "member_package"
}

// PetalPackage 花瓣套餐表
type PetalPackage struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;comment:套餐名称" json:"name"`
	Amount      int       `gorm:"type:int;not null;comment:花瓣数量" json:"amount"`
	Price       float64   `gorm:"type:decimal(10,2);not null;comment:价格" json:"price"`
	ExtraAmount int       `gorm:"type:int;default:0;comment:额外赠送" json:"extraAmount"`
	Description string    `gorm:"type:text;comment:描述" json:"description"`
	Sort        int       `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status      int       `gorm:"type:tinyint;default:1;comment:1上架0下架" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (PetalPackage) TableName() string {
	return "petal_package"
}

// CashExchangeRule 现金兑换规则表
type CashExchangeRule struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PetalAmount int       `gorm:"type:int;not null;comment:花瓣数量" json:"petalAmount"`
	CashAmount  float64   `gorm:"type:decimal(10,2);not null;comment:现金金额" json:"cashAmount"`
	MinWithdraw float64   `gorm:"type:decimal(10,2);comment:最低提现" json:"minWithdraw"`
	Description string    `gorm:"type:text;comment:描述" json:"description"`
	Status      int       `gorm:"type:tinyint;default:1;comment:1启用0禁用" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (CashExchangeRule) TableName() string {
	return "cash_exchange_rule"
}

// Gift 礼物表
type Gift struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;comment:礼物名称" json:"name"`
	Icon        string    `gorm:"type:varchar(500);comment:礼物图标" json:"icon"`
	Price       int       `gorm:"type:int;not null;comment:花瓣价格" json:"price"`
	Category    string    `gorm:"type:varchar(50);comment:分类" json:"category"`
	Description string    `gorm:"type:text;comment:描述" json:"description"`
	Sort        int       `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status      int       `gorm:"type:tinyint;default:1;comment:1上架0下架" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Gift) TableName() string {
	return "gift"
}

// RechargeOrder 充值订单表
type RechargeOrder struct {
	Id        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo   string     `gorm:"type:varchar(50);not null;uniqueIndex;comment:订单号" json:"orderNo"`
	UserId    int64      `gorm:"index;not null;comment:用户ID" json:"userId"`
	Type      int        `gorm:"type:tinyint;not null;comment:充值类型:1会员2花瓣" json:"type"`
	PackageId int64      `gorm:"comment:套餐ID" json:"packageId"`
	Amount    float64    `gorm:"type:decimal(10,2);not null;comment:金额" json:"amount"`
	PayMethod int        `gorm:"type:tinyint;comment:支付方式:1微信2支付宝" json:"payMethod"`
	Status    int        `gorm:"type:tinyint;default:0;comment:0待支付1已支付2已取消3已退款" json:"status"`
	PayTime   *time.Time `gorm:"comment:支付时间" json:"payTime"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (RechargeOrder) TableName() string {
	return "recharge_order"
}

// PetalRecord 花瓣记录表
type PetalRecord struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"index;not null;comment:用户ID" json:"userId"`
	Type      int       `gorm:"type:tinyint;not null;comment:类型:1收入2支出" json:"type"`
	Amount    int       `gorm:"type:int;not null;comment:数量" json:"amount"`
	Reason    string    `gorm:"type:varchar(200);comment:原因" json:"reason"`
	RelatedId int64     `gorm:"comment:关联ID" json:"relatedId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (PetalRecord) TableName() string {
	return "petal_record"
}

// RefundOrder 退款订单表
type RefundOrder struct {
	Id              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo         string     `gorm:"type:varchar(50);not null;uniqueIndex;comment:退款单号" json:"orderNo"`
	RechargeOrderId int64      `gorm:"index;not null;comment:充值订单ID" json:"rechargeOrderId"`
	UserId          int64      `gorm:"index;not null;comment:用户ID" json:"userId"`
	Amount          float64    `gorm:"type:decimal(10,2);not null;comment:退款金额" json:"amount"`
	Reason          string     `gorm:"type:varchar(500);comment:退款原因" json:"reason"`
	Status          int        `gorm:"type:tinyint;default:0;comment:0待审核1已同意2已拒绝3已完成" json:"status"`
	RefundTime      *time.Time `gorm:"comment:退款时间" json:"refundTime"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (RefundOrder) TableName() string {
	return "refund_order"
}

// WithdrawOrder 提现订单表
type WithdrawOrder struct {
	Id           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo      string     `gorm:"type:varchar(50);not null;uniqueIndex;comment:提现单号" json:"orderNo"`
	UserId       int64      `gorm:"index;not null;comment:用户ID" json:"userId"`
	Amount       float64    `gorm:"type:decimal(10,2);not null;comment:提现金额" json:"amount"`
	PetalAmount  int        `gorm:"type:int;not null;comment:花瓣数量" json:"petalAmount"`
	BankInfo     string     `gorm:"type:text;comment:银行信息" json:"bankInfo"`
	Status       int        `gorm:"type:tinyint;default:0;comment:0待审核1已同意2已拒绝3已完成" json:"status"`
	WithdrawTime *time.Time `gorm:"comment:提现时间" json:"withdrawTime"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (WithdrawOrder) TableName() string {
	return "withdraw_order"
}

// Report 举报表
type Report struct {
	Id           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       int64      `gorm:"index;not null;comment:举报人ID" json:"userId"`
	TargetId     int64      `gorm:"index;not null;comment:被举报对象ID" json:"targetId"`
	TargetType   int        `gorm:"type:tinyint;not null;comment:举报类型:1用户2动态3评论" json:"targetType"`
	Reason       string     `gorm:"type:varchar(500);comment:举报原因" json:"reason"`
	Images       string     `gorm:"type:text;comment:证据图片" json:"images"`
	Status       int        `gorm:"type:tinyint;default:0;comment:0待处理1已处理2已驳回" json:"status"`
	HandleResult string     `gorm:"type:text;comment:处理结果" json:"handleResult"`
	HandleTime   *time.Time `gorm:"comment:处理时间" json:"handleTime"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Report) TableName() string {
	return "report"
}
