# Tasks: refactor-agent-first-atomic-api

## Phase 1：核心实现（已完成）

### 1. 契约与文档（原子集与退出码）

- [x] 1.1 在 README 与 capabilities 中写入「原子命令集」列表（命令名、用途、主要参数），作为单一事实来源
- [x] 1.2 在文档与 main 中明确退出码约定：0 成功，1 业务/API 错误，2 用法错误；统一子命令退出行为（ErrUsage → 2）
- [x] 1.3 在 `internal/output` 与各模块确立「格式化器」约定（人类可读 vs JSON）

### 2. 机器可读输出与退出码实现

- [x] 2.1 根命令增加全局 `--json` 标志
- [x] 2.2 实现 `project list` / `project get` 的 JSON 输出与错误时结构化 stderr
- [x] 2.3 实现 `pipeline list` / `get` / `latest` / `check-schedule` 的 JSON 输出
- [x] 2.4 实现 `branch list` / `branch diff` 的 JSON 输出
- [x] 2.5 实现 `mr list` / `mr create` / `mr merge` 的 JSON 输出
- [x] 2.6 实现 `tag list` / `tag create` / `tag delete` 的 JSON 输出
- [x] 2.7 用法错误返回 config.ErrUsage（exit 2）：project list 非法组合、pipeline list 非法 status、pipeline get 非法 ID、mr merge 非法 IID

### 3. Agent 可发现性

- [x] 3.1 实现 `gitlab-tools capabilities`（支持 `--json`）
- [x] 3.2 在 README 与根命令说明 `capabilities` 的用途

### 4. 组合与工作流文档

- [x] 4.1 README 新增「工作流与组合示例」，至少 3 个典型场景
- [x] 4.2 重组 `skills/SKILL.md`：原子命令表 + 推荐组合 + `--json` 说明
- [x] 4.3 README 新增「Agent 与脚本使用」：`--json`、退出码、工作流入口

### 5. 校验与测试

- [x] 5.1 为 `capabilities --json` 增加冒烟测试（main_test.go）
- [x] 5.2 运行 `openspec validate refactor-agent-first-atomic-api --strict` 通过
- [x] 5.3 手动验证：原子命令集与 CLI 一致；退出码符合约定

---

## Phase 2：Spec 对齐（文档与行为）

与细化后的 spec 保持一致：文档化枚举/默认值/排序/语义，补齐 REMOVED 迁移。

### 6. 枚举与默认值文档化

- [x] 6.1 **Pagination**：在 README 或各 list 命令的 `--help` 中文档化 `--limit` 的默认值与上限（agent-interface）
- [x] 6.2 **Pipeline status**：在 `pipeline list --help` 或 README 中明确列出 `--status` 允许值（running, pending, success, failed, canceled, skipped, created, manual）；非法值已返回 exit 2
- [x] 6.3 **MR state**：在 `mr list --help` 或 README 中列出 `--state` 允许值（opened, closed, merged）及默认值（未指定时为 opened）；非法值返回 exit 2（若未实现则补充校验）

### 7. 排序与语义文档化

- [x] 7.1 **List 排序**：在 README 或 `--help` 中说明 pipeline list / tag list / mr list 的排序规则（如 pipeline 按 created_at 降序、newest first），便于 Agent 理解「第一条」
- [x] 7.2 **Project 过滤**：在 README 或 `project list --help` 中说明 `--search`（子串匹配名称/描述）与 `--match`（正则匹配路径/名称）的语义区别
- [x] 7.3 **Tag 创建**：在 README 或 `tag create --help` 中说明 `--ref` 与 `--branch` 的优先级（--ref 优先）及缺省时的默认分支行为
- [x] 7.4 **Branch list 无 project**：在 README 或 `branch list --help` 中注明「不传 project 时列出所有项目分支，可能较慢或结果较大，建议已知 project 时显式传入」

### 8. 错误与 JSON 约定

- [x] 8.1 **Structured error**：确认 `--json` 且发生错误时，stderr 输出结构化错误（如 `{"error":"...","code":1|2}`），并文档化字段（agent-interface）
- [x] 8.2 **Mutation 错误原因**：确认 create/merge/delete 失败时错误信息说明原因（如 MR 已存在、tag 已存在、冲突），便于 Agent 分支（已部分实现，做一次验证即可）

### 9. REMOVED 迁移（project-management）

- [x] 9.1 **Architecture Pattern**：在 `openspec/project.md` 或独立「开发约定」中补充代码组织说明（包结构、internal 划分、main 入口），以便归档时从 project-management spec 中移除该要求后仍有单一来源

### 10. 可选实现与验证

- [x] 10.1 **mr list --state 校验**：若当前未校验 `--state` 非法值，则增加校验并返回 exit 2（与 pipeline list --status 一致）
- [x] 10.2 **List 多项目 pipeline**：确认 `pipeline list <p1> <p2> ...` 的 JSON 输出按项目分组且每组含 project id/path（spec：Pipeline List Filters and Output）。当前实现仅支持单项目，JSON 输出与命令行参数一致；多项目按项目分组为后续扩展。

---

## 执行顺序建议

- Phase 1 已完成，无需改动。
- Phase 2：先做 6.x（枚举与默认值文档化）和 7.x（排序与语义文档化），再做 8.x（验证）、9.1（迁移）、10.x（可选）。同一小节内任务可并行。
