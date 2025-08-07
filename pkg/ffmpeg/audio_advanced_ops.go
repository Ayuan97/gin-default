// Package ffmpeg 提供高级音频处理功能
package ffmpeg

import (
	"context"
	"fmt"
	"strings"
)

// ApplyAudioEqualizer 应用音频均衡器
// inputPath: 输入音频/视频文件路径
// outputPath: 输出文件路径
// options: 均衡器选项
func (f *FFmpeg) ApplyAudioEqualizer(inputPath, outputPath string, options *AudioEqualizerOptions) error {
	return f.ApplyAudioEqualizerWithContext(context.Background(), inputPath, outputPath, options)
}

// ApplyAudioEqualizerWithContext 带上下文的音频均衡器
func (f *FFmpeg) ApplyAudioEqualizerWithContext(ctx context.Context, inputPath, outputPath string, options *AudioEqualizerOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	if options == nil {
		return NewError(ErrInvalidOptions, "音频均衡器选项不能为空", nil)
	}

	// 构建均衡器滤镜
	eqFilter, err := f.buildEqualizerFilter(options)
	if err != nil {
		return err
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath,
		"-af", eqFilter,
		"-c:v", "copy",
		"-y",
		outputPath,
	}

	// 执行命令
	_, err = f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, "音频均衡器处理失败", err)
	}

	f.logger.Info("音频均衡器处理完成: %s -> %s", inputPath, outputPath)
	return nil
}

// ApplyAudioEffect 应用音频效果
// inputPath: 输入音频/视频文件路径
// outputPath: 输出文件路径
// options: 音频效果选项
func (f *FFmpeg) ApplyAudioEffect(inputPath, outputPath string, options *AudioEffectOptions) error {
	return f.ApplyAudioEffectWithContext(context.Background(), inputPath, outputPath, options)
}

// ApplyAudioEffectWithContext 带上下文的音频效果
func (f *FFmpeg) ApplyAudioEffectWithContext(ctx context.Context, inputPath, outputPath string, options *AudioEffectOptions) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	if options == nil {
		return NewError(ErrInvalidOptions, "音频效果选项不能为空", nil)
	}

	// 构建音频效果滤镜
	effectFilter, err := f.buildAudioEffectFilter(options)
	if err != nil {
		return err
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath,
		"-af", effectFilter,
		"-c:v", "copy",
		"-y",
		outputPath,
	}

	// 执行命令
	_, err = f.executeCommand(ctx, args)
	if err != nil {
		return NewError(ErrExecutionFailed, fmt.Sprintf("音频效果处理失败: %s", options.EffectType), err)
	}

	f.logger.Info("音频效果处理完成: %s (%s) -> %s", options.EffectType, inputPath, outputPath)
	return nil
}

// ApplyReverb 应用混响效果
func (f *FFmpeg) ApplyReverb(inputPath, outputPath string, options *ReverbOptions) error {
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

	return f.ApplyAudioEffect(inputPath, outputPath, effectOptions)
}

// ApplyCompressor 应用音频压缩器
func (f *FFmpeg) ApplyCompressor(inputPath, outputPath string, options *CompressorOptions) error {
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
		Intensity:  options.Ratio / 20.0, // 归一化到0-1
		Parameters: map[string]interface{}{
			"threshold":   options.Threshold,
			"ratio":       options.Ratio,
			"attack":      options.Attack,
			"release":     options.Release,
			"makeup_gain": options.MakeupGain,
			"knee_width":  options.KneeWidth,
		},
	}

	return f.ApplyAudioEffect(inputPath, outputPath, effectOptions)
}

// ApplyEcho 应用回声效果
func (f *FFmpeg) ApplyEcho(inputPath, outputPath string, delay float64, decay float64) error {
	effectOptions := &AudioEffectOptions{
		EffectType: AudioEffectEcho,
		Intensity:  decay,
		Parameters: map[string]interface{}{
			"delay": delay,
			"decay": decay,
		},
	}

	return f.ApplyAudioEffect(inputPath, outputPath, effectOptions)
}

