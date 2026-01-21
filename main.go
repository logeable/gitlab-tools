package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/k0kubun/pp/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var (
	gitlabURL   string
	gitlabToken string
)

var rootCmd = &cobra.Command{
	Use:   "gitlab-tools",
	Short: "GitLab 工具集",
	Long:  "一个用于与 GitLab API 交互的命令行工具集",
}

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Pipeline 管理",
	Long:  "查看和管理 GitLab Pipeline",
}

var pipelineGetCmd = &cobra.Command{
	Use:   "get <项目ID> <PipelineID>",
	Short: "获取 pipeline 详细信息",
	Long:  "获取指定项目的 pipeline 状态和详细信息",
	Example: `  gitlab-tools pipeline get 123 456
  gitlab-tools pipeline get my-group/my-project 789`,
	Args: cobra.ExactArgs(2),
	RunE: runPipelineGetCmd,
}

var pipelineListCmd = &cobra.Command{
	Use:   "list <项目ID>...",
	Short: "列出项目的 pipelines",
	Long:  "列出指定项目的 pipeline 列表，每个项目显示最近的几条",
	Example: `  gitlab-tools pipeline list 123
  gitlab-tools pipeline list 123
  gitlab-tools pipeline list my-group/my-project
  gitlab-tools pipeline list 123 --limit 10
  gitlab-tools pipeline list 123 --status success
  gitlab-tools pipeline list 123 --ref main`,
	Args: cobra.ExactArgs(1),
	RunE: runPipelineListCmd,
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "项目管理",
	Long:  "查看和管理 GitLab 项目",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出项目",
	Long:  "列出可访问的项目列表，显示基本信息",
	Example: `  gitlab-tools project list
  gitlab-tools project list --owned
  gitlab-tools project list --search "my-project"
  gitlab-tools project list --match ".*backend.*"
  gitlab-tools project list --match "^my-group/.*"`,
	RunE: runProjectListCmd,
}

var projectGetCmd = &cobra.Command{
	Use:   "get <项目ID>",
	Short: "获取项目详细信息",
	Long:  "获取指定项目的详细信息",
	Example: `  gitlab-tools project get 123
  gitlab-tools project get my-group/my-project
  gitlab-tools project get 123 --detail`,
	Args: cobra.ExactArgs(1),
	RunE: runProjectGetCmd,
}

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "分支管理",
	Long:  "查看和管理 GitLab 项目分支",
}

var branchListCmd = &cobra.Command{
	Use:   "list [项目ID]",
	Short: "列出项目分支",
	Long:  "列出指定项目的分支列表，如果不指定项目 ID 则列出所有可访问项目的分支",
	Example: `  gitlab-tools branch list
  gitlab-tools branch list 123
  gitlab-tools branch list my-group/my-project
  gitlab-tools branch list --search "feature"
  gitlab-tools branch list 123 --search "feature"
  gitlab-tools branch list --hide-empty
  gitlab-tools branch list --quiet
  gitlab-tools branch list --quiet --hide-empty`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBranchListCmd,
}

var branchDiffCmd = &cobra.Command{
	Use:   "diff <项目ID> <源分支> <目标分支>",
	Short: "比较分支差异",
	Long: `比较两个分支之间的差异，显示提交差异和文件变更统计。

源分支（From）：作为比较基准的分支，通常是主分支（如 main 或 master）。
目标分支（To）：要比较的分支，通常是功能分支（如 feature）。

比较结果将显示目标分支相对于源分支的变化，包括：
- 从源分支到目标分支之间的所有提交
- 目标分支中新增、修改、删除的文件

示例：gitlab-tools branch diff 123 main feature
将显示 feature 分支相对于 main 分支的所有变更。`,
	Example: `  gitlab-tools branch diff 123 main feature
  gitlab-tools branch diff my-group/my-project main feature
  gitlab-tools branch diff 123 main feature --stat
  gitlab-tools branch diff 123 main feature --commits`,
	Args: cobra.ExactArgs(3),
	RunE: runBranchDiffCmd,
}

var mrCmd = &cobra.Command{
	Use:   "mr",
	Short: "Merge Request 管理",
	Long:  "查看和管理 GitLab Merge Request",
}

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag 管理",
	Long:  "查看和管理 GitLab 项目标签",
}

var tagListCmd = &cobra.Command{
	Use:   "list <项目ID>",
	Short: "列出项目的标签",
	Long:  "列出指定项目的所有标签列表",
	Example: `  gitlab-tools tag list 123
  gitlab-tools tag list my-group/my-project`,
	Args: cobra.ExactArgs(1),
	RunE: runTagListCmd,
}

var tagCreateCmd = &cobra.Command{
	Use:   "create <项目ID> <标签名>",
	Short: "创建标签",
	Long:  "在指定项目上创建标签，默认在 main 分支上创建",
	Example: `  gitlab-tools tag create 123 v1.0.0
  gitlab-tools tag create my-group/my-project v1.0.0
  gitlab-tools tag create 123 v1.0.0 --branch develop
  gitlab-tools tag create 123 v1.0.0 --ref abc123
  gitlab-tools tag create 123 v1.0.0 --message "版本 1.0.0"`,
	Args: cobra.ExactArgs(2),
	RunE: runTagCreateCmd,
}

var tagDeleteCmd = &cobra.Command{
	Use:   "delete <项目ID> <标签名>",
	Short: "删除标签",
	Long:  "删除指定项目的标签",
	Example: `  gitlab-tools tag delete 123 v1.0.0
  gitlab-tools tag delete my-group/my-project v1.0.0`,
	Args: cobra.ExactArgs(2),
	RunE: runTagDeleteCmd,
}

