// Package ffmpeg 提供高级合成功能
package ffmpeg

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ApplyChromaKey 应用绿幕抠图
// inputPath: 输入视频文件路径
// outputPath: 输出视频文件路径
// options: 绿幕抠图选项
func (f *FFmpeg) ApplyChromaKey(inputPath, outputPath string, options *ChromaKeyOptions) error {
	return f.ApplyChromaKeyWithContext(context.Background(), inputPath, outputPath, options)
}

// ApplyChromaKeyWithContext 带上下文的绿幕抠图
func (f *FFmpeg) ApplyChromaKeyWithContext(ctx context.Context, inputPath, outputPath string, options *ChromaKeyOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	if options == nil {
		return NewError(ErrInvalidOptions, "绿幕抠图选项不能为空", nil)
	}

	// 构建绿幕抠图滤镜
	chromaKeyFilter, err := f.buildChromaKeyFilter(options)
	if err != nil {
		return err
	}

	var args []string

	// 如果有背景，需要两个输入
	if options.BackgroundPath != "" {
		if err := validateInputFile(options.BackgroundPath); err != nil {
			return err
		}

		args = []string{
			"-i", inputPath,
			"-i", options.BackgroundPath,
			"-filter_complex", chromaKeyFilter,
			"-c:a", "copy",
			"-y",
			outputPath,
		}
	} else {
		args = []string{
			"-i", inputPath,
			"-vf", chromaKeyFilter,
			"-c:a", "copy",
			"-y",
			outputPath,
		}
	}

	// 执行命令
	_, err = f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "绿幕抠图处理失败", err)
	}

	f.logger.Info("绿幕抠图处理完成: %s -> %s", inputPath, outputPath)
	return nil
}

// ApplyMask 应用遮罩
// inputPath: 输入视频文件路径
// outputPath: 输出视频文件路径
// options: 遮罩选项
func (f *FFmpeg) ApplyMask(inputPath, outputPath string, options *MaskOptions) error {
	return f.ApplyMaskWithContext(context.Background(), inputPath, outputPath, options)
}

// ApplyMaskWithContext 带上下文的遮罩应用
func (f *FFmpeg) ApplyMaskWithContext(ctx context.Context, inputPath, outputPath string, options *MaskOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	if options == nil {
		return NewError(ErrInvalidOptions, "遮罩选项不能为空", nil)
	}

	// 构建遮罩滤镜
	maskFilter, err := f.buildMaskFilter(options)
	if err != nil {
		return err
	}

	var args []string

	// 如果有遮罩文件，需要两个输入
	if options.MaskPath != "" {
		if err := validateInputFile(options.MaskPath); err != nil {
			return err
		}

		args = []string{
			"-i", inputPath,
			"-i", options.MaskPath,
			"-filter_complex", maskFilter,
			"-c:a", "copy",
			"-y",
			outputPath,
		}
	} else {
		args = []string{
			"-i", inputPath,
			"-vf", maskFilter,
			"-c:a", "copy",
			"-y",
			outputPath,
		}
	}

	// 执行命令
	_, err = f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "遮罩处理失败", err)
	}

	f.logger.Info("遮罩处理完成: %s -> %s", inputPath, outputPath)
	return nil
}

// AddParticleEffect 添加粒子效果
// inputPath: 输入视频文件路径
// outputPath: 输出视频文件路径
// options: 粒子效果选项
func (f *FFmpeg) AddParticleEffect(inputPath, outputPath string, options *ParticleEffectOptions) error {
	return f.AddParticleEffectWithContext(context.Background(), inputPath, outputPath, options)
}

