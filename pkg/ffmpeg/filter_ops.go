// Package ffmpeg 提供高级滤镜操作功能
package ffmpeg

import (
	"context"
	"fmt"
	"strings"
)

// ApplyColorGrading 应用色彩分级
// inputPath: 输入视频文件路径
// outputPath: 输出视频文件路径
// options: 色彩分级选项
func (f *FFmpeg) ApplyColorGrading(inputPath, outputPath string, options *ColorGradingOptions) error {
	return f.ApplyColorGradingWithContext(context.Background(), inputPath, outputPath, options)
}

// ApplyColorGradingWithContext 带上下文的色彩分级
func (f *FFmpeg) ApplyColorGradingWithContext(ctx context.Context, inputPath, outputPath string, options *ColorGradingOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	if options == nil {
		return NewError(ErrInvalidOptions, "色彩分级选项不能为空", nil)
	}

	// 构建色彩分级滤镜
	var filters []string

	// 基础色彩调整
	if options.Brightness != 0 || options.Contrast != 0 || options.Saturation != 0 || options.Gamma != 1.0 {
		eqFilter := fmt.Sprintf("eq=brightness=%f:contrast=%f:saturation=%f:gamma=%f",
			options.Brightness, options.Contrast, options.Saturation, options.Gamma)
		filters = append(filters, eqFilter)
	}

	// 色调调整
	if options.Hue != 0 {
		hueFilter := fmt.Sprintf("hue=h=%f", options.Hue)
		filters = append(filters, hueFilter)
	}

	// 色温调整
	if options.Temperature != 0 {
		tempFilter := f.buildTemperatureFilter(options.Temperature)
		if tempFilter != "" {
			filters = append(filters, tempFilter)
		}
	}

	// 高光和阴影调整
	if options.Highlights != 0 || options.Shadows != 0 {
		shadowHighlightFilter := f.buildShadowHighlightFilter(options.Highlights, options.Shadows)
		if shadowHighlightFilter != "" {
			filters = append(filters, shadowHighlightFilter)
		}
	}

	// 白色和黑色调整
	if options.Whites != 0 || options.Blacks != 0 {
		levelsFilter := f.buildLevelsFilter(options.Whites, options.Blacks)
		if levelsFilter != "" {
			filters = append(filters, levelsFilter)
		}
	}

	// 清晰度调整
	if options.Clarity != 0 {
		clarityFilter := f.buildClarityFilter(options.Clarity)
		if clarityFilter != "" {
			filters = append(filters, clarityFilter)
		}
	}

	// 自然饱和度调整
	if options.Vibrance != 0 {
		vibranceFilter := f.buildVibranceFilter(options.Vibrance)
		if vibranceFilter != "" {
			filters = append(filters, vibranceFilter)
		}
	}

	if len(filters) == 0 {
		return NewError(ErrInvalidOptions, "没有有效的色彩分级参数", nil)
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath,
		"-vf", strings.Join(filters, ","),
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	// 执行命令
	_, err := f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "色彩分级处理失败", err)
	}

	f.logger.Info("色彩分级处理完成: %s -> %s", inputPath, outputPath)
	return nil
}

// ApplyFilter 应用通用滤镜
func (f *FFmpeg) ApplyFilter(inputPath, outputPath string, options *FilterOptions) error {
	return f.ApplyFilterWithContext(context.Background(), inputPath, outputPath, options)
}

// ApplyFilterWithContext 带上下文的通用滤镜应用
func (f *FFmpeg) ApplyFilterWithContext(ctx context.Context, inputPath, outputPath string, options *FilterOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	if options == nil {
		return NewError(ErrInvalidOptions, "滤镜选项不能为空", nil)
	}

	var filterString string
	var err error

	// 根据滤镜类型构建滤镜字符串
	switch options.FilterType {
	case FilterTypeVintage:
		filterString, err = f.buildVintageFilter(options)
	case FilterTypeCinematic:
		filterString, err = f.buildCinematicFilter(options)
	case FilterTypeBeauty:
		filterString, err = f.buildBeautyFilter(options)
	case FilterTypeSharpening:
		filterString, err = f.buildSharpeningFilter(options)
	case FilterTypeDenoising:
		filterString, err = f.buildDenoisingFilter(options)
	case FilterTypeVignette:
		filterString, err = f.buildVignetteFilter(options)
	case FilterTypeGlow:
		filterString, err = f.buildGlowFilter(options)
	case FilterTypeBloom:
		filterString, err = f.buildBloomFilter(options)
	default:
		return NewError(ErrInvalidOptions, fmt.Sprintf("不支持的滤镜类型: %s", options.FilterType), nil)
	}

	if err != nil {
		return err
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath,
		"-vf", filterString,
		"-c:a", "copy",
		"-y",
		outputPath,
	}

	// 添加自定义参数
	args = append(args, options.CustomArgs...)

	// 执行命令
	_, err = f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, fmt.Sprintf("滤镜应用失败: %s", options.FilterType), err)
	}

	f.logger.Info("滤镜应用完成: %s (%s) -> %s", options.FilterType, inputPath, outputPath)
	return nil
}

