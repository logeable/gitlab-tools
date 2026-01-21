# Change: 重构代码文件目录组织

## Why
当前 `main.go` 文件包含所有命令和功能实现，文件行数已超过 1400 行，存在以下问题：
1. **可维护性差**：所有代码集中在一个文件中，难以定位和修改特定功能
2. **可扩展性差**：添加新功能需要修改同一个大文件，容易产生冲突
3. **代码复用困难**：共享逻辑（如客户端创建、配置管理）分散在文件中，难以复用
4. **不符合开源最佳实践**：Go 项目通常采用模块化目录结构，便于协作和维护

参考开源 Go CLI 工具（如 kubectl、helm、gh 等）的最佳实践，需要将代码按功能模块拆分到独立的包中。

## What Changes
- **重构目录结构**：采用标准的 Go 项目布局
  - `cmd/` 目录：存放 `main.go` 和根命令定义
  - `internal/` 目录：存放内部包（不对外暴露）
    - `internal/client`：GitLab 客户端创建和管理
    - `internal/config`：配置管理（Viper 相关）
    - `internal/pipeline`：Pipeline 相关命令和逻辑
    - `internal/project`：Project 相关命令和逻辑
    - `internal/branch`：Branch 相关命令和逻辑
    - `internal/mr`：Merge Request 相关命令和逻辑
    - `internal/tag`：Tag 相关命令和逻辑
    - `internal/output`：输出格式化函数（可选，或放在各自包中）
- **代码拆分**：将 `main.go` 中的代码按功能模块拆分到对应的包中
  - 每个命令组（pipeline, project, branch, mr, tag）独立为一个包
  - 共享的客户端创建逻辑提取到 `internal/client`
  - 配置管理逻辑提取到 `internal/config`
  - 输出格式化函数保留在各自的包中或提取到 `internal/output`
- **保持向后兼容**：重构后 CLI 命令接口和行为保持不变
- **更新项目文档**：更新 `openspec/project.md` 中的架构模式描述

## Impact
- **Affected specs**: `project-management` (架构模式变更)
- **Affected code**: 
  - `main.go`: 大幅简化，只保留根命令定义和 main 函数
  - 新增 `cmd/` 目录结构
  - 新增 `internal/` 目录及其子包
  - 所有功能代码迁移到对应的包中
- **Breaking changes**: 无（CLI 接口保持不变）

