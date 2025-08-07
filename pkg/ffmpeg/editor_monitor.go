// Package ffmpeg 提供进度监控功能
package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ProgressMonitor 进度监控器
type ProgressMonitor struct {
	callback    ProgressCallback // 进度回调函数
	totalTime   time.Duration    // 总时长
	currentTime time.Duration    // 当前时间
	mu          sync.RWMutex     // 读写锁
	cancelled   bool             // 是否已取消
}

// NewProgressMonitor 创建新的进度监控器
func NewProgressMonitor(callback ProgressCallback, totalTime time.Duration) *ProgressMonitor {
	return &ProgressMonitor{
		callback:  callback,
		totalTime: totalTime,
	}
}

// Start 开始监控进度
func (pm *ProgressMonitor) Start(ctx context.Context, cmd *exec.Cmd) error {
	// 创建管道来读取FFmpeg的stderr输出
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("创建stderr管道失败: %w", err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动FFmpeg命令失败: %w", err)
	}

	// 在goroutine中监控进度
	go pm.monitorProgress(stderr)

	// 等待命令完成或上下文取消
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// 上下文取消，终止进程
		pm.Cancel()
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// monitorProgress 监控FFmpeg输出并解析进度
func (pm *ProgressMonitor) monitorProgress(stderr io.Reader) {
	scanner := bufio.NewScanner(stderr)
	
	// 用于匹配FFmpeg进度输出的正则表达式
	timeRegex := regexp.MustCompile(`time=(\d{2}):(\d{2}):(\d{2})\.(\d{2})`)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// 检查是否已取消
		pm.mu.RLock()
		cancelled := pm.cancelled
		pm.mu.RUnlock()
		
		if cancelled {
			break
		}
		
		// 解析时间信息
		if matches := timeRegex.FindStringSubmatch(line); len(matches) == 5 {
			currentTime := pm.parseTimeFromMatches(matches)
			pm.updateProgress(currentTime)
		}
	}
}

// parseTimeFromMatches 从正则匹配结果解析时间
func (pm *ProgressMonitor) parseTimeFromMatches(matches []string) time.Duration {
	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.Atoi(matches[3])
	centiseconds, _ := strconv.Atoi(matches[4])
	
	totalSeconds := float64(hours*3600 + minutes*60 + seconds) + float64(centiseconds)/100.0
	return time.Duration(totalSeconds * float64(time.Second))
}

// updateProgress 更新进度
func (pm *ProgressMonitor) updateProgress(currentTime time.Duration) {
	pm.mu.Lock()
	pm.currentTime = currentTime
	callback := pm.callback
	totalTime := pm.totalTime
	pm.mu.Unlock()
	
	if callback != nil && totalTime > 0 {
		progress := float64(currentTime) / float64(totalTime) * 100
		if progress > 100 {
			progress = 100
		}
		callback(progress, currentTime, totalTime)
	}
}

// Cancel 取消监控
func (pm *ProgressMonitor) Cancel() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.cancelled = true
}

// GetProgress 获取当前进度
func (pm *ProgressMonitor) GetProgress() (float64, time.Duration, time.Duration) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	if pm.totalTime == 0 {
		return 0, pm.currentTime, pm.totalTime
	}
	
	progress := float64(pm.currentTime) / float64(pm.totalTime) * 100
	if progress > 100 {
		progress = 100
	}
	
	return progress, pm.currentTime, pm.totalTime
}

// executeCommandWithProgress 执行带进度监控的FFmpeg命令
func (f *FFmpeg) executeCommandWithProgress(ctx context.Context, args []string, callback ProgressCallback, totalTime time.Duration) error {
	f.mu.RLock()
	execPath := f.execPath
	logger := f.logger
	f.mu.RUnlock()

	logger.Debug("执行带进度监控的FFmpeg命令: %s %s", execPath, strings.Join(args, " "))

	// 创建命令
	cmd := exec.CommandContext(ctx, execPath, args...)
	
	// 如果有回调函数，启用进度监控
	if callback != nil {
		monitor := NewProgressMonitor(callback, totalTime)
		return monitor.Start(ctx, cmd)
	}
	
	// 没有回调函数，直接执行
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("FFmpeg命令执行失败: %s, 输出: %s", err.Error(), string(output))
		return fmt.Errorf("FFmpeg命令执行失败: %w", err)
	}
	
	logger.Debug("FFmpeg命令执行成功")
	return nil
}

