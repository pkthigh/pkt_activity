package model

import "strconv"

// UserLogin 用户登陆日志
type UserLogin struct {
	Id         int    `json:"id"`                                   // ID
	Uid        int    `json:"uid"`                                  // 用户ID
	Cid        int    `json:"cid"`                                  // 客户端编号
	Version    string `json:"version"`                              // 客户端版本
	Channel    int    `json:"channel"`                              // 子渠道号
	Lang       int    `json:"lang"`                                 // 客户端语言
	Network    string `json:"network"`                              // 网络类型(wifi,gsm,4G)
	OsType     string `json:"os_type" gorm:"column:osType"`         // 系统类型(IOS, Android)
	OsVersion  string `json:"os_version" gorm:"column:osVersion"`   // 系统版本号
	DeviceType string `json:"device_type" gorm:"column:deviceType"` // 设备型号(HUAWEI Mate20 ...)
	Imei       string `json:"imei"`                                 // 设备IMEI
	MacAddr    string `json:"mac_addr" gorm:"column:macAddr"`       // 客户端MAC地址
	Ip         string `json:"ip"`                                   // 登录IP
	Time       int    `json:"time"`                                 // 登录时间
}

// TableName ..
func (u *UserLogin) TableName() string {
	return "user_login" + strconv.Itoa(u.Uid%10)
}
