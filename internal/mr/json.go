package mr

import (
	"os"

	"gitlab-tools/internal/output"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// MRItemJSON 用于 mr list --json 的单条输出
type MRItemJSON struct {
	IID           int    `json:"iid"`
	Title         string `json:"title"`
	SourceBranch  string `json:"source_branch"`
	TargetBranch  string `json:"target_branch"`
	State         string `json:"state"`
	WebURL        string `json:"web_url"`
}

// CreateResultJSON 用于 mr create --json
type CreateResultJSON struct {
	IID    int    `json:"iid"`
	WebURL string `json:"web_url"`
}

// MergeResultJSON 用于 mr merge --json
type MergeResultJSON struct {
	IID    int    `json:"iid"`
	WebURL string `json:"web_url"`
}

func mrToItemJSON(m *gitlab.MergeRequest) MRItemJSON {
	return MRItemJSON{
		IID:           m.IID,
		Title:         m.Title,
		SourceBranch:  m.SourceBranch,
		TargetBranch:  m.TargetBranch,
		State:         m.State,
		WebURL:        m.WebURL,
	}
}

// mrBasicToItemJSON 从 List 返回的 BasicMergeRequest 转为 JSON 项
func mrBasicToItemJSON(m *gitlab.BasicMergeRequest) MRItemJSON {
	return MRItemJSON{
		IID:           m.IID,
		Title:         m.Title,
		SourceBranch:  m.SourceBranch,
		TargetBranch:  m.TargetBranch,
		State:         m.State,
		WebURL:        m.WebURL,
	}
}

// WriteMRListJSON 输出 MR 列表为 JSON（接受 List 返回的 BasicMergeRequest）
func WriteMRListJSON(mrs []*gitlab.BasicMergeRequest) error {
	if mrs == nil {
		mrs = []*gitlab.BasicMergeRequest{}
	}
	items := make([]MRItemJSON, 0, len(mrs))
	for _, m := range mrs {
		items = append(items, mrBasicToItemJSON(m))
	}
	return output.WriteJSON(os.Stdout, items)
}

// WriteCreateResultJSON 输出 mr create 结果为 JSON
func WriteCreateResultJSON(m *gitlab.MergeRequest) error {
	return output.WriteJSON(os.Stdout, CreateResultJSON{IID: m.IID, WebURL: m.WebURL})
}

// WriteMergeResultJSON 输出 mr merge 结果为 JSON
func WriteMergeResultJSON(m *gitlab.MergeRequest) error {
	return output.WriteJSON(os.Stdout, MergeResultJSON{IID: m.IID, WebURL: m.WebURL})
}
