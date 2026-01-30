package pipeline

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	pipelineListLimit  int
	pipelineListStatus string
	pipelineListRef    string
)

// NewCommand 创建并返回 pipeline 命令组
func NewCommand() *cobra.Command {
	pipelineCmd := &cobra.Command{
		Use:   "pipeline",
		Short: "查看流水线状态、列表、最新与定时任务",
		Long:  "查看单条流水线、列出一项目的流水线、取某 ref 最新一条、检查定时流水线是否成功。",
	}

	pipelineGetCmd := &cobra.Command{
		Use:   "get <项目ID> <PipelineID>",
		Short: "按项目与 Pipeline ID 获取单条流水线详情",
		Long:  "根据项目与流水线 ID 获取单条流水线详情。",
		Example: `  gitlab-tools pipeline get 123 456
  gitlab-tools pipeline get my-group/my-project 789`,
		Args: cobra.ExactArgs(2),
		RunE: runGetCmd,
	}

	pipelineListCmd := &cobra.Command{
		Use:   "list <项目ID>...",
		Short: "列出指定项目的流水线（可按状态/ref/条数过滤）",
		Long:  "列出项目的流水线，默认按更新时间降序。--status/--ref/--limit 过滤。",
		Example: `  gitlab-tools pipeline list 123
  gitlab-tools pipeline list 123
  gitlab-tools pipeline list my-group/my-project
  gitlab-tools pipeline list 123 --limit 10
  gitlab-tools pipeline list 123 --status success
  gitlab-tools pipeline list 123 --ref main`,
		Args: cobra.ExactArgs(1),
		RunE: runListCmd,
	}

	// pipeline list 标志
	pipelineListCmd.Flags().IntVar(&pipelineListLimit, "limit", 5, "每个项目显示的 pipeline 数量")
	pipelineListCmd.Flags().StringVar(&pipelineListStatus, "status", "", "按状态过滤 Pipeline (running, pending, success, failed, canceled, skipped, created, manual)")
	pipelineListCmd.Flags().StringVar(&pipelineListRef, "ref", "", "按 ref 过滤 Pipeline")

	// 绑定 pipeline list 标志到 Viper
	viper.BindPFlag("pipeline.list.limit", pipelineListCmd.Flags().Lookup("limit"))
	viper.BindPFlag("pipeline.list.status", pipelineListCmd.Flags().Lookup("status"))

	pipelineLatestCmd := &cobra.Command{
		Use:   "latest <项目ID> <分支名>",
		Short: "获取某 ref 的最新一条流水线",
		Long:  "获取指定项目、指定 ref 的最新一条流水线。",
		Example: `  gitlab-tools pipeline latest 123 main
  gitlab-tools pipeline latest my-group/my-project develop`,
		Args: cobra.ExactArgs(2),
		RunE: runLatestCmd,
	}

	pipelineCheckScheduleCmd := &cobra.Command{
		Use:   "check-schedule <项目ID>",
		Short: "检查定时流水线最近一次是否成功（失败则 exit 1）",
		Long:  "检查指定项目最近一次定时流水线是否成功。成功 exit 0，未成功或未配置 exit 1。",
		Example: `  gitlab-tools pipeline check-schedule 123
  gitlab-tools pipeline check-schedule my-group/my-project`,
		Args: cobra.ExactArgs(1),
		RunE: runCheckScheduleCmd,
	}

	pipelineCmd.AddCommand(pipelineGetCmd)
	pipelineCmd.AddCommand(pipelineListCmd)
	pipelineCmd.AddCommand(pipelineLatestCmd)
	pipelineCmd.AddCommand(pipelineCheckScheduleCmd)

	return pipelineCmd
}
