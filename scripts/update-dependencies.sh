#!/bin/bash
# 更新所有子模块间的依赖版本
# 用法：./scripts/update-dependencies.sh v1.0.0

set -e

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "❌ 错误：请提供版本号"
    echo "用法：$0 v1.0.0"
    exit 1
fi

echo "🔄 更新所有子模块依赖版本为：$VERSION"

# 根模块路径
ROOT_MODULE="github.com/kalandramo/lulu-ext"

# 统计更新的文件数
UPDATED=0

# 遍历所有 go.mod 文件
while IFS= read -r -d '' GO_MOD; do
    # 跳过根目录的 go.mod（如果有）
    if [[ "$GO_MOD" == "./go.mod" ]]; then
        continue
    fi

    # 检查是否引用了根模块或其他子模块
    if grep -q "$ROOT_MODULE" "$GO_MOD"; then
        # 备份原文件
        cp "$GO_MOD" "$GO_MOD.bak"

        # 替换所有子模块依赖的版本号为指定版本
        sed -i.bak2 "s|$ROOT_MODULE/[a-zA-Z0-9_-]* v0\.0\.0-[0-9]\{14\}-[a-f0-9]\{12\}|$ROOT_MODULE $VERSION|g" "$GO_MOD"
        sed -i.bak3 "s|$ROOT_MODULE v0\.0\.0-[0-9]\{14\}-[a-f0-9]\{12\}|$ROOT_MODULE $VERSION|g" "$GO_MOD"

        # 清理备份文件
        rm -f "$GO_MOD.bak" "$GO_MOD.bak2" "$GO_MOD.bak3"

        # 检查文件是否真的有变化
        if ! diff -q "$GO_MOD" "$GO_MOD" >/dev/null 2>&1 || grep -q "$VERSION" "$GO_MOD"; then
            echo "  ✅ $(dirname "$GO_MOD")/go.mod"
            UPDATED=$((UPDATED + 1))
        fi
    fi
done < <(find . -name "go.mod" -not -path "./.git/*" -print0)

echo ""
echo "✅ 完成！更新了 $UPDATED 个 go.mod 文件"

# 运行 go mod tidy 清理依赖
echo ""
echo "🧹 清理依赖..."
go mod tidy
echo "✅ 清理完成"
