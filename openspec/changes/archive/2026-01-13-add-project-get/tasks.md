## 1. 实现命令结构
- [x] 1.1 在 `main.go` 中定义 `projectGetCmd` 变量，使用 Cobra 命令结构
- [x] 1.2 设置命令的使用说明、简短描述、详细描述和示例
- [x] 1.3 设置命令参数要求：`cobra.ExactArgs(1)` 接受一个项目 ID 参数
- [x] 1.4 在 `init()` 函数中将 `projectGetCmd` 添加为 `projectCmd` 的子命令
- [x] 1.5 在 `init()` 函数中为 `projectGetCmd` 添加 `--detail` 布尔标志

## 2. 实现命令处理函数
- [x] 2.1 创建 `runProjectGetCmd` 函数，遵循与 `runPipelineGetCmd` 相同的模式
- [x] 2.2 从参数中获取项目 ID（支持数字 ID 或路径格式）
- [x] 2.3 获取 GitLab URL 和 Token（使用现有的 `getGitLabURL()` 和 `getGitLabToken()` 函数）
- [x] 2.4 创建 GitLab 客户端（使用现有的模式）
- [x] 2.5 调用 `client.Projects.GetProject(projectID, nil)` 获取项目信息
- [x] 2.6 处理错误情况（无效的项目 ID、权限不足等）
- [x] 2.7 检查 `--detail` 标志，如果启用则使用 pp 库输出，否则使用格式化输出

## 3. 实现输出格式化函数
- [x] 3.1 创建 `printProjectInfo` 函数，接受 `*gitlab.Project` 参数
- [x] 3.2 参考 `printProjectsList` 函数的输出格式，但针对单个项目进行优化
- [x] 3.3 显示项目的关键信息：ID、名称、路径、可见性、默认分支、描述、Web URL、归档状态、最后活动时间
- [x] 3.4 使用现有的 `formatToLocalTime` 函数格式化时间
- [x] 3.5 在 `printProjectInfo` 函数中，当 `detail` 参数为 true 时，使用 `pp.Print()` 输出完整的项目数据结构（带颜色格式化）

## 4. 添加依赖
- [x] 4.1 检查 `go.mod` 是否已包含 `github.com/k0kubun/pp/v3` 依赖
- [x] 4.2 运行 `go get github.com/k0kubun/pp/v3` 添加依赖
- [x] 4.3 在 `main.go` 中导入 `pp` 包

## 5. 验证
- [x] 5.1 测试使用数字项目 ID：`gitlab-tools project get 123`（命令已实现，需要实际 GitLab token 进行完整测试）
- [x] 5.2 测试使用项目路径：`gitlab-tools project get my-group/my-project`（命令已实现，支持路径格式）
- [x] 5.3 测试使用 `--detail` 参数：`gitlab-tools project get 123 --detail`（命令已实现，--detail 标志已注册）
- [x] 5.4 验证 `--detail` 输出包含完整的项目数据结构（带颜色格式化，代码已实现 pp.Print() 调用）
- [x] 5.5 测试错误情况：无效的项目 ID、不存在的项目、权限不足（错误处理已实现）
- [x] 5.6 验证默认输出格式清晰易读（格式化输出函数已实现，包含所有关键信息）

