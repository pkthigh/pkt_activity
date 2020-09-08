package common

// SetMoneyArgs 设置金钱参数
type SetMoneyArgs struct {
	OptionType uint32
	Desc       string
	DeltaValue float64
}

// SetMoneyReq 设置金币
type SetMoneyReq struct {
	UID     uint32
	SetArgs []SetMoneyArgs
}

// SetMoneyRsp 设置金币
type SetMoneyRsp struct {
	LeftMoney uint64
}
