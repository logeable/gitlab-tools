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
		Short: "列出分支、比较两分支差异",
		Long:  "列出分支（可选项目）、比较两分支的提交与文件变更。",
	}

	branchListCmd := &cobra.Command{
		Use:   "list [项目ID]",
		Short: "列出分支（可选项目；不传则列所有项目，可能较慢）",
		Long:  "列出分支；不传项目则列所有可访问项目（可能较慢）。--search 按分支名过滤，--quiet 只输出名称。",
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
		Short: "比较源分支与目标分支的提交与文件变更",
		Long:  "比较源分支与目标分支的提交与文件变更。--stat 仅文件统计，--commits 仅提交列表。",
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
