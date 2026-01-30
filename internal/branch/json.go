package branch

import (
	"os"

	"gitlab-tools/internal/output"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// BranchItemJSON 用于 --json 时的单条分支输出
type BranchItemJSON struct {
	Name   string `json:"name"`
	Commit string `json:"commit_sha"`
}

// ProjectBranchesJSON 用于 branch list 无项目时的分组输出
type ProjectBranchesJSON struct {
	ProjectID int              `json:"project_id"`
	Path      string           `json:"path"`
	Branches  []BranchItemJSON `json:"branches"`
}

func branchToJSON(b *gitlab.Branch) BranchItemJSON {
	sha := ""
	if b.Commit != nil {
		sha = b.Commit.ID
		if len(sha) > 8 {
			sha = sha[:8]
		}
	}
	return BranchItemJSON{Name: b.Name, Commit: sha}
}

// WriteBranchListJSON 单项目分支列表 JSON
func WriteBranchListJSON(projectID string, branches []*gitlab.Branch) error {
	if branches == nil {
		branches = []*gitlab.Branch{}
	}
	items := make([]BranchItemJSON, 0, len(branches))
	for _, b := range branches {
		items = append(items, branchToJSON(b))
	}
	return output.WriteJSON(os.Stdout, items)
}

// WriteBranchListAllProjectsJSON 多项目分支列表 JSON（按项目分组）
func WriteBranchListAllProjectsJSON(groups []ProjectBranchesJSON) error {
	return output.WriteJSON(os.Stdout, groups)
}

// DiffResultJSON 用于 branch diff --json
type DiffResultJSON struct {
	ProjectID     string   `json:"project_id"`
	SourceBranch  string   `json:"source_branch"`
	TargetBranch  string   `json:"target_branch"`
	CommitCount   int      `json:"commit_count"`
	Commits       []string `json:"commits,omitempty"`
	FilesAdded    int      `json:"files_added"`
	FilesModified int      `json:"files_modified"`
	FilesDeleted  int      `json:"files_deleted"`
}

func WriteBranchDiffJSON(projectID, sourceBranch, targetBranch string, compare *gitlab.Compare) error {
	commits := make([]string, 0, len(compare.Commits))
	for _, c := range compare.Commits {
		sha := c.ID
		if len(sha) > 8 {
			sha = sha[:8]
		}
		commits = append(commits, sha)
	}
	added, modified, deleted := 0, 0, 0
	for _, d := range compare.Diffs {
		if d.NewFile {
			added++
		} else if d.DeletedFile {
			deleted++
		} else {
			modified++
		}
	}
	res := DiffResultJSON{
		ProjectID:     projectID,
		SourceBranch:  sourceBranch,
		TargetBranch:  targetBranch,
		CommitCount:   len(compare.Commits),
		Commits:       commits,
		FilesAdded:    added,
		FilesModified: modified,
		FilesDeleted:  deleted,
	}
	return output.WriteJSON(os.Stdout, res)
}
