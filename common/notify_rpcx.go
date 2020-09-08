package common

// ActivityNotifyArgs 活动通知参数
type ActivityNotifyArgs struct {
	Uid  uint32
	Info *PBActivityInfo
}

// PBActivityInfo 活动信息
type PBActivityInfo struct {
	ActivityId   *uint32 `protobuf:"varint,1,opt,name=activityId" json:"activityId,omitempty"`
	ActivityType *uint32 `protobuf:"varint,2,opt,name=activityType" json:"activityType,omitempty"`
	// HandsProgress 手数活动进度
	HandsProgress *PBActivityInfo_ActivityHandsProgress `protobuf:"bytes,3,opt,name=handsProgress" json:"handsProgress,omitempty"`
	// XXXProgress 以后其他活动进度修改这个结构体(游戏服务器规定)
	// ...
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

// PBActivityInfo_ActivityHandsProgress 手数活动进程
type PBActivityInfo_ActivityHandsProgress struct {
	// 玩家当前手数
	PlayedHandsNum *uint32 `protobuf:"varint,1,opt,name=playedHandsNum" json:"playedHandsNum,omitempty"`
	// 活动需要手数
	DrawHandsNum *uint32 `protobuf:"varint,2,opt,name=drawHandsNum" json:"drawHandsNum,omitempty"`
	// 是否满足条件
	CanDraw              *bool    `protobuf:"varint,3,opt,name=canDraw" json:"canDraw,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
