// Package ffmpeg 提供性能优化功能
package ffmpeg

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// PerformanceOptions 性能优化选项
type PerformanceOptions struct {
	ThreadCount       int     // 线程数量 (0为自动检测)
	HardwareAccel     string  // 硬件加速 ("auto", "cuda", "opencl", "vaapi", "videotoolbox", "none")
	MemoryLimit       int64   // 内存限制 (MB)
	TempDirectory     string  // 临时文件目录
	EnableGPU         bool    // 是否启用GPU加速
	Quality           string  // 质量设置 ("fast", "medium", "slow", "veryslow")
	Preset            string  // 编码预设 ("ultrafast", "superfast", "veryfast", "faster", "fast", "medium", "slow", "slower", "veryslow")
	OptimizeFor       string  // 优化目标 ("speed", "quality", "size")
	EnableMultipass   bool    // 是否启用多遍编码
	BufferSize        int     // 缓冲区大小 (KB)
}

// OptimizeForSpeed 优化编码速度
func (f *FFmpeg) OptimizeForSpeed() *PerformanceOptions {
	return &PerformanceOptions{
		ThreadCount:       runtime.NumCPU(),
		HardwareAccel:     "auto",
		EnableGPU:         true,
		Quality:           "fast",
		Preset:            "ultrafast",
		OptimizeFor:       "speed",
		EnableMultipass:   false,
		BufferSize:        1024,
	}
}

// OptimizeForQuality 优化编码质量
func (f *FFmpeg) OptimizeForQuality() *PerformanceOptions {
	return &PerformanceOptions{
		ThreadCount:       runtime.NumCPU(),
		HardwareAccel:     "auto",
		EnableGPU:         true,
		Quality:           "slow",
		Preset:            "slow",
		OptimizeFor:       "quality",
		EnableMultipass:   true,
		BufferSize:        4096,
	}
}

// OptimizeForSize 优化文件大小
func (f *FFmpeg) OptimizeForSize() *PerformanceOptions {
	return &PerformanceOptions{
		ThreadCount:       runtime.NumCPU() / 2,
		HardwareAccel:     "auto",
		EnableGPU:         false,
		Quality:           "veryslow",
		Preset:            "veryslow",
		OptimizeFor:       "size",
		EnableMultipass:   true,
		BufferSize:        2048,
	}
}

// OptimizeForBalance 平衡优化
func (f *FFmpeg) OptimizeForBalance() *PerformanceOptions {
	return &PerformanceOptions{
		ThreadCount:       runtime.NumCPU(),
		HardwareAccel:     "auto",
		EnableGPU:         true,
		Quality:           "medium",
		Preset:            "medium",
		OptimizeFor:       "balance",
		EnableMultipass:   false,
		BufferSize:        2048,
	}
}

// ApplyPerformanceOptions 应用性能优化选项到FFmpeg参数
func (f *FFmpeg) ApplyPerformanceOptions(args []string, options *PerformanceOptions) []string {
	if options == nil {
		return args
	}

	var optimizedArgs []string

	// 硬件加速设置
	if options.HardwareAccel != "" && options.HardwareAccel != "none" {
		if options.HardwareAccel == "auto" {
			// 自动检测硬件加速
			hwAccel := f.detectHardwareAcceleration()
			if hwAccel != "" {
				optimizedArgs = append(optimizedArgs, "-hwaccel", hwAccel)
			}
		} else {
			optimizedArgs = append(optimizedArgs, "-hwaccel", options.HardwareAccel)
		}
	}

	// 线程数设置
	if options.ThreadCount > 0 {
		optimizedArgs = append(optimizedArgs, "-threads", fmt.Sprintf("%d", options.ThreadCount))
	}

	// 添加原始参数
	optimizedArgs = append(optimizedArgs, args...)

	// 编码预设
	if options.Preset != "" {
		optimizedArgs = append(optimizedArgs, "-preset", options.Preset)
	}

	// 质量设置
	if options.Quality != "" {
		switch options.Quality {
		case "fast":
			optimizedArgs = append(optimizedArgs, "-crf", "28")
		case "medium":
			optimizedArgs = append(optimizedArgs, "-crf", "23")
		case "slow":
			optimizedArgs = append(optimizedArgs, "-crf", "18")
		case "veryslow":
			optimizedArgs = append(optimizedArgs, "-crf", "15")
		}
	}

	// 多遍编码
	if options.EnableMultipass {
		optimizedArgs = append(optimizedArgs, "-pass", "1")
	}

	// 缓冲区大小
	if options.BufferSize > 0 {
		optimizedArgs = append(optimizedArgs, "-bufsize", fmt.Sprintf("%dk", options.BufferSize))
	}

	// 内存限制
	if options.MemoryLimit > 0 {
		// FFmpeg没有直接的内存限制参数，但可以通过其他方式优化
		optimizedArgs = append(optimizedArgs, "-max_muxing_queue_size", "1024")
	}

	return optimizedArgs
}

