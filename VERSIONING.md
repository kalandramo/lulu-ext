# 版本管理规范

## 概述

本项目采用 **统一版本策略**，所有 53 个模块共享同一个版本号。

## 版本号格式

采用 [语义化版本](https://semver.org/lang/zh-CN/) (Semantic Versioning) 格式：

```
v<主版本>.<次版本>.<修订号>
例如：v1.0.0, v2.3.1, v0.1.0
```

### 版本递增规则

| 类型 | 触发条件 | 示例 |
|------|----------|------|
| **主版本 (major)** | 不兼容的 API 变更 | v1.0.0 → v2.0.0 |
| **次版本 (minor)** | 向下兼容的新功能 | v1.0.0 → v1.1.0 |
| **修订号 (patch)** | 向下兼容的问题修复 | v1.0.0 → v1.0.1 |

## 打 Tag 方式

### 方式一：使用脚本（本地）

```bash
# 1. 干运行（预览）
./scripts/tag-all.sh v1.0.0 --dry-run

# 2. 正式打 tag（自动更新子模块依赖）
./scripts/tag-all.sh v1.0.0 --update-deps

# 3. 自动计算下一个版本并打 tag
./scripts/tag-all.sh --update-deps  # 自动 bump patch 版本

# 4. 手动计算版本号
./scripts/bump-version.sh patch  # 输出：v1.0.1
./scripts/bump-version.sh minor  # 输出：v1.1.0
./scripts/bump-version.sh major  # 输出：v2.0.0
```

### 方式二：使用 GitHub Actions（推荐）

1. 进入 Actions 标签页
2. 选择 "Tag All Modules" workflow
3. 点击 "Run workflow"
4. 输入版本号或选择递增类型
5. workflow 会自动：
   - 创建 tag
   - 推送 tag 到远程
   - 更新所有子模块间的依赖版本
   - 提交并推送依赖更新

## 验证 Tag

```bash
# 查看本地 tag
git tag -l

# 查看远程 tag
git ls-remote --tags origin

# 查看特定 tag
git tag -l | grep v1.0.0
```

## 注意事项

1. **Tag 命名**: 必须以 `v` 开头，如 `v1.0.0`
2. **唯一性**: tag 不能重复创建
3. **权限**: 推送 tag 需要 write 权限
4. **不可变**: tag 创建后不应修改或删除

## 回滚 Tag（谨慎使用）

```bash
# 删除本地 tag
git tag -d v1.0.0

# 删除远程 tag
git push origin --delete v1.0.0
```

> ⚠️ 警告：删除已发布的 tag 可能导致依赖此版本的用户出现问题。
