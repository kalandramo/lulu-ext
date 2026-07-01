#!/bin/bash
# 版本号计算工具
# 用法：./scripts/bump-version.sh [major|minor|patch] [current_version]

set -e

BUMP_TYPE=${1:-patch}
CURRENT_VERSION=${2:-$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")}

# 移除 'v' 前缀
CURRENT_VERSION=${CURRENT_VERSION#v}

# 解析版本号
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# 根据类型 bump
case "$BUMP_TYPE" in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
    *)
        echo "❌ 错误：未知的 bump 类型 '$BUMP_TYPE'"
        echo "用法：$0 [major|minor|patch] [current_version]"
        exit 1
        ;;
esac

NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"

echo "$NEW_VERSION"
