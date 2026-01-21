package branch

import (
	"fmt"
	"strings"

	"gitlab-tools/internal/output"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func printBranchesList(projectID string, branches []*gitlab.Branch, singleProject bool, quiet bool) {
	if quiet {
		// quiet 模式：只显示项目名
		if len(branches) > 0 {
			fmt.Println(projectID)
		}
		return
	}

	if singleProject {
		fmt.Printf("项目: %s\n", projectID)
	} else {
		fmt.Printf("项目: %s\n", projectID)
	}

	if len(branches) == 0 {
		fmt.Println("  未找到分支")
		return
	}

	fmt.Printf("  找到 %d 个分支:\n\n", len(branches))
	for i, branch := range branches {
		fmt.Printf("  [%d] %s", i+1, branch.Name)
		if branch.Protected {
			fmt.Printf(" (受保护)")
		}
		if branch.Default {
			fmt.Printf(" (默认分支)")
		}
		fmt.Println()

		if branch.Commit != nil {
			sha := branch.Commit.ID
			if len(sha) > 8 {
				sha = sha[:8]
			}
			fmt.Printf("      最后提交: %s\n", sha)
			if branch.Commit.Message != "" {
				// 只显示提交信息的第一行
				message := strings.Split(branch.Commit.Message, "\n")[0]
				if len(message) > 60 {
					message = message[:60] + "..."
				}
				fmt.Printf("      提交信息: %s\n", message)
			}
			if branch.Commit.CommittedDate != nil {
				fmt.Printf("      提交时间: %s\n", output.FormatToLocalTime(branch.Commit.CommittedDate))
			}
			if branch.Commit.AuthorName != "" {
				fmt.Printf("      提交者: %s\n", branch.Commit.AuthorName)
			}
			if branch.Commit.AuthorEmail != "" {
				fmt.Printf("      提交者邮箱: %s\n", branch.Commit.AuthorEmail)
			}
		}
		fmt.Println()
	}
}

func printBranchDiff(projectID, sourceBranch, targetBranch string, compare *gitlab.Compare, statOnly, commitsOnly bool) {
	fmt.Printf("项目: %s\n", projectID)
	fmt.Printf("源分支: %s\n", sourceBranch)
	fmt.Printf("目标分支: %s\n", targetBranch)
	fmt.Println()

	// 如果仅显示提交列表
	if commitsOnly {
		if len(compare.Commits) == 0 {
			fmt.Println("无提交差异")
			return
		}
		fmt.Printf("提交差异 (%d 个提交):\n\n", len(compare.Commits))
		for i, commit := range compare.Commits {
			sha := commit.ID
			if len(sha) > 8 {
				sha = sha[:8]
			}
			fmt.Printf("  [%d] %s\n", i+1, sha)
			if commit.AuthorName != "" {
				fmt.Printf("      作者: %s", commit.AuthorName)
				if commit.AuthorEmail != "" {
					fmt.Printf(" <%s>", commit.AuthorEmail)
				}
				fmt.Println()
			}
			if commit.Message != "" {
				message := strings.Split(commit.Message, "\n")[0]
				if len(message) > 80 {
					message = message[:80] + "..."
				}
				fmt.Printf("      提交信息: %s\n", message)
			}
			if commit.CommittedDate != nil {
				fmt.Printf("      提交时间: %s\n", output.FormatToLocalTime(commit.CommittedDate))
			}
			fmt.Println()
		}
		return
	}

	// 如果仅显示统计信息
	if statOnly {
		added := 0
		modified := 0
		deleted := 0
		renamed := 0

		for _, diff := range compare.Diffs {
			if diff.NewFile {
				added++
			} else if diff.DeletedFile {
				deleted++
			} else if diff.RenamedFile {
				renamed++
			} else {
				modified++
			}
		}

		fmt.Printf("文件变更统计:\n")
		fmt.Printf("  新增: %d\n", added)
		fmt.Printf("  修改: %d\n", modified)
		fmt.Printf("  删除: %d\n", deleted)
		if renamed > 0 {
			fmt.Printf("  重命名: %d\n", renamed)
		}
		fmt.Printf("  总计: %d\n", len(compare.Diffs))
		return
	}

	// 显示完整信息
	// 提交差异
	if len(compare.Commits) == 0 {
		fmt.Println("无提交差异")
	} else {
		fmt.Printf("提交差异 (%d 个提交):\n\n", len(compare.Commits))
		for i, commit := range compare.Commits {
			sha := commit.ID
			if len(sha) > 8 {
				sha = sha[:8]
			}
			fmt.Printf("  [%d] %s\n", i+1, sha)
			if commit.AuthorName != "" {
				fmt.Printf("      作者: %s", commit.AuthorName)
				if commit.AuthorEmail != "" {
					fmt.Printf(" <%s>", commit.AuthorEmail)
				}
				fmt.Println()
			}
			if commit.Message != "" {
				message := strings.Split(commit.Message, "\n")[0]
				if len(message) > 80 {
					message = message[:80] + "..."
				}
				fmt.Printf("      提交信息: %s\n", message)
			}
			if commit.CommittedDate != nil {
				fmt.Printf("      提交时间: %s\n", output.FormatToLocalTime(commit.CommittedDate))
			}
			fmt.Println()
		}
	}

	// 文件变更统计
	added := 0
	modified := 0
	deleted := 0
	renamed := 0

	for _, diff := range compare.Diffs {
		if diff.NewFile {
			added++
		} else if diff.DeletedFile {
			deleted++
		} else if diff.RenamedFile {
			renamed++
		} else {
			modified++
		}
	}

	fmt.Printf("文件变更统计:\n")
	fmt.Printf("  新增: %d\n", added)
	fmt.Printf("  修改: %d\n", modified)
	fmt.Printf("  删除: %d\n", deleted)
	if renamed > 0 {
		fmt.Printf("  重命名: %d\n", renamed)
	}
	fmt.Printf("  总计: %d\n", len(compare.Diffs))
	fmt.Println()

	// 详细文件差异
	if len(compare.Diffs) > 0 {
		fmt.Printf("文件变更详情:\n\n")
		for i, diff := range compare.Diffs {
			changeType := "修改"
			if diff.NewFile {
				changeType = "新增"
			} else if diff.DeletedFile {
				changeType = "删除"
			} else if diff.RenamedFile {
				changeType = "重命名"
			}

			fmt.Printf("  [%d] %s: %s\n", i+1, changeType, diff.NewPath)
			if diff.RenamedFile && diff.OldPath != diff.NewPath {
				fmt.Printf("      原路径: %s\n", diff.OldPath)
			}
		}
	}
}
