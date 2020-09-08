package common

// ACSTATUS 活动状态
type ACSTATUS int

const (
	// INIT 未开启
	INIT ACSTATUS = 0
	// SCHEDULING 排期中
	SCHEDULING ACSTATUS = 1
	// ONGOING 进行中
	ONGOING ACSTATUS = 2
	// SUSPEND 已中止
	SUSPEND ACSTATUS = 3
	// END 结束
	END ACSTATUS = 4
)

/*
	1. 运营人员> 编辑新的活动入库 = 未开启
	2. 运营人员> 将新的活动纳入排期 = 排期中
	3. 系统检测> 到排期中的活动到活动上线日 = 进行中
	4. 运营人员> 提前下线正在进行中的活动 = 已中止
	5. 系统检测> 正在进行中的活动到活动下线日= 结束
*/

// Int to int
func (status ACSTATUS) Int() int {
	return int(status)
}
