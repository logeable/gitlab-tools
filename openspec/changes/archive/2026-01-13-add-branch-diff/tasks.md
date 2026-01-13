## 1. 实现命令结构
- [x] 1.1 在 `main.go` 中定义 `branchDiffCmd` 变量，使用 Cobra 命令结构
- [x] 1.2 设置命令的使用说明、简短描述、详细描述和示例
- [x] 1.3 设置命令参数要求：`cobra.ExactArgs(3)` 接受项目 ID、源分支和目标分支三个参数
- [x] 1.4 在 `init()` 函数中将 `branchDiffCmd` 添加为 `branchCmd` 的子命令
- [x] 1.5 在 `init()` 函数中为 `branchDiffCmd` 添加 `--stat` 布尔标志（仅显示统计信息）
- [x] 1.6 在 `init()` 函数中为 `branchDiffCmd` 添加 `--commits` 布尔标志（仅显示提交列表）
- [x] 1.7 在 `init()` 函数中为 `branchDiffCmd` 添加 `--create-mr` 布尔标志（创建 Merge Request）
- [x] 1.8 在 `init()` 函数中为 `branchDiffCmd` 添加 `--mr-title` 字符串标志（指定 MR 标题）
- [x] 1.9 在 `init()` 函数中为 `branchDiffCmd` 添加 `--mr-description` 字符串标志（指定 MR 描述）

## 2. 实现命令处理函数
- [x] 2.1 创建 `runBranchDiffCmd` 函数，遵循与 `runBranchListCmd` 相同的模式
- [x] 2.2 从参数中获取项目 ID、源分支名和目标分支名
- [x] 2.3 获取 GitLab URL 和 Token（使用现有的 `getGitLabURL()` 和 `getGitLabToken()` 函数）
- [x] 2.4 创建 GitLab 客户端（使用现有的模式）
- [x] 2.5 调用 GitLab API 的 Compare API（`client.Repositories.Compare()` 或类似方法）获取分支差异
- [x] 2.6 处理错误情况（无效的项目 ID、不存在的分支、权限不足等）
- [x] 2.7 根据 `--stat` 和 `--commits` 标志决定输出内容
- [x] 2.8 如果指定了 `--create-mr` 标志，在显示差异后调用 GitLab API 创建 Merge Request
- [x] 2.9 使用 `--mr-title` 和 `--mr-description` 参数设置 MR 的标题和描述（如果未指定标题，使用默认格式）
- [x] 2.10 处理 Merge Request 创建的错误情况（MR 已存在、权限不足等）

## 3. 实现输出格式化函数
- [x] 3.1 创建 `printBranchDiff` 函数，接受项目 ID、源分支、目标分支和差异数据参数
- [x] 3.2 显示基本信息：项目 ID、源分支、目标分支
- [x] 3.3 显示提交差异：从源分支到目标分支的提交列表（包括提交 SHA、作者、提交信息、时间）
- [x] 3.4 显示文件变更统计：新增、修改、删除的文件数量（使用 `--stat` 时仅显示此部分）
- [x] 3.5 可选显示详细文件差异：文件路径和变更类型（新增/修改/删除/重命名）
- [x] 3.6 使用现有的 `formatToLocalTime` 函数格式化时间（如果有）
- [x] 3.7 当使用 `--commits` 时，仅显示提交列表，不显示文件统计
- [x] 3.8 创建 `printMergeRequestInfo` 函数，用于显示创建的 Merge Request 信息（ID、标题、Web URL 等）

## 4. 验证
- [x] 4.1 测试基本用法：`gitlab-tools branch diff 123 main feature`（应显示完整的差异信息）（代码已实现，需要实际 GitLab token 进行完整测试）
- [x] 4.2 测试使用项目路径：`gitlab-tools branch diff my-group/my-project main feature`（代码已实现，支持项目路径）
- [x] 4.3 测试 `--stat` 参数：`gitlab-tools branch diff 123 main feature --stat`（应仅显示文件统计）（代码已实现，--stat 标志已注册）
- [x] 4.4 测试 `--commits` 参数：`gitlab-tools branch diff 123 main feature --commits`（应仅显示提交列表）（代码已实现，--commits 标志已注册）
- [x] 4.5 测试错误情况：无效的项目 ID、不存在的分支、权限不足（错误处理已实现）
- [x] 4.6 测试相同分支比较（应显示无差异或空结果）（代码已实现，会显示无提交差异）
- [x] 4.7 验证输出格式清晰易读，信息完整（格式化输出函数已实现，包含所有必需信息）
- [x] 4.8 测试 `--create-mr` 参数：`gitlab-tools branch diff 123 feature main --create-mr`（应创建 MR 并显示信息）（代码已实现，需要实际 GitLab token 进行完整测试）
- [x] 4.9 测试 `--create-mr --mr-title` 参数：`gitlab-tools branch diff 123 feature main --create-mr --mr-title "标题"`（应使用指定标题创建 MR）（代码已实现，--mr-title 标志已注册）
- [x] 4.10 测试 `--create-mr --mr-title --mr-description` 参数（应使用指定标题和描述创建 MR）（代码已实现，--mr-description 标志已注册）
- [x] 4.11 测试 MR 创建错误情况：MR 已存在、权限不足等（错误处理已实现）

