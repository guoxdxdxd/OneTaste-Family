#!/bin/bash

# PostgreSQL 备份脚本
# 使用方法: ./backup_postgres.sh

set -e

BACKUP_DIR="../backups/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/onetaste_family_${DATE}.dump"
CONTAINER_NAME="onetaste_postgres"
DB_NAME="onetaste_family"
DB_USER="postgres"

# 创建备份目录
mkdir -p ${BACKUP_DIR}

echo "Starting PostgreSQL backup..."

# 执行备份
docker exec ${CONTAINER_NAME} pg_dump -U ${DB_USER} -F c ${DB_NAME} > ${BACKUP_FILE}

# 检查备份是否成功
if [ $? -eq 0 ]; then
    echo "Backup completed successfully: ${BACKUP_FILE}"
    
    # 压缩备份文件
    gzip ${BACKUP_FILE}
    echo "Backup compressed: ${BACKUP_FILE}.gz"
    
    # 显示备份文件大小
    BACKUP_SIZE=$(du -h "${BACKUP_FILE}.gz" | cut -f1)
    echo "Backup size: ${BACKUP_SIZE}"
    
    # 删除 30 天前的备份
    find ${BACKUP_DIR} -name "*.dump.gz" -mtime +30 -delete
    echo "Old backups (older than 30 days) have been deleted."
else
    echo "Backup failed!"
    exit 1
fi

