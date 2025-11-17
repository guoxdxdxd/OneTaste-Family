#!/bin/bash

# 数据库迁移脚本
# 使用 golang-migrate 工具执行数据库迁移

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 默认配置
MIGRATE_PATH="./migrations"
CONFIG_PATH="${CONFIG_PATH:-config/config.yaml}"

# 从配置文件读取数据库配置的函数
# 使用简单的 bash 方法解析 YAML（适用于简单的配置结构）
load_config_from_yaml() {
    local config_file="$1"
    
    if [ ! -f "$config_file" ]; then
        return 1
    fi
    
    # 方法1: 使用 Python 解析 YAML（如果可用且安装了 PyYAML）
    if command -v python3 &> /dev/null; then
        python3 -c "import yaml" 2>/dev/null
        if [ $? -eq 0 ]; then
            local db_config=$(python3 <<EOF
import yaml
import sys

try:
    with open('$config_file', 'r') as f:
        config = yaml.safe_load(f)
        if config and 'database' in config:
            db = config['database']
            print(f"{db.get('host', 'localhost')}")
            print(f"{db.get('port', 5432)}")
            print(f"{db.get('user', 'postgres')}")
            print(f"{db.get('password', 'postgres')}")
            print(f"{db.get('dbname', 'onetaste_family')}")
            print(f"{db.get('sslmode', 'disable')}")
            sys.exit(0)
except Exception as e:
    sys.exit(1)
EOF
)
            if [ $? -eq 0 ] && [ -n "$db_config" ]; then
                DB_HOST=$(echo "$db_config" | sed -n '1p')
                DB_PORT=$(echo "$db_config" | sed -n '2p')
                DB_USER=$(echo "$db_config" | sed -n '3p')
                DB_PASSWORD=$(echo "$db_config" | sed -n '4p')
                DB_NAME=$(echo "$db_config" | sed -n '5p')
                DB_SSLMODE=$(echo "$db_config" | sed -n '6p')
                return 0
            fi
        fi
    fi
    
    # 方法2: 使用 yq（如果安装了）
    if command -v yq &> /dev/null; then
        DB_HOST=$(yq eval '.database.host // "localhost"' "$config_file" 2>/dev/null || echo "localhost")
        DB_PORT=$(yq eval '.database.port // 5432' "$config_file" 2>/dev/null || echo "5432")
        DB_USER=$(yq eval '.database.user // "postgres"' "$config_file" 2>/dev/null || echo "postgres")
        DB_PASSWORD=$(yq eval '.database.password // "postgres"' "$config_file" 2>/dev/null || echo "postgres")
        DB_NAME=$(yq eval '.database.dbname // "onetaste_family"' "$config_file" 2>/dev/null || echo "onetaste_family")
        DB_SSLMODE=$(yq eval '.database.sslmode // "disable"' "$config_file" 2>/dev/null || echo "disable")
        return 0
    fi
    
    # 方法3: 使用简单的 sed/awk 解析（适用于简单的 YAML 结构）
    # 提取 database 部分的配置值
    local temp_file=$(mktemp)
    
    # 提取 database 部分（包括缩进的行）
    awk '/^database:/{flag=1; next} /^[a-z_]+:/ && flag {flag=0} flag' "$config_file" > "$temp_file"
    
    # 使用更简单的 sed 提取各个配置项的值
    # 处理带引号和不带引号的值
    DB_HOST=$(grep -E "^\s+host:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*host:[[:space:]]*//' | sed -E 's/^["'\''](.*)["'\'']$/\1/' | sed 's/[[:space:]].*$//' | head -1)
    DB_PORT=$(grep -E "^\s+port:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*port:[[:space:]]*//' | sed 's/[[:space:]].*$//' | head -1)
    DB_USER=$(grep -E "^\s+user:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*user:[[:space:]]*//' | sed -E 's/^["'\''](.*)["'\'']$/\1/' | sed 's/[[:space:]].*$//' | head -1)
    DB_PASSWORD=$(grep -E "^\s+password:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*password:[[:space:]]*//' | sed -E 's/^["'\''](.*)["'\'']$/\1/' | sed 's/[[:space:]].*$//' | head -1)
    DB_NAME=$(grep -E "^\s+dbname:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*dbname:[[:space:]]*//' | sed -E 's/^["'\''](.*)["'\'']$/\1/' | sed 's/[[:space:]].*$//' | head -1)
    DB_SSLMODE=$(grep -E "^\s+sslmode:" "$temp_file" 2>/dev/null | sed -E 's/^[[:space:]]*sslmode:[[:space:]]*//' | sed -E 's/^["'\''](.*)["'\'']$/\1/' | sed 's/[[:space:]].*$//' | head -1)
    
    rm -f "$temp_file"
    
    # 如果成功解析了至少一个值，返回成功
    if [ -n "$DB_HOST" ] || [ -n "$DB_USER" ]; then
        return 0
    fi
    
    return 1
}

# 尝试从配置文件加载
if [ -f "$CONFIG_PATH" ]; then
    if load_config_from_yaml "$CONFIG_PATH"; then
        echo -e "${GREEN}已从配置文件加载数据库配置: $CONFIG_PATH${NC}"
    fi
fi

# 环境变量优先级最高，可以覆盖配置文件的值
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-onetaste_family}"
DB_SSLMODE="${DB_SSLMODE:-disable}"

# URL 编码函数（处理特殊字符）
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

# 对密码进行 URL 编码（如果包含特殊字符）
if [[ "$DB_PASSWORD" =~ [^a-zA-Z0-9._~-] ]]; then
    DB_PASSWORD_ENCODED=$(urlencode "$DB_PASSWORD")
else
    DB_PASSWORD_ENCODED="$DB_PASSWORD"
fi

# 构建数据库连接字符串
DB_URL="postgres://${DB_USER}:${DB_PASSWORD_ENCODED}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

# 检查 migrate 命令是否存在
if ! command -v migrate &> /dev/null; then
    echo -e "${RED}错误: migrate 命令未找到${NC}"
    echo -e "${YELLOW}请安装 golang-migrate:${NC}"
    echo "  macOS: brew install golang-migrate"
    echo "  Linux: https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md"
    exit 1
fi

# 显示帮助信息
show_help() {
    echo "数据库迁移工具"
    echo ""
    echo "用法:"
    echo "  $0 [命令]"
    echo ""
    echo "命令:"
    echo "  up         执行所有待执行的迁移"
    echo "  down       回滚最后一次迁移"
    echo "  down-all   回滚所有迁移"
    echo "  version    显示当前迁移版本"
    echo "  create     创建新的迁移文件"
    echo "  force      强制设置迁移版本（危险操作）"
    echo ""
    echo "配置来源（按优先级）:"
    echo "  1. 环境变量（优先级最高）"
    echo "  2. 配置文件: config/config.yaml (或 CONFIG_PATH 指定的路径)"
    echo "  3. 默认值"
    echo ""
    echo "环境变量:"
    echo "  DB_HOST     数据库主机 (默认: localhost)"
    echo "  DB_PORT     数据库端口 (默认: 5432)"
    echo "  DB_USER     数据库用户 (默认: postgres)"
    echo "  DB_PASSWORD 数据库密码 (默认: postgres)"
    echo "  DB_NAME     数据库名称 (默认: onetaste_family)"
    echo "  DB_SSLMODE  SSL模式 (默认: disable)"
    echo "  CONFIG_PATH 配置文件路径 (默认: config/config.yaml)"
    echo ""
    echo "示例:"
    echo "  $0 up                    # 执行迁移"
    echo "  $0 down                  # 回滚一次"
    echo "  $0 create add_new_table  # 创建新迁移文件"
}

# 执行迁移
case "${1:-up}" in
    up)
        echo -e "${GREEN}执行数据库迁移...${NC}"
        migrate -path "${MIGRATE_PATH}" -database "${DB_URL}" up
        echo -e "${GREEN}迁移完成${NC}"
        ;;
    down)
        echo -e "${YELLOW}回滚最后一次迁移...${NC}"
        migrate -path "${MIGRATE_PATH}" -database "${DB_URL}" down 1
        echo -e "${GREEN}回滚完成${NC}"
        ;;
    down-all)
        echo -e "${RED}警告: 将回滚所有迁移${NC}"
        read -p "确认继续? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            migrate -path "${MIGRATE_PATH}" -database "${DB_URL}" down -all
            echo -e "${GREEN}回滚完成${NC}"
        else
            echo "已取消"
        fi
        ;;
    version)
        migrate -path "${MIGRATE_PATH}" -database "${DB_URL}" version
        ;;
    create)
        if [ -z "$2" ]; then
            echo -e "${RED}错误: 请提供迁移文件名${NC}"
            echo "用法: $0 create <迁移文件名>"
            exit 1
        fi
        migrate create -ext sql -dir "${MIGRATE_PATH}" -seq "$2"
        echo -e "${GREEN}迁移文件已创建${NC}"
        ;;
    force)
        if [ -z "$2" ]; then
            echo -e "${RED}错误: 请提供版本号${NC}"
            echo "用法: $0 force <版本号>"
            exit 1
        fi
        echo -e "${RED}警告: 强制设置迁移版本可能导致数据不一致${NC}"
        read -p "确认继续? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            migrate -path "${MIGRATE_PATH}" -database "${DB_URL}" force "$2"
            echo -e "${GREEN}版本已设置${NC}"
        else
            echo "已取消"
        fi
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo -e "${RED}未知命令: $1${NC}"
        show_help
        exit 1
        ;;
esac

