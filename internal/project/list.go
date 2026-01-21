package project

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab-tools/internal/client"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func runListCmd(cmd *cobra.Command, args []string) error {
	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 从 Viper 获取参数（支持从环境变量读取）
	// 优先级：命令行参数 > 环境变量 > 默认值
	owned := viper.GetBool("project.owned") || projectOwned
	archived := viper.GetBool("project.archived") || projectArchived
	search := viper.GetString("project.search")
	if search == "" {
		search = projectSearch
	}
	match := viper.GetString("project.match")
	if match == "" {
		match = projectMatch
	}
	limit := viper.GetInt("project.limit")
	if limit == 0 {
		limit = projectLimit
	}

	// 构建查询选项
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: limit,
			Page:    1,
		},
	}

	if owned {
		ownedPtr := true
		opt.Owned = &ownedPtr
	}

	if archived {
		archivedPtr := true
		opt.Archived = &archivedPtr
	}

	if search != "" {
		opt.Search = &search
	}

	// 获取项目列表
	projects, _, err := client.Projects.ListProjects(opt)
	if err != nil {
		return fmt.Errorf("获取项目列表失败: %v", err)
	}

	// 如果指定了 match 参数，使用正则表达式过滤项目
	if match != "" {
		projects = filterProjectsByRegex(projects, match)
	}

	// 打印项目信息
	printProjectsList(projects)

	return nil
}

func filterProjectsByRegex(projects []*gitlab.Project, pattern string) []*gitlab.Project {
	// 编译正则表达式
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 无效的正则表达式 '%s': %v\n", pattern, err)
		return projects // 如果正则表达式无效，返回所有项目
	}

	// 过滤项目：匹配路径或名称
	var filtered []*gitlab.Project
	for _, project := range projects {
		// 匹配路径（PathWithNamespace）或名称（Name 或 NameWithNamespace）
		if re.MatchString(project.PathWithNamespace) ||
			re.MatchString(project.Name) ||
			re.MatchString(project.NameWithNamespace) {
			filtered = append(filtered, project)
		}
	}

	return filtered
}
