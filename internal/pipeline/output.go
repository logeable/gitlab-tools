package pipeline

import (
	"fmt"
	"time"

	"gitlab-tools/internal/output"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// PrintPipelineInfo 打印 pipeline 信息
func PrintPipelineInfo(pipeline *gitlab.Pipeline) {
	fmt.Printf("Pipeline 信息:\n")
	fmt.Printf("  ID: %d\n", pipeline.ID)
	fmt.Printf("  状态: %s\n", pipeline.Status)
	fmt.Printf("  引用: %s\n", pipeline.Ref)
	fmt.Printf("  SHA: %s\n", pipeline.SHA)
	fmt.Printf("  创建时间: %s\n", output.FormatToLocalTime(pipeline.CreatedAt))
	fmt.Printf("  更新时间: %s\n", output.FormatToLocalTime(pipeline.UpdatedAt))
	fmt.Printf("  源: %s\n", pipeline.Source)
	fmt.Printf("  是否为 tag: %t\n", pipeline.Tag)
	fmt.Printf("  Web URL: %s\n", pipeline.WebURL)

	if pipeline.Duration > 0 {
		dur := time.Duration(pipeline.Duration) * time.Second
		fmt.Printf("  持续时间: %v\n", dur)
	}
	if pipeline.Coverage != "" {
		fmt.Printf("  覆盖率: %s\n", pipeline.Coverage)
	}
}
