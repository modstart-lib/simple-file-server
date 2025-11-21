# 编译目标
TARGET = simple-file-server

# 源文件目录
SRCDIR = .

# 编译输出目录
BINDIR = build

# 构建命令
GO = ~/sdk/go1.22.9/bin/go
BUILD = $(GO) build

# 源文件
SOURCES := $(shell find $(SRCDIR) -name '*.go')

# 默认目标，编译程序
default: clean $(BINDIR)/$(TARGET) sync

# 构建目标
$(BINDIR)/$(TARGET): $(SOURCES)
	@mkdir -p $(BINDIR)
	for GOOS in darwin linux; do \
	   for GOARCH in amd64 arm64 ; do \
		 export GOOS=$$GOOS \
		 	&& export GOARCH=$$GOARCH \
		 	&& $(BUILD) -trimpath -o $(BINDIR)/$(TARGET)-$$GOOS-$$GOARCH $(SRCDIR); \
	   done \
	done

# 清理
clean:
	@rm -rf $(BINDIR)

sync:
	sleep 1

