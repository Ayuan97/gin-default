package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Crop 裁剪视频
// inputPath: 输入视频文件路径
// outputPath: 输出视频文件路径
// options: 裁剪选项
func (f *FFmpeg) Crop(inputPath, outputPath string, options *CropOptions) error {
	return f.CropWithContext(context.Background(), inputPath, outputPath, options)
}

// CropWithContext 带上下文的视频裁剪
func (f *FFmpeg) CropWithContext(ctx context.Context, inputPath, outputPath string, options *CropOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	if options == nil {
		return NewError(ErrInvalidOptions, "裁剪选项不能为空", nil)
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath, // 输入文件
		"-y", // 覆盖输出文件
	}

	// 添加开始时间
	if options.StartTime != "" {
		args = append(args, "-ss", options.StartTime)
	}

	// 添加持续时间或结束时间
	if options.Duration != "" {
		args = append(args, "-t", options.Duration)
	} else if options.EndTime != "" {
		args = append(args, "-to", options.EndTime)
	}

	// 添加编码参数以保持质量
	args = append(args, "-c", "copy") // 尽可能使用流复制以提高速度

	// 添加自定义参数
	args = append(args, options.CustomArgs...)

	// 添加输出文件
	args = append(args, outputPath)

	// 执行命令
	_, err := f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "视频裁剪失败", err)
	}

	f.logger.Info("视频裁剪完成: %s -> %s", inputPath, outputPath)
	return nil
}

// Merge 合并多个视频文件
// inputPaths: 输入视频文件路径列表
// outputPath: 输出视频文件路径
// options: 合并选项
func (f *FFmpeg) Merge(inputPaths []string, outputPath string, options *MergeOptions) error {
	return f.MergeWithContext(context.Background(), inputPaths, outputPath, options)
}

// MergeWithContext 带上下文的视频合并
func (f *FFmpeg) MergeWithContext(ctx context.Context, inputPaths []string, outputPath string, options *MergeOptions) error {
	if len(inputPaths) < 2 {
		return NewError(ErrInvalidInput, "至少需要两个输入文件进行合并", nil)
	}

	// 验证所有输入文件
	for _, inputPath := range inputPaths {
		if err := validateInputFile(inputPath); err != nil {
			return err
		}
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	// 设置默认选项
	if options == nil {
		options = &MergeOptions{
			Method: "concat",
		}
	}

	var args []string

	if options.Method == "filter" {
		// 使用filter方法合并（适用于不同格式的视频）
		args = f.buildFilterMergeArgs(inputPaths, outputPath, options)
	} else {
		// 使用concat方法合并（适用于相同格式的视频）
		args = f.buildConcatMergeArgs(inputPaths, outputPath, options)
	}

	// 执行命令
	_, err := f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "视频合并失败", err)
	}

	f.logger.Info("视频合并完成: %d个文件 -> %s", len(inputPaths), outputPath)
	return nil
}

// buildConcatMergeArgs 构建concat方法的合并参数
func (f *FFmpeg) buildConcatMergeArgs(inputPaths []string, outputPath string, options *MergeOptions) []string {
	args := []string{"-y"} // 覆盖输出文件

	// 添加所有输入文件
	for _, inputPath := range inputPaths {
		args = append(args, "-i", inputPath)
	}

	// 构建filter_complex参数
	var filterParts []string
	for i := range inputPaths {
		filterParts = append(filterParts, fmt.Sprintf("[%d:v][%d:a]", i, i))
	}

	filterComplex := strings.Join(filterParts, "") + fmt.Sprintf("concat=n=%d:v=1:a=1[outv][outa]", len(inputPaths))
	args = append(args, "-filter_complex", filterComplex)
	args = append(args, "-map", "[outv]", "-map", "[outa]")

	// 添加自定义参数
	args = append(args, options.CustomArgs...)

	// 添加输出文件
	args = append(args, outputPath)

	return args
}

// buildFilterMergeArgs 构建filter方法的合并参数
func (f *FFmpeg) buildFilterMergeArgs(inputPaths []string, outputPath string, options *MergeOptions) []string {
	// 创建临时文件列表
	tempFile, err := f.createConcatFile(inputPaths)
	if err != nil {
		f.logger.Error("创建临时文件列表失败: %v", err)
		return nil
	}
	defer os.Remove(tempFile) // 清理临时文件

	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", tempFile,
		"-c", "copy",
		"-y", // 覆盖输出文件
	}

	// 添加自定义参数
	args = append(args, options.CustomArgs...)

	// 添加输出文件
	args = append(args, outputPath)

	return args
}