// ApplyVintageFilter 应用复古滤镜
func (f *FFmpeg) ApplyVintageFilter(inputPath, outputPath string, options *VintageFilterOptions) error {
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

	return f.ApplyFilter(inputPath, outputPath, filterOptions)
}

// ApplyCinematicFilter 应用电影风格滤镜
func (f *FFmpeg) ApplyCinematicFilter(inputPath, outputPath string, options *CinematicFilterOptions) error {
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

	return f.ApplyFilter(inputPath, outputPath, filterOptions)
}

// ApplyBeautyFilter 应用美颜滤镜
func (f *FFmpeg) ApplyBeautyFilter(inputPath, outputPath string, options *BeautyFilterOptions) error {
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

	return f.ApplyFilter(inputPath, outputPath, filterOptions)
}

// ApplySharpening 应用锐化滤镜
func (f *FFmpeg) ApplySharpening(inputPath, outputPath string, options *SharpeningOptions) error {
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

	return f.ApplyFilter(inputPath, outputPath, filterOptions)
}

// ApplyDenoising 应用降噪滤镜
func (f *FFmpeg) ApplyDenoising(inputPath, outputPath string, options *DenoisingOptions) error {
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

	return f.ApplyFilter(inputPath, outputPath, filterOptions)
}

// === 滤镜构建辅助方法 ===

// buildTemperatureFilter 构建色温调整滤镜
func (f *FFmpeg) buildTemperatureFilter(temperature float64) string {
	if temperature == 0 {
		return ""
	}

	// 色温调整通过调整RGB通道实现
	// 正值增加暖色调，负值增加冷色调
	var rGain, bGain float64 = 1.0, 1.0

	if temperature > 0 {
		// 暖色调：增加红色，减少蓝色
		rGain = 1.0 + temperature/100.0*0.3
		bGain = 1.0 - temperature/100.0*0.3
	} else {
		// 冷色调：减少红色，增加蓝色
		rGain = 1.0 + temperature/100.0*0.3
		bGain = 1.0 - temperature/100.0*0.3
	}

	return fmt.Sprintf("colorbalance=rs=%f:bs=%f", (rGain-1.0)*100, (bGain-1.0)*100)
}

// buildShadowHighlightFilter 构建阴影高光调整滤镜
func (f *FFmpeg) buildShadowHighlightFilter(highlights, shadows float64) string {
	if highlights == 0 && shadows == 0 {
		return ""
	}

	// 使用curves滤镜实现阴影高光调整
	var curvePoints []string

	// 阴影调整（影响暗部）
	if shadows != 0 {
		shadowAdjust := shadows / 100.0 * 0.3
		curvePoints = append(curvePoints, fmt.Sprintf("0/%f", 0.0+shadowAdjust))
		curvePoints = append(curvePoints, fmt.Sprintf("0.3/%f", 0.3+shadowAdjust*0.5))
	}

	// 高光调整（影响亮部）
	if highlights != 0 {
		highlightAdjust := highlights / 100.0 * 0.3
		curvePoints = append(curvePoints, fmt.Sprintf("0.7/%f", 0.7+highlightAdjust*0.5))
		curvePoints = append(curvePoints, fmt.Sprintf("1.0/%f", 1.0+highlightAdjust))
	}

	if len(curvePoints) > 0 {
		return fmt.Sprintf("curves=all='%s'", strings.Join(curvePoints, " "))
	}

	return ""
}

// buildLevelsFilter 构建色阶调整滤镜
func (f *FFmpeg) buildLevelsFilter(whites, blacks float64) string {
	if whites == 0 && blacks == 0 {
		return ""
	}

	// 使用curves滤镜实现色阶调整
	blackPoint := blacks / 100.0 * 0.2
	whitePoint := 1.0 + whites/100.0*0.2

	return fmt.Sprintf("curves=all='0/%f 1/%f'", blackPoint, whitePoint)
}

// buildClarityFilter 构建清晰度调整滤镜
func (f *FFmpeg) buildClarityFilter(clarity float64) string {
	if clarity == 0 {
		return ""
	}

	// 清晰度通过unsharp滤镜实现
	amount := 1.0 + clarity/100.0*2.0
	return fmt.Sprintf("unsharp=5:5:%f:5:5:0", amount)
}

