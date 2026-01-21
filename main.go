package main

import (
	"fmt"
	"os"

	"gitlab-tools/internal/branch"
	"gitlab-tools/internal/config"
	"gitlab-tools/internal/mr"
	"gitlab-tools/internal/pipeline"
	"gitlab-tools/internal/project"
	"gitlab-tools/internal/tag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	gitlabURL   string
	gitlabToken string
)

var rootCmd = &cobra.Command{
	Use:   "gitlab-tools",
	Short: "GitLab 工具集",
	Long:  "一个用于与 GitLab API 交互的命令行工具集",
}

func init() {
	// 初始化 Viper
	config.Init()

	// 全局标志
	rootCmd.PersistentFlags().StringVar(&gitlabURL, "url", "", "GitLab 服务器 URL (默认: https://gitlab.com)")
	rootCmd.PersistentFlags().StringVar(&gitlabToken, "token", "", "GitLab 访问令牌 (也可以通过 GITLAB_TOKEN 环境变量设置)")

	// 将 Cobra flags 绑定到 Viper（自动从环境变量读取）
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))

	// 添加子命令
	rootCmd.AddCommand(pipeline.NewCommand())
	rootCmd.AddCommand(project.NewCommand())
	rootCmd.AddCommand(branch.NewCommand())
	rootCmd.AddCommand(mr.NewCommand())
	rootCmd.AddCommand(tag.NewCommand())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
