#!/bin/bash

# 修复迁移 dirty 状态的脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 加载配置（复用 migrate.sh 的配置加载逻辑）
CONFIG_PATH="${CONFIG_PATH:-config/config.yaml}"

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
CURRENT_VERSION=$(migrate -path ./migrations -database "${DB_URL}" version 2>&1 | grep -oE '[0-9]+' | head -1 || echo "0")
IS_DIRTY=$(migrate -path ./migrations -database "${DB_URL}" version 2>&1 | grep -i "dirty" || echo "")

if [ -n "$IS_DIRTY" ]; then
    echo -e "${YELLOW}检测到 dirty 状态，当前版本: $CURRENT_VERSION${NC}"
    echo -e "${YELLOW}正在修复...${NC}"
    
    # 强制设置到上一个版本
    PREV_VERSION=$((CURRENT_VERSION - 1))
    if [ $PREV_VERSION -lt 0 ]; then
        PREV_VERSION=0
    fi
    
    echo -e "${YELLOW}强制设置迁移版本为: $PREV_VERSION${NC}"
    migrate -path ./migrations -database "${DB_URL}" force "$PREV_VERSION"
    
    echo -e "${GREEN}修复完成，现在可以重新执行迁移: ./scripts/migrate.sh up${NC}"
else
    echo -e "${GREEN}迁移状态正常，当前版本: $CURRENT_VERSION${NC}"
fi

