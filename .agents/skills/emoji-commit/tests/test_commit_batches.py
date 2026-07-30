import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock


SCRIPT_DIR = Path(__file__).resolve().parents[1] / "scripts"
if str(SCRIPT_DIR) not in sys.path:
    sys.path.insert(0, str(SCRIPT_DIR))

import commit_batches  # noqa: E402


class CommitBatchesTests(unittest.TestCase):
    maxDiff = None

    def setUp(self):
        self.tempdir = tempfile.TemporaryDirectory()
        self.repo = Path(self.tempdir.name)
        self.git("init", "-q")
        self.git("config", "user.name", "Test User")
        self.git("config", "user.email", "test@example.com")

    def tearDown(self):
        self.tempdir.cleanup()

    def git(self, *args, input_text=None, check=True):
        result = subprocess.run(
            ["git", *args],
            cwd=self.repo,
            text=True,
            input=input_text,
            capture_output=True,
        )
        if check and result.returncode != 0:
            self.fail(
                f"git {' '.join(args)} failed\n"
                f"stdout:\n{result.stdout}\n"
                f"stderr:\n{result.stderr}"
            )
        return result

    def run_cli(self, *args, check=True):
        result = subprocess.run(
            [sys.executable, str(SCRIPT_DIR / "commit_batches.py"), *args],
            cwd=self.repo,
            text=True,
            capture_output=True,
        )
        if check and result.returncode != 0:
            self.fail(
                f"cli {' '.join(args)} failed\n"
                f"stdout:\n{result.stdout}\n"
                f"stderr:\n{result.stderr}"
            )
        return result

    def write_text(self, relative_path, content):
        target = self.repo / relative_path
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(content, encoding="utf-8")

    def write_bytes(self, relative_path, content):
        target = self.repo / relative_path
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_bytes(content)

    def commit_all(self, message="init"):
        self.git("add", "-A")
        self.git("commit", "-m", message)

    def build_plan(self, inventory, commits):
        return {
            "base_head": inventory["base_head"],
            "input_scope": inventory["input_scope"],
            "commits": commits,
        }

    def test_inventory_includes_worktree_change_types_and_untracked(self):
        self.write_text("modify.txt", "base\n")
        self.write_text("delete.txt", "delete me\n")
        self.write_text("rename.txt", "rename me\n")
        self.write_bytes("binary.bin", b"\x00\x01\x02base")
        self.commit_all()

        self.write_text("modify.txt", "base changed\n")
        (self.repo / "delete.txt").unlink()
        self.git("mv", "rename.txt", "renamed.txt")
        self.write_bytes("binary.bin", b"\x00\x01\x02changed")
        self.git("add", "binary.bin", "renamed.txt")
        self.write_text("added.txt", "new file\n")

        inventory = commit_batches.build_inventory(self.repo, "HEAD", "worktree")
        files = {item["path"]: item for item in inventory["files"]}

        self.assertEqual(inventory["input_scope"], "worktree")
        self.assertEqual(inventory["stats"]["file_count"], 5)
        self.assertEqual(files["modify.txt"]["change_type"], "M")
        self.assertEqual(files["delete.txt"]["change_type"], "D")
        self.assertEqual(files["renamed.txt"]["change_type"], "R")
        self.assertEqual(files["renamed.txt"]["old_path"], "rename.txt")
        self.assertEqual(files["renamed.txt"]["new_path"], "renamed.txt")
        self.assertEqual(files["added.txt"]["change_type"], "A")
        self.assertTrue(files["binary.bin"]["binary"])

    def test_inventory_worktree_includes_staged_file_untracked_file_and_symlink(self):
        self.write_text("tracked.txt", "base\n")
        self.write_text("tracked-dir/seed.txt", "seed\n")
        self.commit_all()

        self.write_text("tracked.txt", "base changed\n")
        self.git("add", "tracked.txt")
        self.write_text("loose.txt", "new file\n")
        (self.repo / "tracked-dir-link").symlink_to("tracked-dir", target_is_directory=True)

        inventory = commit_batches.build_inventory(self.repo, "HEAD", "worktree")
        files = {item["path"]: item for item in inventory["files"]}
        symlink_unit = next(item for item in inventory["units"] if item["path"] == "tracked-dir-link")

        self.assertEqual(inventory["stats"]["file_count"], 3)
        self.assertEqual(set(files), {"tracked.txt", "loose.txt", "tracked-dir-link"})
        self.assertEqual(files["tracked.txt"]["change_type"], "M")
        self.assertEqual(files["loose.txt"]["change_type"], "A")
        self.assertEqual(files["tracked-dir-link"]["change_type"], "A")
        self.assertEqual(symlink_unit["kind"], "file")

    def test_non_overlapping_hunks_can_be_split_and_previewed(self):
        self.write_text("app.txt", "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n")
        self.commit_all()

        self.write_text("app.txt", "1x\n2\n3\n4\n5\n6\n7\n8\n9\n10x\n")
        inventory = commit_batches.build_inventory(self.repo, "HEAD", "worktree")
        app_file = inventory["files"][0]

        self.assertTrue(app_file["partial_split_supported"])
        self.assertEqual(len(app_file["unit_ids"]), 2)

        plan = self.build_plan(
            inventory,
            [
                {
                    "id": "first",
                    "reason": "split the first hunk",
                    "split_mode": "hunk",
                    "units": [app_file["unit_ids"][0]],
                    "message": {
                        "header": ":wrench: (emoji-commit) capture the first hunk",
                        "body": ["take the first isolated edit"],
                    },
                },
                {
                    "id": "second",
                    "reason": "split the second hunk",
                    "split_mode": "hunk",
                    "units": [app_file["unit_ids"][1]],
                    "message": {
                        "header": ":wrench: (emoji-commit) capture the second hunk",
                        "body": ["take the second isolated edit"],
                    },
                },
            ],
        )

        validation = commit_batches.validate_plan(self.repo, plan)
        preview = commit_batches.build_preview_text(validation)

        self.assertIn("emoji-commit batch preview", preview)
        self.assertIn("Input scope: worktree", preview)
        self.assertIn("Partial split: yes", preview)
        self.assertIn("capture the first hunk", preview)
        self.assertIn("app.txt @@ -1,4 +1,4 @@", preview)
        self.assertIn("@@ -1,4 +1,4 @@", preview)

    def test_build_commit_message_normalizes_body_and_trailer(self):
        message = commit_batches.build_commit_message(
            {
                "header": ":wrench: (emoji-commit) normalize batched commit text",
                "body": ["  first item  ", "- second item"],
            }
        )

        self.assertIn("- first item\n- second item\n", message)
        self.assertEqual(message.count("AI-Co-Authored-By:"), 1)
        commit_batches.validate_commit_message_text(message)

    def test_build_commit_message_supports_no_scope_and_breaking_footer(self):
        message = commit_batches.build_commit_message(
            {
                "header": ":bug: ! reject duplicate units",
                "body": ["reject duplicate assignment earlier"],
                "breaking_change": "duplicate unit assignments are now rejected earlier",
            }
        )

        self.assertIn(":bug: ! reject duplicate units\n", message)
        self.assertIn(
            "BREAKING CHANGE: duplicate unit assignments are now rejected earlier\n\n",
            message,
        )
        self.assertRegex(message.rstrip().splitlines()[-1], r"^AI-Co-Authored-By: .+$")
        commit_batches.validate_commit_message_text(message)

    def test_build_commit_message_supports_jira_refs_without_breaking_footer(self):
        message = commit_batches.build_commit_message(
            {
                "header": ":memo: (emoji-commit) document jira refs footer behavior",
                "body": ["explain footer ordering with jira refs"],
                "jira_refs": [
                    "https://jira.meitu.com/browse/INTERNAL-1901",
                    "DATA-6755",
                ],
            }
        )

        self.assertIn(
            "Jira-Refs: INTERNAL-1901, DATA-6755\n\nAI-Co-Authored-By:",
            message,
        )
        commit_batches.validate_commit_message_text(message)

    def test_build_commit_message_supports_jira_refs_with_breaking_footer(self):
        message = commit_batches.build_commit_message(
            {
                "header": ":sparkles: (emoji-commit) ! add jira refs footer support",
                "body": ["support jira refs in commit footer"],
                "jira_refs": ["DATA-6755", "TECHPUB-19087"],
                "breaking_change": "footer parsing now recognizes jira refs",
            }
        )

        self.assertIn(
            "Jira-Refs: DATA-6755, TECHPUB-19087\n\n"
            "BREAKING CHANGE: footer parsing now recognizes jira refs\n\n",
            message,
        )
        commit_batches.validate_commit_message_text(message)

    def test_build_commit_message_rejects_more_than_five_body_items(self):
        with self.assertRaisesRegex(
            commit_batches.BatchPlanError,
            "commit body may include at most 5 items",
        ):
            commit_batches.build_commit_message(
                {
                    "header": ":memo: (emoji-commit) keep body output compact",
                    "body": [
                        "item 1",
                        "item 2",
                        "item 3",
                        "item 4",
                        "item 5",
                        "item 6",
                    ],
                }
            )

    def test_validate_commit_message_accepts_optional_scope(self):
        commit_batches.validate_commit_message_text(
            ":memo: document breaking footer behavior\n\nAI-Co-Authored-By: AI Agent\n"
        )
        commit_batches.validate_commit_message_text(
            ":sparkles: (emoji-commit) ! change header grammar\n\n"
            "BREAKING CHANGE: commit headers no longer require a scope\n\n"
            "AI-Co-Authored-By: AI Agent\n"
        )
        commit_batches.validate_commit_message_text(
            ":memo: (emoji-commit) document jira refs footer behavior\n\n"
            "Jira-Refs: INTERNAL-1901, DATA-6755\n\n"
            "AI-Co-Authored-By: AI Agent\n"
        )
        commit_batches.validate_commit_message_text(
            ":sparkles: (emoji-commit) ! add jira refs footer support\n\n"
            "Jira-Refs: DATA-6755, TECHPUB-19087\n\n"
            "BREAKING CHANGE: footer parsing now recognizes jira refs\n\n"
            "AI-Co-Authored-By: AI Agent\n"
        )

    def test_validate_commit_message_rejects_breaking_footer_out_of_order(self):
        message = (
            ":bug: ! reject duplicate units\n\n"
            "BREAKING CHANGE: duplicate unit assignments are now rejected earlier\n"
            "AI-Co-Authored-By: AI Agent\n"
        )

        with self.assertRaisesRegex(
            commit_batches.BatchPlanError,
            "BREAKING CHANGE footer must be separated from the AI-Co-Authored-By trailer by a blank line",
        ):
            commit_batches.validate_commit_message_text(message)

    def test_validate_commit_message_rejects_missing_blank_line_after_breaking_footer(self):
        message = (
            ":bug: ! reject duplicate units\n\n"
            "BREAKING CHANGE: duplicate unit assignments are now rejected earlier\n"
            "AI-Co-Authored-By: AI Agent\n"
        )

        with self.assertRaisesRegex(
            commit_batches.BatchPlanError,
            "BREAKING CHANGE footer must be separated from the AI-Co-Authored-By trailer by a blank line",
        ):
            commit_batches.validate_commit_message_text(message)

    def test_validate_commit_message_rejects_jira_refs_without_blank_line_before_breaking(self):
        message = (
            ":sparkles: (emoji-commit) ! add jira refs footer support\n\n"
            "Jira-Refs: DATA-6755, TECHPUB-19087\n"
            "BREAKING CHANGE: footer parsing now recognizes jira refs\n\n"
            "AI-Co-Authored-By: AI Agent\n"
        )

        with self.assertRaisesRegex(
            commit_batches.BatchPlanError,
            "Jira-Refs footer must be separated from BREAKING CHANGE by a blank line",
        ):
            commit_batches.validate_commit_message_text(message)

    def test_extract_jira_keys_deduplicates_and_preserves_order(self):
        jira_keys = commit_batches.extract_jira_keys(
            [
                "修复了 https://jira.meitu.com/browse/DATA-6755 https://jira.meitu.com/browse/TECHPUB-19087 这两个单子",
                "再补一个 DATA-6755 和 INTERNAL-1901",
            ]
        )

        self.assertEqual(jira_keys, ["DATA-6755", "TECHPUB-19087", "INTERNAL-1901"])

    def test_validate_plan_rejects_duplicate_unit_assignment(self):
        self.write_text("app.txt", "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n")
        self.commit_all()
        self.write_text("app.txt", "1x\n2\n3\n4\n5\n6\n7\n8\n9\n10x\n")

        inventory = commit_batches.build_inventory(self.repo, "HEAD", "worktree")
        unit_ids = inventory["files"][0]["unit_ids"]
        plan = self.build_plan(
            inventory,
            [
                {
                    "id": "first",
                    "reason": "take one hunk",
                    "split_mode": "hunk",
                    "units": [unit_ids[0]],
                    "message": {
                        "header": ":wrench: (emoji-commit) first hunk",
                        "body": [],
                    },
                },
                {
                    "id": "second",
                    "reason": "accidentally reuse a unit",
                    "split_mode": "hunk",
                    "units": [unit_ids[0], unit_ids[1]],
                    "message": {
                        "header": ":wrench: (emoji-commit) duplicate unit",
                        "body": [],
                    },
                },
            ],
        )

        with self.assertRaisesRegex(commit_batches.BatchPlanError, "assigned more than once"):
            commit_batches.validate_plan(self.repo, plan)

    def test_validate_plan_rejects_partial_file_coverage_for_file_mode(self):
        self.write_text("app.txt", "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n")
        self.commit_all()
        self.write_text("app.txt", "1x\n2\n3\n4\n5\n6\n7\n8\n9\n10x\n")

        inventory = commit_batches.build_inventory(self.repo, "HEAD", "worktree")
        unit_ids = inventory["files"][0]["unit_ids"]
        plan = self.build_plan(
            inventory,
            [
                {
                    "id": "first",
                    "reason": "incorrectly claim a whole-file split",
                    "split_mode": "file",
                    "units": [unit_ids[0]],
                    "message": {
                        "header": ":wrench: (emoji-commit) misuse file split",
                        "body": [],
                    },
                },
                {
                    "id": "second",
                    "reason": "cover the remaining hunk",
                    "split_mode": "hunk",
                    "units": [unit_ids[1]],
                    "message": {
                        "header": ":wrench: (emoji-commit) remaining hunk",
                        "body": [],
                    },
                },
            ],
        )

        with self.assertRaisesRegex(commit_batches.BatchPlanError, "does not cover full file"):
            commit_batches.validate_plan(self.repo, plan)

    def test_validate_plan_rejects_binary_partial_split(self):
        self.write_bytes("binary.bin", b"\x00\x01\x02base")
        self.commit_all()
        self.write_bytes("binary.bin", b"\x00\x01\x02changed")

        inventory = commit_batches.build_inventory(self.repo, "HEAD", "worktree")
        unit_id = inventory["units"][0]["id"]
        plan = self.build_plan(
            inventory,
            [
                {
                    "id": "binary",
                    "reason": "try to split a binary patch by hunk",
                    "split_mode": "hunk",
                    "units": [unit_id],
                    "message": {
                        "header": ":wrench: (emoji-commit) reject binary hunk split",
                        "body": [],
                    },
                }
            ],
        )

        with self.assertRaisesRegex(commit_batches.BatchPlanError, "unsupported unit"):
            commit_batches.validate_plan(self.repo, plan)

    def test_cli_inventory_accepts_repo_before_and_after_subcommand(self):
        self.write_text("app.txt", "base\n")
        self.commit_all()
        self.write_text("app.txt", "base changed\n")

        before = self.run_cli("--repo", str(self.repo), "inventory", "--scope", "worktree")
        after = self.run_cli("inventory", "--repo", str(self.repo), "--scope", "worktree")

        before_inventory = json.loads(before.stdout)
        after_inventory = json.loads(after.stdout)

        self.assertEqual(before_inventory["input_scope"], "worktree")
        self.assertEqual(before_inventory["stats"], after_inventory["stats"])
        self.assertEqual(before_inventory["files"], after_inventory["files"])

    def test_cli_preview_plan_accepts_repo_before_and_after_subcommand(self):
        self.write_text("app.txt", "one\ntwo\n")
        self.commit_all()
        self.write_text("app.txt", "one changed\ntwo\n")

        inventory = commit_batches.build_inventory(self.repo, "HEAD", "worktree")
        plan = self.build_plan(
            inventory,
            [
                {
                    "id": "single",
                    "reason": "preview the only unit",
                    "split_mode": "file",
                    "units": [inventory["units"][0]["id"]],
                    "message": {
                        "header": ":wrench: (emoji-commit) preview one change",
                        "body": ["render a stable preview"],
                    },
                }
            ],
        )

        with tempfile.NamedTemporaryFile("w", suffix=".json", delete=False) as handle:
            json.dump(plan, handle)
            plan_path = handle.name

        try:
            before = self.run_cli(
                "--repo",
                str(self.repo),
                "preview-plan",
                "--plan",
                plan_path,
            )
            after = self.run_cli(
                "preview-plan",
                "--repo",
                str(self.repo),
                "--plan",
                plan_path,
            )
        finally:
            Path(plan_path).unlink()

        self.assertIn("emoji-commit batch preview", before.stdout)
        self.assertEqual(before.stdout, after.stdout)

    def test_apply_plan_creates_multiple_commits_from_worktree(self):
        self.write_text("app.txt", "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n")
        self.commit_all()

        self.write_text("app.txt", "1x\n2\n3\n4\n5\n6\n7\n8\n9\n10x\n")
        self.write_text("new.txt", "fresh content\n")

        inventory = commit_batches.build_inventory(self.repo, "HEAD", "worktree")
        app_units = inventory["files"][0]["unit_ids"]
        new_unit = next(
            item["id"] for item in inventory["units"] if item["path"] == "new.txt"
        )

        plan = self.build_plan(
            inventory,
            [
                {
                    "id": "first",
                    "reason": "commit the first isolated edit",
                    "split_mode": "hunk",
                    "units": [app_units[0]],
                    "message": {
                        "header": ":wrench: (emoji-commit) commit the first app hunk",
                        "body": ["capture the first edit in app.txt"],
                    },
                },
                {
                    "id": "second",
                    "reason": "commit the second isolated edit",
                    "split_mode": "hunk",
                    "units": [app_units[1]],
                    "message": {
                        "header": ":wrench: (emoji-commit) commit the second app hunk",
                        "body": ["capture the second edit in app.txt"],
                    },
                },
                {
                    "id": "third",
                    "reason": "commit the added file",
                    "split_mode": "file",
                    "units": [new_unit],
                    "message": {
                        "header": ":sparkles: (emoji-commit) add the new worktree file",
                        "body": ["include the untracked file in its own commit"],
                    },
                },
            ],
        )

        validation = commit_batches.validate_plan(self.repo, plan)
        start_head, final_commit, created_commits = commit_batches.create_commits_with_temp_index(
            self.repo,
            validation,
        )
        commit_batches.apply_commits_transaction(self.repo, start_head, final_commit)

        self.assertEqual(len(created_commits), 3)
        self.assertEqual(self.git("rev-parse", "HEAD").stdout.strip(), final_commit)
        self.assertEqual(self.git("status", "--short").stdout.strip(), "")
        self.assertEqual(
            self.git("log", "--pretty=%s", "-3").stdout.strip().splitlines(),
            [
                ":sparkles: (emoji-commit) add the new worktree file",
                ":wrench: (emoji-commit) commit the second app hunk",
                ":wrench: (emoji-commit) commit the first app hunk",
            ],
        )

    def test_validate_plan_rejects_head_drift(self):
        self.write_text("app.txt", "one\ntwo\n")
        self.commit_all()
        self.write_text("app.txt", "one changed\ntwo\n")

        inventory = commit_batches.build_inventory(self.repo, "HEAD", "worktree")
        plan = self.build_plan(
            inventory,
            [
                {
                    "id": "single",
                    "reason": "take the only unit",
                    "split_mode": "file",
                    "units": [inventory["units"][0]["id"]],
                    "message": {
                        "header": ":wrench: (emoji-commit) apply one change",
                        "body": [],
                    },
                }
            ],
        )

        self.write_text("other.txt", "other\n")
        self.git("add", "other.txt")
        self.git("commit", "-m", "advance head")

        with self.assertRaisesRegex(commit_batches.BatchPlanError, "HEAD changed since preview"):
            commit_batches.validate_plan(self.repo, plan)

    def test_apply_plan_rolls_back_head_and_index_on_failure(self):
        self.write_text("app.txt", "one\ntwo\n")
        self.commit_all()
        self.write_text("app.txt", "one changed\ntwo\n")

        inventory = commit_batches.build_inventory(self.repo, "HEAD", "worktree")
        plan = self.build_plan(
            inventory,
            [
                {
                    "id": "single",
                    "reason": "apply the only unit",
                    "split_mode": "file",
                    "units": [inventory["units"][0]["id"]],
                    "message": {
                        "header": ":wrench: (emoji-commit) roll back on apply failure",
                        "body": ["exercise transactional rollback"],
                    },
                }
            ],
        )

        validation = commit_batches.validate_plan(self.repo, plan)
        start_head, final_commit, _ = commit_batches.create_commits_with_temp_index(
            self.repo,
            validation,
        )
        before_status = self.git("status", "--short").stdout
        real_run_git = commit_batches.run_git

        def flaky_run_git(repo_path, args, **kwargs):
            result = real_run_git(repo_path, args, **kwargs)
            if args == ["read-tree", final_commit]:
                raise commit_batches.BatchPlanError("forced read-tree failure")
            return result

        with mock.patch.object(commit_batches, "run_git", side_effect=flaky_run_git):
            with self.assertRaisesRegex(commit_batches.BatchPlanError, "forced read-tree failure"):
                commit_batches.apply_commits_transaction(self.repo, start_head, final_commit)

        self.assertEqual(self.git("rev-parse", "HEAD").stdout.strip(), start_head)
        self.assertEqual(self.git("status", "--short").stdout, before_status)


if __name__ == "__main__":
    unittest.main()
