package mr

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
	"gitlab-tools/internal/config"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func runCreateCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	sourceBranch := args[1]
	targetBranch := args[2]

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 检查分支之间是否有差异
	compareOpt := &gitlab.CompareOptions{
		From: gitlab.Ptr(targetBranch),
		To:   gitlab.Ptr(sourceBranch),
	}

	compare, _, err := client.Repositories.Compare(projectID, compareOpt)
	if err != nil {
		return fmt.Errorf("获取分支差异失败: %v", err)
	}

	// 检查是否有差异：如果没有文件差异，则不创建 MR
	if len(compare.Diffs) == 0 {
		fmt.Println("提示: 两个分支之间没有文件差异，跳过创建 Merge Request")
		return nil
	}

	// 检查是否已经存在相同源分支和目标分支的开放 MR
	existingMRs, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, &gitlab.ListProjectMergeRequestsOptions{
		SourceBranch: gitlab.Ptr(sourceBranch),
		TargetBranch: gitlab.Ptr(targetBranch),
		State:        gitlab.Ptr("opened"),
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	})
	if err != nil {
		return fmt.Errorf("检查现有 Merge Request 失败: %v", err)
	}

	// 如果已存在开放的 MR，则提示并跳过创建
	if len(existingMRs) > 0 {
		fmt.Printf("提示: 源分支 %s 到目标分支 %s 已存在开放的 Merge Request:\n", sourceBranch, targetBranch)
		for _, existingMR := range existingMRs {
			fmt.Printf("  !%d: %s\n", existingMR.IID, existingMR.Title)
			if existingMR.WebURL != "" {
				fmt.Printf("      %s\n", existingMR.WebURL)
			}
		}
		fmt.Println("跳过创建新的 Merge Request")
		return nil
	}

	// 获取当前用户信息，用于设置 assignee
	currentUser, _, err := client.Users.CurrentUser()
	if err != nil {
		return fmt.Errorf("获取当前用户信息失败: %v", err)
	}

	// 设置 MR 标题
	mrTitle := mrCreateTitle
	if mrTitle == "" {
		// 如果没有指定标题，使用默认格式
		mrTitle = fmt.Sprintf("Merge %s into %s", sourceBranch, targetBranch)
	}

	// 创建 Merge Request（从 sourceBranch 合并到 targetBranch）
	mrOpt := &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.Ptr(mrTitle),
		SourceBranch: gitlab.Ptr(sourceBranch),
		TargetBranch: gitlab.Ptr(targetBranch),
		AssigneeID:   gitlab.Ptr(currentUser.ID),
	}

	if mrCreateDescription != "" {
		mrOpt.Description = gitlab.Ptr(mrCreateDescription)
	}

	mr, _, err := client.MergeRequests.CreateMergeRequest(projectID, mrOpt)
	if err != nil {
		return fmt.Errorf("创建 Merge Request 失败: %v", err)
	}

	if config.GetJSON() {
		return WriteCreateResultJSON(mr)
	}
	// 打印 Merge Request 信息
	printMergeRequestInfo(mr, mrCreateQuiet)

	return nil
}
