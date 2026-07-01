#!/bin/bash
# 统一版本自动打 tag 脚本
# 用法：./scripts/tag-all.sh v1.0.0 [--dry-run]

set -e

VERSION=$1
DRY_RUN=false

# 解析参数
if [ "$2" = "--dry-run" ]; then
    DRY_RUN=true
    echo "🔍 干运行模式 - 不会实际创建或推送 tag"
fi

# 如果未提供版本号，尝试从 bump-version.sh 自动计算
if [ -z "$VERSION" ]; then
    if [ -x "./scripts/bump-version.sh" ]; then
        VERSION=$(./scripts/bump-version.sh patch)
        echo "📦 未提供版本号，自动计算：$VERSION"
    else
        echo "❌ 错误：请提供版本号"
        echo "用法：$0 v1.0.0 [--dry-run]"
        echo "或：$0 [--dry-run]（自动计算下一个 patch 版本）"
        exit 1
    fi
fi

# 验证版本号格式
if ! echo "$VERSION" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$'; then
    echo "❌ 错误：版本号格式应为 v1.0.0 (语义化版本)"
    exit 1
fi

# 检查是否已存在该 tag
if git tag -l "$VERSION" | grep -q "^$VERSION$"; then
    echo "❌ 错误：tag $VERSION 已存在"
    exit 1
fi

# 检查是否有未提交的更改
if ! git diff-index --quiet HEAD -- 2>/dev/null; then
    echo "⚠️  警告：存在未提交的更改"
    if [ "$DRY_RUN" = false ]; then
        read -p "是否继续？[y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
fi

echo "📦 准备为项目打 tag: $VERSION"
echo ""

# 获取当前 commit
CURRENT_COMMIT=$(git rev-parse HEAD)
SHORT_COMMIT=$(git rev-parse --short HEAD)

echo "📋 信息:"
echo "   Commit: $SHORT_COMMIT"
echo "   Version: $VERSION"
echo ""

if [ "$DRY_RUN" = true ]; then
    echo "🔍 预览将要执行的操作:"
    echo "   git tag $VERSION"
    echo "   git push origin $VERSION"
    echo ""
    echo "✅ 预览完成（干运行，未实际执行）"
    exit 0
fi

# 创建 tag
echo "🏷️  创建 tag..."
git tag "$VERSION"
echo "✅ tag $VERSION 创建成功"

# 推送 tag
echo "📤 推送 tag 到远程..."
git push origin "$VERSION"
echo "✅ tag $VERSION 已推送"

echo ""
echo "🎉 完成！"
echo "   Tag: $VERSION"
echo "   Commit: $SHORT_COMMIT"
echo ""
echo "验证命令:"
echo "   git tag -l | grep $VERSION"
echo "   git ls-remote --tags origin | grep $VERSION"
