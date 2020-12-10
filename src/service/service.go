package service

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"pkt_activity/common"
	"pkt_activity/library/config"
	"pkt_activity/library/logger"
	"pkt_activity/model"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"gopkg.in/redis.v5"
)

// ActivityService 活动服务
type ActivityService struct {
	sync.Mutex
	gateway      *http.Server               // http gateway
	service      *server.Server             // rpcx service
	message      chan redis.Message         // redis sub message
	PubSubClient *redis.PubSub              // redis pub & sub
	NotifyClient client.XClient             // rpcx client - notify
	MoneysClient client.XClient             // rpcx client - assets
	activitys    map[common.ACTYPE]Activity // class implements
}

// NewActivityService 新的活动服务
func NewActivityService() ActivityService {

	logger.InfoF("NotifySrvConf: %v  MoneySrvConf: %v", config.GetNotifySrvConf(), config.GetMoneySrvConf())
	if config.GetNotifySrvConf().Addr == "" || config.GetMoneySrvConf().Addr == "" {
		logger.FatalF("Laod Configure File NotifySrvConf & GetMoneySrvConf Address Empty")
	}

	// 发现通知客户端通知服务
	srv := ActivityService{
		service:      server.NewServer(),
		message:      make(chan redis.Message, 1024),
		NotifyClient: client.NewXClient(config.GetNotifySrvConf().Path, client.Failtry, client.RandomSelect, client.NewPeer2PeerDiscovery(config.GetNotifySrvConf().Addr, ""), client.DefaultOption),
		MoneysClient: client.NewXClient(config.GetMoneySrvConf().Path, client.Failtry, client.RandomSelect, client.NewPeer2PeerDiscovery(config.GetMoneySrvConf().Addr, ""), client.DefaultOption),
		activitys:    make(map[common.ACTYPE]Activity),
	}

	srv.gateway = &http.Server{
		Addr:           config.GetServerConf().Address(),
		Handler:        srv.NewActivityGateway(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 注册对应活动类型的活动处理实现类
	srv.activitys[common.HANDNUM] = NewHandNumActivity()

	return srv
}

// Run 运行
func (srv *ActivityService) Run() error {
	/*
		启动顺序
		0.载入正在进行中的活动
		1.启动HTTP网关
		2.启动RPCX网关
		3.连接NotifyRedis
		4.测试UserTokenRedisDB
		5.接收NotifyRedis订阅消息传入有效消息通道
		6.每日调用活动定时任务
		7.读取更新消息分发给各活动实现类处理业务
	*/

	// 0.载入正在进行中的活动
	if err := srv.activitys[common.HANDNUM].LoadOngoing(); err != nil {
		logger.ErrorF("start the load hands activity fail: %v", err)
	} else {
		logger.Info("start the load hands activity successful")
	}

	logger.InfoF("启动步骤: 0/7 成功")

	// 1.启动HTTP网关
	go func() {
		if err := srv.gateway.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.FatalF("start the http gateway fail: %v", err)
		}
		logger.Info("start the http gateway successful")
	}()

	logger.InfoF("启动步骤: 1/7 成功")

	// 2.启动RPCX网关
	go func() {
		if err := srv.service.Register(srv, ""); err != nil {
			logger.FatalF("start the rcpx service RegisterName fail: %v", err)
		}
		if err := srv.service.Serve(config.GetRPCXConf().Protocol, config.GetRPCXConf().Address); err != nil && err != server.ErrServerClosed {
			logger.FatalF("start the rcpx service Serve fail: %v", err)
		}
		logger.Info("start the rcpx service successful")
	}()

	logger.InfoF("启动步骤: 2/7 成功")

	// 3.连接NotifyRedis
	var err error
	srv.PubSubClient, err = Store.Rds(common.UpdateNotification).Subscribe(common.PkcHandOverSubject.String())
	if err != nil {
		logger.FatalF("start the subscribe redis notify fail: %v", err)
	}

	logger.InfoF("启动步骤: 3/7 成功")

	// 4.连接ClientTokenRedisDB
	var flag bool
	for i := 0; i < 3; i++ {
		if err := Store.Rds(common.UserTokenStore).Ping().Err(); err == nil {
			flag = true
			logger.Info("start the token redis db successful")
			break
		}
	}
	if !flag {
		logger.FatalF("start the token redis db fail: %v", err)
	}

	logger.InfoF("启动步骤: 4/7 成功")

	// 5.接收NotifyRedis订阅消息传入有效消息通道
	go func() {
		for {
			receipt, err := srv.PubSubClient.Receive()
			if err != nil {
				// logger.ErrorF("notice receive error: %v", err)
			}
			if receipt != "" {
				switch v := receipt.(type) {
				case *redis.Message:
					srv.message <- *v
					logger.InfoF("notice %s: message: %s\n", v.Channel, v.Payload)

				case error:
					logger.ErrorF("notice receipt to type error: %v", v.Error())
				}
			}
		}
	}()

	logger.InfoF("启动步骤: 5/7 成功")

	// 6.上下线检查 & 每日调用活动定时任务

	go func() {
		for {
			nowtime := time.Now()
			if nowtime.Format("15:04") == "00:00" {
				for _, activity := range srv.activitys {
					if activity.Ongoing() {
						if err := activity.TimingTask(); err != nil {
							logger.ErrorF("[Regularly Check]: 凌晨定时任务 活动: %v 错误: %v", activity.ID(), err)
						}
						logger.InfoF("指定定时任务成功ID: %v", activity.ID())
					}
				}
				time.Sleep(time.Hour)
			} else {
				time.Sleep(2 * time.Second)
			}
		}
	}()

	go func() {
		for {
			nowtime := time.Now()
			mt := nowtime.Format("2006-01-02 15:04")
			for _, activity := range srv.activitys {
				if activity.Ongoing() {
					// 检测当前活动是否下线
					go func() {
						if time.Unix(int64(activity.Info().Info.OfflineTime), 0).Format("2006-01-02 15:04") == mt {

							logger.InfoF("[Regularly Check]: 活动: %v 下线: %v", activity.ID(), mt)
							if err := srv.activitys[common.ACTYPE(activity.Type())].OfflineOngoing(); err != nil {
								logger.ErrorF("[Regularly Check]: 活动: %v 下线: %v 错误: %v", activity.ID, mt, err)
								return
							}

							time.Sleep(1 * time.Second)
							if err := srv.change(); err != nil {
								logger.ErrorF("[Regularly Check]: 活动: %v 下线通知: %v 错误: %v", activity.ID, mt, err)
							}
						}
					}()
				}
			}

			// 查询排期中的活动是否上线
			var activitys []model.AcActivity
			if err := Store.DBs(common.ActivityDsn).Where("`status` = 1 AND `ac_type` = 1").Find(&activitys).Error; err != nil {
				logger.ErrorF("[Regularly Check]: 查询数据库排期活动错误: %v", err)
			}

			for _, a := range activitys {
				if !srv.activitys[common.ACTYPE(a.AcType)].Ongoing() {
					if time.Unix(int64(a.OnlineTime), 0).Format("2006-01-02 15:04") == mt {
						if err := Store.DBs(common.ActivityDsn).Table("ac_activity").Where("`id` = ?", a.ID).Update(model.AcActivity{Status: 2, UpdateAt: time.Now().Unix()}).Error; err != nil {
							logger.Error("[Regularly Check]: 更新数据库活动: %v 上线: %v 错误: %v", a.ID, mt, err)
							return
						}
						if err := srv.activitys[common.ACTYPE(a.AcType)].LoadOngoing(); err != nil {
							logger.Error("[Regularly Check]: 更新内存活动: %v 上线: %v 错误: %v", a.ID, mt, err)
							return
						}
						time.Sleep(1 * time.Second)
						if err := srv.change(); err != nil {
							logger.ErrorF("[Regularly Check]: 活动: %v 上线通知: %v 错误: %v", a.ID, mt, err)
						}
					}
				}
			}
			time.Sleep(time.Second)
		}
	}()

	logger.InfoF("启动步骤: 6/7 成功")

	// 7.读取更新消息分发给各活动实现类处理业务
	go func() {
		for {
			msg := <-srv.message
			switch msg.Channel {

			// 德州玩家更新
			case common.PkcHandOverSubject.String():
				if !srv.activitys[common.HANDNUM].Ongoing() {
					continue
				}

				uid := msg.Payload
				aid := strconv.FormatInt(srv.activitys[common.HANDNUM].ID(), 10)
				if aid == "0" {
					logger.ErrorF("aid = 0")
					continue
				}
				todate := time.Now().Format("2006-01-02")

				// 避免重复主动推送
				if result := Store.Rds(common.UpdateNotification).HGet(fmt.Sprintf("%s(%s)", aid, todate), uid); result.Err() == nil && result.Val() == "ok" {
					continue
				}

				if info, err := srv.activitys[common.HANDNUM].Verification(uid); err == common.Successful {
					if err := srv.notify(info); err != nil {
						logger.ErrorF("用户: %v 满足活动: %v 推送错误: %v", uid, aid, err)
					} else {
						if err := Store.Rds(common.UpdateNotification).HSet(fmt.Sprintf("%s(%s)", aid, todate), uid, "ok").Err(); err != nil {
							logger.ErrorF("handnum activity set UpdateNotification error: %v", err)
						}
					}
				}

				// 钓鱼玩家更新(该版本忽略)
				// case common.FishHandOverSubject.String():
			}
		}
	}()

	logger.InfoF("启动步骤: 7/7 成功")

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	<-signals

	logger.InfoF("侦听到进程关闭信号, 进行关闭中...")

	srv.Close()

	logger.InfoF("进程已关闭, Bye~")
	return nil
}

// Close close
func (srv *ActivityService) Close() {
	if err := srv.PubSubClient.Close(); err != nil {
		logger.ErrorF("Close PubSubClient fail: %v", err)
	}

	if err := srv.gateway.Close(); err != nil {
		logger.ErrorF("Close Gateway fail: %v", err)
	}

	if err := srv.service.Close(); err != nil {
		logger.ErrorF("Close ActivityService fail: %v", err)
	}

	if err := srv.NotifyClient.Close(); err != nil {
		logger.ErrorF("Close NotifyClient fail: %v", err)
	}

	if err := srv.MoneysClient.Close(); err != nil {
		logger.ErrorF("Close MoneysClient fail: %v", err)
	}
	logger.Info(">>> Close Successful")
}