// detectHardwareAcceleration 检测可用的硬件加速
func (f *FFmpeg) detectHardwareAcceleration() string {
	// 检测NVIDIA CUDA
	if f.isHardwareAccelAvailable("cuda") {
		return "cuda"
	}

	// 检测Intel Quick Sync Video
	if f.isHardwareAccelAvailable("qsv") {
		return "qsv"
	}

	// 检测AMD AMF
	if f.isHardwareAccelAvailable("amf") {
		return "amf"
	}

	// 检测Apple VideoToolbox (macOS)
	if runtime.GOOS == "darwin" && f.isHardwareAccelAvailable("videotoolbox") {
		return "videotoolbox"
	}

	// 检测VAAPI (Linux)
	if runtime.GOOS == "linux" && f.isHardwareAccelAvailable("vaapi") {
		return "vaapi"
	}

	// 检测OpenCL
	if f.isHardwareAccelAvailable("opencl") {
		return "opencl"
	}

	return ""
}

// isHardwareAccelAvailable 检查特定硬件加速是否可用
func (f *FFmpeg) isHardwareAccelAvailable(accel string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []string{"-hwaccels"}
	output, err := f.executeCommand(ctx, args)
	if err != nil {
		return false
	}

	return strings.Contains(string(output), accel)
}

// GetOptimalThreadCount 获取最优线程数
func (f *FFmpeg) GetOptimalThreadCount() int {
	cpuCount := runtime.NumCPU()
	
	// 对于视频编码，通常使用CPU核心数
	// 但不超过16个线程，因为收益递减
	if cpuCount > 16 {
		return 16
	}
	
	return cpuCount
}

// EstimateProcessingTime 估算处理时间
func (f *FFmpeg) EstimateProcessingTime(inputDuration time.Duration, options *PerformanceOptions) time.Duration {
	if options == nil {
		options = f.OptimizeForBalance()
	}

	// 基础处理倍数
	var multiplier float64 = 1.0

	// 根据质量设置调整
	switch options.Quality {
	case "fast":
		multiplier = 0.3
	case "medium":
		multiplier = 0.8
	case "slow":
		multiplier = 2.0
	case "veryslow":
		multiplier = 4.0
	}

	// 硬件加速可以显著提升速度
	if options.EnableGPU && options.HardwareAccel != "none" {
		multiplier *= 0.3
	}

	// 多遍编码会增加时间
	if options.EnableMultipass {
		multiplier *= 2.0
	}

	// 线程数影响
	threadEfficiency := float64(options.ThreadCount) / float64(runtime.NumCPU())
	if threadEfficiency > 1.0 {
		threadEfficiency = 1.0
	}
	multiplier /= threadEfficiency

	estimatedTime := time.Duration(float64(inputDuration) * multiplier)
	
	// 最小处理时间为5秒
	if estimatedTime < 5*time.Second {
		estimatedTime = 5 * time.Second
	}

	return estimatedTime
}

// OptimizeCommand 优化FFmpeg命令
func (f *FFmpeg) OptimizeCommand(args []string, options *PerformanceOptions) []string {
	if options == nil {
		options = f.OptimizeForBalance()
	}

	optimizedArgs := f.ApplyPerformanceOptions(args, options)

	// 添加通用优化参数
	optimizedArgs = append(optimizedArgs, 
		"-movflags", "+faststart", // 优化MP4文件结构
		"-fflags", "+genpts",      // 生成PTS
		"-avoid_negative_ts", "make_zero", // 避免负时间戳
	)

	return optimizedArgs
}

// MonitorPerformance 监控性能指标
type PerformanceMonitor struct {
	StartTime    time.Time
	LastUpdate   time.Time
	ProcessedFrames int64
	TotalFrames     int64
	CurrentFPS      float64
	AverageFPS      float64
	EstimatedTimeRemaining time.Duration
}

// NewPerformanceMonitor 创建性能监控器
func (f *FFmpeg) NewPerformanceMonitor(totalFrames int64) *PerformanceMonitor {
	now := time.Now()
	return &PerformanceMonitor{
		StartTime:   now,
		LastUpdate:  now,
		TotalFrames: totalFrames,
	}
}

// Update 更新性能指标
func (pm *PerformanceMonitor) Update(processedFrames int64) {
	now := time.Now()
	elapsed := now.Sub(pm.StartTime)
	
	pm.ProcessedFrames = processedFrames
	pm.LastUpdate = now
	
	if elapsed.Seconds() > 0 {
		pm.AverageFPS = float64(processedFrames) / elapsed.Seconds()
	}
	
	// 计算剩余时间
	if pm.AverageFPS > 0 && pm.TotalFrames > 0 {
		remainingFrames := pm.TotalFrames - processedFrames
		pm.EstimatedTimeRemaining = time.Duration(float64(remainingFrames) / pm.AverageFPS * float64(time.Second))
	}
}

// GetProgress 获取进度百分比
func (pm *PerformanceMonitor) GetProgress() float64 {
	if pm.TotalFrames == 0 {
		return 0
	}
	return float64(pm.ProcessedFrames) / float64(pm.TotalFrames) * 100
}

// GetPerformanceReport 获取性能报告
func (pm *PerformanceMonitor) GetPerformanceReport() string {
	elapsed := time.Since(pm.StartTime)
	progress := pm.GetProgress()
	
	return fmt.Sprintf(
		"进度: %.1f%% | 已处理: %d/%d 帧 | 平均FPS: %.1f | 已用时间: %v | 预计剩余: %v",
		progress,
		pm.ProcessedFrames,
		pm.TotalFrames,
		pm.AverageFPS,
		elapsed.Round(time.Second),
		pm.EstimatedTimeRemaining.Round(time.Second),
	)
}
