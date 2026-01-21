package tag

import (
	"github.com/spf13/cobra"
)

var (
	tagCreateBranch  string
	tagCreateRef     string
	tagCreateMessage string
)

// NewCommand 创建并返回 tag 命令组
func NewCommand() *cobra.Command {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Tag 管理",
		Long:  "查看和管理 GitLab 项目标签",
	}

	tagListCmd := &cobra.Command{
		Use:   "list <项目ID>",
		Short: "列出项目的标签",
		Long:  "列出指定项目的所有标签列表",
		Example: `  gitlab-tools tag list 123
  gitlab-tools tag list my-group/my-project`,
		Args: cobra.ExactArgs(1),
		RunE: runListCmd,
	}

	tagCreateCmd := &cobra.Command{
		Use:   "create <项目ID> <标签名>",
		Short: "创建标签",
		Long:  "在指定项目上创建标签，默认在 main 分支上创建",
		Example: `  gitlab-tools tag create 123 v1.0.0
  gitlab-tools tag create my-group/my-project v1.0.0
  gitlab-tools tag create 123 v1.0.0 --branch develop
  gitlab-tools tag create 123 v1.0.0 --ref abc123
  gitlab-tools tag create 123 v1.0.0 --message "版本 1.0.0"`,
		Args: cobra.ExactArgs(2),
		RunE: runCreateCmd,
	}

	tagDeleteCmd := &cobra.Command{
		Use:   "delete <项目ID> <标签名>",
		Short: "删除标签",
		Long:  "删除指定项目的标签",
		Example: `  gitlab-tools tag delete 123 v1.0.0
  gitlab-tools tag delete my-group/my-project v1.0.0`,
		Args: cobra.ExactArgs(2),
		RunE: runDeleteCmd,
	}

	// tag create 标志
	tagCreateCmd.Flags().StringVar(&tagCreateBranch, "branch", "main", "指定目标分支（默认: main）")
	tagCreateCmd.Flags().StringVar(&tagCreateRef, "ref", "", "指定具体的提交 SHA 或分支名（可选，默认使用分支的最新提交）")
	tagCreateCmd.Flags().StringVar(&tagCreateMessage, "message", "", "指定标签消息（可选）")

	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagCreateCmd)
	tagCmd.AddCommand(tagDeleteCmd)

	return tagCmd
}
