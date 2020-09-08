package service

import (
	"fmt"
	"net/http"
	"pkt_activity/library/logger"
	"pkt_activity/library/storage"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Store 存储模块
var Store *storage.Storage

func init() {
	// 初始化存储模块
	var err error
	Store, err = storage.NewStorage()
	if err != nil {
		logger.FatalF("store init fail: %v", err)
	}
}

// dates 日期计算
func dates(sdate string) []string {
	var dates []string

	start, _ := time.ParseInLocation("2006-01-02", sdate, time.Local)
	nowtime := time.Now().Local()
	if start.Unix() > nowtime.Unix() {
		return dates
	}

	dates = append(dates, start.Format("2006-01-02"))
	for start.Format("2006-01-02") != nowtime.Format("2006-01-02") {
		start = start.AddDate(0, 0, 1)
		dates = append(dates, start.Format("2006-01-02"))
	}

	return dates
}

// Cors cors
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		var keys []string
		for key := range c.Request.Header {
			keys = append(keys, key)
		}

		str := strings.Join(keys, ", ")
		if str != "" {
			str = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", str)
		} else {
			str = "access-control-allow-origin, access-control-allow-headers"
		}

		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma, language")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}

		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
