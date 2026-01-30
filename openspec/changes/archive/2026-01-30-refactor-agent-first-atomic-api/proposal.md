# Change: 以 Agent 优先的原子化 API 重构 GitLab 工具

## Why

当前工具已实现 pipeline、project、branch、mr、tag 等能力，但缺少统一设计：命令边界和组合方式未成体系，实际使用场景与管道组合没有沉淀为可复用经验；Agent/Skill 调用依赖自然语言描述和零散示例，难以稳定、可预测地编排。需要明确「原子化命令层」与「组合/工作流」的边界，使 AI Agent 与 Skills 能基于稳定、可发现的接口快速调用并灵活组合。

## What Changes

- **原子化命令集**：定义并文档化一套最小必要、稳定的 GitLab/git 原子操作命令集（每个命令做一件事、参数与输出可预期），作为 CLI 与 Agent 调用的契约。
- **机器可读输出与退出码**：为上述命令提供可选的机器可读输出（如 `--json`）及稳定的退出码约定（成功/业务失败/用法错误），便于脚本与 Agent 解析与编排。
- **组合与工作流文档**：在文档与 Skill 中沉淀「原子命令 + 组合模式」与典型工作流（如：查项目 → 查分支 → 比较 → 创建 MR），使 Agent 能按场景选择并组合原子命令。
- **Agent/Skill 可发现性**：通过 `gitlab-tools` 的 help/子命令与 Skill 描述，明确列出原子能力与推荐组合，便于 Agent 快速发现与调用。
- 不改变现有命令的默认交互行为；新增能力为可选（如 `--json`、新子命令或文档结构）。

## Impact

- **Affected specs**: 新增 `agent-interface`；现有 `project-management`、`pipeline-management`、`branch-management`、`merge-request-management`、`tag-management` 各增加「机器可读输出」相关要求。
- **Affected code**: 根命令与各子命令（help/输出格式）、`internal/output` 或各模块输出层、`skills/SKILL.md`、README 与可选工作流文档。
