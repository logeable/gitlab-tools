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
		Short: "列出、创建、删除项目标签",
		Long:  "列出、创建、删除项目标签。创建时 --ref 优先于 --branch，缺省默认 main。",
	}

	tagListCmd := &cobra.Command{
		Use:   "list <项目ID>",
		Short: "列出指定项目的所有标签",
		Long:  "列出指定项目的所有标签。",
		Example: `  gitlab-tools tag list 123
  gitlab-tools tag list my-group/my-project`,
		Args: cobra.ExactArgs(1),
		RunE: runListCmd,
	}

	tagCreateCmd := &cobra.Command{
		Use:   "create <项目ID> <标签名>",
		Short: "在指定 ref 或分支上创建标签（--ref 优先于 --branch）",
		Long:  "在指定 ref 或分支上创建标签；--ref 与 --branch 同时存在时以 --ref 为准，缺省默认 main。",
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
		Short: "删除指定项目的指定标签",
		Long:  "删除指定项目的指定标签。",
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