// ApplyChorus 应用合唱效果
func (f *FFmpeg) ApplyChorus(inputPath, outputPath string, intensity float64) error {
	effectOptions := &AudioEffectOptions{
		EffectType: AudioEffectChorus,
		Intensity:  intensity,
		Parameters: map[string]interface{}{
			"in_gain":  0.4,
			"out_gain": 0.4,
			"delays":   "25,40,60",
			"decays":   "0.5,0.3,0.2",
			"speeds":   "0.25,0.4,0.3",
			"depths":   "2,2.3,1.3",
		},
	}

	return f.ApplyAudioEffect(inputPath, outputPath, effectOptions)
}

// ApplyFlanger 应用镶边效果
func (f *FFmpeg) ApplyFlanger(inputPath, outputPath string, intensity float64) error {
	effectOptions := &AudioEffectOptions{
		EffectType: AudioEffectFlanger,
		Intensity:  intensity,
		Parameters: map[string]interface{}{
			"delay":  0,
			"depth":  2,
			"regen":  0,
			"width":  71,
			"speed":  0.5,
			"shape":  "sine",
			"phase":  25,
			"interp": "linear",
		},
	}

	return f.ApplyAudioEffect(inputPath, outputPath, effectOptions)
}

// ApplyPhaser 应用相位器效果
func (f *FFmpeg) ApplyPhaser(inputPath, outputPath string, intensity float64) error {
	effectOptions := &AudioEffectOptions{
		EffectType: AudioEffectPhaser,
		Intensity:  intensity,
		Parameters: map[string]interface{}{
			"in_gain":  0.4,
			"out_gain": 0.74,
			"delay":    3.0,
			"decay":    0.4,
			"speed":    0.5,
			"type":     "triangular",
		},
	}

	return f.ApplyAudioEffect(inputPath, outputPath, effectOptions)
}

// ApplyDistortion 应用失真效果
func (f *FFmpeg) ApplyDistortion(inputPath, outputPath string, intensity float64) error {
	effectOptions := &AudioEffectOptions{
		EffectType: AudioEffectDistortion,
		Intensity:  intensity,
		Parameters: map[string]interface{}{
			"gain":   20,
			"colour": 20,
		},
	}

	return f.ApplyAudioEffect(inputPath, outputPath, effectOptions)
}

// ApplyLimiter 应用限制器
func (f *FFmpeg) ApplyLimiter(inputPath, outputPath string, limit float64) error {
	effectOptions := &AudioEffectOptions{
		EffectType: AudioEffectLimiter,
		Intensity:  1.0,
		Parameters: map[string]interface{}{
			"limit":   limit,
			"attack":  5,
			"release": 50,
		},
	}

	return f.ApplyAudioEffect(inputPath, outputPath, effectOptions)
}

// ApplyNoiseGate 应用噪声门
func (f *FFmpeg) ApplyNoiseGate(inputPath, outputPath string, threshold float64) error {
	effectOptions := &AudioEffectOptions{
		EffectType: AudioEffectGate,
		Intensity:  1.0,
		Parameters: map[string]interface{}{
			"threshold": threshold,
			"ratio":     2,
			"attack":    20,
			"release":   250,
		},
	}

	return f.ApplyAudioEffect(inputPath, outputPath, effectOptions)
}

// ApplyPitchShift 应用变调效果
func (f *FFmpeg) ApplyPitchShift(inputPath, outputPath string, semitones float64) error {
	effectOptions := &AudioEffectOptions{
		EffectType: AudioEffectPitchShift,
		Intensity:  1.0,
		Parameters: map[string]interface{}{
			"pitch": semitones,
		},
	}

	return f.ApplyAudioEffect(inputPath, outputPath, effectOptions)
}

