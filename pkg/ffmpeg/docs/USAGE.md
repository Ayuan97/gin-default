# FFmpeg Go 使用说明

## 快速开始

### 推荐方式：使用 API 入口

```go
import "your-project/pkg/ffmpeg"

// 创建API实例（推荐）
api, err := ffmpeg.NewAPI(&ffmpeg.Config{
    Timeout: 10 * time.Minute,
})

// 便捷方法调用
err = api.Convert("input.mp4", "output.avi", nil)
```

### 链式视频编辑

```go
// 基础编辑
err = api.NewEditor("input.mp4").
    SetOutput("output.mp4").
    CropTime("00:00:10", "00:01:30").
    Resize(1280, 720).
    Execute()
```

### 传统方式：直接使用 FFmpeg 实例

```go
// 创建FFmpeg实例
ff, err := ffmpeg.New(&ffmpeg.Config{
    Timeout: 10 * time.Minute,
})

// 使用链式编辑器
err = ffmpeg.NewVideoEditor(ff, "input.mp4").
    SetOutput("output.mp4").
    CropTime("00:00:10", "00:01:30").
    Execute()
```

### 进度监控

```go
editor.SetProgressCallback(func(progress float64, currentTime, totalTime time.Duration) {
    fmt.Printf("进度: %.1f%% (%v / %v)\n", progress, currentTime, totalTime)
})
```

### 多媒体合成

```go
// 多轨道合成
err = editor.
    AddVideoTrack("overlay.mp4", 10*time.Second, 15*time.Second).
    AddAudioTrack("music.mp3", 0*time.Second, 30*time.Second, 0.5).
    AddTextTrack("标题", 5*time.Second, 3*time.Second, 100, 100, 24, "white").
    Execute()
```

## 核心功能

### 视频编辑

- `CropTime()` - 时间裁剪
- `CropDimension()` - 尺寸裁剪
- `Resize()` - 调整尺寸
- `Rotate()` - 旋转
- `Mirror()` - 镜像翻转

### 视觉效果

- `FadeIn()` / `FadeOut()` - 淡入淡出
- `AdjustBrightness()` - 亮度调整
- `AdjustContrast()` - 对比度调整
- `AddBlur()` - 模糊效果
- `Stabilize()` - 视频防抖

### 音频处理

- `SeparateAudio()` - 音频分离
- `MixAudio()` - 音频混合
- `AddBackgroundMusic()` - 背景音乐

### 图片和文字

- `AddWatermark()` - 添加水印
- `InsertImage()` - 插入图片
- `AddText()` - 添加文字

### 多媒体合成

- `AddVideoTrack()` - 视频轨道
- `AddAudioTrack()` - 音频轨道
- `AddImageTrack()` - 图片轨道
- `AddTextTrack()` - 文字轨道
- `PictureInPicture()` - 画中画
- `SplitScreen()` - 分屏显示

## 配置选项

### Config 配置

```go
&ffmpeg.Config{
    Timeout:     10 * time.Minute,  // 超时时间
    LogLevel:    "info",             // 日志级别
    ExecPath:    "/usr/bin/ffmpeg",  // FFmpeg路径（可选）
}
```

### WatermarkOptions 水印选项

```go
&ffmpeg.WatermarkOptions{
    ImagePath: "watermark.png",
    X:         10,
    Y:         10,
    Scale:     0.2,
    Opacity:   0.8,
}
```

### AudioMixOptions 音频混合选项

```go
&ffmpeg.AudioMixOptions{
    BackgroundPath: "music.mp3",
    Volume:         0.5,
    Loop:           true,
}
```

## 错误处理

```go
if err != nil {
    if ffmpegErr, ok := err.(*ffmpeg.Error); ok {
        switch ffmpegErr.Code {
        case ffmpeg.ErrFileNotFound:
            // 处理文件不存在错误
        case ffmpeg.ErrExecutionFailed:
            // 处理执行失败错误
        }
    }
}
```

## 注意事项

- 确保系统已安装 FFmpeg
- 输入文件路径必须存在且可读
- 输出目录必须存在且可写
- 使用进度回调监控长时间操作
- 链式调用操作按顺序执行