// buildVibranceFilter 构建自然饱和度调整滤镜
func (f *FFmpeg) buildVibranceFilter(vibrance float64) string {
	if vibrance == 0 {
		return ""
	}

	// 自然饱和度调整，保护肤色
	saturation := 1.0 + vibrance/100.0*0.5
	return fmt.Sprintf("eq=saturation=%f", saturation)
}

// buildVintageFilter 构建复古滤镜
func (f *FFmpeg) buildVintageFilter(options *FilterOptions) (string, error) {
	params := options.Parameters
	var filters []string

	// 棕褐色调
	if sepia, ok := params["sepia"].(float64); ok && sepia > 0 {
		filters = append(filters, fmt.Sprintf("colorchannelmixer=.393:.769:.189:0:.349:.686:.168:0:.272:.534:.131"))
	}

	// 去饱和度
	if desaturate, ok := params["desaturate"].(float64); ok && desaturate > 0 {
		saturation := 1.0 - desaturate
		filters = append(filters, fmt.Sprintf("eq=saturation=%f", saturation))
	}

	// 胶片颗粒
	if grain, ok := params["grain"].(float64); ok && grain > 0 {
		filters = append(filters, fmt.Sprintf("noise=alls=%d:allf=t", int(grain*100)))
	}

	// 暗角效果
	if vignette, ok := params["vignette"].(float64); ok && vignette > 0 {
		filters = append(filters, fmt.Sprintf("vignette=PI/4*%f", vignette))
	}

	// 褪色效果
	if fade, ok := params["fade"].(float64); ok && fade > 0 {
		contrast := 1.0 - fade*0.3
		brightness := fade * 0.1
		filters = append(filters, fmt.Sprintf("eq=contrast=%f:brightness=%f", contrast, brightness))
	}

	if len(filters) == 0 {
		return "", NewError(ErrInvalidOptions, "复古滤镜参数无效", nil)
	}

	return strings.Join(filters, ","), nil
}

// buildCinematicFilter 构建电影风格滤镜
func (f *FFmpeg) buildCinematicFilter(options *FilterOptions) (string, error) {
	params := options.Parameters
	var filters []string

	// 宽高比调整（添加黑边）
	if aspectRatio, ok := params["aspect_ratio"].(string); ok && aspectRatio != "" {
		if aspectRatio == "21:9" {
			filters = append(filters, "pad=iw:iw*9/21:(ow-iw)/2:(oh-ih)/2:black")
		}
	}

	// 色彩分级预设
	if colorGrading, ok := params["color_grading"].(string); ok {
		switch colorGrading {
		case "teal_orange":
			filters = append(filters, "colorbalance=rs=0.2:gs=-0.1:bs=-0.3:rm=0.1:gm=0:bm=-0.2:rh=0.3:gh=0.1:bh=-0.1")
		case "bleach_bypass":
			filters = append(filters, "eq=contrast=1.3:saturation=0.7:brightness=0.1")
		case "film_noir":
			filters = append(filters, "eq=saturation=0:contrast=1.5:brightness=-0.2")
		}
	}

	// 胶片颗粒
	if filmGrain, ok := params["film_grain"].(float64); ok && filmGrain > 0 {
		filters = append(filters, fmt.Sprintf("noise=alls=%d:allf=t", int(filmGrain*50)))
	}

	// 光晕效果
	if bloom, ok := params["bloom"].(float64); ok && bloom > 0 {
		filters = append(filters, fmt.Sprintf("gblur=sigma=%f", bloom*3))
	}

	if len(filters) == 0 {
		return "", NewError(ErrInvalidOptions, "电影风格滤镜参数无效", nil)
	}

	return strings.Join(filters, ","), nil
}

// buildBeautyFilter 构建美颜滤镜
func (f *FFmpeg) buildBeautyFilter(options *FilterOptions) (string, error) {
	params := options.Parameters
	var filters []string

	// 磨皮效果（使用模糊）
	if skinSmoothing, ok := params["skin_smoothing"].(float64); ok && skinSmoothing > 0 {
		filters = append(filters, fmt.Sprintf("gblur=sigma=%f", skinSmoothing*2))
	}

	// 美白效果
	if skinBrightening, ok := params["skin_brightening"].(float64); ok && skinBrightening > 0 {
		filters = append(filters, fmt.Sprintf("eq=brightness=%f", skinBrightening*0.3))
	}

	// 眼部增强（增加对比度）
	if eyeEnhancement, ok := params["eye_enhancement"].(float64); ok && eyeEnhancement > 0 {
		filters = append(filters, fmt.Sprintf("eq=contrast=%f", 1.0+eyeEnhancement*0.2))
	}

	if len(filters) == 0 {
		return "", NewError(ErrInvalidOptions, "美颜滤镜参数无效", nil)
	}

	return strings.Join(filters, ","), nil
}

