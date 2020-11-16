package service

import (
	"context"
	"fmt"
	"pkt_activity/common"
	"pkt_activity/library/logger"
	"pkt_activity/model"
	"strconv"
	"time"
)

// change 活动上下线通知
func (srv *ActivityService) change() error {
	type Reply struct{}
	if err := srv.NotifyClient.Call(context.Background(), "BroadcastActivityStatusChange", nil, &Reply{}); err != nil {
		return err
	}
	return nil
}

// notify 活动推送通知
func (srv *ActivityService) notify(info interface{}) error {
	type Reply struct{}
	if err := srv.NotifyClient.Call(context.Background(), "AddActivityNotificationToUsers", info, &Reply{}); err != nil {
		return err
	}
	return nil
}

// assets 用户资产操作
func (srv *ActivityService) assets(uid uint32, deltaValue int64, optionType uint32, desc string) error {
	var args []common.SetMoneyArgs
	args = append(args, common.SetMoneyArgs{OptionType: optionType, Desc: desc, DeltaValue: float64(deltaValue)})
	reply := &common.SetMoneyRsp{}
	return srv.MoneysClient.Call(context.Background(), "SetUserPKC", common.SetMoneyReq{UID: uid, SetArgs: args}, reply)
}

// Participate 参加活动
func (srv *ActivityService) participate(aid int64, atype int, userid string, lang bool) (int64, common.Errs) {
	srv.Lock()
	defer srv.Unlock()

	activity, ok := srv.activitys[common.ACTYPE(atype)]
	if !ok {
		logger.ErrorF("get activitys by atype(%v) = %v", atype, ok)
		return 0, common.ErrNoOngoingActivity
	}

	// 返回活动记录ID
	id, err := activity.Do(userid)
	if err != common.Successful {
		return 0, err
	}

	var record model.AcActivityRecord
	if err := Store.DBs(common.ActivityDsn).Where("`id` =  ?", id).Find(&record).Error; err != nil {
		logger.ErrorF("ErrMysqlQuery recoed:%v error: %v", id, err)
		return 0, common.ErrMysqlQuery
	}

	playerid, ec := strconv.Atoi(userid)
	if ec != nil {
		return 0, common.ErrTypeConversion
	}

	// 修改活动记录
	nowtime := time.Now().Unix()
	record.UpdateAt = nowtime
	record.Status = true
	record.Lang = lang

	if record.Bonus > 0 {
		// 更新用户资产
		if err := srv.assets(uint32(playerid), record.Bonus, 83,
			fmt.Sprintf("活动ID: %v 活动详情ID: %v 命中奖金: %v分",
				record.ActivityID, record.ActivityDetailID, record.Bonus)); err != nil {

			logger.ErrorF("ErrUpdateUserAssets userid:%v recoed_id: %v err: %v", userid, id, err)
			return 0, common.ErrUpdateUserAssets
		}
	}

	if err := Store.DBs(common.ActivityDsn).Table("ac_raffle_record").Where("`id` =  ?", id).Update(&record).Error; err != nil {
		logger.ErrorF("ErrUpdateActivityRecoed userid:%v recoed_id: %v err: %v", userid, id, err)
		return 0, common.ErrUpdateActivityRecoed
	}

	return id, common.Successful
}
