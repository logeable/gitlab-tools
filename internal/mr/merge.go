package mr

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"gitlab-tools/internal/client"
	"gitlab-tools/internal/config"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func runMergeCmd(cmd *cobra.Command, args []string) error {
	var projectID string
	var mrIIDStr string

	if mrMergeLink != "" {
		// example: https://gitlab.rd.revolai.tech/revlo/revlo-admin-backend/-/merge_requests/17
		u, err := url.Parse(mrMergeLink)
		if err != nil {
			return errors.Join(config.ErrUsage, fmt.Errorf("无效的合并链接: %s", mrMergeLink))
		}
		p := u.Path
		parts := strings.Split(p, "/")
		if len(parts) != 6 {
			return errors.Join(config.ErrUsage, fmt.Errorf("无效的合并链接: %s", mrMergeLink))
		}
		projectID = strings.Join([]string{parts[1], parts[2]}, "/")
		mrIIDStr = parts[len(parts)-1]
	} else {
		projectID = args[0]
		mrIIDStr = args[1]
	}

	// 解析 MR IID（用法错误 → 退出码 2）
	mrIID, err := strconv.Atoi(mrIIDStr)
	if err != nil {
		return errors.Join(config.ErrUsage, fmt.Errorf("无效的 MR IID: %s", mrIIDStr))
	}

	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 构建合并选项
	mrOpt := &gitlab.AcceptMergeRequestOptions{}

	if mrMergeDeleteSource {
		mrOpt.ShouldRemoveSourceBranch = gitlab.Ptr(true)
	}

	if mrMergeCommitMessage != "" {
		mrOpt.MergeCommitMessage = gitlab.Ptr(mrMergeCommitMessage)
	}

	// 合并 Merge Request
	mr, _, err := client.MergeRequests.AcceptMergeRequest(projectID, mrIID, mrOpt)
	if err != nil {
		return fmt.Errorf("合并 Merge Request 失败: %v", err)
	}

	if config.GetJSON() {
		return WriteMergeResultJSON(mr)
	}
	// 打印合并结果
	printMergeRequestDetails(mr)

	return nil
}