var mrListCmd = &cobra.Command{
	Use:   "list <项目ID>",
	Short: "列出项目的开放 Merge Request",
	Long:  "列出指定项目的所有开放 Merge Request 列表",
	Example: `  gitlab-tools mr list 123
  gitlab-tools mr list my-group/my-project
  gitlab-tools mr list my-group/my-project --target-branch feature
  gitlab-tools mr list my-group/my-project --state opened
  gitlab-tools mr list my-group/my-project --state closed
  gitlab-tools mr list my-group/my-project --state merged
  gitlab-tools mr list my-group/my-project --with-pipelines`,
	Args: cobra.ExactArgs(1),
	RunE: runMRListCmd,
}

var mrCreateCmd = &cobra.Command{
	Use:   "create <项目ID> <源分支> <目标分支>",
	Short: "创建 Merge Request",
	Long:  "创建新的 Merge Request，从源分支合并到目标分支",
	Example: `  gitlab-tools mr create 123 feature main
  gitlab-tools mr create my-group/my-project feature main
  gitlab-tools mr create 123 feature main --title "我的功能"
  gitlab-tools mr create 123 feature main --title "我的功能" --description "功能描述"
  gitlab-tools mr create 123 feature main --quiet`,
	Args: cobra.ExactArgs(3),
	RunE: runMRCreateCmd,
}

var mrMergeCmd = &cobra.Command{
	Use:   "merge <项目ID> <MR IID>",
	Short: "合并 Merge Request",
	Long:  "合并指定的 Merge Request",
	Example: `  gitlab-tools mr merge 123 456
  gitlab-tools mr merge my-group/my-project 456
  gitlab-tools mr merge 123 456 --delete-source-branch
  gitlab-tools mr merge 123 456 --merge-commit-message "合并信息"`,
	Args: cobra.ExactArgs(2),
	RunE: runMRMergeCmd,
}

var (
	projectOwned         bool
	projectArchived      bool
	projectSearch        string
	projectMatch         string
	projectLimit         int
	projectGetDetail     bool
	pipelineListLimit    int
	pipelineListStatus   string
	branchListSearch     string
	branchListHideEmpty  bool
	branchListQuiet      bool
	branchDiffStat       bool
	branchDiffCommits    bool
	mrCreateTitle        string
	mrCreateDescription  string
	mrCreateQuiet        bool
	mrMergeDeleteSource  bool
	mrMergeCommitMessage string
	mrListTargetBranch   string
	mrListState          string
	mrListWithPipelines  bool
	tagCreateBranch      string
	tagCreateRef         string
	tagCreateMessage     string
	pipelineListRef      string
)

func init() {
	// 初始化 Viper
	initViper()

	// 全局标志
	rootCmd.PersistentFlags().StringVar(&gitlabURL, "url", "", "GitLab 服务器 URL (默认: https://gitlab.com)")
	rootCmd.PersistentFlags().StringVar(&gitlabToken, "token", "", "GitLab 访问令牌 (也可以通过 GITLAB_TOKEN 环境变量设置)")

	// 将 Cobra flags 绑定到 Viper（自动从环境变量读取）
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))

	// project list 标志
	projectListCmd.Flags().BoolVar(&projectOwned, "owned", false, "只显示拥有的项目")
	projectListCmd.Flags().BoolVar(&projectArchived, "archived", false, "包含已归档的项目")
	projectListCmd.Flags().StringVar(&projectSearch, "search", "", "搜索项目名称或描述")
	projectListCmd.Flags().StringVar(&projectMatch, "match", "", "使用正则表达式匹配项目路径或名称")
	projectListCmd.Flags().IntVar(&projectLimit, "limit", 20, "限制返回的项目数量")

	// 绑定 project list 标志到 Viper
	viper.BindPFlag("project.owned", projectListCmd.Flags().Lookup("owned"))
	viper.BindPFlag("project.archived", projectListCmd.Flags().Lookup("archived"))
	viper.BindPFlag("project.search", projectListCmd.Flags().Lookup("search"))
	viper.BindPFlag("project.match", projectListCmd.Flags().Lookup("match"))
	viper.BindPFlag("project.limit", projectListCmd.Flags().Lookup("limit"))

	// pipeline list 标志
	pipelineListCmd.Flags().IntVar(&pipelineListLimit, "limit", 5, "每个项目显示的 pipeline 数量")
	pipelineListCmd.Flags().StringVar(&pipelineListStatus, "status", "", "按状态过滤 Pipeline (running, pending, success, failed, canceled, skipped, created, manual)")

	// 绑定 pipeline list 标志到 Viper
	viper.BindPFlag("pipeline.list.limit", pipelineListCmd.Flags().Lookup("limit"))
	viper.BindPFlag("pipeline.list.status", pipelineListCmd.Flags().Lookup("status"))

	// project get 标志
	projectGetCmd.Flags().BoolVar(&projectGetDetail, "detail", false, "使用详细格式（带颜色）显示完整的项目数据结构")

	// branch list 标志
	branchListCmd.Flags().StringVar(&branchListSearch, "search", "", "按分支名过滤（部分匹配，不区分大小写）")
	branchListCmd.Flags().BoolVar(&branchListHideEmpty, "hide-empty", false, "如果没有分支则隐藏该项目")
	branchListCmd.Flags().BoolVar(&branchListQuiet, "quiet", false, "只显示项目名")

	// branch diff 标志
	branchDiffCmd.Flags().BoolVar(&branchDiffStat, "stat", false, "仅显示文件变更统计信息")
	branchDiffCmd.Flags().BoolVar(&branchDiffCommits, "commits", false, "仅显示提交差异列表")

	// mr create 标志
	mrCreateCmd.Flags().StringVar(&mrCreateTitle, "title", "", "指定 Merge Request 的标题")
	mrCreateCmd.Flags().StringVar(&mrCreateDescription, "description", "", "指定 Merge Request 的描述")
	mrCreateCmd.Flags().BoolVar(&mrCreateQuiet, "quiet", false, "quiet 模式：创建 MR 后只显示链接")

	// mr merge 标志
	mrMergeCmd.Flags().BoolVar(&mrMergeDeleteSource, "delete-source-branch", false, "合并后删除源分支")
	mrMergeCmd.Flags().StringVar(&mrMergeCommitMessage, "merge-commit-message", "", "自定义合并提交信息")

	// mr list 标志
	mrListCmd.Flags().StringVar(&mrListTargetBranch, "target-branch", "", "按目标分支过滤 Merge Request")
	mrListCmd.Flags().StringVar(&mrListState, "state", "", "按状态过滤 Merge Request (opened, closed, merged)")
	mrListCmd.Flags().BoolVar(&mrListWithPipelines, "with-pipelines", false, "显示 Merge Request 的 pipelines")

	// tag create 标志
	tagCreateCmd.Flags().StringVar(&tagCreateBranch, "branch", "main", "指定目标分支（默认: main）")
	tagCreateCmd.Flags().StringVar(&tagCreateRef, "ref", "", "指定具体的提交 SHA 或分支名（可选，默认使用分支的最新提交）")
	tagCreateCmd.Flags().StringVar(&tagCreateMessage, "message", "", "指定标签消息（可选）")

	// pipeline list 标志
	pipelineListCmd.Flags().StringVar(&pipelineListRef, "ref", "", "按 ref 过滤 Pipeline")

	// 添加子命令
	rootCmd.AddCommand(pipelineCmd)
	rootCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(branchCmd)
	rootCmd.AddCommand(mrCmd)
	rootCmd.AddCommand(tagCmd)
	pipelineCmd.AddCommand(pipelineGetCmd)
	pipelineCmd.AddCommand(pipelineListCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectGetCmd)
	branchCmd.AddCommand(branchListCmd)
	branchCmd.AddCommand(branchDiffCmd)
	mrCmd.AddCommand(mrListCmd)
	mrCmd.AddCommand(mrCreateCmd)
	mrCmd.AddCommand(mrMergeCmd)
	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagCreateCmd)
	tagCmd.AddCommand(tagDeleteCmd)
}

