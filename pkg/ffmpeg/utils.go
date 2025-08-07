package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// SupportedVideoFormats 支持的视频格式列表
var SupportedVideoFormats = []string{
	"mp4", "avi", "mov", "mkv", "wmv", "flv", "webm", "m4v", "3gp", "ogv", "ts", "mts", "m2ts",
}

// SupportedAudioFormats 支持的音频格式列表
var SupportedAudioFormats = []string{
	"mp3", "aac", "wav", "flac", "ogg", "m4a", "wma", "ac3", "dts", "opus",
}

// SupportedImageFormats 支持的图片格式列表
var SupportedImageFormats = []string{
	"jpg", "jpeg", "png", "bmp", "gif", "tiff", "webp",
}

// IsVideoFormat 检查文件扩展名是否为支持的视频格式
func IsVideoFormat(filename string) bool {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	for _, format := range SupportedVideoFormats {
		if ext == format {
			return true
		}
	}
	return false
}

// IsAudioFormat 检查文件扩展名是否为支持的音频格式
func IsAudioFormat(filename string) bool {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	for _, format := range SupportedAudioFormats {
		if ext == format {
			return true
		}
	}
	return false
}

// IsImageFormat 检查文件扩展名是否为支持的图片格式
func IsImageFormat(filename string) bool {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	for _, format := range SupportedImageFormats {
		if ext == format {
			return true
		}
	}
	return false
}

// GetFileExtension 获取文件扩展名（不包含点）
func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext != "" {
		return strings.ToLower(strings.TrimPrefix(ext, "."))
	}
	return ""
}

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// GetTempDir 获取临时目录路径
func GetTempDir() string {
	tempDir := os.TempDir()
	ffmpegTempDir := filepath.Join(tempDir, "ffmpeg-go")
	EnsureDir(ffmpegTempDir)
	return ffmpegTempDir
}

// GenerateTempFilename 生成临时文件名
func GenerateTempFilename(prefix, extension string) string {
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%s_%d.%s", prefix, timestamp, extension)
	return filepath.Join(GetTempDir(), filename)
}

// CleanupTempFiles 清理临时文件
func CleanupTempFiles(pattern string) error {
	tempDir := GetTempDir()
	matches, err := filepath.Glob(filepath.Join(tempDir, pattern))
	if err != nil {
		return err
	}

	for _, match := range matches {
		if err := os.Remove(match); err != nil {
			// 记录错误但继续清理其他文件
			fmt.Printf("清理临时文件失败: %s, 错误: %v\n", match, err)
		}
	}

	return nil
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() map[string]string {
	return map[string]string{
		"os":            runtime.GOOS,
		"arch":          runtime.GOARCH,
		"go_version":    runtime.Version(),
		"num_cpu":       fmt.Sprintf("%d", runtime.NumCPU()),
		"num_goroutine": fmt.Sprintf("%d", runtime.NumGoroutine()),
	}
}

// FormatFileSize 格式化文件大小为可读的格式
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration 格式化时长为可读的格式
func FormatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// ParseDuration 解析时长字符串
func ParseDuration(durationStr string) (time.Duration, error) {
	// 支持多种格式: "1h30m", "90m", "5400s", "01:30:00", "90:00"

	// 尝试解析标准Go时长格式
	if duration, err := time.ParseDuration(durationStr); err == nil {
		return duration, nil
	}

	// 尝试解析时间格式
	timeSeconds, err := parseTime(durationStr)
	if err != nil {
		return 0, err
	}

	return time.Duration(timeSeconds * float64(time.Second)), nil
}

// ValidateTimeFormat 验证时间格式
func ValidateTimeFormat(timeStr string) error {
	if timeStr == "" {
		return nil
	}

	_, err := parseTime(timeStr)
	return err
}

// GetVideoQualityPresets 获取视频质量预设
func GetVideoQualityPresets() map[string]map[string]any {
	return map[string]map[string]any{
		"ultra_low": {
			"crf":        35,
			"preset":     "ultrafast",
			"resolution": "480x270",
			"bitrate":    "200k",
		},
		"low": {
			"crf":        30,
			"preset":     "fast",
			"resolution": "854x480",
			"bitrate":    "500k",
		},
		"medium": {
			"crf":        23,
			"preset":     "medium",
			"resolution": "1280x720",
			"bitrate":    "1500k",
		},
		"high": {
			"crf":        18,
			"preset":     "slow",
			"resolution": "1920x1080",
			"bitrate":    "3000k",
		},
		"ultra_high": {
			"crf":        15,
			"preset":     "veryslow",
			"resolution": "3840x2160",
			"bitrate":    "8000k",
		},
	}
}

// GetAudioQualityPresets 获取音频质量预设
func GetAudioQualityPresets() map[string]map[string]any {
	return map[string]map[string]any{
		"low": {
			"bitrate":     "64k",
			"sample_rate": 22050,
			"channels":    1,
		},
		"medium": {
			"bitrate":     "128k",
			"sample_rate": 44100,
			"channels":    2,
		},
		"high": {
			"bitrate":     "192k",
			"sample_rate": 44100,
			"channels":    2,
		},
		"ultra_high": {
			"bitrate":     "320k",
			"sample_rate": 48000,
			"channels":    2,
		},
	}
}

// GetCodecInfo 获取编码器信息
func GetCodecInfo() map[string]map[string]string {
	return map[string]map[string]string{
		"video": {
			"libx264":    "H.264/AVC - 最广泛支持的视频编码器",
			"libx265":    "H.265/HEVC - 更高效的压缩，文件更小",
			"libvpx-vp9": "VP9 - Google开发的开源编码器",
			"libvpx":     "VP8 - 较老的VP编码器",
			"libaom-av1": "AV1 - 最新的开源编码器，压缩效率最高",
			"copy":       "流复制 - 不重新编码，速度最快",
		},
		"audio": {
			"aac":        "AAC - 高质量音频编码器",
			"libmp3lame": "MP3 - 最广泛支持的音频格式",
			"libvorbis":  "Vorbis - 开源音频编码器",
			"flac":       "FLAC - 无损音频编码器",
			"pcm_s16le":  "PCM - 未压缩音频",
			"copy":       "流复制 - 不重新编码",
		},
	}
}

// CheckFFmpegFeatures 检查FFmpeg支持的功能
func (f *FFmpeg) CheckFFmpegFeatures() (map[string]bool, error) {
	features := map[string]bool{
		"libx264":    false,
		"libx265":    false,
		"libvpx":     false,
		"libvpx-vp9": false,
		"libaom-av1": false,
		"libmp3lame": false,
		"aac":        false,
		"libvorbis":  false,
		"flac":       false,
	}

	// 执行ffmpeg -encoders命令获取支持的编码器列表
	output, err := f.executeCommand(context.TODO(), []string{"-encoders"})
	if err != nil {
		return features, err
	}

	outputStr := string(output)
	for codec := range features {
		if strings.Contains(outputStr, codec) {
			features[codec] = true
		}
	}

	return features, nil
}

// GetFFmpegVersion 获取FFmpeg版本信息
func (f *FFmpeg) GetFFmpegVersion() (string, error) {
	output, err := f.executeCommand(context.TODO(), []string{"-version"})
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return lines[0], nil
	}

	return "", fmt.Errorf("无法解析FFmpeg版本信息")
}
