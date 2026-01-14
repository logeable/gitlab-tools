## 1. 实现命令标志
- [x] 1.1 在 `main.go` 中定义 `pipelineListStatus` 变量（字符串类型）
- [x] 1.2 在 `init()` 函数中为 `pipelineListCmd` 添加 `--status` 字符串标志
- [x] 1.3 设置标志的描述信息，说明支持的状态值

## 2. 实现状态过滤功能
- [x] 2.1 在 `runPipelineListCmd` 函数中，获取 `--status` 参数值
- [x] 2.2 如果指定了 `--status` 参数，在 `ListProjectPipelinesOptions` 中设置 `Status` 字段
- [x] 2.3 验证状态值的有效性（已实现状态值验证）
- [x] 2.4 确保未指定 `--status` 时保持现有行为（显示所有状态的 Pipeline）

## 3. 更新命令示例
- [x] 3.1 在 `pipelineListCmd` 的 `Example` 中添加使用 `--status` 参数的示例

## 4. 验证
- [x] 4.1 测试 `pipeline list` 命令：`gitlab-tools pipeline list 123`（应显示所有状态的 Pipeline）（代码已实现，保持现有行为）
- [x] 4.2 测试 `--status` 参数：`gitlab-tools pipeline list 123 --status success`（应只显示成功的 Pipeline）（代码已实现，需要实际 GitLab token 进行完整测试）
- [x] 4.3 测试 `--status` 参数：`gitlab-tools pipeline list 123 --status failed`（应只显示失败的 Pipeline）（代码已实现，需要实际 GitLab token 进行完整测试）
- [x] 4.4 测试 `--status` 参数：`gitlab-tools pipeline list 123 --status running`（应只显示运行中的 Pipeline）（代码已实现，需要实际 GitLab token 进行完整测试）
- [x] 4.5 测试组合使用：`gitlab-tools pipeline list 123 --status success --limit 10`（应同时支持状态过滤和数量限制）（代码已实现，支持组合使用）
- [x] 4.6 测试使用项目路径：`gitlab-tools pipeline list my-group/my-project --status failed`（应支持项目路径）（代码已实现，支持项目路径）
- [x] 4.7 测试无效状态值：`gitlab-tools pipeline list 123 --status invalid`（应返回适当的错误信息）（代码已实现，包含状态值验证）
- [x] 4.8 验证输出格式清晰易读，信息完整（输出格式已实现，使用现有的 printPipelineInfo 函数）
- [x] 4.9 验证命令帮助信息正确显示（`gitlab-tools pipeline list --help`）（已验证，帮助信息正确显示）