// AddParticleEffectWithContext 带上下文的粒子效果添加
func (f *FFmpeg) AddParticleEffectWithContext(ctx context.Context, inputPath, outputPath string, options *ParticleEffectOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	if options == nil {
		return NewError(ErrInvalidOptions, "粒子效果选项不能为空", nil)
	}

	// 构建粒子效果滤镜
	particleFilter, err := f.buildParticleEffectFilter(options)
	if err != nil {
		return err
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath,
		"-vf", particleFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	// 执行命令
	_, err = f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, fmt.Sprintf("粒子效果处理失败: %s", options.ParticleType), err)
	}

	f.logger.Info("粒子效果处理完成: %s (%s) -> %s", options.ParticleType, inputPath, outputPath)
	return nil
}

// AddMotionGraphics 添加动态图形
// inputPath: 输入视频文件路径
// outputPath: 输出视频文件路径
// options: 动态图形选项
func (f *FFmpeg) AddMotionGraphics(inputPath, outputPath string, options *MotionGraphicsOptions) error {
	return f.AddMotionGraphicsWithContext(context.Background(), inputPath, outputPath, options)
}

// AddMotionGraphicsWithContext 带上下文的动态图形添加
func (f *FFmpeg) AddMotionGraphicsWithContext(ctx context.Context, inputPath, outputPath string, options *MotionGraphicsOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	if options == nil {
		return NewError(ErrInvalidOptions, "动态图形选项不能为空", nil)
	}

	// 构建动态图形滤镜
	motionFilter, err := f.buildMotionGraphicsFilter(options)
	if err != nil {
		return err
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath,
		"-vf", motionFilter,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	// 执行命令
	_, err = f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, fmt.Sprintf("动态图形处理失败: %s", options.GraphicsType), err)
	}

	f.logger.Info("动态图形处理完成: %s (%s) -> %s", options.GraphicsType, inputPath, outputPath)
	return nil
}

// CreateSnowEffect 创建雪花效果
func (f *FFmpeg) CreateSnowEffect(inputPath, outputPath string, intensity float64, duration time.Duration) error {
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
		StartTime:    0,
		Duration:     duration,
		EmissionRate: intensity * 10,
		LifeTime:     5.0,
	}

	return f.AddParticleEffect(inputPath, outputPath, options)
}

// CreateRainEffect 创建雨滴效果
func (f *FFmpeg) CreateRainEffect(inputPath, outputPath string, intensity float64, duration time.Duration) error {
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
		StartTime:    0,
		Duration:     duration,
		EmissionRate: intensity * 20,
		LifeTime:     3.0,
	}

	return f.AddParticleEffect(inputPath, outputPath, options)
}

// CreateFireEffect 创建火焰效果
func (f *FFmpeg) CreateFireEffect(inputPath, outputPath string, intensity float64, duration time.Duration) error {
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
		StartTime:    0,
		Duration:     duration,
		EmissionRate: intensity * 15,
		LifeTime:     2.0,
	}

	return f.AddParticleEffect(inputPath, outputPath, options)
}

// CreateSparkleEffect 创建闪光效果
func (f *FFmpeg) CreateSparkleEffect(inputPath, outputPath string, intensity float64, duration time.Duration) error {
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
		StartTime:    0,
		Duration:     duration,
		EmissionRate: intensity * 5,
		LifeTime:     1.5,
	}

	return f.AddParticleEffect(inputPath, outputPath, options)
}

// CreateBubbleEffect 创建气泡效果
func (f *FFmpeg) CreateBubbleEffect(inputPath, outputPath string, intensity float64, duration time.Duration) error {
	options := &ParticleEffectOptions{
		ParticleType: ParticleTypeBubbles,
		Count:        int(intensity * 80),
		Size:         5.0,
		Speed:        20.0,
		Direction:    90, // 向上
		Spread:       20,
		Gravity:      -0.3, // 轻微向上
		Wind:         0.1,
		Opacity:      0.7,
		Color:        "#87CEEB",
		BlendMode:    "overlay",
		StartTime:    0,
		Duration:     duration,
		EmissionRate: intensity * 8,
		LifeTime:     4.0,
	}

	return f.AddParticleEffect(inputPath, outputPath, options)
}

// === 合成滤镜构建辅助方法 ===

