package config

// CONFIG 配置常量
type CONFIG string

const (
	// ConfigDevFileName 配置文件名
	ConfigDevFileName CONFIG = "conf_dev"
	// ConfigProFileName 配置文件名
	ConfigProFileName CONFIG = "conf_pro"
	// ConfigFileType 配置文件类型
	ConfigFileType CONFIG = "json"
	// ConfigFilePath 配置文件路径
	ConfigFilePath CONFIG = "./confs"
)

// String to string
func (config CONFIG) String() string {
	return string(config)
}
