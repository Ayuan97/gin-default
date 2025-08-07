package ffmpeg

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestVideoEditorCreation 测试视频编辑器创建
func TestVideoEditorCreation(t *testing.T) {
	// 创建模拟的FFmpeg实例
	ffmpeg := &FFmpeg{
		execPath: "/usr/bin/ffmpeg", // 假设路径
		timeout:  30 * time.Second,
		logger:   &DefaultLogger{},
	}

	// 创建视频编辑器
	editor := NewVideoEditor(ffmpeg, "test_input.mp4")

	if editor == nil {
		t.Fatal("视频编辑器创建失败")
	}

	if editor.inputPath != "test_input.mp4" {
		t.Errorf("输入路径设置错误，期望: test_input.mp4, 实际: %s", editor.inputPath)
	}

	if editor.ffmpeg != ffmpeg {
		t.Error("FFmpeg实例设置错误")
	}

	if len(editor.operations) != 0 {
		t.Error("初始操作队列应该为空")
	}
}

// TestChainedCalls 测试链式调用
func TestChainedCalls(t *testing.T) {
	ffmpeg := &FFmpeg{
		execPath: "/usr/bin/ffmpeg",
		timeout:  30 * time.Second,
		logger:   &DefaultLogger{},
	}

	editor := NewVideoEditor(ffmpeg, "test_input.mp4")

	// 测试链式调用
	result := editor.
		SetOutput("test_output.mp4").
		CropTime("00:00:10", "00:01:00").
		Resize(1280, 720).
		AddWatermark(&WatermarkOptions{
			ImagePath: "watermark.png",
			X:         10,
			Y:         10,
			Scale:     0.2,
			Opacity:   0.8,
		})

	// 验证返回的是同一个实例
	if result != editor {
		t.Error("链式调用应该返回同一个编辑器实例")
	}

	// 验证输出路径设置
	if editor.outputPath != "test_output.mp4" {
		t.Errorf("输出路径设置错误，期望: test_output.mp4, 实际: %s", editor.outputPath)
	}

	// 验证操作数量
	expectedOperations := 3 // CropTime, Resize, AddWatermark
	if len(editor.operations) != expectedOperations {
		t.Errorf("操作数量错误，期望: %d, 实际: %d", expectedOperations, len(editor.operations))
	}
}

// TestProgressCallback 测试进度回调设置
func TestProgressCallback(t *testing.T) {
	ffmpeg := &FFmpeg{
		execPath: "/usr/bin/ffmpeg",
		timeout:  30 * time.Second,
		logger:   &DefaultLogger{},
	}

	editor := NewVideoEditor(ffmpeg, "test_input.mp4")

	callbackCalled := false
	progressCallback := func(progress float64, currentTime time.Duration, totalTime time.Duration) {
		callbackCalled = true
	}

	editor.SetProgressCallback(progressCallback)

	// 验证回调函数设置
	if editor.progress == nil {
		t.Error("进度回调函数未设置")
	}

	// 模拟调用回调函数
	if editor.progress != nil {
		editor.progress(50.0, 30*time.Second, 60*time.Second)
	}

	if !callbackCalled {
		t.Error("进度回调函数未被调用")
	}
}

// TestOperationTypes 测试不同操作类型
func TestOperationTypes(t *testing.T) {
	ffmpeg := &FFmpeg{
		execPath: "/usr/bin/ffmpeg",
		timeout:  30 * time.Second,
		logger:   &DefaultLogger{},
	}

	editor := NewVideoEditor(ffmpeg, "test_input.mp4")

	// 添加各种类型的操作
	editor.
		CropTime("00:00:10", "00:01:00").
		CropDimension(&CropDimensions{X: 0, Y: 0, Width: 800, Height: 600}).
		Resize(1280, 720).
		Rotate(90).
		Mirror(true).
		AdjustBrightness(0.1).
		AdjustContrast(0.2).
		AddBlur(2.0).
		FadeIn(2 * time.Second).
		FadeOut(2 * time.Second).
		ChangeSpeed(1.5).
		AddText("测试文字", 50, 50, 24, "white", 5*time.Second, 10*time.Second).
		AddWatermark(&WatermarkOptions{
			ImagePath: "watermark.png",
			X: 10, Y: 10, Scale: 0.2, Opacity: 0.8,
		}).
		InsertImage("overlay.png", 15*time.Second, 5*time.Second).
		MixAudio(&AudioMixOptions{
			BackgroundPath: "bg_music.mp3",
			Volume: 0.5, Loop: true,
		}).
		SeparateAudio("extracted_audio.mp3").
		EditFrame(&FrameEditOptions{
			Operation: FrameInsert,
			FrameNumber: 100,
			ImagePath: "frame.png",
		}).
		Stabilize().
		ExtractFrames("/tmp/frames", 1.0).
		CreateFromImages("image_%04d.png", 30.0)

	expectedOperations := 20
	if len(editor.operations) != expectedOperations {
		t.Errorf("操作数量错误，期望: %d, 实际: %d", expectedOperations, len(editor.operations))
	}

	// 验证操作类型
	operations := editor.operations
	if _, ok := operations[0].(*CropTimeOperation); !ok {
		t.Error("第一个操作应该是CropTimeOperation")
	}

	if _, ok := operations[1].(*CropDimensionOperation); !ok {
		t.Error("第二个操作应该是CropDimensionOperation")
	}

	if _, ok := operations[2].(*ResizeOperation); !ok {
		t.Error("第三个操作应该是ResizeOperation")
	}
}

