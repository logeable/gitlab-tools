package branch

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gitlab-tools/internal/client"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func runListCmd(cmd *cobra.Command, args []string) error {
	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 获取搜索参数
	search := branchListSearch

	// 检查是否提供了项目 ID
	if len(args) > 0 && args[0] != "" {
		// 指定了项目 ID，只获取该项目的分支
		projectID := args[0]
		branches, _, err := client.Branches.ListBranches(projectID, nil)
		if err != nil {
			return fmt.Errorf("获取项目 %s 的分支列表失败: %v", projectID, err)
		}

		// 应用搜索过滤
		if search != "" {
			branches = filterBranchesBySearch(branches, search)
		}

		// 如果启用 --hide-empty 且没有分支，则跳过
		if branchListHideEmpty && len(branches) == 0 {
			return nil
		}

		// 打印分支信息
		printBranchesList(projectID, branches, true, branchListQuiet)
	} else {
		// 未指定项目 ID，获取所有项目的分支
		opt := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    1,
			},
		}

		projects, _, err := client.Projects.ListProjects(opt)
		if err != nil {
			return fmt.Errorf("获取项目列表失败: %v", err)
		}

		if len(projects) == 0 {
			fmt.Println("未找到项目")
			return nil
		}

		// 为每个项目获取分支
		for i, project := range projects {
			branches, _, err := client.Branches.ListBranches(project.PathWithNamespace, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "获取项目 %s 的分支列表失败: %v\n", project.PathWithNamespace, err)
				continue
			}

			// 应用搜索过滤
			if search != "" {
				branches = filterBranchesBySearch(branches, search)
			}

			// 如果启用 --hide-empty 且没有分支，则跳过该项目
			if branchListHideEmpty && len(branches) == 0 {
				continue
			}

			if i > 0 && !branchListQuiet {
				fmt.Println() // 项目之间添加空行（quiet 模式下不需要）
			}

			// 打印分支信息
			printBranchesList(project.PathWithNamespace, branches, false, branchListQuiet)
		}
	}

	return nil
}

func filterBranchesBySearch(branches []*gitlab.Branch, searchTerm string) []*gitlab.Branch {
	searchLower := strings.ToLower(searchTerm)
	var filtered []*gitlab.Branch
	for _, branch := range branches {
		if strings.Contains(strings.ToLower(branch.Name), searchLower) {
			filtered = append(filtered, branch)
		}
	}
	return filtered
}
