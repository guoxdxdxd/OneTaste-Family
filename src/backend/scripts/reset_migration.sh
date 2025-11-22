#!/bin/bash

# 清空数据并重新迁移的脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 加载配置（复用 migrate.sh 的配置加载逻辑）
CONFIG_PATH="${CONFIG_PATH:-config/config.yaml}"
MIGRATE_PATH="./migrations"

# 从配置文件读取数据库配置
load_config_from_yaml() {
    local config_file="$1"
    
    if [ ! -f "$config_file" ]; then
        return 1
    fi
    
    # 使用简单的 sed/awk 解析
    local temp_file=$(mktemp)
    awk '/^database:/{flag=1; next} /^[a-z_]+:/ && flag {flag=0} flag' "$config_file" > "$temp_file"
    
    DB_HOST=$(grep -E "^\s+host:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*host:[[:space:]]*//' | sed -E 's/^["'\''](.*)["'\'']$/\1/' | sed 's/[[:space:]].*$//' | head -1)
    DB_PORT=$(grep -E "^\s+port:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*port:[[:space:]]*//' | sed 's/[[:space:]].*$//' | head -1)
    DB_USER=$(grep -E "^\s+user:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*user:[[:space:]]*//' | sed -E 's/^["'\''](.*)["'\'']$/\1/' | sed 's/[[:space:]].*$//' | head -1)
    DB_PASSWORD=$(grep -E "^\s+password:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*password:[[:space:]]*//' | sed -E 's/^["'\''](.*)["'\'']$/\1/' | sed 's/[[:space:]].*$//' | head -1)
    DB_NAME=$(grep -E "^\s+dbname:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*dbname:[[:space:]]*//' | sed -E 's/^["'\''](.*)["'\'']$/\1/' | sed 's/[[:space:]].*$//' | head -1)
    DB_SSLMODE=$(grep -E "^\s+sslmode:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*sslmode:[[:space:]]*//' | sed -E 's/^["'\''](.*)["'\'']$/\1/' | sed 's/[[:space:]].*$//' | head -1)
    
    rm -f "$temp_file"
    
    if [ -n "$DB_HOST" ] || [ -n "$DB_USER" ]; then
        return 0
    fi
    
    return 1
}

# 加载配置
if [ -f "$CONFIG_PATH" ]; then
    if load_config_from_yaml "$CONFIG_PATH"; then
        echo -e "${GREEN}已从配置文件加载数据库配置: $CONFIG_PATH${NC}"
    fi
fi

# 环境变量覆盖
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-onetaste_family}"
DB_SSLMODE="${DB_SSLMODE:-disable}"

# URL 编码函数
urlencode() {
    local string="${1}"
    local strlen=${#string}
    local encoded=""
    local pos c o

    for (( pos=0 ; pos<strlen ; pos++ )); do
        c=${string:$pos:1}
        case "$c" in
            [-_.~a-zA-Z0-9] ) o="${c}" ;;
            * ) printf -v o '%%%02x' "'$c"
        esac
        encoded+="${o}"
    done
    echo "${encoded}"
}

# 对密码进行 URL 编码
if [[ "$DB_PASSWORD" =~ [^a-zA-Z0-9._~-] ]]; then
    DB_PASSWORD_ENCODED=$(urlencode "$DB_PASSWORD")
else
    DB_PASSWORD_ENCODED="$DB_PASSWORD"
fi

DB_URL="postgres://${DB_USER}:${DB_PASSWORD_ENCODED}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

# 检查 migrate 命令
if ! command -v migrate &> /dev/null; then
    echo -e "${RED}错误: migrate 命令未找到${NC}"
    exit 1
fi

# 检查当前迁移状态
CURRENT_VERSION=$(migrate -path "${MIGRATE_PATH}" -database "${DB_URL}" version 2>&1 | grep -oE '[0-9]+' | head -1 || echo "0")
IS_DIRTY=$(migrate -path "${MIGRATE_PATH}" -database "${DB_URL}" version 2>&1 | grep -i "dirty" || echo "")

echo -e "${YELLOW}当前迁移版本: $CURRENT_VERSION${IS_DIRTY:+ (dirty)}${NC}"

# 如果处于dirty状态，先修复
if [ -n "$IS_DIRTY" ]; then
    echo -e "${YELLOW}检测到 dirty 状态，正在修复...${NC}"
    PREV_VERSION=$((CURRENT_VERSION - 1))
    if [ $PREV_VERSION -lt 0 ]; then
        PREV_VERSION=0
    fi
    migrate -path "${MIGRATE_PATH}" -database "${DB_URL}" force "$PREV_VERSION"
    echo -e "${GREEN}已修复 dirty 状态${NC}"
    CURRENT_VERSION=$PREV_VERSION
fi

# 确认操作
if [ "$CURRENT_VERSION" != "0" ]; then
    echo -e "${RED}警告: 将回滚所有迁移并清空数据${NC}"
    echo -e "${YELLOW}当前版本: $CURRENT_VERSION${NC}"
    read -p "确认继续? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "已取消"
        exit 0
    fi
fi

# 回滚所有迁移
if [ "$CURRENT_VERSION" != "0" ]; then
    echo -e "${YELLOW}正在回滚所有迁移...${NC}"
    migrate -path "${MIGRATE_PATH}" -database "${DB_URL}" down -all
    echo -e "${GREEN}回滚完成${NC}"
fi

# 重新执行所有迁移
echo -e "${GREEN}正在执行所有迁移...${NC}"
migrate -path "${MIGRATE_PATH}" -database "${DB_URL}" up
echo -e "${GREEN}迁移完成${NC}"

# 显示最终版本
FINAL_VERSION=$(migrate -path "${MIGRATE_PATH}" -database "${DB_URL}" version 2>&1 | grep -oE '[0-9]+' | head -1 || echo "0")
echo -e "${GREEN}最终迁移版本: $FINAL_VERSION${NC}"

