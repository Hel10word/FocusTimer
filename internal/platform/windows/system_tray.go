package windows

import (
	"FocusTimer/internal/config"
	"FocusTimer/pkg/logger"
	"github.com/getlantern/systray"
	"os"
)

// TrayService 处理系统托盘
type TrayService struct {
	config *config.Config
	logger *logger.Logger
	onExit func()
}

// NewTrayService 创建新的系统托盘服务
func NewTrayService(cfg *config.Config, log *logger.Logger, exitFunc func()) *TrayService {
	return &TrayService{
		config: cfg,
		logger: log,
		onExit: exitFunc,
	}
}

// Run 启动系统托盘
func (ts *TrayService) Run() {
	systray.Run(ts.onReady, ts.onExit)
}

// onReady 系统托盘准备就绪
func (ts *TrayService) onReady() {
	// 设置图标
	iconBytes, err := os.ReadFile("assets/icons/FocusTimer.ico")
	if err != nil {
		ts.logger.Error("无法加载托盘图标: %v", err)
		// 使用默认图标
	} else {
		systray.SetIcon(iconBytes)
	}

	systray.SetTitle("FocusTimer")
	systray.SetTooltip("FocusTimer - 时间管理工具")

	// 添加菜单项
	mStatus := systray.AddMenuItem("状态: 运行中", "显示当前状态")
	mStatus.Disable()

	systray.AddSeparator()

	mPause := systray.AddMenuItem("暂停", "暂停计时器")
	mResume := systray.AddMenuItem("继续", "继续计时器")
	mResume.Disable()

	systray.AddSeparator()

	mShowSettings := systray.AddMenuItem("设置", "打开设置")

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("退出", "退出程序")

	// 处理菜单点击
	go func() {
		for {
			select {
			case <-mPause.ClickedCh:
				mStatus.SetTitle("状态: 已暂停")
				mPause.Disable()
				mResume.Enable()
				// TODO: 实现暂停逻辑

			case <-mResume.ClickedCh:
				mStatus.SetTitle("状态: 运行中")
				mResume.Disable()
				mPause.Enable()
				// TODO: 实现继续逻辑

			case <-mShowSettings.ClickedCh:
				// TODO: 打开设置窗口

			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}
