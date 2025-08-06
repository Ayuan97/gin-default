package ffmpeg_test

import (
	"fmt"
	"log"
	"time"

	"justus/pkg/ffmpeg"
)

// Example_basicUsage 基本使用示例
func Example_basicUsage() {
	// 创建FFmpeg实例
	ff, err := ffmpeg.New(nil)
	if err != nil {
		log.Fatal("初始化FFmpeg失败:", err)
	}

	// 获取FFmpeg版本
	version, err := ff.GetFFmpegVersion()
	if err == nil {
		fmt.Printf("FFmpeg版本: %s\n", version)
	}

	// 检查支持的功能
	features, err := ff.CheckFFmpegFeatures()
	if err == nil {
		fmt.Printf("支持H.264编码: %v\n", features["libx264"])
		fmt.Printf("支持H.265编码: %v\n", features["libx265"])
	}
}

// Example_videoConversion 视频格式转换示例
func Example_videoConversion() {
	ff, err := ffmpeg.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// 基本转换
	err = ff.Convert("input.avi", "output.mp4", nil)
	if err != nil {
		log.Printf("转换失败: %v", err)
		return
	}

	// 高级转换选项
	options := &ffmpeg.ConvertOptions{
		VideoCodec: "libx264",
		AudioCodec: "aac",
		Quality:    "high",
		Metadata: map[string]string{
			"title":   "示例视频",
			"comment": "使用FFmpeg Go包转换",
		},
	}
	err = ff.Convert("input.mov", "output_hq.mp4", options)
	if err != nil {
		log.Printf("高质量转换失败: %v", err)
		return
	}

	// 便捷转换方法
	err = ff.ConvertToWebM("input.mp4", "output.webm")
	if err != nil {
		log.Printf("WebM转换失败: %v", err)
		return
	}

	fmt.Println("视频转换完成")
}

// Example_videoCompression 视频压缩示例
func Example_videoCompression() {
	ff, err := ffmpeg.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// 自定义压缩选项
	options := &ffmpeg.CompressOptions{
		Width:      1280,
		Height:     720,
		CRF:        23,
		Preset:     "medium",
		VideoCodec: "libx264",
		AudioCodec: "aac",
		FrameRate:  30,
	}
	err = ff.Compress("large_video.mp4", "compressed.mp4", options)
	if err != nil {
		log.Printf("压缩失败: %v", err)
		return
	}

	// 预设压缩方法
	err = ff.CompressForWeb("input.mp4", "web_optimized.mp4")
	if err != nil {
		log.Printf("Web优化失败: %v", err)
		return
	}

	// 压缩到指定大小
	err = ff.CompressToSize("input.mp4", "50mb_output.mp4", 50.0)
	if err != nil {
		log.Printf("大小压缩失败: %v", err)
		return
	}

	fmt.Println("视频压缩完成")
}

// Example_audioExtraction 音频提取示例
func Example_audioExtraction() {
	ff, err := ffmpeg.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// 基本音频提取
	err = ff.ExtractAudioToMP3("video.mp4", "audio.mp3")
	if err != nil {
		log.Printf("MP3提取失败: %v", err)
		return
	}

	// 高质量音频提取
	err = ff.ExtractAudioHighQuality("video.mp4", "high_quality.flac")
	if err != nil {
		log.Printf("高质量提取失败: %v", err)
		return
	}

	// 提取音频片段
	err = ff.ExtractAudioSegment("video.mp4", "segment.mp3", "00:01:00", "00:02:00")
	if err != nil {
		log.Printf("片段提取失败: %v", err)
		return
	}

	fmt.Println("音频提取完成")
}

// Example_videoInfo 视频信息获取示例
func Example_videoInfo() {
	ff, err := ffmpeg.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// 获取完整视频信息
	info, err := ff.GetVideoInfo("sample.mp4")
	if err != nil {
		log.Printf("获取信息失败: %v", err)
		return
	}

	fmt.Printf("文件名: %s\n", info.Filename)
	fmt.Printf("时长: %s\n", ffmpeg.FormatDuration(info.Duration))
	fmt.Printf("分辨率: %dx%d\n", info.Width, info.Height)
	fmt.Printf("视频编码: %s\n", info.VideoCodec)
	fmt.Printf("音频编码: %s\n", info.AudioCodec)
	fmt.Printf("比特率: %d kbps\n", info.Bitrate)
	fmt.Printf("帧率: %.2f fps\n", info.FrameRate)
	fmt.Printf("文件大小: %s\n", ffmpeg.FormatFileSize(info.FileSize))
	fmt.Printf("格式: %s\n", info.Format)

	// 检查文件类型
	if ff.IsVideoFile("sample.mp4") {
		fmt.Println("这是一个视频文件")
	}

	if ff.HasAudioStream("sample.mp4") {
		fmt.Println("包含音频流")
	}
}

