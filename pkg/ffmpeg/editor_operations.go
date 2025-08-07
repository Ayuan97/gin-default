// Package ffmpeg 提供视频编辑操作的具体实现
package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CropTimeOperation 时间段裁剪操作
type CropTimeOperation struct {
	StartTime string // 开始时间
	Duration  string // 持续时间
}

func (op *CropTimeOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	// 这里需要实际的输入输出路径，暂时使用临时文件
	tempOutput := generateTempFilePath("crop_time", ".mp4")

	args := []string{
		"-i", "INPUT_PLACEHOLDER", // 将在实际执行时替换
		"-ss", op.StartTime,
		"-t", op.Duration,
		"-c", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *CropTimeOperation) GetDescription() string {
	return fmt.Sprintf("时间段裁剪: %s 开始，持续 %s", op.StartTime, op.Duration)
}

func (op *CropTimeOperation) EstimateDuration() time.Duration {
	return 10 * time.Second // 预估10秒
}

// CropDimensionOperation 尺寸裁剪操作
type CropDimensionOperation struct {
	Dimensions *CropDimensions
}

func (op *CropDimensionOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("crop_dimension", ".mp4")

	cropFilter := fmt.Sprintf("crop=%d:%d:%d:%d",
		op.Dimensions.Width, op.Dimensions.Height,
		op.Dimensions.X, op.Dimensions.Y)

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", cropFilter,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *CropDimensionOperation) GetDescription() string {
	return fmt.Sprintf("尺寸裁剪: %dx%d 从 (%d,%d)",
		op.Dimensions.Width, op.Dimensions.Height,
		op.Dimensions.X, op.Dimensions.Y)
}

func (op *CropDimensionOperation) EstimateDuration() time.Duration {
	return 15 * time.Second
}

// ResizeOperation 调整尺寸操作
type ResizeOperation struct {
	Width  int
	Height int
}

func (op *ResizeOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("resize", ".mp4")

	scaleFilter := fmt.Sprintf("scale=%d:%d", op.Width, op.Height)

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", scaleFilter,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *ResizeOperation) GetDescription() string {
	return fmt.Sprintf("调整尺寸: %dx%d", op.Width, op.Height)
}

func (op *ResizeOperation) EstimateDuration() time.Duration {
	return 20 * time.Second
}

// WatermarkOperation 水印操作
type WatermarkOperation struct {
	Options *WatermarkOptions
}

func (op *WatermarkOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("watermark", ".mp4")

	// 构建水印滤镜
	overlayFilter := fmt.Sprintf("overlay=%d:%d", op.Options.X, op.Options.Y)
	if op.Options.Opacity < 1.0 {
		overlayFilter = fmt.Sprintf("format=rgba,colorchannelmixer=aa=%f[wm];[0:v][wm]overlay=%d:%d",
			op.Options.Opacity, op.Options.X, op.Options.Y)
	}

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-i", op.Options.ImagePath,
		"-filter_complex", overlayFilter,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *WatermarkOperation) GetDescription() string {
	return fmt.Sprintf("添加水印: %s 位置(%d,%d)",
		filepath.Base(op.Options.ImagePath), op.Options.X, op.Options.Y)
}

func (op *WatermarkOperation) EstimateDuration() time.Duration {
	return 25 * time.Second
}

// SeparateAudioOperation 分离音频操作
type SeparateAudioOperation struct {
	AudioOutputPath string
}

func (op *SeparateAudioOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	// 提取音频
	audioArgs := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vn", // 不包含视频
		"-acodec", "copy",
		"-y",
		op.AudioOutputPath,
	}

	if _, err := ffmpeg.executeCommand(ctx, audioArgs); err != nil {
		return fmt.Errorf("提取音频失败: %w", err)
	}

	// 创建无音频视频
	tempOutput := generateTempFilePath("no_audio", ".mp4")
	videoArgs := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-an", // 不包含音频
		"-vcodec", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, videoArgs)
	return err
}

func (op *SeparateAudioOperation) GetDescription() string {
	return fmt.Sprintf("分离音频到: %s", filepath.Base(op.AudioOutputPath))
}