// ApplyAudioNormalization 应用音频标准化
func (f *FFmpeg) ApplyAudioNormalization(inputPath, outputPath string, targetLevel float64) error {
	// 验证输入文件
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// 验证输出文件路径
	if err := validateOutputFile(outputPath); err != nil {
		return err
	}

	// 构建FFmpeg命令参数
	args := []string{
		"-i", inputPath,
		"-af", fmt.Sprintf("loudnorm=I=%f:TP=-1.5:LRA=11", targetLevel),
		"-c:v", "copy",
		"-y",
		outputPath,
	}

	// 执行命令
	_, err := f.executeCommand(context.Background(), args)
	if err != nil {
		return NewError(ErrExecutionFailed, "音频标准化处理失败", err)
	}

	f.logger.Info("音频标准化处理完成: %s -> %s (目标电平: %.1f LUFS)", inputPath, outputPath, targetLevel)
	return nil
}

// === 音频滤镜构建辅助方法 ===

// buildEqualizerFilter 构建均衡器滤镜
func (f *FFmpeg) buildEqualizerFilter(options *AudioEqualizerOptions) (string, error) {
	var filters []string

	// 使用预设
	if options.Preset != "" {
		presetFilter := f.getEqualizerPreset(options.Preset)
		if presetFilter != "" {
			filters = append(filters, presetFilter)
		}
	}

	// 自定义频段
	if len(options.Bands) > 0 {
		for _, band := range options.Bands {
			eqFilter := fmt.Sprintf("equalizer=f=%f:width_type=q:width=%f:g=%f",
				band.Frequency, band.Q, band.Gain)
			filters = append(filters, eqFilter)
		}
	}

	// 主增益
	if options.MasterGain != 0 {
		filters = append(filters, fmt.Sprintf("volume=%fdB", options.MasterGain))
	}

	if len(filters) == 0 {
		return "", NewError(ErrInvalidOptions, "均衡器参数无效", nil)
	}

	return strings.Join(filters, ","), nil
}

// getEqualizerPreset 获取均衡器预设
func (f *FFmpeg) getEqualizerPreset(preset string) string {
	presets := map[string]string{
		"rock":      "equalizer=f=60:width_type=q:width=1:g=4,equalizer=f=170:width_type=q:width=1:g=3,equalizer=f=310:width_type=q:width=1:g=-2,equalizer=f=600:width_type=q:width=1:g=-1,equalizer=f=1000:width_type=q:width=1:g=2,equalizer=f=3000:width_type=q:width=1:g=3,equalizer=f=6000:width_type=q:width=1:g=4,equalizer=f=12000:width_type=q:width=1:g=5,equalizer=f=14000:width_type=q:width=1:g=5",
		"pop":       "equalizer=f=60:width_type=q:width=1:g=2,equalizer=f=170:width_type=q:width=1:g=1,equalizer=f=310:width_type=q:width=1:g=-1,equalizer=f=600:width_type=q:width=1:g=-2,equalizer=f=1000:width_type=q:width=1:g=-1,equalizer=f=3000:width_type=q:width=1:g=1,equalizer=f=6000:width_type=q:width=1:g=2,equalizer=f=12000:width_type=q:width=1:g=3,equalizer=f=14000:width_type=q:width=1:g=3",
		"classical": "equalizer=f=60:width_type=q:width=1:g=0,equalizer=f=170:width_type=q:width=1:g=0,equalizer=f=310:width_type=q:width=1:g=0,equalizer=f=600:width_type=q:width=1:g=0,equalizer=f=1000:width_type=q:width=1:g=0,equalizer=f=3000:width_type=q:width=1:g=0,equalizer=f=6000:width_type=q:width=1:g=-2,equalizer=f=12000:width_type=q:width=1:g=-2,equalizer=f=14000:width_type=q:width=1:g=-2",
		"jazz":      "equalizer=f=60:width_type=q:width=1:g=3,equalizer=f=170:width_type=q:width=1:g=2,equalizer=f=310:width_type=q:width=1:g=1,equalizer=f=600:width_type=q:width=1:g=2,equalizer=f=1000:width_type=q:width=1:g=-2,equalizer=f=3000:width_type=q:width=1:g=-2,equalizer=f=6000:width_type=q:width=1:g=0,equalizer=f=12000:width_type=q:width=1:g=1,equalizer=f=14000:width_type=q:width=1:g=2",
		"vocal":     "equalizer=f=60:width_type=q:width=1:g=-2,equalizer=f=170:width_type=q:width=1:g=-1,equalizer=f=310:width_type=q:width=1:g=2,equalizer=f=600:width_type=q:width=1:g=3,equalizer=f=1000:width_type=q:width=1:g=3,equalizer=f=3000:width_type=q:width=1:g=2,equalizer=f=6000:width_type=q:width=1:g=1,equalizer=f=12000:width_type=q:width=1:g=0,equalizer=f=14000:width_type=q:width=1:g=-1",
	}

	return presets[preset]
}

