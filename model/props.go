package model

// CheckPropsFlow 道具流水
type CheckPropsFlow struct {
	Id            int     `json:"id"`
	UserId        int     `json:"user_id"`         // 用户UID
	Type          string  `json:"type"`            // 类型(interactive_props, delay_thinking, show_hands_cards, show_public_cards)
	PropId        int     `json:"prop_id"`         // 道具ID
	Mid           int     `json:"mid"`             // 牌局MID
	CreateUserId  int     `json:"create_user_id"`  // 创建自建局的用户UID
	CreateClubId  int     `json:"create_club_id"`  // 创建俱乐部局的俱乐部ID
	CreateUnionId int     `json:"create_union_id"` // 创建联盟局的联盟ID
	Num           float64 `json:"num"`             // 使用道具所用的时间币
	CreatedAt     int     `json:"created_at"`      // 道具费产生的时间戳
}

// TableName ..
func (CheckPropsFlow) TableName() string {
	return "pkc_check_props_flow"
}
