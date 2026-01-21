# Design: 代码重构架构设计

## Context
当前项目采用单一文件架构，所有代码集中在 `main.go`（1494 行）。随着功能增加，文件变得难以维护。需要重构为模块化的包结构，提升可维护性和可扩展性。

## Goals / Non-Goals

### Goals
- 将代码按功能模块拆分到独立包中
- 提取共享逻辑（客户端创建、配置管理）到公共包
- 保持 CLI 命令接口和行为完全不变
- 遵循 Go 项目标准目录结构
- 提升代码可读性和可维护性

### Non-Goals
- 不改变 CLI 命令接口
- 不改变功能行为
- 不引入新的依赖
- 不进行功能增强（仅重构）

## Decisions

### Decision: 采用 `internal/` 目录存放内部包
**Rationale**: 
- `internal/` 是 Go 的标准约定，表示这些包只供项目内部使用，不对外暴露
- 符合 Go 项目最佳实践
- 避免外部包意外依赖内部实现

**Alternatives considered**:
- 使用 `pkg/` 目录：通常用于可被外部项目导入的公共包，不符合当前需求
- 使用扁平结构：所有包放在根目录，不够清晰

### Decision: 按功能模块拆分（pipeline, project, branch, mr, tag）
**Rationale**:
- 每个模块对应一个 CLI 命令组，职责清晰
- 便于独立开发和测试
- 符合单一职责原则

**Alternatives considered**:
- 按层次拆分（commands, handlers, formatters）：可能过度设计，当前规模不需要
- 按数据模型拆分：不符合 CLI 工具的使用模式

### Decision: 提取客户端和配置到独立包
**Rationale**:
- 客户端创建逻辑被所有命令复用
- 配置管理逻辑集中，便于维护
- 减少代码重复

### Decision: 输出格式化函数保留在各自包中
**Rationale**:
- 每个模块的输出格式相对独立
- 避免创建不必要的抽象层
- 如果未来需要共享格式化逻辑，可以再提取到 `internal/output`

**Alternatives considered**:
- 统一提取到 `internal/output`：当前各模块输出差异较大，统一抽象可能过度设计

## Directory Structure

```
gitlab-tools/
├── cmd/
│   └── gitlab-tools/
│       └── main.go              # 根命令定义和 main 函数
├── internal/
│   ├── client/
│   │   └── client.go            # GitLab 客户端创建和管理
│   ├── config/
│   │   └── config.go            # 配置管理（Viper）
│   ├── pipeline/
│   │   ├── command.go           # Pipeline 命令定义
│   │   ├── get.go               # pipeline get 命令实现
│   │   ├── list.go              # pipeline list 命令实现
│   │   └── output.go            # Pipeline 输出格式化
│   ├── project/
│   │   ├── command.go           # Project 命令定义
│   │   ├── get.go               # project get 命令实现
│   │   ├── list.go              # project list 命令实现
│   │   └── output.go            # Project 输出格式化
│   ├── branch/
│   │   ├── command.go           # Branch 命令定义
│   │   ├── list.go              # branch list 命令实现
│   │   ├── diff.go              # branch diff 命令实现
│   │   └── output.go            # Branch 输出格式化
│   ├── mr/
│   │   ├── command.go           # MR 命令定义
│   │   ├── list.go              # mr list 命令实现
│   │   ├── create.go            # mr create 命令实现
│   │   ├── merge.go             # mr merge 命令实现
│   │   └── output.go            # MR 输出格式化
│   └── tag/
│       ├── command.go           # Tag 命令定义
│       ├── list.go              # tag list 命令实现
│       ├── create.go            # tag create 命令实现
│       ├── delete.go            # tag delete 命令实现
│       └── output.go            # Tag 输出格式化
├── main.go                      # 保留（向后兼容）或删除
├── go.mod
└── go.sum
```

## Implementation Strategy

### Phase 1: 创建目录结构和共享包
1. 创建 `internal/client` 包，提取客户端创建逻辑
2. 创建 `internal/config` 包，提取配置管理逻辑
3. 创建各功能模块目录结构

### Phase 2: 迁移功能模块（逐个迁移）
1. 迁移 pipeline 模块
2. 迁移 project 模块
3. 迁移 branch 模块
4. 迁移 mr 模块
5. 迁移 tag 模块

### Phase 3: 重构主文件
1. 简化 `main.go`，只保留根命令定义
2. 将根命令定义移到 `cmd/gitlab-tools/main.go`
3. 验证所有命令正常工作

## Risks / Trade-offs

### Risk: 重构过程中可能引入 bug
**Mitigation**: 
- 逐个模块迁移，每迁移一个模块立即测试
- 保持原有代码不变，直到新代码验证通过
- 使用版本控制，便于回滚

### Risk: 包之间的循环依赖
**Mitigation**:
- 共享逻辑（client, config）放在独立的包中
- 功能模块之间不直接依赖，只依赖共享包
- 输出格式化函数放在各自包中，避免跨包依赖

### Trade-off: 文件数量增加 vs 可维护性提升
- **增加**: 文件数量从 1 个增加到约 20+ 个
- **收益**: 每个文件职责单一，代码更易理解和维护
- **结论**: 可维护性提升的收益远大于文件数量增加的成本

## Migration Plan

1. **准备阶段**：创建目录结构，提取共享代码
2. **迁移阶段**：逐个模块迁移，每迁移一个立即测试
3. **验证阶段**：完整测试所有命令，确保行为一致
4. **清理阶段**：删除旧的 `main.go` 代码（如果保留在根目录）

## Open Questions
- 是否需要保留根目录的 `main.go` 作为向后兼容？建议删除，使用标准的 `cmd/` 目录结构。

