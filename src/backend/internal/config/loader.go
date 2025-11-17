package config

import (
	"os"
	"path/filepath"
	"time"
)

// LoadConfig 加载配置文件
// 优先从环境变量获取配置路径，否则使用默认路径
func LoadConfig() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// 默认配置文件路径
		configPath = "config/config.yaml"
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 如果配置文件不存在，使用默认配置
		return loadDefaultConfig()
	}

	return Load(configPath)
}

// loadDefaultConfig 加载默认配置
func loadDefaultConfig() error {
	AppConfig = &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			User:            "postgres",
			Password:        "postgres",
			DBName:          "onetaste_family",
			SSLMode:         "disable",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 300 * time.Second, // 5分钟
			ConnMaxIdleTime: 60 * time.Second,  // 1分钟
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		JWT: JWTConfig{
			Secret:     "your-secret-key-change-in-production",
			Expiration: 24 * time.Hour, // 24小时
		},
	}

	// 从环境变量覆盖配置
	loadFromEnv()

	return nil
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// 尝试多个可能的路径
		possiblePaths := []string{
			"config/config.yaml",
			"./config/config.yaml",
			filepath.Join(os.Getenv("HOME"), ".onetaste_family", "config.yaml"),
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}

		return "config/config.yaml"
	}

	return configPath
}

