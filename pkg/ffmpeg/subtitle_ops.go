// Package ffmpeg 提供高级字幕处理功能
package ffmpeg

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// AddAdvancedSubtitle 添加高级字幕
// inputPath: 输入视频文件路径
// outputPath: 输出视频文件路径
// options: 字幕选项
func (f *FFmpeg) AddAdvancedSubtitle(inputPath, outputPath string, options *AdvancedSubtitleOptions) error {
	return f.AddAdvancedSubtitleWithContext(context.Background(), inputPath, outputPath, options)
}

// AddAdvancedSubtitleWithContext 带上下文的高级字幕添加
func (f *FFmpeg) AddAdvancedSubtitleWithContext(ctx context.Context, inputPath, outputPath string, options *AdvancedSubtitleOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	if options == nil {
		return NewError(ErrInvalidOptions, "字幕选项不能为空", nil)
	}

	// 构建字幕滤镜
	subtitleFilter, err := f.buildAdvancedSubtitleFilter(options)
	if err != nil {
		return err
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath,
		"-vf", subtitleFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	// 执行命令
	_, err = f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "字幕添加失败", err)
	}

	f.logger.Info("字幕添加完成: %s -> %s", inputPath, outputPath)
	return nil
}

// AddSubtitleFromFile 从字幕文件添加字幕
// inputPath: 输入视频文件路径
// subtitlePath: 字幕文件路径 (支持SRT, ASS, VTT等格式)
// outputPath: 输出视频文件路径
func (f *FFmpeg) AddSubtitleFromFile(inputPath, subtitlePath, outputPath string) error {
	return f.AddSubtitleFromFileWithContext(context.Background(), inputPath, subtitlePath, outputPath)
}

// AddSubtitleFromFileWithContext 带上下文的字幕文件添加
func (f *FFmpeg) AddSubtitleFromFileWithContext(ctx context.Context, inputPath, subtitlePath, outputPath string) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}
	if err := validateInputFile(subtitlePath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	// 根据字幕文件格式选择处理方式
	ext := strings.ToLower(filepath.Ext(subtitlePath))
	var args []string

	switch ext {
	case ".srt":
		// SRT字幕文件
		args = []string{
			"-i", inputPath,
			"-vf", fmt.Sprintf("subtitles=%s", subtitlePath),
			"-c:a", "copy",
			"-y",
			outputPath,
		}
	case ".ass", ".ssa":
		// ASS/SSA字幕文件
		args = []string{
			"-i", inputPath,
			"-vf", fmt.Sprintf("ass=%s", subtitlePath),
			"-c:a", "copy",
			"-y",
			outputPath,
		}
	case ".vtt":
		// WebVTT字幕文件
		args = []string{
			"-i", inputPath,
			"-vf", fmt.Sprintf("subtitles=%s", subtitlePath),
			"-c:a", "copy",
			"-y",
			outputPath,
		}
	default:
		return NewError(ErrInvalidOptions, fmt.Sprintf("不支持的字幕文件格式: %s", ext), nil)
	}

	// 执行命令
	_, err := f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "字幕文件添加失败", err)
	}

	f.logger.Info("字幕文件添加完成: %s + %s -> %s", inputPath, subtitlePath, outputPath)
	return nil
}

// AddAnimatedText 添加动画文字
// inputPath: 输入视频文件路径
// outputPath: 输出视频文件路径
// text: 文字内容
// animation: 动画类型
// options: 字幕选项
func (f *FFmpeg) AddAnimatedText(inputPath, outputPath, text string, animation SubtitleAnimationType, options *AdvancedSubtitleOptions) error {
	if options == nil {
		options = &AdvancedSubtitleOptions{
			Text:      text,
			StartTime: 0,
			Duration:  5 * time.Second,
			X:         100,
			Y:         100,
			FontSize:  24,
			Color:     "#FFFFFF",
		}
	}

	options.Text = text
	options.Animation = animation

	return f.AddAdvancedSubtitle(inputPath, outputPath, options)
}