// createConcatFile 创建用于concat的临时文件列表
func (f *FFmpeg) createConcatFile(inputPaths []string) (string, error) {
	tempFile, err := os.CreateTemp("", "ffmpeg_concat_*.txt")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	for _, inputPath := range inputPaths {
		// 转换为绝对路径
		absPath, err := filepath.Abs(inputPath)
		if err != nil {
			return "", err
		}

		// 写入文件列表
		_, err = fmt.Fprintf(tempFile, "file '%s'\n", absPath)
		if err != nil {
			return "", err
		}
	}

	return tempFile.Name(), nil
}

// Screenshot 截取视频帧作为图片
// inputPath: 输入视频文件路径
// outputPath: 输出图片文件路径
// options: 截图选项
func (f *FFmpeg) Screenshot(inputPath, outputPath string, options *ScreenshotOptions) error {
	return f.ScreenshotWithContext(context.Background(), inputPath, outputPath, options)
}

// ScreenshotWithContext 带上下文的视频截图
func (f *FFmpeg) ScreenshotWithContext(ctx context.Context, inputPath, outputPath string, options *ScreenshotOptions) error {
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
		options = &ScreenshotOptions{
			Time:   "00:00:01",
			Format: "jpg",
		}
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath, // 输入文件
		"-y", // 覆盖输出文件
	}

	// 添加时间点
	if options.Time != "" {
		args = append(args, "-ss", options.Time)
	}

	// 添加帧数限制
	args = append(args, "-frames:v", "1")

	// 添加分辨率设置
	if options.Width > 0 || options.Height > 0 {
		var scaleFilter string
		if options.Width > 0 && options.Height > 0 {
			scaleFilter = fmt.Sprintf("scale=%d:%d", options.Width, options.Height)
		} else if options.Width > 0 {
			scaleFilter = fmt.Sprintf("scale=%d:-1", options.Width)
		} else {
			scaleFilter = fmt.Sprintf("scale=-1:%d", options.Height)
		}
		args = append(args, "-vf", scaleFilter)
	}

	// 添加质量设置
	if options.Quality > 0 {
		args = append(args, "-q:v", fmt.Sprintf("%d", options.Quality))
	}

	// 添加自定义参数
	args = append(args, options.CustomArgs...)

	// 添加输出文件
	args = append(args, outputPath)

	// 执行命令
	_, err := f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "视频截图失败", err)
	}

	f.logger.Info("视频截图完成: %s -> %s", inputPath, outputPath)
	return nil
}

// CropTimeRange 裁剪指定时间范围的视频片段
func (f *FFmpeg) CropTimeRange(inputPath, outputPath, startTime, endTime string) error {
	options := &CropOptions{
		StartTime: startTime,
		EndTime:   endTime,
	}
	return f.Crop(inputPath, outputPath, options)
}

// CropDuration 裁剪指定时长的视频片段
func (f *FFmpeg) CropDuration(inputPath, outputPath, startTime, duration string) error {
	options := &CropOptions{
		StartTime: startTime,
		Duration:  duration,
	}
	return f.Crop(inputPath, outputPath, options)
}

// ScreenshotAtTime 在指定时间点截图
func (f *FFmpeg) ScreenshotAtTime(inputPath, outputPath, timePoint string) error {
	options := &ScreenshotOptions{
		Time:   timePoint,
		Format: "jpg",
	}
	return f.Screenshot(inputPath, outputPath, options)
}

// ScreenshotMultiple 在多个时间点截图
func (f *FFmpeg) ScreenshotMultiple(inputPath string, timePoints []string, outputDir string) error {
	if len(timePoints) == 0 {
		return NewError(ErrInvalidInput, "时间点列表不能为空", nil)
	}

	for i, timePoint := range timePoints {
		filename := fmt.Sprintf("screenshot_%03d.jpg", i+1)
		outputPath := filepath.Join(outputDir, filename)

		if err := f.ScreenshotAtTime(inputPath, outputPath, timePoint); err != nil {
			return fmt.Errorf("截图失败，时间点: %s, 错误: %w", timePoint, err)
		}

		f.logger.Info("多点截图进度: %d/%d - %s 完成", i+1, len(timePoints), timePoint)
	}

	f.logger.Info("多点截图完成，共生成 %d 张图片", len(timePoints))
	return nil
}