// Example_videoEditing 视频编辑示例
func Example_videoEditing() {
	ff, err := ffmpeg.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// 视频裁剪
	err = ff.CropTimeRange("input.mp4", "cropped.mp4", "00:01:00", "00:03:00")
	if err != nil {
		log.Printf("裁剪失败: %v", err)
		return
	}

	// 视频合并
	inputFiles := []string{"part1.mp4", "part2.mp4", "part3.mp4"}
	err = ff.Merge(inputFiles, "merged.mp4", nil)
	if err != nil {
		log.Printf("合并失败: %v", err)
		return
	}

	// 视频截图
	err = ff.ScreenshotAtTime("input.mp4", "thumbnail.jpg", "00:00:30")
	if err != nil {
		log.Printf("截图失败: %v", err)
		return
	}

	// 多点截图
	timePoints := []string{"00:00:10", "00:01:00", "00:02:00", "00:03:00"}
	err = ff.ScreenshotMultiple("input.mp4", timePoints, "./thumbnails/")
	if err != nil {
		log.Printf("多点截图失败: %v", err)
		return
	}

	fmt.Println("视频编辑完成")
}

// Example_batchProcessing 批量处理示例
func Example_batchProcessing() {
	ff, err := ffmpeg.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	inputFiles := []string{
		"video1.avi",
		"video2.mov",
		"video3.mkv",
	}

	// 批量转换为MP4
	convertOptions := &ffmpeg.ConvertOptions{
		OutputFormat: "mp4",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		Quality:      "medium",
	}
	err = ff.BatchConvert(inputFiles, "./converted/", convertOptions)
	if err != nil {
		log.Printf("批量转换失败: %v", err)
		return
	}

	// 批量压缩
	compressOptions := &ffmpeg.CompressOptions{
		Width:      1280,
		Height:     720,
		CRF:        25,
		Preset:     "fast",
		VideoCodec: "libx264",
	}
	err = ff.BatchCompress(inputFiles, "./compressed/", compressOptions)
	if err != nil {
		log.Printf("批量压缩失败: %v", err)
		return
	}

	// 批量音频提取
	audioOptions := &ffmpeg.ExtractAudioOptions{
		Format:  "mp3",
		Bitrate: "192k",
	}
	err = ff.BatchExtractAudio(inputFiles, "./audio/", audioOptions)
	if err != nil {
		log.Printf("批量音频提取失败: %v", err)
		return
	}

	fmt.Println("批量处理完成")
}

// customLogger 自定义日志记录器实现
type customLogger struct{}

func (l *customLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("[自定义INFO] "+msg+"\n", args...)
}

func (l *customLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("[自定义ERROR] "+msg+"\n", args...)
}

func (l *customLogger) Debug(msg string, args ...interface{}) {
	fmt.Printf("[自定义DEBUG] "+msg+"\n", args...)
}

// Example_customLogger 自定义日志示例
func Example_customLogger() {

	// 使用自定义日志创建FFmpeg实例
	config := &ffmpeg.Config{
		Logger:  &customLogger{},
		Timeout: 10 * time.Minute,
	}

	ff, err := ffmpeg.New(config)
	if err != nil {
		log.Fatal(err)
	}

	// 执行操作，将使用自定义日志
	err = ff.ConvertToMP4("input.avi", "output.mp4")
	if err != nil {
		log.Printf("转换失败: %v", err)
	}
}

// Example_errorHandling 错误处理示例
func Example_errorHandling() {
	ff, err := ffmpeg.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ff.Convert("nonexistent.avi", "output.mp4", nil)
	if err != nil {
		// 检查具体错误类型
		if ffmpegErr, ok := err.(*ffmpeg.Error); ok {
			switch ffmpegErr.Code {
			case ffmpeg.ErrFileNotFound:
				fmt.Println("错误：输入文件不存在")
			case ffmpeg.ErrFFmpegNotFound:
				fmt.Println("错误：未找到FFmpeg可执行文件")
			case ffmpeg.ErrExecutionFailed:
				fmt.Println("错误：FFmpeg执行失败")
			case ffmpeg.ErrTimeout:
				fmt.Println("错误：操作超时")
			case ffmpeg.ErrUnsupportedFormat:
				fmt.Println("错误：不支持的格式")
			case ffmpeg.ErrInvalidOptions:
				fmt.Println("错误：无效的选项")
			default:
				fmt.Printf("未知错误: %v\n", err)
			}
		} else {
			fmt.Printf("其他错误: %v\n", err)
		}
	}
}