// buildChromaKeyFilter 构建绿幕抠图滤镜
func (f *FFmpeg) buildChromaKeyFilter(options *ChromaKeyOptions) (string, error) {
	// 解析颜色
	keyColor := strings.TrimPrefix(options.KeyColor, "#")
	if len(keyColor) != 6 {
		return "", NewError(ErrInvalidOptions, "无效的抠图颜色格式", nil)
	}

	var filters []string

	// 基础chromakey滤镜
	chromaFilter := fmt.Sprintf("chromakey=0x%s:similarity=%f:blend=%f",
		keyColor, options.Tolerance, options.Softness)
	filters = append(filters, chromaFilter)

	// 溢色抑制
	if options.SpillSuppression > 0 {
		spillFilter := fmt.Sprintf("despill=type=green:mix=%f:expand=%f",
			options.SpillSuppression, options.SpillSuppression*0.5)
		filters = append(filters, spillFilter)
	}

	// 边缘羽化
	if options.EdgeFeather > 0 {
		featherFilter := fmt.Sprintf("boxblur=%f:%f", options.EdgeFeather*2, options.EdgeFeather*2)
		filters = append(filters, featherFilter)
	}

	// 如果有背景，进行合成
	if options.BackgroundPath != "" {
		// 使用overlay滤镜合成
		baseFilter := strings.Join(filters, ",")
		compositeFilter := fmt.Sprintf("[0:v]%s[fg];[1:v][fg]overlay=0:0[v]", baseFilter)
		return compositeFilter, nil
	}

	return strings.Join(filters, ","), nil
}

// buildMaskFilter 构建遮罩滤镜
func (f *FFmpeg) buildMaskFilter(options *MaskOptions) (string, error) {
	var filters []string

	switch options.MaskType {
	case MaskTypeAlpha:
		if options.MaskPath != "" {
			// 使用外部遮罩文件
			maskFilter := "[0:v][1:v]alphamerge[masked]"
			if options.Invert {
				maskFilter = "[1:v]negate[inverted];[0:v][inverted]alphamerge[masked]"
			}
			filters = append(filters, maskFilter)
		}
	case MaskTypeLuma:
		if options.MaskPath != "" {
			// 使用亮度遮罩
			maskFilter := "[1:v]format=gray[mask];[0:v][mask]maskedmerge[masked]"
			if options.Invert {
				maskFilter = "[1:v]format=gray,negate[mask];[0:v][mask]maskedmerge[masked]"
			}
			filters = append(filters, maskFilter)
		}
	case MaskTypeColor:
		// 颜色遮罩（类似chromakey）
		maskFilter := fmt.Sprintf("chromakey=0x00FF00:similarity=0.3:blend=0.1")
		filters = append(filters, maskFilter)
	case MaskTypeShape:
		// 形状遮罩（圆形、矩形等）
		shapeFilter := f.buildShapeMask(options)
		if shapeFilter != "" {
			filters = append(filters, shapeFilter)
		}
	case MaskTypeGradient:
		// 渐变遮罩
		gradientFilter := f.buildGradientMask(options)
		if gradientFilter != "" {
			filters = append(filters, gradientFilter)
		}
	}

	// 羽化效果
	if options.Feather > 0 {
		featherFilter := fmt.Sprintf("boxblur=%f:%f", options.Feather*5, options.Feather*5)
		filters = append(filters, featherFilter)
	}

	// 透明度调整
	if options.Opacity > 0 && options.Opacity < 1 {
		opacityFilter := fmt.Sprintf("format=yuva420p,colorchannelmixer=aa=%f", options.Opacity)
		filters = append(filters, opacityFilter)
	}

	if len(filters) == 0 {
		return "", NewError(ErrInvalidOptions, "无效的遮罩参数", nil)
	}

	return strings.Join(filters, ","), nil
}