// ProgressTracker 进度跟踪器，用于多步骤操作
type ProgressTracker struct {
	steps       []ProgressStep   // 步骤列表
	currentStep int              // 当前步骤
	callback    ProgressCallback // 总体进度回调
	mu          sync.RWMutex     // 读写锁
}

// ProgressStep 进度步骤
type ProgressStep struct {
	Name        string        // 步骤名称
	Weight      float64       // 权重（占总进度的比例）
	Duration    time.Duration // 预估时长
	Progress    float64       // 当前进度 (0-100)
	CurrentTime time.Duration // 当前时间
}

// NewProgressTracker 创建新的进度跟踪器
func NewProgressTracker(callback ProgressCallback) *ProgressTracker {
	return &ProgressTracker{
		steps:    make([]ProgressStep, 0),
		callback: callback,
	}
}

// AddStep 添加步骤
func (pt *ProgressTracker) AddStep(name string, weight float64, duration time.Duration) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	
	step := ProgressStep{
		Name:     name,
		Weight:   weight,
		Duration: duration,
		Progress: 0,
	}
	pt.steps = append(pt.steps, step)
}

// StartStep 开始执行步骤
func (pt *ProgressTracker) StartStep(stepIndex int) ProgressCallback {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	
	if stepIndex < 0 || stepIndex >= len(pt.steps) {
		return nil
	}
	
	pt.currentStep = stepIndex
	
	// 返回该步骤的进度回调函数
	return func(progress float64, currentTime time.Duration, totalTime time.Duration) {
		pt.updateStepProgress(stepIndex, progress, currentTime)
	}
}

// updateStepProgress 更新步骤进度
func (pt *ProgressTracker) updateStepProgress(stepIndex int, progress float64, currentTime time.Duration) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	
	if stepIndex < 0 || stepIndex >= len(pt.steps) {
		return
	}
	
	pt.steps[stepIndex].Progress = progress
	pt.steps[stepIndex].CurrentTime = currentTime
	
	// 计算总体进度
	totalProgress := pt.calculateTotalProgress()
	totalCurrentTime := pt.calculateTotalCurrentTime()
	totalDuration := pt.calculateTotalDuration()
	
	// 调用总体进度回调
	if pt.callback != nil {
		pt.callback(totalProgress, totalCurrentTime, totalDuration)
	}
}

// calculateTotalProgress 计算总体进度
func (pt *ProgressTracker) calculateTotalProgress() float64 {
	var totalProgress float64
	var totalWeight float64
	
	for _, step := range pt.steps {
		totalProgress += step.Progress * step.Weight
		totalWeight += step.Weight
	}
	
	if totalWeight == 0 {
		return 0
	}
	
	return totalProgress / totalWeight
}

// calculateTotalCurrentTime 计算总体当前时间
func (pt *ProgressTracker) calculateTotalCurrentTime() time.Duration {
	var totalTime time.Duration
	
	for i, step := range pt.steps {
		if i < pt.currentStep {
			// 已完成的步骤
			totalTime += step.Duration
		} else if i == pt.currentStep {
			// 当前步骤
			totalTime += step.CurrentTime
		}
		// 未开始的步骤不计算
	}
	
	return totalTime
}

// calculateTotalDuration 计算总体预估时长
func (pt *ProgressTracker) calculateTotalDuration() time.Duration {
	var totalDuration time.Duration
	
	for _, step := range pt.steps {
		totalDuration += step.Duration
	}
	
	return totalDuration
}

// GetCurrentStep 获取当前步骤信息
func (pt *ProgressTracker) GetCurrentStep() (int, ProgressStep) {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	
	if pt.currentStep < 0 || pt.currentStep >= len(pt.steps) {
		return -1, ProgressStep{}
	}
	
	return pt.currentStep, pt.steps[pt.currentStep]
}

// GetAllSteps 获取所有步骤信息
func (pt *ProgressTracker) GetAllSteps() []ProgressStep {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	
	steps := make([]ProgressStep, len(pt.steps))
	copy(steps, pt.steps)
	return steps
}