// AddLowerThird 添加下三分之一标题
// inputPath: 输入视频文件路径
// outputPath: 输出视频文件路径
// title: 主标题
// subtitle: 副标题
// duration: 显示时长
func (f *FFmpeg) AddLowerThird(inputPath, outputPath, title, subtitle string, duration time.Duration) error {
	// 主标题选项
	titleOptions := &AdvancedSubtitleOptions{
		Text:              title,
		StartTime:         0,
		Duration:          duration,
		X:                 50,
		Y:                 -150, // 距离底部150像素
		FontSize:          32,
		FontWeight:        "bold",
		Color:             "#FFFFFF",
		BackgroundColor:   "#000000AA",
		Alignment:         "left",
		VerticalAlign:     "bottom",
		Animation:         SubtitleAnimationSlideIn,
		AnimationDuration: 500 * time.Millisecond,
	}

	// 先添加主标题
	tempPath := outputPath + ".temp.mp4"
	err := f.AddAdvancedSubtitle(inputPath, tempPath, titleOptions)
	if err != nil {
		return err
	}

	// 副标题选项
	subtitleOptions := &AdvancedSubtitleOptions{
		Text:              subtitle,
		StartTime:         0,
		Duration:          duration,
		X:                 50,
		Y:                 -110, // 距离底部110像素
		FontSize:          20,
		FontWeight:        "normal",
		Color:             "#CCCCCC",
		BackgroundColor:   "#000000AA",
		Alignment:         "left",
		VerticalAlign:     "bottom",
		Animation:         SubtitleAnimationSlideIn,
		AnimationDuration: 500 * time.Millisecond,
	}

	// 添加副标题
	err = f.AddAdvancedSubtitle(tempPath, outputPath, subtitleOptions)
	if err != nil {
		return err
	}

	// 删除临时文件
	// 注意：这里应该使用适当的文件删除方法
	f.logger.Info("下三分之一标题添加完成: %s -> %s", inputPath, outputPath)
	return nil
}

// CreateSubtitleTemplate 创建字幕模板
func (f *FFmpeg) CreateSubtitleTemplate(name, description string, style *AdvancedSubtitleOptions, animation SubtitleAnimationType) *SubtitleTemplate {
	return &SubtitleTemplate{
		Name:        name,
		Description: description,
		Style:       style,
		Animation:   animation,
	}
}

// GetBuiltinSubtitleTemplates 获取内置字幕模板
func (f *FFmpeg) GetBuiltinSubtitleTemplates() map[string]*SubtitleTemplate {
	templates := make(map[string]*SubtitleTemplate)

	// 标准模板
	templates["standard"] = &SubtitleTemplate{
		Name:        "标准字幕",
		Description: "简洁的白色字幕，适合大多数场景",
		Style: &AdvancedSubtitleOptions{
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
		},
		Animation: SubtitleAnimationFadeIn,
	}

	// 电影模板
	templates["cinematic"] = &SubtitleTemplate{
		Name:        "电影字幕",
		Description: "电影风格的字幕，带有优雅的动画效果",
		Style: &AdvancedSubtitleOptions{
			FontSize:        28,
			FontWeight:      "normal",
			Color:           "#F0F0F0",
			BackgroundColor: "#00000080",
			OutlineColor:    "#000000",
			OutlineWidth:    1,
			Alignment:       "center",
			VerticalAlign:   "bottom",
			LineSpacing:     1.2,
		},
		Animation: SubtitleAnimationTypewriter,
	}

	// 游戏模板
	templates["gaming"] = &SubtitleTemplate{
		Name:        "游戏字幕",
		Description: "适合游戏视频的动感字幕",
		Style: &AdvancedSubtitleOptions{
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
		},
		Animation: SubtitleAnimationGlow,
	}

	// 教育模板
	templates["educational"] = &SubtitleTemplate{
		Name:        "教育字幕",
		Description: "清晰易读的教育视频字幕",
		Style: &AdvancedSubtitleOptions{
			FontSize:        22,
			FontWeight:      "normal",
			Color:           "#333333",
			BackgroundColor: "#FFFFFFCC",
			OutlineColor:    "#FFFFFF",
			OutlineWidth:    1,
			Alignment:       "center",
			VerticalAlign:   "bottom",
			LineSpacing:     1.3,
		},
		Animation: SubtitleAnimationSlideIn,
	}

	// 新闻模板
	templates["news"] = &SubtitleTemplate{
		Name:        "新闻字幕",
		Description: "专业的新闻节目字幕样式",
		Style: &AdvancedSubtitleOptions{
			FontSize:        20,
			FontWeight:      "bold",
			Color:           "#FFFFFF",
			BackgroundColor: "#0066CC",
			Alignment:       "left",
			VerticalAlign:   "bottom",
			LineSpacing:     1.1,
		},
		Animation: SubtitleAnimationSlideIn,
	}

	return templates
}

