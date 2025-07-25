package logger

import (
	"justus/pkg/setting"
	"justus/pkg/zincsearch"
)

// NewZincHook 创建 ZincSearch 日志钩子
func NewZincHook() (*zincsearch.ZincHook, error) {
	// 使用统一的 ZincSearch 配置
	config := setting.ZincSearchSetting

	// 创建 ZincSearch 客户端
	client := zincsearch.NewCustomClient(
		config.Host,
		config.Username,
		config.Password,
		config.Timeout,
	)

	// 检查连接
	if err := client.Ping(); err != nil {
		return nil, err
	}

	// 创建钩子，使用默认索引
	hook := zincsearch.NewZincHook(client, config.DefaultIndex)

	return hook, nil
}
