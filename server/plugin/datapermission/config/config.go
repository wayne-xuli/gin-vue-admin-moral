package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DataPermission 数据权限配置
type DataPermission struct {
	// 是否启用数据权限
	Enabled bool `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	// 默认数据范围
	DefaultDataScope string `mapstructure:"default-data-scope" json:"defaultDataScope" yaml:"default-data-scope"`
	// 是否启用字段权限
	FieldPermissionEnabled bool `mapstructure:"field-permission-enabled" json:"fieldPermissionEnabled" yaml:"field-permission-enabled"`
	// 是否启用SQL注入检查
	SqlInjectionCheck bool `mapstructure:"sql-injection-check" json:"sqlInjectionCheck" yaml:"sql-injection-check"`
	// 缓存配置
	Cache CacheConfig `mapstructure:"cache" json:"cache" yaml:"cache"`
	// 日志配置
	Log LogConfig `mapstructure:"log" json:"log" yaml:"log"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	// 是否启用缓存
	Enabled bool `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	// 缓存过期时间（秒）
	Expiration int `mapstructure:"expiration" json:"expiration" yaml:"expiration"`
	// 缓存键前缀
	KeyPrefix string `mapstructure:"key-prefix" json:"keyPrefix" yaml:"key-prefix"`
}

// LogConfig 日志配置
type LogConfig struct {
	// 是否启用权限日志
	Enabled bool `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	// 日志级别
	Level string `mapstructure:"level" json:"level" yaml:"level"`
	// 是否记录SQL语句
	LogSQL bool `mapstructure:"log-sql" json:"logSQL" yaml:"log-sql"`
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() DataPermission {
	return DataPermission{
		Enabled:                true,
		DefaultDataScope:       "self",
		FieldPermissionEnabled: true,
		SqlInjectionCheck:      true,
		Cache: CacheConfig{
			Enabled:    true,
			Expiration: 300, // 5分钟
			KeyPrefix:  "data_permission:",
		},
		Log: LogConfig{
			Enabled: true,
			Level:   "info",
			LogSQL:  false,
		},
	}
}

// LoadConfig 从独立配置文件加载配置
func LoadConfig() DataPermission {
	// 获取当前文件所在目录
	currentDir, _ := os.Getwd()
	configPath := filepath.Join(currentDir, "plugin", "datapermission", "config", "config.yaml")

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return GetDefaultConfig()
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return GetDefaultConfig()
	}

	// 解析配置文件
	var configFile struct {
		DataPermission DataPermission `yaml:"datapermission"`
	}

	if err := yaml.Unmarshal(data, &configFile); err != nil {
		return GetDefaultConfig()
	}

	return configFile.DataPermission
}
