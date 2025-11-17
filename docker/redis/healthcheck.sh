#!/bin/sh
# Redis 健康检查脚本

# 从环境变量读取密码，如果没有则使用默认值
REDIS_PASSWORD=${REDIS_PASSWORD:-ChangeMe123!}

# 执行 ping 命令
redis-cli -a "${REDIS_PASSWORD}" ping

