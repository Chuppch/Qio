#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
import os
import re
import shlex
import shutil
import subprocess
import sys
import tempfile
from collections import defaultdict
from pathlib import Path
from typing import Any


HEADER_PATTERN = re.compile(r"^:[a-z0-9_+-]+:(?: \([^)]+\))?(?: !)? .+")
BREAKING_CHANGE_PATTERN = re.compile(r"^BREAKING CHANGE:[ \t].+$")
JIRA_REFS_PATTERN = re.compile(r"^Jira-Refs:[ \t].+$")
TRAILER_PATTERN = re.compile(r"^AI-Co-Authored-By:[ \t].+$", re.IGNORECASE)
FORBIDDEN_COAUTHOR_PATTERN = re.compile(r"^Co-authored-by:|^Co-Authored-By:")
JIRA_URL_PATTERN = re.compile(r"https?://jira\.meitu\.com/browse/([A-Z][A-Z0-9]+-\d+)")
JIRA_KEY_PATTERN = re.compile(r"\b([A-Z][A-Z0-9]+-\d+)\b")
HUNK_HEADER_PATTERN = re.compile(
    r"^@@ -(?P<old_start>\d+)(?:,(?P<old_count>\d+))? "
    r"\+(?P<new_start>\d+)(?:,(?P<new_count>\d+))? @@"
)


class BatchPlanError(Exception):
    """批次计划无法安全生成或应用时抛出的统一错误。"""


def run_git(
    repo_path: str | Path,
    args: list[str],
    *,
    input_text: str | None = None,
    env: dict[str, str] | None = None,
    allow_returncodes: tuple[int, ...] = (0,),
) -> str:
    """在指定仓库中执行 git 命令，并把失败统一转成 BatchPlanError。"""
    merged_env = os.environ.copy()
    if env:
        merged_env.update(env)

    result = subprocess.run(
        ["git", *args],
        cwd=repo_path,
        env=merged_env,
        input=input_text,
        text=True,
        capture_output=True,
    )

    if result.returncode not in allow_returncodes:
        message = (
            result.stderr.strip()
            or result.stdout.strip()
            or f"git {' '.join(args)} failed"
        )
        raise BatchPlanError(message)

    return result.stdout


def trim_path_label(label: str | None) -> str | None:
    """把 diff 头里的 a/、b/ 或 /dev/null 标记收敛成真实路径。"""
    if not label:
        return None
    if label == "/dev/null":
        return None
    if label.startswith(("a/", "b/")):
        return label[2:]
    return label


def parse_diff_git_paths(header_line: str) -> tuple[str | None, str | None]:
    """解析 `diff --git` 头，拿到旧路径和新路径。"""
    if not header_line.startswith("diff --git "):
        raise BatchPlanError(f"unsupported diff header: {header_line}")

    parts = shlex.split(header_line[len("diff --git ") :])
    if len(parts) < 2:
        raise BatchPlanError(f"unable to parse diff header: {header_line}")

    return trim_path_label(parts[0]), trim_path_label(parts[1])


def stable_id(kind: str, path: str, patch: str) -> str:
    """基于 kind、path 和 patch 内容生成稳定 unit id。"""
    payload = f"{kind}\0{path}\0".encode("utf-8") + patch.encode("utf-8")
    return f"{kind}-{hashlib.sha256(payload).hexdigest()[:12]}"


def detect_change_type(lines: list[str]) -> str:
    """从单个文件 patch 头部推断 Git 变更类型。"""
    stripped_lines = [line.rstrip("\n") for line in lines]
    if any(line.startswith("new file mode ") for line in stripped_lines):
        return "A"
    if any(line.startswith("deleted file mode ") for line in stripped_lines):
        return "D"
    if any(line.startswith("rename from ") for line in stripped_lines):
        return "R"
    return "M"


def parse_hunk_header(header_line: str) -> tuple[int, int, int, int]:
    """解析 hunk 头中的旧文件/新文件起始行与行数。"""
    match = HUNK_HEADER_PATTERN.match(header_line.strip())
    if not match:
        raise BatchPlanError(f"invalid hunk header: {header_line}")

    old_count = int(match.group("old_count") or "1")
    new_count = int(match.group("new_count") or "1")
    return (
        int(match.group("old_start")),
        old_count,
        int(match.group("new_start")),
        new_count,
    )


