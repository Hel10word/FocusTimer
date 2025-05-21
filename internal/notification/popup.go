package notification

import (
	"fmt"
	"runtime"

	"FocusTimer/internal/config"
	"FocusTimer/internal/platform/windows"
	// 未来可以添加其他平台的导入
)

// PopupNotifier 表示一个弹窗通知器
type PopupNotifier struct {
	config   *config.PopupConfig
	platform Notifier
}

// NewPopupNotifier 创建一个新的弹窗通知器
func NewPopupNotifier(cfg *config.PopupConfig) Notifier {
	var platform Notifier

	// 根据操作系统选择不同的实现
	switch runtime.GOOS {
	case "windows":
		platform = windows.NewPopupNotifier(cfg)
	case "darwin":
		// TODO: 实现macOS弹窗
		platform = &unsupportedNotifier{os: "macOS", notificationType: "popup"}
	case "linux":
		// TODO: 实现Linux弹窗
		platform = &unsupportedNotifier{os: "Linux", notificationType: "popup"}
	default:
		return &unsupportedNotifier{os: runtime.GOOS, notificationType: "popup"}
	}

	return &PopupNotifier{
		config:   cfg,
		platform: platform,
	}
}

// Notify 显示弹窗
func (p *PopupNotifier) Notify() error {
	return p.platform.Notify()
}

// Stop 关闭弹窗
func (p *PopupNotifier) Stop() error {
	return p.platform.Stop()
}

// unsupportedNotifier 表示不支持的通知器
type unsupportedNotifier struct {
	os               string
	notificationType string
}

func (u *unsupportedNotifier) Notify() error {
	return fmt.Errorf("不支持在 %s 上的 %s 通知", u.os, u.notificationType)
}

func (u *unsupportedNotifier) Stop() error {
	return nil
}