func initViper() {
	// 设置环境变量前缀（可选，如果设置则环境变量需要以 GITLAB_TOOLS_ 开头）
	// viper.SetEnvPrefix("GITLAB_TOOLS")

	// 自动读取环境变量
	viper.AutomaticEnv()

	// 将环境变量名中的点号和横线替换为下划线（Viper 的默认行为）
	// 例如：GITLAB_URL, GITLAB_TOKEN
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// 设置默认值
	viper.SetDefault("url", "https://gitlab.com")
	viper.SetDefault("project.limit", 20)
	viper.SetDefault("pipeline.list.limit", 5)
}

func runPipelineGetCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	pipelineIDStr := args[1]

	// 解析 pipeline ID
	pipelineID, err := strconv.Atoi(pipelineIDStr)
	if err != nil {
		return fmt.Errorf("无效的 pipeline ID: %s", pipelineIDStr)
	}

	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	// 获取 pipeline 状态
	pipeline, _, err := client.Pipelines.GetPipeline(projectID, pipelineID)
	if err != nil {
		return fmt.Errorf("获取 pipeline 状态失败: %v", err)
	}

	// 打印 pipeline 信息
	printPipelineInfo(pipeline)

	return nil
}

func runPipelineListCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]

	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	// 从 Viper 获取 limit（支持从环境变量和命令行参数读取）
	// 优先级：命令行参数 > 环境变量 > 默认值
	limit := viper.GetInt("pipeline.list.limit")
	if limit <= 0 {
		// 如果 viper 返回 0 或负数，使用命令行参数的默认值或变量值
		if pipelineListLimit > 0 {
			limit = pipelineListLimit
		} else {
			limit = 5 // 最终默认值
		}
	}

	// 构建查询选项
	opt := &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: limit,
			Page:    1,
		},
		OrderBy: gitlab.Ptr("updated_at"),
		Sort:    gitlab.Ptr("desc"),
	}

	// 如果指定了 --status 参数，设置状态过滤
	status := viper.GetString("pipeline.list.status")
	if status == "" {
		status = pipelineListStatus
	}
	if status != "" {
		// 验证状态值（GitLab API 支持的状态值）
		validStatuses := map[string]bool{
			"running":  true,
			"pending":  true,
			"success":  true,
			"failed":   true,
			"canceled": true,
			"skipped":  true,
			"created":  true,
			"manual":   true,
		}
		if !validStatuses[status] {
			return fmt.Errorf("无效的状态值: %s。支持的状态值: running, pending, success, failed, canceled, skipped, created, manual", status)
		}
		statusValue := gitlab.BuildStateValue(status)
		opt.Status = &statusValue
	}

	if pipelineListRef != "" {
		opt.Ref = &pipelineListRef
	}

	// 获取 pipeline 列表
	pipelines, _, err := client.Pipelines.ListProjectPipelines(projectID, opt)
	if err != nil {
		return fmt.Errorf("获取项目 %s 的 pipeline 列表失败: %v", projectID, err)
	}

	fmt.Printf("项目: %s\n", projectID)
	fmt.Printf("  找到 %d 条 pipeline:\n\n", len(pipelines))
	for i, pipeline := range pipelines {
		pipeline, _, err := client.Pipelines.GetPipeline(projectID, pipeline.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取项目 %s 的 pipeline %d 失败: %v\n", projectID, pipeline.ID, err)
			continue
		}
		if i > 0 {
			fmt.Println()
		}
		printPipelineInfo(pipeline)
	}

	return nil
}