def split_patch_blocks(diff_text: str) -> list[str]:
    """把整段 diff 文本按文件切成独立 patch block。"""
    if not diff_text.strip():
        return []

    blocks: list[str] = []
    current: list[str] = []

    for line in diff_text.splitlines(keepends=True):
        # Git 会把多个文件 patch 串成一段输出，这里先切回单文件 block，
        # 后面的覆盖率校验和 unit 分配才能按文件精确判断。
        if line.startswith("diff --git "):
            if current:
                blocks.append("".join(current))
            current = [line]
            continue

        if current:
            current.append(line)

    if current:
        blocks.append("".join(current))

    return blocks


def build_untracked_patch_blocks(
    repo_path: str | Path,
    base_head: str,
    untracked_paths: list[str],
) -> list[str]:
    """在临时 index 中 stage 未跟踪路径，并导出稳定的新增 patch block。"""
    paths = sorted(filter(None, untracked_paths))
    if not paths:
        return []

    with tempfile.TemporaryDirectory(prefix="emoji-commit-untracked-index-") as temp_dir:
        temp_index = Path(temp_dir) / "index"
        env = {"GIT_INDEX_FILE": str(temp_index)}
        # `git diff HEAD` 无法包含 untracked 内容，这里用隔离 index 暂存它们，
        # 再通过 `git diff --cached <base>` 导出成标准新增 patch。
        run_git(repo_path, ["read-tree", base_head], env=env)
        run_git(repo_path, ["add", "--all", "--", *paths], env=env)
        diff_text = run_git(
            repo_path,
            [
                "diff",
                "--cached",
                "--binary",
                "--full-index",
                "--no-color",
                "--find-renames",
                "--no-ext-diff",
                base_head,
                "--",
                *paths,
            ],
            env=env,
        )

    return split_patch_blocks(diff_text)


def change_type_summary(
    change_type: str,
    *,
    binary: bool = False,
    old_path: str | None = None,
    new_path: str | None = None,
) -> str:
    """把 Git 变更类型转成更适合预览文本的摘要。"""
    if binary:
        return "binary patch"
    if change_type == "A":
        return "new file"
    if change_type == "D":
        return "delete file"
    if change_type == "R":
        return f"rename {old_path} -> {new_path}"
    return "file patch"


def parse_file_patch(block: str) -> tuple[dict[str, Any], list[dict[str, Any]]]:
    """把单文件 patch 解析成 file record 与可分配的 units。"""
    if not block.strip():
        raise BatchPlanError("encountered empty diff block")

    lines = block.splitlines(keepends=True)
    old_path, new_path = parse_diff_git_paths(lines[0].rstrip("\n"))
    change_type = detect_change_type(lines)
    binary = "GIT binary patch" in block or any(
        line.startswith("Binary files ") for line in lines
    )

    rename_from = None
    rename_to = None
    for line in lines:
        stripped = line.rstrip("\n")
        if stripped.startswith("rename from "):
            rename_from = stripped[len("rename from ") :]
        elif stripped.startswith("rename to "):
            rename_to = stripped[len("rename to ") :]

    if rename_from:
        old_path = rename_from
    if rename_to:
        new_path = rename_to

    path = new_path if change_type != "D" else old_path
    if not path:
        path = new_path or old_path
    if not path:
        raise BatchPlanError(f"unable to resolve patch path: {lines[0].rstrip()}")

    hunk_indices = [
        index for index, line in enumerate(lines) if line.startswith("@@ ")
    ]
    header_lines = lines[: hunk_indices[0]] if hunk_indices else lines[:]

    # 只有“普通文本修改 + 多个互不重叠的 hunk”才允许按 hunk 拆分。
    # 新增、删除、重命名、二进制 patch 都保持 file 级，避免生成 Git
    # 无法安全重放的半成品状态。
    partial_split_supported = (
        change_type == "M"
        and not binary
        and len(hunk_indices) > 1
        and old_path == new_path
    )

    units: list[dict[str, Any]] = []

    if binary or change_type in {"A", "D", "R"} or not hunk_indices:
        patch = block if block.endswith("\n") else f"{block}\n"
        units.append(
            {
                "id": stable_id("file", path, patch),
                "kind": "file",
                "path": path,
                "change_type": change_type,
                "patch": patch,
                "summary": change_type_summary(
                    change_type,
                    binary=binary,
                    old_path=old_path,
                    new_path=new_path,
                ),
                "partial_split_supported": False,
                "binary": binary,
            }
        )
    else:
        slice_boundaries = hunk_indices + [len(lines)]
        for start, end in zip(slice_boundaries, slice_boundaries[1:]):
            hunk_lines = lines[start:end]
            hunk_header = hunk_lines[0].rstrip("\n")
            old_start, old_count, new_start, new_count = parse_hunk_header(hunk_header)
            patch = "".join(header_lines + hunk_lines)
            units.append(
                {
                    "id": stable_id("hunk", path, patch),
                    "kind": "hunk",
                    "path": path,
                    "change_type": change_type,
                    "patch": patch,
                    "summary": hunk_header,
                    "partial_split_supported": partial_split_supported,
                    "old_start": old_start,
                    "old_count": old_count,
                    "new_start": new_start,
                    "new_count": new_count,
                }
            )

    file_record = {
        "path": path,
        "old_path": old_path or path,
        "new_path": new_path or path,
        "change_type": change_type,
        "binary": binary,
        "partial_split_supported": partial_split_supported,
        "unit_ids": [unit["id"] for unit in units],
    }

    return file_record, units


