package tag

import (
	"fmt"
	"strings"

	"gitlab-tools/internal/output"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func printTagsList(projectID string, tags []*gitlab.Tag) {
	fmt.Printf("项目: %s\n", projectID)
	if len(tags) == 0 {
		fmt.Println("  未找到标签")
		return
	}

	fmt.Printf("  找到 %d 个标签:\n\n", len(tags))
	for i, tag := range tags {
		fmt.Printf("  [%d] %s\n", i+1, tag.Name)
		if tag.Commit != nil {
			sha := tag.Commit.ID
			if len(sha) > 8 {
				sha = sha[:8]
			}
			fmt.Printf("      提交: %s\n", sha)
			if tag.Commit.Message != "" {
				message := strings.Split(tag.Commit.Message, "\n")[0]
				if len(message) > 60 {
					message = message[:60] + "..."
				}
				fmt.Printf("      提交信息: %s\n", message)
			}
			if tag.Commit.CommittedDate != nil {
				fmt.Printf("      提交时间: %s\n", output.FormatToLocalTime(tag.Commit.CommittedDate))
			}
			if tag.Commit.AuthorName != "" {
				fmt.Printf("      提交者: %s\n", tag.Commit.AuthorName)
			}
		}
		if tag.Message != "" {
			fmt.Printf("      标签消息: %s\n", tag.Message)
		}
		if tag.Release != nil && tag.Release.Description != "" {
			fmt.Printf("      发布说明: %s\n", tag.Release.Description)
		}
		fmt.Println()
	}
}

func printTagInfo(tag *gitlab.Tag) {
	fmt.Println()
	fmt.Println("标签已创建:")
	fmt.Printf("  标签名: %s\n", tag.Name)
	if tag.Commit != nil {
		fmt.Printf("  提交: %s\n", tag.Commit.ID)
		if tag.Commit.Message != "" {
			message := strings.Split(tag.Commit.Message, "\n")[0]
			if len(message) > 80 {
				message = message[:80] + "..."
			}
			fmt.Printf("  提交信息: %s\n", message)
		}
		if tag.Commit.CommittedDate != nil {
			fmt.Printf("  提交时间: %s\n", output.FormatToLocalTime(tag.Commit.CommittedDate))
		}
		if tag.Commit.AuthorName != "" {
			fmt.Printf("  提交者: %s\n", tag.Commit.AuthorName)
		}
	}
	if tag.Message != "" {
		fmt.Printf("  标签消息: %s\n", tag.Message)
	}
}
