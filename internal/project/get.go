package project

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
	"gitlab-tools/internal/config"
	"gitlab-tools/internal/output"
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

	// 输出：--json 时输出 JSON，否则人类可读
	if config.GetJSON() {
		return output.WriteJSON(os.Stdout, projectToGetJSON(project))
	}
	printProjectInfo(project, projectGetDetail)

	return nil
}
