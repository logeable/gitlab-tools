package branch

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func runDiffCmd(cmd *cobra.Command, args []string) error {
	projectID := args[0]
	sourceBranch := args[1]
	targetBranch := args[2]

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 构建比较选项
	opt := &gitlab.CompareOptions{
		From: gitlab.Ptr(sourceBranch),
		To:   gitlab.Ptr(targetBranch),
	}

	// 获取分支差异
	compare, _, err := client.Repositories.Compare(projectID, opt)
	if err != nil {
		return fmt.Errorf("获取分支差异失败: %v", err)
	}

	// 打印分支差异信息
	printBranchDiff(projectID, sourceBranch, targetBranch, compare, branchDiffStat, branchDiffCommits)

	return nil
}
