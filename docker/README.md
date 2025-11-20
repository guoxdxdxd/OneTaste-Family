# Docker 部署快速指南

## 快速开始

### 1. 初始化环境

```bash
cd docker
chmod +x init.sh
./init.sh
```

### 2. 配置环境变量

编辑 `.env` 文件，设置强密码：

```bash
nano .env
```

**重要**：生产环境必须修改默认密码！

### 3. 启动服务

```bash
docker-compose up -d
```

### 4. 验证服务

```bash
# 查看容器状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 测试 PostgreSQL 连接
docker-compose exec postgres psql -U postgres -d onetaste_family -c "SELECT version();"

# 测试 Redis 连接
docker-compose exec redis redis-cli -a YourRedisPassword ping

# 检查 MinIO 状态
curl -f http://localhost:9000/minio/health/live
# 或打开 http://localhost:9090 使用 MINIO_ROOT_USER/MINIO_ROOT_PASSWORD 登录控制台
```

## 常用命令

```bash
# 启动服务
docker compose up -d

# 停止服务
docker compose stop

# 重启服务
docker compose restart

# 查看日志
docker compose logs -f [service_name]

# 进入容器
docker compose exec postgres bash
docker compose exec redis sh
docker compose exec minio sh

# 查看单个服务日志
docker compose logs -f postgres
docker compose logs -f redis
docker compose logs -f minio

# 备份 PostgreSQL
./backup_postgres.sh

# 备份 Redis
./backup_redis.sh

# 使用 mc 临时容器查看 MinIO 桶
MINIO_ROOT_USER=$(grep MINIO_ROOT_USER .env | cut -d '=' -f2)
MINIO_ROOT_PASSWORD=$(grep MINIO_ROOT_PASSWORD .env | cut -d '=' -f2)
docker run --rm --network host \
  -e MC_HOST_local="http://${MINIO_ROOT_USER}:${MINIO_ROOT_PASSWORD}@localhost:9000" \
  minio/mc ls local
```

## 故障排查

### PostgreSQL 容器启动失败

如果遇到 `container onetaste_postgres is unhealthy` 错误：

```bash
# 方法 1: 使用自动修复脚本（推荐）
./fix_postgres.sh

# 方法 2: 手动排查
./troubleshoot.sh

# 方法 3: 查看详细日志
docker compose logs postgres

# 方法 4: 清理并重新启动
docker compose down
rm -rf ../data/postgres/*
chmod 700 ../data/postgres
docker compose up -d postgres
```

### 常见问题

1. **数据目录权限问题**
   ```bash
   chmod 700 ../data/postgres
   ```

2. **端口被占用**
   ```bash
   # 检查端口占用
   ss -tuln | grep 5432
   # 或修改 docker-compose.yml 中的端口
   ```

3. **数据目录损坏**
   ```bash
   docker compose down
   rm -rf ../data/postgres/*
   docker compose up -d postgres
   ```

4. **MinIO 控制台无法访问**
   ```bash
   docker compose logs minio
   curl http://localhost:9000/minio/health/live
   ```

## 文件说明

- `docker-compose.yml` - Docker Compose 配置文件
- `.env.example` - 环境变量示例文件
- `.env` - 环境变量文件（需要创建，不要提交到 Git）
- `postgres/init.sql` - PostgreSQL 初始化脚本
- `redis/redis.conf` - Redis 配置文件
- `backup_postgres.sh` - PostgreSQL 备份脚本
- `backup_redis.sh` - Redis 备份脚本
- `init.sh` - 初始化脚本
- `redis/healthcheck.sh` - Redis 健康检查脚本（被容器挂载使用）

## 数据目录

- `../data/postgres/` - PostgreSQL 数据目录
- `../data/redis/` - Redis 数据目录
- `../data/minio/` - MinIO 数据目录
- `../backups/postgres/` - PostgreSQL 备份目录
- `../backups/redis/` - Redis 备份目录
- `../backups/minio/` - MinIO 备份目录（手动创建，用于 mc 同步备份）

## 连接信息

### PostgreSQL
- 主机: `localhost` (容器内: `postgres`)
- 端口: `5432`
- 用户: `postgres`
- 数据库: `onetaste_family`
- 密码: `.env` 文件中的 `POSTGRES_PASSWORD`

### Redis
- 主机: `localhost` (容器内: `redis`)
- 端口: `6379`
- 密码: `.env` 文件中的 `REDIS_PASSWORD`

### MinIO
- API: `http://localhost:9000` (容器内: `http://minio:9000`)
- Console: `http://localhost:9090`
- 访问凭据: `.env` 文件中的 `MINIO_ROOT_USER` / `MINIO_ROOT_PASSWORD`
- MC 示例地址: `MC_HOST_local="http://USER:PASSWORD@localhost:9000"`

## 安全提示

1. **必须修改默认密码**：生产环境使用强密码
2. **不要提交 .env 文件**：添加到 `.gitignore`
3. **定期备份**：设置自动备份任务
4. **限制网络访问**：生产环境不要暴露端口到公网

## 详细文档

查看 [DOCKER_DEPLOYMENT.md](../docs/DOCKER_DEPLOYMENT.md) 获取完整的部署文档。