// buildAudioEffectFilter 构建音频效果滤镜
func (f *FFmpeg) buildAudioEffectFilter(options *AudioEffectOptions) (string, error) {
	switch options.EffectType {
	case AudioEffectReverb:
		return f.buildReverbFilter(options)
	case AudioEffectEcho:
		return f.buildEchoFilter(options)
	case AudioEffectChorus:
		return f.buildChorusFilter(options)
	case AudioEffectFlanger:
		return f.buildFlangerFilter(options)
	case AudioEffectPhaser:
		return f.buildPhaserFilter(options)
	case AudioEffectDistortion:
		return f.buildDistortionFilter(options)
	case AudioEffectCompressor:
		return f.buildCompressorFilter(options)
	case AudioEffectLimiter:
		return f.buildLimiterFilter(options)
	case AudioEffectGate:
		return f.buildGateFilter(options)
	case AudioEffectPitchShift:
		return f.buildPitchShiftFilter(options)
	default:
		return "", NewError(ErrInvalidOptions, fmt.Sprintf("不支持的音频效果类型: %s", options.EffectType), nil)
	}
}

// buildReverbFilter 构建混响滤镜
func (f *FFmpeg) buildReverbFilter(options *AudioEffectOptions) (string, error) {
	params := options.Parameters

	roomSize, _ := params["room_size"].(float64)
	damping, _ := params["damping"].(float64)
	wetLevel, _ := params["wet_level"].(float64)
	dryLevel, _ := params["dry_level"].(float64)

	// 使用freeverb滤镜
	filter := fmt.Sprintf("afreeverb=room_size=%f:damping=%f:wet_gain=%f:dry_gain=%f",
		roomSize, damping, wetLevel, dryLevel)

	return filter, nil
}

// buildEchoFilter 构建回声滤镜
func (f *FFmpeg) buildEchoFilter(options *AudioEffectOptions) (string, error) {
	params := options.Parameters

	delay, _ := params["delay"].(float64)
	decay, _ := params["decay"].(float64)

	// 使用aecho滤镜
	filter := fmt.Sprintf("aecho=0.8:0.9:%f:%f", delay*1000, decay)

	return filter, nil
}

// buildChorusFilter 构建合唱滤镜
func (f *FFmpeg) buildChorusFilter(options *AudioEffectOptions) (string, error) {
	params := options.Parameters

	inGain, _ := params["in_gain"].(float64)
	outGain, _ := params["out_gain"].(float64)
	delays, _ := params["delays"].(string)
	decays, _ := params["decays"].(string)
	speeds, _ := params["speeds"].(string)
	depths, _ := params["depths"].(string)

	// 使用achorus滤镜
	filter := fmt.Sprintf("achorus=%f:%f:%s:%s:%s:%s",
		inGain, outGain, delays, decays, speeds, depths)

	return filter, nil
}

