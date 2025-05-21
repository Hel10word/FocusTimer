package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"FocusTimer/internal/config"
	"FocusTimer/internal/timer"
	"FocusTimer/pkg/logger"
)

// Service 是应用的主服务
type Service struct {
	config    *config.Config
	scheduler *timer.Scheduler
	logger    *logger.Logger
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

	// 创建调度器
	scheduler := timer.NewScheduler(cfg, log)

	return &Service{
		config:    cfg,
		scheduler: scheduler,
		logger:    log,
	}, nil
}

// Run 运行服务
func (s *Service) Run() error {

	// 设置信号处理
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动调度器
	s.scheduler.Start()
	s.logger.Info("Service started")

	// 等待终止信号
	select {
	case <-signalCh:
		s.logger.Info("Received termination signal")
	case <-ctx.Done():
		s.logger.Info("Context cancelled")
	}

	// 优雅停止
	s.logger.Info("Shutting down...")
	s.scheduler.Stop()

	return nil
}
