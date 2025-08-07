// Package ffmpeg 提供链式调用的视频编辑器
package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// VideoEditor 链式调用的视频编辑器
type VideoEditor struct {
	ffmpeg     *FFmpeg              // FFmpeg实例
	operations []VideoEditOperation // 操作队列
	timeline   *Timeline            // 时间轴
	inputPath  string               // 输入文件路径
	outputPath string               // 输出文件路径
	mu         sync.RWMutex         // 读写锁
	cancelled  bool                 // 是否已取消
	progress   ProgressCallback     // 进度回调
}

// NewVideoEditor 创建新的视频编辑器实例
func NewVideoEditor(ffmpeg *FFmpeg, inputPath string) *VideoEditor {
	return &VideoEditor{
		ffmpeg:     ffmpeg,
		operations: make([]VideoEditOperation, 0),
		timeline: &Timeline{
			Tracks: make(map[string][]*TimelineElement),
		},
		inputPath: inputPath,
	}
}

// SetOutput 设置输出文件路径
func (ve *VideoEditor) SetOutput(outputPath string) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()
	ve.outputPath = outputPath
	return ve
}

// SetProgressCallback 设置进度回调函数
func (ve *VideoEditor) SetProgressCallback(callback ProgressCallback) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()
	ve.progress = callback
	return ve
}

// AddOperation 添加操作到队列
func (ve *VideoEditor) AddOperation(op VideoEditOperation) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()
	ve.operations = append(ve.operations, op)
	return ve
}

// CropTime 时间段裁剪
func (ve *VideoEditor) CropTime(startTime, duration string) *VideoEditor {
	op := &CropTimeOperation{
		StartTime: startTime,
		Duration:  duration,
	}
	return ve.AddOperation(op)
}

// CropDimension 尺寸裁剪
func (ve *VideoEditor) CropDimension(dimensions *CropDimensions) *VideoEditor {
	op := &CropDimensionOperation{
		Dimensions: dimensions,
	}
	return ve.AddOperation(op)
}

// Resize 调整视频尺寸
func (ve *VideoEditor) Resize(width, height int) *VideoEditor {
	op := &ResizeOperation{
		Width:  width,
		Height: height,
	}
	return ve.AddOperation(op)
}

// AddWatermark 添加水印
func (ve *VideoEditor) AddWatermark(options *WatermarkOptions) *VideoEditor {
	op := &WatermarkOperation{
		Options: options,
	}
	return ve.AddOperation(op)
}

// SeparateAudio 分离音频
func (ve *VideoEditor) SeparateAudio(audioOutputPath string) *VideoEditor {
	op := &SeparateAudioOperation{
		AudioOutputPath: audioOutputPath,
	}
	return ve.AddOperation(op)
}

// MixAudio 混合音频
func (ve *VideoEditor) MixAudio(options *AudioMixOptions) *VideoEditor {
	op := &AudioMixOperation{
		Options: options,
	}
	return ve.AddOperation(op)
}

// InsertImage 插入图片
func (ve *VideoEditor) InsertImage(imagePath string, startTime, duration time.Duration) *VideoEditor {
	op := &InsertImageOperation{
		ImagePath: imagePath,
		StartTime: startTime,
		Duration:  duration,
	}
	return ve.AddOperation(op)
}

// EditFrame 编辑帧
func (ve *VideoEditor) EditFrame(options *FrameEditOptions) *VideoEditor {
	op := &FrameEditOperation{
		Options: options,
	}
	return ve.AddOperation(op)
}

// CropTimeRange 裁剪时间范围（便捷方法）
func (ve *VideoEditor) CropTimeRange(startTime, endTime string) *VideoEditor {
	// 计算持续时间
	start, _ := parseTime(startTime)
	end, _ := parseTime(endTime)
	duration := fmt.Sprintf("%.2f", end-start)

	return ve.CropTime(startTime, duration)
}

// AddText 添加文字（使用drawtext滤镜）
func (ve *VideoEditor) AddText(text string, x, y int, fontSize int, color string, startTime, duration time.Duration) *VideoEditor {
	op := &TextOperation{
		Text:      text,
		X:         x,
		Y:         y,
		FontSize:  fontSize,
		Color:     color,
		StartTime: startTime,
		Duration:  duration,
	}
	return ve.AddOperation(op)
}

// FadeIn 添加淡入效果
func (ve *VideoEditor) FadeIn(duration time.Duration) *VideoEditor {
	op := &FadeOperation{
		Type:     "in",
		Duration: duration,
	}
	return ve.AddOperation(op)
}

// FadeOut 添加淡出效果
func (ve *VideoEditor) FadeOut(duration time.Duration) *VideoEditor {
	op := &FadeOperation{
		Type:     "out",
		Duration: duration,
	}
	return ve.AddOperation(op)
}

// ChangeSpeed 改变播放速度
func (ve *VideoEditor) ChangeSpeed(factor float64) *VideoEditor {
	op := &SpeedOperation{
		Factor: factor,
	}
	return ve.AddOperation(op)
}

// Rotate 旋转视频
func (ve *VideoEditor) Rotate(angle int) *VideoEditor {
	op := &RotateOperation{
		Angle: angle,
	}
	return ve.AddOperation(op)
}

// Mirror 镜像翻转
func (ve *VideoEditor) Mirror(horizontal bool) *VideoEditor {
	op := &MirrorOperation{
		Horizontal: horizontal,
	}
	return ve.AddOperation(op)
}

// AdjustBrightness 调整亮度
func (ve *VideoEditor) AdjustBrightness(brightness float64) *VideoEditor {
	op := &BrightnessOperation{
		Brightness: brightness,
	}
	return ve.AddOperation(op)
}

// AdjustContrast 调整对比度
func (ve *VideoEditor) AdjustContrast(contrast float64) *VideoEditor {
	op := &ContrastOperation{
		Contrast: contrast,
	}
	return ve.AddOperation(op)
}

// AddBlur 添加模糊效果
func (ve *VideoEditor) AddBlur(radius float64) *VideoEditor {
	op := &BlurOperation{
		Radius: radius,
	}
	return ve.AddOperation(op)
}

// Stabilize 视频防抖
func (ve *VideoEditor) Stabilize() *VideoEditor {
	op := &StabilizeOperation{}
	return ve.AddOperation(op)
}

// ExtractFrames 提取帧为图片序列
func (ve *VideoEditor) ExtractFrames(outputDir string, fps float64) *VideoEditor {
	op := &ExtractFramesOperation{
		OutputDir: outputDir,
		FPS:       fps,
	}
	return ve.AddOperation(op)
}

// CreateFromImages 从图片序列创建视频
func (ve *VideoEditor) CreateFromImages(imagePattern string, fps float64) *VideoEditor {
	op := &CreateFromImagesOperation{
		ImagePattern: imagePattern,
		FPS:          fps,
	}
	return ve.AddOperation(op)
}

// === 多媒体合成链式调用方法 ===

// AddVideoTrack 添加视频轨道
func (ve *VideoEditor) AddVideoTrack(videoPath string, startTime, duration time.Duration) *VideoEditor {
	op := &AddTrackOperation{
		TrackType:  "video",
		SourcePath: videoPath,
		StartTime:  startTime,
		Duration:   duration,
	}
	return ve.AddOperation(op)
}

// AddAudioTrack 添加音频轨道
func (ve *VideoEditor) AddAudioTrack(audioPath string, startTime, duration time.Duration, volume float64) *VideoEditor {
	op := &AddTrackOperation{
		TrackType:  "audio",
		SourcePath: audioPath,
		StartTime:  startTime,
		Duration:   duration,
		Volume:     volume,
	}
	return ve.AddOperation(op)
}

