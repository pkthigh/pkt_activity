package model

import "strconv"

// UserInfo 用户信息
type UserInfo struct {
	Uid     int    `gorm:"primary_key"`
	Name    string `gorm:"type:varchar(64);not null"`
	Sex     int    `gorm:"type:tinyint(3);not null"`
	Head    string `gorm:"type:text;not null"`
	Vip     int    `gorm:"type:tinyint(3);not null"`
	Viptime int    `gorm:"type:int(10);not null"`
	Area    int    `gorm:"type:int(10);not null"`
	Brief   string `gorm:"type:varchar(256);not null"`
	Atoken  string `gorm:"type:varchar(256);not null"`
	Cid     int    `gorm:"type:smallint(5);not null"`
	Ctime   int64  `gorm:"type:int(10);not null"`
	Atime   int    `gorm:"type:int(10);not null"`
	Status  int    `gorm:"type:tinyint(3);not null"`
	Ip      int    `gorm:"type:int(10);not null"`
}

// TableName ..
func (u *UserInfo) TableName() string {
	return "user_info" + strconv.Itoa(u.Uid%10)
}