func (op *SeparateAudioOperation) EstimateDuration() time.Duration {
	return 15 * time.Second
}

// AudioMixOperation 音频混合操作
type AudioMixOperation struct {
	Options *AudioMixOptions
}

func (op *AudioMixOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("audio_mix", ".mp4")

	// 构建音频混合滤镜
	var filterComplex string
	if op.Options.Loop {
		filterComplex = fmt.Sprintf("[1:a]aloop=loop=-1:size=2e+09[bg];[0:a][bg]amix=inputs=2:duration=first:dropout_transition=2,volume=%f[a]",
			op.Options.Volume)
	} else {
		filterComplex = fmt.Sprintf("[0:a][1:a]amix=inputs=2:duration=first:dropout_transition=2,volume=%f[a]",
			op.Options.Volume)
	}

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-i", op.Options.BackgroundPath,
		"-filter_complex", filterComplex,
		"-map", "0:v",
		"-map", "[a]",
		"-c:v", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *AudioMixOperation) GetDescription() string {
	return fmt.Sprintf("混合音频: %s (音量: %.1f)",
		filepath.Base(op.Options.BackgroundPath), op.Options.Volume)
}

func (op *AudioMixOperation) EstimateDuration() time.Duration {
	return 30 * time.Second
}

// InsertImageOperation 插入图片操作
type InsertImageOperation struct {
	ImagePath string
	StartTime time.Duration
	Duration  time.Duration
}

func (op *InsertImageOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("insert_image", ".mp4")

	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	filterComplex := fmt.Sprintf("[1:v]scale=iw:ih[img];[0:v][img]overlay=enable='between(t,%f,%f)'",
		startSec, startSec+durationSec)

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-i", op.ImagePath,
		"-filter_complex", filterComplex,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *InsertImageOperation) GetDescription() string {
	return fmt.Sprintf("插入图片: %s (%.1fs-%.1fs)",
		filepath.Base(op.ImagePath), op.StartTime.Seconds(),
		(op.StartTime + op.Duration).Seconds())
}

func (op *InsertImageOperation) EstimateDuration() time.Duration {
	return 20 * time.Second
}

// FrameEditOperation 帧编辑操作
type FrameEditOperation struct {
	Options *FrameEditOptions
}

func (op *FrameEditOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("frame_edit", ".mp4")

	switch op.Options.Operation {
	case FrameInsert:
		return op.executeInsertFrame(ctx, ffmpeg, tempOutput)
	case FrameDelete:
		return op.executeDeleteFrame(ctx, ffmpeg, tempOutput)
	case FrameReplace:
		return op.executeReplaceFrame(ctx, ffmpeg, tempOutput)
	default:
		return fmt.Errorf("不支持的帧操作: %s", op.Options.Operation)
	}
}

func (op *FrameEditOperation) executeInsertFrame(ctx context.Context, ffmpeg *FFmpeg, output string) error {
	// 帧插入的复杂实现，这里简化处理
	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-i", op.Options.ImagePath,
		"-filter_complex", "[0:v][1:v]concat=n=2:v=1:a=0[v]",
		"-map", "[v]",
		"-map", "0:a",
		"-y",
		output,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *FrameEditOperation) executeDeleteFrame(ctx context.Context, ffmpeg *FFmpeg, output string) error {
	// 帧删除实现
	frameTime := float64(op.Options.FrameNumber) / 30.0 // 假设30fps

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", fmt.Sprintf("select='not(between(t,%f,%f))',setpts=N/FRAME_RATE/TB",
			frameTime, frameTime+0.033),
		"-af", "aselect='not(between(t,%f,%f))',asetpts=N/SR/TB",
		"-y",
		output,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *FrameEditOperation) executeReplaceFrame(ctx context.Context, ffmpeg *FFmpeg, output string) error {
	// 帧替换实现
	frameTime := float64(op.Options.FrameNumber) / 30.0

	filterComplex := fmt.Sprintf("[1:v]scale=iw:ih[img];[0:v][img]overlay=enable='between(t,%f,%f)'",
		frameTime, frameTime+0.033)

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-i", op.Options.ImagePath,
		"-filter_complex", filterComplex,
		"-c:a", "copy",
		"-y",
		output,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *FrameEditOperation) GetDescription() string {
	return fmt.Sprintf("帧编辑: %s 帧号%d", op.Options.Operation, op.Options.FrameNumber)
}

