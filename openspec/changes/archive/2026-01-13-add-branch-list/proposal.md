# Change: 添加列出项目分支命令

## Why
当前工具支持管理项目和 Pipeline，但不支持查看项目的分支信息。用户需要查看项目的分支列表时，只能通过 GitLab Web 界面或 Git 命令，不够便捷。添加 `branch list` 命令可以提供命令行方式快速查看项目的所有分支信息。

## What Changes
- 添加 `branch list [项目ID]` 命令（项目 ID 为可选参数）
- 如果指定项目 ID，支持使用项目 ID（数字）或项目路径（如 `my-group/my-project`）来标识项目
- 如果不指定项目 ID，则列出所有可访问项目的分支
- 显示分支的详细信息，包括分支名、是否受保护、是否默认分支、最后提交信息等
- 添加 `--search` 参数，支持按分支名进行过滤（支持部分匹配）
- 遵循现有的命令结构和错误处理模式

## Impact
- **Affected specs**: `branch-management` (新增能力)
- **Affected code**: 
  - `main.go`: 添加 `branchCmd` 和 `branchListCmd` 命令定义
  - `main.go`: 添加 `runBranchListCmd` 函数（支持可选项目 ID，如果不指定则获取所有项目）
  - `main.go`: 添加 `printBranchesList` 函数用于格式化输出分支信息
  - `main.go`: 添加 `--search` 标志变量和处理逻辑
  - `main.go`: 在 `init()` 函数中注册新命令和标志

