package project

import (
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// filterProjectsBySchedule 过滤出有活跃 pipeline schedule 配置的项目，并返回每个项目的 schedules
func filterProjectsBySchedule(client *gitlab.Client, projects []*gitlab.Project) ([]*gitlab.Project, map[string][]*gitlab.PipelineSchedule, error) {
	var filtered []*gitlab.Project
	schedulesMap := make(map[string][]*gitlab.PipelineSchedule)

	for _, project := range projects {
		// 查询项目的 pipeline schedules
		schedules, _, err := client.PipelineSchedules.ListPipelineSchedules(project.PathWithNamespace, nil)
		if err != nil {
			// 如果查询失败（可能是权限问题或项目不存在），跳过该项目
			continue
		}

		// 过滤出活跃的 schedules 并获取详细信息
		var activeSchedules []*gitlab.PipelineSchedule
		for _, schedule := range schedules {
			if !schedule.Active {
				continue
			}

			// 获取 schedule 详细信息
			detail, _, err := client.PipelineSchedules.GetPipelineSchedule(project.PathWithNamespace, schedule.ID)
			if err != nil {
				// 如果获取详情失败，使用基本信息
				activeSchedules = append(activeSchedules, schedule)
				continue
			}
			activeSchedules = append(activeSchedules, detail)
		}

		// 如果有活跃的 pipeline schedule 配置，添加到过滤结果中
		if len(activeSchedules) > 0 {
			filtered = append(filtered, project)
			schedulesMap[project.PathWithNamespace] = activeSchedules
		}
	}

	return filtered, schedulesMap, nil
}