func (op *FrameEditOperation) EstimateDuration() time.Duration {
	return 25 * time.Second
}

// TextOperation 文字添加操作
type TextOperation struct {
	Text      string
	X         int
	Y         int
	FontSize  int
	Color     string
	StartTime time.Duration
	Duration  time.Duration
}

func (op *TextOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("text", ".mp4")

	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	drawTextFilter := fmt.Sprintf("drawtext=text='%s':x=%d:y=%d:fontsize=%d:fontcolor=%s:enable='between(t,%f,%f)'",
		op.Text, op.X, op.Y, op.FontSize, op.Color, startSec, startSec+durationSec)

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", drawTextFilter,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *TextOperation) GetDescription() string {
	return fmt.Sprintf("添加文字: %s 位置(%d,%d)", op.Text, op.X, op.Y)
}

func (op *TextOperation) EstimateDuration() time.Duration {
	return 15 * time.Second
}

// FadeOperation 淡入淡出操作
type FadeOperation struct {
	Type     string // "in" 或 "out"
	Duration time.Duration
}

func (op *FadeOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("fade", ".mp4")

	var fadeFilter string
	if op.Type == "in" {
		fadeFilter = fmt.Sprintf("fade=in:0:%d", int(op.Duration.Seconds()*30)) // 假设30fps
	} else {
		fadeFilter = fmt.Sprintf("fade=out:st=0:d=%f", op.Duration.Seconds())
	}

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", fadeFilter,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *FadeOperation) GetDescription() string {
	return fmt.Sprintf("淡%s效果: %.1fs", op.Type, op.Duration.Seconds())
}

func (op *FadeOperation) EstimateDuration() time.Duration {
	return 10 * time.Second
}

// SpeedOperation 速度调整操作
type SpeedOperation struct {
	Factor float64 // 速度因子，>1加速，<1减速
}

func (op *SpeedOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("speed", ".mp4")

	videoFilter := fmt.Sprintf("setpts=%f*PTS", 1.0/op.Factor)
	audioFilter := fmt.Sprintf("atempo=%f", op.Factor)

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", videoFilter,
		"-af", audioFilter,
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *SpeedOperation) GetDescription() string {
	return fmt.Sprintf("速度调整: %.2fx", op.Factor)
}

func (op *SpeedOperation) EstimateDuration() time.Duration {
	return 20 * time.Second
}

// RotateOperation 旋转操作
type RotateOperation struct {
	Angle int // 旋转角度
}

func (op *RotateOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("rotate", ".mp4")

	var rotateFilter string
	switch op.Angle {
	case 90:
		rotateFilter = "transpose=1"
	case 180:
		rotateFilter = "transpose=2,transpose=2"
	case 270:
		rotateFilter = "transpose=2"
	default:
		rotateFilter = fmt.Sprintf("rotate=%f*PI/180", float64(op.Angle))
	}

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", rotateFilter,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *RotateOperation) GetDescription() string {
	return fmt.Sprintf("旋转: %d度", op.Angle)
}

func (op *RotateOperation) EstimateDuration() time.Duration {
	return 15 * time.Second
}

// MirrorOperation 镜像翻转操作
type MirrorOperation struct {
	Horizontal bool // true为水平翻转，false为垂直翻转
}

func (op *MirrorOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("mirror", ".mp4")

	var flipFilter string
	if op.Horizontal {
		flipFilter = "hflip"
	} else {
		flipFilter = "vflip"
	}

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", flipFilter,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *MirrorOperation) GetDescription() string {
	if op.Horizontal {
		return "水平镜像翻转"
	}
	return "垂直镜像翻转"
}

func (op *MirrorOperation) EstimateDuration() time.Duration {
	return 10 * time.Second
}

// BrightnessOperation 亮度调整操作
type BrightnessOperation struct {
	Brightness float64 // 亮度值，-1.0到1.0
}