// TestOperationDescriptions 测试操作描述
func TestOperationDescriptions(t *testing.T) {
	// 测试各种操作的描述
	cropOp := &CropTimeOperation{StartTime: "00:00:10", Duration: "00:01:00"}
	if cropOp.GetDescription() == "" {
		t.Error("CropTimeOperation描述不能为空")
	}

	resizeOp := &ResizeOperation{Width: 1280, Height: 720}
	if resizeOp.GetDescription() == "" {
		t.Error("ResizeOperation描述不能为空")
	}

	watermarkOp := &WatermarkOperation{
		Options: &WatermarkOptions{
			ImagePath: "watermark.png",
			X: 10, Y: 10, Scale: 0.2, Opacity: 0.8,
		},
	}
	if watermarkOp.GetDescription() == "" {
		t.Error("WatermarkOperation描述不能为空")
	}

	textOp := &TextOperation{
		Text: "测试文字", X: 50, Y: 50,
		FontSize: 24, Color: "white",
		StartTime: 5 * time.Second,
		Duration: 10 * time.Second,
	}
	if textOp.GetDescription() == "" {
		t.Error("TextOperation描述不能为空")
	}
}

// TestOperationEstimation 测试操作时间估算
func TestOperationEstimation(t *testing.T) {
	cropOp := &CropTimeOperation{StartTime: "00:00:10", Duration: "00:01:00"}
	if cropOp.EstimateDuration() <= 0 {
		t.Error("操作时间估算应该大于0")
	}

	resizeOp := &ResizeOperation{Width: 1280, Height: 720}
	if resizeOp.EstimateDuration() <= 0 {
		t.Error("操作时间估算应该大于0")
	}

	stabilizeOp := &StabilizeOperation{}
	if stabilizeOp.EstimateDuration() <= resizeOp.EstimateDuration() {
		t.Error("防抖操作应该比调整尺寸操作耗时更长")
	}
}

// TestEditorMethods 测试编辑器方法
func TestEditorMethods(t *testing.T) {
	ffmpeg := &FFmpeg{
		execPath: "/usr/bin/ffmpeg",
		timeout:  30 * time.Second,
		logger:   &DefaultLogger{},
	}

	editor := NewVideoEditor(ffmpeg, "test_input.mp4")

	// 添加一些操作
	editor.CropTime("00:00:10", "00:01:00").Resize(1280, 720)

	// 测试获取操作数量
	if editor.GetOperationCount() != 2 {
		t.Errorf("操作数量错误，期望: 2, 实际: %d", editor.GetOperationCount())
	}

	// 测试清空操作
	editor.Clear()
	if editor.GetOperationCount() != 0 {
		t.Error("清空操作后，操作数量应该为0")
	}

	// 测试取消操作
	editor.Cancel()
	if !editor.cancelled {
		t.Error("取消操作后，cancelled标志应该为true")
	}

	// 测试克隆
	clone := editor.Clone()
	if clone == editor {
		t.Error("克隆应该返回新的实例")
	}
	if clone.inputPath != editor.inputPath {
		t.Error("克隆的输入路径应该相同")
	}
	if clone.GetOperationCount() != 0 {
		t.Error("克隆的操作队列应该为空")
	}
}

// TestValidation 测试输入验证
func TestValidation(t *testing.T) {
	ffmpeg := &FFmpeg{
		execPath: "/usr/bin/ffmpeg",
		timeout:  30 * time.Second,
		logger:   &DefaultLogger{},
	}

	editor := NewVideoEditor(ffmpeg, "test_input.mp4")

	// 测试没有设置输出路径的情况
	err := editor.CropTime("00:00:10", "00:01:00").Execute()
	if err == nil {
		t.Error("没有设置输出路径应该返回错误")
	}

	// 测试没有操作的情况
	editor.SetOutput("test_output.mp4").Clear()
	err = editor.Execute()
	if err == nil {
		t.Error("没有操作应该返回错误")
	}
}

// TestContextCancellation 测试上下文取消
func TestContextCancellation(t *testing.T) {
	ffmpeg := &FFmpeg{
		execPath: "/usr/bin/ffmpeg",
		timeout:  30 * time.Second,
		logger:   &DefaultLogger{},
	}

	editor := NewVideoEditor(ffmpeg, "test_input.mp4")

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	
	// 立即取消
	cancel()

	// 尝试执行（应该被取消）
	err := editor.
		SetOutput("test_output.mp4").
		CropTime("00:00:10", "00:01:00").
		ExecuteWithContext(ctx)

	if err == nil {
		t.Error("上下文取消后应该返回错误")
	}

	if err != context.Canceled {
		t.Errorf("应该返回context.Canceled错误，实际: %v", err)
	}
}

// BenchmarkChainedCalls 性能测试：链式调用
func BenchmarkChainedCalls(b *testing.B) {
	ffmpeg := &FFmpeg{
		execPath: "/usr/bin/ffmpeg",
		timeout:  30 * time.Second,
		logger:   &DefaultLogger{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		editor := NewVideoEditor(ffmpeg, "test_input.mp4")
		editor.
			SetOutput("test_output.mp4").
			CropTime("00:00:10", "00:01:00").
			Resize(1280, 720).
			AddWatermark(&WatermarkOptions{
				ImagePath: "watermark.png",
				X: 10, Y: 10, Scale: 0.2, Opacity: 0.8,
			}).
			FadeIn(2 * time.Second).
			FadeOut(2 * time.Second)
	}
}

// TestMain 测试主函数
func TestMain(m *testing.M) {
	// 设置测试环境
	// 这里可以添加测试前的准备工作

	// 运行测试
	code := m.Run()

	// 清理测试环境
	// 这里可以添加测试后的清理工作

	os.Exit(code)
}
