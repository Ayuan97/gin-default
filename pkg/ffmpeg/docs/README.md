# FFmpeg Go 封装库

一个功能强大的 FFmpeg Go 语言封装库，提供链式调用 API 和多媒体合成功能。

## ✨ 特性

- 🔗 **链式调用设计** - 流畅的 API，支持连续操作
- 📊 **进度监控** - 实时进度回调和取消支持
- 🎬 **多媒体合成** - 支持多轨道视频、音频、图片合成
- ⚡ **高性能** - 基于 FFmpeg 的高效处理
- 🛠️ **易于使用** - 统一的 API 入口和便捷方法

## 🚀 快速开始

### 安装

确保系统已安装 FFmpeg：

```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt update && sudo apt install ffmpeg
```

### 基本使用

```go
import "your-project/pkg/ffmpeg"

// 创建API实例
api, err := ffmpeg.NewAPI(nil)
if err != nil {
    log.Fatal(err)
}

// 链式视频编辑
err = api.NewEditor("input.mp4").
    SetOutput("output.mp4").
    CropTime("00:00:10", "00:01:30").
    Resize(1280, 720).
    AddWatermark(&ffmpeg.WatermarkOptions{
        ImagePath: "logo.png",
        X: 10, Y: 10, Scale: 0.2, Opacity: 0.8,
    }).
    Execute()
```

### 多媒体合成

```go
// 复杂的多媒体合成
err = api.NewEditor("main.mp4").
    SetOutput("composition.mp4").
    AddVideoTrack("overlay.mp4", 10*time.Second, 15*time.Second).
    AddAudioTrack("music.mp3", 0*time.Second, 30*time.Second, 0.5).
    AddTextTrack("标题", 5*time.Second, 3*time.Second, 100, 100, 24, "white").
    PictureInPicture("pip.mp4", 20*time.Second, 10*time.Second, 1200, 700, 400, 225).
    Execute()
```

## 📁 目录结构

```
pkg/ffmpeg/
├── 📋 ffmpeg_api.go           # 统一API入口
├── 🔧 核心文件
│   ├── ffmpeg.go              # FFmpeg核心功能
│   ├── types.go               # 类型定义
│   └── utils.go               # 工具函数
├── 🎬 链式编辑器
│   ├── editor.go              # 视频编辑器
│   ├── operations.go          # 操作实现
│   └── monitor.go             # 进度监控
├── ⚙️ 功能模块
│   ├── audio_ops.go           # 音频操作
│   ├── compress_ops.go        # 压缩功能
│   ├── convert_ops.go         # 格式转换
│   ├── edit_ops.go            # 编辑操作
│   └── info_ops.go            # 信息获取
├── 📁 docs/                   # 文档目录
└── 🧪 tests/                  # 测试目录
```

## 📖 文档

- [使用说明](docs/USAGE.md) - 详细的使用指南和API参考
- [结构说明](docs/STRUCTURE.md) - 包结构和文件功能说明

## 🎯 核心功能

### 视频编辑
- 时间裁剪、尺寸裁剪、调整尺寸
- 旋转、镜像翻转、淡入淡出
- 亮度/对比度调整、模糊效果、防抖

### 音频处理
- 音频分离、混合、背景音乐添加
- 音量调节、淡入淡出效果

### 多媒体合成
- 多轨道视频、音频、图片合成
- 画中画、分屏显示、字幕添加
- 水印、文字、转场效果

### 进度监控
- 实时进度回调
- 操作取消支持
- 多步骤进度跟踪

## 🔧 配置选项

```go
config := &ffmpeg.Config{
    Timeout:  10 * time.Minute,  // 超时时间
    LogLevel: "info",             // 日志级别
    ExecPath: "/usr/bin/ffmpeg",  // FFmpeg路径（可选）
}

api, err := ffmpeg.NewAPI(config)
```

## 📝 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**注意**: 使用前请确保系统已正确安装 FFmpeg。
