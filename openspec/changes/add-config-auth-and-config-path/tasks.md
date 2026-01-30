# Tasks: add-config-auth-and-config-path

## 1. 全局 -c/--config
- [x] 1.1 在根命令上增加 `-c`/`--config` 持久化标志，文档说明为「配置文件路径」
- [x] 1.2 在 `config.Init()` 或调用前，若 `-c` 已设置则用该路径作为 Viper 的配置文件（SetConfigFile），否则保持现有搜索逻辑
- [x] 1.3 确保子命令（含 `config auth`）均能使用该配置路径读写

## 2. config auth 子命令
- [x] 2.1 新增顶层子命令 `config`，其下子命令 `auth`
- [x] 2.2 `config auth` 支持用户输入 GitLab URL 与 access token（交互式提示或标志），并写入当前生效的配置文件（默认或 -c 指定）
- [x] 2.3 若指定了 `-c`，则将 url/token 写入该路径；若文件不存在则创建（含父目录）
- [x] 2.4 写入后给出简短成功提示；失败时返回明确错误并退出码 1
- [x] 2.5 在 `gitlab-tools --help` 与 `gitlab-tools config auth --help` 中补充说明与 -c 的配合

## 3. 文档与规格
- [x] 3.1 更新 README 或使用说明：如何通过 `config auth` 与 `-c` 配置
- [x] 3.2 更新 capabilities 输出（若有）：加入 `config auth` 及 -c 的说明

## 4. 验收
- [x] 4.1 手动验证：未指定 -c 时 `config auth` 写入默认搜索路径下的 config；指定 -c 时写入该文件
- [x] 4.2 手动验证：`gitlab-tools -c /path/to/config.yaml project list` 使用该文件中的 url/token
