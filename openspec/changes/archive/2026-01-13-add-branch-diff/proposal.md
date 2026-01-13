# Change: 添加分支差异比较命令

## Why
当前工具支持列出项目分支，但不支持比较不同分支之间的差异。用户需要查看两个分支之间的差异时（如提交差异、文件变更等），只能通过 GitLab Web 界面或 Git 命令，不够便捷。添加 `branch diff` 命令可以提供命令行方式快速查看分支之间的差异信息。

## What Changes
- 添加 `branch diff <项目ID> <源分支> <目标分支>` 命令
- 支持使用项目 ID（数字）或项目路径（如 `my-group/my-project`）来标识项目
- 显示两个分支之间的差异信息，包括：
  - 提交差异（从源分支到目标分支的提交列表）
  - 文件变更统计（新增、修改、删除的文件数量）
  - 可选的详细文件差异信息
- 添加 `--stat` 参数，仅显示文件变更统计信息
- 添加 `--commits` 参数，仅显示提交差异列表
- 添加 `--create-mr` 参数，在显示差异后创建 Merge Request
- 添加 `--mr-title` 参数，指定 Merge Request 的标题（与 `--create-mr` 一起使用）
- 添加 `--mr-description` 参数，指定 Merge Request 的描述（与 `--create-mr` 一起使用）
- 遵循现有的命令结构和错误处理模式

## Impact
- **Affected specs**: `branch-management` (新增能力)
- **Affected code**: 
  - `main.go`: 添加 `branchDiffCmd` 命令定义
  - `main.go`: 添加 `runBranchDiffCmd` 函数处理分支差异比较
  - `main.go`: 添加 `printBranchDiff` 函数用于格式化输出差异信息
  - `main.go`: 添加 `--stat`、`--commits`、`--create-mr`、`--mr-title`、`--mr-description` 标志变量和处理逻辑
  - `main.go`: 添加创建 Merge Request 的功能（使用 GitLab API 的 MergeRequests.CreateMergeRequest 方法）
  - `main.go`: 在 `init()` 函数中注册新命令和标志