func (op *BrightnessOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("brightness", ".mp4")

	brightnessFilter := fmt.Sprintf("eq=brightness=%f", op.Brightness)

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", brightnessFilter,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *BrightnessOperation) GetDescription() string {
	return fmt.Sprintf("亮度调整: %.2f", op.Brightness)
}

func (op *BrightnessOperation) EstimateDuration() time.Duration {
	return 12 * time.Second
}

// ContrastOperation 对比度调整操作
type ContrastOperation struct {
	Contrast float64 // 对比度值，-1.0到1.0
}

func (op *ContrastOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("contrast", ".mp4")

	contrastFilter := fmt.Sprintf("eq=contrast=%f", op.Contrast)

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", contrastFilter,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *ContrastOperation) GetDescription() string {
	return fmt.Sprintf("对比度调整: %.2f", op.Contrast)
}

func (op *ContrastOperation) EstimateDuration() time.Duration {
	return 12 * time.Second
}

// BlurOperation 模糊效果操作
type BlurOperation struct {
	Radius float64 // 模糊半径
}

func (op *BlurOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("blur", ".mp4")

	blurFilter := fmt.Sprintf("boxblur=%f", op.Radius)

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", blurFilter,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *BlurOperation) GetDescription() string {
	return fmt.Sprintf("模糊效果: 半径%.1f", op.Radius)
}

func (op *BlurOperation) EstimateDuration() time.Duration {
	return 18 * time.Second
}

// StabilizeOperation 视频防抖操作
type StabilizeOperation struct{}

func (op *StabilizeOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("stabilize", ".mp4")

	// 使用vidstabdetect和vidstabtransform进行防抖
	stabilizeFilter := "vidstabtransform=smoothing=30"

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", stabilizeFilter,
		"-c:a", "copy",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *StabilizeOperation) GetDescription() string {
	return "视频防抖"
}

func (op *StabilizeOperation) EstimateDuration() time.Duration {
	return 45 * time.Second // 防抖处理较慢
}

// ExtractFramesOperation 提取帧操作
type ExtractFramesOperation struct {
	OutputDir string  // 输出目录
	FPS       float64 // 提取帧率
}

func (op *ExtractFramesOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	outputPattern := filepath.Join(op.OutputDir, "frame_%04d.png")

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", fmt.Sprintf("fps=%f", op.FPS),
		"-y",
		outputPattern,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *ExtractFramesOperation) GetDescription() string {
	return fmt.Sprintf("提取帧: %.1f fps 到 %s", op.FPS, op.OutputDir)
}

func (op *ExtractFramesOperation) EstimateDuration() time.Duration {
	return 30 * time.Second
}

// CreateFromImagesOperation 从图片创建视频操作
type CreateFromImagesOperation struct {
	ImagePattern string  // 图片模式，如 "image_%04d.png"
	FPS          float64 // 输出帧率
}

