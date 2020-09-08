package service

import (
	"context"
	"pkt_activity/common"
	"pkt_activity/library/logger"
	"strconv"
)

// HandActivityProgressReq 手数活动请求
type HandActivityProgressReq struct {
	Aid int
	Uid uint32
}

// HandActivityProgressResp 手数活动响应
type HandActivityProgressResp struct {
	PlayedHandsNum int  // 当前进行的手数
	DrawHandsNum   int  // 抽奖需要的手数
	CanDraw        bool // 是否可以抽奖
}

// GetHandActivityProgress 获取手数活动进度
func (srv *ActivityService) GetHandActivityProgress(ctx context.Context, req *HandActivityProgressReq, resp *HandActivityProgressResp) error {
	// 当前正在进行的活动等于Request.AID
	if srv.activitys[common.HANDNUM].Ongoing() && srv.activitys[common.HANDNUM].ID() == int64(req.Aid) {

		userid := strconv.Itoa(int(req.Uid))

		val, err := srv.activitys[common.HANDNUM].Verification(userid)
		if err != common.Successful && err != common.ErrVerifyNoPass && err != common.ErrActivityCompleted {
			return err.Error()
		}

		result := val.(common.ActivityNotifyArgs)
		resp.PlayedHandsNum = int(*result.Info.HandsProgress.PlayedHandsNum)
		resp.DrawHandsNum = int(*result.Info.HandsProgress.DrawHandsNum)
		if err == common.ErrActivityCompleted || err == common.ErrVerifyNoPass {
			resp.CanDraw = false
		} else {
			resp.CanDraw = *result.Info.HandsProgress.CanDraw
		}
		logger.InfoF("RPCX GetHandActivityProgress AID: %v, UID: %v > Resp: %v ", req.Aid, req.Uid, *resp)
	} else {
		logger.ErrorF("RPCX GetHandActivityProgress AID: %v, UID: %v 没有正在进行中的活动", req.Aid, req.Uid)
		return common.ErrNoOngoingActivity.Error()
	}

	return nil
}

// OfflineReq 下线请求
type OfflineReq struct {
	AID int64
}

// OfflineResp 下线响应
type OfflineResp struct {
}

// Offline 下线活动
func (srv *ActivityService) Offline(ctx context.Context, req *OfflineReq, resp *OfflineResp) error {
	for _, activity := range srv.activitys {
		if activity.Ongoing() {
			if activity.ID() == req.AID {
				logger.InfoF("Offline RPCX id: %v", req.AID)
				if err := srv.activitys[common.ACTYPE(activity.Type())].OfflineOngoing(); err != nil {
					logger.ErrorF("Offline RPCX id: %v, err: %v", req.AID, err)
					return err
				}
				if err := srv.change(); err != nil {
					logger.ErrorF("online change error: %v", err)
				}
				return nil
			}
		}
	}
	return nil
}
