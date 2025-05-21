package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Load 从指定路径加载配置文件
func Load(configPath string) (*Config, error) {
	// 检查文件是否存在
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取文件内容
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// validateConfig 验证配置有效性
func validateConfig(config *Config) error {
	// 验证循环配置
	if len(config.Cycles) == 0 {
		return fmt.Errorf("至少需要一个时间循环")
	}

	for i, cycle := range config.Cycles {
		if cycle.Duration <= 0 {
			return fmt.Errorf("循环 %d 的持续时间必须大于0", i+1)
		}

		if len(cycle.Prompts) == 0 {
			return fmt.Errorf("循环 %d 至少需要一个提示", i+1)
		}

		// 验证提示配置
		for j, prompt := range cycle.Prompts {
			if prompt.Interval.Min <= 0 {
				return fmt.Errorf("循环 %d 的提示 %d 的最小间隔必须大于0", i+1, j+1)
			}

			if prompt.Interval.Max < prompt.Interval.Min {
				return fmt.Errorf("循环 %d 的提示 %d 的最大间隔必须大于或等于最小间隔", i+1, j+1)
			}

			if prompt.Type == "popup" && prompt.Popup == nil {
				return fmt.Errorf("循环 %d 的提示 %d 类型为弹窗但未提供弹窗配置", i+1, j+1)
			}

			if prompt.Type == "sound" && prompt.Sound == nil {
				return fmt.Errorf("循环 %d 的提示 %d 类型为声音但未提供声音配置", i+1, j+1)
			}

			if prompt.Type == "both" && (prompt.Popup == nil || prompt.Sound == nil) {
				return fmt.Errorf("循环 %d 的提示 %d 类型为两者但未提供完整配置", i+1, j+1)
			}
		}
	}

	return nil
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, configPath string) error {
	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 转换为YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 写入文件
	if err := ioutil.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}
