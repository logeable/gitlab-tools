package mr

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
	"gitlab-tools/internal/output"
	"gitlab-tools/internal/pipeline"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func runListCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	var targetBranch *string
	if mrListTargetBranch != "" {
		targetBranch = &mrListTargetBranch
	}

	var state *string
	if mrListState != "" {
		state = &mrListState
	}

	opt := &gitlab.ListProjectMergeRequestsOptions{
		State: state,
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		TargetBranch: targetBranch,
	}

	// 获取 Merge Request 列表
	mrs, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, opt)
	if err != nil {
		return fmt.Errorf("获取项目 %s 的 Merge Request 列表失败: %v", projectID, err)
	}

	// 打印 Merge Request 列表
	fmt.Printf("项目: %s\n", projectID)
	if len(mrs) == 0 {
		fmt.Println("  未找到 Merge Request")
		return nil
	}

	fmt.Printf("  找到 %d 个 Merge Request:\n\n", len(mrs))
	for i, mr := range mrs {
		fmt.Printf("  [%d] !%d: %s\n", i+1, mr.IID, mr.Title)
		fmt.Printf("      源分支: %s -> 目标分支: %s\n", mr.SourceBranch, mr.TargetBranch)
		fmt.Printf("      状态: %s\n", mr.State)
		fmt.Printf("      合并状态: %s\n", mr.DetailedMergeStatus)
		if mr.Author != nil {
			fmt.Printf("      创建者: %s", mr.Author.Name)
			if mr.Author.Username != "" {
				fmt.Printf(" (@%s)", mr.Author.Username)
			}
			fmt.Println()
		}
		if mr.CreatedAt != nil {
			fmt.Printf("      创建时间: %s\n", output.FormatToLocalTime(mr.CreatedAt))
		}
		if mr.WebURL != "" {
			fmt.Printf("      Web URL: %s\n", mr.WebURL)
		}

		if mrListWithPipelines {
			pipelines, _, err := client.MergeRequests.ListMergeRequestPipelines(projectID, mr.IID)
			if err != nil {
				return fmt.Errorf("获取 Merge Request %d 的 pipelines 失败: %v", mr.IID, err)
			}
			for _, p := range pipelines {
				p, _, err := client.Pipelines.GetPipeline(projectID, p.ID)
				if err != nil {
					return fmt.Errorf("获取 pipeline %d 失败: %v", p.ID, err)
				}
				pipeline.PrintPipelineInfo(p)
				fmt.Println()
			}
		}
		fmt.Println()
	}

	return nil
}
