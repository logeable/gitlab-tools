package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Init 初始化 Viper 配置
func Init() {
	// 设置环境变量前缀（可选，如果设置则环境变量需要以 GITLAB_TOOLS_ 开头）
	// viper.SetEnvPrefix("GITLAB_TOOLS")

	// 自动读取环境变量
	viper.AutomaticEnv()

	// 将环境变量名中的点号和横线替换为下划线（Viper 的默认行为）
	// 例如：GITLAB_URL, GITLAB_TOKEN
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// 设置默认值
	viper.SetDefault("url", "https://gitlab.com")
	viper.SetDefault("project.limit", 20)
	viper.SetDefault("pipeline.list.limit", 5)
}

// GetGitLabURL 获取 GitLab 服务器 URL
func GetGitLabURL() string {
	// Viper 的优先级：命令行参数 > 环境变量 > 默认值
	// 如果命令行参数设置了，viper.GetString 会返回命令行参数的值
	if url := viper.GetString("url"); url != "" {
		return url
	}
	return "https://gitlab.com"
}

// GetGitLabToken 获取 GitLab 访问令牌
func GetGitLabToken() string {
	// Viper 自动从环境变量 GITLAB_TOKEN 或命令行参数 --token 读取
	return viper.GetString("token")
}
