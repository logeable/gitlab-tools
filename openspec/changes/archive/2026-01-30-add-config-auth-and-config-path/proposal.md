# Change: 支持 config auth 与 -c 指定配置路径

## Why
用户当前只能通过环境变量、命令行 `--url`/`--token` 或手动编辑配置文件来设置 GitLab URL 与 access token，缺少一次性交互式配置入口；同时无法在命令中指定配置文件路径，不利于多环境或多项目使用不同配置。

## What Changes
- 新增顶层子命令 `config auth`，支持用户输入 GitLab URL 与 access token 并写入配置文件。
- 支持全局参数 `-c`/`--config` 指定配置文件路径；指定后该次运行读写均使用该路径（`config auth` 写入、其他命令读取）。
- 将「配置管理」作为新能力纳入规格（config-management），并在 agent-interface 中补充 config 子命令与全局 -c 的约定。

## Impact
- Affected specs: agent-interface（原子命令集与全局 -c）、新增 config-management
- Affected code: `main.go`（根命令 -c、config 子命令注册）、`internal/config/config.go`（Init 支持显式 config 路径）、新增 `internal/config/` 下 auth 相关逻辑或命令