func getGitLabURL() string {
	// Viper 的优先级：命令行参数 > 环境变量 > 默认值
	// 如果命令行参数设置了，viper.GetString 会返回命令行参数的值
	if url := viper.GetString("url"); url != "" {
		return url
	}
	return "https://gitlab.com"
}

func getGitLabToken() string {
	// Viper 自动从环境变量 GITLAB_TOKEN 或命令行参数 --token 读取
	return viper.GetString("token")
}

func runProjectGetCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]

	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	// 获取项目信息
	project, _, err := client.Projects.GetProject(projectID, nil)
	if err != nil {
		return fmt.Errorf("获取项目信息失败: %v", err)
	}

	// 打印项目信息
	printProjectInfo(project, projectGetDetail)

	return nil
}

func runBranchListCmd(cmd *cobra.Command, args []string) error {
	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	// 获取搜索参数
	search := branchListSearch

	// 检查是否提供了项目 ID
	if len(args) > 0 && args[0] != "" {
		// 指定了项目 ID，只获取该项目的分支
		projectID := args[0]
		branches, _, err := client.Branches.ListBranches(projectID, nil)
		if err != nil {
			return fmt.Errorf("获取项目 %s 的分支列表失败: %v", projectID, err)
		}

		// 应用搜索过滤
		if search != "" {
			branches = filterBranchesBySearch(branches, search)
		}

		// 如果启用 --hide-empty 且没有分支，则跳过
		if branchListHideEmpty && len(branches) == 0 {
			return nil
		}

		// 打印分支信息
		printBranchesList(projectID, branches, true, branchListQuiet)
	} else {
		// 未指定项目 ID，获取所有项目的分支
		opt := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    1,
			},
		}

		projects, _, err := client.Projects.ListProjects(opt)
		if err != nil {
			return fmt.Errorf("获取项目列表失败: %v", err)
		}

		if len(projects) == 0 {
			fmt.Println("未找到项目")
			return nil
		}

		// 为每个项目获取分支
		for i, project := range projects {

			branches, _, err := client.Branches.ListBranches(project.PathWithNamespace, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "获取项目 %s 的分支列表失败: %v\n", project.PathWithNamespace, err)
				continue
			}

			// 应用搜索过滤
			if search != "" {
				branches = filterBranchesBySearch(branches, search)
			}

			// 如果启用 --hide-empty 且没有分支，则跳过该项目
			if branchListHideEmpty && len(branches) == 0 {
				continue
			}

			if i > 0 && !branchListQuiet {
				fmt.Println() // 项目之间添加空行（quiet 模式下不需要）
			}

			// 打印分支信息
			printBranchesList(project.PathWithNamespace, branches, false, branchListQuiet)
		}
	}

	return nil
}

func runBranchDiffCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	sourceBranch := args[1]
	targetBranch := args[2]

	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	// 构建比较选项
	opt := &gitlab.CompareOptions{
		From: gitlab.Ptr(sourceBranch),
		To:   gitlab.Ptr(targetBranch),
	}

	// 获取分支差异
	compare, _, err := client.Repositories.Compare(projectID, opt)
	if err != nil {
		return fmt.Errorf("获取分支差异失败: %v", err)
	}

	// 打印分支差异信息
	printBranchDiff(projectID, sourceBranch, targetBranch, compare, branchDiffStat, branchDiffCommits)

	return nil
}

func runMRCreateCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	sourceBranch := args[1]
	targetBranch := args[2]

	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	// 检查分支之间是否有差异
	compareOpt := &gitlab.CompareOptions{
		From: gitlab.Ptr(targetBranch),
		To:   gitlab.Ptr(sourceBranch),
	}

	compare, _, err := client.Repositories.Compare(projectID, compareOpt)
	if err != nil {
		return fmt.Errorf("获取分支差异失败: %v", err)
	}

	// 检查是否有差异：如果没有文件差异，则不创建 MR
	if len(compare.Diffs) == 0 {
		fmt.Println("提示: 两个分支之间没有文件差异，跳过创建 Merge Request")
		return nil
	}

	// 检查是否已经存在相同源分支和目标分支的开放 MR
	existingMRs, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, &gitlab.ListProjectMergeRequestsOptions{
		SourceBranch: gitlab.Ptr(sourceBranch),
		TargetBranch: gitlab.Ptr(targetBranch),
		State:        gitlab.Ptr("opened"),
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	})
	if err != nil {
		return fmt.Errorf("检查现有 Merge Request 失败: %v", err)
	}

	// 如果已存在开放的 MR，则提示并跳过创建
	if len(existingMRs) > 0 {
		fmt.Printf("提示: 源分支 %s 到目标分支 %s 已存在开放的 Merge Request:\n", sourceBranch, targetBranch)
		for _, existingMR := range existingMRs {
			fmt.Printf("  !%d: %s\n", existingMR.IID, existingMR.Title)
			if existingMR.WebURL != "" {
				fmt.Printf("      %s\n", existingMR.WebURL)
			}
		}
		fmt.Println("跳过创建新的 Merge Request")
		return nil
	}

	// 获取当前用户信息，用于设置 assignee
	currentUser, _, err := client.Users.CurrentUser()
	if err != nil {
		return fmt.Errorf("获取当前用户信息失败: %v", err)
	}

	// 设置 MR 标题
	mrTitle := mrCreateTitle
	if mrTitle == "" {
		// 如果没有指定标题，使用默认格式
		mrTitle = fmt.Sprintf("Merge %s into %s", sourceBranch, targetBranch)
	}

	// 创建 Merge Request（从 sourceBranch 合并到 targetBranch）
	mrOpt := &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.Ptr(mrTitle),
		SourceBranch: gitlab.Ptr(sourceBranch),
		TargetBranch: gitlab.Ptr(targetBranch),
		AssigneeID:   gitlab.Ptr(currentUser.ID),
	}

	if mrCreateDescription != "" {
		mrOpt.Description = gitlab.Ptr(mrCreateDescription)
	}

	mr, _, err := client.MergeRequests.CreateMergeRequest(projectID, mrOpt)
	if err != nil {
		return fmt.Errorf("创建 Merge Request 失败: %v", err)
	}

	// 打印 Merge Request 信息
	printMergeRequestInfo(mr, mrCreateQuiet)

	return nil
}

