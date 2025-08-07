# FFmpeg Go 包结构说明

## 文件组织结构

```
pkg/ffmpeg/
├── 📋 README.md               # 包入口文档和快速开始指南
├── 🚀 api.go                  # 统一API入口和便捷方法
│
├── 🔧 核心文件层 (core_*)
│   ├── core.go                # 核心FFmpeg结构体和基础功能
│   ├── types.go               # 类型定义和数据结构
│   └── utils.go               # 工具函数和辅助方法
│
├── 🎬 链式编辑器层 (editor_*)
│   ├── editor_core.go         # 链式调用视频编辑器
│   ├── editor_operations.go   # 具体操作实现和多媒体合成
│   └── editor_monitor.go      # 进度监控系统
│
├── ⚙️ 功能模块层 (*_ops.go)
│   ├── audio_ops.go           # 音频操作功能
│   ├── compress_ops.go        # 压缩操作功能
│   ├── convert_ops.go         # 格式转换功能
│   ├── edit_ops.go            # 编辑操作功能
│   └── info_ops.go            # 信息获取功能
│
├── 📁 docs/                   # 文档目录
│   ├── USAGE.md               # 使用说明
│   └── STRUCTURE.md           # 结构说明（本文件）
│
└── 🧪 tests/                  # 测试目录
    └── video_editor_test.go   # 单元测试
```

### 文件命名约定

为了在单一目录中保持清晰的逻辑分组，采用以下命名约定：

- **API 入口**: `api.go` - 统一的 API 入口
- **核心文件**: `core.go`, `types.go`, `utils.go` - 基础功能
- **链式编辑器**: `editor_*.go` - 链式调用相关功能
- **功能模块**: `*_ops.go` - 各种专业功能实现
- **文档**: `docs/` 目录 - 使用说明和结构文档
- **测试**: `tests/` 目录 - 单元测试文件

## 文件功能说明

### API 入口

#### `api.go`

- **功能**: 统一的 API 入口和便捷方法
- **主要内容**:
  - `API` 结构体 - 统一的 API 入口
  - `NewAPI()` 构造函数
  - 便捷方法封装 (`Convert`, `Compress`, `ExtractAudio` 等)
  - 快速编辑方法 (`QuickEdit`, `QuickCompose`)
  - 批量处理方法 (`BatchProcess`)
- **依赖**: `core.go`, `types.go`

### 核心文件

#### `core.go`

- **功能**: FFmpeg 核心结构体定义和基础功能
- **主要内容**:
  - `FFmpeg` 结构体
  - `New()` 构造函数
  - 基础命令执行方法
  - 配置管理
- **依赖**: `types.go`, `utils.go`

#### `types.go`

- **功能**: 所有类型定义和数据结构
- **主要内容**:
  - 配置结构体 (`Config`)
  - 错误类型定义 (`Error`, `ErrorCode`)
  - 操作选项结构体 (`WatermarkOptions`, `AudioMixOptions` 等)
  - 多媒体轨道结构体 (`VideoTrack`, `AudioTrack` 等)
  - 接口定义 (`VideoEditOperation`)
- **依赖**: 无

#### `utils.go`

- **功能**: 工具函数和辅助方法
- **主要内容**:
  - 文件验证函数
  - 路径处理函数
  - 时间解析函数
  - 日志记录器
- **依赖**: `types.go`

### 链式编辑器

#### `editor.go`

- **功能**: 链式调用视频编辑器核心实现
- **主要内容**:
  - `VideoEditor` 结构体
  - 链式调用方法 (`CropTime`, `Resize`, `AddWatermark` 等)
  - 多媒体合成方法 (`AddVideoTrack`, `PictureInPicture` 等)
  - 执行引擎 (`Execute`, `ExecuteWithContext`)
- **依赖**: `ffmpeg.go`, `types.go`, `operations.go`, `monitor.go`

#### `operations.go`

- **功能**: 具体操作的实现类
- **主要内容**:
  - 各种操作结构体 (`CropTimeOperation`, `ResizeOperation` 等)
  - 多媒体合成操作 (`AddTrackOperation`, `PictureInPictureOperation` 等)
  - 操作执行方法 (`Execute`)
  - FFmpeg 命令构建逻辑
- **依赖**: `types.go`

#### `monitor.go`

- **功能**: 进度监控和取消机制
- **主要内容**:
  - `ProgressMonitor` 结构体
  - `ProgressTracker` 多步骤进度跟踪
  - 进度解析和回调机制
  - 带进度监控的命令执行方法
- **依赖**: `types.go`

### 功能模块

#### `audio_ops.go`

- **功能**: 音频相关操作
- **主要内容**: 音频提取、混合、格式转换等
- **依赖**: `ffmpeg.go`, `types.go`

#### `compress_ops.go`

- **功能**: 视频压缩操作
- **主要内容**: 不同质量级别的压缩方法
- **依赖**: `ffmpeg.go`, `types.go`

#### `convert_ops.go`

- **功能**: 格式转换操作
- **主要内容**: 视频格式转换、编码器选择
- **依赖**: `ffmpeg.go`, `types.go`

#### `edit_ops.go`

- **功能**: 基础编辑操作
- **主要内容**: 裁剪、调整尺寸、旋转等基础功能
- **依赖**: `ffmpeg.go`, `types.go`

#### `info_ops.go`

- **功能**: 媒体信息获取
- **主要内容**: 视频信息解析、元数据提取
- **依赖**: `ffmpeg.go`, `types.go`

## 依赖关系图

```
types.go (基础类型)
    ↓
ffmpeg.go (核心功能) ← utils.go (工具函数)
    ↓
[功能模块] audio_ops.go, compress_ops.go, convert_ops.go, edit_ops.go, info_ops.go
    ↓
operations.go (操作实现) ← progress_monitor.go (进度监控)
    ↓
video_editor.go (链式编辑器)
    ↓
video_editor_test.go (测试)
```

## 推荐使用方式

### 基础使用

```go
import "your-project/pkg/ffmpeg"

// 1. 创建FFmpeg实例
ff, err := ffmpeg.New(nil)

// 2. 使用功能模块（传统方式）
err = ff.Convert("input.mp4", "output.avi", nil)

// 3. 使用链式编辑器（推荐方式）
err = ffmpeg.NewVideoEditor(ff, "input.mp4").
    SetOutput("output.mp4").
    CropTime("00:00:10", "00:01:00").
    Execute()
```

### 高级使用

```go
// 复杂的多媒体合成
editor := ffmpeg.NewVideoEditor(ff, "main.mp4")
err = editor.
    SetProgressCallback(progressCallback).
    AddVideoTrack("overlay.mp4", 10*time.Second, 15*time.Second).
    AddAudioTrack("music.mp3", 0*time.Second, 30*time.Second, 0.5).
    PictureInPicture("pip.mp4", 20*time.Second, 10*time.Second, 1200, 700, 400, 225).
    Execute()
```

## 扩展指南

### 添加新操作

1. 在 `types.go` 中定义操作选项结构体
2. 在 `operations.go` 中实现操作结构体和 `Execute` 方法
3. 在 `video_editor.go` 中添加链式调用方法
4. 在测试文件中添加相应测试

### 添加新功能模块

1. 创建新的 `*_ops.go` 文件
2. 实现具体功能方法
3. 在 `ffmpeg.go` 中添加对应的公开方法（如需要）

## 设计原则

- **单一职责**: 每个文件专注于特定功能领域
- **依赖倒置**: 核心功能不依赖具体实现
- **接口隔离**: 使用接口定义操作契约
- **开闭原则**: 易于扩展新功能，无需修改现有代码
