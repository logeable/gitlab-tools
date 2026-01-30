package pipeline

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"gitlab-tools/internal/client"
	"gitlab-tools/internal/config"
)

func runListCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 从 Viper 获取 limit（支持从环境变量和命令行参数读取）
	// 优先级：命令行参数 > 环境变量 > 默认值
	limit := viper.GetInt("pipeline.list.limit")
	if limit <= 0 {
		// 如果 viper 返回 0 或负数，使用命令行参数的默认值或变量值
		if pipelineListLimit > 0 {
			limit = pipelineListLimit
		} else {
			limit = 5 // 最终默认值
		}
	}

	// 构建查询选项
	opt := &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: limit,
			Page:    1,
		},
		OrderBy: gitlab.Ptr("updated_at"),
		Sort:    gitlab.Ptr("desc"),
	}

	// 如果指定了 --status 参数，设置状态过滤
	status := viper.GetString("pipeline.list.status")
	if status == "" {
		status = pipelineListStatus
	}
	if status != "" {
		// 验证状态值（GitLab API 支持的状态值）
		validStatuses := map[string]bool{
			"running":  true,
			"pending":  true,
			"success":  true,
			"failed":   true,
			"canceled": true,
			"skipped":  true,
			"created":  true,
			"manual":   true,
		}
		if !validStatuses[status] {
			return errors.Join(config.ErrUsage, fmt.Errorf("无效的状态值: %s。支持的状态值: running, pending, success, failed, canceled, skipped, created, manual", status))
		}
		statusValue := gitlab.BuildStateValue(status)
		opt.Status = &statusValue
	}

	if pipelineListRef != "" {
		opt.Ref = &pipelineListRef
	}

	// 获取 pipeline 列表
	pipelines, _, err := client.Pipelines.ListProjectPipelines(projectID, opt)
	if err != nil {
		return fmt.Errorf("获取项目 %s 的 pipeline 列表失败: %v", projectID, err)
	}

	if config.GetJSON() {
		// 获取完整 pipeline 信息用于 JSON
		fullPipelines := make([]*gitlab.Pipeline, 0, len(pipelines))
		for _, p := range pipelines {
			full, _, err := client.Pipelines.GetPipeline(projectID, p.ID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "获取 pipeline %d 失败: %v\n", p.ID, err)
				continue
			}
			fullPipelines = append(fullPipelines, full)
		}
		return WritePipelineListJSON(projectID, fullPipelines)
	}

	fmt.Printf("项目: %s\n", projectID)
	fmt.Printf("  找到 %d 条 pipeline:\n\n", len(pipelines))
	for i, pipeline := range pipelines {
		pipeline, _, err := client.Pipelines.GetPipeline(projectID, pipeline.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取项目 %s 的 pipeline %d 失败: %v\n", projectID, pipeline.ID, err)
			continue
		}
		if i > 0 {
			fmt.Println()
		}
		PrintPipelineInfo(pipeline)
	}

	return nil
}
