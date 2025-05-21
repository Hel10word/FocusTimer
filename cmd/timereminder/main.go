package main

import (
	"flag"
	"fmt"
	"os"

	"FocusTimer/internal/app"
)

func main() {
	// 添加调试输出
	fmt.Println("程序启动...")

	// 命令行参数
	configPath := flag.String("config", "configs/default.yaml", "配置文件路径")
	flag.Parse()

	fmt.Printf("使用配置文件: %s\n", *configPath)

	// 创建服务
	service, err := app.NewService(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化服务错误: %v\n", err)
		// 暂停以便查看错误
		fmt.Println("按Enter键退出...")
		fmt.Scanln()
		os.Exit(1)
	}

	// 运行服务
	fmt.Println("服务开始运行...")
	if err := service.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "运行服务错误: %v\n", err)
		fmt.Println("按Enter键退出...")
		fmt.Scanln()
		os.Exit(1)
	}
}