def build_inventory(
    repo_path: str | Path,
    base_ref: str,
    input_scope: str = "worktree",
) -> dict[str, Any]:
    """基于 base_ref 收集 staged 或整个 worktree 的变更清单。"""
    if input_scope not in {"staged", "worktree"}:
        raise BatchPlanError(f"unsupported input_scope: {input_scope}")

    base_head = run_git(repo_path, ["rev-parse", base_ref]).strip()

    if input_scope == "staged":
        diff_text = run_git(
            repo_path,
            [
                "diff",
                "--cached",
                "--binary",
                "--full-index",
                "--no-color",
                "--find-renames",
                "--no-ext-diff",
                base_ref,
                "--",
            ],
        )
    else:
        diff_text = run_git(
            repo_path,
            [
                "diff",
                "--binary",
                "--full-index",
                "--no-color",
                "--find-renames",
                "--no-ext-diff",
                base_ref,
                "--",
            ],
        )

    patch_blocks = split_patch_blocks(diff_text)

    if input_scope == "worktree":
        untracked = run_git(
            repo_path,
            ["ls-files", "--others", "--exclude-standard", "-z"],
        )
        patch_blocks.extend(
            build_untracked_patch_blocks(
                repo_path,
                base_head,
                untracked.split("\0"),
            )
        )

    files: list[dict[str, Any]] = []
    units: list[dict[str, Any]] = []
    for block in patch_blocks:
        file_record, file_units = parse_file_patch(block)
        files.append(file_record)
        units.extend(file_units)

    return {
        "base_ref": base_ref,
        "base_head": base_head,
        "input_scope": input_scope,
        "files": files,
        "units": units,
        "stats": {
            "file_count": len(files),
            "unit_count": len(units),
        },
    }


def load_json_file(path: str | Path) -> dict[str, Any]:
    """读取计划文件，并把常见文件/JSON 错误转成业务错误。"""
    try:
        with Path(path).open("r", encoding="utf-8") as handle:
            return json.load(handle)
    except FileNotFoundError as exc:
        raise BatchPlanError(f"plan file does not exist: {path}") from exc
    except json.JSONDecodeError as exc:
        raise BatchPlanError(f"invalid plan JSON: {exc}") from exc


def resolve_agent_name() -> str:
    """按既定优先级推断 AI-Co-Authored-By 的代理名称。"""
    for key in ("COMMIT_AI_AGENT_NAME", "AI_AGENT_NAME", "AGENT_NAME"):
        value = os.getenv(key)
        if value:
            return value

    prefix_map = (
        ("OPENAI_", "Codex"),
        ("CODEX_", "Codex"),
        ("ANTHROPIC_", "Claude"),
        ("CLAUDE_", "Claude"),
        ("MINMAX_", "MinMax"),
        ("GOOGLE_", "Gemini"),
        ("GEMINI_", "Gemini"),
    )
    for key in os.environ:
        for prefix, display_name in prefix_map:
            if key.startswith(prefix):
                return display_name

    return "AI Agent"


