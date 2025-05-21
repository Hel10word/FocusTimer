package config

import (
	"time"
)

// Config 代表应用的总体配置
type Config struct {
	Cycles         []Cycle `yaml:"cycles"`
	General        General `yaml:"general"`
	AutoStart      bool    `yaml:"autoStart"`
	MinimizeToTray bool    `yaml:"minimizeToTray"`
}

// General 包含应用的全局设置
type General struct {
	Language string `yaml:"language"`
	LogLevel string `yaml:"logLevel"`
}

// Cycle 代表一个时间循环
type Cycle struct {
	Name        string        `yaml:"name"`
	Duration    time.Duration `yaml:"duration"`
	Prompts     []Prompt      `yaml:"prompts"`
	EndPrompt   *Prompt       `yaml:"endPrompt"`
	RepeatCount int           `yaml:"repeatCount"` // -1表示无限重复
}

// Prompt 代表一个提示
type Prompt struct {
	Type     string       `yaml:"type"` // popup, sound, both
	Interval Interval     `yaml:"interval"`
	Popup    *PopupConfig `yaml:"popup,omitempty"`
	Sound    *SoundConfig `yaml:"sound,omitempty"`
}

// Interval 表示提示间隔
type Interval struct {
	Min time.Duration `yaml:"min"`
	Max time.Duration `yaml:"max"`
}

// PopupConfig 配置弹窗样式
type PopupConfig struct {
	Text            string        `yaml:"text"`
	Duration        time.Duration `yaml:"duration"`
	AlwaysOnTop     bool          `yaml:"alwaysOnTop"`
	BackgroundColor string        `yaml:"backgroundColor"`
	TextColor       string        `yaml:"textColor"`
	FontSize        int           `yaml:"fontSize"`
}

// SoundConfig 配置声音提示
type SoundConfig struct {
	FilePath string        `yaml:"filePath"`
	Loop     bool          `yaml:"loop"`
	Duration time.Duration `yaml:"duration"`
	Volume   float64       `yaml:"volume"`
}