func (op *CreateFromImagesOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("from_images", ".mp4")

	args := []string{
		"-framerate", fmt.Sprintf("%f", op.FPS),
		"-i", op.ImagePattern,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-y",
		tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *CreateFromImagesOperation) GetDescription() string {
	return fmt.Sprintf("从图片创建视频: %s (%.1f fps)", op.ImagePattern, op.FPS)
}

func (op *CreateFromImagesOperation) EstimateDuration() time.Duration {
	return 25 * time.Second
}

// === 多媒体合成操作类型 ===

// AddTrackOperation 添加轨道操作
type AddTrackOperation struct {
	TrackType  string        // 轨道类型: video, audio, image, overlay, text
	SourcePath string        // 源文件路径
	Text       string        // 文字内容（用于text轨道）
	StartTime  time.Duration // 开始时间
	Duration   time.Duration // 持续时间
	X          int           // X坐标
	Y          int           // Y坐标
	Volume     float64       // 音量（用于audio轨道）
	Opacity    float64       // 透明度（用于overlay轨道）
	FontSize   int           // 字体大小（用于text轨道）
	Color      string        // 颜色（用于text轨道）
}

func (op *AddTrackOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("add_track", ".mp4")

	var args []string

	switch op.TrackType {
	case "video":
		args = op.buildVideoTrackArgs(tempOutput)
	case "audio":
		args = op.buildAudioTrackArgs(tempOutput)
	case "image":
		args = op.buildImageTrackArgs(tempOutput)
	case "overlay":
		args = op.buildOverlayTrackArgs(tempOutput)
	case "text":
		args = op.buildTextTrackArgs(tempOutput)
	default:
		return fmt.Errorf("不支持的轨道类型: %s", op.TrackType)
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *AddTrackOperation) buildVideoTrackArgs(output string) []string {
	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	return []string{
		"-i", "INPUT_PLACEHOLDER",
		"-i", op.SourcePath,
		"-filter_complex", fmt.Sprintf("[1:v]setpts=PTS-STARTPTS+%f/TB[v1];[0:v][v1]overlay=enable='between(t,%f,%f)'",
			startSec, startSec, startSec+durationSec),
		"-c:a", "copy",
		"-y", output,
	}
}

func (op *AddTrackOperation) buildAudioTrackArgs(output string) []string {
	startSec := op.StartTime.Seconds()

	return []string{
		"-i", "INPUT_PLACEHOLDER",
		"-i", op.SourcePath,
		"-filter_complex", fmt.Sprintf("[1:a]adelay=%f|%f,volume=%f[a1];[0:a][a1]amix=inputs=2:duration=first[a]",
			startSec*1000, startSec*1000, op.Volume),
		"-map", "0:v",
		"-map", "[a]",
		"-y", output,
	}
}

func (op *AddTrackOperation) buildImageTrackArgs(output string) []string {
	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	return []string{
		"-i", "INPUT_PLACEHOLDER",
		"-i", op.SourcePath,
		"-filter_complex", fmt.Sprintf("[1:v]scale=iw:ih[img];[0:v][img]overlay=%d:%d:enable='between(t,%f,%f)'",
			op.X, op.Y, startSec, startSec+durationSec),
		"-c:a", "copy",
		"-y", output,
	}
}

func (op *AddTrackOperation) buildOverlayTrackArgs(output string) []string {
	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	return []string{
		"-i", "INPUT_PLACEHOLDER",
		"-i", op.SourcePath,
		"-filter_complex", fmt.Sprintf("[1:v]format=rgba,colorchannelmixer=aa=%f[ovl];[0:v][ovl]overlay=%d:%d:enable='between(t,%f,%f)'",
			op.Opacity, op.X, op.Y, startSec, startSec+durationSec),
		"-c:a", "copy",
		"-y", output,
	}
}

func (op *AddTrackOperation) buildTextTrackArgs(output string) []string {
	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	drawTextFilter := fmt.Sprintf("drawtext=text='%s':x=%d:y=%d:fontsize=%d:fontcolor=%s:enable='between(t,%f,%f)'",
		op.Text, op.X, op.Y, op.FontSize, op.Color, startSec, startSec+durationSec)

	return []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", drawTextFilter,
		"-c:a", "copy",
		"-y", output,
	}
}

func (op *AddTrackOperation) GetDescription() string {
	return fmt.Sprintf("添加%s轨道: %s (%.1fs-%.1fs)",
		op.TrackType, op.getSourceDescription(),
		op.StartTime.Seconds(), (op.StartTime + op.Duration).Seconds())
}

func (op *AddTrackOperation) getSourceDescription() string {
	switch op.TrackType {
	case "text":
		return op.Text
	default:
		return filepath.Base(op.SourcePath)
	}
}

func (op *AddTrackOperation) EstimateDuration() time.Duration {
	return 20 * time.Second
}

// TransitionOperation 转场操作
type TransitionOperation struct {
	Type     string        // 转场类型: fade, dissolve, wipe, slide
	Duration time.Duration // 转场时长
}

func (op *TransitionOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("transition", ".mp4")

	var transitionFilter string
	durationSec := op.Duration.Seconds()

	switch op.Type {
	case "fade":
		transitionFilter = fmt.Sprintf("fade=in:0:%d,fade=out:st=%f:d=%f",
			int(durationSec*30), durationSec, durationSec)
	case "dissolve":
		transitionFilter = fmt.Sprintf("fade=in:0:%d:alpha=1", int(durationSec*30))
	default:
		transitionFilter = fmt.Sprintf("fade=in:0:%d", int(durationSec*30))
	}

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", transitionFilter,
		"-c:a", "copy",
		"-y", tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *TransitionOperation) GetDescription() string {
	return fmt.Sprintf("转场效果: %s (%.1fs)", op.Type, op.Duration.Seconds())
}

func (op *TransitionOperation) EstimateDuration() time.Duration {
	return 15 * time.Second
}

// ComposeOperation 合成操作
type ComposeOperation struct{}

func (op *ComposeOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("compose", ".mp4")

	// 这里实现复杂的多轨道合成逻辑
	// 简化实现，实际应该根据时间轴信息进行合成
	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-c", "copy",
		"-y", tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *ComposeOperation) GetDescription() string {
	return "多媒体合成"
}

func (op *ComposeOperation) EstimateDuration() time.Duration {
	return 30 * time.Second
}

// PictureInPictureOperation 画中画操作
type PictureInPictureOperation struct {
	PipVideoPath string
	StartTime    time.Duration
	Duration     time.Duration
	X            int
	Y            int
	Width        int
	Height       int
}

func (op *PictureInPictureOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("pip", ".mp4")

	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	filterComplex := fmt.Sprintf("[1:v]scale=%d:%d[pip];[0:v][pip]overlay=%d:%d:enable='between(t,%f,%f)'",
		op.Width, op.Height, op.X, op.Y, startSec, startSec+durationSec)

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-i", op.PipVideoPath,
		"-filter_complex", filterComplex,
		"-c:a", "copy",
		"-y", tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *PictureInPictureOperation) GetDescription() string {
	return fmt.Sprintf("画中画: %s (%dx%d)", filepath.Base(op.PipVideoPath), op.Width, op.Height)
}

func (op *PictureInPictureOperation) EstimateDuration() time.Duration {
	return 25 * time.Second
}

// SplitScreenOperation 分屏操作
type SplitScreenOperation struct {
	Videos []string
	Layout string // "2x1", "1x2", "2x2"
}

func (op *SplitScreenOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("split_screen", ".mp4")

	var filterComplex string
	var inputs []string

	// 构建输入参数
	inputs = append(inputs, "-i", "INPUT_PLACEHOLDER")
	for _, video := range op.Videos {
		inputs = append(inputs, "-i", video)
	}

	// 根据布局构建滤镜
	switch op.Layout {
	case "2x1":
		filterComplex = "[0:v]scale=iw/2:ih[v0];[1:v]scale=iw/2:ih[v1];[v0][v1]hstack=inputs=2[v]"
	case "1x2":
		filterComplex = "[0:v]scale=iw:ih/2[v0];[1:v]scale=iw:ih/2[v1];[v0][v1]vstack=inputs=2[v]"
	case "2x2":
		if len(op.Videos) >= 3 {
			filterComplex = "[0:v]scale=iw/2:ih/2[v0];[1:v]scale=iw/2:ih/2[v1];[2:v]scale=iw/2:ih/2[v2];[3:v]scale=iw/2:ih/2[v3];[v0][v1]hstack[top];[v2][v3]hstack[bottom];[top][bottom]vstack[v]"
		}
	default:
		filterComplex = "[0:v]scale=iw/2:ih[v0];[1:v]scale=iw/2:ih[v1];[v0][v1]hstack=inputs=2[v]"
	}

	args := append(inputs, "-filter_complex", filterComplex, "-map", "[v]", "-c:a", "copy", "-y", tempOutput)

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *SplitScreenOperation) GetDescription() string {
	return fmt.Sprintf("分屏显示: %s (%d个视频)", op.Layout, len(op.Videos))
}

func (op *SplitScreenOperation) EstimateDuration() time.Duration {
	return 35 * time.Second
}

// BackgroundMusicOperation 背景音乐操作
type BackgroundMusicOperation struct {
	MusicPath string
	StartTime time.Duration
	Duration  time.Duration
	Volume    float64
	FadeIn    time.Duration
	FadeOut   time.Duration
}

func (op *BackgroundMusicOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("bg_music", ".mp4")

	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()
	fadeInSec := op.FadeIn.Seconds()
	fadeOutSec := op.FadeOut.Seconds()

	audioFilter := fmt.Sprintf("[1:a]adelay=%f|%f,volume=%f", startSec*1000, startSec*1000, op.Volume)

	if fadeInSec > 0 {
		audioFilter += fmt.Sprintf(",afade=in:st=%f:d=%f", startSec, fadeInSec)
	}
	if fadeOutSec > 0 {
		audioFilter += fmt.Sprintf(",afade=out:st=%f:d=%f", startSec+durationSec-fadeOutSec, fadeOutSec)
	}

	audioFilter += "[bg];[0:a][bg]amix=inputs=2:duration=first[a]"

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-i", op.MusicPath,
		"-filter_complex", audioFilter,
		"-map", "0:v",
		"-map", "[a]",
		"-y", tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *BackgroundMusicOperation) GetDescription() string {
	return fmt.Sprintf("背景音乐: %s (音量%.1f)", filepath.Base(op.MusicPath), op.Volume)
}

func (op *BackgroundMusicOperation) EstimateDuration() time.Duration {
	return 20 * time.Second
}

// SubtitleOperation 字幕操作
type SubtitleOperation struct {
	SubtitleFile string
	Style        *SubtitleStyle
}

func (op *SubtitleOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("subtitle", ".mp4")

	// 构建字幕滤镜
	subtitleFilter := fmt.Sprintf("subtitles=%s", op.SubtitleFile)

	if op.Style != nil {
		// 添加样式参数
		if op.Style.FontSize > 0 {
			subtitleFilter += fmt.Sprintf(":force_style='FontSize=%d", op.Style.FontSize)
		}
		if op.Style.FontColor != "" {
			subtitleFilter += fmt.Sprintf(",PrimaryColour=%s", op.Style.FontColor)
		}
		if op.Style.FontFamily != "" {
			subtitleFilter += fmt.Sprintf(",FontName=%s", op.Style.FontFamily)
		}
		subtitleFilter += "'"
	}

	args := []string{
		"-i", "INPUT_PLACEHOLDER",
		"-vf", subtitleFilter,
		"-c:a", "copy",
		"-y", tempOutput,
	}

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *SubtitleOperation) GetDescription() string {
	return fmt.Sprintf("添加字幕: %s", filepath.Base(op.SubtitleFile))
}

func (op *SubtitleOperation) EstimateDuration() time.Duration {
	return 25 * time.Second
}

// SlideshowOperation 幻灯片操作
type SlideshowOperation struct {
	Images     []string
	Duration   time.Duration // 每张图片的持续时间
	Transition string        // 转场效果
}

func (op *SlideshowOperation) Execute(ctx context.Context, ffmpeg *FFmpeg) error {
	tempOutput := generateTempFilePath("slideshow", ".mp4")

	// 构建输入参数
	var inputs []string
	for _, image := range op.Images {
		inputs = append(inputs, "-loop", "1", "-t", fmt.Sprintf("%.2f", op.Duration.Seconds()), "-i", image)
	}

	// 构建滤镜复合体
	var filterParts []string
	for i := range op.Images {
		filterParts = append(filterParts, fmt.Sprintf("[%d:v]", i))
	}

	filterComplex := strings.Join(filterParts, "") + fmt.Sprintf("concat=n=%d:v=1:a=0[v]", len(op.Images))

	args := append(inputs, "-filter_complex", filterComplex, "-map", "[v]", "-c:v", "libx264", "-pix_fmt", "yuv420p", "-y", tempOutput)

	_, err := ffmpeg.executeCommand(ctx, args)
	return err
}

func (op *SlideshowOperation) GetDescription() string {
	return fmt.Sprintf("幻灯片: %d张图片 (每张%.1fs)", len(op.Images), op.Duration.Seconds())
}

func (op *SlideshowOperation) EstimateDuration() time.Duration {
	return time.Duration(len(op.Images)) * 5 * time.Second // 每张图片估算5秒处理时间
}

// generateTempFilePath 生成临时文件路径
func generateTempFilePath(prefix, ext string) string {
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%s_%d%s", prefix, timestamp, ext)
	return filepath.Join(os.TempDir(), filename)
}
