package project

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
)

func runGetCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 获取项目信息
	project, _, err := client.Projects.GetProject(projectID, nil)
	if err != nil {
		return fmt.Errorf("获取项目信息失败: %v", err)
	}

	// 打印项目信息
	printProjectInfo(project, projectGetDetail)

	return nil
}
