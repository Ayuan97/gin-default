package ffmpeg

import (
	"context"
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

// CancelFunc 取消操作函数类型
type CancelFunc func()

// VideoEditOperation 视频编辑操作接口
type VideoEditOperation interface {
	Execute(ctx context.Context, ffmpeg *FFmpeg) error
	GetDescription() string
	EstimateDuration() time.Duration
}

// Timeline 时间轴结构，用于管理多媒体素材的时间安排
type Timeline struct {
	Duration time.Duration                 // 总时长
	Tracks   map[string][]*TimelineElement // 轨道映射（video, audio, overlay等）
}

// TimelineElement 时间轴元素
type TimelineElement struct {
	StartTime time.Duration // 开始时间
	Duration  time.Duration // 持续时间
	Source    string        // 源文件路径
	Type      ElementType   // 元素类型
	Options   interface{}   // 特定选项
}

// ElementType 元素类型
type ElementType string

const (
	ElementTypeVideo     ElementType = "video"     // 视频元素
	ElementTypeAudio     ElementType = "audio"     // 音频元素
	ElementTypeImage     ElementType = "image"     // 图片元素
	ElementTypeWatermark ElementType = "watermark" // 水印元素
	ElementTypeText      ElementType = "text"      // 文字元素
)

// CropDimensions 裁剪尺寸选项
type CropDimensions struct {
	X      int // 裁剪起始X坐标
	Y      int // 裁剪起始Y坐标
	Width  int // 裁剪宽度
	Height int // 裁剪高度
}

// WatermarkOptions 水印选项
type WatermarkOptions struct {
	ImagePath string  // 水印图片路径
	X         int     // X坐标位置
	Y         int     // Y坐标位置
	Scale     float64 // 缩放比例 (0.1-1.0)
	Opacity   float64 // 透明度 (0.0-1.0)
}

// AudioMixOptions 音频混合选项
type AudioMixOptions struct {
	BackgroundPath string  // 背景音频路径
	Volume         float64 // 音量调节 (0.0-2.0)
	FadeIn         string  // 淡入时间
	FadeOut        string  // 淡出时间
	Loop           bool    // 是否循环
}

// FrameOperation 帧操作类型
type FrameOperation string

const (
	FrameInsert  FrameOperation = "insert"  // 插入帧
	FrameDelete  FrameOperation = "delete"  // 删除帧
	FrameReplace FrameOperation = "replace" // 替换帧
)

// FrameEditOptions 帧编辑选项
type FrameEditOptions struct {
	Operation   FrameOperation // 操作类型
	FrameNumber int            // 帧号
	ImagePath   string         // 图片路径（用于插入或替换）
	Count       int            // 操作帧数
}

// SubtitleStyle 字幕样式
type SubtitleStyle struct {
	FontSize   int    // 字体大小
	FontColor  string // 字体颜色
	FontFamily string // 字体族
	Position   string // 位置 ("bottom", "top", "center")
	Alignment  string // 对齐方式 ("left", "center", "right")
	Outline    int    // 描边宽度
	Shadow     int    // 阴影偏移
	MarginV    int    // 垂直边距
	MarginL    int    // 左边距
	MarginR    int    // 右边距
}

// MultiTrackComposer 多轨道合成器配置
type MultiTrackComposer struct {
	VideoTracks   []*VideoTrack   // 视频轨道
	AudioTracks   []*AudioTrack   // 音频轨道
	ImageTracks   []*ImageTrack   // 图片轨道
	TextTracks    []*TextTrack    // 文字轨道
	OverlayTracks []*OverlayTrack // 叠加轨道
	OutputWidth   int             // 输出宽度
	OutputHeight  int             // 输出高度
	OutputFPS     float64         // 输出帧率
	TotalDuration time.Duration   // 总时长
}

// VideoTrack 视频轨道
type VideoTrack struct {
	ID        string        // 轨道ID
	Source    string        // 视频源路径
	StartTime time.Duration // 开始时间
	Duration  time.Duration // 持续时间
	X         int           // X坐标
	Y         int           // Y坐标
	Width     int           // 宽度
	Height    int           // 高度
	Opacity   float64       // 透明度
	ZIndex    int           // 层级
}

// AudioTrack 音频轨道
type AudioTrack struct {
	ID        string        // 轨道ID
	Source    string        // 音频源路径
	StartTime time.Duration // 开始时间
	Duration  time.Duration // 持续时间
	Volume    float64       // 音量
	FadeIn    time.Duration // 淡入时长
	FadeOut   time.Duration // 淡出时长
	Loop      bool          // 是否循环
}

// ImageTrack 图片轨道
type ImageTrack struct {
	ID        string        // 轨道ID
	Source    string        // 图片源路径
	StartTime time.Duration // 开始时间
	Duration  time.Duration // 持续时间
	X         int           // X坐标
	Y         int           // Y坐标
	Width     int           // 宽度
	Height    int           // 高度
	Opacity   float64       // 透明度
	ZIndex    int           // 层级
}

// TextTrack 文字轨道
type TextTrack struct {
	ID         string        // 轨道ID
	Text       string        // 文字内容
	StartTime  time.Duration // 开始时间
	Duration   time.Duration // 持续时间
	X          int           // X坐标
	Y          int           // Y坐标
	FontSize   int           // 字体大小
	FontColor  string        // 字体颜色
	FontFamily string        // 字体族
	Opacity    float64       // 透明度
	ZIndex     int           // 层级
}

// OverlayTrack 叠加轨道
type OverlayTrack struct {
	ID        string        // 轨道ID
	Source    string        // 叠加源路径
	StartTime time.Duration // 开始时间
	Duration  time.Duration // 持续时间
	X         int           // X坐标
	Y         int           // Y坐标
	Width     int           // 宽度
	Height    int           // 高度
	Opacity   float64       // 透明度
	BlendMode string        // 混合模式
	ZIndex    int           // 层级
}

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

// === 高级滤镜系统类型定义 ===

// FilterType 滤镜类型
type FilterType string

const (
	FilterTypeColorGrading FilterType = "color_grading" // 色彩分级
	FilterTypeVintage      FilterType = "vintage"       // 复古风格
	FilterTypeCinematic    FilterType = "cinematic"     // 电影风格
	FilterTypeBeauty       FilterType = "beauty"        // 美颜滤镜
	FilterTypeSharpening   FilterType = "sharpening"    // 锐化
	FilterTypeDenoising    FilterType = "denoising"     // 降噪
	FilterTypeVignette     FilterType = "vignette"      // 暗角效果
	FilterTypeGlow         FilterType = "glow"          // 发光效果
	FilterTypeBloom        FilterType = "bloom"         // 光晕效果
)

// ColorGradingOptions 色彩分级选项
type ColorGradingOptions struct {
	Brightness  float64 // 亮度 (-1.0 到 1.0)
	Contrast    float64 // 对比度 (-1.0 到 1.0)
	Saturation  float64 // 饱和度 (-1.0 到 1.0)
	Hue         float64 // 色调 (-180 到 180)
	Gamma       float64 // 伽马值 (0.1 到 3.0)
	Temperature float64 // 色温 (-100 到 100)
	Tint        float64 // 色调偏移 (-100 到 100)
	Highlights  float64 // 高光 (-100 到 100)
	Shadows     float64 // 阴影 (-100 到 100)
	Whites      float64 // 白色 (-100 到 100)
	Blacks      float64 // 黑色 (-100 到 100)
	Clarity     float64 // 清晰度 (-100 到 100)
	Vibrance    float64 // 自然饱和度 (-100 到 100)
}

// FilterOptions 通用滤镜选项
type FilterOptions struct {
	FilterType FilterType             // 滤镜类型
	Intensity  float64                // 强度 (0.0 到 1.0)
	StartTime  time.Duration          // 开始时间
	Duration   time.Duration          // 持续时间
	Parameters map[string]interface{} // 自定义参数
	CustomArgs []string               // 自定义FFmpeg参数
}

// VintageFilterOptions 复古滤镜选项
type VintageFilterOptions struct {
	Sepia      float64 // 棕褐色调 (0.0 到 1.0)
	Grain      float64 // 胶片颗粒 (0.0 到 1.0)
	Vignette   float64 // 暗角强度 (0.0 到 1.0)
	Fade       float64 // 褪色效果 (0.0 到 1.0)
	Scratches  bool    // 是否添加划痕
	DustSpots  bool    // 是否添加灰尘斑点
	ColorShift float64 // 色彩偏移 (0.0 到 1.0)
	Desaturate float64 // 去饱和度 (0.0 到 1.0)
}

// CinematicFilterOptions 电影风格滤镜选项
type CinematicFilterOptions struct {
	AspectRatio    string  // 宽高比 ("21:9", "16:9", "4:3")
	LetterboxColor string  // 黑边颜色
	ColorGrading   string  // 色彩分级预设 ("teal_orange", "bleach_bypass", "film_noir")
	FilmGrain      float64 // 胶片颗粒 (0.0 到 1.0)
	Bloom          float64 // 光晕效果 (0.0 到 1.0)
	LensFlare      bool    // 镜头光晕
	MotionBlur     float64 // 运动模糊 (0.0 到 1.0)
}

// BeautyFilterOptions 美颜滤镜选项
type BeautyFilterOptions struct {
	SkinSmoothing   float64 // 磨皮强度 (0.0 到 1.0)
	SkinBrightening float64 // 美白强度 (0.0 到 1.0)
	EyeEnhancement  float64 // 眼部增强 (0.0 到 1.0)
	TeethWhitening  float64 // 牙齿美白 (0.0 到 1.0)
	FaceSlimming    float64 // 瘦脸效果 (0.0 到 1.0)
	EyeEnlarging    float64 // 大眼效果 (0.0 到 1.0)
	NoseReshaping   float64 // 鼻子重塑 (0.0 到 1.0)
	LipEnhancement  float64 // 唇部增强 (0.0 到 1.0)
}

// SharpeningOptions 锐化选项
type SharpeningOptions struct {
	Amount    float64 // 锐化强度 (0.0 到 2.0)
	Radius    float64 // 锐化半径 (0.1 到 5.0)
	Threshold float64 // 锐化阈值 (0.0 到 1.0)
	Method    string  // 锐化方法 ("unsharp", "lanczos", "spline")
}

// DenoisingOptions 降噪选项
type DenoisingOptions struct {
	Strength     float64 // 降噪强度 (0.0 到 1.0)
	Method       string  // 降噪方法 ("nlmeans", "bm3d", "hqdn3d")
	TemporalNR   bool    // 时域降噪
	SpatialNR    bool    // 空域降噪
	PreserveEdge bool    // 保护边缘
}

// VignetteOptions 暗角效果选项
type VignetteOptions struct {
	Intensity float64 // 暗角强度 (0.0 到 1.0)
	Size      float64 // 暗角大小 (0.0 到 1.0)
	Softness  float64 // 边缘柔和度 (0.0 到 1.0)
	Shape     string  // 形状 ("circle", "ellipse", "rectangle")
	CenterX   float64 // 中心X坐标 (0.0 到 1.0)
	CenterY   float64 // 中心Y坐标 (0.0 到 1.0)
}

// GlowOptions 发光效果选项
type GlowOptions struct {
	Intensity  float64 // 发光强度 (0.0 到 2.0)
	Radius     float64 // 发光半径 (1.0 到 50.0)
	Threshold  float64 // 发光阈值 (0.0 到 1.0)
	Color      string  // 发光颜色 (hex格式)
	BlendMode  string  // 混合模式 ("screen", "overlay", "soft_light")
	Saturation float64 // 发光饱和度 (0.0 到 2.0)
}

// BloomOptions 光晕效果选项
type BloomOptions struct {
	Intensity  float64 // 光晕强度 (0.0 到 2.0)
	Radius     float64 // 光晕半径 (1.0 到 100.0)
	Threshold  float64 // 光晕阈值 (0.0 到 1.0)
	Iterations int     // 迭代次数 (1 到 10)
	Quality    string  // 质量设置 ("low", "medium", "high")
}

// LUTOptions LUT色彩查找表选项
type LUTOptions struct {
	LUTPath       string  // LUT文件路径
	Intensity     float64 // 应用强度 (0.0 到 1.0)
	Interpolation string  // 插值方法 ("nearest", "linear", "cubic")
	Format        string  // LUT格式 ("cube", "3dl", "dat")
}

// === 转场效果系统类型定义 ===

// AdvancedTransitionType 高级转场类型
type AdvancedTransitionType string

const (
	TransitionFade     AdvancedTransitionType = "fade"     // 淡入淡出
	TransitionDissolve AdvancedTransitionType = "dissolve" // 溶解
	TransitionWipe     AdvancedTransitionType = "wipe"     // 擦除
	TransitionSlide    AdvancedTransitionType = "slide"    // 滑动
	TransitionPush     AdvancedTransitionType = "push"     // 推拉
	TransitionZoom     AdvancedTransitionType = "zoom"     // 缩放
	TransitionRotate   AdvancedTransitionType = "rotate"   // 旋转
	TransitionFlip     AdvancedTransitionType = "flip"     // 翻转
	TransitionCube     AdvancedTransitionType = "cube"     // 立方体
	TransitionSphere   AdvancedTransitionType = "sphere"   // 球体
	TransitionRipple   AdvancedTransitionType = "ripple"   // 波纹
	TransitionMosaic   AdvancedTransitionType = "mosaic"   // 马赛克
	TransitionPixelate AdvancedTransitionType = "pixelate" // 像素化
	TransitionGlitch   AdvancedTransitionType = "glitch"   // 故障效果
	TransitionBurn     AdvancedTransitionType = "burn"     // 燃烧效果
	TransitionShatter  AdvancedTransitionType = "shatter"  // 破碎效果
)

// AdvancedTransitionOptions 高级转场选项
type AdvancedTransitionOptions struct {
	Type       AdvancedTransitionType // 转场类型
	Duration   time.Duration          // 转场时长
	Direction  string                 // 方向 ("left", "right", "up", "down", "center")
	Easing     string                 // 缓动函数 ("linear", "ease_in", "ease_out", "ease_in_out")
	Intensity  float64                // 强度 (0.0 到 1.0)
	Color      string                 // 转场颜色
	Feather    float64                // 羽化程度 (0.0 到 1.0)
	Reverse    bool                   // 是否反向
	CustomArgs []string               // 自定义参数
}

// WipeTransitionOptions 擦除转场选项
type WipeTransitionOptions struct {
	Direction string  // 擦除方向 ("left_to_right", "right_to_left", "top_to_bottom", "bottom_to_top")
	Angle     float64 // 擦除角度 (0 到 360)
	Softness  float64 // 边缘柔和度 (0.0 到 1.0)
	Shape     string  // 擦除形状 ("linear", "radial", "diamond", "heart")
}

// SlideTransitionOptions 滑动转场选项
type SlideTransitionOptions struct {
	Direction string  // 滑动方向
	Distance  float64 // 滑动距离 (0.0 到 2.0)
	Bounce    bool    // 是否有弹跳效果
	Rotation  float64 // 滑动时的旋转角度
}

// ZoomTransitionOptions 缩放转场选项
type ZoomTransitionOptions struct {
	ZoomIn   bool    // true为放大，false为缩小
	CenterX  float64 // 缩放中心X (0.0 到 1.0)
	CenterY  float64 // 缩放中心Y (0.0 到 1.0)
	MaxScale float64 // 最大缩放比例 (1.0 到 10.0)
	Blur     float64 // 缩放时的模糊程度 (0.0 到 1.0)
}

// GlitchTransitionOptions 故障效果转场选项
type GlitchTransitionOptions struct {
	Intensity    float64 // 故障强度 (0.0 到 1.0)
	BlockSize    int     // 故障块大小 (1 到 50)
	ColorShift   float64 // 色彩偏移 (0.0 到 1.0)
	DigitalNoise float64 // 数字噪声 (0.0 到 1.0)
	Scanlines    bool    // 是否添加扫描线
	Distortion   float64 // 扭曲程度 (0.0 到 1.0)
}

// === 高级音频处理类型定义 ===

// AudioEqualizerBand 音频均衡器频段
type AudioEqualizerBand struct {
	Frequency float64 // 频率 (Hz)
	Gain      float64 // 增益 (dB) (-20 到 20)
	Q         float64 // 品质因数 (0.1 到 10)
}

// AudioEqualizerOptions 音频均衡器选项
type AudioEqualizerOptions struct {
	Bands      []AudioEqualizerBand // 频段设置
	Preset     string               // 预设 ("rock", "pop", "classical", "jazz", "vocal")
	MasterGain float64              // 主增益 (dB)
}

// AudioEffectType 音频效果类型
type AudioEffectType string

const (
	AudioEffectReverb     AudioEffectType = "reverb"      // 混响
	AudioEffectEcho       AudioEffectType = "echo"        // 回声
	AudioEffectChorus     AudioEffectType = "chorus"      // 合唱
	AudioEffectFlanger    AudioEffectType = "flanger"     // 镶边
	AudioEffectPhaser     AudioEffectType = "phaser"      // 相位器
	AudioEffectDistortion AudioEffectType = "distortion"  // 失真
	AudioEffectCompressor AudioEffectType = "compressor"  // 压缩器
	AudioEffectLimiter    AudioEffectType = "limiter"     // 限制器
	AudioEffectGate       AudioEffectType = "gate"        // 噪声门
	AudioEffectPitchShift AudioEffectType = "pitch_shift" // 变调
)

// AudioEffectOptions 音频效果选项
type AudioEffectOptions struct {
	EffectType AudioEffectType        // 效果类型
	Intensity  float64                // 强度 (0.0 到 1.0)
	Parameters map[string]interface{} // 效果参数
	StartTime  time.Duration          // 开始时间
	Duration   time.Duration          // 持续时间
}

// ReverbOptions 混响选项
type ReverbOptions struct {
	RoomSize   float64 // 房间大小 (0.0 到 1.0)
	Damping    float64 // 阻尼 (0.0 到 1.0)
	WetLevel   float64 // 湿声级别 (0.0 到 1.0)
	DryLevel   float64 // 干声级别 (0.0 到 1.0)
	PreDelay   float64 // 预延迟 (ms)
	Diffusion  float64 // 扩散 (0.0 到 1.0)
	ReverbType string  // 混响类型 ("hall", "room", "plate", "spring")
}

// CompressorOptions 压缩器选项
type CompressorOptions struct {
	Threshold  float64 // 阈值 (dB)
	Ratio      float64 // 压缩比 (1.0 到 20.0)
	Attack     float64 // 启动时间 (ms)
	Release    float64 // 释放时间 (ms)
	MakeupGain float64 // 补偿增益 (dB)
	KneeWidth  float64 // 拐点宽度 (dB)
}

// === 高级字幕系统类型定义 ===

// SubtitleAnimationType 字幕动画类型
type SubtitleAnimationType string

const (
	SubtitleAnimationFadeIn     SubtitleAnimationType = "fade_in"    // 淡入
	SubtitleAnimationFadeOut    SubtitleAnimationType = "fade_out"   // 淡出
	SubtitleAnimationSlideIn    SubtitleAnimationType = "slide_in"   // 滑入
	SubtitleAnimationSlideOut   SubtitleAnimationType = "slide_out"  // 滑出
	SubtitleAnimationZoomIn     SubtitleAnimationType = "zoom_in"    // 缩放进入
	SubtitleAnimationZoomOut    SubtitleAnimationType = "zoom_out"   // 缩放退出
	SubtitleAnimationTypewriter SubtitleAnimationType = "typewriter" // 打字机效果
	SubtitleAnimationBounce     SubtitleAnimationType = "bounce"     // 弹跳
	SubtitleAnimationRotate     SubtitleAnimationType = "rotate"     // 旋转
	SubtitleAnimationGlow       SubtitleAnimationType = "glow"       // 发光
)

// AdvancedSubtitleOptions 高级字幕选项
type AdvancedSubtitleOptions struct {
	Text              string                // 字幕文本
	StartTime         time.Duration         // 开始时间
	Duration          time.Duration         // 持续时间
	X                 int                   // X坐标
	Y                 int                   // Y坐标
	FontFamily        string                // 字体族
	FontSize          int                   // 字体大小
	FontWeight        string                // 字体粗细 ("normal", "bold", "light")
	FontStyle         string                // 字体样式 ("normal", "italic")
	Color             string                // 字体颜色
	BackgroundColor   string                // 背景颜色
	OutlineColor      string                // 描边颜色
	OutlineWidth      int                   // 描边宽度
	ShadowColor       string                // 阴影颜色
	ShadowOffsetX     int                   // 阴影X偏移
	ShadowOffsetY     int                   // 阴影Y偏移
	ShadowBlur        int                   // 阴影模糊
	Alignment         string                // 对齐方式 ("left", "center", "right")
	VerticalAlign     string                // 垂直对齐 ("top", "middle", "bottom")
	LineSpacing       float64               // 行间距
	LetterSpacing     float64               // 字符间距
	Opacity           float64               // 透明度 (0.0 到 1.0)
	Rotation          float64               // 旋转角度
	ScaleX            float64               // X轴缩放
	ScaleY            float64               // Y轴缩放
	Animation         SubtitleAnimationType // 动画类型
	AnimationDuration time.Duration         // 动画时长
}

// SubtitleTemplate 字幕模板
type SubtitleTemplate struct {
	Name        string                   // 模板名称
	Description string                   // 模板描述
	Style       *AdvancedSubtitleOptions // 样式设置
	Animation   SubtitleAnimationType    // 默认动画
}

// === 高级合成功能类型定义 ===

// ChromaKeyOptions 绿幕抠图选项
type ChromaKeyOptions struct {
	KeyColor         string  // 抠图颜色 (hex格式)
	Tolerance        float64 // 容差 (0.0 到 1.0)
	Softness         float64 // 边缘柔和度 (0.0 到 1.0)
	SpillSuppression float64 // 溢色抑制 (0.0 到 1.0)
	EdgeFeather      float64 // 边缘羽化 (0.0 到 1.0)
	LightWrap        float64 // 光线包裹 (0.0 到 1.0)
	ColorCorrection  bool    // 是否进行色彩校正
	BackgroundPath   string  // 背景视频/图片路径
}

// MaskType 遮罩类型
type MaskType string

const (
	MaskTypeAlpha    MaskType = "alpha"    // Alpha遮罩
	MaskTypeLuma     MaskType = "luma"     // 亮度遮罩
	MaskTypeColor    MaskType = "color"    // 颜色遮罩
	MaskTypeShape    MaskType = "shape"    // 形状遮罩
	MaskTypeGradient MaskType = "gradient" // 渐变遮罩
)

// MaskOptions 遮罩选项
type MaskOptions struct {
	MaskType  MaskType      // 遮罩类型
	MaskPath  string        // 遮罩文件路径
	Invert    bool          // 是否反转遮罩
	Feather   float64       // 羽化程度 (0.0 到 1.0)
	Opacity   float64       // 遮罩透明度 (0.0 到 1.0)
	BlendMode string        // 混合模式
	StartTime time.Duration // 开始时间
	Duration  time.Duration // 持续时间
}

// ParticleType 粒子类型
type ParticleType string

const (
	ParticleTypeSnow     ParticleType = "snow"     // 雪花
	ParticleTypeRain     ParticleType = "rain"     // 雨滴
	ParticleTypeFire     ParticleType = "fire"     // 火焰
	ParticleTypeSmoke    ParticleType = "smoke"    // 烟雾
	ParticleTypeSparkle  ParticleType = "sparkle"  // 闪光
	ParticleTypeBubbles  ParticleType = "bubbles"  // 气泡
	ParticleTypeLeaves   ParticleType = "leaves"   // 落叶
	ParticleTypeStars    ParticleType = "stars"    // 星星
	ParticleTypeHearts   ParticleType = "hearts"   // 爱心
	ParticleTypeConfetti ParticleType = "confetti" // 彩纸
)

// ParticleEffectOptions 粒子效果选项
type ParticleEffectOptions struct {
	ParticleType ParticleType  // 粒子类型
	Count        int           // 粒子数量
	Size         float64       // 粒子大小
	Speed        float64       // 粒子速度
	Direction    float64       // 粒子方向 (角度)
	Spread       float64       // 扩散角度
	Gravity      float64       // 重力影响
	Wind         float64       // 风力影响
	Opacity      float64       // 粒子透明度
	Color        string        // 粒子颜色
	BlendMode    string        // 混合模式
	StartTime    time.Duration // 开始时间
	Duration     time.Duration // 持续时间
	EmissionRate float64       // 发射速率 (粒子/秒)
	LifeTime     float64       // 粒子生命周期 (秒)
}

// MotionGraphicsType 动态图形类型
type MotionGraphicsType string

const (
	MotionGraphicsLowerThird MotionGraphicsType = "lower_third" // 下三分之一标题
	MotionGraphicsCallout    MotionGraphicsType = "callout"     // 标注
	MotionGraphicsProgress   MotionGraphicsType = "progress"    // 进度条
	MotionGraphicsCounter    MotionGraphicsType = "counter"     // 计数器
	MotionGraphicsChart      MotionGraphicsType = "chart"       // 图表
	MotionGraphicsLogo       MotionGraphicsType = "logo"        // Logo动画
)

// MotionGraphicsOptions 动态图形选项
type MotionGraphicsOptions struct {
	GraphicsType MotionGraphicsType     // 图形类型
	Template     string                 // 模板名称
	Text         string                 // 文本内容
	StartTime    time.Duration          // 开始时间
	Duration     time.Duration          // 持续时间
	X            int                    // X坐标
	Y            int                    // Y坐标
	Width        int                    // 宽度
	Height       int                    // 高度
	Color        string                 // 主色调
	AccentColor  string                 // 强调色
	Animation    string                 // 动画类型
	Parameters   map[string]interface{} // 自定义参数
}

// === 新增操作类型定义 ===

// OperationType 操作类型
type OperationType string

const (
	OpTypeColorGrading     OperationType = "color_grading"     // 色彩分级操作
	OpTypeFilter           OperationType = "filter"            // 滤镜操作
	OpTypeTransition       OperationType = "transition"        // 转场操作
	OpTypeAudioEffect      OperationType = "audio_effect"      // 音频效果操作
	OpTypeAdvancedSubtitle OperationType = "advanced_subtitle" // 高级字幕操作
	OpTypeChromaKey        OperationType = "chroma_key"        // 绿幕抠图操作
	OpTypeMask             OperationType = "mask"              // 遮罩操作
	OpTypeParticleEffect   OperationType = "particle_effect"   // 粒子效果操作
	OpTypeMotionGraphics   OperationType = "motion_graphics"   // 动态图形操作
)

// BaseOperation 基础操作结构
type BaseOperation struct {
	Type      OperationType // 操作类型
	StartTime time.Duration // 开始时间
	Duration  time.Duration // 持续时间
}

// ColorGradingOperation 色彩分级操作
type ColorGradingOperation struct {
	BaseOperation
	Options *ColorGradingOptions
}

func (op *ColorGradingOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	// 实现色彩分级操作
	return nil
}

func (op *ColorGradingOperation) GetDescription() string {
	return "色彩分级操作"
}

func (op *ColorGradingOperation) EstimateDuration() time.Duration {
	return 30 * time.Second
}

// FilterOperation 滤镜操作
type FilterOperation struct {
	BaseOperation
	Options *FilterOptions
}

func (op *FilterOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	// 实现滤镜操作
	return nil
}

func (op *FilterOperation) GetDescription() string {
	return fmt.Sprintf("滤镜操作: %s", op.Options.FilterType)
}

func (op *FilterOperation) EstimateDuration() time.Duration {
	return 20 * time.Second
}

// AdvancedTransitionOperation 高级转场操作
type AdvancedTransitionOperation struct {
	BaseOperation
	Options *AdvancedTransitionOptions
}

func (op *AdvancedTransitionOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	// 实现高级转场操作
	return nil
}

func (op *AdvancedTransitionOperation) GetDescription() string {
	return fmt.Sprintf("高级转场操作: %s", op.Options.Type)
}

func (op *AdvancedTransitionOperation) EstimateDuration() time.Duration {
	return 15 * time.Second
}

// AudioEffectOperation 音频效果操作
type AudioEffectOperation struct {
	BaseOperation
	Options *AudioEffectOptions
}

func (op *AudioEffectOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	// 实现音频效果操作
	return nil
}

func (op *AudioEffectOperation) GetDescription() string {
	return fmt.Sprintf("音频效果操作: %s", op.Options.EffectType)
}

func (op *AudioEffectOperation) EstimateDuration() time.Duration {
	return 25 * time.Second
}

// AdvancedSubtitleOperation 高级字幕操作
type AdvancedSubtitleOperation struct {
	BaseOperation
	Options *AdvancedSubtitleOptions
}

func (op *AdvancedSubtitleOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	// 实现高级字幕操作
	return nil
}

func (op *AdvancedSubtitleOperation) GetDescription() string {
	return fmt.Sprintf("高级字幕操作: %s", op.Options.Text)
}

func (op *AdvancedSubtitleOperation) EstimateDuration() time.Duration {
	return 10 * time.Second
}

// ChromaKeyOperation 绿幕抠图操作
type ChromaKeyOperation struct {
	BaseOperation
	Options *ChromaKeyOptions
}

func (op *ChromaKeyOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	// 实现绿幕抠图操作
	return nil
}

func (op *ChromaKeyOperation) GetDescription() string {
	return "绿幕抠图操作"
}

func (op *ChromaKeyOperation) EstimateDuration() time.Duration {
	return 40 * time.Second
}

// MaskOperation 遮罩操作
type MaskOperation struct {
	BaseOperation
	Options *MaskOptions
}

func (op *MaskOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	// 实现遮罩操作
	return nil
}

func (op *MaskOperation) GetDescription() string {
	return fmt.Sprintf("遮罩操作: %s", op.Options.MaskType)
}

func (op *MaskOperation) EstimateDuration() time.Duration {
	return 20 * time.Second
}

// ParticleEffectOperation 粒子效果操作
type ParticleEffectOperation struct {
	BaseOperation
	Options *ParticleEffectOptions
}

func (op *ParticleEffectOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	// 实现粒子效果操作
	return nil
}

func (op *ParticleEffectOperation) GetDescription() string {
	return fmt.Sprintf("粒子效果操作: %s", op.Options.ParticleType)
}

func (op *ParticleEffectOperation) EstimateDuration() time.Duration {
	return 35 * time.Second
}

// MotionGraphicsOperation 动态图形操作
type MotionGraphicsOperation struct {
	BaseOperation
	Options *MotionGraphicsOptions
}

func (op *MotionGraphicsOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	// 实现动态图形操作
	return nil
}

func (op *MotionGraphicsOperation) GetDescription() string {
	return fmt.Sprintf("动态图形操作: %s", op.Options.GraphicsType)
}

func (op *MotionGraphicsOperation) EstimateDuration() time.Duration {
	return 15 * time.Second
}
