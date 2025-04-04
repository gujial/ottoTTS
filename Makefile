# Makefile for building ottoTTScli project

# 定义变量
GO = go
MAIN_FILE = cli/ottoTTScli/ottoTTScli.go
BINARY_NAME = ottoTTScli
BUILD_DIR = build
ASSETS_DIR = assets

# 默认目标
all: build

# 编译 Go 程序到 build 目录
build: $(MAIN_FILE)
	@echo "Building Go program..."
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# 复制 assets 文件夹到 build 目录
copy-assets:
	@echo "Copying assets folder to build directory..."
	cp -r $(ASSETS_DIR) $(BUILD_DIR)/

# 清理生成的文件
clean:
	@echo "Cleaning up build directory..."
	rm -rf $(BUILD_DIR)

# 完整的构建流程，包括编译和复制 assets
full-build: build copy-assets

# 运行 Go 程序
run: $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Running the application..."
	$(BUILD_DIR)/$(BINARY_NAME)

# 安装 Go 程序
install:
	$(GO) install $(MAIN_FILE)