def sanitize_agent_name(agent_name: str) -> str:
    """清洗代理名称，避免 trailer 被换行或冒号污染。"""
    value = str(agent_name or "").strip()
    value = re.sub(r"[\r\n\x00-\x1f\x7f]+", " ", value)
    value = value.replace(":", "-")
    value = re.sub(r"\s+", " ", value).strip()
    return value or "AI Agent"


def normalize_body_item(item: Any) -> str:
    """把 body 条目标准化成无前缀、单行的 bullet 内容。"""
    text = str(item).strip()
    text = re.sub(r"^[-*]\s*", "", text)
    return text


def extract_jira_keys(value: Any) -> list[str]:
    """从 Jira URL / Jira Key 混合输入中提取、去重并保持首次出现顺序。"""
    if value is None:
        raw_items: list[Any] = []
    elif isinstance(value, str):
        raw_items = [value]
    else:
        try:
            raw_items = list(value)
        except TypeError:
            raw_items = [value]

    ordered_keys: list[str] = []
    seen: set[str] = set()

    for item in raw_items:
        text = str(item).strip()
        if not text:
            continue

        matches: list[tuple[int, str]] = []
        for match in JIRA_URL_PATTERN.finditer(text):
            matches.append((match.start(1), match.group(1)))
        for match in JIRA_KEY_PATTERN.finditer(text):
            matches.append((match.start(1), match.group(1)))

        for _, key in sorted(matches, key=lambda pair: pair[0]):
            if key not in seen:
                seen.add(key)
                ordered_keys.append(key)

    return ordered_keys


def normalize_jira_refs(value: Any) -> str:
    """把 Jira 输入归一成单行 Jira-Refs trailer；没有 key 时返回空字符串。"""
    jira_keys = extract_jira_keys(value)
    if not jira_keys:
        return ""
    return f"Jira-Refs: {', '.join(jira_keys)}"


def build_commit_message(message_data: dict[str, Any]) -> str:
    """把 header/body 组装成最终 commit message，并追加 AI trailer。"""
    if not isinstance(message_data, dict):
        raise BatchPlanError("message must be an object")

    header = str(message_data.get("header", "")).strip()
    if not header:
        raise BatchPlanError("message.header is required")

    body_value = message_data.get("body", [])
    if body_value is None:
        raw_body_items: list[Any] = []
    elif isinstance(body_value, str):
        raw_body_items = [body_value]
    else:
        try:
            raw_body_items = list(body_value)
        except TypeError:
            raw_body_items = [body_value]

    body_items = [
        normalize_body_item(item)
        for item in raw_body_items
        if str(item).strip()
    ]
    # 这是当前执行器的紧凑输出限制，不是 header/footer 语法本身的一部分。
    # 继续保留它，避免 batch preview 与最终消息在 CLI 中膨胀成过长正文。
    if len(body_items) > 5:
        raise BatchPlanError("commit body may include at most 5 items")

    breaking_change = str(message_data.get("breaking_change", "")).strip()
    if breaking_change:
        breaking_change = re.sub(r"[\r\n]+", " ", breaking_change)
        breaking_change = re.sub(r"\s+", " ", breaking_change).strip()
        if not breaking_change:
            raise BatchPlanError("message.breaking_change must not be empty")

    jira_refs = normalize_jira_refs(message_data.get("jira_refs", []))

    trailer = f"AI-Co-Authored-By: {sanitize_agent_name(resolve_agent_name())}"
    lines = [header, ""]
    if body_items:
        lines.extend(f"- {item}" for item in body_items)
        lines.append("")
    if jira_refs:
        lines.append(jira_refs)
        lines.append("")
    if breaking_change:
        lines.append(f"BREAKING CHANGE: {breaking_change}")
        lines.append("")
    lines.append(trailer)
    full_text = "\n".join(lines) + "\n"
    validate_commit_message_text(full_text)
    return full_text


