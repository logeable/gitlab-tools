package pipeline

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
)

func runGetCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	pipelineIDStr := args[1]

	// 解析 pipeline ID
	pipelineID, err := strconv.Atoi(pipelineIDStr)
	if err != nil {
		return fmt.Errorf("无效的 pipeline ID: %s", pipelineIDStr)
	}

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 获取 pipeline 状态
	pipeline, _, err := client.Pipelines.GetPipeline(projectID, pipelineID)
	if err != nil {
		return fmt.Errorf("获取 pipeline 状态失败: %v", err)
	}

	// 打印 pipeline 信息
	PrintPipelineInfo(pipeline)

	return nil
}
