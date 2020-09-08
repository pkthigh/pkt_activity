package common

import "fmt"

// Errs 错误定义
type Errs string

const (
	// Successful 成功
	Successful Errs = "Successful"
	// ErrParamInvalid 参数不合法
	ErrParamInvalid Errs = "ErrParamInvalid"
	// ErrPermissions 权限认证失败
	ErrPermissions Errs = "ErrPermissions"
	// ErrActivityCompleted 已经完成活动
	ErrActivityCompleted Errs = "ErrActivityCompleted"
	// ErrNoOngoingActivity 没有正在进行中的活动
	ErrNoOngoingActivity Errs = "ErrNoOngoingActivity"
	// ErrVerifyNoPass 不符合抽奖校验
	ErrVerifyNoPass Errs = "ErrVerifyNoPass"
	// ErrDeviceBeenRaffled 该设备已经抽过奖
	ErrDeviceBeenRaffled Errs = "ErrDeviceBeenRaffled"
	// ErrNoMatchActivityDetail 没有匹配到活动条件
	ErrNoMatchActivityDetail Errs = "ErrNoMatchActivityDetail"
	// ErrRepeatPushNotifications 重复推送通知
	ErrRepeatPushNotifications Errs = "ErrRepeatPushNotifications"
	// ErrLotteryFailed 抽奖失败
	ErrLotteryFailed Errs = "ErrLotteryFailed"
	// ErrTypeConversion 类型转换错误
	ErrTypeConversion Errs = "ErrTypeConversion"
	// ErrRedisQuery Redis查询错误
	ErrRedisQuery Errs = "ErrRedisQuery"
	// ErrMysqlQuery Mysql查询错误
	ErrMysqlQuery Errs = "ErrMysqlQuery"
	// ErrUpdateUserAssets 更新用户资产失败
	ErrUpdateUserAssets Errs = "ErrUpdateUserAssets"
	// ErrUpdateActivityRecoed 更新用户活动记录失败
	ErrUpdateActivityRecoed Errs = "ErrUpdateActivityRecoed"
	// ErrSaveLotteryResult Mysql保存抽奖结果失败
	ErrSaveLotteryResult Errs = "ErrSaveLotteryResult"
)

func (errs Errs) String() string {
	return string(errs)
}

func (errs Errs) Error() error {
	return fmt.Errorf("%v", errs)
}
