package common

// RetObj 返回对象
type RetObj struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// RetObjSuccessful 成功
func RetObjSuccessful(data interface{}) RetObj {
	return RetObj{
		Code: 1,
		Data: data,
		Msg:  Successful.String(),
	}
}

// RetObjParamFailure 参数失败
func RetObjParamFailure(err Errs) RetObj {
	return RetObj{
		Code: 100,
		Data: nil,
		Msg:  err.String(),
	}
}

// RetObjPermissionsFailure 权限认证失败
func RetObjPermissionsFailure() RetObj {
	return RetObj{
		Code: 200,
		Data: nil,
		Msg:  ErrPermissions.String(),
	}
}

// RetObjDeviceBeenRaffledFailureZH 设备重复抽奖(中文)
func RetObjDeviceBeenRaffledFailureZH() RetObj {
	return RetObj{
		Code: 201,
		Data: nil,
		Msg:  "每台手机仅限一个账号参与",
	}
}

// RetObjDeviceBeenRaffledFailureEN 设备重复抽奖(英文)
func RetObjDeviceBeenRaffledFailureEN() RetObj {
	return RetObj{
		Code: 201,
		Data: nil,
		Msg:  "Limit one account per mobile device",
	}
}

// RetObjServerFailure 服务器失败
func RetObjServerFailure(err Errs) RetObj {
	return RetObj{
		Code: 300,
		Data: nil,
		Msg:  err.String(),
	}
}
