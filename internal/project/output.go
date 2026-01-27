package project

import (
	"fmt"
	"time"

	"gitlab-tools/internal/output"

	"github.com/k0kubun/pp/v3"
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

func printProjectsList(projects []*gitlab.Project, opts listOptions, schedulesMap map[string][]*gitlab.PipelineSchedule, client *gitlab.Client) {
	if len(projects) == 0 {
		if !opts.quiet {
			fmt.Println("未找到项目")
		}
		return
	}

	if opts.quiet {
		printProjectsListQuiet(projects)
		return
	}

	// 正常模式：输出详细信息
	fmt.Printf("找到 %d 个项目:\n\n", len(projects))
	for i, project := range projects {
		printProjectItem(project, i+1)
		if opts.scheduleDetail && schedulesMap != nil {
			printProjectSchedules(project.PathWithNamespace, schedulesMap[project.PathWithNamespace], client)
		}
		fmt.Println()
	}
}

func printProjectsListQuiet(projects []*gitlab.Project) {
	for _, project := range projects {
		fmt.Println(project.PathWithNamespace)
	}
}

func printProjectItem(project *gitlab.Project, index int) {
	fmt.Printf("[%d] %s\n", index, project.NameWithNamespace)
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
}

func printProjectSchedules(projectID string, schedules []*gitlab.PipelineSchedule, client *gitlab.Client) {
	if len(schedules) == 0 {
		return
	}

	fmt.Printf("    Pipeline Schedules (%d):\n", len(schedules))
	for j, schedule := range schedules {
		printScheduleInfo(schedule, j+1)
		if schedule.LastPipeline != nil {
			printScheduleLastPipeline(projectID, schedule.LastPipeline, client)
		}
	}
}

func printScheduleInfo(schedule *gitlab.PipelineSchedule, index int) {
	fmt.Printf("      [%d] %s\n", index, schedule.Description)
	fmt.Printf("        ID: %d\n", schedule.ID)
	fmt.Printf("        引用: %s\n", schedule.Ref)
	fmt.Printf("        定时: %s\n", schedule.Cron)
	if schedule.CronTimezone != "" {
		fmt.Printf("        时区: %s\n", schedule.CronTimezone)
	}
	fmt.Printf("        激活: %t\n", schedule.Active)
	if schedule.NextRunAt != nil {
		fmt.Printf("        下次运行: %s\n", output.FormatToLocalTime(schedule.NextRunAt))
	}
}

func printScheduleLastPipeline(projectID string, lastPipeline *gitlab.LastPipeline, client *gitlab.Client) {
	fmt.Printf("        最后 Pipeline:\n")
	fmt.Printf("          ID: %d\n", lastPipeline.ID)
	fmt.Printf("          状态: %s\n", lastPipeline.Status)

	// 获取完整的 pipeline 信息以显示更多详情
	if client == nil {
		return
	}

	pipeline, _, err := client.Pipelines.GetPipeline(projectID, lastPipeline.ID)
	if err != nil {
		return
	}

	fmt.Printf("          引用: %s\n", pipeline.Ref)
	fmt.Printf("          SHA: %s\n", pipeline.SHA)
	if pipeline.CreatedAt != nil {
		fmt.Printf("          创建时间: %s\n", output.FormatToLocalTime(pipeline.CreatedAt))
	}
	if pipeline.UpdatedAt != nil {
		fmt.Printf("          更新时间: %s\n", output.FormatToLocalTime(pipeline.UpdatedAt))
	}
	if pipeline.WebURL != "" {
		fmt.Printf("          Web URL: %s\n", pipeline.WebURL)
	}
	if pipeline.Duration > 0 {
		dur := time.Duration(pipeline.Duration) * time.Second
		fmt.Printf("          持续时间: %v\n", dur)
	}
}
