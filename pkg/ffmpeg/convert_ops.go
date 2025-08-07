package ffmpeg

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// Convert 转换视频格式
// inputPath: 输入文件路径
// outputPath: 输出文件路径
// options: 转换选项，可以为nil使用默认设置
func (f *FFmpeg) Convert(inputPath, outputPath string, options *ConvertOptions) error {
	return f.ConvertWithContext(context.Background(), inputPath, outputPath, options)
}

// ConvertWithContext 带上下文的视频格式转换
func (f *FFmpeg) ConvertWithContext(ctx context.Context, inputPath, outputPath string, options *ConvertOptions) error {
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
		options = &ConvertOptions{}
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath, // 输入文件
		"-y",            // 覆盖输出文件
	}

	// 添加视频编码器
	if options.VideoCodec != "" {
		args = append(args, "-c:v", options.VideoCodec)
	} else {
		// 根据输出格式自动选择编码器
		if codec := getDefaultVideoCodec(outputPath); codec != "" {
			args = append(args, "-c:v", codec)
		}
	}

	// 添加音频编码器
	if options.AudioCodec != "" {
		args = append(args, "-c:a", options.AudioCodec)
	} else {
		// 根据输出格式自动选择编码器
		if codec := getDefaultAudioCodec(outputPath); codec != "" {
			args = append(args, "-c:a", codec)
		}
	}

	// 添加质量设置
	if options.Quality != "" {
		args = append(args, getQualityArgs(options.Quality)...)
	}

	// 添加元数据
	for key, value := range options.Metadata {
		args = append(args, "-metadata", fmt.Sprintf("%s=%s", key, value))
	}

	// 添加自定义参数
	args = append(args, options.CustomArgs...)

	// 添加输出文件
	args = append(args, outputPath)

	// 执行命令
	_, err := f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "视频格式转换失败", err)
	}

	f.logger.Info("视频格式转换完成: %s -> %s", inputPath, outputPath)
	return nil
}

// getDefaultVideoCodec 根据输出文件扩展名获取默认视频编码器
func getDefaultVideoCodec(outputPath string) string {
	ext := strings.ToLower(filepath.Ext(outputPath))
	switch ext {
	case ".mp4":
		return "libx264"
	case ".avi":
		return "libx264"
	case ".mov":
		return "libx264"
	case ".mkv":
		return "libx264"
	case ".webm":
		return "libvpx-vp9"
	case ".flv":
		return "libx264"
	case ".wmv":
		return "wmv2"
	case ".3gp":
		return "libx264"
	default:
		return "libx264" // 默认使用H.264编码器
	}
}

// getDefaultAudioCodec 根据输出文件扩展名获取默认音频编码器
func getDefaultAudioCodec(outputPath string) string {
	ext := strings.ToLower(filepath.Ext(outputPath))
	switch ext {
	case ".mp4":
		return "aac"
	case ".avi":
		return "mp3"
	case ".mov":
		return "aac"
	case ".mkv":
		return "aac"
	case ".webm":
		return "libvorbis"
	case ".flv":
		return "aac"
	case ".wmv":
		return "wmav2"
	case ".3gp":
		return "aac"
	default:
		return "aac" // 默认使用AAC编码器
	}
}

// getQualityArgs 根据质量设置获取FFmpeg参数
func getQualityArgs(quality string) []string {
	switch strings.ToLower(quality) {
	case "high":
		return []string{"-crf", "18", "-preset", "slow"}
	case "medium":
		return []string{"-crf", "23", "-preset", "medium"}
	case "low":
		return []string{"-crf", "28", "-preset", "fast"}
	default:
		// 如果是数字，假设是CRF值
		if strings.Contains(quality, "crf") {
			parts := strings.Split(quality, ":")
			if len(parts) == 2 {
				return []string{"-crf", parts[1]}
			}
		}
		// 如果是纯数字，直接作为CRF值
		return []string{"-crf", quality}
	}
}

// ConvertToMP4 转换为MP4格式的便捷方法
func (f *FFmpeg) ConvertToMP4(inputPath, outputPath string) error {
	options := &ConvertOptions{
		OutputFormat: "mp4",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		Quality:      "medium",
	}
	return f.Convert(inputPath, outputPath, options)
}

// ConvertToAVI 转换为AVI格式的便捷方法
func (f *FFmpeg) ConvertToAVI(inputPath, outputPath string) error {
	options := &ConvertOptions{
		OutputFormat: "avi",
		VideoCodec:   "libx264",
		AudioCodec:   "mp3",
		Quality:      "medium",
	}
	return f.Convert(inputPath, outputPath, options)
}

// ConvertToMOV 转换为MOV格式的便捷方法
func (f *FFmpeg) ConvertToMOV(inputPath, outputPath string) error {
	options := &ConvertOptions{
		OutputFormat: "mov",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		Quality:      "medium",
	}
	return f.Convert(inputPath, outputPath, options)
}

// ConvertToMKV 转换为MKV格式的便捷方法
func (f *FFmpeg) ConvertToMKV(inputPath, outputPath string) error {
	options := &ConvertOptions{
		OutputFormat: "mkv",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		Quality:      "medium",
	}
	return f.Convert(inputPath, outputPath, options)
}

// ConvertToWebM 转换为WebM格式的便捷方法
func (f *FFmpeg) ConvertToWebM(inputPath, outputPath string) error {
	options := &ConvertOptions{
		OutputFormat: "webm",
		VideoCodec:   "libvpx-vp9",
		AudioCodec:   "libvorbis",
		Quality:      "medium",
	}
	return f.Convert(inputPath, outputPath, options)
}

// BatchConvert 批量转换视频格式
func (f *FFmpeg) BatchConvert(inputPaths []string, outputDir string, options *ConvertOptions) error {
	if len(inputPaths) == 0 {
		return NewError(ErrInvalidInput, "输入文件列表不能为空", nil)
	}

	if outputDir == "" {
		return NewError(ErrInvalidInput, "输出目录不能为空", nil)
	}

	// 验证输出目录
	if err := validateOutputFile(filepath.Join(outputDir, "test")); err != nil {
		return err
	}

	for _, inputPath := range inputPaths {
		// 生成输出文件名
		filename := filepath.Base(inputPath)
		ext := filepath.Ext(filename)
		nameWithoutExt := strings.TrimSuffix(filename, ext)
		
		var outputExt string
		if options != nil && options.OutputFormat != "" {
			outputExt = "." + options.OutputFormat
		} else {
			outputExt = ".mp4" // 默认转换为MP4
		}
		
		outputPath := filepath.Join(outputDir, nameWithoutExt+outputExt)

		// 转换单个文件
		if err := f.Convert(inputPath, outputPath, options); err != nil {
			f.logger.Error("批量转换失败，文件: %s, 错误: %v", inputPath, err)
			return fmt.Errorf("批量转换失败，文件: %s, 错误: %w", inputPath, err)
		}

		f.logger.Info("批量转换进度: %s 完成", inputPath)
	}

	f.logger.Info("批量转换完成，共处理 %d 个文件", len(inputPaths))
	return nil
}