func filterBranchesBySearch(branches []*gitlab.Branch, searchTerm string) []*gitlab.Branch {
	searchLower := strings.ToLower(searchTerm)
	var filtered []*gitlab.Branch
	for _, branch := range branches {
		if strings.Contains(strings.ToLower(branch.Name), searchLower) {
			filtered = append(filtered, branch)
		}
	}
	return filtered
}

func runProjectListCmd(cmd *cobra.Command, args []string) error {
	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	// 从 Viper 获取参数（支持从环境变量读取）
	// 优先级：命令行参数 > 环境变量 > 默认值
	owned := viper.GetBool("project.owned") || projectOwned
	archived := viper.GetBool("project.archived") || projectArchived
	search := viper.GetString("project.search")
	if search == "" {
		search = projectSearch
	}
	match := viper.GetString("project.match")
	if match == "" {
		match = projectMatch
	}
	limit := viper.GetInt("project.limit")
	if limit == 0 {
		limit = projectLimit
	}

	// 构建查询选项
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: limit,
			Page:    1,
		},
	}

	if owned {
		ownedPtr := true
		opt.Owned = &ownedPtr
	}

	if archived {
		archivedPtr := true
		opt.Archived = &archivedPtr
	}

	if search != "" {
		opt.Search = &search
	}

	// 获取项目列表
	projects, _, err := client.Projects.ListProjects(opt)
	if err != nil {
		return fmt.Errorf("获取项目列表失败: %v", err)
	}

	// 如果指定了 match 参数，使用正则表达式过滤项目
	if match != "" {
		projects = filterProjectsByRegex(projects, match)
	}

	// 打印项目信息
	printProjectsList(projects)

	return nil
}

func filterProjectsByRegex(projects []*gitlab.Project, pattern string) []*gitlab.Project {
	// 编译正则表达式
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 无效的正则表达式 '%s': %v\n", pattern, err)
		return projects // 如果正则表达式无效，返回所有项目
	}

	// 过滤项目：匹配路径或名称
	var filtered []*gitlab.Project
	for _, project := range projects {
		// 匹配路径（PathWithNamespace）或名称（Name 或 NameWithNamespace）
		if re.MatchString(project.PathWithNamespace) ||
			re.MatchString(project.Name) ||
			re.MatchString(project.NameWithNamespace) {
			filtered = append(filtered, project)
		}
	}

	return filtered
}

func printProjectInfo(project *gitlab.Project, useDetail bool) {
	if useDetail {
		pp.Print(project)
		return
	}

	fmt.Printf("项目信息:\n")
	fmt.Printf("  ID: %d\n", project.ID)
	fmt.Printf("  名称: %s\n", project.Name)
	fmt.Printf("  路径: %s\n", project.PathWithNamespace)
	fmt.Printf("  可见性: %s\n", project.Visibility)
	if project.DefaultBranch != "" {
		fmt.Printf("  默认分支: %s\n", project.DefaultBranch)
	}
	if project.Description != "" {
		fmt.Printf("  描述: %s\n", project.Description)
	}
	fmt.Printf("  Web URL: %s\n", project.WebURL)
	if project.Archived {
		fmt.Printf("  状态: 已归档\n")
	}
	if project.LastActivityAt != nil {
		fmt.Printf("  最后活动: %s\n", formatToLocalTime(project.LastActivityAt))
	}
	if project.CreatedAt != nil {
		fmt.Printf("  创建时间: %s\n", formatToLocalTime(project.CreatedAt))
	}
}

func printProjectsList(projects []*gitlab.Project) {
	if len(projects) == 0 {
		fmt.Println("未找到项目")
		return
	}

	fmt.Printf("找到 %d 个项目:\n\n", len(projects))
	for i, project := range projects {
		fmt.Printf("[%d] %s\n", i+1, project.NameWithNamespace)
		fmt.Printf("    ID: %d\n", project.ID)
		fmt.Printf("    路径: %s\n", project.PathWithNamespace)
		fmt.Printf("    可见性: %s\n", project.Visibility)
		fmt.Printf("    默认分支: %s\n", project.DefaultBranch)
		if project.Description != "" {
			fmt.Printf("    描述: %s\n", project.Description)
		}
		fmt.Printf("    Web URL: %s\n", project.WebURL)
		if project.Archived {
			fmt.Printf("    状态: 已归档\n")
		}
		if project.LastActivityAt != nil {
			fmt.Printf("    最后活动: %s\n", formatToLocalTime(project.LastActivityAt))
		}
		fmt.Println()
	}
}

