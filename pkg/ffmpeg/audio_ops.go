package ffmpeg

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// ExtractAudio 从视频中提取音频
// inputPath: 输入视频文件路径
// outputPath: 输出音频文件路径
// options: 提取选项，可以为nil使用默认设置
func (f *FFmpeg) ExtractAudio(inputPath, outputPath string, options *ExtractAudioOptions) error {
	return f.ExtractAudioWithContext(context.Background(), inputPath, outputPath, options)
}

// ExtractAudioWithContext 带上下文的音频提取
func (f *FFmpeg) ExtractAudioWithContext(ctx context.Context, inputPath, outputPath string, options *ExtractAudioOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	// 设置默认选项
	if options == nil {
		options = &ExtractAudioOptions{
			Format:     getAudioFormatFromPath(outputPath),
			Bitrate:    "128k",
			SampleRate: 44100,
			Channels:   2,
		}
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath, // 输入文件
		"-y",            // 覆盖输出文件
		"-vn",           // 不包含视频流
	}

	// 添加开始时间
	if options.StartTime != "" {
		args = append(args, "-ss", options.StartTime)
	}

	// 添加持续时间
	if options.Duration != "" {
		args = append(args, "-t", options.Duration)
	}

	// 添加音频编码器
	if options.Format != "" {
		codec := getAudioCodecByFormat(options.Format)
		if codec != "" {
			args = append(args, "-c:a", codec)
		}
	}

	// 添加比特率
	if options.Bitrate != "" {
		args = append(args, "-b:a", options.Bitrate)
	}

	// 添加采样率
	if options.SampleRate > 0 {
		args = append(args, "-ar", fmt.Sprintf("%d", options.SampleRate))
	}

	// 添加声道数
	if options.Channels > 0 {
		args = append(args, "-ac", fmt.Sprintf("%d", options.Channels))
	}

	// 添加自定义参数
	args = append(args, options.CustomArgs...)

	// 添加输出文件
	args = append(args, outputPath)

	// 执行命令
	_, err := f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "音频提取失败", err)
	}

	f.logger.Info("音频提取完成: %s -> %s", inputPath, outputPath)
	return nil
}

// getAudioFormatFromPath 从文件路径获取音频格式
func getAudioFormatFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp3":
		return "mp3"
	case ".aac":
		return "aac"
	case ".wav":
		return "wav"
	case ".flac":
		return "flac"
	case ".ogg":
		return "ogg"
	case ".m4a":
		return "m4a"
	case ".wma":
		return "wma"
	default:
		return "mp3" // 默认格式
	}
}

// getAudioCodecByFormat 根据音频格式获取编码器
func getAudioCodecByFormat(format string) string {
	switch strings.ToLower(format) {
	case "mp3":
		return "libmp3lame"
	case "aac":
		return "aac"
	case "wav":
		return "pcm_s16le"
	case "flac":
		return "flac"
	case "ogg":
		return "libvorbis"
	case "m4a":
		return "aac"
	case "wma":
		return "wmav2"
	default:
		return "libmp3lame" // 默认编码器
	}
}

// ExtractAudioToMP3 提取音频为MP3格式的便捷方法
func (f *FFmpeg) ExtractAudioToMP3(inputPath, outputPath string) error {
	options := &ExtractAudioOptions{
		Format:     "mp3",
		Bitrate:    "192k",
		SampleRate: 44100,
		Channels:   2,
	}
	return f.ExtractAudio(inputPath, outputPath, options)
}

// ExtractAudioToAAC 提取音频为AAC格式的便捷方法
func (f *FFmpeg) ExtractAudioToAAC(inputPath, outputPath string) error {
	options := &ExtractAudioOptions{
		Format:     "aac",
		Bitrate:    "128k",
		SampleRate: 44100,
		Channels:   2,
	}
	return f.ExtractAudio(inputPath, outputPath, options)
}

// ExtractAudioToWAV 提取音频为WAV格式的便捷方法
func (f *FFmpeg) ExtractAudioToWAV(inputPath, outputPath string) error {
	options := &ExtractAudioOptions{
		Format:     "wav",
		SampleRate: 44100,
		Channels:   2,
	}
	return f.ExtractAudio(inputPath, outputPath, options)
}

// ExtractAudioToFLAC 提取音频为FLAC格式的便捷方法（无损）
func (f *FFmpeg) ExtractAudioToFLAC(inputPath, outputPath string) error {
	options := &ExtractAudioOptions{
		Format:     "flac",
		SampleRate: 44100,
		Channels:   2,
	}
	return f.ExtractAudio(inputPath, outputPath, options)
}

// ExtractAudioSegment 提取音频片段
func (f *FFmpeg) ExtractAudioSegment(inputPath, outputPath, startTime, duration string) error {
	options := &ExtractAudioOptions{
		Format:     getAudioFormatFromPath(outputPath),
		Bitrate:    "192k",
		SampleRate: 44100,
		Channels:   2,
		StartTime:  startTime,
		Duration:   duration,
	}
	return f.ExtractAudio(inputPath, outputPath, options)
}

// ExtractAudioHighQuality 提取高质量音频
func (f *FFmpeg) ExtractAudioHighQuality(inputPath, outputPath string) error {
	format := getAudioFormatFromPath(outputPath)
	options := &ExtractAudioOptions{
		Format:     format,
		SampleRate: 48000,
		Channels:   2,
	}

	// 根据格式设置高质量参数
	switch format {
	case "mp3":
		options.Bitrate = "320k"
	case "aac":
		options.Bitrate = "256k"
	case "flac":
		// FLAC是无损格式，不需要设置比特率
	default:
		options.Bitrate = "320k"
	}

	return f.ExtractAudio(inputPath, outputPath, options)
}

// ExtractAudioLowQuality 提取低质量音频（小文件）
func (f *FFmpeg) ExtractAudioLowQuality(inputPath, outputPath string) error {
	options := &ExtractAudioOptions{
		Format:     getAudioFormatFromPath(outputPath),
		Bitrate:    "64k",
		SampleRate: 22050,
		Channels:   1, // 单声道
	}
	return f.ExtractAudio(inputPath, outputPath, options)
}

// BatchExtractAudio 批量提取音频
func (f *FFmpeg) BatchExtractAudio(inputPaths []string, outputDir string, options *ExtractAudioOptions) error {
	if len(inputPaths) == 0 {
		return NewError(ErrInvalidInput, "输入文件列表不能为空", nil)
	}

	if outputDir == "" {
		return NewError(ErrInvalidInput, "输出目录不能为空", nil)
	}

	for i, inputPath := range inputPaths {
		// 生成输出文件名
		filename := filepath.Base(inputPath)
		ext := filepath.Ext(filename)
		nameWithoutExt := strings.TrimSuffix(filename, ext)
		
		var outputExt string
		if options != nil && options.Format != "" {
			outputExt = "." + options.Format
		} else {
			outputExt = ".mp3" // 默认提取为MP3
		}
		
		outputPath := filepath.Join(outputDir, nameWithoutExt+outputExt)

		// 提取单个文件的音频
		if err := f.ExtractAudio(inputPath, outputPath, options); err != nil {
			f.logger.Error("批量音频提取失败，文件: %s, 错误: %v", inputPath, err)
			return fmt.Errorf("批量音频提取失败，文件: %s, 错误: %w", inputPath, err)
		}

		f.logger.Info("批量音频提取进度: %d/%d - %s 完成", i+1, len(inputPaths), inputPath)
	}

	f.logger.Info("批量音频提取完成，共处理 %d 个文件", len(inputPaths))
	return nil
}
