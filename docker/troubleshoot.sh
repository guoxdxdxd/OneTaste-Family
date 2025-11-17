#!/bin/bash

# Docker 故障排查脚本
# 用于诊断 PostgreSQL 和 Redis 容器启动问题

echo "=========================================="
echo "Docker 容器故障排查工具"
echo "=========================================="
echo ""

# 1. 检查 Docker 服务状态
echo "1. 检查 Docker 服务状态..."
if systemctl is-active --quiet docker; then
    echo "   ✓ Docker 服务正在运行"
else
    echo "   ✗ Docker 服务未运行"
    echo "   请执行: sudo systemctl start docker"
    exit 1
fi
echo ""

# 2. 检查容器状态
echo "2. 检查容器状态..."
docker compose ps
echo ""

# 3. 查看 PostgreSQL 日志
echo "3. 查看 PostgreSQL 容器日志（最近 50 行）..."
docker compose logs --tail=50 postgres
echo ""

# 4. 检查数据目录
echo "4. 检查数据目录..."
if [ -d "../data/postgres" ]; then
    echo "   ✓ data/postgres 目录存在"
    ls -la ../data/postgres | head -5
    echo "   目录权限:"
    stat -c "%a %n" ../data/postgres 2>/dev/null || stat -f "%OLp %N" ../data/postgres
else
    echo "   ✗ data/postgres 目录不存在"
    echo "   请执行: mkdir -p ../data/postgres && chmod 700 ../data/postgres"
fi
echo ""

# 5. 检查端口占用
echo "5. 检查端口占用..."
if command -v ss &> /dev/null; then
    echo "   PostgreSQL 端口 5432:"
    ss -tuln | grep 5432 || echo "   端口 5432 未被占用"
    echo "   Redis 端口 6379:"
    ss -tuln | grep 6379 || echo "   端口 6379 未被占用"
else
    echo "   使用 netstat 检查端口..."
    netstat -tuln 2>/dev/null | grep -E "5432|6379" || echo "   端口未被占用"
fi
echo ""

# 6. 检查磁盘空间
echo "6. 检查磁盘空间..."
df -h . | tail -1
echo ""

# 7. 检查环境变量文件
echo "7. 检查环境变量文件..."
if [ -f ".env" ]; then
    echo "   ✓ .env 文件存在"
    if grep -q "POSTGRES_PASSWORD" .env && grep -q "REDIS_PASSWORD" .env; then
        echo "   ✓ 密码已配置"
    else
        echo "   ✗ 密码未配置"
    fi
else
    echo "   ✗ .env 文件不存在"
    echo "   请执行: cp .env.example .env && nano .env"
fi
echo ""

# 8. 检查 Docker 网络
echo "8. 检查 Docker 网络..."
docker network ls | grep onetaste_network || echo "   网络不存在"
echo ""

# 9. 尝试进入容器（如果容器存在）
echo "9. 尝试检查容器内部..."
if docker ps -a | grep -q onetaste_postgres; then
    echo "   尝试检查 PostgreSQL 容器..."
    docker compose exec postgres pg_isready -U postgres 2>&1 || echo "   无法连接到容器"
fi
echo ""

echo "=========================================="
echo "排查完成！"
echo "=========================================="
echo ""
echo "常见解决方案："
echo "1. 如果数据目录权限问题："
echo "   chmod 700 ../data/postgres"
echo ""
echo "2. 如果端口被占用："
echo "   修改 docker-compose.yml 中的端口映射"
echo ""
echo "3. 如果容器启动失败："
echo "   docker compose down"
echo "   docker compose up -d"
echo ""
echo "4. 如果数据目录损坏："
echo "   docker compose down"
echo "   rm -rf ../data/postgres/*"
echo "   docker compose up -d postgres"

