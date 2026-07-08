package main

import (
	"fmt"
	"os"
	"os/signal"
	_ "runtime"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yml"
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	// 1. 初始化 Cron
	c := cron.New(cron.WithSeconds())

	_, err = c.AddFunc(cfg.Cron, func() {
		Backup(cfg)
	})
	if err != nil {
		panic(err)
	}

	// 2. 启动 Cron
	c.Start()
	fmt.Printf("🚀 服务启动成功！当前 Cron 表达式: [%s]\n", cfg.Cron)
	fmt.Printf("当前时间: [%s]\n", time.Now())
	fmt.Println(cfg.Des)
	if cfg.Init == true {
		fmt.Println("您已设置首次启动备份功能。")
		Backup(cfg)
	}
	// 3. 设置信号监听以实现优雅关闭
	// 创建一个通道来接收系统信号
	quit := make(chan os.Signal, 1)
	// 监听中断 (Ctrl+C) 和 终止 (kill) 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞在这里，直到收到信号
	<-quit

	fmt.Println("\n正在关闭服务...")

	// 4. 停止 Cron 任务（这会等待当前正在运行的任务执行完毕）
	ctx := c.Stop()
	<-ctx.Done() // 等待停止确认

	fmt.Println("✅ 服务已安全关闭。")
}
