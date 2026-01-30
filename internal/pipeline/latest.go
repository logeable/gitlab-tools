package pipeline

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
	"gitlab-tools/internal/config"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func runLatestCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	ref := args[1]

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 构建查询选项，按更新时间降序排列，只取第一条
	opt := &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 1,
			Page:    1,
		},
		OrderBy: gitlab.Ptr("updated_at"),
		Sort:    gitlab.Ptr("desc"),
		Ref:     &ref,
	}

	// 获取 pipeline 列表
	pipelines, _, err := client.Pipelines.ListProjectPipelines(projectID, opt)
	if err != nil {
		return fmt.Errorf("获取项目 %s 分支 %s 的最新 pipeline 失败: %v", projectID, ref, err)
	}

	if len(pipelines) == 0 {
		return fmt.Errorf("项目 %s 的分支 %s 没有找到 pipeline", projectID, ref)
	}

	// 获取完整的 pipeline 信息
	pipeline, _, err := client.Pipelines.GetPipeline(projectID, pipelines[0].ID)
	if err != nil {
		return fmt.Errorf("获取 pipeline %d 详细信息失败: %v", pipelines[0].ID, err)
	}

	if config.GetJSON() {
		return WritePipelineGetJSON(pipeline)
	}
	PrintPipelineInfo(pipeline)

	return nil
}
