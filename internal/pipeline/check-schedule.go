package pipeline

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
	"gitlab-tools/internal/output"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func runCheckScheduleCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 构建查询选项，只查询 scheduled pipelines，按更新时间降序排列，只取第一条
	source := "schedule"
	opt := &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 1,
			Page:    1,
		},
		OrderBy: gitlab.Ptr("updated_at"),
		Sort:    gitlab.Ptr("desc"),
		Source:  &source,
	}

	// 获取 scheduled pipeline 列表
	pipelines, _, err := client.Pipelines.ListProjectPipelines(projectID, opt)
	if err != nil {
		return fmt.Errorf("获取项目 %s 的 scheduled pipeline 失败: %v", projectID, err)
	}

	if len(pipelines) == 0 {
		fmt.Fprintf(os.Stderr, "项目 %s 没有找到 scheduled pipeline\n", projectID)
		os.Exit(1)
		return nil
	}

	// 获取完整的 pipeline 信息
	pipeline, _, err := client.Pipelines.GetPipeline(projectID, pipelines[0].ID)
	if err != nil {
		return fmt.Errorf("获取 pipeline %d 详细信息失败: %v", pipelines[0].ID, err)
	}

	// 检查状态
	if pipeline.Status == "success" {
		fmt.Printf("✓ 项目 %s 最近的 scheduled pipeline 成功\n", projectID)
		fmt.Printf("  Pipeline ID: %d\n", pipeline.ID)
		fmt.Printf("  状态: %s\n", pipeline.Status)
		fmt.Printf("  引用: %s\n", pipeline.Ref)
		fmt.Printf("  SHA: %s\n", pipeline.SHA)
		fmt.Printf("  更新时间: %s\n", output.FormatToLocalTime(pipeline.UpdatedAt))
		fmt.Printf("  Web URL: %s\n", pipeline.WebURL)
		return nil
	} else {
		fmt.Fprintf(os.Stderr, "✗ 项目 %s 最近的 scheduled pipeline 未成功\n", projectID)
		fmt.Fprintf(os.Stderr, "  Pipeline ID: %d\n", pipeline.ID)
		fmt.Fprintf(os.Stderr, "  状态: %s\n", pipeline.Status)
		fmt.Fprintf(os.Stderr, "  引用: %s\n", pipeline.Ref)
		fmt.Fprintf(os.Stderr, "  SHA: %s\n", pipeline.SHA)
		fmt.Fprintf(os.Stderr, "  更新时间: %s\n", output.FormatToLocalTime(pipeline.UpdatedAt))
		fmt.Fprintf(os.Stderr, "  Web URL: %s\n", pipeline.WebURL)
		os.Exit(1)
		return nil
	}
}
