package notification

import (
	"fmt"
	"time"

	"FocusTimer/internal/config"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"os"
	"path/filepath"
	"strings"
)

// SoundNotifier 表示声音通知
type SoundNotifier struct {
	config   *config.SoundConfig
	streamer beep.StreamSeekCloser
	format   beep.Format
	done     chan bool
}

// NewSoundNotifier 创建声音通知器
func NewSoundNotifier(cfg *config.SoundConfig) Notifier {
	return &SoundNotifier{
		config: cfg,
		done:   make(chan bool),
	}
}

// Notify 播放声音
func (s *SoundNotifier) Notify() error {
	// 检查文件是否存在
	_, err := os.Stat(s.config.FilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("声音文件不存在: %s", s.config.FilePath)
	}

	// 打开文件
	f, err := os.Open(s.config.FilePath)
	if err != nil {
		return fmt.Errorf("无法打开声音文件: %v", err)
	}

	// 根据文件扩展名解码
	ext := strings.ToLower(filepath.Ext(s.config.FilePath))
	switch ext {
	case ".mp3":
		s.streamer, s.format, err = mp3.Decode(f)
	case ".wav":
		s.streamer, s.format, err = wav.Decode(f)
	default:
		f.Close()
		return fmt.Errorf("不支持的音频格式: %s", ext)
	}

	if err != nil {
		f.Close()
		return fmt.Errorf("解码音频文件失败: %v", err)
	}

	// 初始化音频播放器
	err = speaker.Init(s.format.SampleRate, s.format.SampleRate.N(time.Second/10))
	if err != nil {
		s.streamer.Close()
		return fmt.Errorf("初始化音频播放器失败: %v", err)
	}

	// 播放声音
	if s.config.Loop {
		// 循环播放
		loopStreamer := beep.Loop(-1, s.streamer)
		speaker.Play(loopStreamer)
	} else {
		// 单次播放
		speaker.Play(s.streamer)
	}

	// 如果设置了持续时间，计划停止
	if s.config.Duration > 0 {
		go func() {
			select {
			case <-time.After(s.config.Duration):
				s.Stop()
			case <-s.done:
				// 已经通过Stop()停止
			}
		}()
	}

	return nil
}

// Stop 停止声音
func (s *SoundNotifier) Stop() error {
	if s.streamer != nil {
		speaker.Clear()
		err := s.streamer.Close()
		s.streamer = nil
		close(s.done)
		return err
	}
	return nil
}
