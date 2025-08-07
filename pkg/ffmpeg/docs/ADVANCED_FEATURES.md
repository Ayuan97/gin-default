# FFmpeg 高级功能使用指南

本文档介绍了 FFmpeg 封装库中新增的高级视频编辑功能，这些功能提供了与剪映（CapCut）类似的专业视频编辑能力。

## 🎨 高级滤镜系统

### 色彩分级

```go
// 创建色彩分级选项
colorOptions := &ffmpeg.ColorGradingOptions{
    Brightness:  0.1,   // 亮度调整
    Contrast:    0.2,   // 对比度调整
    Saturation:  0.15,  // 饱和度调整
    Hue:         10,    // 色调调整
    Temperature: 20,    // 色温调整（暖色调）
    Highlights:  -10,   // 高光调整
    Shadows:     15,    // 阴影调整
    Clarity:     25,    // 清晰度调整
    Vibrance:    20,    // 自然饱和度调整
}

// 应用色彩分级
err := ffmpeg.ApplyColorGrading("input.mp4", "output.mp4", colorOptions)

// 链式调用
editor := ffmpeg.NewVideoEditor("input.mp4")
editor.ApplyColorGrading(colorOptions).
       Export("output.mp4")
```

### 风格滤镜

```go
// 复古滤镜
vintageOptions := &ffmpeg.VintageFilterOptions{
    Sepia:      0.6,  // 棕褐色调
    Grain:      0.4,  // 胶片颗粒
    Vignette:   0.3,  // 暗角效果
    Fade:       0.2,  // 褪色效果
    Scratches:  true, // 添加划痕
    DustSpots:  true, // 添加灰尘斑点
}

editor.ApplyVintageFilter(vintageOptions)

// 电影风格滤镜
cinematicOptions := &ffmpeg.CinematicFilterOptions{
    AspectRatio:    "21:9",
    ColorGrading:   "teal_orange",
    FilmGrain:      0.3,
    Bloom:          0.4,
    LensFlare:      true,
}

editor.ApplyCinematicFilter(cinematicOptions)

// 美颜滤镜
beautyOptions := &ffmpeg.BeautyFilterOptions{
    SkinSmoothing:   0.6,
    SkinBrightening: 0.4,
    EyeEnhancement:  0.5,
    TeethWhitening:  0.3,
}

editor.ApplyBeautyFilter(beautyOptions)
```

## 🔄 转场效果系统

### 基础转场

```go
// 淡入淡出转场
editor.AddFadeTransition(2 * time.Second)

// 溶解转场
editor.AddDissolveTransition(1.5 * time.Second)

// 擦除转场
editor.AddWipeTransition(1 * time.Second, "left_to_right")

// 滑动转场
editor.AddSlideTransition(1.2 * time.Second, "left")
```

### 高级转场

```go
// 缩放转场
editor.AddZoomTransition(1.5 * time.Second, 2.0) // 2倍缩放

// 故障效果转场
editor.AddGlitchTransition(1 * time.Second, 0.7)

// 自定义转场选项
transitionOptions := &ffmpeg.AdvancedTransitionOptions{
    Type:      ffmpeg.TransitionCube,
    Duration:  2 * time.Second,
    Direction: "up",
    Easing:    "ease_in_out",
    Intensity: 0.8,
}

editor.AddAdvancedTransition(ffmpeg.TransitionCube, 2*time.Second, transitionOptions)
```

## 🎵 高级音频处理

### 音频均衡器

```go
// 预设均衡器
eqOptions := &ffmpeg.AudioEqualizerOptions{
    Preset:     "rock",
    MasterGain: 2.0,
}

editor.ApplyAudioEqualizer(eqOptions)

// 自定义频段
customEQ := &ffmpeg.AudioEqualizerOptions{
    Bands: []ffmpeg.AudioEqualizerBand{
        {Frequency: 60, Gain: 3, Q: 1.0},    // 低频增强
        {Frequency: 1000, Gain: 2, Q: 1.0},  // 中频增强
        {Frequency: 8000, Gain: 4, Q: 1.0},  // 高频增强
    },
    MasterGain: 1.0,
}

editor.ApplyAudioEqualizer(customEQ)
```

### 音频效果

```go
// 混响效果
reverbOptions := &ffmpeg.ReverbOptions{
    RoomSize:   0.7,
    Damping:    0.4,
    WetLevel:   0.3,
    DryLevel:   0.7,
    ReverbType: "hall",
}

editor.ApplyReverb(reverbOptions)

// 压缩器
compressorOptions := &ffmpeg.CompressorOptions{
    Threshold:  -18,
    Ratio:      4.0,
    Attack:     5,
    Release:    50,
    MakeupGain: 2,
}

editor.ApplyCompressor(compressorOptions)
```

## 📝 高级字幕系统

### 动画字幕

```go
// 基础动画字幕
editor.AddAnimatedText("Hello World!", ffmpeg.SubtitleAnimationTypewriter, 
                      0, 5*time.Second)

// 高级字幕选项
subtitleOptions := &ffmpeg.AdvancedSubtitleOptions{
    Text:              "专业字幕效果",
    StartTime:         2 * time.Second,
    Duration:          4 * time.Second,
    X:                 100,
    Y:                 100,
    FontSize:          32,
    FontWeight:        "bold",
    Color:             "#FFFFFF",
    OutlineColor:      "#000000",
    OutlineWidth:      2,
    ShadowColor:       "#333333",
    ShadowOffsetX:     3,
    ShadowOffsetY:     3,
    ShadowBlur:        5,
    Animation:         ffmpeg.SubtitleAnimationGlow,
    AnimationDuration: 800 * time.Millisecond,
}

editor.AddAdvancedSubtitle(subtitleOptions)
```

