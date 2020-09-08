package common

import "pkt_activity/library/storage"

const (
	// ActivityDsn 活动数据库
	ActivityDsn storage.MYSQL = "activity_dsn"
	// DataDsn 流水数据库(只读)
	DataDsn storage.MYSQL = "data_dsn"
	// UserDsn 用户数据库(只读)
	UserDsn storage.MYSQL = "user_dsn"
	// LogDsn 日志数据库(只读)
	LogDsn storage.MYSQL = "log_dsn"
)
