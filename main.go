package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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
  gitlab-tools pipeline list 123 456 789
  gitlab-tools pipeline list my-group/my-project
  gitlab-tools pipeline list 123 --limit 10`,
	Args: cobra.MinimumNArgs(1),
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

var (
	projectOwned      bool
	projectArchived   bool
	projectSearch     string
	projectMatch      string
	projectLimit      int
	pipelineListLimit int
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

	// 绑定 pipeline list 标志到 Viper
	viper.BindPFlag("pipeline.list.limit", pipelineListCmd.Flags().Lookup("limit"))

	// 添加子命令
	rootCmd.AddCommand(pipelineCmd)
	rootCmd.AddCommand(projectCmd)
	pipelineCmd.AddCommand(pipelineGetCmd)
	pipelineCmd.AddCommand(pipelineListCmd)
	projectCmd.AddCommand(projectListCmd)
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

	// 获取项目 ID 列表
	projectIDs := args

	// 为每个项目获取 pipeline 列表
	for i, projectID := range projectIDs {
		if i > 0 {
			fmt.Println() // 项目之间添加空行
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

		// 获取 pipeline 列表
		pipelines, _, err := client.Pipelines.ListProjectPipelines(projectID, opt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取项目 %s 的 pipeline 列表失败: %v\n", projectID, err)
			continue
		}

		// 打印项目 pipelines
		printProjectPipelines(projectID, pipelines)
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

func printProjectPipelines(projectID string, pipelines []*gitlab.PipelineInfo) {
	fmt.Printf("项目: %s\n", projectID)
	if len(pipelines) == 0 {
		fmt.Println("  未找到 pipeline")
		return
	}

	fmt.Printf("  找到 %d 条 pipeline:\n\n", len(pipelines))
	for i, pipeline := range pipelines {
		fmt.Printf("  [%d] Pipeline #%d (IID: %d)\n", i+1, pipeline.ID, pipeline.IID)
		fmt.Printf("      状态: %s\n", pipeline.Status)
		fmt.Printf("      引用: %s\n", pipeline.Ref)
		sha := pipeline.SHA
		if len(sha) > 8 {
			sha = sha[:8]
		}
		fmt.Printf("      SHA: %s\n", sha)
		fmt.Printf("      源: %s\n", pipeline.Source)
		if pipeline.Name != "" {
			fmt.Printf("      名称: %s\n", pipeline.Name)
		}
		if pipeline.CreatedAt != nil {
			fmt.Printf("      创建时间: %s\n", formatToLocalTime(pipeline.CreatedAt))
		}
		if pipeline.UpdatedAt != nil {
			fmt.Printf("      更新时间: %s\n", formatToLocalTime(pipeline.UpdatedAt))
		}
		fmt.Printf("      Web URL: %s\n", pipeline.WebURL)
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
	fmt.Printf("  Web URL: %s\n", pipeline.WebURL)

	if pipeline.Duration > 0 {
		dur := time.Duration(pipeline.Duration) * time.Second
		fmt.Printf("  持续时间: %v\n", dur)
	}
	if pipeline.Coverage != "" {
		fmt.Printf("  覆盖率: %s\n", pipeline.Coverage)
	}
}

func formatToLocalTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
