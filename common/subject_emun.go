package common

// Subject 订阅常量
type Subject string

const (
	// FishHandOverSubject 钓鱼每手结束信息
	FishHandOverSubject Subject = "fish_hand_over"
	// PkcHandOverSubject 微币每手结束信息
	PkcHandOverSubject Subject = "pkc_hand_over"
	// ItemSubject 道具使用
	ItemSubject Subject = "item"
)

func (subject Subject) String() string {
	return string(subject)
}
