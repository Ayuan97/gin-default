package ffmpeg

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// VideoInfo 视频文件信息
type VideoInfo struct {
	Filename   string        `json:"filename"`    // 文件名
	Duration   time.Duration `json:"duration"`    // 视频时长
	Width      int           `json:"width"`       // 视频宽度
	Height     int           `json:"height"`      // 视频高度
	Bitrate    int           `json:"bitrate"`     // 比特率 (kbps)
	FrameRate  float64       `json:"frame_rate"`  // 帧率
	VideoCodec string        `json:"video_codec"` // 视频编码格式
	AudioCodec string        `json:"audio_codec"` // 音频编码格式
	FileSize   int64         `json:"file_size"`   // 文件大小 (bytes)
	Format     string        `json:"format"`      // 容器格式
}

// ConvertOptions 视频格式转换选项
type ConvertOptions struct {
	OutputFormat string            // 输出格式 (mp4, avi, mov, mkv等)
	VideoCodec   string            // 视频编码器 (libx264, libx265, copy等)
	AudioCodec   string            // 音频编码器 (aac, mp3, copy等)
	Quality      string            // 质量设置 (high, medium, low 或 CRF值)
	CustomArgs   []string          // 自定义FFmpeg参数
	Metadata     map[string]string // 元数据信息
}

// CompressOptions 视频压缩选项
type CompressOptions struct {
	Width      int      // 输出宽度，0表示保持原始比例
	Height     int      // 输出高度，0表示保持原始比例
	Bitrate    string   // 目标比特率 (如 "1000k", "2M")
	FrameRate  float64  // 目标帧率，0表示保持原始帧率
	CRF        int      // 恒定质量因子 (0-51，越小质量越好)
	Preset     string   // 编码预设 (ultrafast, fast, medium, slow, veryslow)
	VideoCodec string   // 视频编码器
	AudioCodec string   // 音频编码器
	CustomArgs []string // 自定义FFmpeg参数
}

// ExtractAudioOptions 音频提取选项
type ExtractAudioOptions struct {
	Format     string   // 输出音频格式 (mp3, aac, wav, flac等)
	Bitrate    string   // 音频比特率 (如 "128k", "320k")
	SampleRate int      // 采样率 (如 44100, 48000)
	Channels   int      // 声道数 (1=单声道, 2=立体声)
	StartTime  string   // 开始时间 (如 "00:01:30")
	Duration   string   // 持续时间 (如 "00:02:00")
	CustomArgs []string // 自定义FFmpeg参数
}

// CropOptions 视频裁剪选项
type CropOptions struct {
	StartTime  string   // 开始时间 (如 "00:01:30" 或 "90")
	Duration   string   // 持续时间 (如 "00:02:00" 或 "120")
	EndTime    string   // 结束时间 (如 "00:03:30")，与Duration二选一
	CustomArgs []string // 自定义FFmpeg参数
}

// MergeOptions 视频合并选项
type MergeOptions struct {
	Method     string   // 合并方法 ("concat" 或 "filter")
	CustomArgs []string // 自定义FFmpeg参数
}

// ScreenshotOptions 截图选项
type ScreenshotOptions struct {
	Time       string   // 截图时间点 (如 "00:01:30" 或 "90")
	Width      int      // 截图宽度，0表示保持原始尺寸
	Height     int      // 截图高度，0表示保持原始尺寸
	Quality    int      // 截图质量 (1-31，数值越小质量越好)
	Format     string   // 输出格式 (jpg, png, bmp等)
	CustomArgs []string // 自定义FFmpeg参数
}

// ProgressCallback 进度回调函数类型
type ProgressCallback func(progress float64, currentTime time.Duration, totalTime time.Duration)

// Error 自定义错误类型
type Error struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// ErrorCode 错误代码类型
type ErrorCode string

const (
	// ErrInvalidInput 无效输入错误
	ErrInvalidInput ErrorCode = "INVALID_INPUT"
	// ErrFileNotFound 文件未找到错误
	ErrFileNotFound ErrorCode = "FILE_NOT_FOUND"
	// ErrFFmpegNotFound FFmpeg未找到错误
	ErrFFmpegNotFound ErrorCode = "FFMPEG_NOT_FOUND"
	// ErrExecutionFailed 执行失败错误
	ErrExecutionFailed ErrorCode = "EXECUTION_FAILED"
	// ErrTimeout 超时错误
	ErrTimeout ErrorCode = "TIMEOUT"
	// ErrUnsupportedFormat 不支持的格式错误
	ErrUnsupportedFormat ErrorCode = "UNSUPPORTED_FORMAT"
	// ErrInvalidOptions 无效选项错误
	ErrInvalidOptions ErrorCode = "INVALID_OPTIONS"
)

// NewError 创建新的错误
func NewError(code ErrorCode, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// parseTime 解析时间字符串为秒数
func parseTime(timeStr string) (float64, error) {
	if timeStr == "" {
		return 0, nil
	}

	// 如果是纯数字，直接返回秒数
	if seconds, err := strconv.ParseFloat(timeStr, 64); err == nil {
		return seconds, nil
	}

	// 解析 HH:MM:SS 或 MM:SS 格式
	parts := strings.Split(timeStr, ":")
	var hours, minutes, seconds float64
	var err error

	switch len(parts) {
	case 1:
		// 只有秒数
		seconds, err = strconv.ParseFloat(parts[0], 64)
	case 2:
		// MM:SS 格式
		minutes, err = strconv.ParseFloat(parts[0], 64)
		if err == nil {
			seconds, err = strconv.ParseFloat(parts[1], 64)
		}
	case 3:
		// HH:MM:SS 格式
		hours, err = strconv.ParseFloat(parts[0], 64)
		if err == nil {
			minutes, err = strconv.ParseFloat(parts[1], 64)
		}
		if err == nil {
			seconds, err = strconv.ParseFloat(parts[2], 64)
		}
	default:
		return 0, fmt.Errorf("无效的时间格式: %s", timeStr)
	}

	if err != nil {
		return 0, fmt.Errorf("解析时间失败: %s", timeStr)
	}

	return hours*3600 + minutes*60 + seconds, nil
}

// formatTime 将秒数格式化为时间字符串
// 注意：此函数暂时未使用，但保留以备将来使用
func formatTime(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := int(seconds) % 3600 / 60
	secs := int(seconds) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

// validateInputFile 验证输入文件
func validateInputFile(path string) error {
	if path == "" {
		return NewError(ErrInvalidInput, "输入文件路径不能为空", nil)
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return NewError(ErrFileNotFound, fmt.Sprintf("输入文件不存在: %s", path), err)
		}
		return NewError(ErrInvalidInput, fmt.Sprintf("无法访问输入文件: %s", path), err)
	}

	return nil
}

// validateOutputFile 验证输出文件路径
func validateOutputFile(path string) error {
	if path == "" {
		return NewError(ErrInvalidInput, "输出文件路径不能为空", nil)
	}

	// 检查输出目录是否存在
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return NewError(ErrInvalidInput, fmt.Sprintf("输出目录不存在: %s", dir), err)
		}
		return NewError(ErrInvalidInput, fmt.Sprintf("无法访问输出目录: %s", dir), err)
	}

	return nil
}
