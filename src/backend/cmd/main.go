package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "onetaste-family/backend/docs/swagger" // 导入生成的 Swagger 文档
	"onetaste-family/backend/internal/config"
	"onetaste-family/backend/internal/handlers"
	"onetaste-family/backend/pkg/cache"
	"onetaste-family/backend/pkg/database"
	"onetaste-family/backend/pkg/storage"
)

// @title           OneTaste Family API
// @version         1.0
// @description     家庭饮食管理系统的后端API接口文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.onetastefamily.com/support
// @contact.email  support@onetastefamily.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 使用 "Bearer {token}" 格式，token 通过登录接口获取

func main() {
	fmt.Println("OneTasteFamily Backend Service")
	log.Println("Server starting...")

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化 MinIO 客户端
	minioCfg := storage.MinIOConfig{
		Endpoint:  config.AppConfig.MinIO.Endpoint,
		AccessKey: config.AppConfig.MinIO.AccessKey,
		SecretKey: config.AppConfig.MinIO.SecretKey,
		Bucket:    config.AppConfig.MinIO.Bucket,
		Region:    config.AppConfig.MinIO.Region,
		UseSSL:    config.AppConfig.MinIO.UseSSL,
		BaseURL:   config.AppConfig.MinIO.BaseURL,
	}

	if err := storage.InitMinIO(minioCfg); err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	// 初始化数据库连接
	dbCfg := database.Config{
		Host:            config.AppConfig.Database.Host,
		Port:            config.AppConfig.Database.Port,
		User:            config.AppConfig.Database.User,
		Password:        config.AppConfig.Database.Password,
		DBName:          config.AppConfig.Database.DBName,
		SSLMode:         config.AppConfig.Database.SSLMode,
		MaxOpenConns:    config.AppConfig.Database.MaxOpenConns,
		MaxIdleConns:    config.AppConfig.Database.MaxIdleConns,
		ConnMaxLifetime: config.AppConfig.Database.ConnMaxLifetime,
		ConnMaxIdleTime: config.AppConfig.Database.ConnMaxIdleTime,
	}

	if err := database.Connect(dbCfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Database connected successfully")

	// 初始化 Redis
	redisCfg := cache.RedisConfig{
		Host:     config.AppConfig.Redis.Host,
		Port:     config.AppConfig.Redis.Port,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	}

	if err := cache.InitRedis(redisCfg); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer cache.CloseRedis()
	log.Println("Redis connected successfully")

	// 设置Gin模式
	if config.AppConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 添加CORS中间件
	r.Use(corsMiddleware())

	// 注册所有路由
	handlers.RegisterAllRoutes(r)

	// 健康检查
	r.GET("/health", healthCheck)

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// healthCheck 健康检查
// @Summary 健康检查
// @Description 检查服务是否正常运行
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
