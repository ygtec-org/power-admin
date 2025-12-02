# CMS 服务 Makefile

.PHONY: gen build run clean

# 生成代码
gen:
	@echo "Generating code from cms.api..."
	goctl api go --api ./api/cms.api --home ./common/template --dir . --style gozero

# 编译
build: gen
	@echo "Building power-cms-server..."
	go build -o bin/power-cms-server.exe .

# 运行
run: build
	@echo "Running power-cms-server..."
	./bin/power-cms-server.exe

# 清理
clean:
	@echo "Cleaning..."
	rm -rf bin
	go clean