// buildShapeMask 构建形状遮罩
func (f *FFmpeg) buildShapeMask(options *MaskOptions) string {
	// 创建圆形遮罩示例
	return "geq=lum='if(lt(sqrt(pow(X-W/2,2)+pow(Y-H/2,2)),min(W,H)/3),255,0)':cb=128:cr=128"
}

// buildGradientMask 构建渐变遮罩
func (f *FFmpeg) buildGradientMask(options *MaskOptions) string {
	// 创建线性渐变遮罩示例
	return "geq=lum='255*X/W':cb=128:cr=128"
}

// buildParticleEffectFilter 构建粒子效果滤镜
func (f *FFmpeg) buildParticleEffectFilter(options *ParticleEffectOptions) (string, error) {
	switch options.ParticleType {
	case ParticleTypeSnow:
		return f.buildSnowParticleFilter(options), nil
	case ParticleTypeRain:
		return f.buildRainParticleFilter(options), nil
	case ParticleTypeFire:
		return f.buildFireParticleFilter(options), nil
	case ParticleTypeSmoke:
		return f.buildSmokeParticleFilter(options), nil
	case ParticleTypeSparkle:
		return f.buildSparkleParticleFilter(options), nil
	case ParticleTypeBubbles:
		return f.buildBubbleParticleFilter(options), nil
	case ParticleTypeLeaves:
		return f.buildLeavesParticleFilter(options), nil
	case ParticleTypeStars:
		return f.buildStarsParticleFilter(options), nil
	case ParticleTypeHearts:
		return f.buildHeartsParticleFilter(options), nil
	case ParticleTypeConfetti:
		return f.buildConfettiParticleFilter(options), nil
	default:
		return "", NewError(ErrInvalidOptions, fmt.Sprintf("不支持的粒子类型: %s", options.ParticleType), nil)
	}
}

// buildSnowParticleFilter 构建雪花粒子滤镜
func (f *FFmpeg) buildSnowParticleFilter(options *ParticleEffectOptions) string {
	// 使用noise和其他滤镜模拟雪花效果
	var filters []string

	// 生成噪声作为雪花
	filters = append(filters, "noise=alls=20:allf=t")

	// 添加运动模糊模拟下落
	filters = append(filters, "minterpolate=fps=25:mi_mode=mci:mc_mode=aobmc:me_mode=bidir:vsbmc=1")

	// 调整透明度
	if options.Opacity < 1.0 {
		filters = append(filters, fmt.Sprintf("format=yuva420p,colorchannelmixer=aa=%f", options.Opacity))
	}

	return strings.Join(filters, ",")
}

// buildRainParticleFilter 构建雨滴粒子滤镜
func (f *FFmpeg) buildRainParticleFilter(options *ParticleEffectOptions) string {
	// 使用线条和运动模糊模拟雨滴
	var filters []string

	// 生成垂直线条
	filters = append(filters, "geq=lum='if(mod(X,8)<2,255,0)':cb=128:cr=128")

	// 添加运动模糊
	filters = append(filters, "minterpolate=fps=25:mi_mode=mci")

	// 调整颜色
	color := strings.TrimPrefix(options.Color, "#")
	if color != "" {
		filters = append(filters, fmt.Sprintf("colorbalance=rs=0.2:gs=0.2:bs=0.6"))
	}

	return strings.Join(filters, ",")
}

// buildFireParticleFilter 构建火焰粒子滤镜
func (f *FFmpeg) buildFireParticleFilter(options *ParticleEffectOptions) string {
	// 使用噪声和颜色调整模拟火焰
	var filters []string

	// 生成噪声
	filters = append(filters, "noise=alls=30:allf=t")

	// 调整为火焰颜色
	filters = append(filters, "colorchannelmixer=rr=1.5:gg=0.8:bb=0.2")

	// 添加模糊和扭曲
	filters = append(filters, "gblur=sigma=2")

	return strings.Join(filters, ",")
}

