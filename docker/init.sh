#!/bin/bash

# Docker 初始化脚本
# 用于创建必要的目录和设置权限

set -e

echo "Initializing Docker environment..."

# 创建数据目录
echo "Creating data directories..."
mkdir -p ../data/postgres
mkdir -p ../data/redis
mkdir -p ../backups/postgres
mkdir -p ../backups/redis

# 设置目录权限和所有者
echo "Setting directory permissions and ownership..."

# PostgreSQL 容器内的 postgres 用户通常是 UID 999
# 尝试设置目录所有者为 UID 999
POSTGRES_OWNER_SET=false
if id 999 >/dev/null 2>&1; then
    # 系统中有 UID 999 的用户
    if chown 999:999 ../data/postgres 2>/dev/null; then
        echo "   ✓ 已设置 postgres 目录所有者为 UID 999"
        POSTGRES_OWNER_SET=true
    else
        echo "   ⚠ 无法设置所有者，尝试使用 sudo..."
        if sudo chown 999:999 ../data/postgres 2>/dev/null; then
            echo "   ✓ 已使用 sudo 设置 postgres 目录所有者"
            POSTGRES_OWNER_SET=true
        else
            echo "   ⚠ 无法设置所有者，将使用权限 777（临时方案，仅开发环境）"
            chmod 777 ../data/postgres
        fi
    fi
else
    # 系统中没有 UID 999 的用户（如 macOS），尝试使用 sudo
    echo "   ⚠ 系统中没有 UID 999 的用户，尝试使用 sudo..."
    if sudo chown 999:999 ../data/postgres 2>/dev/null; then
        echo "   ✓ 已使用 sudo 设置 postgres 目录所有者"
        POSTGRES_OWNER_SET=true
    else
        echo "   ⚠ 无法设置所有者，将使用权限 777（临时方案，仅开发环境）"
        chmod 777 ../data/postgres
    fi
fi

# 如果成功设置了所有者，使用 700 权限；否则保持 777
if [ "$POSTGRES_OWNER_SET" = true ]; then
    chmod 700 ../data/postgres
else
    chmod 777 ../data/postgres
fi

# Redis 容器内的 redis 用户通常是 UID 999 或 1000
REDIS_OWNER_SET=false
if id 999 >/dev/null 2>&1; then
    if chown 999:999 ../data/redis 2>/dev/null; then
        REDIS_OWNER_SET=true
    elif sudo chown 999:999 ../data/redis 2>/dev/null; then
        REDIS_OWNER_SET=true
    elif id 1000 >/dev/null 2>&1; then
        if chown 1000:1000 ../data/redis 2>/dev/null; then
            REDIS_OWNER_SET=true
        elif sudo chown 1000:1000 ../data/redis 2>/dev/null; then
            REDIS_OWNER_SET=true
        else
            chmod 777 ../data/redis
        fi
    else
        if sudo chown 1000:1000 ../data/redis 2>/dev/null; then
            REDIS_OWNER_SET=true
        else
            chmod 777 ../data/redis
        fi
    fi
else
    if sudo chown 999:999 ../data/redis 2>/dev/null; then
        REDIS_OWNER_SET=true
    elif sudo chown 1000:1000 ../data/redis 2>/dev/null; then
        REDIS_OWNER_SET=true
    else
        chmod 777 ../data/redis
    fi
fi
# Redis 使用 755 权限即可，因为容器内的 redis 用户只需要读写权限
chmod 755 ../data/redis

chmod 755 ../backups/postgres
chmod 755 ../backups/redis

# 创建 .env 文件（如果不存在）
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo "Please edit .env file and set secure passwords!"
fi

# 设置脚本执行权限
echo "Setting script permissions..."
chmod +x backup_postgres.sh
chmod +x backup_redis.sh

echo "Initialization completed!"
echo ""
echo "Next steps:"
echo "1. Edit .env file and set secure passwords"
echo "2. Run: docker compose up -d"
echo "3. Check status: docker compose ps"