def validate_commit_message_text(message: str) -> None:
    """校验 commit message 是否满足 emoji-commit 的格式约束。"""
    lines = message.replace("\r\n", "\n").replace("\r", "\n").split("\n")
    while lines and lines[-1] == "":
        lines.pop()

    if not lines:
        raise BatchPlanError("commit message is empty")

    header = lines[0]
    if not HEADER_PATTERN.match(header):
        raise BatchPlanError(f"invalid commit header: {header}")

    if len(lines) < 2 or lines[1] != "":
        raise BatchPlanError("commit message must contain a blank line after the header")

    trailer_indexes = [
        index for index, line in enumerate(lines) if TRAILER_PATTERN.match(line)
    ]
    if len(trailer_indexes) != 1:
        raise BatchPlanError(
            "commit message must contain exactly one AI-Co-Authored-By trailer"
        )

    breaking_indexes = [
        index for index, line in enumerate(lines) if BREAKING_CHANGE_PATTERN.match(line)
    ]
    if len(breaking_indexes) > 1:
        raise BatchPlanError(
            "commit message may contain at most one BREAKING CHANGE footer"
        )

    jira_refs_indexes = [
        index for index, line in enumerate(lines) if JIRA_REFS_PATTERN.match(line)
    ]
    if len(jira_refs_indexes) > 1:
        raise BatchPlanError(
            "commit message may contain at most one Jira-Refs footer"
        )

    if any(FORBIDDEN_COAUTHOR_PATTERN.match(line) for line in lines):
        raise BatchPlanError(
            "commit message must not contain Co-authored-by trailers"
        )

    trailer_index = trailer_indexes[0]
    jira_refs_index = jira_refs_indexes[0] if jira_refs_indexes else None
    breaking_index = breaking_indexes[0] if breaking_indexes else None

    if jira_refs_index is not None:
        if jira_refs_index == 0 or lines[jira_refs_index - 1] != "":
            raise BatchPlanError(
                "commit message must contain a blank line before the footer block"
            )
        if breaking_index is not None:
            if jira_refs_index >= breaking_index:
                raise BatchPlanError(
                    "Jira-Refs footer must appear before BREAKING CHANGE"
                )
            if jira_refs_index + 1 != breaking_index - 1 or lines[jira_refs_index + 1] != "":
                raise BatchPlanError(
                    "Jira-Refs footer must be separated from BREAKING CHANGE by a blank line"
                )
        elif jira_refs_index != trailer_index - 2 or lines[jira_refs_index + 1] != "":
            raise BatchPlanError(
                "Jira-Refs footer must be separated from the AI-Co-Authored-By trailer by a blank line"
            )

    if breaking_index is not None:
        if breaking_index != trailer_index - 2:
            raise BatchPlanError(
                "BREAKING CHANGE footer must be separated from the AI-Co-Authored-By trailer by a blank line"
            )
        if breaking_index == 0 or lines[breaking_index - 1] != "":
            raise BatchPlanError(
                "commit message must contain a blank line before the footer block"
            )
        if lines[breaking_index + 1] != "":
            raise BatchPlanError(
                "commit message must contain a blank line between BREAKING CHANGE and AI-Co-Authored-By"
            )
    elif trailer_index == 0 or lines[trailer_index - 1] != "":
        raise BatchPlanError(
            "commit message must contain a blank line before the trailer"
        )

    if trailer_index != len(lines) - 1:
        raise BatchPlanError("AI-Co-Authored-By trailer must be the last line")


