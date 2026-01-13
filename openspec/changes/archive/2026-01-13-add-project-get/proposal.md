# Change: 添加获取单个项目信息命令

## Why
当前工具支持列出项目列表（`project list`），但不支持获取单个项目的详细信息。用户需要查看特定项目的完整信息时，只能从列表中查找，不够便捷。添加 `project get` 命令可以提供与 `pipeline get` 类似的体验，让用户能够快速查看单个项目的详细信息。

## What Changes
- 添加 `project get <项目ID>` 命令
- 支持使用项目 ID（数字）或项目路径（如 `my-group/my-project`）来标识项目
- 显示项目的详细信息，包括 ID、名称、路径、可见性、默认分支、描述、Web URL、归档状态、最后活动时间等
- 添加 `--detail` 参数，使用 `github.com/k0kubun/pp/v3` 库以带颜色的详细格式展示项目的完整数据结构（适用于调试和查看所有字段）
- 遵循现有的命令结构和错误处理模式

## Impact
- **Affected specs**: `project-management` (新增能力)
- **Affected code**: 
  - `main.go`: 添加 `projectGetCmd` 命令定义和 `runProjectGetCmd` 函数
  - `main.go`: 添加 `printProjectInfo` 函数用于格式化输出项目信息
  - `main.go`: 添加 `--detail` 标志变量和处理逻辑
  - `main.go`: 在 `init()` 函数中注册新命令和标志
  - `go.mod`: 添加 `github.com/k0kubun/pp/v3` 依赖

