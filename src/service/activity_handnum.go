package service

import (
	"fmt"
	"pkt_activity/common"
	"pkt_activity/library/logger"
	"pkt_activity/model"
	"sort"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/redis.v5"
)

/*
	手数活动
*/

// HandNumActivity 手数活动
type HandNumActivity struct {
	actype  common.ACTYPE
	ongoing *ActivityInfos
}

// NewHandNumActivity 手数活动
func NewHandNumActivity() Activity {
	return &HandNumActivity{
		actype: common.HANDNUM,
	}
}

// ID 活动ID
func (activity *HandNumActivity) ID() int64 {
	if activity.Ongoing() {
		return activity.ongoing.Info.ID
	}
	return 0
}

// Type 活动类型
func (activity *HandNumActivity) Type() int {
	return activity.actype.Int()
}

// Info 活动信息
func (activity *HandNumActivity) Info() ActivityInfos {
	return *activity.ongoing
}

// TimingTask 凌晨定时任务
func (activity *HandNumActivity) TimingTask() error {

	// 更新每日剩余礼包
	if activity.Ongoing() {
		for i, detail := range activity.ongoing.Details {
			activity.ongoing.Details[i].BonusRemainng = activity.ongoing.Details[i].BonusNum
			if err := Store.DBs(common.ActivityDsn).Table("ac_activity_detail").Where("`id` = ?", detail.ID).
				Update("bonus_remainng", detail.BonusNum).Error; err != nil {
				logger.ErrorF("upate mysql bonus_remainng aid: %v details id: %v error: %v", activity.ID(), activity.ongoing.Details[i].ID, err)
			}
		}
	}

	return nil
}

// Ongoing 活动进行
func (activity *HandNumActivity) Ongoing() bool {
	if activity.ongoing == nil {
		return false
	}
	if activity.ongoing.Info.Status != 2 {
		return false
	}
	return true
}

// Online 上线活动
func (activity *HandNumActivity) Online() error {
	db := Store.DBs(common.ActivityDsn)
	var info model.AcActivity

	today, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
	if err := db.Where("`status` = ? AND `online_time` = ?", common.SCHEDULING, today.Unix).Find(&info).Error; err != nil {
		logger.ErrorF("HandNumActivity Update ONGOING as error: %v", err)
		return err
	}

	var details []model.AcActivityDetail
	if err := db.Where("`activity_id` = ?", info.ID).Find(&details).Error; err != nil {
		logger.ErrorF("HandNumActivity Update ONGOING ds error: %v", err)
		return err
	}

	sort.Slice(details, func(i, j int) bool {
		if details[i].Index < details[j].Index {
			return true
		}
		return false
	})

	info.Status = common.ONGOING.Int()
	info.UpdateAt = time.Now().Unix()
	activity.ongoing = &ActivityInfos{Info: info, Details: details}

	if err := db.Table("ac_activity").Where("`id` = ?", info.ID).Update(info).Error; err != nil {
		logger.ErrorF("Update ac_activity id: %v status to onging fail: %v", info.ID, err)
	}

	return nil
}

// LoadOngoing 载入正在进行中的活动
func (activity *HandNumActivity) LoadOngoing() error {
	db := Store.DBs(common.ActivityDsn)

	var info model.AcActivity
	if err := db.Where("`status` = ?", common.ONGOING).Find(&info).Error; err != nil {
		logger.ErrorF("HandNumActivity Update ONGOING as error: %v", err)
		return err
	}

	var details []model.AcActivityDetail
	if err := db.Where("`activity_id` = ?", info.ID).Find(&details).Error; err != nil {
		logger.ErrorF("HandNumActivity Update ONGOING ds error: %v", err)
		return err
	}

	sort.Slice(details, func(i, j int) bool {
		if details[i].Index < details[j].Index {
			return true
		}
		return false
	})

	activity.ongoing = &ActivityInfos{Info: info, Details: details}
	return nil
}

// OfflineOngoing 下线正在进行的活动
func (activity *HandNumActivity) OfflineOngoing() error {
	// 有正在进行中的活动
	if ongoing := activity.Ongoing(); ongoing {
		// 提前下线
		if time.Now().Unix() < int64(activity.ongoing.Info.OfflineTime) {
			activity.ongoing.Info.Status = 3

		} else { // 正常下线
			activity.ongoing.Info.Status = 4
		}
		db := Store.DBs(common.ActivityDsn)
		if err := db.Table("ac_activity").Where("`id` = ?", activity.ongoing.Info.ID).Update(model.AcActivity{Status: activity.ongoing.Info.Status, UpdateAt: time.Now().Unix()}).Error; err != nil {
			logger.ErrorF("HandNumActivity Offline error: %v", err)
			return err
		}
		activity.ongoing = nil
	}
	return nil
}

