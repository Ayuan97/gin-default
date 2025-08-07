package ffmpeg

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// Compress 压缩视频文件
// inputPath: 输入文件路径
// outputPath: 输出文件路径
// options: 压缩选项，可以为nil使用默认设置
func (f *FFmpeg) Compress(inputPath, outputPath string, options *CompressOptions) error {
	return f.CompressWithContext(context.Background(), inputPath, outputPath, options)
}

// CompressWithContext 带上下文的视频压缩
func (f *FFmpeg) CompressWithContext(ctx context.Context, inputPath, outputPath string, options *CompressOptions) error {
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
		options = &CompressOptions{
			CRF:        23,
			Preset:     "medium",
			VideoCodec: "libx264",
			AudioCodec: "aac",
		}
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath, // 输入文件
		"-y", // 覆盖输出文件
	}

	// 添加视频编码器
	if options.VideoCodec != "" {
		args = append(args, "-c:v", options.VideoCodec)
	} else {
		args = append(args, "-c:v", "libx264")
	}

	// 添加音频编码器
	if options.AudioCodec != "" {
		args = append(args, "-c:a", options.AudioCodec)
	} else {
		args = append(args, "-c:a", "aac")
	}

	// 添加分辨率设置
	if options.Width > 0 || options.Height > 0 {
		var scaleFilter string
		if options.Width > 0 && options.Height > 0 {
			scaleFilter = fmt.Sprintf("scale=%d:%d", options.Width, options.Height)
		} else if options.Width > 0 {
			scaleFilter = fmt.Sprintf("scale=%d:-2", options.Width)
		} else {
			scaleFilter = fmt.Sprintf("scale=-2:%d", options.Height)
		}
		args = append(args, "-vf", scaleFilter)
	}

	// 添加比特率设置
	if options.Bitrate != "" {
		args = append(args, "-b:v", options.Bitrate)
	}

	// 添加帧率设置
	if options.FrameRate > 0 {
		args = append(args, "-r", fmt.Sprintf("%.2f", options.FrameRate))
	}

	// 添加CRF设置（恒定质量因子）
	if options.CRF > 0 {
		args = append(args, "-crf", strconv.Itoa(options.CRF))
	}

	// 添加编码预设
	if options.Preset != "" {
		args = append(args, "-preset", options.Preset)
	}

	// 添加自定义参数
	args = append(args, options.CustomArgs...)

	// 添加输出文件
	args = append(args, outputPath)

	// 执行命令
	_, err := f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "视频压缩失败", err)
	}

	f.logger.Info("视频压缩完成: %s -> %s", inputPath, outputPath)
	return nil
}

// CompressToSize 压缩视频到指定文件大小（近似）
// targetSizeMB: 目标文件大小（MB）
func (f *FFmpeg) CompressToSize(inputPath, outputPath string, targetSizeMB float64) error {
	// 首先获取视频时长
	duration, err := f.GetVideoDuration(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频时长失败: %w", err)
	}

	// 计算目标比特率
	durationSeconds := duration.Seconds()
	targetBitrate := int((targetSizeMB * 8 * 1024 * 1024) / durationSeconds * 0.9) // 0.9是安全系数

	options := &CompressOptions{
		Bitrate:    fmt.Sprintf("%dk", targetBitrate),
		VideoCodec: "libx264",
		AudioCodec: "aac",
		Preset:     "medium",
	}

	return f.Compress(inputPath, outputPath, options)
}

// CompressLowQuality 低质量压缩（高压缩比）
func (f *FFmpeg) CompressLowQuality(inputPath, outputPath string) error {
	options := &CompressOptions{
		CRF:        28,
		Preset:     "fast",
		VideoCodec: "libx264",
		AudioCodec: "aac",
		Width:      854, // 480p宽度
		Height:     480, // 480p高度
		FrameRate:  24,  // 降低帧率
	}
	return f.Compress(inputPath, outputPath, options)
}

