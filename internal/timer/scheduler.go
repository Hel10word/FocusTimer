package timer

import (
	"FocusTimer/internal/config"
	"FocusTimer/pkg/logger"
	"context"
)

// Scheduler 管理时间循环和提示调度
type Scheduler struct {
	cycles       []*Cycle
	currentCycle *Cycle
	ctx          context.Context
	cancel       context.CancelFunc
	logger       *logger.Logger
}

// NewScheduler 创建新的调度器
func NewScheduler(cfg *config.Config, log *logger.Logger) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	var cycles []*Cycle
	for _, cycleCfg := range cfg.Cycles {
		cycles = append(cycles, NewCycle(cycleCfg))
	}

	return &Scheduler{
		cycles: cycles,
		ctx:    ctx,
		cancel: cancel,
		logger: log,
	}
}

// Start 开始调度
func (s *Scheduler) Start() {
	go s.run()
}

// Stop 停止调度
func (s *Scheduler) Stop() {
	s.cancel()
	if s.currentCycle != nil {
		s.currentCycle.Stop()
	}
}

// run 运行调度循环
func (s *Scheduler) run() {
	for {
		for _, cycle := range s.cycles {
			select {
			case <-s.ctx.Done():
				return
			default:
				s.currentCycle = cycle
				s.logger.Info("Starting cycle: %s", cycle.config.Name)
				cycle.Run(s.ctx)
			}
		}
	}
}
