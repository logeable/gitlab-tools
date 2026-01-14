# Change: 为 Pipeline List 添加按状态查询功能

## Why
当前 `pipeline list` 命令支持列出项目的 Pipeline 列表，但不支持按状态过滤。在实际使用中，用户经常需要：
1. 只查看失败的 Pipeline，便于快速定位问题
2. 只查看成功的 Pipeline，确认构建状态
3. 只查看运行中的 Pipeline，了解当前构建进度

添加按状态查询功能可以提升 Pipeline 管理的效率，特别是在需要快速定位问题或查看特定状态的 Pipeline 时。

## What Changes
- 为 `pipeline list` 命令添加 `--status` 参数，支持按状态过滤 Pipeline
- 支持的状态值包括：running, pending, success, failed, canceled, skipped, created, manual
- 如果不指定 `--status` 参数，则显示所有状态的 Pipeline（保持现有行为）
- 支持使用项目 ID（数字）或项目路径（如 `my-group/my-project`）来标识项目
- 遵循现有的命令结构和错误处理模式

## Impact
- **Affected specs**: `pipeline-management` (新增能力)
- **Affected code**: 
  - `main.go`: 添加 `pipelineListStatus` 变量
  - `main.go`: 在 `init()` 函数中为 `pipelineListCmd` 添加 `--status` 字符串标志
  - `main.go`: 在 `runPipelineListCmd` 函数中，如果指定了 `--status` 参数，则在 `ListProjectPipelinesOptions` 中设置 `Status` 字段

