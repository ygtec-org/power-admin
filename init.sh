#!/bin/bash

# Power Admin 项目初始化脚本
# 用于快速设置开发环境

set -e

echo "========================================="
echo "  Power Admin - 项目初始化脚本"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 检查环境
echo -e "${BLUE}[1] 检查系统环境...${NC}"

# 检查Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go 未安装${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Go $(go version | awk '{print $3}')${NC}"

# 检查MySQL
if ! command -v mysql &> /dev/null; then
    echo -e "${RED}✗ MySQL 未安装${NC}"
    exit 1
fi
echo -e "${GREEN}✓ MySQL 已安装${NC}"

# 检查Redis
if ! command -v redis-cli &> /dev/null; then
    echo -e "${YELLOW}⚠ Redis 未安装（但可继续）${NC}"
fi

# 检查Node
if ! command -v node &> /dev/null; then
    echo -e "${RED}✗ Node.js 未安装${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Node.js $(node --version)${NC}"

echo ""
echo -e "${BLUE}[2] 初始化数据库...${NC}"

# 提示用户输入MySQL密码
read -sp "请输入 MySQL root 密码: " MYSQL_PASSWORD
echo ""

# 创建数据库并导入数据
if [ -z "$MYSQL_PASSWORD" ]; then
    mysql -u root < power-admin-server/db/init.sql
else
    mysql -u root -p"$MYSQL_PASSWORD" < power-admin-server/db/init.sql
fi

echo -e "${GREEN}✓ 数据库初始化完成${NC}"

echo ""
echo -e "${BLUE}[3] 编译后端...${NC}"

cd power-admin-server
go mod tidy
go build -o power-admin.exe

echo -e "${GREEN}✓ 后端编译完成${NC}"

cd ..

echo ""
echo -e "${BLUE}[4] 初始化前端...${NC}"

cd power-admin-web
npm install

echo -e "${GREEN}✓ 前端依赖安装完成${NC}"

cd ..

echo ""
echo "========================================="
echo -e "${GREEN}✓ 初始化完成！${NC}"
echo "========================================="
echo ""
echo "后续步骤："
echo ""
echo "1. 启动后端服务："
echo -e "   ${BLUE}cd power-admin-server${NC}"
echo -e "   ${BLUE}./power-admin.exe -f etc/power-api.yaml${NC}"
echo ""
echo "2. 启动前端开发服务器（新终端）："
echo -e "   ${BLUE}cd power-admin-web${NC}"
echo -e "   ${BLUE}npm run dev${NC}"
echo ""
echo "3. 打开浏览器访问："
echo -e "   ${BLUE}http://localhost:5173${NC}"
echo ""
echo "4. 使用以下默认凭证登录："
echo -e "   ${YELLOW}手机号：13800138000${NC}"
echo -e "   ${YELLOW}密码：admin123${NC}"
echo ""
echo "更多信息请参考："
echo "  - 快速启动指南: QUICKSTART.md"
echo "  - 开发指南: DEVELOPMENT_GUIDE.md"
echo "  - 架构设计: ARCHITECTURE.md"
echo ""
