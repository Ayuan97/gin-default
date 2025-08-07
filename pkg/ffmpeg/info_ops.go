package ffmpeg

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GetVideoInfo 获取视频文件信息
func (f *FFmpeg) GetVideoInfo(inputPath string) (*VideoInfo, error) {
	return f.GetVideoInfoWithContext(context.Background(), inputPath)
}

// GetVideoInfoWithContext 带上下文的获取视频文件信息
func (f *FFmpeg) GetVideoInfoWithContext(ctx context.Context, inputPath string) (*VideoInfo, error) {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return nil, err
	}

	// 使用ffprobe获取详细信息
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		inputPath,
	}

	// 使用ffprobe而不是ffmpeg
	probePath := strings.Replace(f.execPath, "ffmpeg", "ffprobe", 1)
	if _, err := os.Stat(probePath); err != nil {
		// 如果ffprobe不存在，回退到使用ffmpeg
		return f.getVideoInfoWithFFmpeg(ctx, inputPath)
	}

	// 执行ffprobe命令
	cmd := exec.CommandContext(ctx, probePath, args...)
	output, err := cmd.Output()
	if err != nil {
		f.logger.Error("ffprobe命令执行失败: %v", err)
		// 回退到使用ffmpeg
		return f.getVideoInfoWithFFmpeg(ctx, inputPath)
	}

	// 解析JSON输出
	var probeResult struct {
		Format struct {
			Filename   string `json:"filename"`
			Duration   string `json:"duration"`
			Size       string `json:"size"`
			BitRate    string `json:"bit_rate"`
			FormatName string `json:"format_name"`
		} `json:"format"`
		Streams []struct {
			CodecType  string `json:"codec_type"`
			CodecName  string `json:"codec_name"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			RFrameRate string `json:"r_frame_rate"`
			BitRate    string `json:"bit_rate"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(output, &probeResult); err != nil {
		f.logger.Error("解析ffprobe输出失败: %v", err)
		return f.getVideoInfoWithFFmpeg(ctx, inputPath)
	}

	// 构建VideoInfo结构
	info := &VideoInfo{
		Filename: filepath.Base(inputPath),
		Format:   probeResult.Format.FormatName,
	}

	// 解析时长
	if duration, err := strconv.ParseFloat(probeResult.Format.Duration, 64); err == nil {
		info.Duration = time.Duration(duration * float64(time.Second))
	}

	// 解析文件大小
	if size, err := strconv.ParseInt(probeResult.Format.Size, 10, 64); err == nil {
		info.FileSize = size
	}

	// 解析比特率
	if bitrate, err := strconv.Atoi(probeResult.Format.BitRate); err == nil {
		info.Bitrate = bitrate / 1000 // 转换为kbps
	}

	// 解析流信息
	for _, stream := range probeResult.Streams {
		if stream.CodecType == "video" {
			info.VideoCodec = stream.CodecName
			info.Width = stream.Width
			info.Height = stream.Height

			// 解析帧率
			if stream.RFrameRate != "" {
				if frameRate := parseFrameRate(stream.RFrameRate); frameRate > 0 {
					info.FrameRate = frameRate
				}
			}
		} else if stream.CodecType == "audio" && info.AudioCodec == "" {
			info.AudioCodec = stream.CodecName
		}
	}

	return info, nil
}

// getVideoInfoWithFFmpeg 使用ffmpeg获取视频信息（回退方法）
func (f *FFmpeg) getVideoInfoWithFFmpeg(ctx context.Context, inputPath string) (*VideoInfo, error) {
	args := []string{
		"-i", inputPath,
		"-f", "null",
		"-",
	}

	output, err := f.executeCommand(ctx, args)
	if err != nil {
		// FFmpeg在这种情况下会返回错误，但输出包含我们需要的信息
		output = []byte(err.Error())
	}

	return parseFFmpegOutput(string(output), inputPath)
}

// parseFFmpegOutput 解析FFmpeg输出获取视频信息
func parseFFmpegOutput(output, inputPath string) (*VideoInfo, error) {
	info := &VideoInfo{
		Filename: filepath.Base(inputPath),
	}

	// 获取文件大小
	if stat, err := os.Stat(inputPath); err == nil {
		info.FileSize = stat.Size()
	}

	// 解析时长
	durationRegex := regexp.MustCompile(`Duration: (\d{2}):(\d{2}):(\d{2})\.(\d{2})`)
	if matches := durationRegex.FindStringSubmatch(output); len(matches) == 5 {
		hours, _ := strconv.Atoi(matches[1])
		minutes, _ := strconv.Atoi(matches[2])
		seconds, _ := strconv.Atoi(matches[3])
		centiseconds, _ := strconv.Atoi(matches[4])

		totalSeconds := hours*3600 + minutes*60 + seconds
		info.Duration = time.Duration(totalSeconds)*time.Second + time.Duration(centiseconds*10)*time.Millisecond
	}

	// 解析比特率
	bitrateRegex := regexp.MustCompile(`bitrate: (\d+) kb/s`)
	if matches := bitrateRegex.FindStringSubmatch(output); len(matches) == 2 {
		if bitrate, err := strconv.Atoi(matches[1]); err == nil {
			info.Bitrate = bitrate
		}
	}

	// 解析视频流信息
	videoRegex := regexp.MustCompile(`Stream #\d+:\d+.*?: Video: (\w+).*?, (\d+)x(\d+).*?, (\d+(?:\.\d+)?) fps`)
	if matches := videoRegex.FindStringSubmatch(output); len(matches) == 5 {
		info.VideoCodec = matches[1]
		info.Width, _ = strconv.Atoi(matches[2])
		info.Height, _ = strconv.Atoi(matches[3])
		info.FrameRate, _ = strconv.ParseFloat(matches[4], 64)
	}

	// 解析音频流信息
	audioRegex := regexp.MustCompile(`Stream #\d+:\d+.*?: Audio: (\w+)`)
	if matches := audioRegex.FindStringSubmatch(output); len(matches) == 2 {
		info.AudioCodec = matches[1]
	}

	// 尝试从文件扩展名推断格式
	ext := strings.ToLower(filepath.Ext(inputPath))
	if ext != "" {
		info.Format = strings.TrimPrefix(ext, ".")
	}

	return info, nil
}

// parseFrameRate 解析帧率字符串
func parseFrameRate(frameRateStr string) float64 {
	parts := strings.Split(frameRateStr, "/")
	if len(parts) == 2 {
		numerator, err1 := strconv.ParseFloat(parts[0], 64)
		denominator, err2 := strconv.ParseFloat(parts[1], 64)
		if err1 == nil && err2 == nil && denominator != 0 {
			return numerator / denominator
		}
	}

	// 如果不是分数格式，尝试直接解析
	if frameRate, err := strconv.ParseFloat(frameRateStr, 64); err == nil {
		return frameRate
	}

	return 0
}

// GetVideoDuration 获取视频时长
func (f *FFmpeg) GetVideoDuration(inputPath string) (time.Duration, error) {
	info, err := f.GetVideoInfo(inputPath)
	if err != nil {
		return 0, err
	}
	return info.Duration, nil
}

// GetVideoResolution 获取视频分辨率
func (f *FFmpeg) GetVideoResolution(inputPath string) (width, height int, err error) {
	info, err := f.GetVideoInfo(inputPath)
	if err != nil {
		return 0, 0, err
	}
	return info.Width, info.Height, nil
}

// GetVideoCodec 获取视频编码格式
func (f *FFmpeg) GetVideoCodec(inputPath string) (string, error) {
	info, err := f.GetVideoInfo(inputPath)
	if err != nil {
		return "", err
	}
	return info.VideoCodec, nil
}

// GetAudioCodec 获取音频编码格式
func (f *FFmpeg) GetAudioCodec(inputPath string) (string, error) {
	info, err := f.GetVideoInfo(inputPath)
	if err != nil {
		return "", err
	}
	return info.AudioCodec, nil
}

// GetVideoFrameRate 获取视频帧率
func (f *FFmpeg) GetVideoFrameRate(inputPath string) (float64, error) {
	info, err := f.GetVideoInfo(inputPath)
	if err != nil {
		return 0, err
	}
	return info.FrameRate, nil
}

// GetVideoBitrate 获取视频比特率
func (f *FFmpeg) GetVideoBitrate(inputPath string) (int, error) {
	info, err := f.GetVideoInfo(inputPath)
	if err != nil {
		return 0, err
	}
	return info.Bitrate, nil
}

// GetFileSize 获取文件大小
func (f *FFmpeg) GetFileSize(inputPath string) (int64, error) {
	info, err := f.GetVideoInfo(inputPath)
	if err != nil {
		return 0, err
	}
	return info.FileSize, nil
}

// IsVideoFile 检查文件是否为视频文件
func (f *FFmpeg) IsVideoFile(inputPath string) bool {
	info, err := f.GetVideoInfo(inputPath)
	if err != nil {
		return false
	}
	return info.VideoCodec != ""
}

// HasAudioStream 检查视频是否包含音频流
func (f *FFmpeg) HasAudioStream(inputPath string) bool {
	info, err := f.GetVideoInfo(inputPath)
	if err != nil {
		return false
	}
	return info.AudioCodec != ""
}