// CompressMediumQuality 中等质量压缩
func (f *FFmpeg) CompressMediumQuality(inputPath, outputPath string) error {
	options := &CompressOptions{
		CRF:        23,
		Preset:     "medium",
		VideoCodec: "libx264",
		AudioCodec: "aac",
		Width:      1280, // 720p宽度
		Height:     720,  // 720p高度
	}
	return f.Compress(inputPath, outputPath, options)
}

// CompressHighQuality 高质量压缩（低压缩比）
func (f *FFmpeg) CompressHighQuality(inputPath, outputPath string) error {
	options := &CompressOptions{
		CRF:        18,
		Preset:     "slow",
		VideoCodec: "libx264",
		AudioCodec: "aac",
		Width:      1920, // 1080p宽度
		Height:     1080, // 1080p高度
	}
	return f.Compress(inputPath, outputPath, options)
}

// CompressFor4K 4K视频压缩
func (f *FFmpeg) CompressFor4K(inputPath, outputPath string) error {
	options := &CompressOptions{
		CRF:        20,
		Preset:     "slow",
		VideoCodec: "libx265", // 使用H.265编码器获得更好的压缩效果
		AudioCodec: "aac",
		Width:      3840, // 4K宽度
		Height:     2160, // 4K高度
	}
	return f.Compress(inputPath, outputPath, options)
}

// CompressForWeb 为Web优化的压缩
func (f *FFmpeg) CompressForWeb(inputPath, outputPath string) error {
	options := &CompressOptions{
		CRF:        25,
		Preset:     "fast",
		VideoCodec: "libx264",
		AudioCodec: "aac",
		Width:      1280, // 720p
		Height:     720,
		FrameRate:  30,
		CustomArgs: []string{
			"-movflags", "+faststart", // 优化Web播放
			"-profile:v", "baseline", // 兼容性更好的配置
			"-level", "3.0",
		},
	}
	return f.Compress(inputPath, outputPath, options)
}

// CompressForMobile 为移动设备优化的压缩
func (f *FFmpeg) CompressForMobile(inputPath, outputPath string) error {
	options := &CompressOptions{
		CRF:        26,
		Preset:     "fast",
		VideoCodec: "libx264",
		AudioCodec: "aac",
		Width:      854, // 480p
		Height:     480,
		FrameRate:  24,
		CustomArgs: []string{
			"-movflags", "+faststart",
			"-profile:v", "baseline",
			"-level", "3.0",
			"-b:a", "96k", // 降低音频比特率
		},
	}
	return f.Compress(inputPath, outputPath, options)
}

// BatchCompress 批量压缩视频
func (f *FFmpeg) BatchCompress(inputPaths []string, outputDir string, options *CompressOptions) error {
	if len(inputPaths) == 0 {
		return NewError(ErrInvalidInput, "输入文件列表不能为空", nil)
	}

	if outputDir == "" {
		return NewError(ErrInvalidInput, "输出目录不能为空", nil)
	}

	for i, inputPath := range inputPaths {
		// 生成输出文件名
		outputPath := generateOutputPath(inputPath, outputDir, "_compressed")

		// 压缩单个文件
		if err := f.Compress(inputPath, outputPath, options); err != nil {
			f.logger.Error("批量压缩失败，文件: %s, 错误: %v", inputPath, err)
			return fmt.Errorf("批量压缩失败，文件: %s, 错误: %w", inputPath, err)
		}

		f.logger.Info("批量压缩进度: %d/%d - %s 完成", i+1, len(inputPaths), inputPath)
	}

	f.logger.Info("批量压缩完成，共处理 %d 个文件", len(inputPaths))
	return nil
}

// generateOutputPath 生成输出文件路径
func generateOutputPath(inputPath, outputDir, suffix string) string {
	filename := filepath.Base(inputPath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return filepath.Join(outputDir, nameWithoutExt+suffix+ext)
}
