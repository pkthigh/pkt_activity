package model

// InsuranceFlowLogV2 保险
type InsuranceFlowLogV2 struct {
	Id             int    `json:"id"`
	CreatedAt      int64  `json:"created_at"`
	UserId         int    `json:"user_id"`
	ClubId         int    `json:"club_id"`
	MatchId        int    `json:"match_id"`
	MatchGameId    int    `json:"match_game_id"`
	MatchType      string `json:"match_type"`
	MatchKind      int    `json:"match_kind"`
	Outs           int    `json:"outs"`
	Buy            int    `json:"buy"`
	Pay            int    `json:"pay"`
	Round          string `json:"round"`
	PotId          int    `json:"pot_id"`
	PotChips       int    `json:"pot_chips"`
	HoleCards      string `json:"hole_cards"`
	Board          string `json:"board"`
	AllInUserIds   string `json:"all_in_user_ids"`
	AllInUserCount int    `json:"all_in_user_count"`
}

func (InsuranceFlowLogV2) TableName() string {
	return "pkc_insurance_flow_log_v2"
}
