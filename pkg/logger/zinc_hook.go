package logger

import (
	"justus/pkg/setting"
	"justus/pkg/zincsearch"
)

// NewZincHook 创建 ZincSearch 日志钩子
func NewZincHook(s *setting.LoggerSettingS) (*zincsearch.ZincHook, error) {
	// 创建 ZincSearch 客户端
	client := zincsearch.NewCustomClient(
		s.LogZincHost,
		s.LogZincUser,
		s.LogZincPassword,
		30, // 默认30秒超时
	)

	// 检查连接
	if err := client.Ping(); err != nil {
		return nil, err
	}

	// 创建钩子
	hook := zincsearch.NewZincHook(client, s.LogZincIndex)

	return hook, nil
}