func printBranchesList(projectID string, branches []*gitlab.Branch, singleProject bool, quiet bool) {
	if quiet {
		// quiet 模式：只显示项目名
		if len(branches) > 0 {
			fmt.Println(projectID)
		}
		return
	}

	if singleProject {
		fmt.Printf("项目: %s\n", projectID)
	} else {
		fmt.Printf("项目: %s\n", projectID)
	}

	if len(branches) == 0 {
		fmt.Println("  未找到分支")
		return
	}

	fmt.Printf("  找到 %d 个分支:\n\n", len(branches))
	for i, branch := range branches {
		fmt.Printf("  [%d] %s", i+1, branch.Name)
		if branch.Protected {
			fmt.Printf(" (受保护)")
		}
		if branch.Default {
			fmt.Printf(" (默认分支)")
		}
		fmt.Println()

		if branch.Commit != nil {
			sha := branch.Commit.ID
			if len(sha) > 8 {
				sha = sha[:8]
			}
			fmt.Printf("      最后提交: %s\n", sha)
			if branch.Commit.Message != "" {
				// 只显示提交信息的第一行
				message := strings.Split(branch.Commit.Message, "\n")[0]
				if len(message) > 60 {
					message = message[:60] + "..."
				}
				fmt.Printf("      提交信息: %s\n", message)
			}
			if branch.Commit.CommittedDate != nil {
				fmt.Printf("      提交时间: %s\n", formatToLocalTime(branch.Commit.CommittedDate))
			}
			if branch.Commit.AuthorName != "" {
				fmt.Printf("      提交者: %s\n", branch.Commit.AuthorName)
			}
			if branch.Commit.AuthorEmail != "" {
				fmt.Printf("      提交者邮箱: %s\n", branch.Commit.AuthorEmail)
			}
		}
		fmt.Println()
	}
}

func printPipelineInfo(pipeline *gitlab.Pipeline) {
	fmt.Printf("Pipeline 信息:\n")
	fmt.Printf("  ID: %d\n", pipeline.ID)
	fmt.Printf("  状态: %s\n", pipeline.Status)
	fmt.Printf("  引用: %s\n", pipeline.Ref)
	fmt.Printf("  SHA: %s\n", pipeline.SHA)
	fmt.Printf("  创建时间: %s\n", formatToLocalTime(pipeline.CreatedAt))
	fmt.Printf("  更新时间: %s\n", formatToLocalTime(pipeline.UpdatedAt))
	fmt.Printf("  源: %s\n", pipeline.Source)
	fmt.Printf("  是否为 tag: %t\n", pipeline.Tag)
	fmt.Printf("  Web URL: %s\n", pipeline.WebURL)

	if pipeline.Duration > 0 {
		dur := time.Duration(pipeline.Duration) * time.Second
		fmt.Printf("  持续时间: %v\n", dur)
	}
	if pipeline.Coverage != "" {
		fmt.Printf("  覆盖率: %s\n", pipeline.Coverage)
	}
}

func printBranchDiff(projectID, sourceBranch, targetBranch string, compare *gitlab.Compare, statOnly, commitsOnly bool) {
	fmt.Printf("项目: %s\n", projectID)
	fmt.Printf("源分支: %s\n", sourceBranch)
	fmt.Printf("目标分支: %s\n", targetBranch)
	fmt.Println()

	// 如果仅显示提交列表
	if commitsOnly {
		if len(compare.Commits) == 0 {
			fmt.Println("无提交差异")
			return
		}
		fmt.Printf("提交差异 (%d 个提交):\n\n", len(compare.Commits))
		for i, commit := range compare.Commits {
			sha := commit.ID
			if len(sha) > 8 {
				sha = sha[:8]
			}
			fmt.Printf("  [%d] %s\n", i+1, sha)
			if commit.AuthorName != "" {
				fmt.Printf("      作者: %s", commit.AuthorName)
				if commit.AuthorEmail != "" {
					fmt.Printf(" <%s>", commit.AuthorEmail)
				}
				fmt.Println()
			}
			if commit.Message != "" {
				message := strings.Split(commit.Message, "\n")[0]
				if len(message) > 80 {
					message = message[:80] + "..."
				}
				fmt.Printf("      提交信息: %s\n", message)
			}
			if commit.CommittedDate != nil {
				fmt.Printf("      提交时间: %s\n", formatToLocalTime(commit.CommittedDate))
			}
			fmt.Println()
		}
		return
	}

	// 如果仅显示统计信息
	if statOnly {
		added := 0
		modified := 0
		deleted := 0
		renamed := 0

		for _, diff := range compare.Diffs {
			if diff.NewFile {
				added++
			} else if diff.DeletedFile {
				deleted++
			} else if diff.RenamedFile {
				renamed++
			} else {
				modified++
			}
		}

		fmt.Printf("文件变更统计:\n")
		fmt.Printf("  新增: %d\n", added)
		fmt.Printf("  修改: %d\n", modified)
		fmt.Printf("  删除: %d\n", deleted)
		if renamed > 0 {
			fmt.Printf("  重命名: %d\n", renamed)
		}
		fmt.Printf("  总计: %d\n", len(compare.Diffs))
		return
	}

	// 显示完整信息
	// 提交差异
	if len(compare.Commits) == 0 {
		fmt.Println("无提交差异")
	} else {
		fmt.Printf("提交差异 (%d 个提交):\n\n", len(compare.Commits))
		for i, commit := range compare.Commits {
			sha := commit.ID
			if len(sha) > 8 {
				sha = sha[:8]
			}
			fmt.Printf("  [%d] %s\n", i+1, sha)
			if commit.AuthorName != "" {
				fmt.Printf("      作者: %s", commit.AuthorName)
				if commit.AuthorEmail != "" {
					fmt.Printf(" <%s>", commit.AuthorEmail)
				}
				fmt.Println()
			}
			if commit.Message != "" {
				message := strings.Split(commit.Message, "\n")[0]
				if len(message) > 80 {
					message = message[:80] + "..."
				}
				fmt.Printf("      提交信息: %s\n", message)
			}
			if commit.CommittedDate != nil {
				fmt.Printf("      提交时间: %s\n", formatToLocalTime(commit.CommittedDate))
			}
			fmt.Println()
		}
	}

	// 文件变更统计
	added := 0
	modified := 0
	deleted := 0
	renamed := 0

	for _, diff := range compare.Diffs {
		if diff.NewFile {
			added++
		} else if diff.DeletedFile {
			deleted++
		} else if diff.RenamedFile {
			renamed++
		} else {
			modified++
		}
	}

	fmt.Printf("文件变更统计:\n")
	fmt.Printf("  新增: %d\n", added)
	fmt.Printf("  修改: %d\n", modified)
	fmt.Printf("  删除: %d\n", deleted)
	if renamed > 0 {
		fmt.Printf("  重命名: %d\n", renamed)
	}
	fmt.Printf("  总计: %d\n", len(compare.Diffs))
	fmt.Println()

	// 详细文件差异
	if len(compare.Diffs) > 0 {
		fmt.Printf("文件变更详情:\n\n")
		for i, diff := range compare.Diffs {
			changeType := "修改"
			if diff.NewFile {
				changeType = "新增"
			} else if diff.DeletedFile {
				changeType = "删除"
			} else if diff.RenamedFile {
				changeType = "重命名"
			}

			fmt.Printf("  [%d] %s: %s\n", i+1, changeType, diff.NewPath)
			if diff.RenamedFile && diff.OldPath != diff.NewPath {
				fmt.Printf("      原路径: %s\n", diff.OldPath)
			}
		}
	}
}

