#!/bin/bash

# 项目名和输出目录
PROJECT_NAME=${1:-myapp}
OUTPUT_DIR=${2:-bin}

# 创建输出目录
mkdir -p $OUTPUT_DIR

# 本地编译
echo "编译本地执行文件..."
GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) go build -o $OUTPUT_DIR/$PROJECT_NAME-$GOOS-$GOARCH

# 交叉编译
#PLATFORMS=("windows/amd64" "linux/amd64" "darwin/amd64")
PLATFORMS=("windows/amd64" "linux/amd64")
echo "开始交叉编译..."
for PLATFORM in "${PLATFORMS[@]}"
do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    OUTPUT_NAME=$PROJECT_NAME-$GOOS-$GOARCH
    if [ $GOOS = "windows" ]; then
        OUTPUT_NAME+='.exe'
    fi
    
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT_DIR/$OUTPUT_NAME
    if [ $? -ne 0 ]; then
        echo '交叉编译失败:' $PLATFORM
        exit 1
    fi
done

echo "编译完成."

