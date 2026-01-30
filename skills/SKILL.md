---
name: gitlab-tools
description: Operate the gitlab-tools CLI to query and manage GitLab projects, branches, tags, pipelines, and merge requests. Use when the user asks to perform GitLab actions via gitlab-tools, retrieve repository metadata, or needs command help/output for GitLab operations.
---

# Gitlab Tools

## Overview
Use gitlab-tools to discover projects, inspect branches/tags/pipelines, and manage merge requests or tags via the GitLab API. All commands support `--json` for machine-readable output and script/Agent chaining.

## Quick Start
1. Resolve the project ID or full path (PathWithNamespace) using `gitlab-tools project list` (or `project list --json` to parse id/path).
2. Run the specific subcommand for branches, tags, pipelines, or merge requests.
3. Summarize the result, including status, IDs, URLs, and timestamps when relevant.

## Atomic Commands (discover with `gitlab-tools capabilities` or `capabilities --json`)

| Domain   | Command | Typical use        | Main args |
|----------|---------|--------------------|-----------|
| project  | list    | 列出项目           | [--search] [--match] [--limit] [--quiet] |
| project  | get     | 获取单项目详情     | \<project\> |
| pipeline | list    | 列出 Pipeline      | \<project\> [--status] [--ref] [--limit] |
| pipeline | get     | 获取单条 Pipeline  | \<project\> \<pipeline-id\> |
| pipeline | latest  | 指定 ref 的最新 Pipeline | \<project\> \<ref\> |
| pipeline | check-schedule | 检查 Scheduled Pipeline | \<project\> |
| branch   | list    | 列出分支           | [project] [--search] [--hide-empty] [--quiet] |
| branch   | diff    | 比较两分支差异     | \<project\> \<source\> \<target\> [--stat] [--commits] |
| mr       | list    | 列出 MR            | \<project\> [--state opened\|closed\|merged] |
| mr       | create  | 创建 MR            | \<project\> \<source\> \<target\> [--title] [--description] |
| mr       | merge   | 合并 MR            | \<project\> \<iid\> [--delete-source-branch] |
| tag      | list    | 列出标签           | \<project\> |
| tag      | create  | 创建标签           | \<project\> \<name\> [--branch] [--ref] [--message] |
| tag      | delete  | 删除标签           | \<project\> \<name\> |

## 使用场景

- **项目**：发现/搜索项目（--search 子串、--match 正则）、拿到 path/id 后作为后续命令入参；--has-schedule 找有定时流水线的项目。
- **Pipeline**：本地/CI 查流水线状态；`pipeline latest` 快速看 main/develop 是否绿；`pipeline list --status failed` 筛失败单；`check-schedule` 做定时任务健康检查、脚本根据 exit 码告警。
- **分支**：看某项目分支列表；`branch diff` 合并前看变更范围，确认后再建 MR；不传项目会列所有可访问项目（较慢，建议已知项目时显式传参）。
- **标签**：发版打 tag（--ref/--branch，缺省 main）；误打或重打时 delete；list 看已有 tag 再 create/delete。
- **MR**：查待合并 MR（list 默认 opened）；branch diff 确认后 create；list 取 iid 后 merge，可选 --delete-source-branch。
- **capabilities**：Agent/脚本发现可用命令与参数；加 --json 得到机器可读列表便于解析与链式执行。

## 组合用法

- **发现项目 → 列分支 → 比较差异 → 建 MR**：`project list --json` 解析 id/path → `branch list <project> --json` → `branch diff <project> main feature` → `mr create <project> feature main --title "..."`.
- **看某分支最新流水线**：`pipeline latest <project> <ref> --json`（成功 0，无流水线或失败 1）；配合 `pipeline list <project> --ref <ref>` 看历史。
- **列 open MR 并合并一个**：`mr list <project> --state opened --json` 解析 iid → `mr merge <project> <iid>` [--delete-source-branch].
- **发版打 tag**：`tag list <project>` 看已有 tag → `tag create <project> v1.0.0 --branch main --message "Release 1.0"`；误打则 `tag delete <project> v1.0.0`.
- **定时流水线健康检查**：`pipeline check-schedule <project>`（成功 0，未成功或未配置 1），脚本中根据 exit 码告警。
- **从 list 得到 id 后查详情**：`pipeline list <project> --json` / `pipeline latest ... --json` 取 id → `pipeline get <project> <id>`；`project list --json` 取 path → `project get <path>` 或直接作为后续命令的 <project> 参数。

