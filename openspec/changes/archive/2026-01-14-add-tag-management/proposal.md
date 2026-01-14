# Change: 添加 Tag 管理功能

## Why
当前工具支持分支管理、Merge Request 管理等功能，但不支持创建和管理 Git 标签（Tag）。在软件发布流程中，经常需要为特定版本创建标签。添加 Tag 管理功能可以：
1. 通过命令行快速创建标签，无需打开 GitLab Web 界面
2. 支持在指定分支上打标签，便于版本管理
3. 提升发布流程的自动化程度

## What Changes
- 添加 `tag` 命令组，用于管理 Git 标签
- 添加 `tag list <项目ID>` 命令，查询指定项目的标签列表
- 添加 `tag create <项目ID> <标签名>` 命令，在指定项目上创建标签
- 支持使用项目 ID（数字）或项目路径（如 `my-group/my-project`）来标识项目
- 默认在 `main` 分支上创建标签
- 支持 `--branch` 参数指定目标分支
- 支持 `--ref` 参数指定具体的提交 SHA 或分支名（可选，默认使用分支的最新提交）
- 支持 `--message` 参数指定标签消息（可选）
- 显示标签列表和创建的标签信息，包括标签名、提交 SHA、Web URL 等
- 遵循现有的命令结构和错误处理模式

## Impact
- **Affected specs**: `tag-management` (新增能力)
- **Affected code**: 
  - `main.go`: 添加 `tagCmd` 命令组定义
  - `main.go`: 添加 `tagListCmd` 命令定义和处理函数
  - `main.go`: 添加 `tagCreateCmd` 命令定义和处理函数
  - `main.go`: 添加 `printTagsList` 函数用于格式化输出标签列表
  - `main.go`: 添加 `printTagInfo` 函数用于格式化输出标签信息
  - `main.go`: 在 `init()` 函数中注册新命令组和命令