func printMergeRequestInfo(mr *gitlab.MergeRequest, quiet bool) {
	if quiet {
		// quiet 模式：只显示链接
		if mr.WebURL != "" {
			fmt.Println(mr.WebURL)
		}
		return
	}

	// 正常模式：显示完整信息
	fmt.Println()
	fmt.Println("Merge Request 已创建:")
	fmt.Printf("  ID: %d\n", mr.IID)
	fmt.Printf("  标题: %s\n", mr.Title)
	if mr.Description != "" {
		fmt.Printf("  描述: %s\n", mr.Description)
	}
	fmt.Printf("  源分支: %s\n", mr.SourceBranch)
	fmt.Printf("  目标分支: %s\n", mr.TargetBranch)
	fmt.Printf("  状态: %s\n", mr.State)
	if mr.WebURL != "" {
		fmt.Printf("  Web URL: %s\n", mr.WebURL)
	}
	if mr.CreatedAt != nil {
		fmt.Printf("  创建时间: %s\n", formatToLocalTime(mr.CreatedAt))
	}
}

func formatToLocalTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

func runMRListCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]

	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	var targetBranch *string
	if mrListTargetBranch != "" {
		targetBranch = &mrListTargetBranch
	}

	var state *string
	if mrListState != "" {
		state = &mrListState
	}

	opt := &gitlab.ListProjectMergeRequestsOptions{
		State: state,
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		TargetBranch: targetBranch,
	}

	// 获取 Merge Request 列表
	mrs, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, opt)
	if err != nil {
		return fmt.Errorf("获取项目 %s 的 Merge Request 列表失败: %v", projectID, err)
	}

	// 打印 Merge Request 列表
	fmt.Printf("项目: %s\n", projectID)
	if len(mrs) == 0 {
		fmt.Println("  未找到 Merge Request")
		return nil
	}

	fmt.Printf("  找到 %d 个 Merge Request:\n\n", len(mrs))
	for i, mr := range mrs {
		fmt.Printf("  [%d] !%d: %s\n", i+1, mr.IID, mr.Title)
		fmt.Printf("      源分支: %s -> 目标分支: %s\n", mr.SourceBranch, mr.TargetBranch)
		fmt.Printf("      状态: %s\n", mr.State)
		fmt.Printf("      合并状态: %s\n", mr.DetailedMergeStatus)
		if mr.Author != nil {
			fmt.Printf("      创建者: %s", mr.Author.Name)
			if mr.Author.Username != "" {
				fmt.Printf(" (@%s)", mr.Author.Username)
			}
			fmt.Println()
		}
		if mr.CreatedAt != nil {
			fmt.Printf("      创建时间: %s\n", formatToLocalTime(mr.CreatedAt))
		}
		if mr.WebURL != "" {
			fmt.Printf("      Web URL: %s\n", mr.WebURL)
		}

		if mrListWithPipelines {
			pipelines, _, err := client.MergeRequests.ListMergeRequestPipelines(projectID, mr.IID)
			if err != nil {
				return fmt.Errorf("获取 Merge Request %d 的 pipelines 失败: %v", mr.IID, err)
			}
			for _, pipeline := range pipelines {
				pipeline, _, err := client.Pipelines.GetPipeline(projectID, pipeline.ID)
				if err != nil {
					return fmt.Errorf("获取 pipeline %d 失败: %v", pipeline.ID, err)
				}
				printPipelineInfo(pipeline)
				fmt.Println()
			}
		}
		fmt.Println()
	}

	return nil
}

func runMRMergeCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	mrIIDStr := args[1]

	// 解析 MR IID
	mrIID, err := strconv.Atoi(mrIIDStr)
	if err != nil {
		return fmt.Errorf("无效的 MR IID: %s", mrIIDStr)
	}

	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	// 构建合并选项
	mrOpt := &gitlab.AcceptMergeRequestOptions{}

	if mrMergeDeleteSource {
		mrOpt.ShouldRemoveSourceBranch = gitlab.Ptr(true)
	}

	if mrMergeCommitMessage != "" {
		mrOpt.MergeCommitMessage = gitlab.Ptr(mrMergeCommitMessage)
	}

	// 合并 Merge Request
	mr, _, err := client.MergeRequests.AcceptMergeRequest(projectID, mrIID, mrOpt)
	if err != nil {
		return fmt.Errorf("合并 Merge Request 失败: %v", err)
	}

	// 打印合并结果
	printMergeRequestDetails(mr)

	return nil
}

func printMergeRequestDetails(mr *gitlab.MergeRequest) {
	fmt.Println("Merge Request 已合并:")
	fmt.Printf("  ID: %d\n", mr.IID)
	fmt.Printf("  标题: %s\n", mr.Title)
	if mr.Description != "" {
		fmt.Printf("  描述: %s\n", mr.Description)
	}
	fmt.Printf("  源分支: %s\n", mr.SourceBranch)
	fmt.Printf("  目标分支: %s\n", mr.TargetBranch)
	fmt.Printf("  状态: %s\n", mr.State)
	if mr.MergedAt != nil {
		fmt.Printf("  合并时间: %s\n", formatToLocalTime(mr.MergedAt))
	}
	if mr.WebURL != "" {
		fmt.Printf("  Web URL: %s\n", mr.WebURL)
	}
}

func runTagListCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]

	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	// 获取标签列表
	tags, _, err := client.Tags.ListTags(projectID, nil)
	if err != nil {
		return fmt.Errorf("获取项目 %s 的标签列表失败: %v", projectID, err)
	}

	// 打印标签列表
	printTagsList(projectID, tags)

	return nil
}

func runTagCreateCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	tagName := args[1]

	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	// 获取分支参数（默认值为 "main"）
	branch := tagCreateBranch
	if branch == "" {
		branch = "main"
	}

	// 确定 ref（如果指定了 --ref 则使用该值，否则使用分支名）
	ref := tagCreateRef
	if ref == "" {
		ref = branch
	}

	// 构建创建标签选项
	tagOpt := &gitlab.CreateTagOptions{
		TagName: gitlab.Ptr(tagName),
		Ref:     gitlab.Ptr(ref),
	}

	if tagCreateMessage != "" {
		tagOpt.Message = gitlab.Ptr(tagCreateMessage)
	}

	// 创建标签
	tag, _, err := client.Tags.CreateTag(projectID, tagOpt)
	if err != nil {
		return fmt.Errorf("创建标签失败: %v", err)
	}

	// 打印标签信息
	printTagInfo(tag)

	return nil
}

func printTagsList(projectID string, tags []*gitlab.Tag) {
	fmt.Printf("项目: %s\n", projectID)
	if len(tags) == 0 {
		fmt.Println("  未找到标签")
		return
	}

	fmt.Printf("  找到 %d 个标签:\n\n", len(tags))
	for i, tag := range tags {
		fmt.Printf("  [%d] %s\n", i+1, tag.Name)
		if tag.Commit != nil {
			sha := tag.Commit.ID
			if len(sha) > 8 {
				sha = sha[:8]
			}
			fmt.Printf("      提交: %s\n", sha)
			if tag.Commit.Message != "" {
				message := strings.Split(tag.Commit.Message, "\n")[0]
				if len(message) > 60 {
					message = message[:60] + "..."
				}
				fmt.Printf("      提交信息: %s\n", message)
			}
			if tag.Commit.CommittedDate != nil {
				fmt.Printf("      提交时间: %s\n", formatToLocalTime(tag.Commit.CommittedDate))
			}
			if tag.Commit.AuthorName != "" {
				fmt.Printf("      提交者: %s\n", tag.Commit.AuthorName)
			}
		}
		if tag.Message != "" {
			fmt.Printf("      标签消息: %s\n", tag.Message)
		}
		if tag.Release != nil && tag.Release.Description != "" {
			fmt.Printf("      发布说明: %s\n", tag.Release.Description)
		}
		fmt.Println()
	}
}

func printTagInfo(tag *gitlab.Tag) {
	fmt.Println()
	fmt.Println("标签已创建:")
	fmt.Printf("  标签名: %s\n", tag.Name)
	if tag.Commit != nil {
		fmt.Printf("  提交: %s\n", tag.Commit.ID)
		if tag.Commit.Message != "" {
			message := strings.Split(tag.Commit.Message, "\n")[0]
			if len(message) > 80 {
				message = message[:80] + "..."
			}
			fmt.Printf("  提交信息: %s\n", message)
		}
		if tag.Commit.CommittedDate != nil {
			fmt.Printf("  提交时间: %s\n", formatToLocalTime(tag.Commit.CommittedDate))
		}
		if tag.Commit.AuthorName != "" {
			fmt.Printf("  提交者: %s\n", tag.Commit.AuthorName)
		}
	}
	if tag.Message != "" {
		fmt.Printf("  标签消息: %s\n", tag.Message)
	}
}

func runTagDeleteCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	tagName := args[1]

	// 获取配置
	url := getGitLabURL()
	token := getGitLabToken()
	if token == "" {
		return fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	// 创建 GitLab 客户端
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	// 删除标签
	_, err = client.Tags.DeleteTag(projectID, tagName)
	if err != nil {
		return fmt.Errorf("删除标签失败: %v", err)
	}

	fmt.Printf("标签 %s 已删除\n", tagName)

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