// Verification 校验玩家
func (activity *HandNumActivity) Verification(userid string) (interface{}, common.Errs) {
	// userid 转换错误
	uid, err := strconv.Atoi(userid)
	if err != nil {
		logger.ErrorF("HandNumActivity Verification strconv.Atoi(%v) error: %v", userid, err)
		return nil, common.ErrTypeConversion
	}
	todate := time.Now().Format("2006-01-02")

	// 查询今日玩家完成德州手数
	var phn uint32
	num, err := Store.Rds(common.TexasHandOverRecordStore).HGet(userid, todate).Uint64()
	if err != nil && err != redis.Nil {
		logger.ErrorF("HandNumActivity Verification get date: %v palyer: %v hands to int64 error: %v", todate, userid, err)
		return nil, common.ErrRedisQuery
	}
	phn = uint32(num)

	// 封装PB消息返回发给客户端
	var aid uint32 = uint32(activity.ID())
	var atype uint32 = uint32(activity.Type())
	var dhn uint32 = uint32(activity.ongoing.Info.HandNum)

	// 超出门槛, 前端显示时, 只显示门槛的MAX就行
	var cd bool
	if phn >= dhn {
		cd = true
		phn = dhn
	}

	info := common.ActivityNotifyArgs{
		Uid: uint32(uid),
		Info: &common.PBActivityInfo{
			ActivityId:   &aid,
			ActivityType: &atype,
			HandsProgress: &common.PBActivityInfo_ActivityHandsProgress{
				PlayedHandsNum: &phn,
				DrawHandsNum:   &dhn,
				CanDraw:        &cd,
			},
		},
	}
	var f bool = false
	if cd {

		// 查询玩家今日是否已经完成
		result := Store.Rds(common.UserActivityResultCache).HGet(fmt.Sprintf("%s(%s)", strconv.FormatInt(activity.ID(), 10), time.Now().Format("2006-01-02")), userid)
		// 查询错误 并且查询不是为KeyNil错误
		if err := result.Err(); err != nil && err != redis.Nil {
			logger.ErrorF("HandNumActivity Verification query palyer: %v hands over error: %v", userid, err)
			return nil, common.ErrRedisQuery

			// 玩家今日已经完成了活动
		} else if err == nil {
			info.Info.HandsProgress.CanDraw = &f
			return info, common.ErrActivityCompleted
		}

		return info, common.Successful
	}
	info.Info.HandsProgress.CanDraw = &f
	return info, common.ErrVerifyNoPass
}