def validate_plan(repo_path: str | Path, plan: dict[str, Any]) -> dict[str, Any]:
    """校验批次计划，并补齐 apply/preview 所需的规范化数据。"""
    commits = plan.get("commits")
    if not isinstance(commits, list) or not commits:
        raise BatchPlanError("plan must contain a non-empty commits array")

    current_head = run_git(repo_path, ["rev-parse", "HEAD"]).strip()
    plan_head = str(plan.get("base_head", "")).strip()
    if not plan_head:
        raise BatchPlanError("plan.base_head is required")
    if current_head != plan_head:
        raise BatchPlanError(
            f"HEAD changed since preview: expected {plan_head}, got {current_head}"
        )

    input_scope = str(plan.get("input_scope", "")).strip()
    if input_scope not in {"staged", "worktree"}:
        raise BatchPlanError("plan.input_scope must be either 'staged' or 'worktree'")

    inventory = build_inventory(repo_path, "HEAD", input_scope)
    units_by_id = {unit["id"]: unit for unit in inventory["units"]}
    files_by_path = {file_record["path"]: file_record for file_record in inventory["files"]}

    assigned: dict[str, str] = {}
    commit_ids: set[str] = set()
    normalized_commits: list[dict[str, Any]] = []

    for commit in commits:
        commit_id = str(commit.get("id", "")).strip()
        if not commit_id:
            raise BatchPlanError("each commit requires a stable id")
        if commit_id in commit_ids:
            raise BatchPlanError(f"duplicate commit id: {commit_id}")
        commit_ids.add(commit_id)

        split_mode = commit.get("split_mode")
        if split_mode not in {"file", "hunk"}:
            raise BatchPlanError(
                f"unsupported split_mode for {commit_id}: {split_mode}"
            )

        reason = str(commit.get("reason", "")).strip()
        if not reason:
            raise BatchPlanError(f"commit {commit_id} must include a reason")

        units = commit.get("units")
        if not isinstance(units, list) or not units:
            raise BatchPlanError(
                f"commit {commit_id} must include at least one unit"
            )

        normalized_units: list[str] = []
        for unit_id in units:
            normalized_unit_id = str(unit_id).strip()
            if normalized_unit_id not in units_by_id:
                raise BatchPlanError(
                    f"commit {commit_id} references unknown unit: {normalized_unit_id}"
                )
            if normalized_unit_id in assigned:
                raise BatchPlanError(
                    f"unit {normalized_unit_id} is assigned more than once"
                )
            assigned[normalized_unit_id] = commit_id
            normalized_units.append(normalized_unit_id)

        message_data = commit.get("message")
        if not isinstance(message_data, dict):
            raise BatchPlanError(
                f"commit {commit_id} must include message.header and message.body"
            )

        full_message = build_commit_message(message_data)
        normalized_body = []
        body_value = message_data.get("body", [])
        if body_value is None:
            body_value = []
        if isinstance(body_value, str):
            body_value = [body_value]
        for item in body_value:
            if str(item).strip():
                normalized_body.append(normalize_body_item(item))

        normalized_commits.append(
            {
                "id": commit_id,
                "reason": reason,
                "split_mode": split_mode,
                "units": normalized_units,
                "message": {
                    "header": str(message_data.get("header", "")).strip(),
                    "body": normalized_body,
                    "jira_refs": extract_jira_keys(message_data.get("jira_refs", [])),
                    "full_text": full_message,
                },
            }
        )

    inventory_unit_ids = set(units_by_id)
    assigned_ids = set(assigned)
    missing_units = sorted(inventory_unit_ids - assigned_ids)
    extra_units = sorted(assigned_ids - inventory_unit_ids)
    if missing_units or extra_units:
        details = []
        if missing_units:
            details.append(f"missing units: {', '.join(missing_units)}")
        if extra_units:
            details.append(f"unexpected units: {', '.join(extra_units)}")
        raise BatchPlanError("; ".join(details))

    for commit in normalized_commits:
        if commit["split_mode"] == "file":
            units_by_path: dict[str, set[str]] = defaultdict(set)
            for unit_id in commit["units"]:
                units_by_path[units_by_id[unit_id]["path"]].add(unit_id)
            for path, commit_unit_ids in units_by_path.items():
                file_record = files_by_path[path]
                file_unit_ids = set(file_record["unit_ids"])
                # file 模式是严格全量覆盖：既然声称“这个 commit 接管整个文件”，
                # 那就不能留下同文件的其他 unit 给后续 commit 捡漏。
                if commit_unit_ids != file_unit_ids:
                    raise BatchPlanError(
                        f"commit {commit['id']} uses split_mode=file but does not cover full file: {path}"
                    )
        else:
            for unit_id in commit["units"]:
                unit = units_by_id[unit_id]
                file_record = files_by_path[unit["path"]]
                # hunk 模式只允许作用在 inventory 阶段已经标记为“可安全局部拆分”
                # 的文本修改上，避免把危险 patch 强拆开。
                if unit["kind"] != "hunk" or not file_record["partial_split_supported"]:
                    raise BatchPlanError(
                        f"commit {commit['id']} uses split_mode=hunk with unsupported unit: {unit_id}"
                    )

    return {
        "input_scope": input_scope,
        "inventory": inventory,
        "units_by_id": units_by_id,
        "files_by_path": files_by_path,
        "commits": normalized_commits,
    }


