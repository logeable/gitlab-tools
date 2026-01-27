# GitLab Tools

一个功能强大的 GitLab 命令行工具集，用于与 GitLab API 交互，提供项目管理、Pipeline 管理、分支管理、Merge Request 管理和标签管理等核心功能。

## ✨ 特性

- 🚀 **项目管理** - 列出、搜索、获取项目信息，支持 Pipeline Schedule 查询
- 🔄 **Pipeline 管理** - 查看 Pipeline 状态、列表和最新 Pipeline，支持 Scheduled Pipeline 检查
- 🌿 **分支管理** - 列出项目分支，比较分支差异，查看提交和文件变更统计
- 🔀 **Merge Request 管理** - 创建、列出、合并 Merge Request，支持 Pipeline 状态显示
- 🏷️ **标签管理** - 创建、删除、列出项目标签
- 🤖 **Agent Skills 支持** - 提供 AI Agent Skill，让 AI 助手能够直接使用 gitlab-tools 进行 GitLab 操作

## 📋 要求

- Go 1.25.2 或更高版本
- GitLab 访问令牌（Personal Access Token 或 Project Access Token）

## 🚀 安装

### 从源码构建

```bash
git clone https://github.com/your-username/gitlab-tools.git
cd gitlab-tools
go build -o gitlab-tools
```

### 使用 Go 安装

```bash
go install gitlab-tools@latest
```

## ⚙️ 配置

### 配置文件

复制示例配置文件并填入你的配置：

```bash
cp config.yaml.example config.yaml
```

编辑 `config.yaml`：

```yaml
# GitLab 服务器 URL（默认: https://gitlab.com）
url: https://gitlab.com

# GitLab 访问令牌
token: your-gitlab-token-here
```

### 环境变量

你也可以通过环境变量配置：

```bash
export GITLAB_URL=https://gitlab.com
export GITLAB_TOKEN=your-gitlab-token-here
```

### 命令行参数

所有配置项也可以通过命令行参数覆盖：

```bash
gitlab-tools --url https://gitlab.com --token your-token project list
```

## 📖 使用示例

### 项目管理

```bash
# 列出所有可访问的项目
gitlab-tools project list

# 只显示拥有的项目
gitlab-tools project list --owned

# 搜索项目
gitlab-tools project list --search "my-project"

# 使用正则表达式匹配项目
gitlab-tools project list --match ".*backend.*"

# 只显示配置了 Pipeline Schedule 的项目
gitlab-tools project list --has-schedule

# 获取项目详细信息
gitlab-tools project get 123
gitlab-tools project get my-group/my-project
```

### Pipeline 管理

```bash
# 列出项目的 Pipelines
gitlab-tools pipeline list 123

# 按状态过滤
gitlab-tools pipeline list 123 --status success

# 按分支过滤
gitlab-tools pipeline list 123 --ref main

# 获取 Pipeline 详细信息
gitlab-tools pipeline get 123 456

# 获取指定分支的最新 Pipeline
gitlab-tools pipeline latest 123 main

# 检查最近的 Scheduled Pipeline 是否成功
gitlab-tools pipeline check-schedule 123
```

### 分支管理

```bash
# 列出所有项目的分支
gitlab-tools branch list

# 列出指定项目的分支
gitlab-tools branch list 123

# 搜索分支
gitlab-tools branch list --search "feature"

# 比较分支差异
gitlab-tools branch diff 123 main feature

# 只显示文件变更统计
gitlab-tools branch diff 123 main feature --stat

# 只显示提交差异
gitlab-tools branch diff 123 main feature --commits
```

### Merge Request 管理

```bash
# 列出项目的 Merge Request
gitlab-tools mr list 123

# 按目标分支过滤
gitlab-tools mr list 123 --target-branch main

# 按状态过滤
gitlab-tools mr list 123 --state opened

# 显示 Pipeline 状态
gitlab-tools mr list 123 --with-pipelines

# 创建 Merge Request
gitlab-tools mr create 123 feature main --title "新功能" --description "功能描述"

# 合并 Merge Request
gitlab-tools mr merge 123 456

# 合并后删除源分支
gitlab-tools mr merge 123 456 --delete-source-branch
```

### 标签管理

```bash
# 列出项目的标签
gitlab-tools tag list 123

# 创建标签
gitlab-tools tag create 123 v1.0.0

# 在指定分支创建标签
gitlab-tools tag create 123 v1.0.0 --branch develop

# 在指定提交创建标签
gitlab-tools tag create 123 v1.0.0 --ref abc123

# 创建带消息的标签
gitlab-tools tag create 123 v1.0.0 --message "版本 1.0.0"

# 删除标签
gitlab-tools tag delete 123 v1.0.0
```

## 🤖 Agent Skills

本项目提供了 Agent Skill，让 AI 助手（如 Claude、Cursor 等）能够直接使用 gitlab-tools 进行 GitLab 操作。

### 安装 Skill

将 `skills/SKILL.md` 文件添加到你的 AI 助手技能目录中：

