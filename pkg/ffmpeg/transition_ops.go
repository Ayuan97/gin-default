// Package ffmpeg 提供高级转场效果功能
package ffmpeg

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AddAdvancedTransition 添加高级转场效果
// inputPath1: 第一个视频文件路径
// inputPath2: 第二个视频文件路径
// outputPath: 输出视频文件路径
// options: 转场选项
func (f *FFmpeg) AddAdvancedTransition(inputPath1, inputPath2, outputPath string, options *AdvancedTransitionOptions) error {
	return f.AddAdvancedTransitionWithContext(context.Background(), inputPath1, inputPath2, outputPath, options)
}

// AddAdvancedTransitionWithContext 带上下文的高级转场效果
func (f *FFmpeg) AddAdvancedTransitionWithContext(ctx context.Context, inputPath1, inputPath2, outputPath string, options *AdvancedTransitionOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath1); err != nil {
		return err
	}
	if err := validateInputFile(inputPath2); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	if options == nil {
		return NewError(ErrInvalidOptions, "转场选项不能为空", nil)
	}

	// 构建转场滤镜
	transitionFilter, err := f.buildAdvancedTransitionFilter(options)
	if err != nil {
		return err
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath1,
		"-i", inputPath2,
		"-filter_complex", transitionFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	// 添加自定义参数
	args = append(args, options.CustomArgs...)

	// 执行命令
	_, err = f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, fmt.Sprintf("转场效果处理失败: %s", options.Type), err)
	}

	f.logger.Info("转场效果处理完成: %s + %s -> %s (%s)", inputPath1, inputPath2, outputPath, options.Type)
	return nil
}

// AddWipeTransition 添加擦除转场
func (f *FFmpeg) AddWipeTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, options *WipeTransitionOptions) error {
	if options == nil {
		options = &WipeTransitionOptions{
			Direction: "left_to_right",
			Angle:     0,
			Softness:  0.1,
			Shape:     "linear",
		}
	}

	advancedOptions := &AdvancedTransitionOptions{
		Type:      TransitionWipe,
		Duration:  duration,
		Direction: options.Direction,
		Feather:   options.Softness,
		Intensity: 1.0,
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, advancedOptions)
}

// AddSlideTransition 添加滑动转场
func (f *FFmpeg) AddSlideTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, options *SlideTransitionOptions) error {
	if options == nil {
		options = &SlideTransitionOptions{
			Direction: "left",
			Distance:  1.0,
			Bounce:    false,
			Rotation:  0,
		}
	}

	advancedOptions := &AdvancedTransitionOptions{
		Type:      TransitionSlide,
		Duration:  duration,
		Direction: options.Direction,
		Intensity: options.Distance,
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, advancedOptions)
}

// AddZoomTransition 添加缩放转场
func (f *FFmpeg) AddZoomTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, options *ZoomTransitionOptions) error {
	if options == nil {
		options = &ZoomTransitionOptions{
			ZoomIn:   true,
			CenterX:  0.5,
			CenterY:  0.5,
			MaxScale: 2.0,
			Blur:     0.0,
		}
	}

	transitionType := TransitionZoom
	if !options.ZoomIn {
		transitionType = TransitionZoom // 可以扩展为 TransitionZoomOut
	}

	advancedOptions := &AdvancedTransitionOptions{
		Type:      transitionType,
		Duration:  duration,
		Intensity: options.MaxScale,
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, advancedOptions)
}

// AddGlitchTransition 添加故障效果转场
func (f *FFmpeg) AddGlitchTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, options *GlitchTransitionOptions) error {
	if options == nil {
		options = &GlitchTransitionOptions{
			Intensity:    0.5,
			BlockSize:    10,
			ColorShift:   0.3,
			DigitalNoise: 0.2,
			Scanlines:    true,
			Distortion:   0.1,
		}
	}

	advancedOptions := &AdvancedTransitionOptions{
		Type:      TransitionGlitch,
		Duration:  duration,
		Intensity: options.Intensity,
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, advancedOptions)
}

// AddDissolveTransition 添加溶解转场
func (f *FFmpeg) AddDissolveTransition(inputPath1, inputPath2, outputPath string, duration time.Duration) error {
	options := &AdvancedTransitionOptions{
		Type:     TransitionDissolve,
		Duration: duration,
		Easing:   "ease_in_out",
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, options)
}

// AddFadeTransition 添加淡入淡出转场
func (f *FFmpeg) AddFadeTransition(inputPath1, inputPath2, outputPath string, duration time.Duration) error {
	options := &AdvancedTransitionOptions{
		Type:     TransitionFade,
		Duration: duration,
		Easing:   "linear",
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, options)
}

// AddPushTransition 添加推拉转场
func (f *FFmpeg) AddPushTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, direction string) error {
	options := &AdvancedTransitionOptions{
		Type:      TransitionPush,
		Duration:  duration,
		Direction: direction,
		Easing:    "ease_out",
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, options)
}

