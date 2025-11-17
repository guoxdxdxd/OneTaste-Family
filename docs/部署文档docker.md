# Docker 部署文档

## 概述

本文档介绍如何使用 Docker 容器部署 OneTaste Family 项目的数据库服务（PostgreSQL 和 Redis）。

## 系统要求

### 操作系统
- **推荐**：Ubuntu 22.04 LTS 64位（UEFI）
- 其他支持 Docker 的 Linux 发行版

### 硬件要求
- CPU：2 核或以上
- 内存：4GB 或以上（推荐 8GB）
- 磁盘空间：至少 20GB 可用空间
- 网络：稳定的网络连接

### 软件要求
- Docker 20.10+
- Docker Compose 2.0+
- Git（用于克隆项目）

## Ubuntu 22.04 安装 Docker

### 1. 更新系统包

```bash
sudo apt update
sudo apt upgrade -y
```

### 2. 安装必要的依赖

```bash
sudo apt install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
```

### 3. 添加 Docker 官方 GPG 密钥

```bash
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
```

### 4. 设置 Docker 仓库

```bash
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

### 5. 安装 Docker Engine

```bash
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

### 6. 验证安装

```bash
sudo docker --version
sudo docker compose version
```

### 7. 将当前用户添加到 docker 组（可选，避免每次使用 sudo）

```bash
sudo usermod -aG docker $USER
# 重新登录或执行以下命令使组权限生效
newgrp docker
```

### 8. 测试 Docker

```bash
docker run hello-world
```

## 安装 Docker Compose（如果使用独立版本）

如果系统安装的是 Docker Compose 独立版本（非插件版本），可以这样安装：

```bash
# 下载最新版本
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# 设置执行权限
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker-compose --version
```

## Ubuntu 22.04 特定配置

### 1. 配置 Docker 服务自启动

```bash
sudo systemctl enable docker
sudo systemctl start docker
```

### 2. 配置防火墙（如果使用 UFW）

```bash
# 允许 Docker 网络
sudo ufw allow 5432/tcp  # PostgreSQL
sudo ufw allow 6379/tcp  # Redis

# 或者允许 Docker 自动管理防火墙规则
sudo ufw allow from 172.16.0.0/12
```

### 3. 优化系统参数（可选，提升性能）

编辑 `/etc/sysctl.conf`：

```bash
sudo nano /etc/sysctl.conf
```

添加以下配置：

```conf
# 网络优化
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 2048

# 内存优化
vm.overcommit_memory = 1
vm.swappiness = 10
```

应用配置：

```bash
sudo sysctl -p
```

### 4. 配置 Docker 日志大小限制

编辑 `/etc/docker/daemon.json`：

```bash
sudo mkdir -p /etc/docker
sudo nano /etc/docker/daemon.json
```

添加以下配置：

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "storage-driver": "overlay2"
}
```

重启 Docker 服务：

```bash
sudo systemctl restart docker
```

## 目录结构

```
OneTasteFamily/
├── docker/
│   ├── docker-compose.yml          # Docker Compose 配置文件
│   ├── postgres/
│   │   └── init.sql                 # PostgreSQL 初始化脚本（可选）
│   └── redis/
│       └── redis.conf               # Redis 配置文件（可选）
├── data/
│   ├── postgres/                    # PostgreSQL 数据目录（自动创建）
│   └── redis/                       # Redis 数据目录（自动创建）
└── docs/
    └── DOCKER_DEPLOYMENT.md         # 本文档
```

## Docker Compose 配置

创建 `docker/docker-compose.yml` 文件：

```yaml
version: '3.8'

