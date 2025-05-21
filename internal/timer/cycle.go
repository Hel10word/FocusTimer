package timer

import (
	"context"
	"math/rand"
	"time"

	"FocusTimer/internal/config"
	"FocusTimer/internal/notification"
)

// Cycle 管理单个时间循环
type Cycle struct {
	config       config.Cycle
	notifiers    []notification.Notifier
	endNotifiers []notification.Notifier
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewCycle 创建新循环
func NewCycle(cfg config.Cycle) *Cycle {
	ctx, cancel := context.WithCancel(context.Background())
	return &Cycle{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Run 运行循环
func (c *Cycle) Run(parentCtx context.Context) {
	// 创建带有周期时长的上下文
	ctx, cancel := context.WithTimeout(parentCtx, c.config.Duration)
	defer cancel()

	c.ctx = ctx

	// 运行提示循环
	cycleEnd := time.After(c.config.Duration)

	for {
		for _, prompt := range c.config.Prompts {
			// 随机等待间隔
			interval := c.getRandomInterval(prompt.Interval)
			select {
			case <-time.After(interval):
				c.triggerPrompt(prompt)
			case <-cycleEnd:
				goto CycleComplete
			case <-ctx.Done():
				return
			}
		}
	}

CycleComplete:
	// 触发循环结束提示
	if c.config.EndPrompt != nil {
		c.triggerPrompt(*c.config.EndPrompt)
	}
}

// Stop 停止循环
func (c *Cycle) Stop() {
	c.cancel()
}

// getRandomInterval 根据配置获取随机间隔
func (c *Cycle) getRandomInterval(interval config.Interval) time.Duration {
	if interval.Min == interval.Max {
		return interval.Min
	}

	diff := interval.Max - interval.Min
	return interval.Min + time.Duration(rand.Int63n(int64(diff)))
}

// triggerPrompt 触发提示
func (c *Cycle) triggerPrompt(prompt config.Prompt) {
	notifiers := notification.Factory(&prompt)

	for _, n := range notifiers {
		go func(notifier notification.Notifier) {
			err := notifier.Notify()
			if err != nil {
				// 处理错误
			}

			// 如果提示配置了持续时间，在持续时间后停止
			if prompt.Type == "popup" && prompt.Popup != nil {
				time.Sleep(prompt.Popup.Duration)
				notifier.Stop()
			} else if prompt.Type == "sound" && prompt.Sound != nil {
				time.Sleep(prompt.Sound.Duration)
				notifier.Stop()
			}
		}(n)
	}
}
