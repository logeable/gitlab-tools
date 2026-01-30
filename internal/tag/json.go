package tag

import (
	"os"

	"gitlab-tools/internal/output"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// TagItemJSON 用于 tag list --json 的单条输出
type TagItemJSON struct {
	Name    string `json:"name"`
	Commit  string `json:"commit_sha"`
	WebURL  string `json:"web_url"`
}

// CreateResultJSON 用于 tag create --json
type CreateResultJSON struct {
	Name   string `json:"name"`
	Commit string `json:"commit_sha"`
	WebURL string `json:"web_url"`
}

func tagToItemJSON(t *gitlab.Tag) TagItemJSON {
	sha := ""
	if t.Commit != nil {
		sha = t.Commit.ID
		if len(sha) > 8 {
			sha = sha[:8]
		}
	}
	webURL := ""
	// GitLab Tag 可能无 WebURL 字段，留空或由调用方构造
	return TagItemJSON{
		Name:   t.Name,
		Commit: sha,
		WebURL: webURL,
	}
}

// WriteTagListJSON 输出 tag 列表为 JSON
func WriteTagListJSON(tags []*gitlab.Tag) error {
	if tags == nil {
		tags = []*gitlab.Tag{}
	}
	items := make([]TagItemJSON, 0, len(tags))
	for _, t := range tags {
		items = append(items, tagToItemJSON(t))
	}
	return output.WriteJSON(os.Stdout, items)
}

// WriteCreateResultJSON 输出 tag create 结果为 JSON
func WriteCreateResultJSON(t *gitlab.Tag) error {
	sha := ""
	if t.Commit != nil {
		sha = t.Commit.ID
		if len(sha) > 8 {
			sha = sha[:8]
		}
	}
	return output.WriteJSON(os.Stdout, CreateResultJSON{
		Name:   t.Name,
		Commit: sha,
		WebURL: "",
	})
}
