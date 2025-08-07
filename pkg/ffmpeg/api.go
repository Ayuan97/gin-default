// Package ffmpeg 提供统一的API入口
package ffmpeg

import (
	"fmt"
	"time"
)

// 重新导出核心类型（这些类型在子目录中定义，但属于同一个包）
// 由于Go的包系统，同一包中的类型在任何目录中都是可见的

// API 统一的API入口结构
type API struct {
	ffmpeg *FFmpeg
}

// NewAPI 创建新的API实例
func NewAPI(config *Config) (*API, error) {
	if config == nil {
		config = &Config{
			Timeout: 10 * time.Minute,
		}
	}

	ffmpeg, err := New(config)
	if err != nil {
		return nil, err
	}

	return &API{
		ffmpeg: ffmpeg,
	}, nil
}

// NewEditor 创建新的视频编辑器
func (api *API) NewEditor(inputPath string) *VideoEditor {
	return NewVideoEditor(api.ffmpeg, inputPath)
}

// GetFFmpeg 获取底层FFmpeg实例（用于直接调用传统方法）
func (api *API) GetFFmpeg() *FFmpeg {
	return api.ffmpeg
}

// === 便捷方法：直接调用常用功能 ===

// Convert 视频格式转换
func (api *API) Convert(inputPath, outputPath string, options *ConvertOptions) error {
	return api.ffmpeg.Convert(inputPath, outputPath, options)
}

// GetVideoInfo 获取视频信息
func (api *API) GetVideoInfo(inputPath string) (*VideoInfo, error) {
	return api.ffmpeg.GetVideoInfo(inputPath)
}

// === 快速链式编辑方法 ===

// QuickEdit 快速编辑（常用操作的组合）
func (api *API) QuickEdit(inputPath, outputPath string) *VideoEditor {
	return api.NewEditor(inputPath).SetOutput(outputPath)
}

// QuickCompose 快速合成（多媒体合成的起点）
func (api *API) QuickCompose(mainVideoPath, outputPath string) *VideoEditor {
	return api.NewEditor(mainVideoPath).
		SetOutput(outputPath).
		SetProgressCallback(func(progress float64, currentTime, totalTime time.Duration) {
			// 默认的简单进度显示
			// 用户可以通过SetProgressCallback覆盖
		})
}

// === 批量处理方法 ===

// BatchProcess 批量处理多个文件
func (api *API) BatchProcess(inputPaths []string, outputDir string, processor func(*VideoEditor) *VideoEditor) error {
	for i, inputPath := range inputPaths {
		editor := api.NewEditor(inputPath)

		// 生成输出路径
		outputPath := generateBatchOutputPath(outputDir, inputPath, i)
		editor.SetOutput(outputPath)

		// 应用用户定义的处理逻辑
		editor = processor(editor)

		// 执行处理
		if err := editor.Execute(); err != nil {
			return err
		}
	}
	return nil
}

// generateBatchOutputPath 生成批量处理的输出路径
func generateBatchOutputPath(outputDir, inputPath string, index int) string {
	// 简化实现，实际应该更智能地处理文件名和扩展名
	return fmt.Sprintf("%s/processed_%d.mp4", outputDir, index+1)
}
