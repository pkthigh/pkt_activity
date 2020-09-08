package common

// ACTYPE 活动类型
type ACTYPE int

const (
	// DEFAULT 非法类型
	DEFAULT ACTYPE = 0
	// HANDNUM 手数活动
	HANDNUM ACTYPE = 1
)

// Int to int
func (ty ACTYPE) Int() int {
	return int(ty)
}