// AddImageTrack 添加图片轨道（在指定时间显示图片）
func (ve *VideoEditor) AddImageTrack(imagePath string, startTime, duration time.Duration, x, y int) *VideoEditor {
	op := &AddTrackOperation{
		TrackType:  "image",
		SourcePath: imagePath,
		StartTime:  startTime,
		Duration:   duration,
		X:          x,
		Y:          y,
	}
	return ve.AddOperation(op)
}

// AddOverlayTrack 添加叠加轨道（水印、贴图等）
func (ve *VideoEditor) AddOverlayTrack(overlayPath string, startTime, duration time.Duration, x, y int, opacity float64) *VideoEditor {
	op := &AddTrackOperation{
		TrackType:  "overlay",
		SourcePath: overlayPath,
		StartTime:  startTime,
		Duration:   duration,
		X:          x,
		Y:          y,
		Opacity:    opacity,
	}
	return ve.AddOperation(op)
}

// AddTextTrack 添加文字轨道
func (ve *VideoEditor) AddTextTrack(text string, startTime, duration time.Duration, x, y int, fontSize int, color string) *VideoEditor {
	op := &AddTrackOperation{
		TrackType: "text",
		Text:      text,
		StartTime: startTime,
		Duration:  duration,
		X:         x,
		Y:         y,
		FontSize:  fontSize,
		Color:     color,
	}
	return ve.AddOperation(op)
}

// SetTransition 设置轨道间的转场效果
func (ve *VideoEditor) SetTransition(transitionType string, duration time.Duration) *VideoEditor {
	op := &TransitionOperation{
		Type:     transitionType,
		Duration: duration,
	}
	return ve.AddOperation(op)
}

// ComposeMultimedia 执行多媒体合成（将所有轨道合成为最终视频）
func (ve *VideoEditor) ComposeMultimedia() *VideoEditor {
	op := &ComposeOperation{}
	return ve.AddOperation(op)
}

// === 高级合成方法 ===

// PictureInPicture 画中画效果
func (ve *VideoEditor) PictureInPicture(pipVideoPath string, startTime, duration time.Duration, x, y, width, height int) *VideoEditor {
	op := &PictureInPictureOperation{
		PipVideoPath: pipVideoPath,
		StartTime:    startTime,
		Duration:     duration,
		X:            x,
		Y:            y,
		Width:        width,
		Height:       height,
	}
	return ve.AddOperation(op)
}

// SplitScreen 分屏效果
func (ve *VideoEditor) SplitScreen(videos []string, layout string) *VideoEditor {
	op := &SplitScreenOperation{
		Videos: videos,
		Layout: layout, // "2x1", "1x2", "2x2" 等
	}
	return ve.AddOperation(op)
}

// AddBackgroundMusic 添加背景音乐（支持多段音乐）
func (ve *VideoEditor) AddBackgroundMusic(musicPath string, startTime, duration time.Duration, volume float64, fadeIn, fadeOut time.Duration) *VideoEditor {
	op := &BackgroundMusicOperation{
		MusicPath: musicPath,
		StartTime: startTime,
		Duration:  duration,
		Volume:    volume,
		FadeIn:    fadeIn,
		FadeOut:   fadeOut,
	}
	return ve.AddOperation(op)
}

// AddSubtitles 添加字幕
func (ve *VideoEditor) AddSubtitles(subtitleFile string, style *SubtitleStyle) *VideoEditor {
	op := &SubtitleOperation{
		SubtitleFile: subtitleFile,
		Style:        style,
	}
	return ve.AddOperation(op)
}

// CreateSlideshow 创建幻灯片（从多张图片）
func (ve *VideoEditor) CreateSlideshow(images []string, duration time.Duration, transition string) *VideoEditor {
	op := &SlideshowOperation{
		Images:     images,
		Duration:   duration,
		Transition: transition,
	}
	return ve.AddOperation(op)
}

// Execute 执行所有操作
func (ve *VideoEditor) Execute() error {
	return ve.ExecuteWithContext(context.Background())
}

// ExecuteWithContext 带上下文执行所有操作
func (ve *VideoEditor) ExecuteWithContext(ctx context.Context) error {
	ve.mu.RLock()
	operations := make([]VideoEditOperation, len(ve.operations))
	copy(operations, ve.operations)
	outputPath := ve.outputPath
	progress := ve.progress
	inputPath := ve.inputPath
	ve.mu.RUnlock()

	if outputPath == "" {
		return NewError(ErrInvalidOptions, "输出路径未设置", nil)
	}

	if len(operations) == 0 {
		return NewError(ErrInvalidOptions, "没有要执行的操作", nil)
	}

	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	// 创建进度跟踪器
	tracker := NewProgressTracker(progress)

	// 添加步骤到跟踪器
	for i, op := range operations {
		stepName := fmt.Sprintf("步骤%d: %s", i+1, op.GetDescription())
		weight := 1.0 // 每个步骤权重相等
		duration := op.EstimateDuration()
		tracker.AddStep(stepName, weight, duration)
	}

	// 执行操作链
	return ve.executeOperationChain(ctx, operations, tracker, inputPath, outputPath)
}

// executeOperationChain 执行操作链
func (ve *VideoEditor) executeOperationChain(ctx context.Context, operations []VideoEditOperation,
	tracker *ProgressTracker, inputPath, outputPath string) error {

	currentInput := inputPath
	tempFiles := make([]string, 0) // 记录临时文件用于清理

	defer func() {
		// 清理临时文件
		ve.cleanupTempFiles(tempFiles)
	}()

	for i, op := range operations {
		// 检查是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		ve.mu.RLock()
		cancelled := ve.cancelled
		ve.mu.RUnlock()

		if cancelled {
			return NewError(ErrExecutionFailed, "操作已取消", nil)
		}

		// 获取步骤进度回调
		stepCallback := tracker.StartStep(i)

		// 确定输出路径
		var currentOutput string
		if i == len(operations)-1 {
			// 最后一个操作，输出到最终路径
			currentOutput = outputPath
		} else {
			// 中间操作，输出到临时文件
			currentOutput = generateTempFilePath(fmt.Sprintf("step_%d", i), ".mp4")
			tempFiles = append(tempFiles, currentOutput)
		}

		// 执行操作
		if err := ve.executeOperation(ctx, op, currentInput, currentOutput, stepCallback); err != nil {
			return fmt.Errorf("执行操作 %d (%s) 失败: %w", i+1, op.GetDescription(), err)
		}

		// 更新输入路径为当前输出路径
		currentInput = currentOutput

		ve.ffmpeg.logger.Info("操作完成: %s (%d/%d)", op.GetDescription(), i+1, len(operations))
	}

	ve.ffmpeg.logger.Info("所有视频编辑操作完成: %s", outputPath)
	return nil
}

