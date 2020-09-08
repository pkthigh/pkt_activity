package service

import (
	"fmt"
	"math/rand"
	"pkt_activity/common"
	"pkt_activity/model"
	"strconv"
	"time"
)

// Activity 活动抽象类
type Activity interface {
	ID() int64
	Type() int
	Info() ActivityInfos

	Online() error                                         // 判断 排期中的活动是否今日上线.
	Ongoing() bool                                         // 是否 正在进行中的活动.
	LoadOngoing() error                                    // 载入 正在进行中的活动 从数据库.
	OfflineOngoing() error                                 // 下线 正在进行中的活动 至数据库 (正常下线或提前离线都由该函数处理).
	TimingTask() error                                     // 凌晨定时任务(由Service凌晨统一调用)
	Verification(userid string) (interface{}, common.Errs) // 验证 玩家 是否 符合 正在进行的活动.
	Do(userid string) (int64, common.Errs)                 // 进行 抽奖 返回 抽奖 记录ID & 错误.
}

// ActivityInfos 活动信息
type ActivityInfos struct {
	Info    model.AcActivity         `json:"info"`
	Details []model.AcActivityDetail `json:"details"`
}

// Lottery lottery
func Lottery(odds ...float32) (int, error) {
	var pool []int

	if len(odds) == 1 {
		return 0, nil
	}

	for i, f := range odds {
		if f < 0 || f > 1 {
			return i, fmt.Errorf("odds param error")
		}

		value, err := strconv.ParseFloat(fmt.Sprintf("%.3f", f), 64)
		if err != nil {
			return i, err
		}

		for num := 0; num < int(value*1000); num++ {
			pool = append(pool, i)
		}
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := len(pool)
	for i := 0; i < length; i++ {
		number := random.Intn(length)
		pool[length-1], pool[number] = pool[number], pool[length-1]
	}

	return pool[random.Intn(length)-1], nil
}
