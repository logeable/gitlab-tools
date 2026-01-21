package branch

import (
	"github.com/spf13/cobra"
)

var (
	branchListSearch    string
	branchListHideEmpty bool
	branchListQuiet     bool
	branchDiffStat      bool
	branchDiffCommits   bool
)

// NewCommand 创建并返回 branch 命令组
func NewCommand() *cobra.Command {
	branchCmd := &cobra.Command{
		Use:   "branch",
		Short: "分支管理",
		Long:  "查看和管理 GitLab 项目分支",
	}

	branchListCmd := &cobra.Command{
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
		RunE: runListCmd,
	}

	branchDiffCmd := &cobra.Command{
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
		RunE: runDiffCmd,
	}

	// branch list 标志
	branchListCmd.Flags().StringVar(&branchListSearch, "search", "", "按分支名过滤（部分匹配，不区分大小写）")
	branchListCmd.Flags().BoolVar(&branchListHideEmpty, "hide-empty", false, "如果没有分支则隐藏该项目")
	branchListCmd.Flags().BoolVar(&branchListQuiet, "quiet", false, "只显示项目名")

	// branch diff 标志
	branchDiffCmd.Flags().BoolVar(&branchDiffStat, "stat", false, "仅显示文件变更统计信息")
	branchDiffCmd.Flags().BoolVar(&branchDiffCommits, "commits", false, "仅显示提交差异列表")

	branchCmd.AddCommand(branchListCmd)
	branchCmd.AddCommand(branchDiffCmd)

	return branchCmd
}
