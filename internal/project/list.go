package project

import (
	"fmt"
	"os"
	"regexp"

	"gitlab-tools/internal/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	hasSchedule := viper.GetBool("project.has-schedule") || projectHasSchedule
	scheduleDetail := viper.GetBool("project.schedule-detail") || projectScheduleDetail
	quiet := viper.GetBool("project.quiet") || projectQuiet

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

	// 如果指定了 schedule-detail 但没有指定 has-schedule，给出提示
	if scheduleDetail && !hasSchedule {
		return fmt.Errorf("--schedule-detail 需要与 --has-schedule 一起使用")
	}

	// 如果指定了 has-schedule 参数，过滤有 pipeline schedule 配置的项目
	var projectSchedules map[string][]*gitlab.PipelineSchedule
	if hasSchedule {
		var filtered []*gitlab.Project
		projectSchedules, err = filterProjectsBySchedule(client, projects, &filtered)
		if err != nil {
			return fmt.Errorf("过滤 pipeline schedule 失败: %v", err)
		}
		projects = filtered
	}

	// 打印项目信息
	printProjectsList(projects, quiet, scheduleDetail, projectSchedules, client)

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

// filterProjectsBySchedule 过滤出有 pipeline schedule 配置的项目，并返回每个项目的 schedules
// 默认只返回活跃的 schedules
func filterProjectsBySchedule(client *gitlab.Client, projects []*gitlab.Project, filtered *[]*gitlab.Project) (map[string][]*gitlab.PipelineSchedule, error) {
	schedulesMap := make(map[string][]*gitlab.PipelineSchedule)

	for _, project := range projects {
		// 查询项目的 pipeline schedules
		schedules, _, err := client.PipelineSchedules.ListPipelineSchedules(project.PathWithNamespace, nil)
		if err != nil {
			// 如果查询失败（可能是权限问题或项目不存在），跳过该项目
			continue
		}
		// 过滤出活跃的 schedules
		var activeSchedules []*gitlab.PipelineSchedule
		for _, schedule := range schedules {
			if schedule.Active {
				activeSchedules = append(activeSchedules, schedule)
			}

			// get schedule detail
			schedule, _, err := client.PipelineSchedules.GetPipelineSchedule(project.PathWithNamespace, schedule.ID)
			if err != nil {
				continue
			}
			activeSchedules = append(activeSchedules, schedule)
		}
		// 如果有活跃的 pipeline schedule 配置，添加到过滤结果中
		if len(activeSchedules) > 0 {
			*filtered = append(*filtered, project)
			schedulesMap[project.PathWithNamespace] = activeSchedules
		}
	}

	return schedulesMap, nil
}