// executeOperation 执行单个操作
func (ve *VideoEditor) executeOperation(ctx context.Context, op VideoEditOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	// 这里需要根据具体操作类型来构建和执行FFmpeg命令
	// 暂时使用简化的实现
	switch operation := op.(type) {
	case *CropTimeOperation:
		return ve.executeCropTimeOperation(ctx, operation, inputPath, outputPath, callback)
	case *CropDimensionOperation:
		return ve.executeCropDimensionOperation(ctx, operation, inputPath, outputPath, callback)
	case *ResizeOperation:
		return ve.executeResizeOperation(ctx, operation, inputPath, outputPath, callback)
	case *WatermarkOperation:
		return ve.executeWatermarkOperation(ctx, operation, inputPath, outputPath, callback)
	case *SeparateAudioOperation:
		return ve.executeSeparateAudioOperation(ctx, operation, inputPath, outputPath, callback)
	case *AudioMixOperation:
		return ve.executeAudioMixOperation(ctx, operation, inputPath, outputPath, callback)
	case *InsertImageOperation:
		return ve.executeInsertImageOperation(ctx, operation, inputPath, outputPath, callback)
	case *FrameEditOperation:
		return ve.executeFrameEditOperation(ctx, operation, inputPath, outputPath, callback)
	case *TextOperation:
		return ve.executeTextOperation(ctx, operation, inputPath, outputPath, callback)
	case *FadeOperation:
		return ve.executeFadeOperation(ctx, operation, inputPath, outputPath, callback)
	case *SpeedOperation:
		return ve.executeSpeedOperation(ctx, operation, inputPath, outputPath, callback)
	case *RotateOperation:
		return ve.executeRotateOperation(ctx, operation, inputPath, outputPath, callback)
	case *MirrorOperation:
		return ve.executeMirrorOperation(ctx, operation, inputPath, outputPath, callback)
	case *BrightnessOperation:
		return ve.executeBrightnessOperation(ctx, operation, inputPath, outputPath, callback)
	case *ContrastOperation:
		return ve.executeContrastOperation(ctx, operation, inputPath, outputPath, callback)
	case *BlurOperation:
		return ve.executeBlurOperation(ctx, operation, inputPath, outputPath, callback)
	case *StabilizeOperation:
		return ve.executeStabilizeOperation(ctx, operation, inputPath, outputPath, callback)
	case *ExtractFramesOperation:
		return ve.executeExtractFramesOperation(ctx, operation, inputPath, outputPath, callback)
	case *CreateFromImagesOperation:
		return ve.executeCreateFromImagesOperation(ctx, operation, inputPath, outputPath, callback)
	case *AddTrackOperation:
		return ve.executeAddTrackOperation(ctx, operation, inputPath, outputPath, callback)
	case *TransitionOperation:
		return ve.executeTransitionOperation(ctx, operation, inputPath, outputPath, callback)
	case *ComposeOperation:
		return ve.executeComposeOperation(ctx, operation, inputPath, outputPath, callback)
	case *PictureInPictureOperation:
		return ve.executePictureInPictureOperation(ctx, operation, inputPath, outputPath, callback)
	case *SplitScreenOperation:
		return ve.executeSplitScreenOperation(ctx, operation, inputPath, outputPath, callback)
	case *BackgroundMusicOperation:
		return ve.executeBackgroundMusicOperation(ctx, operation, inputPath, outputPath, callback)
	case *SubtitleOperation:
		return ve.executeSubtitleOperation(ctx, operation, inputPath, outputPath, callback)
	case *SlideshowOperation:
		return ve.executeSlideshowOperation(ctx, operation, inputPath, outputPath, callback)
	default:
		return fmt.Errorf("不支持的操作类型: %T", op)
	}
}

// cleanupTempFiles 清理临时文件
func (ve *VideoEditor) cleanupTempFiles(tempFiles []string) {
	for _, file := range tempFiles {
		if err := os.Remove(file); err != nil {
			ve.ffmpeg.logger.Error("清理临时文件失败: %s, 错误: %v", file, err)
		}
	}
}

// Cancel 取消操作
func (ve *VideoEditor) Cancel() {
	ve.mu.Lock()
	defer ve.mu.Unlock()
	ve.cancelled = true
}

// GetOperationCount 获取操作数量
func (ve *VideoEditor) GetOperationCount() int {
	ve.mu.RLock()
	defer ve.mu.RUnlock()
	return len(ve.operations)
}

// GetTimeline 获取时间轴
func (ve *VideoEditor) GetTimeline() *Timeline {
	ve.mu.RLock()
	defer ve.mu.RUnlock()
	return ve.timeline
}

// Clear 清空所有操作
func (ve *VideoEditor) Clear() *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()
	ve.operations = make([]VideoEditOperation, 0)
	ve.cancelled = false
	return ve
}

// Clone 克隆编辑器（不包括操作队列）
func (ve *VideoEditor) Clone() *VideoEditor {
	ve.mu.RLock()
	defer ve.mu.RUnlock()

	return &VideoEditor{
		ffmpeg:     ve.ffmpeg,
		timeline:   ve.timeline,
		inputPath:  ve.inputPath,
		operations: make([]VideoEditOperation, 0),
	}
}