任意命令加 `--json` 可解析 stdout 做链式调用；退出码：0 成功，1 业务/API 错误，2 用法错误。

## Project Discovery
- List projects: `gitlab-tools project list`
- Search by name or description: `gitlab-tools project list --search "app-backend"`
- Match with regex: `gitlab-tools project list --match "^group/.*backend"`
- JSON for chaining: `gitlab-tools project list --json`

If multiple matches appear, ask the user to confirm the exact project path or ID.

## Branches
- List branches: `gitlab-tools branch list <project>` or `branch list` (all projects)
- Compare branches: `gitlab-tools branch diff <project> <source> <target>`
- With `--json`: output includes name, commit_sha; when listing all projects, each group has project_id/path.

## Tags
- List tags: `gitlab-tools tag list <project>`
- Create tag: `gitlab-tools tag create <project> <name>` [--branch] [--ref] [--message]
- Delete tag: `gitlab-tools tag delete <project> <name>`

## Pipelines
- Latest pipeline: `gitlab-tools pipeline latest <project> <ref>`
- List pipelines: `gitlab-tools pipeline list <project>` [--status success|failed|...] [--limit N]
- Pipeline detail: `gitlab-tools pipeline get <project> <pipeline-id>`
- Scheduled check: `gitlab-tools pipeline check-schedule <project>` (exit 0 = success, 1 = failure or no schedule)

Pipeline status values: running, pending, success, failed, canceled, skipped, created, manual.

## Merge Requests
- List MRs: `gitlab-tools mr list <project>` [--state opened|closed|merged]
- Create MR: `gitlab-tools mr create <project> <source> <target> --title "..." --description "..."`
- Merge MR: `gitlab-tools mr merge <project> <iid>` [--delete-source-branch]

## Global Options
- Token: `--token <token>` or `GITLAB_TOKEN`
- URL: `--url <gitlab-url>` or `GITLAB_URL`
- JSON: `--json` for machine-readable output (snake_case fields)

Use `gitlab-tools <command> --help` to confirm flags before running uncommon options.

## Best Practices & Common Patterns

### Finding projects with specific branch
Use `branch list --search` instead of checking each project individually:
```bash
# Efficient - find all projects with develop branch
gitlab-tools branch list --search develop
```

### Branch diff parameter order
**Important**: The parameter order is `gitlab-tools branch diff <project> <source> <target>`
- **source branch**: comparison baseline (usually `main`)
- **target branch**: branch to compare (usually `develop`)
- Example: `gitlab-tools branchable diff project main develop --stat` (shows develop changes relative to main)

### Project path reliability
Project paths returned by `branch list --search` are verified and usable. Prefer these paths for subsequent operations.

### JSON output
Several commands support `--json` for machine-readable output:
```bash
gitlab-tools project list --json
gitlab-tools pipeline latest project develop --json
```

### Check develop vs main diff workflow
```bash
# 1. Find projects with develop branch
gitlab-tools branch list --search develop

# 2. Check develop relative to main
gitlab-tools branch diff project main develop --stat
```

### Batch operations
```bash
# Check multiple projects for branch differences
for p in "project1" "project2" "project3"; do
  gitlab-tools branch diff "$p" main develop --stat
done

# Check latest pipeline status for multiple projects
for p in "project1" "project2" "project3"; do
  gitlab-tools pipeline latest "$p" develop
done
```
