package common

import "pkt_activity/library/storage"

const (
	// UserTokenStore 用户Token
	UserTokenStore storage.AREA = 0

	// ItemRecordStore 道具记录存储
	ItemRecordStore storage.AREA = 1
	// TexasHandOverRecordStore 德州每手记录存储
	TexasHandOverRecordStore storage.AREA = 2
	// InsuranceRecordStore 保险记录存储
	InsuranceRecordStore storage.AREA = 3
	// FishingHandOverRecordStore 钓鱼每手记录存储
	FishingHandOverRecordStore storage.AREA = 4
	// FishingUserHandOverRecordStore 钓鱼用户每手记录存储
	FishingUserHandOverRecordStore storage.AREA = 5

	// UserActivityResultCache 用户每日活动结果缓存
	UserActivityResultCache storage.AREA = 10
	// UserDeviceIDResultCache 用户抽奖结果设备ID缓存
	UserDeviceIDResultCache storage.AREA = 11

	// UpdateNotification 更新通知 & 玩家每日是否已推送
	UpdateNotification storage.AREA = 15
)
