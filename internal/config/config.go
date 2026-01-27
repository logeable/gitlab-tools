package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Init 初始化 Viper 配置
// 配置优先级：命令行参数 > 环境变量 > 配置文件 > 默认值
func Init() {
	// 设置环境变量前缀（可选，如果设置则环境变量需要以 GITLAB_TOOLS_ 开头）
	// viper.SetEnvPrefix("GITLAB_TOOLS")

	// 自动读取环境变量
	viper.AutomaticEnv()

	// 将环境变量名中的点号和横线替换为下划线（Viper 的默认行为）
	// 例如：GITLAB_URL, GITLAB_TOKEN
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetEnvPrefix("GITLAB")

	// 配置文件设置
	viper.SetConfigName("config") // 配置文件名称（不含扩展名）

	// 添加配置文件搜索路径（按优先级顺序）
	// 1. 当前工作目录
	viper.AddConfigPath(".")

	// 2. 用户主目录下的 .config 目录
	if home, err := os.UserHomeDir(); err == nil {
		configDir := filepath.Join(home, ".config", "gitlab-tools")
		viper.AddConfigPath(configDir)
		// 也支持直接在用户主目录
		viper.AddConfigPath(home)
	}

	// 3. 用户主目录下的 .gitlab-tools 目录
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(filepath.Join(home, ".gitlab-tools"))
	}

	// 读取配置文件（如果存在，忽略文件不存在的错误）
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			fmt.Fprintf(os.Stderr, "读取配置文件失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 设置默认值
	viper.SetDefault("url", "https://gitlab.com")
	viper.SetDefault("project.limit", 20)
	viper.SetDefault("pipeline.list.limit", 5)

}

// GetGitLabURL 获取 GitLab 服务器 URL
// 配置优先级：命令行参数 > 环境变量 > 配置文件 > 默认值
func GetGitLabURL() string {
	if url := viper.GetString("url"); url != "" {
		return url
	}
	return "https://gitlab.com"
}

// GetGitLabToken 获取 GitLab 访问令牌
// 配置优先级：命令行参数 > 环境变量 > 配置文件 > 默认值
func GetGitLabToken() string {
	return viper.GetString("token")
}
