@echo off
chcp 65001 > nul
echo.
echo ╔════════════════════════════════════════════════════════╗
echo ║         Power Admin 项目 - 初始化脚本                  ║
echo ╚════════════════════════════════════════════════════════╝
echo.

setlocal enabledelayedexpansion

:: 检查数据库
echo [1/5] 检查数据库连接...
mysql -u root -proot -e "SELECT 1" >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误: 无法连接到 MySQL 数据库
    echo 请确保:
    echo   1. MySQL 服务正在运行
    echo   2. 用户名和密码正确 (默认: root/root)
    echo.
    pause
    exit /b 1
)
echo ✅ MySQL 连接成功

:: 创建数据库和表
echo [2/5] 初始化数据库...
mysql -u root -proot -e "CREATE DATABASE IF NOT EXISTS power_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"
if errorlevel 1 (
    echo ❌ 错误: 创建数据库失败
    pause
    exit /b 1
)

:: 导入 SQL 脚本
echo [3/5] 导入数据库脚本...
mysql -u root -proot power_admin < power-admin-server\db\init.sql
if errorlevel 1 (
    echo ❌ 错误: 导入 SQL 脚本失败
    echo 请检查 SQL 文件是否存在: power-admin-server\db\init.sql
    pause
    exit /b 1
)
echo ✅ 数据库初始化成功

:: 检查 Go 环境
echo [4/5] 检查 Go 环境...
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误: 未找到 Go 编译器
    echo 请先安装 Go (https://golang.org)
    pause
    exit /b 1
)
echo ✅ Go 环境就绪

:: 编译后端
echo [5/5] 编译后端程序...
cd power-admin-server
go mod download
go build -o power-admin.exe
if errorlevel 1 (
    echo ❌ 错误: 后端编译失败
    cd ..
    pause
    exit /b 1
)
cd ..
echo ✅ 后端编译成功

echo.
echo ╔════════════════════════════════════════════════════════╗
echo ║              初始化完成! 可以开始运行了                ║
echo ╚════════════════════════════════════════════════════════╝
echo.
echo 下一步操作:
echo.
echo 【终端 1 - 启动后端】
echo   cd power-admin-server
echo   .\power-admin.exe -f etc\power-api.yaml
echo.
echo 【终端 2 - 启动前端】
echo   cd power-admin-web
echo   npm install
echo   npm run dev
echo.
echo 【浏览器访问】
echo   http://localhost:5173
echo.
echo 【默认登录凭证】
echo   手机号: 13800138000
echo   密码: admin123
echo.
pause
