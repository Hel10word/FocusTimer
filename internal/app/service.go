package app

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"FocusTimer/internal/config"
	"FocusTimer/internal/timer"
	"FocusTimer/pkg/logger"

	// 导入平台特定代码
	"FocusTimer/internal/platform/windows"
)

// Service 是应用的主服务
type Service struct {
	config    *config.Config
	scheduler *timer.Scheduler
	logger    *logger.Logger
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewService 创建新的服务实例
func NewService(configPath string) (*Service, error) {
	// 加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	// 初始化日志
	log := logger.New(cfg.General.LogLevel)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())

	return &Service{
		config: cfg,
		logger: log,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Run 运行服务
func (s *Service) Run() error {
	// 创建并启动调度器
	s.scheduler = timer.NewScheduler(s.config, s.logger)

	// 启动系统托盘
	go s.startSystemTray()

	// 设置信号处理
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动调度器
	s.scheduler.Start()
	s.logger.Info("服务已启动")

	// 显示一个启动通知
	s.showStartupNotification()

	// 等待终止信号
	select {
	case <-signalCh:
		s.logger.Info("收到终止信号")
	case <-s.ctx.Done():
		s.logger.Info("上下文已取消")
	}

	// 优雅停止
	s.logger.Info("正在关闭...")
	s.scheduler.Stop()

	return nil
}

// Stop 停止服务
func (s *Service) Stop() {
	s.cancel()
}

// startSystemTray 启动系统托盘
func (s *Service) startSystemTray() {
	switch runtime.GOOS {
	case "windows":
		tray := windows.NewTrayService(s.config, s.logger, s.Stop)
		tray.Run()
	case "darwin":
		// macOS版本
		s.logger.Info("macOS系统托盘尚未实现")
	case "linux":
		// Linux版本
		s.logger.Info("Linux系统托盘尚未实现")
	default:
		s.logger.Warn("不支持的操作系统: %s", runtime.GOOS)
	}
}

// showStartupNotification 显示启动通知
func (s *Service) showStartupNotification() {
	// 创建一个启动提示
	popup := &config.PopupConfig{
		Text:            "FocusTimer已启动并在后台运行",
		Duration:        5 * time.Second,
		AlwaysOnTop:     true,
		BackgroundColor: "#4CAF50",
		TextColor:       "#FFFFFF",
		FontSize:        14,
	}

	notifier := windows.NewPopupNotifier(popup)
	notifier.Notify()
}
