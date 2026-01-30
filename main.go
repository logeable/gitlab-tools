package main

import (
	"errors"
	"fmt"
	"os"

	"gitlab-tools/internal/branch"
	"gitlab-tools/internal/config"
	"gitlab-tools/internal/mr"
	"gitlab-tools/internal/output"
	"gitlab-tools/internal/pipeline"
	"gitlab-tools/internal/project"
	"gitlab-tools/internal/tag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	gitlabURL   string
	gitlabToken string
	jsonOutput  bool
)

var rootCmd = &cobra.Command{
	Use:   "gitlab-tools",
	Short: "与 GitLab API 交互的 CLI 工具集",
	Long:  "管理项目、Pipeline、分支、MR 和标签。支持 --json 输出。退出码：0 成功，1 业务/API 错误，2 用法错误。",
}

func init() {
	// 初始化 Viper
	config.Init()

	// 全局标志
	rootCmd.PersistentFlags().StringVar(&gitlabURL, "url", "", "GitLab 服务器 URL (默认: https://gitlab.com，也可通过配置文件或 GITLAB_URL 环境变量设置)")
	rootCmd.PersistentFlags().StringVar(&gitlabToken, "token", "", "GitLab 访问令牌 (可通过配置文件或 GITLAB_TOKEN 环境变量设置)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "以 JSON 格式输出结果，便于脚本与 Agent 解析")

	// 将 Cobra flags 绑定到 Viper（自动从环境变量读取）
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))

	// 添加子命令
	rootCmd.AddCommand(pipeline.NewCommand())
	rootCmd.AddCommand(project.NewCommand())
	rootCmd.AddCommand(branch.NewCommand())
	rootCmd.AddCommand(mr.NewCommand())
	rootCmd.AddCommand(tag.NewCommand())
	rootCmd.AddCommand(newCapabilitiesCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		exitCode := 1
		if errors.Is(err, config.ErrUsage) {
			exitCode = 2
		}
		if config.GetJSON() {
			_ = output.WriteJSONError(os.Stderr, err.Error(), exitCode)
		} else {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		}
		os.Exit(exitCode)
	}
}

// capabilities 子命令：输出原子能力列表，便于 Agent 发现
func newCapabilitiesCmd() *cobra.Command {
	type cap struct {
		Domain string `json:"domain"`
		Cmd    string `json:"cmd"`
		Use    string `json:"use"`
		Args   string `json:"args"`
	}
	caps := []cap{
		{"project", "list", "列出项目", "[--owned] [--search] [--match] [--limit] [--quiet]"},
		{"project", "get", "获取单项目详情", "<项目ID或路径>"},
		{"pipeline", "list", "列出 Pipeline", "<项目> [--status] [--ref] [--limit]"},
		{"pipeline", "get", "获取单条 Pipeline", "<项目> <PipelineID>"},
		{"pipeline", "latest", "指定 ref 的最新 Pipeline", "<项目> <ref>"},
		{"pipeline", "check-schedule", "检查 Scheduled Pipeline", "<项目>"},
		{"branch", "list", "列出分支", "[项目] [--search] [--hide-empty] [--quiet]"},
		{"branch", "diff", "比较两分支差异", "<项目> <源分支> <目标分支> [--stat] [--commits]"},
		{"mr", "list", "列出 MR", "<项目> [--state] [--target-branch]"},
		{"mr", "create", "创建 MR", "<项目> <源分支> <目标分支> [--title] [--description]"},
		{"mr", "merge", "合并 MR", "<项目> <MR IID> [--delete-source-branch]"},
		{"tag", "list", "列出标签", "<项目>"},
		{"tag", "create", "创建标签", "<项目> <标签名> [--branch] [--ref] [--message]"},
		{"tag", "delete", "删除标签", "<项目> <标签名>"},
	}
	cmd := &cobra.Command{
		Use:   "capabilities",
		Short: "列出所有原子命令与主要参数",
		Long:  "列出所有原子命令与主要参数。加 --json 输出机器可读列表。",
		RunE: func(cmd *cobra.Command, args []string) error {
			if config.GetJSON() {
				return output.WriteJSON(os.Stdout, caps)
			}
			for _, c := range caps {
				fmt.Printf("%s %s\t%s\t%s\n", c.Domain, c.Cmd, c.Use, c.Args)
			}
			return nil
		},
	}
	return cmd
}
