package pipeline

import (
	"os"
	"time"

	"gitlab-tools/internal/output"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// PipelineItemJSON 用于 --json 时的 pipeline 输出（snake_case）
type PipelineItemJSON struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	Ref       string    `json:"ref"`
	SHA       string    `json:"sha"`
	WebURL    string    `json:"web_url"`
	CreatedAt time.Time `json:"created_at"`
}

func pipelineToJSON(p *gitlab.Pipeline) PipelineItemJSON {
	createdAt := time.Time{}
	if p.CreatedAt != nil {
		createdAt = *p.CreatedAt
	}
	return PipelineItemJSON{
		ID:        p.ID,
		Status:    p.Status,
		Ref:       p.Ref,
		SHA:       p.SHA,
		WebURL:    p.WebURL,
		CreatedAt: createdAt,
	}
}

// WritePipelineListJSON 输出 pipeline 列表为 JSON（单项目）
func WritePipelineListJSON(projectID string, pipelines []*gitlab.Pipeline) error {
	if pipelines == nil {
		pipelines = []*gitlab.Pipeline{}
	}
	items := make([]PipelineItemJSON, 0, len(pipelines))
	for _, p := range pipelines {
		items = append(items, pipelineToJSON(p))
	}
	return output.WriteJSON(os.Stdout, items)
}

// WritePipelineGetJSON 输出单条 pipeline 为 JSON
func WritePipelineGetJSON(p *gitlab.Pipeline) error {
	return output.WriteJSON(os.Stdout, pipelineToJSON(p))
}

// CheckScheduleResultJSON 用于 pipeline check-schedule --json
type CheckScheduleResultJSON struct {
	Success   bool   `json:"success"`
	ProjectID string `json:"project_id"`
	Message   string `json:"message,omitempty"`
	PipelineID int   `json:"pipeline_id,omitempty"`
	Status    string `json:"status,omitempty"`
	WebURL    string `json:"web_url,omitempty"`
}
