package config

import "fmt"

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// Address 服务器运行地址 ip:port
func (server ServerConfig) Address() string {
	if server.Port == "" {
		return server.Host
	}
	return fmt.Sprintf("%s:%s", server.Host, server.Port)
}

// StorageConfig 存储配置
type StorageConfig struct {
	SQL struct {
		DBs map[string]string `json:"dbs"`
	} `json:"sql"`
	Mgo struct {
		URI      string `json:"uri"`
		DataBase string `json:"database"`
	} `json:"mgo"`
	Rds struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
	} `json:"rds"`
}

// NatsConfig 消息中间件集群配置
type NatsConfig struct {
	Client  string `json:"client"`  // 客户ID
	Cluster string `json:"cluster"` // 集群ID
	URLs    string `json:"urls"`    // 集群地址
}

// RPCXSrvConfig RPCX服务配置
type RPCXSrvConfig struct {
	Name     string `json:"name"`     // 服务名
	Protocol string `json:"protocol"` // 通讯协议
	Address  string `json:"address"`  // 通讯地址
}

// NotifySrvConfig 通知服配置
type NotifySrvConfig struct {
	Addr string `json:"addr"`
	Path string `json:"path"`
}

// MoneySrvConfig 资产服配置
type MoneySrvConfig struct {
	Addr string `json:"addr"`
	Path string `json:"path"`
}

// TokenRdsAreaConfig  用户TokenRedis配置
type TokenRdsAreaConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}
