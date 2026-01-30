package project

import (
	"os"

	"gitlab-tools/internal/output"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// ProjectItemJSON 用于 --json 时的单条项目输出（snake_case）
type ProjectItemJSON struct {
	ID               int    `json:"id"`
	Path             string `json:"path"`
	PathWithNamespace string `json:"path_with_namespace"`
	Visibility       string `json:"visibility"`
	WebURL           string `json:"web_url"`
	DefaultBranch    string `json:"default_branch"`
	Archived         bool   `json:"archived"`
}

// ProjectGetJSON 用于 project get --json 的单个项目输出
type ProjectGetJSON struct {
	ID               int    `json:"id"`
	Path             string `json:"path"`
	PathWithNamespace string `json:"path_with_namespace"`
	Visibility       string `json:"visibility"`
	WebURL           string `json:"web_url"`
	DefaultBranch   string `json:"default_branch"`
	Archived         bool   `json:"archived"`
}

func projectToItemJSON(p *gitlab.Project) ProjectItemJSON {
	path := ""
	if p.PathWithNamespace != "" {
		path = p.PathWithNamespace
	} else if p.Path != "" {
		path = p.Path
	}
	return ProjectItemJSON{
		ID:                p.ID,
		Path:              path,
		PathWithNamespace: p.PathWithNamespace,
		Visibility:        string(p.Visibility),
		WebURL:            p.WebURL,
		DefaultBranch:     p.DefaultBranch,
		Archived:          p.Archived,
	}
}

func projectToGetJSON(p *gitlab.Project) ProjectGetJSON {
	path := p.PathWithNamespace
	if path == "" {
		path = p.Path
	}
	return ProjectGetJSON{
		ID:                p.ID,
		Path:              path,
		PathWithNamespace: p.PathWithNamespace,
		WebURL:            p.WebURL,
		DefaultBranch:    p.DefaultBranch,
		Visibility:        string(p.Visibility),
		Archived:          p.Archived,
	}
}

func writeProjectsJSON(projects []*gitlab.Project) error {
	if projects == nil {
		projects = []*gitlab.Project{}
	}
	items := make([]ProjectItemJSON, 0, len(projects))
	for _, p := range projects {
		items = append(items, projectToItemJSON(p))
	}
	return output.WriteJSON(os.Stdout, items)
}
