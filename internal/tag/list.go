package tag

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
)

func runListCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 获取标签列表
	tags, _, err := client.Tags.ListTags(projectID, nil)
	if err != nil {
		return fmt.Errorf("获取项目 %s 的标签列表失败: %v", projectID, err)
	}

	// 打印标签列表
	printTagsList(projectID, tags)

	return nil
}
