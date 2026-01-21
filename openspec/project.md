# Project Context

## Purpose
自用的命令行工具，用于 GitLab 的相关操作。提供便捷的 CLI 接口来管理和查询 GitLab 项目、Pipeline 等信息。

主要功能包括：
- Pipeline 管理：查看 Pipeline 详细信息、列出项目的 Pipeline 列表
- 项目管理：列出项目、搜索和过滤项目

## Tech Stack
- **Go 1.25.2** - 主要编程语言
- **Cobra** (`github.com/spf13/cobra`) - CLI 框架，用于构建命令结构
- **Viper** (`github.com/spf13/viper`) - 配置管理，支持环境变量和命令行参数
- **GitLab API Client** (`gitlab.com/gitlab-org/api/client-go`) - GitLab API 官方客户端库

## Project Conventions

### Code Style
- 使用 Go 标准代码风格和 `gofmt` 格式化
- 变量命名使用驼峰命名法（camelCase）
- 函数命名使用驼峰命名法，导出函数首字母大写
- 错误处理使用 `fmt.Errorf` 包装错误信息
- 所有用户可见的输出和注释使用中文
- 命令行帮助信息使用中文

### Architecture Patterns
- **模块化包结构**：代码按功能模块组织到独立的包中
  - `internal/client`：GitLab 客户端创建和管理
  - `internal/config`：配置管理（Viper 相关）
  - `internal/pipeline`：Pipeline 相关命令和逻辑
  - `internal/project`：Project 相关命令和逻辑
  - `internal/branch`：Branch 相关命令和逻辑
  - `internal/mr`：Merge Request 相关命令和逻辑
  - `internal/tag`：Tag 相关命令和逻辑
  - `internal/output`：共享的输出格式化函数
- **命令结构**：使用 Cobra 构建层次化命令结构
  - 根命令：`gitlab-tools`
  - 子命令：`pipeline`、`project`、`branch`、`mr`、`tag`
  - 子子命令：`pipeline get`、`pipeline list`、`project list` 等
- **配置管理**：使用 Viper 统一管理配置，优先级为：命令行参数 > 环境变量 > 默认值
- **客户端模式**：通过 GitLab API 客户端与 GitLab 服务器交互

### Testing Strategy
- 当前项目未包含测试代码
- 未来可考虑添加单元测试和集成测试

### Git Workflow
- 自用项目，遵循标准的 Git 工作流
- 提交信息使用中文或英文均可

## Domain Context
- **GitLab API**：工具通过 GitLab REST API 与 GitLab 服务器交互
- **Pipeline**：GitLab CI/CD 的构建流水线，包含状态、引用、SHA 等信息
- **Project**：GitLab 项目，包含 ID、路径、可见性、默认分支等信息
- **认证**：使用 GitLab Personal Access Token 进行 API 认证
- **项目标识**：支持使用项目 ID（数字）或项目路径（如 `my-group/my-project`）来标识项目

## Important Constraints
- 需要有效的 GitLab Personal Access Token 才能使用
- 默认连接到 `https://gitlab.com`，可通过 `--url` 参数或环境变量指定其他 GitLab 实例
- Token 可通过 `GITLAB_TOKEN` 环境变量或 `--token` 命令行参数提供
- 工具设计为自用，不包含复杂的错误恢复或重试机制

## External Dependencies
- **GitLab API**：依赖 GitLab 服务器的 REST API
- **GitLab API Client Go**：官方 Go 语言客户端库，处理 API 请求和响应
- **Cobra**：提供 CLI 命令解析和帮助生成
- **Viper**：提供配置管理，支持多种配置源（环境变量、命令行参数等）