// ApplySubtitleTemplate 应用字幕模板
func (f *FFmpeg) ApplySubtitleTemplate(inputPath, outputPath, text string, templateName string, startTime, duration time.Duration) error {
	templates := f.GetBuiltinSubtitleTemplates()
	template, exists := templates[templateName]
	if !exists {
		return NewError(ErrInvalidOptions, fmt.Sprintf("字幕模板不存在: %s", templateName), nil)
	}

	// 复制模板样式
	options := *template.Style
	options.Text = text
	options.StartTime = startTime
	options.Duration = duration
	options.Animation = template.Animation

	return f.AddAdvancedSubtitle(inputPath, outputPath, &options)
}

// === 字幕滤镜构建辅助方法 ===

// buildAdvancedSubtitleFilter 构建高级字幕滤镜
func (f *FFmpeg) buildAdvancedSubtitleFilter(options *AdvancedSubtitleOptions) (string, error) {
	// 构建基础drawtext滤镜
	var filterParts []string

	// 文本内容
	text := strings.ReplaceAll(options.Text, "'", "\\'")
	text = strings.ReplaceAll(text, ":", "\\:")
	filterParts = append(filterParts, fmt.Sprintf("text='%s'", text))

	// 字体设置
	if options.FontFamily != "" {
		filterParts = append(filterParts, fmt.Sprintf("fontfile='%s'", options.FontFamily))
	}
	if options.FontSize > 0 {
		filterParts = append(filterParts, fmt.Sprintf("fontsize=%d", options.FontSize))
	}
	if options.FontWeight != "" && options.FontWeight != "normal" {
		// FFmpeg的drawtext不直接支持fontweight，需要通过字体文件实现
	}

	// 颜色设置
	if options.Color != "" {
		color := strings.TrimPrefix(options.Color, "#")
		filterParts = append(filterParts, fmt.Sprintf("fontcolor=0x%s", color))
	}

	// 位置设置
	x := f.buildPositionExpression(options.X, options.Alignment, "w")
	y := f.buildPositionExpression(options.Y, options.VerticalAlign, "h")
	filterParts = append(filterParts, fmt.Sprintf("x=%s", x))
	filterParts = append(filterParts, fmt.Sprintf("y=%s", y))

	// 描边设置
	if options.OutlineWidth > 0 {
		filterParts = append(filterParts, fmt.Sprintf("borderw=%d", options.OutlineWidth))
		if options.OutlineColor != "" {
			color := strings.TrimPrefix(options.OutlineColor, "#")
			filterParts = append(filterParts, fmt.Sprintf("bordercolor=0x%s", color))
		}
	}

	// 阴影设置
	if options.ShadowOffsetX != 0 || options.ShadowOffsetY != 0 {
		filterParts = append(filterParts, fmt.Sprintf("shadowx=%d", options.ShadowOffsetX))
		filterParts = append(filterParts, fmt.Sprintf("shadowy=%d", options.ShadowOffsetY))
		if options.ShadowColor != "" {
			color := strings.TrimPrefix(options.ShadowColor, "#")
			filterParts = append(filterParts, fmt.Sprintf("shadowcolor=0x%s", color))
		}
	}

	// 背景设置
	if options.BackgroundColor != "" {
		color := strings.TrimPrefix(options.BackgroundColor, "#")
		filterParts = append(filterParts, fmt.Sprintf("box=1:boxcolor=0x%s", color))
		filterParts = append(filterParts, "boxborderw=5")
	}

	// 时间设置
	if options.StartTime > 0 {
		startSeconds := options.StartTime.Seconds()
		filterParts = append(filterParts, fmt.Sprintf("enable='gte(t,%.2f)'", startSeconds))
	}
	if options.Duration > 0 {
		startSeconds := options.StartTime.Seconds()
		endSeconds := startSeconds + options.Duration.Seconds()
		filterParts = append(filterParts, fmt.Sprintf("enable='between(t,%.2f,%.2f)'", startSeconds, endSeconds))
	}

	// 透明度设置
	if options.Opacity > 0 && options.Opacity < 1 {
		alpha := int(options.Opacity * 255)
		filterParts = append(filterParts, fmt.Sprintf("alpha=%d", alpha))
	}

	// 旋转设置
	if options.Rotation != 0 {
		filterParts = append(filterParts, fmt.Sprintf("angle=%f*PI/180", options.Rotation))
	}

	// 构建基础滤镜
	baseFilter := fmt.Sprintf("drawtext=%s", strings.Join(filterParts, ":"))

	// 添加动画效果
	if options.Animation != "" {
		animatedFilter, err := f.buildSubtitleAnimation(baseFilter, options)
		if err != nil {
			return "", err
		}
		return animatedFilter, nil
	}

	return baseFilter, nil
}

