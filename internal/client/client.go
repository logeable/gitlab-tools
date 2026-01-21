package client

import (
	"fmt"

	"gitlab-tools/internal/config"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// NewClient 创建并返回一个新的 GitLab 客户端
func NewClient() (*gitlab.Client, error) {
	url := config.GetGitLabURL()
	token := config.GetGitLabToken()
	if token == "" {
		return nil, fmt.Errorf("错误: 请设置 GITLAB_TOKEN 环境变量或使用 --token 标志")
	}

	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return nil, fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	return client, nil
}
