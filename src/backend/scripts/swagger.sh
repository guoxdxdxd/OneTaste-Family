#!/bin/bash

# Swagger 文档生成脚本

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"

cd "$PROJECT_DIR"

# 检查 swag 是否安装
if ! command -v swag &> /dev/null; then
    echo "swag 未安装，正在安装..."
    go install github.com/swaggo/swag/cmd/swag@latest
    
    # 检查安装是否成功
    if ! command -v swag &> /dev/null; then
        echo "❌ swag 安装失败，请检查 GOPATH/bin 是否在 PATH 中"
        exit 1
    fi
fi

# 生成 Swagger 文档
echo "正在生成 Swagger 文档..."
swag init -g cmd/main.go -o docs/swagger

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Swagger 文档生成成功！"
    echo "📄 文档位置: docs/swagger/"
    echo "🌐 访问地址: http://localhost:8080/swagger/index.html"
    echo ""
    echo "提示："
    echo "  1. 启动服务器后访问 http://localhost:8080/swagger/index.html"
    echo "  2. 修改代码后重新运行此脚本更新文档"
else
    echo "❌ Swagger 文档生成失败！"
    exit 1
fi

