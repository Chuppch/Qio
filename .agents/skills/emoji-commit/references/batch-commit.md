# Batch Commit 详细流程

当用户希望“先看整个未提交工作树，再按逻辑拆成多条 commit”，或当前仍存在 unstaged / untracked 改动时读取本文。

## 适用场景

以下表达默认进入 batch commit，而不是只盯着 staged diff：

- “看看 git 里没提交的代码，分类一下，然后分批提交代码”
- “把这些改动拆成几次 commit”
- “按逻辑分批提交”
- `split commits`
- `batch commit`

默认输入范围是整个未提交工作树（`worktree`），包含：

- 已暂存变更
- 未暂存变更
- 未跟踪文件

## 0. 先确认输入范围

```bash
git status --short
```

- 只要存在 unstaged 或 untracked 改动，默认用 `--scope worktree`。
- 只有工作区除 staged 外已经干净，或用户明确要求“只提交暂存区”，才用 `--scope staged`。
- 若遇到未跟踪 symlink / 目录被旧版 inventory 漏收，修复动作见 `troubleshooting.md`。

## 1. 收集 inventory

使用 `<skill_root>/scripts/commit_batches.py` 输出当前工作树 inventory：

```bash
python3 <skill_root>/scripts/commit_batches.py \
  --repo <repo-path> \
  inventory \
  --scope worktree > /tmp/emoji_commit_inventory.json
```

输出 JSON 至少包含：

- `base_ref`
- `base_head`
- `input_scope`
- `files`
- `units`
- `stats`

## 2. 生成批次计划

根据 inventory 组织计划 JSON，固定结构如下：

```json
{
  "base_head": "<HEAD hash from inventory>",
  "input_scope": "worktree",
  "commits": [
    {
      "id": "commit-1",
      "reason": "why this batch exists",
      "split_mode": "file",
      "units": ["file-1234567890ab"],
      "message": {
        "header": ":wrench: (scope) subject",
        "body": ["bullet 1", "bullet 2"]
      }
    }
  ]
}
```

规则：

- `input_scope` 在 v1 默认写 `worktree`。
- 每个 `unit` 必须恰好归属一个 commit，禁止遗漏和重复。
- `split_mode=file` 表示该 commit 必须完整覆盖某个文件的全部 unit。
- `split_mode=hunk` 只允许用于 `partial_split_supported=true` 的文本 patch。
- `message.header`、`message.body`、`message.jira_refs` 继续服从主文档和 `fex-conventional-commits.md` 的约束。
- 计划文件建议写到仓库外部临时路径，例如 `/tmp/emoji_commit_plan.json`，避免把 plan 本身当成新的未跟踪改动。

## 3. 预览计划

不要直接 apply 多条 commit，先渲染预览：

```bash
python3 <skill_root>/scripts/commit_batches.py \
  --repo <repo-path> \
  preview-plan \
  --plan /tmp/emoji_commit_plan.json
```

预览时至少确认：

- 拆分边界是否合理
- 每条 commit 的 header / body 是否符合语义
- 是否存在本应合并或本应拆开的改动

## 4. 等待用户确认

batch commit 的默认边界是：

- 先展示 preview
- 等用户确认
- 再执行 apply

v1 默认不自动创建多条 commit。

## 5. 确认后 apply

```bash
python3 <skill_root>/scripts/commit_batches.py \
  --repo <repo-path> \
  apply-plan \
  --plan /tmp/emoji_commit_plan.json
```

`apply-plan` 的行为约束：

- 先校验 `base_head` 是否仍与 preview 时一致
- 再校验当前 inventory 与计划中的 unit 分配是否仍匹配
- 使用临时 index 与 `commit-tree` 事务式生成 commits
- 成功后一次性更新 `HEAD` 与 index
- 失败时回滚 `HEAD` 与 index，不保留半成品提交

## 6. 逐条提交后校验

apply 完成后，仍需对最终提交执行消息校验：

- Header shortcode 合规
- `AI-Co-Authored-By:` 恰好一行
- Header / Body / footer 的空行位置正确
- `Jira-Refs:` 若存在，必须保持单行聚合格式与正确顺序

如需具体自检命令，回看 `single-commit.md`；如遇 preview / apply 失败、范围判断错误或 inventory 边缘问题，转到 `troubleshooting.md`。