// buildPositionExpression 构建位置表达式
func (f *FFmpeg) buildPositionExpression(pos int, align string, dimension string) string {
	if pos < 0 {
		// 负值表示从右边或底部计算
		return fmt.Sprintf("%s%d", dimension, pos)
	}

	switch align {
	case "center":
		if dimension == "w" {
			return fmt.Sprintf("(w-text_w)/2+%d", pos)
		} else {
			return fmt.Sprintf("(h-text_h)/2+%d", pos)
		}
	case "right", "bottom":
		if dimension == "w" {
			return fmt.Sprintf("w-text_w-%d", pos)
		} else {
			return fmt.Sprintf("h-text_h-%d", pos)
		}
	default: // left, top
		return fmt.Sprintf("%d", pos)
	}
}

// buildSubtitleAnimation 构建字幕动画
func (f *FFmpeg) buildSubtitleAnimation(baseFilter string, options *AdvancedSubtitleOptions) (string, error) {
	startTime := options.StartTime.Seconds()
	animDuration := options.AnimationDuration.Seconds()
	if animDuration == 0 {
		animDuration = 0.5 // 默认动画时长
	}

	switch options.Animation {
	case SubtitleAnimationFadeIn:
		return f.buildFadeInAnimation(baseFilter, startTime, animDuration), nil
	case SubtitleAnimationFadeOut:
		return f.buildFadeOutAnimation(baseFilter, startTime, options.Duration.Seconds(), animDuration), nil
	case SubtitleAnimationSlideIn:
		return f.buildSlideInAnimation(baseFilter, startTime, animDuration, options), nil
	case SubtitleAnimationSlideOut:
		return f.buildSlideOutAnimation(baseFilter, startTime, options.Duration.Seconds(), animDuration, options), nil
	case SubtitleAnimationZoomIn:
		return f.buildZoomInAnimation(baseFilter, startTime, animDuration), nil
	case SubtitleAnimationZoomOut:
		return f.buildZoomOutAnimation(baseFilter, startTime, options.Duration.Seconds(), animDuration), nil
	case SubtitleAnimationTypewriter:
		return f.buildTypewriterAnimation(baseFilter, startTime, options), nil
	case SubtitleAnimationBounce:
		return f.buildBounceAnimation(baseFilter, startTime, animDuration), nil
	case SubtitleAnimationRotate:
		return f.buildRotateAnimation(baseFilter, startTime, animDuration), nil
	case SubtitleAnimationGlow:
		return f.buildGlowAnimation(baseFilter, startTime, animDuration), nil
	default:
		return baseFilter, nil
	}
}

// buildFadeInAnimation 构建淡入动画
func (f *FFmpeg) buildFadeInAnimation(baseFilter string, startTime, duration float64) string {
	// 在原有滤镜基础上添加alpha动画
	alphaExpr := fmt.Sprintf("if(lt(t,%.2f),0,if(lt(t,%.2f),(t-%.2f)/%.2f,1))",
		startTime, startTime+duration, startTime, duration)
	return strings.Replace(baseFilter, "drawtext=", fmt.Sprintf("drawtext=alpha='%s':", alphaExpr), 1)
}

// buildFadeOutAnimation 构建淡出动画
func (f *FFmpeg) buildFadeOutAnimation(baseFilter string, startTime, totalDuration, animDuration float64) string {
	fadeStartTime := startTime + totalDuration - animDuration
	alphaExpr := fmt.Sprintf("if(lt(t,%.2f),1,if(lt(t,%.2f),1-(t-%.2f)/%.2f,0))",
		fadeStartTime, fadeStartTime+animDuration, fadeStartTime, animDuration)
	return strings.Replace(baseFilter, "drawtext=", fmt.Sprintf("drawtext=alpha='%s':", alphaExpr), 1)
}

// buildSlideInAnimation 构建滑入动画
func (f *FFmpeg) buildSlideInAnimation(baseFilter string, startTime, duration float64, options *AdvancedSubtitleOptions) string {
	// 从左侧滑入
	xExpr := fmt.Sprintf("if(lt(t,%.2f),-text_w,if(lt(t,%.2f),-text_w+(t-%.2f)/%.2f*(text_w+%d),%d))",
		startTime, startTime+duration, startTime, duration, options.X, options.X)

	// 替换x坐标表达式
	newFilter := strings.Replace(baseFilter, fmt.Sprintf("x=%d", options.X), fmt.Sprintf("x='%s'", xExpr), 1)
	return newFilter
}

