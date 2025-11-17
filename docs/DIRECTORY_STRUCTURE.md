# 项目目录结构

本文档描述了 OneTasteFamily 项目的完整目录结构。

## 根目录结构

```
OneTasteFamily/
├── src/                    # 源代码目录
│   ├── backend/            # Go后端服务
│   ├── frontend/           # 前端项目
│   │   ├── admin/          # 管理后台（Vue）
│   │   └── miniapp/        # 小程序（Vue）
│   └── ai-service/         # Python AI服务
├── docker/                 # Docker配置和脚本
│   ├── postgres/           # PostgreSQL配置
│   └── redis/              # Redis配置
├── docs/                   # 项目文档
└── README.md              # 项目说明文档
```

## 详细目录结构

### docker/ - Docker配置目录

```
docker/
├── postgres/              # PostgreSQL相关配置
│   └── init.sql          # 数据库初始化脚本
│
├── redis/                 # Redis相关配置
│   ├── redis.conf        # Redis配置文件
│   └── healthcheck.sh    # Redis健康检查脚本
│
├── docker-compose.yml     # Docker Compose配置文件
├── init.sh                # 初始化脚本
├── backup_postgres.sh     # PostgreSQL备份脚本
├── backup_redis.sh        # Redis备份脚本
├── fix_postgres.sh        # PostgreSQL修复脚本
├── troubleshoot.sh        # 故障排查脚本
├── .env                   # 环境变量文件（不提交到版本控制）
├── .env.example           # 环境变量示例文件
└── README.md              # Docker部署说明文档
```

### docs/ - 文档目录

```
docs/
├── API.md                 # API接口文档
├── DATABASE.md            # 数据库设计文档
├── DEPLOYMENT.md          # 部署文档
├── DOCKER_DEPLOYMENT.md   # Docker部署文档
├── DIRECTORY_STRUCTURE.md # 目录结构文档（本文件）
├── REQUIREMENTS.md        # 需求文档
├── TECH_STACK.md          # 技术栈文档
└── UBUNTU_SETUP.md        # Ubuntu环境搭建文档
```

## 各模块说明

### Docker配置

**位置**: `docker/`

**主要职责**:
- 容器化部署配置
- 数据库初始化
- 服务编排
- 备份和恢复脚本

**包含服务**:
- PostgreSQL数据库
- Redis缓存
- Nginx反向代理（配置在docker-compose.yml中引用）

## 文件命名规范

### Go文件
- 使用小写字母和下划线：`user_service.go`
- 测试文件：`user_service_test.go`
- 主程序：`main.go`

### Python文件
- 使用小写字母和下划线：`user_service.py`
- 测试文件：`test_user_service.py`
- 主程序：`main.py`

### Vue文件
- 组件文件：PascalCase，如 `UserProfile.vue`
- 页面文件：PascalCase，如 `HomePage.vue`
- 工具文件：camelCase，如 `apiClient.js`

### 配置文件
- YAML配置：`config.yaml`
- 环境变量：`.env`（不提交），`.env.example`（提交）
- Docker配置：`docker-compose.yml`

## 代码模块详细目录结构

各代码模块的详细目录结构请参考以下文档：

- [Backend 目录结构](./BACKEND_DIRECTORY.md) - Go后端服务详细目录结构
- [Frontend Admin 目录结构](./FRONTEND_ADMIN_DIRECTORY.md) - 管理后台详细目录结构
- [Frontend Miniapp 目录结构](./FRONTEND_MINIAPP_DIRECTORY.md) - 小程序详细目录结构
- [AI Service 目录结构](./AI_SERVICE_DIRECTORY.md) - Python AI服务详细目录结构

## 注意事项

1. **环境变量文件**: `.env` 文件包含敏感信息，不应提交到版本控制系统
2. **配置文件**: 各服务应使用配置文件或环境变量管理配置，避免硬编码
3. **日志文件**: 日志文件应统一管理，建议使用日志目录或日志服务
4. **临时文件**: 临时文件和构建产物不应提交到版本控制
5. **文档更新**: 当目录结构发生变化时，应及时更新本文档

## 版本历史

- **v1.0.0** (2024-01-XX): 初始目录结构创建

