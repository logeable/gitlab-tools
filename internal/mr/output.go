package mr

import (
	"fmt"

	"gitlab-tools/internal/output"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func printMergeRequestInfo(mr *gitlab.MergeRequest, quiet bool) {
	if quiet {
		// quiet 模式：只显示链接
		if mr.WebURL != "" {
			fmt.Println(mr.WebURL)
		}
		return
	}

	// 正常模式：显示完整信息
	fmt.Println()
	fmt.Println("Merge Request 已创建:")
	fmt.Printf("  ID: %d\n", mr.IID)
	fmt.Printf("  标题: %s\n", mr.Title)
	if mr.Description != "" {
		fmt.Printf("  描述: %s\n", mr.Description)
	}
	fmt.Printf("  源分支: %s\n", mr.SourceBranch)
	fmt.Printf("  目标分支: %s\n", mr.TargetBranch)
	fmt.Printf("  状态: %s\n", mr.State)
	if mr.WebURL != "" {
		fmt.Printf("  Web URL: %s\n", mr.WebURL)
	}
	if mr.CreatedAt != nil {
		fmt.Printf("  创建时间: %s\n", output.FormatToLocalTime(mr.CreatedAt))
	}
}

func printMergeRequestDetails(mr *gitlab.MergeRequest) {
	fmt.Println("Merge Request 已合并:")
	fmt.Printf("  ID: %d\n", mr.IID)
	fmt.Printf("  标题: %s\n", mr.Title)
	if mr.Description != "" {
		fmt.Printf("  描述: %s\n", mr.Description)
	}
	fmt.Printf("  源分支: %s\n", mr.SourceBranch)
	fmt.Printf("  目标分支: %s\n", mr.TargetBranch)
	fmt.Printf("  状态: %s\n", mr.State)
	if mr.MergedAt != nil {
		fmt.Printf("  合并时间: %s\n", output.FormatToLocalTime(mr.MergedAt))
	}
	if mr.WebURL != "" {
		fmt.Printf("  Web URL: %s\n", mr.WebURL)
	}
}