// buildFlangerFilter 构建镶边滤镜
func (f *FFmpeg) buildFlangerFilter(options *AudioEffectOptions) (string, error) {
	params := options.Parameters

	delay, _ := params["delay"].(float64)
	depth, _ := params["depth"].(float64)
	regen, _ := params["regen"].(float64)
	width, _ := params["width"].(float64)
	speed, _ := params["speed"].(float64)

	// 使用aflanger滤镜
	filter := fmt.Sprintf("aflanger=delay=%f:depth=%f:regen=%f:width=%f:speed=%f",
		delay, depth, regen, width, speed)

	return filter, nil
}

// buildPhaserFilter 构建相位器滤镜
func (f *FFmpeg) buildPhaserFilter(options *AudioEffectOptions) (string, error) {
	params := options.Parameters

	inGain, _ := params["in_gain"].(float64)
	outGain, _ := params["out_gain"].(float64)
	delay, _ := params["delay"].(float64)
	decay, _ := params["decay"].(float64)
	speed, _ := params["speed"].(float64)

	// 使用aphaser滤镜
	filter := fmt.Sprintf("aphaser=in_gain=%f:out_gain=%f:delay=%f:decay=%f:speed=%f",
		inGain, outGain, delay, decay, speed)

	return filter, nil
}

// buildDistortionFilter 构建失真滤镜
func (f *FFmpeg) buildDistortionFilter(options *AudioEffectOptions) (string, error) {
	params := options.Parameters

	gain, _ := params["gain"].(float64)
	colour, _ := params["colour"].(float64)

	// 使用overdrive滤镜
	filter := fmt.Sprintf("overdrive=gain=%f:colour=%f", gain, colour)

	return filter, nil
}

// buildCompressorFilter 构建压缩器滤镜
func (f *FFmpeg) buildCompressorFilter(options *AudioEffectOptions) (string, error) {
	params := options.Parameters

	threshold, _ := params["threshold"].(float64)
	ratio, _ := params["ratio"].(float64)
	attack, _ := params["attack"].(float64)
	release, _ := params["release"].(float64)
	makeupGain, _ := params["makeup_gain"].(float64)

	// 使用acompressor滤镜
	filter := fmt.Sprintf("acompressor=threshold=%fdB:ratio=%f:attack=%f:release=%f:makeup=%fdB",
		threshold, ratio, attack, release, makeupGain)

	return filter, nil
}

// buildLimiterFilter 构建限制器滤镜
func (f *FFmpeg) buildLimiterFilter(options *AudioEffectOptions) (string, error) {
	params := options.Parameters

	limit, _ := params["limit"].(float64)
	attack, _ := params["attack"].(float64)
	release, _ := params["release"].(float64)

	// 使用alimiter滤镜
	filter := fmt.Sprintf("alimiter=limit=%f:attack=%f:release=%f", limit, attack, release)

	return filter, nil
}

// buildGateFilter 构建噪声门滤镜
func (f *FFmpeg) buildGateFilter(options *AudioEffectOptions) (string, error) {
	params := options.Parameters

	threshold, _ := params["threshold"].(float64)
	ratio, _ := params["ratio"].(float64)
	attack, _ := params["attack"].(float64)
	release, _ := params["release"].(float64)

	// 使用agate滤镜
	filter := fmt.Sprintf("agate=threshold=%fdB:ratio=%f:attack=%f:release=%f",
		threshold, ratio, attack, release)

	return filter, nil
}

// buildPitchShiftFilter 构建变调滤镜
func (f *FFmpeg) buildPitchShiftFilter(options *AudioEffectOptions) (string, error) {
	params := options.Parameters

	pitch, _ := params["pitch"].(float64)

	// 计算变调比例 (每个半音约为1.059463倍)
	ratio := 1.0
	if pitch != 0 {
		ratio = 1.059463 * pitch
	}

	// 使用asetrate和atempo实现变调
	filter := fmt.Sprintf("asetrate=44100*%f,aresample=44100,atempo=%f", ratio, 1.0/ratio)

	return filter, nil
}