// buildSmokeParticleFilter 构建烟雾粒子滤镜
func (f *FFmpeg) buildSmokeParticleFilter(options *ParticleEffectOptions) string {
	var filters []string

	// 生成烟雾纹理
	filters = append(filters, "noise=alls=15:allf=t")
	filters = append(filters, "gblur=sigma=5")
	filters = append(filters, "colorchannelmixer=rr=0.5:gg=0.5:bb=0.5")

	return strings.Join(filters, ",")
}

// buildSparkleParticleFilter 构建闪光粒子滤镜
func (f *FFmpeg) buildSparkleParticleFilter(options *ParticleEffectOptions) string {
	var filters []string

	// 生成随机闪光点
	filters = append(filters, "noise=alls=5:allf=t")
	filters = append(filters, "threshold=0.95")

	// 添加发光效果
	filters = append(filters, "gblur=sigma=3")

	// 调整为金色
	filters = append(filters, "colorchannelmixer=rr=1.5:gg=1.2:bb=0.3")

	return strings.Join(filters, ",")
}

// buildBubbleParticleFilter 构建气泡粒子滤镜
func (f *FFmpeg) buildBubbleParticleFilter(options *ParticleEffectOptions) string {
	var filters []string

	// 生成圆形气泡
	filters = append(filters, "geq=lum='if(lt(mod(sqrt(pow(X-W/2,2)+pow(Y-H/2,2)),50),5),255,0)':cb=128:cr=128")

	// 添加透明度
	filters = append(filters, "format=yuva420p,colorchannelmixer=aa=0.7")

	return strings.Join(filters, ",")
}

// buildLeavesParticleFilter 构建落叶粒子滤镜
func (f *FFmpeg) buildLeavesParticleFilter(options *ParticleEffectOptions) string {
	var filters []string

	// 生成叶子形状
	filters = append(filters, "noise=alls=10:allf=t")
	filters = append(filters, "colorchannelmixer=rr=0.8:gg=1.2:bb=0.3")

	return strings.Join(filters, ",")
}

// buildStarsParticleFilter 构建星星粒子滤镜
func (f *FFmpeg) buildStarsParticleFilter(options *ParticleEffectOptions) string {
	var filters []string

	// 生成星形
	filters = append(filters, "noise=alls=3:allf=t")
	filters = append(filters, "threshold=0.98")
	filters = append(filters, "gblur=sigma=1")

	return strings.Join(filters, ",")
}

// buildHeartsParticleFilter 构建爱心粒子滤镜
func (f *FFmpeg) buildHeartsParticleFilter(options *ParticleEffectOptions) string {
	var filters []string

	// 生成心形（简化版）
	filters = append(filters, "noise=alls=8:allf=t")
	filters = append(filters, "colorchannelmixer=rr=1.5:gg=0.3:bb=0.5")

	return strings.Join(filters, ",")
}

// buildConfettiParticleFilter 构建彩纸粒子滤镜
func (f *FFmpeg) buildConfettiParticleFilter(options *ParticleEffectOptions) string {
	var filters []string

	// 生成彩色碎片
	filters = append(filters, "noise=alls=25:allf=t")
	filters = append(filters, "hue=h=n*360/25:s=2")

	return strings.Join(filters, ",")
}

// buildMotionGraphicsFilter 构建动态图形滤镜
func (f *FFmpeg) buildMotionGraphicsFilter(options *MotionGraphicsOptions) (string, error) {
	switch options.GraphicsType {
	case MotionGraphicsLowerThird:
		return f.buildLowerThirdFilter(options), nil
	case MotionGraphicsCallout:
		return f.buildCalloutFilter(options), nil
	case MotionGraphicsProgress:
		return f.buildProgressFilter(options), nil
	case MotionGraphicsCounter:
		return f.buildCounterFilter(options), nil
	case MotionGraphicsChart:
		return f.buildChartFilter(options), nil
	case MotionGraphicsLogo:
		return f.buildLogoFilter(options), nil
	default:
		return "", NewError(ErrInvalidOptions, fmt.Sprintf("不支持的动态图形类型: %s", options.GraphicsType), nil)
	}
}

