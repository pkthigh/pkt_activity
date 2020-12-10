package service

import (
	"pkt_activity/common"
	"pkt_activity/library/logger"
	"pkt_activity/model"

	"github.com/gin-gonic/gin"
)

// NewActivityGateway 新的活动网关
func (srv *ActivityService) NewActivityGateway() *gin.Engine {
	router := gin.Default()
	router.Use(Cors())

	v1 := router.Group("/v1")
	{
		activity := v1.Group("/activity")
		{
			// 当前活动列表
			activity.GET("/list", srv.ActivityList)
			// 用户活动信息
			activity.GET("/info", srv.UserActivityInfo)
			// 用户中奖记录
			activity.GET("/records", srv.UserActivityRecords)
			// 活动结果轮转
			activity.GET("/results", srv.UserActivityResults)
			// 用户参与活动
			activity.POST("/participate", srv.UserActivityParticipate)
		}

	}
	return router
}

// GetTokenToUserID tk -> uid
func (srv *ActivityService) GetTokenToUserID(c *gin.Context) string {
	tk := c.Request.Header.Get("token")
	if tk == "" {
		return ""
	}
	result := Store.Rds(common.UserTokenStore).HGet("O_"+tk, "1")
	if result.Err() != nil {
		logger.ErrorF("Find UserTokenStore token(%v) error: %v", tk, result.Err())
		return ""

	}

	logger.InfoF("tk: %v > uid: %v", tk, result.Val())
	return result.Val()
}

// GetUserLanguage 活动用户语言
func (srv *ActivityService) GetUserLanguage(c *gin.Context) int {
	var code int
	language := c.Request.Header.Get("language")
	if language == "zh" {
		code = 0
	}
	if language == "en" {
		code = 1
	}
	return code
}

// ActivityList 活动列表
func (srv *ActivityService) ActivityList(c *gin.Context) {
	var list []ActivityInfos
	for _, activity := range srv.activitys {
		if activity.Ongoing() {
			list = append(list, activity.Info())
		}
	}
	c.JSON(200, common.RetObjSuccessful(list))
	return
}

// UserActivityInfo 用户活动信息
func (srv *ActivityService) UserActivityInfo(c *gin.Context) {
	type ParamModel struct {
		AID int64 `form:"aid" json:"aid"`
	}

	var param ParamModel
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(200, common.RetObjParamFailure(common.ErrParamInvalid))
		return
	}

	userid := srv.GetTokenToUserID(c)
	if userid == "" {
		c.JSON(200, common.RetObjPermissionsFailure())
		return
	}

	if srv.activitys[common.HANDNUM].Ongoing() && srv.activitys[common.HANDNUM].ID() == param.AID {

		type Resp struct {
			ID int64 `json:"id"`

			OnlineTime  int `json:"online_time"`  // 活动开始时间
			OfflineTime int `json:"offline_time"` // 活动结束时间

			NameZh    string `json:"name_zh"`    // 活动名字 中文
			ContentZh string `json:"content_zh"` // 活动内容 中文
			PicZhURL  string `json:"pic_zh_url"` // 活动海报 中文

			NameEn    string `json:"name_en"`    // 活动名字 英文
			ContentEn string `json:"content_en"` // 活动内容 英文
			PicEnURL  string `json:"pic_en_url"` // 活动海报 英文

			Status int `json:"status"` // [0:未开启 / 1:排期中 / 2:进行中 / 3:提前下线 4: 已结束]

			HandNum      int `json:"hand_num"`      // 活动门槛手数
			PlayerHands  int `json:"player_hands"`  // 用户当前手数
			PlayerStatus int `json:"palyer_status"` // [0: 不满足抽奖条件 1: 满足抽奖条件未抽奖 2: 已抽奖]
		}

		var resp Resp
		if err := Store.DBs(common.ActivityDsn).Table("ac_activity").Where("`id` = ? AND `status` = ?", param.AID, 2).Find(&resp).Error; err != nil {
			logger.Error("Find ActivityDsn error: %v", err)
			c.JSON(200, common.ErrMysqlQuery)
			return
		}

		var (
			result interface{}
			err    common.Errs
		)
		result, err = srv.activitys[common.HANDNUM].Verification(userid)
		if err != common.Successful && err != common.ErrVerifyNoPass && err != common.ErrActivityCompleted {
			c.JSON(200, common.RetObjServerFailure(err))
			return
		}

		resp.PlayerHands = int(*result.(common.ActivityNotifyArgs).Info.HandsProgress.PlayedHandsNum)

		if err == common.ErrVerifyNoPass {
			resp.PlayerStatus = 0
		}
		if err == common.Successful {
			resp.PlayerStatus = 1
		}
		if err == common.ErrActivityCompleted {
			resp.PlayerStatus = 2
		}

		c.JSON(200, common.RetObjSuccessful(resp))
		return
	}

	c.JSON(200, common.RetObjParamFailure(common.ErrNoOngoingActivity))
	return
}

