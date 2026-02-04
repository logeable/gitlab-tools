package mr

import (
	"github.com/spf13/cobra"
)

var (
	mrCreateTitle        string
	mrCreateDescription  string
	mrCreateQuiet        bool
	mrMergeDeleteSource  bool
	mrMergeCommitMessage string
	mrMergeLink          string
	mrListTargetBranch   string
	mrListState          string
	mrListWithPipelines  bool
)

// NewCommand 创建并返回 mr 命令组
func NewCommand() *cobra.Command {
	mrCmd := &cobra.Command{
		Use:   "mr",
		Short: "列出、创建、合并 Merge Request",
		Long:  "列出 MR（默认 opened）、创建 MR、按 IID 合并。",
	}

	mrListCmd := &cobra.Command{
		Use:   "list <项目ID>",
		Short: "列出 MR（默认 opened；可按 state/target-branch 过滤）",
		Long:  "列出项目的 MR；--state 为 opened/closed/merged（默认 opened），--target-branch 按目标分支过滤。",
		Example: `  gitlab-tools mr list 123
  gitlab-tools mr list my-group/my-project
  gitlab-tools mr list my-group/my-project --target-branch feature
  gitlab-tools mr list my-group/my-project --state opened
  gitlab-tools mr list my-group/my-project --state closed
  gitlab-tools mr list my-group/my-project --state merged
  gitlab-tools mr list my-group/my-project --with-pipelines`,
		Args: cobra.ExactArgs(1),
		RunE: runListCmd,
	}

	mrCreateCmd := &cobra.Command{
		Use:   "create <项目ID> <源分支> <目标分支>",
		Short: "从源分支向目标分支创建 MR",
		Long:  "从源分支向目标分支创建 MR。--title/--description 可选，--quiet 只输出链接。",
		Example: `  gitlab-tools mr create 123 feature main
  gitlab-tools mr create my-group/my-project feature main
  gitlab-tools mr create 123 feature main --title "我的功能"
  gitlab-tools mr create 123 feature main --title "我的功能" --description "功能描述"
  gitlab-tools mr create 123 feature main --quiet`,
		Args: cobra.ExactArgs(3),
		RunE: runCreateCmd,
	}

	mrMergeCmd := &cobra.Command{
		Use:   "merge <项目ID> <MR IID>",
		Short: "按项目与 IID 合并 MR（可选合并后删源分支）",
		Long:  "按项目与 MR IID 合并。--delete-source-branch 合并后删除源分支。失败（如已合并、冲突）exit 1。",
		Example: `  gitlab-tools mr merge 123 456
  gitlab-tools mr merge my-group/my-project 456
  gitlab-tools mr merge 123 456 --delete-source-branch
  gitlab-tools mr merge 123 456 --merge-commit-message "合并信息"`,
		RunE: runMergeCmd,
	}

	// mr create 标志
	mrCreateCmd.Flags().StringVar(&mrCreateTitle, "title", "", "指定 Merge Request 的标题")
	mrCreateCmd.Flags().StringVar(&mrCreateDescription, "description", "", "指定 Merge Request 的描述")
	mrCreateCmd.Flags().BoolVar(&mrCreateQuiet, "quiet", false, "quiet 模式：创建 MR 后只显示链接")

	// mr merge 标志
	mrMergeCmd.Flags().BoolVar(&mrMergeDeleteSource, "delete-source-branch", false, "合并后删除源分支")
	mrMergeCmd.Flags().StringVar(&mrMergeCommitMessage, "merge-commit-message", "", "自定义合并提交信息")
	mrMergeCmd.Flags().StringVar(&mrMergeLink, "link", "", "合并链接")

	// mr list 标志
	mrListCmd.Flags().StringVar(&mrListTargetBranch, "target-branch", "", "按目标分支过滤 Merge Request")
	mrListCmd.Flags().StringVar(&mrListState, "state", "", "按状态过滤 Merge Request (opened, closed, merged)")
	mrListCmd.Flags().BoolVar(&mrListWithPipelines, "with-pipelines", false, "显示 Merge Request 的 pipelines")

	mrCmd.AddCommand(mrListCmd)
	mrCmd.AddCommand(mrCreateCmd)
	mrCmd.AddCommand(mrMergeCmd)

	return mrCmd
}
