package notification

import (
	"FocusTimer/internal/config"
)

// Notifier 定义通知的接口
type Notifier interface {
	// Notify 发送通知
	Notify() error
	// Stop 停止通知
	Stop() error
}

// Factory 创建Notifier的工厂方法
func Factory(promptConfig *config.Prompt) []Notifier {
	var notifiers []Notifier

	if promptConfig.Popup != nil {
		notifiers = append(notifiers, NewPopupNotifier(promptConfig.Popup))
	}

	if promptConfig.Sound != nil {
		notifiers = append(notifiers, NewSoundNotifier(promptConfig.Sound))
	}

	return notifiers
}
