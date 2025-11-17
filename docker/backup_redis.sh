#!/bin/bash

# Redis 备份脚本
# 使用方法: ./backup_redis.sh

set -e

BACKUP_DIR="../backups/redis"
DATE=$(date +%Y%m%d_%H%M%S)
CONTAINER_NAME="onetaste_redis"

# 从环境变量读取 Redis 密码，如果没有则使用默认值
REDIS_PASSWORD=${REDIS_PASSWORD:-ChangeMe123!}

# 创建备份目录
mkdir -p ${BACKUP_DIR}

echo "Starting Redis backup..."

# 创建 RDB 快照
docker exec ${CONTAINER_NAME} redis-cli -a ${REDIS_PASSWORD} SAVE

# 检查 SAVE 是否成功
if [ $? -eq 0 ]; then
    echo "RDB snapshot created successfully"
    
    # 复制 RDB 文件
    docker cp ${CONTAINER_NAME}:/data/dump.rdb ${BACKUP_DIR}/dump_${DATE}.rdb
    
    # 复制 AOF 文件（如果存在）
    if docker exec ${CONTAINER_NAME} test -f /data/appendonly.aof; then
        docker cp ${CONTAINER_NAME}:/data/appendonly.aof ${BACKUP_DIR}/appendonly_${DATE}.aof
        echo "AOF file backed up: appendonly_${DATE}.aof"
    fi
    
    # 压缩备份文件
    cd ${BACKUP_DIR}
    tar -czf redis_backup_${DATE}.tar.gz dump_${DATE}.rdb appendonly_${DATE}.aof 2>/dev/null || tar -czf redis_backup_${DATE}.tar.gz dump_${DATE}.rdb
    cd - > /dev/null
    
    # 删除未压缩的文件
    rm -f ${BACKUP_DIR}/dump_${DATE}.rdb ${BACKUP_DIR}/appendonly_${DATE}.aof
    
    echo "Backup completed successfully: ${BACKUP_DIR}/redis_backup_${DATE}.tar.gz"
    
    # 显示备份文件大小
    BACKUP_SIZE=$(du -h "${BACKUP_DIR}/redis_backup_${DATE}.tar.gz" | cut -f1)
    echo "Backup size: ${BACKUP_SIZE}"
    
    # 删除 30 天前的备份
    find ${BACKUP_DIR} -name "redis_backup_*.tar.gz" -mtime +30 -delete
    echo "Old backups (older than 30 days) have been deleted."
else
    echo "Backup failed!"
    exit 1
fi