services:
  # PostgreSQL 数据库服务
  postgres:
    image: postgres:latest
    container_name: onetaste_postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-ChangeMe123!}
      POSTGRES_DB: onetaste_family
      PGDATA: /var/lib/postgresql/data/pgdata
      TZ: Asia/Shanghai
    ports:
      - "${POSTGRES_PORT:-5432}:5432"
    volumes:
      # 挂载数据目录到宿主机
      - ../data/postgres:/var/lib/postgresql/data
      # 可选：挂载初始化脚本
      - ./postgres/init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    networks:
      - onetaste_network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    command:
      - "postgres"
      - "-c"
      - "max_connections=200"
      - "-c"
      - "shared_buffers=256MB"
      - "-c"
      - "effective_cache_size=1GB"
      - "-c"
      - "maintenance_work_mem=64MB"
      - "-c"
      - "checkpoint_completion_target=0.9"
      - "-c"
      - "wal_buffers=16MB"
      - "-c"
      - "default_statistics_target=100"
      - "-c"
      - "random_page_cost=1.1"
      - "-c"
      - "effective_io_concurrency=200"
      - "-c"
      - "work_mem=4MB"
      - "-c"
      - "min_wal_size=1GB"
      - "-c"
      - "max_wal_size=4GB"

  # Redis 缓存服务
  redis:
    image: redis:latest
    container_name: onetaste_redis
    restart: unless-stopped
    command: redis-server --appendonly yes --requirepass ${REDIS_PASSWORD:-ChangeMe123!}
    ports:
      - "${REDIS_PORT:-6379}:6379"
    volumes:
      # 挂载数据目录到宿主机
      - ../data/redis:/data
      # 可选：挂载 Redis 配置文件
      - ./redis/redis.conf:/usr/local/etc/redis/redis.conf:ro
    networks:
      - onetaste_network
    healthcheck:
      test: ["CMD", "redis-cli", "--raw", "incr", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5
    depends_on:
      postgres:
        condition: service_healthy

networks:
  onetaste_network:
    driver: bridge
    name: onetaste_network
```

## 环境变量配置

创建 `docker/.env` 文件（不要提交到版本控制）：

```env
# PostgreSQL 配置
POSTGRES_PASSWORD=YourSecurePassword123!
POSTGRES_PORT=5432

# Redis 配置
REDIS_PASSWORD=YourSecurePassword123!
REDIS_PORT=6379
```

**安全提示**：
- 生产环境必须使用强密码
- 不要将 `.env` 文件提交到 Git 仓库
- 建议使用密码管理工具生成随机密码

## 部署步骤

### 1. 克隆项目（如果还没有）

```bash
git clone <repository-url>
cd OneTasteFamily
```

### 2. 创建必要的目录

```bash
# 在项目根目录执行
mkdir -p data/postgres
mkdir -p data/redis
mkdir -p backups/postgres
mkdir -p backups/redis

# 设置目录权限（PostgreSQL 需要特定权限）
chmod 700 data/postgres
chmod 755 data/redis
chmod 755 backups/postgres
chmod 755 backups/redis
```

**注意**：在 Ubuntu 22.04 上，如果使用非 root 用户，确保有足够的权限创建和访问这些目录。

### 3. 创建环境变量文件

```bash
cd docker
cp .env.example .env  # 如果有示例文件
# 或者直接创建
nano .env
```

编辑 `.env` 文件，设置强密码。在 Ubuntu 22.04 上，可以使用以下命令生成随机密码：

```bash
# 生成 PostgreSQL 密码
openssl rand -base64 32

# 生成 Redis 密码
openssl rand -base64 32
```

### 4. 启动服务

```bash
cd docker
docker compose up -d
```

**注意**：在 Ubuntu 22.04 上，如果使用 Docker Compose 插件版本，使用 `docker compose`（没有连字符）。如果是独立版本，使用 `docker-compose`（有连字符）。

### 5. 验证服务状态

```bash
# 查看容器状态
docker compose ps

# 查看日志
docker compose logs -f

# 检查 PostgreSQL 健康状态
docker compose exec postgres pg_isready -U postgres

# 检查 Redis 健康状态（从 .env 文件读取密码）
docker compose exec redis redis-cli -a $(grep REDIS_PASSWORD .env | cut -d '=' -f2) ping
```

### 6. 测试连接

#### PostgreSQL 连接测试

```bash
# 进入 PostgreSQL 容器
docker compose exec postgres psql -U postgres -d onetaste_family

# 在 psql 中执行
\dt  # 查看表列表
SELECT version();  # 查看版本信息
\q   # 退出
```

#### Redis 连接测试

```bash
# 从环境变量读取密码并连接
REDIS_PASSWORD=$(grep REDIS_PASSWORD .env | cut -d '=' -f2)
docker compose exec redis redis-cli -a "${REDIS_PASSWORD}"

# 在 redis-cli 中执行
PING  # 应该返回 PONG
INFO  # 查看 Redis 信息
SET test_key "test_value"
GET test_key
EXIT  # 退出
```

或者直接使用：

```bash
docker compose exec redis redis-cli -a $(grep REDIS_PASSWORD .env | cut -d '=' -f2)
```

## 数据目录说明

### PostgreSQL 数据目录

- **挂载路径**：`../data/postgres:/var/lib/postgresql/data`
- **数据位置**：`data/postgres/`
- **包含内容**：
  - 数据库文件
  - WAL 日志
  - 配置文件
  - 事务日志

### Redis 数据目录

- **挂载路径**：`../data/redis:/data`
- **数据位置**：`data/redis/`
- **包含内容**：
  - AOF 文件（appendonly.aof）
  - RDB 快照文件（dump.rdb，如果启用）

## 常用操作

### 启动服务

```bash
cd docker
docker compose up -d
```

### 停止服务

```bash
cd docker
docker compose stop
```

### 停止并删除容器（保留数据）

```bash
cd docker
docker compose down
```

### 停止并删除容器和数据卷（危险操作）

```bash
cd docker
docker compose down -v
```

### 查看日志

```bash
# 查看所有服务日志
docker compose logs -f

# 查看 PostgreSQL 日志
docker compose logs -f postgres

# 查看 Redis 日志
docker compose logs -f redis

# 查看最近 100 行日志
docker compose logs --tail=100 -f
```

### 重启服务

```bash
cd docker
docker compose restart

# 重启特定服务
docker compose restart postgres
docker compose restart redis
```

### 进入容器

```bash
# 进入 PostgreSQL 容器
docker compose exec postgres bash

# 进入 Redis 容器
docker compose exec redis sh
```

## 数据备份

### PostgreSQL 备份

#### 1. 使用 pg_dump 备份

```bash
# 备份单个数据库
docker compose exec postgres pg_dump -U postgres onetaste_family > backup_$(date +%Y%m%d_%H%M%S).sql

# 备份为自定义格式（推荐，支持压缩）
docker compose exec postgres pg_dump -U postgres -F c -f /tmp/backup.dump onetaste_family
docker compose cp postgres:/tmp/backup.dump ./backup_$(date +%Y%m%d_%H%M%S).dump
```

#### 2. 备份整个数据目录

```bash
# 停止 PostgreSQL 容器
docker compose stop postgres

# 备份数据目录（在项目根目录执行）
cd ..
tar -czf postgres_backup_$(date +%Y%m%d_%H%M%S).tar.gz data/postgres/

# 启动 PostgreSQL 容器
cd docker
docker compose start postgres
```

#### 3. 自动备份脚本

创建 `docker/backup_postgres.sh`：

```bash
#!/bin/bash

BACKUP_DIR="../backups/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/onetaste_family_${DATE}.dump"

mkdir -p ${BACKUP_DIR}

# 执行备份
docker compose exec -T postgres pg_dump -U postgres -F c onetaste_family > ${BACKUP_FILE}

# 压缩备份文件
gzip ${BACKUP_FILE}

# 删除 30 天前的备份
find ${BACKUP_DIR} -name "*.dump.gz" -mtime +30 -delete

echo "Backup completed: ${BACKUP_FILE}.gz"
```

设置执行权限：

```bash
chmod +x docker/backup_postgres.sh
```

添加到 crontab（每天凌晨 2 点执行）：

```bash
0 2 * * * cd /path/to/OneTasteFamily/docker && ./backup_postgres.sh
```

### Redis 备份

#### 1. 使用 redis-cli 备份

```bash
# 从环境变量读取密码
REDIS_PASSWORD=$(grep REDIS_PASSWORD .env | cut -d '=' -f2)

# 创建 RDB 快照
docker compose exec redis redis-cli -a "${REDIS_PASSWORD}" SAVE

# 复制 RDB 文件
docker compose cp redis:/data/dump.rdb ./redis_backup_$(date +%Y%m%d_%H%M%S).rdb
```

#### 2. 备份 AOF 文件

```bash
# 复制 AOF 文件
docker compose cp redis:/data/appendonly.aof ./redis_backup_$(date +%Y%m%d_%H%M%S).aof
```

#### 3. 备份整个数据目录

```bash
# 停止 Redis 容器
docker compose stop redis

# 备份数据目录（在项目根目录执行）
cd ..
tar -czf redis_backup_$(date +%Y%m%d_%H%M%S).tar.gz data/redis/

# 启动 Redis 容器
cd docker
docker compose start redis
```

## 数据恢复

### PostgreSQL 恢复

#### 1. 从 SQL 文件恢复

```bash
# 恢复 SQL 备份
docker compose exec -T postgres psql -U postgres -d onetaste_family < backup.sql
```

#### 2. 从自定义格式恢复

```bash
# 复制备份文件到容器
docker compose cp backup.dump postgres:/tmp/backup.dump

# 恢复数据库
docker compose exec postgres pg_restore -U postgres -d onetaste_family -c /tmp/backup.dump
```

#### 3. 从数据目录恢复

```bash
# 停止 PostgreSQL
docker compose stop postgres

# 恢复数据目录（在项目根目录执行）
cd ..
rm -rf data/postgres/*
tar -xzf postgres_backup_20240115_020000.tar.gz

# 启动 PostgreSQL
cd docker
docker compose start postgres
```

### Redis 恢复

#### 1. 从 RDB 文件恢复

```bash
# 停止 Redis
docker compose stop redis

# 复制 RDB 文件到数据目录（在项目根目录执行）
cd ..
cp redis_backup.rdb data/redis/dump.rdb

# 启动 Redis
cd docker
docker compose start redis
```

#### 2. 从 AOF 文件恢复

```bash
# 停止 Redis
docker compose stop redis

# 复制 AOF 文件到数据目录（在项目根目录执行）
cd ..
cp redis_backup.aof data/redis/appendonly.aof

# 启动 Redis
cd docker
docker compose start redis
```

## 性能优化

### PostgreSQL 优化

在 `docker-compose.yml` 中已经包含了一些性能优化参数，可以根据服务器资源调整：

```yaml
command:
  - "postgres"
  - "-c"
  - "max_connections=200"              # 最大连接数
  - "-c"
  - "shared_buffers=256MB"             # 共享缓冲区（建议为内存的 25%）
  - "-c"
  - "effective_cache_size=1GB"         # 有效缓存大小（建议为内存的 50-75%）
  - "-c"
  - "maintenance_work_mem=64MB"         # 维护工作内存
  - "-c"
  - "checkpoint_completion_target=0.9"  # 检查点完成目标
  - "-c"
  - "wal_buffers=16MB"                  # WAL 缓冲区
  - "-c"
  - "default_statistics_target=100"     # 默认统计目标
  - "-c"
  - "random_page_cost=1.1"              # 随机页面成本（SSD 建议 1.1）
  - "-c"
  - "effective_io_concurrency=200"      # 有效 IO 并发数（SSD 建议 200）
  - "-c"
  - "work_mem=4MB"                      # 工作内存
  - "-c"
  - "min_wal_size=1GB"                  # 最小 WAL 大小
  - "-c"
  - "max_wal_size=4GB"                  # 最大 WAL 大小
```

### Redis 优化

创建 `docker/redis/redis.conf`：

```conf
# 内存限制（根据实际情况调整）
maxmemory 2gb
maxmemory-policy allkeys-lru

# 持久化配置
appendonly yes
appendfsync everysec

# 性能优化
tcp-backlog 511
timeout 0
tcp-keepalive 300

# 日志级别
loglevel notice
```

在 `docker-compose.yml` 中挂载配置文件：

```yaml
volumes:
  - ../data/redis:/data
  - ./redis/redis.conf:/usr/local/etc/redis/redis.conf:ro
command: redis-server /usr/local/etc/redis/redis.conf --requirepass ${REDIS_PASSWORD}
```

## 监控和维护

### 查看资源使用情况

```bash
# 查看容器资源使用
docker stats onetaste_postgres onetaste_redis

# 查看磁盘使用（在项目根目录执行）
cd ..
du -sh data/postgres data/redis

# 查看详细的磁盘使用情况
du -h --max-depth=1 data/postgres
du -h --max-depth=1 data/redis
```

### 清理未使用的数据

```bash
# PostgreSQL 清理
docker compose exec postgres psql -U postgres -d onetaste_family -c "VACUUM ANALYZE;"

# Redis 清理（如果设置了 maxmemory-policy）
REDIS_PASSWORD=$(grep REDIS_PASSWORD .env | cut -d '=' -f2)
docker compose exec redis redis-cli -a "${REDIS_PASSWORD}" FLUSHDB  # 危险操作，谨慎使用
```

### 查看数据库大小

```bash
# PostgreSQL 数据库大小
docker compose exec postgres psql -U postgres -d onetaste_family -c "
SELECT 
    pg_database.datname,
    pg_size_pretty(pg_database_size(pg_database.datname)) AS size
FROM pg_database
WHERE datname = 'onetaste_family';
"

# 查看表大小
docker compose exec postgres psql -U postgres -d onetaste_family -c "
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
"
```

## 故障排查

### PostgreSQL 无法启动

1. **检查数据目录权限**
```bash
# 在项目根目录执行
ls -la data/postgres
# 确保目录权限为 700
# 如果权限不对，执行：
chmod 700 data/postgres
```

2. **查看日志**
```bash
cd docker
docker compose logs postgres
```

3. **检查端口占用（Ubuntu 22.04）**
```bash
# 使用 ss 命令（推荐）
ss -tuln | grep 5432

# 或使用 netstat（需要安装 net-tools）
sudo apt install net-tools
netstat -tuln | grep 5432

# 或使用 lsof
sudo apt install lsof
sudo lsof -i :5432
```

4. **检查 SELinux（Ubuntu 通常不启用，但可以检查）**
```bash
# Ubuntu 默认不启用 SELinux，但如果启用了，可能需要调整
getenforce  # 如果返回 Disabled，则不需要处理
```

5. **重置数据目录（会丢失数据）**
```bash
cd docker
docker compose down
cd ..
rm -rf data/postgres/*
cd docker
docker compose up -d postgres
```

### Redis 无法启动

1. **查看日志**
```bash
cd docker
docker compose logs redis
```

2. **检查端口占用**
```bash
ss -tuln | grep 6379
# 或
sudo lsof -i :6379
```

3. **检查数据目录权限**
```bash
# 在项目根目录执行
ls -la data/redis
# 确保目录权限为 755
chmod 755 data/redis
```

4. **检查 Redis 配置**
```bash
cd docker
docker compose exec redis cat /usr/local/etc/redis/redis.conf
```

### 连接问题

1. **检查网络**
```bash
docker network inspect onetaste_network
```

2. **测试连接**
```bash
# PostgreSQL
cd docker
docker compose exec postgres psql -U postgres -d onetaste_family -c "SELECT 1;"

# Redis
REDIS_PASSWORD=$(grep REDIS_PASSWORD .env | cut -d '=' -f2)
docker compose exec redis redis-cli -a "${REDIS_PASSWORD}" ping
```

3. **检查防火墙（Ubuntu 22.04 使用 UFW）**
```bash
# 查看防火墙状态
sudo ufw status

# 如果防火墙启用，开放端口
sudo ufw allow 5432/tcp
sudo ufw allow 6379/tcp

# 或者允许 Docker 网络
sudo ufw allow from 172.16.0.0/12
```

4. **检查 Docker 服务状态**
```bash
sudo systemctl status docker
# 如果服务未运行，启动它
sudo systemctl start docker
```

## 安全建议

1. **使用强密码**
   - 生产环境必须使用复杂的随机密码
   - 定期更换密码

2. **限制网络访问**
   - 生产环境不要将端口暴露到公网
   - 使用防火墙限制访问来源

3. **定期更新镜像**
   ```bash
   docker-compose pull
   docker-compose up -d
   ```

4. **启用 SSL/TLS**
   - PostgreSQL 可以配置 SSL 连接
   - Redis 6.0+ 支持 TLS

5. **定期备份**
   - 设置自动备份任务
   - 测试备份恢复流程

## 升级指南

### 升级 PostgreSQL

```bash
cd docker

# 1. 备份数据
./backup_postgres.sh

# 2. 停止服务
docker compose stop postgres

# 3. 更新镜像版本（在 docker-compose.yml 中）
# image: postgres:15  # 指定版本
# 或者拉取最新镜像
docker compose pull postgres

# 4. 启动新版本
docker compose up -d postgres

# 5. 验证
docker compose exec postgres psql -U postgres -d onetaste_family -c "SELECT version();"
```

### 升级 Redis

```bash
cd docker

# 1. 备份数据
./backup_redis.sh

# 2. 停止服务
docker compose stop redis

# 3. 更新镜像版本
# image: redis:7  # 指定版本
# 或者拉取最新镜像
docker compose pull redis

# 4. 启动新版本
docker compose up -d redis

# 5. 验证
REDIS_PASSWORD=$(grep REDIS_PASSWORD .env | cut -d '=' -f2)
docker compose exec redis redis-cli -a "${REDIS_PASSWORD}" ping
```

## 快速参考

### 常用命令

```bash
cd docker

# 启动所有服务
docker compose up -d

# 停止所有服务
docker compose stop

# 查看状态
docker compose ps

# 查看日志
docker compose logs -f

# 重启服务
docker compose restart

# 进入 PostgreSQL
docker compose exec postgres psql -U postgres -d onetaste_family

# 进入 Redis（自动读取密码）
REDIS_PASSWORD=$(grep REDIS_PASSWORD .env | cut -d '=' -f2)
docker compose exec redis redis-cli -a "${REDIS_PASSWORD}"

# 或者创建一个别名（添加到 ~/.bashrc）
alias redis-cli-docker='docker compose exec redis redis-cli -a $(grep REDIS_PASSWORD .env | cut -d "=" -f2)'
```

### 连接信息

- **PostgreSQL**
  - 主机：`localhost`（容器内：`postgres`）
  - 端口：`5432`
  - 用户：`postgres`
  - 数据库：`onetaste_family`
  - 密码：`.env` 文件中的 `POSTGRES_PASSWORD`

- **Redis**
  - 主机：`localhost`（容器内：`redis`）
  - 端口：`6379`
  - 密码：`.env` 文件中的 `REDIS_PASSWORD`

## 注意事项

1. **数据持久化**：数据目录已挂载到宿主机，删除容器不会丢失数据
2. **端口冲突**：确保 5432 和 6379 端口未被占用
3. **权限问题**：PostgreSQL 数据目录需要正确的权限设置
4. **密码安全**：生产环境必须使用强密码，不要使用默认密码
5. **备份策略**：定期备份数据，测试恢复流程
6. **资源限制**：根据服务器资源调整 PostgreSQL 和 Redis 的配置参数

## 相关文档

- [数据库设计文档](./DATABASE.md)
- [部署文档](./DEPLOYMENT.md)
- [API 文档](./API.md)

