---
name: emoji-commit
description: |
  自动分析 git staged 变更，或将整个未提交工作树拆分为多条使用 cz-emoji shortcode 风格的 commit。
  当用户要求：commit code, generate commit message, auto commit, 帮我提交, 自动提交, 生成 commit, write commit message, 代码提交, 提交代码, split commits, batch commit，或表达“看看 git 里没提交的代码并分类后分批提交”时使用此技能。
---

# Emoji Commit

分析 git staged 变更，或将整个未提交工作树拆分为可预览、可确认、可事务式应用的 commit 批次。

## 安装

```bash
npx -y @meitu/skills add git@git.meitu.com:fex/internal-skills.git --skill emoji-commit
```

**提示：前置环境配置**

使用 `npx @meitu/skills` 脚本，需确保全局 `.npmrc` 包含以下配置：

```bash
# 一键配置脚本
npmrc="$HOME/.npmrc"; block=$'@meitu:registry=http://npm.meitu-int.com\nregistry=https://registry.npmmirror.com'; grep -Fqx "$block" "$npmrc" 2>/dev/null || printf '%s\n' "$block" >> "$npmrc"
```

详细配置说明参考 [`.npmrc` 配置说明](https://cf.meitu.com/confluence/x/XTuOEg)。

## 默认路由

先看完整工作区：

```bash
git status --short
```

- **Single commit**：工作区除 staged 外已经干净，或用户明确要求“只提交暂存区”。
- **Batch commit**：用户要求 `split commits` / `batch commit` / “按逻辑分批提交”，或当前仍存在 unstaged / untracked 改动。
- **Batch 安全边界**：先 preview，再等待用户确认，最后才 apply。

如果当前请求已经明显落在某条分支，直接读取对应 guide；否则先根据上面的路由判断。

## 默认结果约束

完整标题/页脚规范、顺序关系和示例以 [fex-conventional-commits.md](references/fex-conventional-commits.md) 为准。主 `SKILL.md` 默认只保留最小硬结果约束：

### Header

- 允许的标题骨架：
  - `:emoji: subject`
  - `:emoji: ! subject`
  - `:emoji: (scope) subject`
  - `:emoji: (scope) ! subject`
- Header 中的 emoji 必须使用 shortcode，例如 `:bug:`，不要使用 Unicode emoji。
- `scope` 是可选项；出现时必须写成 `(scope)`。
- `!` 是独立的 breaking marker，位于 emoji 或 `(scope)` 之后、subject 之前。

### Body

- Body 保持 bullet 风格，并尽量简洁。
- Header 和正文之间必须保留一个空行。
- 若存在 footer block，正文和 footer block 之间必须保留一个空行。

### Footer

- 若存在 Jira 上下文，输出单行 `Jira-Refs: KEY1, KEY2`，不要保留原始 Jira URL。
- 若存在 `BREAKING CHANGE:`，它必须位于 `AI-Co-Authored-By:` 之前。
- 最终消息必须包含且仅包含一行 `AI-Co-Authored-By:`，并且它必须是最后一行。
- 禁止输出 `Co-authored-by` / `Co-Authored-By` 之类的标准 co-author trailer。

## 何时读取哪份 Reference

| Reference | 何时读取 | 主要内容 |
|---|---|---|
| [single-commit.md](references/single-commit.md) | 已确定当前是 single commit 时 | staged 分析、提交命令、AI trailer、自检 |
| [batch-commit.md](references/batch-commit.md) | 已确定当前是 batch commit 时 | inventory、plan、preview、apply |
| [fex-conventional-commits.md](references/fex-conventional-commits.md) | 需要 canonical 规范时 | header、footer、Jira-Refs、完整示例 |
| [cz-emoji-types.md](references/cz-emoji-types.md) | 常用类型不足以覆盖当前语义时 | 完整 emoji 类型表 |
| [troubleshooting.md](references/troubleshooting.md) | 遇到异常或恢复场景时 | 自检失败、hook、fallback、恢复动作 |

## 快速提示

- Batch commit 使用 `<skill_root>/scripts/commit_batches.py` 作为执行器。
- [fex-conventional-commits.md](references/fex-conventional-commits.md) 描述的是**规范结果**；命令写法、风格建议和故障处理放在按需 guide 中。
- 如果你只需要快速选 emoji，可先看 [cz-emoji-types.md](references/cz-emoji-types.md)；如果你需要 canonical 顺序或 Jira footer 示例，优先回到 [fex-conventional-commits.md](references/fex-conventional-commits.md)。