**对于 Claude Desktop:**
```bash
# 复制 skill 文件到 Claude 技能目录
cp skills/SKILL.md ~/.claude/skills/gitlab-tools/SKILL.md
```

**对于 Cursor:**
```bash
# 复制 skill 文件到 Cursor 技能目录
cp skills/SKILL.md ~/.cursor/skills/gitlab-tools/SKILL.md
```

### 使用方式

安装后，AI 助手将能够：

- 自动识别 GitLab 相关的操作请求
- 使用 gitlab-tools 执行项目管理、Pipeline 查询、分支比较等操作
- 理解项目路径和 ID，自动解析和转换
- 提供结构化的结果输出，包括状态、ID、URL 和时间戳

### 示例对话

安装 skill 后，你可以直接与 AI 助手对话：

```
用户: "帮我查看 my-group/my-project 项目的最新 pipeline 状态"
AI: [自动使用 gitlab-tools pipeline latest my-group/my-project main]
```

```
用户: "列出所有包含 'backend' 的项目"
AI: [自动使用 gitlab-tools project list --match ".*backend.*"]
```

### Skill 功能

Agent Skill 支持以下操作：

- **项目发现** - 搜索和匹配项目
- **分支管理** - 列出分支、比较分支差异
- **标签操作** - 创建、删除、列出标签
- **Pipeline 查询** - 获取最新 Pipeline、列表和详细信息
- **Merge Request** - 创建、列出、合并 MR

更多详细信息请查看 [skills/SKILL.md](skills/SKILL.md)。

## 📚 命令参考

### 全局参数

- `--url`: GitLab 服务器 URL（默认: https://gitlab.com）
- `--token`: GitLab 访问令牌

### 项目命令 (`project`)

- `list`: 列出项目
  - `--owned`: 只显示拥有的项目
  - `--archived`: 包含已归档的项目
  - `--search`: 搜索项目名称或描述
  - `--match`: 使用正则表达式匹配项目路径或名称
  - `--limit`: 限制返回的项目数量（默认: 20）
  - `--has-schedule`: 只显示配置了 Pipeline Schedule 的项目
  - `--schedule-detail`: 输出 Pipeline Schedule 的详细信息
  - `--quiet`: 只输出项目名称
- `get <项目ID>`: 获取项目详细信息
  - `--detail`: 使用详细格式显示完整的项目数据结构

### Pipeline 命令 (`pipeline`)

- `list <项目ID>`: 列出项目的 Pipelines
  - `--limit`: 每个项目显示的 Pipeline 数量（默认: 5）
  - `--status`: 按状态过滤（running, pending, success, failed, canceled, skipped, created, manual）
  - `--ref`: 按 ref 过滤
- `get <项目ID> <PipelineID>`: 获取 Pipeline 详细信息
- `latest <项目ID> <分支名>`: 获取指定分支的最新 Pipeline
- `check-schedule <项目ID>`: 检查最近的 Scheduled Pipeline 是否成功

### 分支命令 (`branch`)

- `list [项目ID]`: 列出项目分支
  - `--search`: 按分支名过滤（部分匹配，不区分大小写）
  - `--hide-empty`: 如果没有分支则隐藏该项目
  - `--quiet`: 只显示项目名
- `diff <项目ID> <源分支> <目标分支>`: 比较分支差异
  - `--stat`: 仅显示文件变更统计信息
  - `--commits`: 仅显示提交差异列表

### Merge Request 命令 (`mr`)

- `list <项目ID>`: 列出项目的 Merge Request
  - `--target-branch`: 按目标分支过滤
  - `--state`: 按状态过滤（opened, closed, merged）
  - `--with-pipelines`: 显示 Merge Request 的 Pipelines
- `create <项目ID> <源分支> <目标分支>`: 创建 Merge Request
  - `--title`: 指定 Merge Request 的标题
  - `--description`: 指定 Merge Request 的描述
  - `--quiet`: 创建 MR 后只显示链接
- `merge <项目ID> <MR IID>`: 合并 Merge Request
  - `--delete-source-branch`: 合并后删除源分支
  - `--merge-commit-message`: 自定义合并提交信息

### 标签命令 (`tag`)

- `list <项目ID>`: 列出项目的标签
- `create <项目ID> <标签名>`: 创建标签
  - `--branch`: 指定目标分支（默认: main）
  - `--ref`: 指定具体的提交 SHA 或分支名
  - `--message`: 指定标签消息
- `delete <项目ID> <标签名>`: 删除标签

## 🔧 开发

### 构建

```bash
go build -o gitlab-tools
```

### 运行测试

```bash
go test ./...
```

## 🤝 贡献

欢迎贡献！请遵循以下步骤：

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 📝 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🙏 致谢

- [GitLab Go API Client](https://gitlab.com/gitlab-org/api/client-go) - GitLab API 客户端库
- [Cobra](https://github.com/spf13/cobra) - 命令行框架
- [Viper](https://github.com/spf13/viper) - 配置管理

## 📮 问题反馈

如果你遇到任何问题或有功能建议，请在 [Issues](https://github.com/your-username/gitlab-tools/issues) 中提交。