def build_preview_text(validation: dict[str, Any]) -> str:
    """把校验后的计划渲染成给用户确认的纯文本预览。"""
    inventory = validation["inventory"]
    units_by_id = validation["units_by_id"]
    files_by_path = validation["files_by_path"]

    lines = [
        "emoji-commit batch preview",
        f"Base HEAD: {inventory['base_head']}",
        f"Input scope: {validation['input_scope']}",
        "Status: waiting for confirmation, no commits have been created.",
        "",
    ]

    for index, commit in enumerate(validation["commits"], start=1):
        lines.append(f"{index}. {commit['message']['header']}")
        lines.append(f"Reason: {commit['reason']}")
        lines.append(f"Split mode: {commit['split_mode']}")
        has_partial_split = any(
            units_by_id[unit_id]["partial_split_supported"]
            for unit_id in commit["units"]
        )
        lines.append(f"Partial split: {'yes' if has_partial_split else 'no'}")
        lines.append("Coverage:")
        for unit_id in commit["units"]:
            unit = units_by_id[unit_id]
            file_record = files_by_path[unit["path"]]
            if unit["kind"] == "hunk":
                lines.append(f"- {unit['path']} {unit['summary']}")
                continue

            if file_record["change_type"] == "R":
                label = f"{file_record['old_path']} -> {file_record['new_path']}"
            else:
                label = unit["path"]
            lines.append(f"- {label} ({file_record['change_type']})")

        lines.append("Body:")
        if commit["message"]["body"]:
            lines.extend(f"- {item}" for item in commit["message"]["body"])
        else:
            lines.append("- (no body items)")
        lines.append("")

    return "\n".join(lines) + "\n"


def create_commits_with_temp_index(
    repo_path: str | Path, validation: dict[str, Any]
) -> tuple[str, str, list[dict[str, str]]]:
    """在临时 index 中构造整组 commit，但先不真正移动 HEAD。"""
    start_head = validation["inventory"]["base_head"]
    final_commit = start_head
    created_commits: list[dict[str, str]] = []

    with tempfile.TemporaryDirectory(prefix="emoji-commit-index-") as temp_dir:
        temp_index = Path(temp_dir) / "index"
        env = {"GIT_INDEX_FILE": str(temp_index)}
        # 在隔离的 index 里组装整组 commit，这样 preview/apply 可以先把树拼好，
        # 不会提前污染用户当前真实 index。
        run_git(repo_path, ["read-tree", start_head], env=env)

        for commit in validation["commits"]:
            for unit_id in commit["units"]:
                patch = validation["units_by_id"][unit_id]["patch"]
                run_git(
                    repo_path,
                    ["apply", "--cached", "--binary", "--whitespace=nowarn", "-"],
                    input_text=patch,
                    env=env,
                )

            tree = run_git(repo_path, ["write-tree"], env=env).strip()
            final_commit = run_git(
                repo_path,
                ["commit-tree", tree, "-p", final_commit],
                input_text=commit["message"]["full_text"],
            ).strip()
            stored_message = run_git(
                repo_path,
                ["log", "-1", "--pretty=%B", final_commit],
            )
            validate_commit_message_text(stored_message)
            created_commits.append(
                {
                    "id": commit["id"],
                    "commit": final_commit,
                    "header": commit["message"]["header"],
                }
            )

    return start_head, final_commit, created_commits


def restore_index(index_path: str | Path, backup_path: Path | None) -> None:
    """把真实 index 恢复到 apply 前的状态。"""
    if backup_path and backup_path.exists():
        shutil.copyfile(backup_path, index_path)
        return

    if Path(index_path).exists():
        Path(index_path).unlink()


