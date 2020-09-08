package model

// AcActivity 活动
type AcActivity struct {
	ID int64 `json:"id"`

	CreateAt  int64 `json:"create_at"` // 创建时间
	UpdateAt  int64 `json:"update_at"` // 更新时间
	CreatedBy int64 `json:"-"`         // 创建人员
	UpdatedBy int64 `json:"-"`         // 更新人员

	OnlineTime  int `json:"online_time"`  // 活动开始时间
	OfflineTime int `json:"offline_time"` // 活动结束时间

	NameZh    string `json:"name_zh"`    // 活动名字 中文
	ContentZh string `json:"content_zh"` // 活动内容 中文
	PicZhURL  string `json:"pic_zh_url"` // 活动海报 中文

	NameEn    string `json:"name_en"`    // 活动名字 英文
	ContentEn string `json:"content_en"` // 活动内容 英文
	PicEnURL  string `json:"pic_en_url"` // 活动海报 英文

	AcType int `json:"ac_type"` // [0:其他 / 1:手数活动]
	Status int `json:"status"`  // [0:未开启 / 1:排期中 / 2:进行中 / 3:提前下线 4: 已结束]

	PageURL string `json:"page_url"` // 活动页面链接
	HandNum int    `json:"hand_num"` // 活动门槛手数

	Remark string `json:"remark"` // 备注
}

// TableName ...
func (AcActivity) TableName() string {
	return "ac_activity"
}

// AcActivityDetail 活动详情
type AcActivityDetail struct {
	ID         int64 `json:"id"`
	ActivityID int64 `json:"activity_id"` // 关联活动ID

	CreateAt  int   `json:"-"` // 创建时间
	UpdateAt  int   `json:"-"` // 更新时间
	CreatedBy int64 `json:"-"` // 创建人员
	UpdatedBy int64 `json:"-"` // 更新人员

	Index         int     `json:"index"` // 条件顺序
	Prop          int64   `json:"-"`     // 近7日道具费
	Insurance     int64   `json:"-"`     // 近7日保险费
	Bonus         int64   `json:"bonus"` // 礼包奖金
	BonusNum      int     `json:"-"`     // 每日礼包份数
	BonusRemainng int     `json:"-"`     // 剩余礼包份数
	Odds          float32 `json:"-"`     // 中奖概率
	BonusAll      int64   `json:"-"`     // 总奖金
}

// TableName ...
func (AcActivityDetail) TableName() string {
	return "ac_activity_detail"
}

// AcActivityRecord 活动记录
type AcActivityRecord struct {
	ID               int64  `json:"id"`
	CreateAt         int64  `json:"create_at"`          // 创建时间
	UpdateAt         int64  `json:"update_at"`          // 更新时间
	PlayerID         int64  `json:"player_id"`          // 玩家标识
	PlayerName       string `json:"player_name"`        // 玩家标识
	RaffleTime       int64  `json:"raffle_time"`        // 抽奖时间
	ActivityID       int64  `json:"activity_id"`        // 活动标识
	ActivityDetailID int64  `json:"activity_detail_id"` // 活动详情标识
	Bonus            int64  `json:"bonus"`              // 奖金
	Status           bool   `json:"status"`             // 派奖状态[0:未派发 1:已派发]
	DeviceID         string `json:"device_id"`          // 设备ID
	Lang             bool   `json:"lang"`               // 语言[0:中文 1:英文]
}

// TableName ..
func (AcActivityRecord) TableName() string {
	return "ac_raffle_record"
}
