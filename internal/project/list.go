package project

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"gitlab-tools/internal/client"
	"gitlab-tools/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// listOptions 封装项目列表的选项
type listOptions struct {
	owned          bool
	archived       bool
	search         string
	match          string
	limit          int
	hasSchedule    bool
	scheduleDetail bool
	quiet          bool
}

// getListOptions 从 Viper 和命令行参数获取列表选项
func getListOptions() listOptions {
	return listOptions{
		owned:          viper.GetBool("project.owned") || projectOwned,
		archived:       viper.GetBool("project.archived") || projectArchived,
		search:         getStringOption("project.search", projectSearch),
		match:          getStringOption("project.match", projectMatch),
		limit:          getIntOption("project.limit", projectLimit),
		hasSchedule:    viper.GetBool("project.has-schedule") || projectHasSchedule,
		scheduleDetail: viper.GetBool("project.schedule-detail") || projectScheduleDetail,
		quiet:          viper.GetBool("project.quiet") || projectQuiet,
	}
}

func getStringOption(key string, fallback string) string {
	if val := viper.GetString(key); val != "" {
		return val
	}
	return fallback
}

func getIntOption(key string, fallback int) int {
	if val := viper.GetInt(key); val != 0 {
		return val
	}
	return fallback
}

func runListCmd(cmd *cobra.Command, args []string) error {
	// 创建 GitLab 客户端
	client, err := client.NewClient()
	if err != nil {
		return err
	}

	// 获取选项
	opts := getListOptions()

	// 验证选项（用法错误 → 退出码 2）
	if opts.scheduleDetail && !opts.hasSchedule {
		return errors.Join(config.ErrUsage, fmt.Errorf("--schedule-detail 需要与 --has-schedule 一起使用"))
	}

	// 构建查询选项
	opt := buildListProjectsOptions(opts)

	// 获取项目列表
	projects, _, err := client.Projects.ListProjects(opt)
	if err != nil {
		return fmt.Errorf("获取项目列表失败: %v", err)
	}

	// 应用过滤并获取 schedules（如果需要）
	projects, schedulesMap := applyFilters(client, projects, opts)

	// 输出：--json 时输出 JSON，否则人类可读
	if config.GetJSON() {
		return writeProjectsJSON(projects)
	}
	printProjectsList(projects, opts, schedulesMap, client)

	return nil
}

// buildListProjectsOptions 构建 GitLab API 查询选项
func buildListProjectsOptions(opts listOptions) *gitlab.ListProjectsOptions {
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: opts.limit,
			Page:    1,
		},
	}

	if opts.owned {
		opt.Owned = gitlab.Ptr(true)
	}

	if opts.archived {
		opt.Archived = gitlab.Ptr(true)
	}

	if opts.search != "" {
		opt.Search = gitlab.Ptr(opts.search)
	}

	return opt
}

// applyFilters 应用所有过滤条件，返回过滤后的项目和 schedules map
func applyFilters(client *gitlab.Client, projects []*gitlab.Project, opts listOptions) ([]*gitlab.Project, map[string][]*gitlab.PipelineSchedule) {
	// 正则表达式过滤
	if opts.match != "" {
		projects = filterProjectsByRegex(projects, opts.match)
	}

	// Schedule 过滤（如果需要）
	var schedulesMap map[string][]*gitlab.PipelineSchedule
	if opts.hasSchedule {
		filtered, schedules, err := filterProjectsBySchedule(client, projects)
		if err == nil {
			projects = filtered
			schedulesMap = schedules
		}
		// 如果过滤失败，继续使用原始列表（不中断执行）
	}

	return projects, schedulesMap
}

func filterProjectsByRegex(projects []*gitlab.Project, pattern string) []*gitlab.Project {
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 无效的正则表达式 '%s': %v\n", pattern, err)
		return projects
	}

	var filtered []*gitlab.Project
	for _, project := range projects {
		if re.MatchString(project.PathWithNamespace) ||
			re.MatchString(project.Name) ||
			re.MatchString(project.NameWithNamespace) {
			filtered = append(filtered, project)
		}
	}

	return filtered
}
