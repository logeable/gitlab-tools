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
		Short: "Pipeline 管理",
		Long:  "查看和管理 GitLab Pipeline",
	}

	pipelineGetCmd := &cobra.Command{
		Use:   "get <项目ID> <PipelineID>",
		Short: "获取 pipeline 详细信息",
		Long:  "获取指定项目的 pipeline 状态和详细信息",
		Example: `  gitlab-tools pipeline get 123 456
  gitlab-tools pipeline get my-group/my-project 789`,
		Args: cobra.ExactArgs(2),
		RunE: runGetCmd,
	}

	pipelineListCmd := &cobra.Command{
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
		Short: "获取指定分支的最新 pipeline",
		Long:  "获取指定项目的指定分支的最新 pipeline",
		Example: `  gitlab-tools pipeline latest 123 main
  gitlab-tools pipeline latest my-group/my-project develop`,
		Args: cobra.ExactArgs(2),
		RunE: runLatestCmd,
	}

	pipelineCheckScheduleCmd := &cobra.Command{
		Use:   "check-schedule <项目ID>",
		Short: "检查最近的 scheduled pipeline 是否成功",
		Long:  "检查指定项目的最近 scheduled pipeline 是否成功，如果未成功则退出码为 1",
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