// buildSharpeningFilter 构建锐化滤镜
func (f *FFmpeg) buildSharpeningFilter(options *FilterOptions) (string, error) {
	params := options.Parameters
	var filters []string

	amount, _ := params["amount"].(float64)
	radius, _ := params["radius"].(float64)
	threshold, _ := params["threshold"].(float64)
	method, _ := params["method"].(string)

	switch method {
	case "unsharp":
		filters = append(filters, fmt.Sprintf("unsharp=%f:%f:%f:%f:%f:%f",
			radius, radius, amount, radius, radius, threshold))
	case "lanczos":
		filters = append(filters, fmt.Sprintf("scale=iw:ih:flags=lanczos"))
	case "spline":
		filters = append(filters, fmt.Sprintf("scale=iw:ih:flags=spline"))
	default:
		filters = append(filters, fmt.Sprintf("unsharp=%f:%f:%f:%f:%f:%f",
			radius, radius, amount, radius, radius, threshold))
	}

	return strings.Join(filters, ","), nil
}

// buildDenoisingFilter 构建降噪滤镜
func (f *FFmpeg) buildDenoisingFilter(options *FilterOptions) (string, error) {
	params := options.Parameters
	var filters []string

	strength, _ := params["strength"].(float64)
	method, _ := params["method"].(string)
	temporalNR, _ := params["temporal_nr"].(bool)
	spatialNR, _ := params["spatial_nr"].(bool)

	switch method {
	case "nlmeans":
		if spatialNR {
			filters = append(filters, fmt.Sprintf("nlmeans=s=%f", strength*10))
		}
	case "bm3d":
		filters = append(filters, fmt.Sprintf("bm3d=sigma=%f", strength*25))
	case "hqdn3d":
		if temporalNR && spatialNR {
			filters = append(filters, fmt.Sprintf("hqdn3d=%f:%f:%f:%f",
				strength*4, strength*3, strength*6, strength*4.5))
		} else if spatialNR {
			filters = append(filters, fmt.Sprintf("hqdn3d=%f:%f", strength*4, strength*3))
		}
	default:
		filters = append(filters, fmt.Sprintf("nlmeans=s=%f", strength*10))
	}

	return strings.Join(filters, ","), nil
}

// buildVignetteFilter 构建暗角滤镜
func (f *FFmpeg) buildVignetteFilter(options *FilterOptions) (string, error) {
	params := options.Parameters

	intensity, _ := params["intensity"].(float64)
	size, _ := params["size"].(float64)
	softness, _ := params["softness"].(float64)

	if intensity == 0 {
		intensity = 0.5
	}
	if size == 0 {
		size = 0.5
	}
	if softness == 0 {
		softness = 0.5
	}

	vignetteFilter := fmt.Sprintf("vignette=PI/4*%f:mode=forward", intensity)
	return vignetteFilter, nil
}

// buildGlowFilter 构建发光滤镜
func (f *FFmpeg) buildGlowFilter(options *FilterOptions) (string, error) {
	params := options.Parameters
	var filters []string

	intensity, _ := params["intensity"].(float64)
	radius, _ := params["radius"].(float64)
	threshold, _ := params["threshold"].(float64)

	if intensity == 0 {
		intensity = 0.5
	}
	if radius == 0 {
		radius = 10.0
	}
	if threshold == 0 {
		threshold = 0.8
	}

	// 使用高斯模糊创建发光效果
	filters = append(filters, fmt.Sprintf("gblur=sigma=%f", radius))
	filters = append(filters, fmt.Sprintf("eq=brightness=%f", intensity))

	return strings.Join(filters, ","), nil
}

// buildBloomFilter 构建光晕滤镜
func (f *FFmpeg) buildBloomFilter(options *FilterOptions) (string, error) {
	params := options.Parameters
	var filters []string

	intensity, _ := params["intensity"].(float64)
	radius, _ := params["radius"].(float64)
	threshold, _ := params["threshold"].(float64)

	if intensity == 0 {
		intensity = 0.5
	}
	if radius == 0 {
		radius = 20.0
	}
	if threshold == 0 {
		threshold = 0.7
	}

	// 光晕效果通过多层模糊实现
	filters = append(filters, fmt.Sprintf("gblur=sigma=%f", radius*0.5))
	filters = append(filters, fmt.Sprintf("gblur=sigma=%f", radius))
	filters = append(filters, fmt.Sprintf("eq=brightness=%f", intensity*0.3))

	return strings.Join(filters, ","), nil
}
