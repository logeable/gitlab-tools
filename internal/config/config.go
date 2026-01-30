package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// ErrUsage 表示用法错误（缺少必填参数、非法标志等），main 应退出码 2。
var ErrUsage = errors.New("usage error")

var initOnce sync.Once

// Init 初始化 Viper 配置，应在根命令 PersistentPreRunE 中调用并传入 -c/--config 的解析值。
// configFilePath 为空时使用默认搜索路径，非空时仅从该文件读取。
// 配置优先级：命令行参数 > 环境变量 > 配置文件 > 默认值
func Init(configFilePath string) {
	initOnce.Do(func() {
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		viper.SetEnvPrefix("GITLAB")

		if configFilePath != "" {
			viper.SetConfigFile(configFilePath)
		} else {
			viper.SetConfigName("config")
			viper.AddConfigPath(".")
			if home, err := os.UserHomeDir(); err == nil {
				viper.AddConfigPath(filepath.Join(home, ".config", "gitlab-tools"))
				viper.AddConfigPath(home)
				viper.AddConfigPath(filepath.Join(home, ".gitlab-tools"))
			}
		}

		err := viper.ReadInConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			} else {
				fmt.Fprintf(os.Stderr, "读取配置文件失败: %v\n", err)
				os.Exit(1)
			}
		}

		viper.SetDefault("url", "https://gitlab.com")
		viper.SetDefault("project.limit", 20)
		viper.SetDefault("pipeline.list.limit", 5)
	})
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

// GetJSON 是否以 JSON 格式输出（全局 --json 标志）
func GetJSON() bool {
	return viper.GetBool("json")
}
