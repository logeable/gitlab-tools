package project

import (
	"fmt"

	"github.com/k0kubun/pp/v3"
	"gitlab-tools/internal/output"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func printProjectInfo(project *gitlab.Project, useDetail bool) {
	if useDetail {
		pp.Print(project)
		return
	}

	fmt.Printf("项目信息:\n")
	fmt.Printf("  ID: %d\n", project.ID)
	fmt.Printf("  名称: %s\n", project.Name)
	fmt.Printf("  路径: %s\n", project.PathWithNamespace)
	fmt.Printf("  可见性: %s\n", project.Visibility)
	if project.DefaultBranch != "" {
		fmt.Printf("  默认分支: %s\n", project.DefaultBranch)
	}
	if project.Description != "" {
		fmt.Printf("  描述: %s\n", project.Description)
	}
	fmt.Printf("  Web URL: %s\n", project.WebURL)
	if project.Archived {
		fmt.Printf("  状态: 已归档\n")
	}
	if project.LastActivityAt != nil {
		fmt.Printf("  最后活动: %s\n", output.FormatToLocalTime(project.LastActivityAt))
	}
	if project.CreatedAt != nil {
		fmt.Printf("  创建时间: %s\n", output.FormatToLocalTime(project.CreatedAt))
	}
}

func printProjectsList(projects []*gitlab.Project) {
	if len(projects) == 0 {
		fmt.Println("未找到项目")
		return
	}

	fmt.Printf("找到 %d 个项目:\n\n", len(projects))
	for i, project := range projects {
		fmt.Printf("[%d] %s\n", i+1, project.NameWithNamespace)
		fmt.Printf("    ID: %d\n", project.ID)
		fmt.Printf("    路径: %s\n", project.PathWithNamespace)
		fmt.Printf("    可见性: %s\n", project.Visibility)
		fmt.Printf("    默认分支: %s\n", project.DefaultBranch)
		if project.Description != "" {
			fmt.Printf("    描述: %s\n", project.Description)
		}
		fmt.Printf("    Web URL: %s\n", project.WebURL)
		if project.Archived {
			fmt.Printf("    状态: 已归档\n")
		}
		if project.LastActivityAt != nil {
			fmt.Printf("    最后活动: %s\n", output.FormatToLocalTime(project.LastActivityAt))
		}
		fmt.Println()
	}
}
