# Change: 添加 Merge Request 管理功能

## Why
当前工具支持创建 Merge Request（通过 `branch diff --create-mr`），但不支持查询和管理已存在的 Merge Request。用户需要：
1. 查询项目中开放的 Merge Request 列表，以便快速了解待处理的合并请求
2. 通过命令行直接合并 Merge Request，无需打开 GitLab Web 界面

添加独立的 Merge Request 管理命令可以提升工作效率，特别是在需要批量处理或自动化场景中。

## What Changes
- 添加 `mr` 命令组，用于管理 Merge Request
- 添加 `mr list <项目ID>` 命令，查询指定项目的开放 Merge Request 列表
- 支持使用项目 ID（数字）或项目路径（如 `my-group/my-project`）来标识项目
- 显示 Merge Request 的基本信息，包括：IID、标题、源分支、目标分支、状态、创建者、创建时间、Web URL 等
- 添加 `mr merge <项目ID> <MR IID>` 命令，合并指定的 Merge Request
- 支持可选的合并参数，如合并提交信息、是否删除源分支等
- 遵循现有的命令结构和错误处理模式

## Impact
- **Affected specs**: `merge-request-management` (新增能力)
- **Affected code**: 
  - `main.go`: 添加 `mrCmd` 命令组定义
  - `main.go`: 添加 `mrListCmd` 命令定义和处理函数
  - `main.go`: 添加 `mrMergeCmd` 命令定义和处理函数
  - `main.go`: 添加 `printMergeRequestsList` 函数用于格式化输出 MR 列表
  - `main.go`: 添加 `printMergeRequestDetails` 函数用于显示 MR 详细信息
  - `main.go`: 在 `init()` 函数中注册新命令组和命令