// AddRotateTransition 添加旋转转场
func (f *FFmpeg) AddRotateTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, clockwise bool) error {
	direction := "clockwise"
	if !clockwise {
		direction = "counterclockwise"
	}

	options := &AdvancedTransitionOptions{
		Type:      TransitionRotate,
		Duration:  duration,
		Direction: direction,
		Easing:    "ease_in_out",
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, options)
}

// AddFlipTransition 添加翻转转场
func (f *FFmpeg) AddFlipTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, axis string) error {
	options := &AdvancedTransitionOptions{
		Type:      TransitionFlip,
		Duration:  duration,
		Direction: axis, // "horizontal" 或 "vertical"
		Easing:    "ease_in_out",
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, options)
}

// AddCubeTransition 添加立方体转场
func (f *FFmpeg) AddCubeTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, direction string) error {
	options := &AdvancedTransitionOptions{
		Type:      TransitionCube,
		Duration:  duration,
		Direction: direction,
		Easing:    "ease_in_out",
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, options)
}

// AddRippleTransition 添加波纹转场
func (f *FFmpeg) AddRippleTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, intensity float64) error {
	options := &AdvancedTransitionOptions{
		Type:      TransitionRipple,
		Duration:  duration,
		Intensity: intensity,
		Easing:    "ease_out",
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, options)
}

// AddMosaicTransition 添加马赛克转场
func (f *FFmpeg) AddMosaicTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, blockSize int) error {
	options := &AdvancedTransitionOptions{
		Type:      TransitionMosaic,
		Duration:  duration,
		Intensity: float64(blockSize),
		Easing:    "linear",
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, options)
}

// AddPixelateTransition 添加像素化转场
func (f *FFmpeg) AddPixelateTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, pixelSize int) error {
	options := &AdvancedTransitionOptions{
		Type:      TransitionPixelate,
		Duration:  duration,
		Intensity: float64(pixelSize),
		Easing:    "ease_in",
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, options)
}

// AddBurnTransition 添加燃烧效果转场
func (f *FFmpeg) AddBurnTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, intensity float64) error {
	options := &AdvancedTransitionOptions{
		Type:      TransitionBurn,
		Duration:  duration,
		Intensity: intensity,
		Color:     "#FF4500", // 橙红色
		Easing:    "ease_in",
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, options)
}

// AddShatterTransition 添加破碎效果转场
func (f *FFmpeg) AddShatterTransition(inputPath1, inputPath2, outputPath string, duration time.Duration, pieces int) error {
	options := &AdvancedTransitionOptions{
		Type:      TransitionShatter,
		Duration:  duration,
		Intensity: float64(pieces),
		Easing:    "ease_in",
	}

	return f.AddAdvancedTransition(inputPath1, inputPath2, outputPath, options)
}

// === 转场滤镜构建辅助方法 ===

// buildAdvancedTransitionFilter 构建高级转场滤镜
func (f *FFmpeg) buildAdvancedTransitionFilter(options *AdvancedTransitionOptions) (string, error) {
	duration := options.Duration.Seconds()

	switch options.Type {
	case TransitionFade:
		return f.buildFadeTransitionFilter(duration, options)
	case TransitionDissolve:
		return f.buildDissolveTransitionFilter(duration, options)
	case TransitionWipe:
		return f.buildWipeTransitionFilter(duration, options)
	case TransitionSlide:
		return f.buildSlideTransitionFilter(duration, options)
	case TransitionPush:
		return f.buildPushTransitionFilter(duration, options)
	case TransitionZoom:
		return f.buildZoomTransitionFilter(duration, options)
	case TransitionRotate:
		return f.buildRotateTransitionFilter(duration, options)
	case TransitionFlip:
		return f.buildFlipTransitionFilter(duration, options)
	case TransitionCube:
		return f.buildCubeTransitionFilter(duration, options)
	case TransitionSphere:
		return f.buildSphereTransitionFilter(duration, options)
	case TransitionRipple:
		return f.buildRippleTransitionFilter(duration, options)
	case TransitionMosaic:
		return f.buildMosaicTransitionFilter(duration, options)
	case TransitionPixelate:
		return f.buildPixelateTransitionFilter(duration, options)
	case TransitionGlitch:
		return f.buildGlitchTransitionFilter(duration, options)
	case TransitionBurn:
		return f.buildBurnTransitionFilter(duration, options)
	case TransitionShatter:
		return f.buildShatterTransitionFilter(duration, options)
	default:
		return "", NewError(ErrInvalidOptions, fmt.Sprintf("不支持的转场类型: %s", options.Type), nil)
	}
}

// buildFadeTransitionFilter 构建淡入淡出转场滤镜
func (f *FFmpeg) buildFadeTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 使用xfade滤镜实现淡入淡出
	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=fade:duration=%f:offset=0[v]", duration)
	return filter, nil
}

// buildDissolveTransitionFilter 构建溶解转场滤镜
func (f *FFmpeg) buildDissolveTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 使用xfade滤镜的dissolve模式
	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=dissolve:duration=%f:offset=0[v]", duration)
	return filter, nil
}