// buildLowerThirdFilter 构建下三分之一标题滤镜
func (f *FFmpeg) buildLowerThirdFilter(options *MotionGraphicsOptions) string {
	// 创建下三分之一标题背景和文字
	var filters []string

	// 背景矩形
	bgColor := strings.TrimPrefix(options.Color, "#")
	bgFilter := fmt.Sprintf("drawbox=x=%d:y=%d:w=%d:h=%d:color=0x%s@0.8:t=fill",
		options.X, options.Y, options.Width, options.Height, bgColor)
	filters = append(filters, bgFilter)

	// 主标题文字
	textFilter := fmt.Sprintf("drawtext=text='%s':x=%d:y=%d:fontsize=24:fontcolor=white",
		options.Text, options.X+20, options.Y+10)
	filters = append(filters, textFilter)

	return strings.Join(filters, ",")
}

// buildCalloutFilter 构建标注滤镜
func (f *FFmpeg) buildCalloutFilter(options *MotionGraphicsOptions) string {
	// 创建标注气泡和文字
	var filters []string

	// 标注背景
	bgFilter := fmt.Sprintf("drawbox=x=%d:y=%d:w=%d:h=%d:color=white@0.9:t=fill",
		options.X, options.Y, options.Width, options.Height)
	filters = append(filters, bgFilter)

	// 标注文字
	textFilter := fmt.Sprintf("drawtext=text='%s':x=%d:y=%d:fontsize=18:fontcolor=black",
		options.Text, options.X+10, options.Y+10)
	filters = append(filters, textFilter)

	return strings.Join(filters, ",")
}

// buildProgressFilter 构建进度条滤镜
func (f *FFmpeg) buildProgressFilter(options *MotionGraphicsOptions) string {
	// 创建进度条
	var filters []string

	// 进度条背景
	bgFilter := fmt.Sprintf("drawbox=x=%d:y=%d:w=%d:h=%d:color=gray@0.5:t=fill",
		options.X, options.Y, options.Width, options.Height)
	filters = append(filters, bgFilter)

	// 进度条前景（动态宽度）
	progressFilter := fmt.Sprintf("drawbox=x=%d:y=%d:w='%d*t/10':h=%d:color=blue:t=fill",
		options.X, options.Y, options.Width, options.Height)
	filters = append(filters, progressFilter)

	return strings.Join(filters, ",")
}

// buildCounterFilter 构建计数器滤镜
func (f *FFmpeg) buildCounterFilter(options *MotionGraphicsOptions) string {
	// 创建数字计数器
	counterFilter := fmt.Sprintf("drawtext=text='%%{eif\\:t\\:d}':x=%d:y=%d:fontsize=48:fontcolor=white",
		options.X, options.Y)

	return counterFilter
}

// buildChartFilter 构建图表滤镜
func (f *FFmpeg) buildChartFilter(options *MotionGraphicsOptions) string {
	// 创建简单的柱状图
	var filters []string

	// 图表背景
	bgFilter := fmt.Sprintf("drawbox=x=%d:y=%d:w=%d:h=%d:color=white@0.8:t=fill",
		options.X, options.Y, options.Width, options.Height)
	filters = append(filters, bgFilter)

	// 示例柱状图
	for i := 0; i < 5; i++ {
		barHeight := 20 + i*10
		barFilter := fmt.Sprintf("drawbox=x=%d:y=%d:w=30:h=%d:color=blue:t=fill",
			options.X+10+i*40, options.Y+options.Height-barHeight, barHeight)
		filters = append(filters, barFilter)
	}

	return strings.Join(filters, ",")
}

// buildLogoFilter 构建Logo动画滤镜
func (f *FFmpeg) buildLogoFilter(options *MotionGraphicsOptions) string {
	// 创建Logo动画（缩放进入）
	logoFilter := fmt.Sprintf("scale='iw*min(1,t/2)':'ih*min(1,t/2)'")

	return logoFilter
}
