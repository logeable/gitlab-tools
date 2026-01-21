package tag

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
)

func runDeleteCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	tagName := args[1]

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 删除标签
	_, err = client.Tags.DeleteTag(projectID, tagName)
	if err != nil {
		return fmt.Errorf("删除标签失败: %v", err)
	}

	fmt.Printf("标签 %s 已删除\n", tagName)

	return nil
}
