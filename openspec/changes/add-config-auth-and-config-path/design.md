# Design: config auth 与 -c 配置路径

## Context
- 当前配置通过 Viper 管理，Init() 使用 SetConfigName + AddConfigPath 多路径搜索，无显式单文件路径。
- 需要新增 `config auth` 写入 url/token，以及全局 `-c` 指定配置文件，且两者需一致：-c 指定的文件既被读取也被 config auth 写入。

## Goals / Non-Goals
- **Goals**: 提供 `config auth` 交互式/标志式配置入口；全局 `-c` 指定配置文件路径；读写同一文件语义一致。
- **Non-Goals**: 多 profile、加密存储、GUI 配置不在本次范围。

## Decisions
- **-c 生效时机**：在 Cobra 执行子命令前解析根命令的 PersistentFlags，在 config.Init() 中或 Init 前根据 -c 调用 Viper.SetConfigFile(path)，使后续 ReadInConfig / 写入均针对该文件。若未传 -c，保持现有 AddConfigPath 搜索逻辑。
- **config auth 写入目标**：与当前 Viper 使用的配置文件一致。即若通过 -c 指定了文件，则写入该文件；否则写入 Viper 当前选中的 config 文件（若从未 ReadInConfig 成功，则需约定默认路径，例如 `~/.config/gitlab-tools/config.yaml`）。
- **Alternatives considered**：仅环境变量/仅手动编辑——不利于首次体验；单独「写路径」与「读路径」——增加复杂度，本次不采用。

## Risks / Trade-offs
- **Init 顺序**：-c 必须在 Viper.ReadInConfig 之前生效，因此根命令 init() 中先绑定 -c，再在 Init() 内根据 -c 设置 SetConfigFile 再 ReadInConfig。
- **config auth 无 -c 且无现有文件**：约定写入默认路径（如 `~/.config/gitlab-tools/config.yaml`）并创建父目录，与 project.md 中「用户主目录下的 .config 目录」一致。

## Migration Plan
- 无破坏性变更：未使用 -c 时行为与现有一致；现有环境变量与 --url/--token 仍优先。
- 无需数据迁移。

## Open Questions
- 无。
