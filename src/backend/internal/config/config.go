package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	MinIO    MinIOConfig    `yaml:"minio"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"` // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	SSLMode         string        `yaml:"sslmode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`  // 支持 "300s", "5m" 等格式
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"` // 支持 "60s", "1m" 等格式
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `yaml:"secret"`
	Expiration time.Duration `yaml:"expiration"`
}

// MinIOConfig 对象存储配置
type MinIOConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	UseSSL    bool   `yaml:"use_ssl"`
	BaseURL   string `yaml:"base_url"`
}

var AppConfig *Config

// Load 加载配置文件
func Load(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	AppConfig = &Config{}
	if err := yaml.Unmarshal(data, AppConfig); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// 从环境变量覆盖配置（如果存在）
	loadFromEnv()
	ensureMinioDefaults()

	return nil
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv() {
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		AppConfig.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		fmt.Sscanf(dbPort, "%d", &AppConfig.Database.Port)
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		AppConfig.Database.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		AppConfig.Database.Password = dbPassword
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		AppConfig.Database.DBName = dbName
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		AppConfig.JWT.Secret = jwtSecret
	}
	if minioEndpoint := os.Getenv("MINIO_ENDPOINT"); minioEndpoint != "" {
		AppConfig.MinIO.Endpoint = minioEndpoint
	}
	if minioAccessKey := os.Getenv("MINIO_ROOT_USER"); minioAccessKey != "" {
		AppConfig.MinIO.AccessKey = minioAccessKey
	} else if minioAccessKey := os.Getenv("MINIO_ACCESS_KEY"); minioAccessKey != "" {
		AppConfig.MinIO.AccessKey = minioAccessKey
	}
	if minioSecret := os.Getenv("MINIO_ROOT_PASSWORD"); minioSecret != "" {
		AppConfig.MinIO.SecretKey = minioSecret
	} else if minioSecret := os.Getenv("MINIO_SECRET_KEY"); minioSecret != "" {
		AppConfig.MinIO.SecretKey = minioSecret
	}
	if minioBucket := os.Getenv("MINIO_BUCKET"); minioBucket != "" {
		AppConfig.MinIO.Bucket = minioBucket
	}
	if minioRegion := os.Getenv("MINIO_REGION"); minioRegion != "" {
		AppConfig.MinIO.Region = minioRegion
	}
	if minioUseSSL := os.Getenv("MINIO_USE_SSL"); minioUseSSL != "" {
		if parsed, err := strconv.ParseBool(minioUseSSL); err == nil {
			AppConfig.MinIO.UseSSL = parsed
		}
	}
	if minioBaseURL := os.Getenv("MINIO_SERVER_URL"); minioBaseURL != "" {
		AppConfig.MinIO.BaseURL = minioBaseURL
	} else if minioBaseURL := os.Getenv("MINIO_BASE_URL"); minioBaseURL != "" {
		AppConfig.MinIO.BaseURL = minioBaseURL
	}
}

// ensureMinioDefaults 确保 MinIO 配置存在合理默认值
func ensureMinioDefaults() {
	if AppConfig.MinIO.Endpoint == "" {
		AppConfig.MinIO.Endpoint = "localhost:9000"
	}
	if AppConfig.MinIO.Bucket == "" {
		AppConfig.MinIO.Bucket = "onetaste-media"
	}
	if AppConfig.MinIO.BaseURL == "" && AppConfig.MinIO.Endpoint != "" {
		scheme := "http"
		if AppConfig.MinIO.UseSSL {
			scheme = "https"
		}
		AppConfig.MinIO.BaseURL = fmt.Sprintf("%s://%s", scheme, AppConfig.MinIO.Endpoint)
	}
}

// GetDatabaseDSN 获取数据库连接字符串
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