// buildSlideOutAnimation 构建滑出动画
func (f *FFmpeg) buildSlideOutAnimation(baseFilter string, startTime, totalDuration, animDuration float64, options *AdvancedSubtitleOptions) string {
	slideStartTime := startTime + totalDuration - animDuration
	xExpr := fmt.Sprintf("if(lt(t,%.2f),%d,if(lt(t,%.2f),%d+(t-%.2f)/%.2f*(w+text_w),w))",
		slideStartTime, options.X, slideStartTime+animDuration, options.X, slideStartTime, animDuration)

	// 替换x坐标表达式
	newFilter := strings.Replace(baseFilter, fmt.Sprintf("x=%d", options.X), fmt.Sprintf("x='%s'", xExpr), 1)
	return newFilter
}

// buildZoomInAnimation 构建缩放进入动画
func (f *FFmpeg) buildZoomInAnimation(baseFilter string, startTime, duration float64) string {
	// 添加缩放效果（通过fontsize实现）
	sizeExpr := fmt.Sprintf("if(lt(t,%.2f),0,if(lt(t,%.2f),(t-%.2f)/%.2f*fontsize,fontsize))",
		startTime, startTime+duration, startTime, duration)
	return strings.Replace(baseFilter, "fontsize=", fmt.Sprintf("fontsize='%s':", sizeExpr), 1)
}

// buildZoomOutAnimation 构建缩放退出动画
func (f *FFmpeg) buildZoomOutAnimation(baseFilter string, startTime, totalDuration, animDuration float64) string {
	zoomStartTime := startTime + totalDuration - animDuration
	sizeExpr := fmt.Sprintf("if(lt(t,%.2f),fontsize,if(lt(t,%.2f),fontsize*(1-(t-%.2f)/%.2f),0))",
		zoomStartTime, zoomStartTime+animDuration, zoomStartTime, animDuration)
	return strings.Replace(baseFilter, "fontsize=", fmt.Sprintf("fontsize='%s':", sizeExpr), 1)
}

// buildTypewriterAnimation 构建打字机动画
func (f *FFmpeg) buildTypewriterAnimation(baseFilter string, startTime float64, options *AdvancedSubtitleOptions) string {
	textLen := len(options.Text)
	duration := options.Duration.Seconds()
	charDuration := duration / float64(textLen)

	// 使用text表达式实现打字机效果
	textExpr := fmt.Sprintf("substr('%s',0,max(0,min(%d,floor((t-%.2f)/%.2f))))",
		options.Text, textLen, startTime, charDuration)

	return strings.Replace(baseFilter, fmt.Sprintf("text='%s'", options.Text), fmt.Sprintf("text='%s'", textExpr), 1)
}

// buildBounceAnimation 构建弹跳动画
func (f *FFmpeg) buildBounceAnimation(baseFilter string, startTime, duration float64) string {
	// 使用sin函数创建弹跳效果
	yOffset := fmt.Sprintf("sin((t-%.2f)*2*PI/%.2f)*20", startTime, duration)

	// 在y坐标中添加偏移
	if strings.Contains(baseFilter, "y=") {
		// 找到y=部分并添加偏移
		return strings.Replace(baseFilter, "y=", fmt.Sprintf("y='%s'+", yOffset), 1)
	}

	return baseFilter
}

// buildRotateAnimation 构建旋转动画
func (f *FFmpeg) buildRotateAnimation(baseFilter string, startTime, duration float64) string {
	// 添加旋转角度表达式
	angleExpr := fmt.Sprintf("(t-%.2f)*2*PI/%.2f", startTime, duration)

	if strings.Contains(baseFilter, "angle=") {
		return strings.Replace(baseFilter, "angle=", fmt.Sprintf("angle='%s'+", angleExpr), 1)
	} else {
		return baseFilter + fmt.Sprintf(":angle='%s'", angleExpr)
	}
}

// buildGlowAnimation 构建发光动画
func (f *FFmpeg) buildGlowAnimation(baseFilter string, startTime, duration float64) string {
	// 使用sin函数创建发光效果（通过alpha变化）
	alphaExpr := fmt.Sprintf("0.5+0.5*sin((t-%.2f)*4*PI/%.2f)", startTime, duration)

	if strings.Contains(baseFilter, "alpha=") {
		return strings.Replace(baseFilter, "alpha=", fmt.Sprintf("alpha='%s'*", alphaExpr), 1)
	} else {
		return baseFilter + fmt.Sprintf(":alpha='%s'", alphaExpr)
	}
}
