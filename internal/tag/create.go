package tag

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
	"gitlab-tools/internal/config"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func runCreateCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	tagName := args[1]

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
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

	if config.GetJSON() {
		return WriteCreateResultJSON(tag)
	}
	// 打印标签信息
	printTagInfo(tag)

	return nil
}
