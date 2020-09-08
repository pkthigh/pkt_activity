package config

import (
	"os"
	"pkt_activity/library/logger"

	"github.com/spf13/viper"
)

// config 配置
var config *Config

// Config 配置文件结构
type Config struct {
	Server          ServerConfig       `json:"server"`
	Storage         StorageConfig      `json:"storage"`
	Nats            NatsConfig         `json:"nats"`
	RPCX            RPCXSrvConfig      `json:"rpcx"`
	NotifySrvConfig NotifySrvConfig    `json:"NotifySrvConfig"`
	MoneySrvConfig  MoneySrvConfig     `json:"MoneySrvConfig"`
	TokenRdsArea    TokenRdsAreaConfig `json:"TokenRdsArea"`
}

func init() {
	config = new(Config)

	vpr := viper.New()

	if os.Getenv("GIN_MODE") == "release" {
		vpr.SetConfigName(ConfigProFileName.String())
	} else {
		vpr.SetConfigName(ConfigDevFileName.String())
	}
	vpr.SetConfigType(ConfigFileType.String())
	vpr.AddConfigPath(ConfigFilePath.String())
	if err := vpr.ReadInConfig(); err != nil {
		logger.FatalF("ReadInConfig", err)
	}

	if err := vpr.Unmarshal(&config); err != nil {
		logger.FatalF("Unmarshal Config", err)
	}
}

// GetServerConf 获取服务器配置
func GetServerConf() ServerConfig {
	return config.Server
}

// GetStorageConf 获取存储配置
func GetStorageConf() StorageConfig {
	return config.Storage
}

// GetNatsConf 获取Nats配置
func GetNatsConf() NatsConfig {
	return config.Nats
}

// GetRPCXConf 获取RPCX配置
func GetRPCXConf() RPCXSrvConfig {
	return config.RPCX
}

// GetNotifySrvConf 获取通知服RPCX配置
func GetNotifySrvConf() NotifySrvConfig {
	return config.NotifySrvConfig
}

// GetMoneySrvConf 获取资产服RPCX配置
func GetMoneySrvConf() MoneySrvConfig {
	return config.MoneySrvConfig
}

// GetTokenRdsAreaConfig 获取用户TokenRedis配置
func GetTokenRdsAreaConfig() TokenRdsAreaConfig {
	return config.TokenRdsArea
}