// executeCropTimeOperation 执行时间段裁剪操作
func (ve *VideoEditor) executeCropTimeOperation(ctx context.Context, op *CropTimeOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	args := []string{
		"-i", inputPath,
		"-ss", op.StartTime,
		"-t", op.Duration,
		"-c", "copy",
		"-y",
		outputPath,
	}

	// 获取视频信息以估算总时长
	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeCropDimensionOperation 执行尺寸裁剪操作
func (ve *VideoEditor) executeCropDimensionOperation(ctx context.Context, op *CropDimensionOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	cropFilter := fmt.Sprintf("crop=%d:%d:%d:%d",
		op.Dimensions.Width, op.Dimensions.Height,
		op.Dimensions.X, op.Dimensions.Y)

	args := []string{
		"-i", inputPath,
		"-vf", cropFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeResizeOperation 执行调整尺寸操作
func (ve *VideoEditor) executeResizeOperation(ctx context.Context, op *ResizeOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	scaleFilter := fmt.Sprintf("scale=%d:%d", op.Width, op.Height)

	args := []string{
		"-i", inputPath,
		"-vf", scaleFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeWatermarkOperation 执行水印操作
func (ve *VideoEditor) executeWatermarkOperation(ctx context.Context, op *WatermarkOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	overlayFilter := fmt.Sprintf("overlay=%d:%d", op.Options.X, op.Options.Y)
	if op.Options.Opacity < 1.0 {
		overlayFilter = fmt.Sprintf("format=rgba,colorchannelmixer=aa=%f[wm];[0:v][wm]overlay=%d:%d",
			op.Options.Opacity, op.Options.X, op.Options.Y)
	}

	args := []string{
		"-i", inputPath,
		"-i", op.Options.ImagePath,
		"-filter_complex", overlayFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeSeparateAudioOperation 执行分离音频操作
func (ve *VideoEditor) executeSeparateAudioOperation(ctx context.Context, op *SeparateAudioOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	// 提取音频
	audioArgs := []string{
		"-i", inputPath,
		"-vn", // 不包含视频
		"-acodec", "copy",
		"-y",
		op.AudioOutputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	// 先提取音频
	if err := ve.ffmpeg.executeCommandWithProgress(ctx, audioArgs, nil, videoInfo.Duration/2); err != nil {
		return fmt.Errorf("提取音频失败: %w", err)
	}

	// 创建无音频视频
	videoArgs := []string{
		"-i", inputPath,
		"-an", // 不包含音频
		"-vcodec", "copy",
		"-y",
		outputPath,
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, videoArgs, callback, videoInfo.Duration/2)
}

// executeAudioMixOperation 执行音频混合操作
func (ve *VideoEditor) executeAudioMixOperation(ctx context.Context, op *AudioMixOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

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
		"-i", inputPath,
		"-i", op.Options.BackgroundPath,
		"-filter_complex", filterComplex,
		"-map", "0:v",
		"-map", "[a]",
		"-c:v", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeInsertImageOperation 执行插入图片操作
func (ve *VideoEditor) executeInsertImageOperation(ctx context.Context, op *InsertImageOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	filterComplex := fmt.Sprintf("[1:v]scale=iw:ih[img];[0:v][img]overlay=enable='between(t,%f,%f)'",
		startSec, startSec+durationSec)

	args := []string{
		"-i", inputPath,
		"-i", op.ImagePath,
		"-filter_complex", filterComplex,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeFrameEditOperation 执行帧编辑操作
func (ve *VideoEditor) executeFrameEditOperation(ctx context.Context, op *FrameEditOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	switch op.Options.Operation {
	case FrameInsert:
		return ve.executeInsertFrame(ctx, op.Options, inputPath, outputPath, callback, videoInfo.Duration)
	case FrameDelete:
		return ve.executeDeleteFrame(ctx, op.Options, inputPath, outputPath, callback, videoInfo.Duration)
	case FrameReplace:
		return ve.executeReplaceFrame(ctx, op.Options, inputPath, outputPath, callback, videoInfo.Duration)
	default:
		return fmt.Errorf("不支持的帧操作: %s", op.Options.Operation)
	}
}

// executeInsertFrame 执行插入帧操作
func (ve *VideoEditor) executeInsertFrame(ctx context.Context, options *FrameEditOptions,
	inputPath, outputPath string, callback ProgressCallback, totalDuration time.Duration) error {

	args := []string{
		"-i", inputPath,
		"-i", options.ImagePath,
		"-filter_complex", "[0:v][1:v]concat=n=2:v=1:a=0[v]",
		"-map", "[v]",
		"-map", "0:a",
		"-y",
		outputPath,
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, totalDuration)
}

// executeDeleteFrame 执行删除帧操作
func (ve *VideoEditor) executeDeleteFrame(ctx context.Context, options *FrameEditOptions,
	inputPath, outputPath string, callback ProgressCallback, totalDuration time.Duration) error {

	frameTime := float64(options.FrameNumber) / 30.0 // 假设30fps

	args := []string{
		"-i", inputPath,
		"-vf", fmt.Sprintf("select='not(between(t,%f,%f))',setpts=N/FRAME_RATE/TB",
			frameTime, frameTime+0.033),
		"-af", "aselect='not(between(t,%f,%f))',asetpts=N/SR/TB",
		"-y",
		outputPath,
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, totalDuration)
}

// executeReplaceFrame 执行替换帧操作
func (ve *VideoEditor) executeReplaceFrame(ctx context.Context, options *FrameEditOptions,
	inputPath, outputPath string, callback ProgressCallback, totalDuration time.Duration) error {

	frameTime := float64(options.FrameNumber) / 30.0

	filterComplex := fmt.Sprintf("[1:v]scale=iw:ih[img];[0:v][img]overlay=enable='between(t,%f,%f)'",
		frameTime, frameTime+0.033)

	args := []string{
		"-i", inputPath,
		"-i", options.ImagePath,
		"-filter_complex", filterComplex,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, totalDuration)
}

// executeTextOperation 执行文字添加操作
func (ve *VideoEditor) executeTextOperation(ctx context.Context, op *TextOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	drawTextFilter := fmt.Sprintf("drawtext=text='%s':x=%d:y=%d:fontsize=%d:fontcolor=%s:enable='between(t,%f,%f)'",
		op.Text, op.X, op.Y, op.FontSize, op.Color, startSec, startSec+durationSec)

	args := []string{
		"-i", inputPath,
		"-vf", drawTextFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeFadeOperation 执行淡入淡出操作
func (ve *VideoEditor) executeFadeOperation(ctx context.Context, op *FadeOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	var fadeFilter string
	if op.Type == "in" {
		fadeFilter = fmt.Sprintf("fade=in:0:%d", int(op.Duration.Seconds()*30)) // 假设30fps
	} else {
		fadeFilter = fmt.Sprintf("fade=out:st=0:d=%f", op.Duration.Seconds())
	}

	args := []string{
		"-i", inputPath,
		"-vf", fadeFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeSpeedOperation 执行速度调整操作
func (ve *VideoEditor) executeSpeedOperation(ctx context.Context, op *SpeedOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	videoFilter := fmt.Sprintf("setpts=%f*PTS", 1.0/op.Factor)
	audioFilter := fmt.Sprintf("atempo=%f", op.Factor)

	args := []string{
		"-i", inputPath,
		"-vf", videoFilter,
		"-af", audioFilter,
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeRotateOperation 执行旋转操作
func (ve *VideoEditor) executeRotateOperation(ctx context.Context, op *RotateOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

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
		"-i", inputPath,
		"-vf", rotateFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeMirrorOperation 执行镜像翻转操作
func (ve *VideoEditor) executeMirrorOperation(ctx context.Context, op *MirrorOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	var flipFilter string
	if op.Horizontal {
		flipFilter = "hflip"
	} else {
		flipFilter = "vflip"
	}

	args := []string{
		"-i", inputPath,
		"-vf", flipFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeBrightnessOperation 执行亮度调整操作
func (ve *VideoEditor) executeBrightnessOperation(ctx context.Context, op *BrightnessOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	brightnessFilter := fmt.Sprintf("eq=brightness=%f", op.Brightness)

	args := []string{
		"-i", inputPath,
		"-vf", brightnessFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeContrastOperation 执行对比度调整操作
func (ve *VideoEditor) executeContrastOperation(ctx context.Context, op *ContrastOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	contrastFilter := fmt.Sprintf("eq=contrast=%f", op.Contrast)

	args := []string{
		"-i", inputPath,
		"-vf", contrastFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeBlurOperation 执行模糊效果操作
func (ve *VideoEditor) executeBlurOperation(ctx context.Context, op *BlurOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	blurFilter := fmt.Sprintf("boxblur=%f", op.Radius)

	args := []string{
		"-i", inputPath,
		"-vf", blurFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeStabilizeOperation 执行视频防抖操作
func (ve *VideoEditor) executeStabilizeOperation(ctx context.Context, op *StabilizeOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	stabilizeFilter := "vidstabtransform=smoothing=30"

	args := []string{
		"-i", inputPath,
		"-vf", stabilizeFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeExtractFramesOperation 执行提取帧操作
func (ve *VideoEditor) executeExtractFramesOperation(ctx context.Context, op *ExtractFramesOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	outputPattern := filepath.Join(op.OutputDir, "frame_%04d.png")

	args := []string{
		"-i", inputPath,
		"-vf", fmt.Sprintf("fps=%f", op.FPS),
		"-y",
		outputPattern,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeCreateFromImagesOperation 执行从图片创建视频操作
func (ve *VideoEditor) executeCreateFromImagesOperation(ctx context.Context, op *CreateFromImagesOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	args := []string{
		"-framerate", fmt.Sprintf("%f", op.FPS),
		"-i", op.ImagePattern,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-y",
		outputPath,
	}

	// 对于从图片创建视频，无法准确估算时长，使用固定值
	estimatedDuration := 30 * time.Second

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, estimatedDuration)
}

// === 多媒体合成操作执行方法 ===

// executeAddTrackOperation 执行添加轨道操作
func (ve *VideoEditor) executeAddTrackOperation(ctx context.Context, op *AddTrackOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	var args []string

	switch op.TrackType {
	case "video":
		args = ve.buildVideoTrackArgs(op, inputPath, outputPath)
	case "audio":
		args = ve.buildAudioTrackArgs(op, inputPath, outputPath)
	case "image":
		args = ve.buildImageTrackArgs(op, inputPath, outputPath)
	case "overlay":
		args = ve.buildOverlayTrackArgs(op, inputPath, outputPath)
	case "text":
		args = ve.buildTextTrackArgs(op, inputPath, outputPath)
	default:
		return fmt.Errorf("不支持的轨道类型: %s", op.TrackType)
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

func (ve *VideoEditor) buildVideoTrackArgs(op *AddTrackOperation, inputPath, outputPath string) []string {
	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	return []string{
		"-i", inputPath,
		"-i", op.SourcePath,
		"-filter_complex", fmt.Sprintf("[1:v]setpts=PTS-STARTPTS+%f/TB[v1];[0:v][v1]overlay=enable='between(t,%f,%f)'",
			startSec, startSec, startSec+durationSec),
		"-c:a", "copy",
		"-y", outputPath,
	}
}

func (ve *VideoEditor) buildAudioTrackArgs(op *AddTrackOperation, inputPath, outputPath string) []string {
	startSec := op.StartTime.Seconds()

	return []string{
		"-i", inputPath,
		"-i", op.SourcePath,
		"-filter_complex", fmt.Sprintf("[1:a]adelay=%f|%f,volume=%f[a1];[0:a][a1]amix=inputs=2:duration=first[a]",
			startSec*1000, startSec*1000, op.Volume),
		"-map", "0:v",
		"-map", "[a]",
		"-y", outputPath,
	}
}

func (ve *VideoEditor) buildImageTrackArgs(op *AddTrackOperation, inputPath, outputPath string) []string {
	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	return []string{
		"-i", inputPath,
		"-i", op.SourcePath,
		"-filter_complex", fmt.Sprintf("[1:v]scale=iw:ih[img];[0:v][img]overlay=%d:%d:enable='between(t,%f,%f)'",
			op.X, op.Y, startSec, startSec+durationSec),
		"-c:a", "copy",
		"-y", outputPath,
	}
}

func (ve *VideoEditor) buildOverlayTrackArgs(op *AddTrackOperation, inputPath, outputPath string) []string {
	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	return []string{
		"-i", inputPath,
		"-i", op.SourcePath,
		"-filter_complex", fmt.Sprintf("[1:v]format=rgba,colorchannelmixer=aa=%f[ovl];[0:v][ovl]overlay=%d:%d:enable='between(t,%f,%f)'",
			op.Opacity, op.X, op.Y, startSec, startSec+durationSec),
		"-c:a", "copy",
		"-y", outputPath,
	}
}

func (ve *VideoEditor) buildTextTrackArgs(op *AddTrackOperation, inputPath, outputPath string) []string {
	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	drawTextFilter := fmt.Sprintf("drawtext=text='%s':x=%d:y=%d:fontsize=%d:fontcolor=%s:enable='between(t,%f,%f)'",
		op.Text, op.X, op.Y, op.FontSize, op.Color, startSec, startSec+durationSec)

	return []string{
		"-i", inputPath,
		"-vf", drawTextFilter,
		"-c:a", "copy",
		"-y", outputPath,
	}
}

// executeTransitionOperation 执行转场操作
func (ve *VideoEditor) executeTransitionOperation(ctx context.Context, op *TransitionOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

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
		"-i", inputPath,
		"-vf", transitionFilter,
		"-c:a", "copy",
		"-y", outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeComposeOperation 执行合成操作
func (ve *VideoEditor) executeComposeOperation(ctx context.Context, op *ComposeOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	// 简化的合成实现，实际应该根据时间轴信息进行复杂合成
	args := []string{
		"-i", inputPath,
		"-c", "copy",
		"-y", outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executePictureInPictureOperation 执行画中画操作
func (ve *VideoEditor) executePictureInPictureOperation(ctx context.Context, op *PictureInPictureOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	startSec := op.StartTime.Seconds()
	durationSec := op.Duration.Seconds()

	filterComplex := fmt.Sprintf("[1:v]scale=%d:%d[pip];[0:v][pip]overlay=%d:%d:enable='between(t,%f,%f)'",
		op.Width, op.Height, op.X, op.Y, startSec, startSec+durationSec)

	args := []string{
		"-i", inputPath,
		"-i", op.PipVideoPath,
		"-filter_complex", filterComplex,
		"-c:a", "copy",
		"-y", outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeSplitScreenOperation 执行分屏操作
func (ve *VideoEditor) executeSplitScreenOperation(ctx context.Context, op *SplitScreenOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

	var filterComplex string
	var inputs []string

	// 构建输入参数
	inputs = append(inputs, "-i", inputPath)
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

	args := append(inputs, "-filter_complex", filterComplex, "-map", "[v]", "-c:a", "copy", "-y", outputPath)

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeBackgroundMusicOperation 执行背景音乐操作
func (ve *VideoEditor) executeBackgroundMusicOperation(ctx context.Context, op *BackgroundMusicOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

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
		"-i", inputPath,
		"-i", op.MusicPath,
		"-filter_complex", audioFilter,
		"-map", "0:v",
		"-map", "[a]",
		"-y", outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeSubtitleOperation 执行字幕操作
func (ve *VideoEditor) executeSubtitleOperation(ctx context.Context, op *SubtitleOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

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
		"-i", inputPath,
		"-vf", subtitleFilter,
		"-c:a", "copy",
		"-y", outputPath,
	}

	videoInfo, err := ve.ffmpeg.GetVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("获取视频信息失败: %w", err)
	}

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, videoInfo.Duration)
}

// executeSlideshowOperation 执行幻灯片操作
func (ve *VideoEditor) executeSlideshowOperation(ctx context.Context, op *SlideshowOperation,
	inputPath, outputPath string, callback ProgressCallback) error {

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

	args := append(inputs, "-filter_complex", filterComplex, "-map", "[v]", "-c:v", "libx264", "-pix_fmt", "yuv420p", "-y", outputPath)

	// 对于幻灯片，估算总时长
	estimatedDuration := time.Duration(len(op.Images)) * op.Duration

	return ve.ffmpeg.executeCommandWithProgress(ctx, args, callback, estimatedDuration)
}

// === 新增高级滤镜链式调用方法 ===

// ApplyColorGrading 应用色彩分级
func (ve *VideoEditor) ApplyColorGrading(options *ColorGradingOptions) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if ve.cancelled {
		return ve
	}

	operation := &ColorGradingOperation{
		BaseOperation: BaseOperation{
			Type:      OpTypeColorGrading,
			StartTime: 0,
			Duration:  0, // 应用到整个视频
		},
		Options: options,
	}

	ve.operations = append(ve.operations, operation)
	return ve
}

// ApplyFilter 应用通用滤镜
func (ve *VideoEditor) ApplyFilter(filterType FilterType, options *FilterOptions) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if ve.cancelled {
		return ve
	}

	if options == nil {
		options = &FilterOptions{
			FilterType: filterType,
			Intensity:  1.0,
		}
	} else {
		options.FilterType = filterType
	}

	operation := &FilterOperation{
		BaseOperation: BaseOperation{
			Type:      OpTypeFilter,
			StartTime: options.StartTime,
			Duration:  options.Duration,
		},
		Options: options,
	}

	ve.operations = append(ve.operations, operation)
	return ve
}

// ApplyVintageFilter 应用复古滤镜
func (ve *VideoEditor) ApplyVintageFilter(options *VintageFilterOptions) *VideoEditor {
	if options == nil {
		options = &VintageFilterOptions{
			Sepia:      0.5,
			Grain:      0.3,
			Vignette:   0.4,
			Fade:       0.2,
			ColorShift: 0.1,
			Desaturate: 0.3,
		}
	}

	filterOptions := &FilterOptions{
		FilterType: FilterTypeVintage,
		Intensity:  1.0,
		Parameters: map[string]interface{}{
			"sepia":       options.Sepia,
			"grain":       options.Grain,
			"vignette":    options.Vignette,
			"fade":        options.Fade,
			"scratches":   options.Scratches,
			"dust_spots":  options.DustSpots,
			"color_shift": options.ColorShift,
			"desaturate":  options.Desaturate,
		},
	}

	return ve.ApplyFilter(FilterTypeVintage, filterOptions)
}

// ApplyCinematicFilter 应用电影风格滤镜
func (ve *VideoEditor) ApplyCinematicFilter(options *CinematicFilterOptions) *VideoEditor {
	if options == nil {
		options = &CinematicFilterOptions{
			AspectRatio:    "21:9",
			LetterboxColor: "black",
			ColorGrading:   "teal_orange",
			FilmGrain:      0.2,
			Bloom:          0.3,
			LensFlare:      false,
			MotionBlur:     0.0,
		}
	}

	filterOptions := &FilterOptions{
		FilterType: FilterTypeCinematic,
		Intensity:  1.0,
		Parameters: map[string]interface{}{
			"aspect_ratio":    options.AspectRatio,
			"letterbox_color": options.LetterboxColor,
			"color_grading":   options.ColorGrading,
			"film_grain":      options.FilmGrain,
			"bloom":           options.Bloom,
			"lens_flare":      options.LensFlare,
			"motion_blur":     options.MotionBlur,
		},
	}

	return ve.ApplyFilter(FilterTypeCinematic, filterOptions)
}

// ApplyBeautyFilter 应用美颜滤镜
func (ve *VideoEditor) ApplyBeautyFilter(options *BeautyFilterOptions) *VideoEditor {
	if options == nil {
		options = &BeautyFilterOptions{
			SkinSmoothing:   0.5,
			SkinBrightening: 0.3,
			EyeEnhancement:  0.4,
			TeethWhitening:  0.2,
			FaceSlimming:    0.0,
			EyeEnlarging:    0.0,
			NoseReshaping:   0.0,
			LipEnhancement:  0.2,
		}
	}

	filterOptions := &FilterOptions{
		FilterType: FilterTypeBeauty,
		Intensity:  1.0,
		Parameters: map[string]interface{}{
			"skin_smoothing":   options.SkinSmoothing,
			"skin_brightening": options.SkinBrightening,
			"eye_enhancement":  options.EyeEnhancement,
			"teeth_whitening":  options.TeethWhitening,
			"face_slimming":    options.FaceSlimming,
			"eye_enlarging":    options.EyeEnlarging,
			"nose_reshaping":   options.NoseReshaping,
			"lip_enhancement":  options.LipEnhancement,
		},
	}

	return ve.ApplyFilter(FilterTypeBeauty, filterOptions)
}

// ApplySharpening 应用锐化滤镜
func (ve *VideoEditor) ApplySharpening(options *SharpeningOptions) *VideoEditor {
	if options == nil {
		options = &SharpeningOptions{
			Amount:    1.0,
			Radius:    1.0,
			Threshold: 0.0,
			Method:    "unsharp",
		}
	}

	filterOptions := &FilterOptions{
		FilterType: FilterTypeSharpening,
		Intensity:  1.0,
		Parameters: map[string]interface{}{
			"amount":    options.Amount,
			"radius":    options.Radius,
			"threshold": options.Threshold,
			"method":    options.Method,
		},
	}

	return ve.ApplyFilter(FilterTypeSharpening, filterOptions)
}

// ApplyDenoising 应用降噪滤镜
func (ve *VideoEditor) ApplyDenoising(options *DenoisingOptions) *VideoEditor {
	if options == nil {
		options = &DenoisingOptions{
			Strength:     0.5,
			Method:       "nlmeans",
			TemporalNR:   true,
			SpatialNR:    true,
			PreserveEdge: true,
		}
	}

	filterOptions := &FilterOptions{
		FilterType: FilterTypeDenoising,
		Intensity:  1.0,
		Parameters: map[string]interface{}{
			"strength":      options.Strength,
			"method":        options.Method,
			"temporal_nr":   options.TemporalNR,
			"spatial_nr":    options.SpatialNR,
			"preserve_edge": options.PreserveEdge,
		},
	}

	return ve.ApplyFilter(FilterTypeDenoising, filterOptions)
}

// === 转场效果链式调用方法 ===

// AddAdvancedTransition 添加高级转场效果
func (ve *VideoEditor) AddAdvancedTransition(transitionType AdvancedTransitionType, duration time.Duration, options *AdvancedTransitionOptions) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if ve.cancelled {
		return ve
	}

	if options == nil {
		options = &AdvancedTransitionOptions{
			Type:     transitionType,
			Duration: duration,
			Easing:   "ease_in_out",
		}
	} else {
		options.Type = transitionType
		options.Duration = duration
	}

	operation := &AdvancedTransitionOperation{
		BaseOperation: BaseOperation{
			Type:      OpTypeTransition,
			StartTime: 0,
			Duration:  duration,
		},
		Options: options,
	}

	ve.operations = append(ve.operations, operation)
	return ve
}

// AddFadeTransition 添加淡入淡出转场
func (ve *VideoEditor) AddFadeTransition(duration time.Duration) *VideoEditor {
	return ve.AddAdvancedTransition(TransitionFade, duration, nil)
}

// AddDissolveTransition 添加溶解转场
func (ve *VideoEditor) AddDissolveTransition(duration time.Duration) *VideoEditor {
	return ve.AddAdvancedTransition(TransitionDissolve, duration, nil)
}

// AddWipeTransition 添加擦除转场
func (ve *VideoEditor) AddWipeTransition(duration time.Duration, direction string) *VideoEditor {
	options := &AdvancedTransitionOptions{
		Direction: direction,
		Easing:    "linear",
	}
	return ve.AddAdvancedTransition(TransitionWipe, duration, options)
}

// AddSlideTransition 添加滑动转场
func (ve *VideoEditor) AddSlideTransition(duration time.Duration, direction string) *VideoEditor {
	options := &AdvancedTransitionOptions{
		Direction: direction,
		Easing:    "ease_out",
	}
	return ve.AddAdvancedTransition(TransitionSlide, duration, options)
}

// AddZoomTransition 添加缩放转场
func (ve *VideoEditor) AddZoomTransition(duration time.Duration, intensity float64) *VideoEditor {
	options := &AdvancedTransitionOptions{
		Intensity: intensity,
		Easing:    "ease_in_out",
	}
	return ve.AddAdvancedTransition(TransitionZoom, duration, options)
}

// AddGlitchTransition 添加故障效果转场
func (ve *VideoEditor) AddGlitchTransition(duration time.Duration, intensity float64) *VideoEditor {
	options := &AdvancedTransitionOptions{
		Intensity: intensity,
		Easing:    "linear",
	}
	return ve.AddAdvancedTransition(TransitionGlitch, duration, options)
}

// === 音频效果链式调用方法 ===

// ApplyAudioEqualizer 应用音频均衡器
func (ve *VideoEditor) ApplyAudioEqualizer(options *AudioEqualizerOptions) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if ve.cancelled {
		return ve
	}

	if options == nil {
		options = &AudioEqualizerOptions{
			Preset:     "balanced",
			MasterGain: 0,
		}
	}

	effectOptions := &AudioEffectOptions{
		EffectType: AudioEffectCompressor, // 使用压缩器作为示例
		Intensity:  1.0,
		Parameters: map[string]interface{}{
			"preset":      options.Preset,
			"bands":       options.Bands,
			"master_gain": options.MasterGain,
		},
	}

	operation := &AudioEffectOperation{
		BaseOperation: BaseOperation{
			Type:      OpTypeAudioEffect,
			StartTime: 0,
			Duration:  0, // 应用到整个音频
		},
		Options: effectOptions,
	}

	ve.operations = append(ve.operations, operation)
	return ve
}

// ApplyAudioEffect 应用音频效果
func (ve *VideoEditor) ApplyAudioEffect(effectType AudioEffectType, options *AudioEffectOptions) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if ve.cancelled {
		return ve
	}

	if options == nil {
		options = &AudioEffectOptions{
			EffectType: effectType,
			Intensity:  0.5,
		}
	} else {
		options.EffectType = effectType
	}

	operation := &AudioEffectOperation{
		BaseOperation: BaseOperation{
			Type:      OpTypeAudioEffect,
			StartTime: options.StartTime,
			Duration:  options.Duration,
		},
		Options: options,
	}

	ve.operations = append(ve.operations, operation)
	return ve
}

// ApplyReverb 应用混响效果
func (ve *VideoEditor) ApplyReverb(options *ReverbOptions) *VideoEditor {
	if options == nil {
		options = &ReverbOptions{
			RoomSize:   0.5,
			Damping:    0.5,
			WetLevel:   0.3,
			DryLevel:   0.7,
			PreDelay:   20,
			Diffusion:  0.5,
			ReverbType: "hall",
		}
	}

	effectOptions := &AudioEffectOptions{
		EffectType: AudioEffectReverb,
		Intensity:  options.WetLevel,
		Parameters: map[string]interface{}{
			"room_size":   options.RoomSize,
			"damping":     options.Damping,
			"wet_level":   options.WetLevel,
			"dry_level":   options.DryLevel,
			"pre_delay":   options.PreDelay,
			"diffusion":   options.Diffusion,
			"reverb_type": options.ReverbType,
		},
	}

	return ve.ApplyAudioEffect(AudioEffectReverb, effectOptions)
}

// ApplyCompressor 应用音频压缩器
func (ve *VideoEditor) ApplyCompressor(options *CompressorOptions) *VideoEditor {
	if options == nil {
		options = &CompressorOptions{
			Threshold:  -20,
			Ratio:      4.0,
			Attack:     5,
			Release:    50,
			MakeupGain: 0,
			KneeWidth:  2.5,
		}
	}

	effectOptions := &AudioEffectOptions{
		EffectType: AudioEffectCompressor,
		Intensity:  options.Ratio / 20.0,
		Parameters: map[string]interface{}{
			"threshold":   options.Threshold,
			"ratio":       options.Ratio,
			"attack":      options.Attack,
			"release":     options.Release,
			"makeup_gain": options.MakeupGain,
			"knee_width":  options.KneeWidth,
		},
	}

	return ve.ApplyAudioEffect(AudioEffectCompressor, effectOptions)
}

// === 高级字幕链式调用方法 ===

// AddAdvancedSubtitle 添加高级字幕
func (ve *VideoEditor) AddAdvancedSubtitle(options *AdvancedSubtitleOptions) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if ve.cancelled {
		return ve
	}

	if options == nil {
		return ve
	}

	operation := &AdvancedSubtitleOperation{
		BaseOperation: BaseOperation{
			Type:      OpTypeAdvancedSubtitle,
			StartTime: options.StartTime,
			Duration:  options.Duration,
		},
		Options: options,
	}

	ve.operations = append(ve.operations, operation)
	return ve
}

// AddAnimatedText 添加动画文字
func (ve *VideoEditor) AddAnimatedText(text string, animation SubtitleAnimationType, startTime, duration time.Duration) *VideoEditor {
	options := &AdvancedSubtitleOptions{
		Text:              text,
		StartTime:         startTime,
		Duration:          duration,
		X:                 100,
		Y:                 100,
		FontSize:          24,
		Color:             "#FFFFFF",
		Animation:         animation,
		AnimationDuration: 500 * time.Millisecond,
	}

	return ve.AddAdvancedSubtitle(options)
}

// AddLowerThird 添加下三分之一标题
func (ve *VideoEditor) AddLowerThird(title, subtitle string, startTime, duration time.Duration) *VideoEditor {
	// 主标题
	titleOptions := &AdvancedSubtitleOptions{
		Text:              title,
		StartTime:         startTime,
		Duration:          duration,
		X:                 50,
		Y:                 -150,
		FontSize:          32,
		FontWeight:        "bold",
		Color:             "#FFFFFF",
		BackgroundColor:   "#000000AA",
		Alignment:         "left",
		VerticalAlign:     "bottom",
		Animation:         SubtitleAnimationSlideIn,
		AnimationDuration: 500 * time.Millisecond,
	}

	ve.AddAdvancedSubtitle(titleOptions)

	// 副标题
	subtitleOptions := &AdvancedSubtitleOptions{
		Text:              subtitle,
		StartTime:         startTime,
		Duration:          duration,
		X:                 50,
		Y:                 -110,
		FontSize:          20,
		FontWeight:        "normal",
		Color:             "#CCCCCC",
		BackgroundColor:   "#000000AA",
		Alignment:         "left",
		VerticalAlign:     "bottom",
		Animation:         SubtitleAnimationSlideIn,
		AnimationDuration: 500 * time.Millisecond,
	}

	return ve.AddAdvancedSubtitle(subtitleOptions)
}

// ApplySubtitleTemplate 应用字幕模板
func (ve *VideoEditor) ApplySubtitleTemplate(text, templateName string, startTime, duration time.Duration) *VideoEditor {
	// 这里可以根据模板名称设置不同的样式
	var options *AdvancedSubtitleOptions

	switch templateName {
	case "standard":
		options = &AdvancedSubtitleOptions{
			Text:          text,
			StartTime:     startTime,
			Duration:      duration,
			FontSize:      24,
			FontWeight:    "normal",
			Color:         "#FFFFFF",
			OutlineColor:  "#000000",
			OutlineWidth:  2,
			ShadowColor:   "#000000",
			ShadowOffsetX: 2,
			ShadowOffsetY: 2,
			ShadowBlur:    3,
			Alignment:     "center",
			VerticalAlign: "bottom",
			Animation:     SubtitleAnimationFadeIn,
		}
	case "cinematic":
		options = &AdvancedSubtitleOptions{
			Text:            text,
			StartTime:       startTime,
			Duration:        duration,
			FontSize:        28,
			FontWeight:      "normal",
			Color:           "#F0F0F0",
			BackgroundColor: "#00000080",
			OutlineColor:    "#000000",
			OutlineWidth:    1,
			Alignment:       "center",
			VerticalAlign:   "bottom",
			LineSpacing:     1.2,
			Animation:       SubtitleAnimationTypewriter,
		}
	case "gaming":
		options = &AdvancedSubtitleOptions{
			Text:          text,
			StartTime:     startTime,
			Duration:      duration,
			FontSize:      26,
			FontWeight:    "bold",
			Color:         "#00FF00",
			OutlineColor:  "#000000",
			OutlineWidth:  3,
			ShadowColor:   "#00AA00",
			ShadowOffsetX: 0,
			ShadowOffsetY: 0,
			ShadowBlur:    5,
			Alignment:     "center",
			VerticalAlign: "bottom",
			Animation:     SubtitleAnimationGlow,
		}
	default:
		// 默认标准模板
		options = &AdvancedSubtitleOptions{
			Text:          text,
			StartTime:     startTime,
			Duration:      duration,
			FontSize:      24,
			Color:         "#FFFFFF",
			Alignment:     "center",
			VerticalAlign: "bottom",
			Animation:     SubtitleAnimationFadeIn,
		}
	}

	return ve.AddAdvancedSubtitle(options)
}

// === 高级合成链式调用方法 ===

// ApplyChromaKey 应用绿幕抠图
func (ve *VideoEditor) ApplyChromaKey(options *ChromaKeyOptions) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if ve.cancelled {
		return ve
	}

	if options == nil {
		options = &ChromaKeyOptions{
			KeyColor:         "#00FF00",
			Tolerance:        0.3,
			Softness:         0.1,
			SpillSuppression: 0.1,
			EdgeFeather:      0.0,
			LightWrap:        0.0,
			ColorCorrection:  false,
		}
	}

	operation := &ChromaKeyOperation{
		BaseOperation: BaseOperation{
			Type:      OpTypeChromaKey,
			StartTime: 0,
			Duration:  0, // 应用到整个视频
		},
		Options: options,
	}

	ve.operations = append(ve.operations, operation)
	return ve
}

// ApplyMask 应用遮罩
func (ve *VideoEditor) ApplyMask(options *MaskOptions) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if ve.cancelled {
		return ve
	}

	if options == nil {
		return ve
	}

	operation := &MaskOperation{
		BaseOperation: BaseOperation{
			Type:      OpTypeMask,
			StartTime: options.StartTime,
			Duration:  options.Duration,
		},
		Options: options,
	}

	ve.operations = append(ve.operations, operation)
	return ve
}

// AddParticleEffect 添加粒子效果
func (ve *VideoEditor) AddParticleEffect(options *ParticleEffectOptions) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if ve.cancelled {
		return ve
	}

	if options == nil {
		return ve
	}

	operation := &ParticleEffectOperation{
		BaseOperation: BaseOperation{
			Type:      OpTypeParticleEffect,
			StartTime: options.StartTime,
			Duration:  options.Duration,
		},
		Options: options,
	}

	ve.operations = append(ve.operations, operation)
	return ve
}

// CreateSnowEffect 创建雪花效果
func (ve *VideoEditor) CreateSnowEffect(intensity float64, startTime, duration time.Duration) *VideoEditor {
	options := &ParticleEffectOptions{
		ParticleType: ParticleTypeSnow,
		Count:        int(intensity * 200),
		Size:         2.0,
		Speed:        50.0,
		Direction:    270, // 向下
		Spread:       30,
		Gravity:      0.8,
		Wind:         0.1,
		Opacity:      0.8,
		Color:        "#FFFFFF",
		BlendMode:    "screen",
		StartTime:    startTime,
		Duration:     duration,
		EmissionRate: intensity * 10,
		LifeTime:     5.0,
	}

	return ve.AddParticleEffect(options)
}

// CreateRainEffect 创建雨滴效果
func (ve *VideoEditor) CreateRainEffect(intensity float64, startTime, duration time.Duration) *VideoEditor {
	options := &ParticleEffectOptions{
		ParticleType: ParticleTypeRain,
		Count:        int(intensity * 300),
		Size:         1.0,
		Speed:        100.0,
		Direction:    260, // 稍微倾斜
		Spread:       10,
		Gravity:      1.2,
		Wind:         0.3,
		Opacity:      0.6,
		Color:        "#87CEEB",
		BlendMode:    "overlay",
		StartTime:    startTime,
		Duration:     duration,
		EmissionRate: intensity * 20,
		LifeTime:     3.0,
	}

	return ve.AddParticleEffect(options)
}

// CreateFireEffect 创建火焰效果
func (ve *VideoEditor) CreateFireEffect(intensity float64, startTime, duration time.Duration) *VideoEditor {
	options := &ParticleEffectOptions{
		ParticleType: ParticleTypeFire,
		Count:        int(intensity * 100),
		Size:         3.0,
		Speed:        30.0,
		Direction:    90, // 向上
		Spread:       45,
		Gravity:      -0.5, // 负重力，向上飘
		Wind:         0.2,
		Opacity:      0.9,
		Color:        "#FF4500",
		BlendMode:    "screen",
		StartTime:    startTime,
		Duration:     duration,
		EmissionRate: intensity * 15,
		LifeTime:     2.0,
	}

	return ve.AddParticleEffect(options)
}

// CreateSparkleEffect 创建闪光效果
func (ve *VideoEditor) CreateSparkleEffect(intensity float64, startTime, duration time.Duration) *VideoEditor {
	options := &ParticleEffectOptions{
		ParticleType: ParticleTypeSparkle,
		Count:        int(intensity * 50),
		Size:         4.0,
		Speed:        10.0,
		Direction:    0,
		Spread:       360, // 全方向
		Gravity:      0.0,
		Wind:         0.0,
		Opacity:      1.0,
		Color:        "#FFD700",
		BlendMode:    "screen",
		StartTime:    startTime,
		Duration:     duration,
		EmissionRate: intensity * 5,
		LifeTime:     1.5,
	}

	return ve.AddParticleEffect(options)
}

// AddMotionGraphics 添加动态图形
func (ve *VideoEditor) AddMotionGraphics(options *MotionGraphicsOptions) *VideoEditor {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if ve.cancelled {
		return ve
	}

	if options == nil {
		return ve
	}

	operation := &MotionGraphicsOperation{
		BaseOperation: BaseOperation{
			Type:      OpTypeMotionGraphics,
			StartTime: options.StartTime,
			Duration:  options.Duration,
		},
		Options: options,
	}

	ve.operations = append(ve.operations, operation)
	return ve
}

// AddLowerThirdGraphics 添加下三分之一图形
func (ve *VideoEditor) AddLowerThirdGraphics(title string, startTime, duration time.Duration) *VideoEditor {
	options := &MotionGraphicsOptions{
		GraphicsType: MotionGraphicsLowerThird,
		Template:     "default",
		Text:         title,
		StartTime:    startTime,
		Duration:     duration,
		X:            50,
		Y:            -150,
		Width:        400,
		Height:       80,
		Color:        "#0066CC",
		AccentColor:  "#FFFFFF",
		Animation:    "slide_in",
	}

	return ve.AddMotionGraphics(options)
}

// AddProgressBar 添加进度条
func (ve *VideoEditor) AddProgressBar(startTime, duration time.Duration, x, y, width, height int) *VideoEditor {
	options := &MotionGraphicsOptions{
		GraphicsType: MotionGraphicsProgress,
		Template:     "default",
		Text:         "",
		StartTime:    startTime,
		Duration:     duration,
		X:            x,
		Y:            y,
		Width:        width,
		Height:       height,
		Color:        "#0066CC",
		AccentColor:  "#FFFFFF",
		Animation:    "linear",
	}

	return ve.AddMotionGraphics(options)
}

// AddCounter 添加计数器
func (ve *VideoEditor) AddCounter(startTime, duration time.Duration, x, y int) *VideoEditor {
	options := &MotionGraphicsOptions{
		GraphicsType: MotionGraphicsCounter,
		Template:     "default",
		Text:         "0",
		StartTime:    startTime,
		Duration:     duration,
		X:            x,
		Y:            y,
		Width:        100,
		Height:       50,
		Color:        "#FFFFFF",
		AccentColor:  "#0066CC",
		Animation:    "count_up",
	}

	return ve.AddMotionGraphics(options)
}