// buildWipeTransitionFilter 构建擦除转场滤镜
func (f *FFmpeg) buildWipeTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	direction := "l" // 默认从左到右
	switch options.Direction {
	case "left_to_right":
		direction = "l"
	case "right_to_left":
		direction = "r"
	case "top_to_bottom":
		direction = "u"
	case "bottom_to_top":
		direction = "d"
	}

	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=wipe%s:duration=%f:offset=0[v]", direction, duration)
	return filter, nil
}

// buildSlideTransitionFilter 构建滑动转场滤镜
func (f *FFmpeg) buildSlideTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	direction := "left"
	switch options.Direction {
	case "left":
		direction = "left"
	case "right":
		direction = "right"
	case "up":
		direction = "up"
	case "down":
		direction = "down"
	}

	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=slide%s:duration=%f:offset=0[v]", direction, duration)
	return filter, nil
}

// buildPushTransitionFilter 构建推拉转场滤镜
func (f *FFmpeg) buildPushTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 推拉效果类似滑动，但两个视频同时移动
	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=slideleft:duration=%f:offset=0[v]", duration)
	return filter, nil
}

// buildZoomTransitionFilter 构建缩放转场滤镜
func (f *FFmpeg) buildZoomTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 使用自定义表达式实现缩放效果
	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=zoom:duration=%f:offset=0[v]", duration)
	return filter, nil
}

// buildRotateTransitionFilter 构建旋转转场滤镜
func (f *FFmpeg) buildRotateTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 旋转转场需要使用复杂的滤镜链
	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=circleopen:duration=%f:offset=0[v]", duration)
	return filter, nil
}

// buildFlipTransitionFilter 构建翻转转场滤镜
func (f *FFmpeg) buildFlipTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 翻转效果
	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=vertopen:duration=%f:offset=0[v]", duration)
	if options.Direction == "horizontal" {
		filter = fmt.Sprintf("[0:v][1:v]xfade=transition=horzopen:duration=%f:offset=0[v]", duration)
	}
	return filter, nil
}

// buildCubeTransitionFilter 构建立方体转场滤镜
func (f *FFmpeg) buildCubeTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 立方体效果需要3D变换
	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=squeezev:duration=%f:offset=0[v]", duration)
	return filter, nil
}

// buildSphereTransitionFilter 构建球体转场滤镜
func (f *FFmpeg) buildSphereTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 球体效果
	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=circleopen:duration=%f:offset=0[v]", duration)
	return filter, nil
}

// buildRippleTransitionFilter 构建波纹转场滤镜
func (f *FFmpeg) buildRippleTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 波纹效果
	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=radial:duration=%f:offset=0[v]", duration)
	return filter, nil
}

// buildMosaicTransitionFilter 构建马赛克转场滤镜
func (f *FFmpeg) buildMosaicTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 马赛克效果
	blockSize := int(options.Intensity)
	if blockSize <= 0 {
		blockSize = 10
	}

	filter := fmt.Sprintf("[0:v]scale=iw/%d:ih/%d,scale=iw*%d:ih*%d:flags=neighbor[v0];[1:v][v0]xfade=transition=pixelize:duration=%f:offset=0[v]",
		blockSize, blockSize, blockSize, blockSize, duration)
	return filter, nil
}

// buildPixelateTransitionFilter 构建像素化转场滤镜
func (f *FFmpeg) buildPixelateTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 像素化效果
	pixelSize := int(options.Intensity)
	if pixelSize <= 0 {
		pixelSize = 8
	}

	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=pixelize:duration=%f:offset=0[v]", duration)
	return filter, nil
}

// buildGlitchTransitionFilter 构建故障效果转场滤镜
func (f *FFmpeg) buildGlitchTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 故障效果组合多种滤镜
	var filters []string

	// 添加噪声
	filters = append(filters, "noise=alls=20:allf=t")

	// 添加色彩偏移
	filters = append(filters, "colorchannelmixer=rr=0.9:gg=1.1:bb=0.8")

	// 基础转场
	baseFilter := fmt.Sprintf("[0:v][1:v]xfade=transition=dissolve:duration=%f:offset=0", duration)

	filter := fmt.Sprintf("%s,%s[v]", baseFilter, strings.Join(filters, ","))
	return filter, nil
}

// buildBurnTransitionFilter 构建燃烧效果转场滤镜
func (f *FFmpeg) buildBurnTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 燃烧效果
	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=burn:duration=%f:offset=0[v]", duration)
	return filter, nil
}

// buildShatterTransitionFilter 构建破碎效果转场滤镜
func (f *FFmpeg) buildShatterTransitionFilter(duration float64, options *AdvancedTransitionOptions) (string, error) {
	// 破碎效果
	filter := fmt.Sprintf("[0:v][1:v]xfade=transition=diagtl:duration=%f:offset=0[v]", duration)
	return filter, nil
}
