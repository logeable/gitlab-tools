package tag

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
	"gitlab-tools/internal/config"
	"gitlab-tools/internal/output"
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

	if config.GetJSON() {
		return output.WriteJSON(os.Stdout, map[string]string{"deleted": tagName})
	}
	fmt.Printf("标签 %s 已删除\n", tagName)

	return nil
}
