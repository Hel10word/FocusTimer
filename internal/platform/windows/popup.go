package windows

import (
	"FocusTimer/internal/config"
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

// PopupNotifier 实现Windows平台的弹窗
type PopupNotifier struct {
	config *config.PopupConfig
	form   *walk.MainWindow
}

// NewPopupNotifier 创建Windows弹窗通知
func NewPopupNotifier(cfg *config.PopupConfig) *PopupNotifier {
	return &PopupNotifier{
		config: cfg,
	}
}

// Notify 显示弹窗
func (p *PopupNotifier) Notify() error {
	// 创建Windows弹窗
	form, err := walk.NewMainWindow()
	if err != nil {
		return err
	}

	p.form = form

	// 设置窗口属性
	form.SetTitle("TimeReminder")
	form.SetSize(walk.Size{Width: 300, Height: 200})
	form.SetLayout(walk.NewVBoxLayout())

	if p.config.AlwaysOnTop {
		win.SetWindowPos(form.Handle(), win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)
	}

	// 创建文本标签
	textLabel, err := walk.NewTextLabel(form)
	if err != nil {
		return err
	}
	textLabel.SetText(p.config.Text)

	// 显示窗口
	form.Show()

	return nil
}

// Stop 关闭弹窗
func (p *PopupNotifier) Stop() error {
	if p.form != nil {
		p.form.Close()
	}
	return nil
}
