#!/bin/bash

# PostgreSQL 容器修复脚本

set -e

echo "=========================================="
echo "PostgreSQL 容器修复工具"
echo "=========================================="
echo ""

# 1. 停止所有容器
echo "1. 停止所有容器..."
docker compose down
echo ""

# 2. 检查数据目录
echo "2. 检查数据目录..."
if [ -d "../data/postgres" ]; then
    echo "   数据目录存在，检查权限..."
    ls -ld ../data/postgres
    
    # 检查是否有损坏的数据
    if [ -d "../data/postgres/pgdata" ] && [ -f "../data/postgres/pgdata/PG_VERSION" ]; then
        echo "   检测到已有数据，询问是否清理..."
        read -p "   是否清理现有数据？(y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "   清理数据目录..."
            rm -rf ../data/postgres/*
            echo "   ✓ 数据目录已清理"
        else
            echo "   保留现有数据"
        fi
    fi
    
    # 确保目录权限正确
    echo "   设置目录权限..."
    chmod 700 ../data/postgres
    echo "   ✓ 权限已设置"
else
    echo "   创建数据目录..."
    mkdir -p ../data/postgres
    chmod 700 ../data/postgres
    echo "   ✓ 数据目录已创建"
fi
echo ""

# 3. 检查环境变量
echo "3. 检查环境变量..."
if [ ! -f ".env" ]; then
    echo "   .env 文件不存在，从示例创建..."
    if [ -f ".env.example" ]; then
        cp .env.example .env
        echo "   ✓ .env 文件已创建"
        echo "   请编辑 .env 文件设置密码："
        echo "   nano .env"
    else
        echo "   ✗ .env.example 文件不存在"
        echo "   请手动创建 .env 文件"
        exit 1
    fi
else
    echo "   ✓ .env 文件存在"
fi
echo ""

# 4. 检查端口占用
echo "4. 检查端口占用..."
if command -v ss &> /dev/null; then
    if ss -tuln | grep -q ":5432 "; then
        echo "   ⚠ 端口 5432 已被占用"
        echo "   占用端口的进程："
        sudo lsof -i :5432 2>/dev/null || ss -tulnp | grep 5432
        echo ""
        read -p "   是否继续？(y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        echo "   ✓ 端口 5432 可用"
    fi
fi
echo ""

# 5. 启动 PostgreSQL 容器
echo "5. 启动 PostgreSQL 容器..."
docker compose up -d postgres
echo ""

# 6. 等待容器启动
echo "6. 等待容器启动（最多 30 秒）..."
for i in {1..30}; do
    if docker compose exec postgres pg_isready -U postgres > /dev/null 2>&1; then
        echo "   ✓ PostgreSQL 已就绪"
        break
    fi
    echo "   等待中... ($i/30)"
    sleep 1
done

# 7. 检查容器状态
echo ""
echo "7. 检查容器状态..."
docker compose ps postgres
echo ""

# 8. 查看日志
if ! docker compose exec postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo "8. 容器未就绪，查看日志："
    docker compose logs --tail=50 postgres
    echo ""
    echo "=========================================="
    echo "修复失败，请检查日志"
    echo "=========================================="
    exit 1
fi

# 9. 测试连接
echo "8. 测试数据库连接..."
docker compose exec postgres psql -U postgres -d onetaste_family -c "SELECT version();" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "   ✓ 数据库连接成功"
else
    echo "   ⚠ 数据库连接失败，但容器已启动"
fi
echo ""

echo "=========================================="
echo "修复完成！"
echo "=========================================="
echo ""
echo "下一步："
echo "1. 启动所有服务: docker compose up -d"
echo "2. 查看状态: docker compose ps"
echo "3. 查看日志: docker compose logs -f"

