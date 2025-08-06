// Package ffmpeg 提供了对FFmpeg命令行工具的Go语言封装
// 支持视频格式转换、压缩、音频提取、视频信息获取和基本编辑操作
package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// FFmpeg 主要的FFmpeg操作接口
type FFmpeg struct {
	execPath string        // FFmpeg可执行文件路径
	timeout  time.Duration // 命令执行超时时间
	logger   Logger        // 日志记录器
}

// Logger 日志记录接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// defaultLogger 默认日志记录器实现
type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("[INFO] "+msg+"\n", args...)
}

func (l *defaultLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("[ERROR] "+msg+"\n", args...)
}

func (l *defaultLogger) Debug(msg string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+msg+"\n", args...)
}

// Config FFmpeg配置选项
type Config struct {
	FFmpegPath string        // FFmpeg可执行文件路径，为空时自动检测
	Timeout    time.Duration // 命令执行超时时间，默认30分钟
	Logger     Logger        // 日志记录器，为空时使用默认记录器
}

// New 创建新的FFmpeg实例
func New(config *Config) (*FFmpeg, error) {
	if config == nil {
		config = &Config{}
	}

	// 设置默认超时时间
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Minute
	}

	// 设置默认日志记录器
	if config.Logger == nil {
		config.Logger = &defaultLogger{}
	}

	// 检测FFmpeg可执行文件路径
	execPath := config.FFmpegPath
	if execPath == "" {
		var err error
		execPath, err = detectFFmpegPath()
		if err != nil {
			return nil, fmt.Errorf("无法检测到FFmpeg可执行文件: %w", err)
		}
	}

	// 验证FFmpeg可执行文件是否存在且可执行
	if err := validateFFmpegPath(execPath); err != nil {
		return nil, fmt.Errorf("FFmpeg可执行文件验证失败: %w", err)
	}

	return &FFmpeg{
		execPath: execPath,
		timeout:  config.Timeout,
		logger:   config.Logger,
	}, nil
}

// detectFFmpegPath 自动检测FFmpeg可执行文件路径
func detectFFmpegPath() (string, error) {
	// 根据操作系统确定可执行文件名
	execName := "ffmpeg"
	if runtime.GOOS == "windows" {
		execName = "ffmpeg.exe"
	}

	// 首先尝试在PATH中查找
	if path, err := exec.LookPath(execName); err == nil {
		return path, nil
	}

	// 尝试常见的安装路径
	commonPaths := getCommonFFmpegPaths()
	for _, path := range commonPaths {
		fullPath := filepath.Join(path, execName)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("未找到FFmpeg可执行文件，请确保已安装FFmpeg并添加到PATH环境变量中")
}

// getCommonFFmpegPaths 获取常见的FFmpeg安装路径
func getCommonFFmpegPaths() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			"C:\\ffmpeg\\bin",
			"C:\\Program Files\\ffmpeg\\bin",
			"C:\\Program Files (x86)\\ffmpeg\\bin",
		}
	case "darwin": // macOS
		return []string{
			"/usr/local/bin",
			"/opt/homebrew/bin",
			"/usr/bin",
		}
	case "linux":
		return []string{
			"/usr/bin",
			"/usr/local/bin",
			"/opt/ffmpeg/bin",
		}
	default:
		return []string{"/usr/bin", "/usr/local/bin"}
	}
}

// validateFFmpegPath 验证FFmpeg可执行文件路径
func validateFFmpegPath(path string) error {
	// 检查文件是否存在
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("文件不存在: %s", path)
	}

	// 尝试执行版本命令来验证是否为有效的FFmpeg可执行文件
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, path, "-version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("无法执行FFmpeg命令: %w", err)
	}

	// 检查输出是否包含FFmpeg标识
	if !strings.Contains(string(output), "ffmpeg version") {
		return fmt.Errorf("不是有效的FFmpeg可执行文件")
	}

	return nil
}

// executeCommand 执行FFmpeg命令
func (f *FFmpeg) executeCommand(ctx context.Context, args []string) ([]byte, error) {
	f.logger.Debug("执行FFmpeg命令: %s %s", f.execPath, strings.Join(args, " "))

	// 创建带超时的上下文
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), f.timeout)
		defer cancel()
	}

	// 执行命令
	cmd := exec.CommandContext(ctx, f.execPath, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		f.logger.Error("FFmpeg命令执行失败: %s, 输出: %s", err.Error(), string(output))
		return output, fmt.Errorf("FFmpeg命令执行失败: %w", err)
	}

	f.logger.Debug("FFmpeg命令执行成功")
	return output, nil
}