### 字幕模板

```go
// 使用内置模板
editor.ApplySubtitleTemplate("新闻标题", "news", 0, 3*time.Second)
editor.ApplySubtitleTemplate("游戏解说", "gaming", 5*time.Second, 4*time.Second)
editor.ApplySubtitleTemplate("教育内容", "educational", 10*time.Second, 6*time.Second)

// 下三分之一标题
editor.AddLowerThird("主标题", "副标题", 0, 5*time.Second)
```

## 🎬 高级合成功能

### 绿幕抠图

```go
// 绿幕抠图选项
chromaOptions := &ffmpeg.ChromaKeyOptions{
    KeyColor:         "#00FF00",  // 绿色
    Tolerance:        0.3,        // 容差
    Softness:         0.1,        // 边缘柔和度
    SpillSuppression: 0.2,        // 溢色抑制
    EdgeFeather:      0.05,       // 边缘羽化
    BackgroundPath:   "background.mp4", // 背景视频
}

editor.ApplyChromaKey(chromaOptions)
```

### 粒子效果

```go
// 雪花效果
editor.CreateSnowEffect(0.8, 0, 10*time.Second) // 强度0.8，持续10秒

// 雨滴效果
editor.CreateRainEffect(0.6, 5*time.Second, 8*time.Second)

// 火焰效果
editor.CreateFireEffect(0.7, 2*time.Second, 6*time.Second)

// 闪光效果
editor.CreateSparkleEffect(0.5, 0, 15*time.Second)

// 自定义粒子效果
particleOptions := &ffmpeg.ParticleEffectOptions{
    ParticleType: ffmpeg.ParticleTypeHearts,
    Count:        100,
    Size:         3.0,
    Speed:        25.0,
    Direction:    90,  // 向上
    Spread:       45,
    Gravity:      -0.2,
    Opacity:      0.8,
    Color:        "#FF69B4",
    BlendMode:    "screen",
    StartTime:    0,
    Duration:     8 * time.Second,
    EmissionRate: 10,
    LifeTime:     3.0,
}

editor.AddParticleEffect(particleOptions)
```

### 动态图形

```go
// 下三分之一图形
editor.AddLowerThirdGraphics("新闻标题", 0, 5*time.Second)

// 进度条
editor.AddProgressBar(2*time.Second, 8*time.Second, 50, 400, 300, 20)

// 计数器
editor.AddCounter(0, 10*time.Second, 100, 100)
```

## ⚡ 性能优化

### 性能预设

```go
// 速度优化
speedOptions := ffmpeg.OptimizeForSpeed()
editor.SetPerformanceOptions(speedOptions)

// 质量优化
qualityOptions := ffmpeg.OptimizeForQuality()
editor.SetPerformanceOptions(qualityOptions)

// 文件大小优化
sizeOptions := ffmpeg.OptimizeForSize()
editor.SetPerformanceOptions(sizeOptions)

// 平衡优化
balanceOptions := ffmpeg.OptimizeForBalance()
editor.SetPerformanceOptions(balanceOptions)
```

### 硬件加速

```go
// 自定义性能选项
perfOptions := &ffmpeg.PerformanceOptions{
    ThreadCount:     8,
    HardwareAccel:   "cuda",  // 使用NVIDIA GPU加速
    EnableGPU:       true,
    Quality:         "medium",
    Preset:          "fast",
    OptimizeFor:     "speed",
    EnableMultipass: false,
    BufferSize:      2048,
}

editor.SetPerformanceOptions(perfOptions)
```

## 🔗 链式调用示例

```go
// 完整的视频编辑流程
editor := ffmpeg.NewVideoEditor("input.mp4")

result := editor.
    // 基础编辑
    CropTime("00:00:10", "00:01:30").
    Resize(1920, 1080).
    
    // 色彩调整
    ApplyColorGrading(&ffmpeg.ColorGradingOptions{
        Brightness: 0.1,
        Contrast:   0.2,
        Saturation: 0.15,
    }).
    
    // 应用滤镜
    ApplyVintageFilter(&ffmpeg.VintageFilterOptions{
        Sepia:    0.4,
        Grain:    0.3,
        Vignette: 0.2,
    }).
    
    // 添加转场
    AddFadeTransition(1 * time.Second).
    
    // 音频处理
    ApplyReverb(&ffmpeg.ReverbOptions{
        RoomSize: 0.6,
        WetLevel: 0.3,
    }).
    
    // 添加字幕
    AddAnimatedText("精彩内容", ffmpeg.SubtitleAnimationTypewriter, 
                   5*time.Second, 3*time.Second).
    
    // 添加粒子效果
    CreateSnowEffect(0.5, 0, 10*time.Second).
    
    // 绿幕抠图
    ApplyChromaKey(&ffmpeg.ChromaKeyOptions{
        KeyColor:  "#00FF00",
        Tolerance: 0.3,
    }).
    
    // 性能优化
    SetPerformanceOptions(ffmpeg.OptimizeForBalance()).
    
    // 导出
    Export("final_output.mp4")

if result.Error != nil {
    log.Fatal("视频处理失败:", result.Error)
}

fmt.Printf("视频处理完成！耗时: %v\n", result.Duration)
```

## 📊 进度监控

```go
// 带进度回调的处理
editor.ExportWithProgress("output.mp4", func(progress float64, info string) {
    fmt.Printf("进度: %.1f%% - %s\n", progress, info)
})
```

这些新功能大大扩展了 FFmpeg 封装库的能力，使其能够处理专业级的视频编辑任务，提供与商业视频编辑软件相媲美的功能。