def apply_commits_transaction(
    repo_path: str | Path, start_head: str, final_commit: str
) -> None:
    """以事务方式把临时构造好的 commit 链正式落到当前仓库。"""
    raw_index_path = Path(
        run_git(repo_path, ["rev-parse", "--git-path", "index"]).strip()
    )
    index_path = raw_index_path if raw_index_path.is_absolute() else Path(repo_path) / raw_index_path
    backup_path: Path | None = None

    if index_path.exists():
        with tempfile.NamedTemporaryFile(
            prefix="emoji-commit-index-backup-",
            delete=False,
        ) as handle:
            backup_path = Path(handle.name)
        shutil.copyfile(index_path, backup_path)

    updated_ref = False
    try:
        # 只有整组 commit 都构造完成后才真正移动 HEAD。
        # 如果后续 read-tree 刷新工作区快照失败，就把 HEAD 和磁盘 index 一起回滚。
        run_git(
            repo_path,
            ["update-ref", "-m", "emoji-commit batch apply", "HEAD", final_commit, start_head],
        )
        updated_ref = True
        run_git(repo_path, ["read-tree", final_commit])
    except Exception:
        if updated_ref:
            try:
                run_git(
                    repo_path,
                    [
                        "update-ref",
                        "-m",
                        "emoji-commit batch rollback",
                        "HEAD",
                        start_head,
                        final_commit,
                    ],
                )
            except Exception:
                pass
        restore_index(index_path, backup_path)
        raise
    finally:
        if backup_path and backup_path.exists():
            backup_path.unlink()


def cmd_inventory(args: argparse.Namespace) -> None:
    """CLI 子命令：输出当前 inventory JSON。"""
    inventory = build_inventory(args.repo, args.base, args.scope)
    print(json.dumps(inventory, ensure_ascii=False, indent=2))


def cmd_preview_plan(args: argparse.Namespace) -> None:
    """CLI 子命令：校验计划并输出预览文本。"""
    plan = load_json_file(args.plan)
    validation = validate_plan(args.repo, plan)
    sys.stdout.write(build_preview_text(validation))


def cmd_apply_plan(args: argparse.Namespace) -> None:
    """CLI 子命令：校验计划并事务式应用整组 commit。"""
    plan = load_json_file(args.plan)
    validation = validate_plan(args.repo, plan)
    start_head, final_commit, created_commits = create_commits_with_temp_index(
        args.repo,
        validation,
    )
    apply_commits_transaction(args.repo, start_head, final_commit)
    print(
        json.dumps(
            {
                "base_head": start_head,
                "final_head": final_commit,
                "created_commits": created_commits,
            },
            ensure_ascii=False,
            indent=2,
        )
    )


def build_parser() -> argparse.ArgumentParser:
    """构建 `inventory / preview-plan / apply-plan` 三个 CLI 入口。"""
    def add_repo_argument(target: argparse.ArgumentParser) -> None:
        target.add_argument(
            "--repo",
            default=argparse.SUPPRESS,
            help="Git repository path",
        )

    parser = argparse.ArgumentParser(
        description="Split staged or worktree changes into previewable emoji-commit batches."
    )
    add_repo_argument(parser)
    subparsers = parser.add_subparsers(dest="command", required=True)

    inventory_parser = subparsers.add_parser(
        "inventory",
        help="Output the current change inventory as JSON",
    )
    add_repo_argument(inventory_parser)
    inventory_parser.add_argument(
        "--base",
        default="HEAD",
        help="Base ref for inventory generation",
    )
    inventory_parser.add_argument(
        "--scope",
        choices=("worktree", "staged"),
        default="worktree",
        help="Inventory scope: full worktree or staged-only",
    )
    inventory_parser.set_defaults(func=cmd_inventory)

    preview_parser = subparsers.add_parser(
        "preview-plan",
        help="Render a human-readable preview for a batch plan",
    )
    add_repo_argument(preview_parser)
    preview_parser.add_argument("--plan", required=True, help="Path to the plan JSON file")
    preview_parser.set_defaults(func=cmd_preview_plan)

    apply_parser = subparsers.add_parser(
        "apply-plan",
        help="Apply a confirmed batch plan transactionally",
    )
    add_repo_argument(apply_parser)
    apply_parser.add_argument("--plan", required=True, help="Path to the plan JSON file")
    apply_parser.set_defaults(func=cmd_apply_plan)

    return parser


def main(argv: list[str] | None = None) -> int:
    """命令行主入口，统一处理 BatchPlanError 的退出码。"""
    parser = build_parser()
    args = parser.parse_args(argv)
    if not hasattr(args, "repo"):
        args.repo = "."
    try:
        args.func(args)
    except BatchPlanError as exc:
        print(str(exc), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
