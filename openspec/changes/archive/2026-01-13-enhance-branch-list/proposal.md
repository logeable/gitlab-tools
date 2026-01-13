# Change: 增强分支列表命令功能

## Why
基础的 `branch list` 命令已经实现，但在实际使用中发现需要以下增强功能：
1. 当列出所有项目时，有些项目可能没有分支，显示这些项目会干扰查看，需要能够隐藏空项目
2. 在某些场景下（如脚本处理），只需要项目名列表，不需要详细的分支信息
3. 分支信息中缺少提交者信息，这对于了解分支的维护者很有用

## What Changes
- 添加 `--hide-empty` 参数：如果没有分支则隐藏该项目（适用于列出所有项目时）
- 添加 `--quiet` 参数：只显示项目名，不显示详细的分支信息（适用于脚本处理和快速查看）
- 在分支详细信息中显示提交者信息：包括提交者姓名和邮箱
- 优化输出格式：在多项目模式下，空行分隔的位置调整到检查 `--hide-empty` 之后

## Impact
- **Affected specs**: `branch-management` (增强现有能力)
- **Affected code**: 
  - `main.go`: 添加 `branchListHideEmpty` 和 `branchListQuiet` 标志变量
  - `main.go`: 在 `init()` 函数中注册新标志
  - `main.go`: 在 `runBranchListCmd` 函数中添加 `--hide-empty` 和 `--quiet` 的处理逻辑
  - `main.go`: 在 `printBranchesList` 函数中添加 quiet 模式支持和提交者信息显示

