package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var authURL, authToken string

// NewCommand 创建并返回 config 命令组
func NewCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "配置管理",
		Long:  "管理 GitLab URL 与访问令牌等配置。使用 -c/--config 可指定配置文件路径。",
	}

	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "设置 GitLab URL 与访问令牌",
		Long:  "将 GitLab 服务器 URL 与 access token 写入配置文件。可通过 --url/--token 传入，或未传时交互式输入。使用 -c/--config 可指定写入的配置文件路径。",
		Example: `  gitlab-tools config auth
  gitlab-tools config auth --url https://gitlab.example.com --token glpat-xxx
  gitlab-tools -c /path/to/config.yaml config auth`,
		RunE: runAuthCmd,
	}

	authCmd.Flags().StringVar(&authURL, "url", "", "GitLab 服务器 URL")
	authCmd.Flags().StringVar(&authToken, "token", "", "GitLab 访问令牌")

	configCmd.AddCommand(authCmd)
	return configCmd
}

func runAuthCmd(cmd *cobra.Command, args []string) error {
	urlVal := strings.TrimSpace(authURL)
	tokenVal := strings.TrimSpace(authToken)

	if urlVal == "" || tokenVal == "" {
		reader := bufio.NewReader(os.Stdin)
		if urlVal == "" {
			fmt.Fprint(os.Stderr, "GitLab URL (例如 https://gitlab.com): ")
			line, _ := reader.ReadString('\n')
			urlVal = strings.TrimSpace(line)
		}
		if tokenVal == "" {
			fmt.Fprint(os.Stderr, "GitLab Access Token: ")
			line, _ := reader.ReadString('\n')
			tokenVal = strings.TrimSpace(line)
		}
	}

	if urlVal == "" {
		return fmt.Errorf("未提供 GitLab URL")
	}
	if tokenVal == "" {
		return fmt.Errorf("未提供 Access Token")
	}

	targetPath := viper.GetString("config")
	if targetPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("无法获取用户主目录: %w", err)
		}
		targetPath = filepath.Join(home, ".config", "gitlab-tools", "config.yaml")
	}

	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	viper.Set("url", urlVal)
	viper.Set("token", tokenVal)
	if err := viper.WriteConfigAs(targetPath); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	if !GetJSON() {
		fmt.Fprintf(os.Stderr, "已写入配置: %s\n", targetPath)
	}
	return nil
}