// Do 进行活动
func (activity *HandNumActivity) Do(userid string) (int64, common.Errs) {
	// 1.判断当前活动类是否有正在进行中的活动
	if ongoing := activity.Ongoing(); !ongoing {
		return 0, common.ErrNoOngoingActivity
	}

	// 2.判断该玩家是否满足抽奖资格
	if _, err := activity.Verification(userid); err != common.Successful {
		return 0, common.ErrVerifyNoPass
	}

	// 查询用户登陆设备ID
	uid, _ := strconv.Atoi(userid)
	var log model.UserLogin
	if err := Store.DBs(common.LogDsn).Table("user_login"+strconv.Itoa(uid%10)).Where("`uid` = ?", uid).Find(&log).Error; err != nil {
		logger.ErrorF("Find user login did error: %v", err)
		return 0, common.ErrMysqlQuery
	}

	var did string
	if log.MacAddr != "" {
		did = log.MacAddr
	}

	// 3.转换玩家ID类型
	playerid, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		return 0, common.ErrTypeConversion
	}

	var (
		itemfrees int64 // 道具费用
		insufrees int64 // 保险费用
	)

	itemstore := Store.Rds(common.ItemRecordStore)
	insustore := Store.Rds(common.InsuranceRecordStore)

	// 4.获取该玩家7天相关数据
	start := time.Now().AddDate(0, 0, -6).Format("2006-01-02")
	for _, date := range dates(start) {
		start, _ := time.ParseInLocation("2006-01-02", date, time.Local)

		var (
			itemfree, insufree int64
			err                error
		)

		type Free struct {
			Itemfree int64
			Insufree int64
		}

		var free Free

		// 道具费用
		itemfree, err = itemstore.HGet(userid, date).Int64()
		if err != nil || err == redis.Nil {
			if err := Store.DBs(common.DataDsn).
				Raw("SELECT COUNT(`num`) as 'itemfree' FROM pk_data.pkc_check_props_flow WHERE(`user_id` = ? AND `created_at` >= ? AND `created_at` < ?)",
					userid, start.Unix(), start.Unix()+60*60*24).Scan(&free).Error; err != nil {
				itemfree = free.Itemfree
				itemstore.HSet(userid, date, itemfree)
			}
		}
		itemfrees += itemfree

		// 保险费用
		insufree, _ = insustore.HGet(userid, date).Int64()
		if err != nil || err == redis.Nil {
			if err := Store.DBs(common.DataDsn).
				Raw("SELECT COUNT(`buy`) as 'insufree' FROM pk_data.insurance_flow_log_v2 WHERE(`user_id` = ? AND `created_at` >= ? AND `created_at` < ?)",
					userid, start.Unix(), start.Unix()+60*60*24).Scan(&free).Error; err != nil {
				insufree = free.Insufree
				insustore.HSet(userid, date, insufree)
			}
		}
		insufrees += insufree
	}

	// 5.筛选出满足该玩家的活动条件
	var meets []int
	for _, detail := range activity.ongoing.Details {
		// 满足 道具费 保险费 剩余礼包
		if itemfrees >= detail.Prop && insufrees >= detail.Insurance && detail.BonusRemainng > 0 {
			meets = append(meets, detail.Index)
		}
	}

	// 6.对满足的活动条件进行排序
	sort.Slice(meets, func(i, j int) bool {
		if i > j {
			return true
		}
		return false
	})

	// 7.是否有满足的活动条件
	if len(meets) < 1 {
		return 0, common.ErrNoMatchActivityDetail
	}

	// 8.累计抽奖分母
	var odds []float32
	for _, idx := range meets {
		for _, detail := range activity.ongoing.Details {
			if detail.Index == idx {
				odds = append(odds, detail.Odds)
			}
		}
	}

	// 9.进行抽奖 返回抽中的条件下标
	idx, err := Lottery(odds...)
	if err != nil {
		logger.ErrorF("ErrLotteryFailed: %v", err)
		return 0, common.ErrLotteryFailed
	}
	// 10.获取活动条件的index
	index := meets[idx]

	nowtime := time.Now().Unix()
	for i, detail := range activity.ongoing.Details {

		if detail.Index == index {

			var userinfo model.UserInfo
			if err := Store.DBs(common.UserDsn).Table("user_info"+strconv.Itoa(int(playerid)%10)).Where("`uid`=?", playerid).First(&userinfo).Error; err != nil {
				logger.ErrorF("find userinfo name by userid(%v) error: %v", playerid, err)
			}

			// 保存用户结果
			record := model.AcActivityRecord{
				CreateAt:         nowtime,
				UpdateAt:         nowtime,
				PlayerID:         playerid,
				RaffleTime:       nowtime,
				ActivityID:       activity.ID(),
				PlayerName:       userinfo.Name,
				ActivityDetailID: detail.ID,
				Bonus:            detail.Bonus,
				DeviceID:         did,
				Status:           false,
			}
			if err := Store.DBs(common.ActivityDsn).Save(&record).Error; err != nil {
				return 0, common.ErrSaveLotteryResult
			}

			// 减去今日剩余礼包
			activity.ongoing.Details[i].BonusRemainng--
			if err := Store.DBs(common.ActivityDsn).Table("ac_activity_detail").Where("`id` = ?", activity.ongoing.Details[i].ID).
				Update("bonus_remainng", gorm.Expr("bonus_remainng - ?", 1)).Error; err != nil {
				logger.ErrorF("upate mysql bonus_remainng id: %v error: %v", activity.ongoing.Details[i].ID, err)
			}

			// 缓存 活动ID(当日) 用户 活动条件ID
			if err := Store.Rds(common.UserActivityResultCache).HSet(fmt.Sprintf("%s(%s)", strconv.FormatInt(activity.ID(), 10), time.Now().Format("2006-01-02")), userid, record.ActivityDetailID).Err(); err != nil {
				logger.ErrorF("upate redis cache user activity result error: %v", err)
			}

			return record.ID, common.Successful
		}
	}

	return 0, common.ErrNoMatchActivityDetail
}