// UserActivityRecords 用户活动记录
func (srv *ActivityService) UserActivityRecords(c *gin.Context) {
	type ParamModel struct {
		AID int64 `form:"aid" json:"aid"`
	}

	var param ParamModel
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(200, common.RetObjParamFailure(common.ErrParamInvalid))
		return
	}

	userid := srv.GetTokenToUserID(c)
	if userid == "" {
		c.JSON(200, common.RetObjPermissionsFailure())
		return
	}

	var records []model.AcActivityRecord
	if err := Store.DBs(common.ActivityDsn).Table("ac_raffle_record").Where("`player_id` = ? AND `activity_id` = ?", userid, param.AID).Order("`create_at` desc").Find(&records).Error; err != nil {
		logger.ErrorF("Find player %v AcActivityRecord error: %v", userid, err)
		c.JSON(200, common.RetObjServerFailure(common.ErrMysqlQuery))
		return
	}

	c.JSON(200, common.RetObjSuccessful(records))
	return
}

// UserActivityParticipate 用户活动参与
func (srv *ActivityService) UserActivityParticipate(c *gin.Context) {
	type ParamModel struct {
		AID   int64 `form:"aid" json:"aid"`
		AType int   `form:"atype" json:"atype"`
		Lang  bool  `form:"lang" json:"lang"` // 语言 0:ZH 1:EN
	}
	srv.GetUserLanguage(c)

	var param ParamModel
	if err := c.BindJSON(&param); err != nil {
		c.JSON(200, common.RetObjParamFailure(common.ErrParamInvalid))
		return
	}

	userid := srv.GetTokenToUserID(c)
	if userid == "" {
		c.JSON(200, common.RetObjPermissionsFailure())
		return
	}

	id, err := srv.participate(param.AID, param.AType, userid, param.Lang)
	if err != common.Successful {

		if err == common.ErrDeviceBeenRaffled {
			if srv.GetUserLanguage(c) == 1 {
				c.JSON(200, common.RetObjDeviceBeenRaffledFailureEN())
			} else {
				c.JSON(200, common.RetObjDeviceBeenRaffledFailureZH())
			}
			return
		}

		// 用户资产更新失败 活动记录更新失败
		if err != common.ErrUpdateUserAssets && err != common.ErrUpdateActivityRecoed {
			c.JSON(200, common.RetObjServerFailure(err))
			return
		}
	}

	var record model.AcActivityRecord
	if err := Store.DBs(common.ActivityDsn).Where("`id` = ?", id).Find(&record).Error; err != nil {
		logger.ErrorF("ErrMysqlQuery AcActivityRecord id: %v err: %v", id, err)
		c.JSON(200, common.RetObjServerFailure(common.ErrMysqlQuery))
		return
	}

	type Resp struct {
		DetailID int64 `json:"detail_id"`
		Bonus    int64 `json:"bonus"` // 奖金
	}

	c.JSON(200, common.RetObjSuccessful(Resp{DetailID: record.ActivityDetailID, Bonus: record.Bonus}))
	return
}

// UserActivityResults 结果
func (srv *ActivityService) UserActivityResults(c *gin.Context) {
	type ParamModel struct {
		AID   int64 `form:"aid" json:"aid"`
		Count int   `form:"count" json:"count"`
	}

	var param ParamModel
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(200, common.RetObjParamFailure(common.ErrParamInvalid))
		return
	}

	if param.Count == 0 {
		param.Count = 50
	}

	type Item struct {
		PlayerName string `json:"player_name"`
		Bonus      int    `json:"bonus"`
	}

	var records []Item
	if err := Store.DBs(common.ActivityDsn).Table("ac_raffle_record").Where("`activity_id` = ? AND `bonus` > 0", param.AID).Limit(param.Count).Find(&records).Error; err != nil {
		logger.ErrorF("UserActivityResults ErrMysqlQuery AcActivityRecord id: %v err: %v", param.AID, err)
		c.JSON(200, common.RetObjServerFailure(common.ErrMysqlQuery))
		return
	}

	c.JSON(200, common.RetObjSuccessful(records))
	return
}
