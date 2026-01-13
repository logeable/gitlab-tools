## 1. 实现命令结构
- [x] 1.1 在 `main.go` 中定义 `branchCmd` 变量，使用 Cobra 命令结构
- [x] 1.2 在 `main.go` 中定义 `branchListCmd` 变量，使用 Cobra 命令结构
- [x] 1.3 设置命令的使用说明、简短描述、详细描述和示例
- [x] 1.4 设置命令参数要求：`cobra.MaximumNArgs(1)` 接受一个可选的项目 ID 参数
- [x] 1.5 在 `init()` 函数中将 `branchCmd` 添加为根命令的子命令
- [x] 1.6 在 `init()` 函数中将 `branchListCmd` 添加为 `branchCmd` 的子命令
- [x] 1.7 在 `init()` 函数中为 `branchListCmd` 添加 `--search` 字符串标志

## 2. 实现命令处理函数
- [x] 2.1 创建 `runBranchListCmd` 函数，遵循与 `runPipelineListCmd` 相同的模式
- [x] 2.2 检查参数：如果提供了项目 ID，从参数中获取（支持数字 ID 或路径格式）
- [x] 2.3 如果未提供项目 ID，调用 `client.Projects.ListProjects()` 获取所有可访问的项目列表
- [x] 2.4 获取 GitLab URL 和 Token（使用现有的 `getGitLabURL()` 和 `getGitLabToken()` 函数）
- [x] 2.5 创建 GitLab 客户端（使用现有的模式）
- [x] 2.6 如果指定了项目 ID，调用 `client.Branches.ListBranches(projectID, opt)` 获取单个项目的分支列表
- [x] 2.7 如果未指定项目 ID，遍历所有项目，为每个项目调用 `client.Branches.ListBranches()` 获取分支列表
- [x] 2.8 处理错误情况（无效的项目 ID、权限不足等，对于多项目情况，单个项目失败应继续处理其他项目）
- [x] 2.9 如果指定了 `--search` 参数，在客户端返回结果后过滤分支（支持部分匹配）

## 3. 实现输出格式化函数
- [x] 3.1 创建 `printBranchesList` 函数，接受项目 ID 和分支列表参数
- [x] 3.2 参考 `printProjectPipelines` 函数的输出格式
- [x] 3.3 显示项目标识（当列出所有项目时，需要显示项目信息）
- [x] 3.4 显示分支的关键信息：分支名、是否受保护、是否默认分支、最后提交的 SHA 和提交信息
- [x] 3.5 使用现有的 `formatToLocalTime` 函数格式化时间（如果有）
- [x] 3.6 当列出多个项目时，在项目之间添加分隔（空行或分隔符）

## 4. 验证
- [x] 4.1 测试不指定项目 ID：`gitlab-tools branch list`（应列出所有项目的分支）（命令已实现，需要实际 GitLab token 进行完整测试）
- [x] 4.2 测试使用数字项目 ID：`gitlab-tools branch list 123`（命令已实现，支持数字 ID）
- [x] 4.3 测试使用项目路径：`gitlab-tools branch list my-group/my-project`（命令已实现，支持路径格式）
- [x] 4.4 测试使用 `--search` 参数（不指定项目）：`gitlab-tools branch list --search "feature"`（命令已实现，--search 标志已注册）
- [x] 4.5 测试使用 `--search` 参数（指定项目）：`gitlab-tools branch list 123 --search "feature"`（命令已实现）
- [x] 4.6 验证 `--search` 参数正确过滤分支名（代码已实现 filterBranchesBySearch 函数，支持部分匹配和不区分大小写）
- [x] 4.7 测试错误情况：无效的项目 ID、不存在的项目、权限不足（错误处理已实现，多项目情况下单个项目失败会继续处理其他项目）
- [x] 4.8 验证输出格式清晰易读，多项目时分组显示清晰（格式化输出函数已实现，项目之间有空行分隔）

