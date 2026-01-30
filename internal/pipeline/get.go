package pipeline

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
	"gitlab-tools/internal/config"
)

func runGetCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	pipelineIDStr := args[1]

	// 解析 pipeline ID（用法错误 → 退出码 2）
	pipelineID, err := strconv.Atoi(pipelineIDStr)
	if err != nil {
		return errors.Join(config.ErrUsage, fmt.Errorf("无效的 pipeline ID: %s", pipelineIDStr))
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

	if config.GetJSON() {
		return WritePipelineGetJSON(pipeline)
	}
	PrintPipelineInfo(pipeline)

	return nil
}
