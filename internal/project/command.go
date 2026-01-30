package project

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	projectOwned          bool
	projectArchived       bool
	projectSearch         string
	projectMatch          string
	projectLimit          int
	projectGetDetail      bool
	projectHasSchedule    bool
	projectScheduleDetail bool
	projectQuiet          bool
)

// NewCommand 创建并返回 project 命令组
func NewCommand() *cobra.Command {
	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "列出、搜索、获取项目",
		Long:  "列出可访问项目、按名称/路径过滤、获取单项目详情。",
	}

	projectListCmd := &cobra.Command{
		Use:   "list",
		Short: "列出项目（支持过滤与分页）",
		Long:  "列出当前可访问的项目。--search 子串匹配名称/描述，--match 正则匹配路径/名称；--has-schedule 只显示有定时流水线的项目。",
		Example: `  gitlab-tools project list
  gitlab-tools project list --owned
  gitlab-tools project list --search "my-project"
  gitlab-tools project list --match ".*backend.*"
  gitlab-tools project list --match "^my-group/.*"
  gitlab-tools project list --has-schedule
  gitlab-tools project list --has-schedule --schedule-detail
  gitlab-tools project list --has-schedule --quiet`,
		RunE: runListCmd,
	}

	projectGetCmd := &cobra.Command{
		Use:   "get <项目ID>",
		Short: "按 ID 或路径获取单项目详情",
		Long:  "根据项目 ID 或路径获取单项目详情。",
		Example: `  gitlab-tools project get 123
  gitlab-tools project get my-group/my-project
  gitlab-tools project get 123 --detail`,
		Args: cobra.ExactArgs(1),
		RunE: runGetCmd,
	}

	// project list 标志
	projectListCmd.Flags().BoolVar(&projectOwned, "owned", false, "只显示拥有的项目")
	projectListCmd.Flags().BoolVar(&projectArchived, "archived", false, "包含已归档的项目")
	projectListCmd.Flags().StringVar(&projectSearch, "search", "", "搜索项目名称或描述")
	projectListCmd.Flags().StringVar(&projectMatch, "match", "", "使用正则表达式匹配项目路径或名称")
	projectListCmd.Flags().IntVar(&projectLimit, "limit", 20, "限制返回的项目数量")
	projectListCmd.Flags().BoolVar(&projectHasSchedule, "has-schedule", false, "只显示配置了 pipeline schedule 的项目")
	projectListCmd.Flags().BoolVar(&projectScheduleDetail, "schedule-detail", false, "输出 pipeline schedule 的详细信息（需要与 --has-schedule 一起使用）")
	projectListCmd.Flags().BoolVar(&projectQuiet, "quiet", false, "只输出项目名称（PathWithNamespace）")

	// 绑定 project list 标志到 Viper
	viper.BindPFlag("project.owned", projectListCmd.Flags().Lookup("owned"))
	viper.BindPFlag("project.archived", projectListCmd.Flags().Lookup("archived"))
	viper.BindPFlag("project.search", projectListCmd.Flags().Lookup("search"))
	viper.BindPFlag("project.match", projectListCmd.Flags().Lookup("match"))
	viper.BindPFlag("project.limit", projectListCmd.Flags().Lookup("limit"))
	viper.BindPFlag("project.has-schedule", projectListCmd.Flags().Lookup("has-schedule"))
	viper.BindPFlag("project.schedule-detail", projectListCmd.Flags().Lookup("schedule-detail"))
	viper.BindPFlag("project.quiet", projectListCmd.Flags().Lookup("quiet"))

	// project get 标志
	projectGetCmd.Flags().BoolVar(&projectGetDetail, "detail", false, "使用详细格式（带颜色）显示完整的项目数据结构")

	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectGetCmd)

	return projectCmd
}
