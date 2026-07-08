package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Backup(cfg *Config) {
	// 1. 获取当前日期作为子目录名 (例如: 2026-01-30)
	dateDir := time.Now().Format("2006-01-02")
	targetPath := filepath.Join(cfg.BackupPath, dateDir)

	// 创建当天的备份文件夹
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		fmt.Println("创建日期备份目录失败:", err)
		return
	}

	// 2. 执行备份
	for _, db := range cfg.Databases {
		backupDatabase(cfg, db.Name, targetPath)
	}

	// 3. 清理指定天数前的旧备份
	cleanOldBackups(cfg.BackupPath, cfg.Clear)
}

func backupDatabase(cfg *Config, dbName string, targetPath string) {
	timestamp := time.Now().Format("150405") // 文件夹已按日期分类，文件名只需保留时间戳
	fileName := fmt.Sprintf("%s_%s.sql", dbName, timestamp)
	filePath := filepath.Join(targetPath, fileName)

	args := []string{
		"-h", cfg.MySQL.Host,
		"-P", fmt.Sprintf("%d", cfg.MySQL.Port),
		"-u", cfg.MySQL.User,
		fmt.Sprintf("-p%s", cfg.MySQL.Password),
		"--ssl=0",
		"--single-transaction",
		"--flush-privileges",
		"--quick",
		"--routines",
		"--events",
		dbName,
	}

	cmd := exec.Command("mariadb-dump", args...)

	outFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("创建备份文件失败:", err)
		return
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	fmt.Printf("[%s] 开始备份数据库: %s\n", time.Now().Format("15:04:05"), dbName)
	if err := cmd.Run(); err != nil {
		fmt.Println("备份失败:", dbName, err)
		return
	}
	fmt.Printf("[%s] 备份数据库成功...: %s\n", time.Now().Format("15:04:05"), filePath)
}

// cleanOldBackups 用于删除指定目录下超过指定天数的文件或目录(根据mod时间而不是文件夹名字)
func cleanOldBackups(backupRoot string, clear int) {
	fmt.Println("检查并清理过期备份...")
	files, err := os.ReadDir(backupRoot)
	if err != nil {
		fmt.Println("读取备份根目录失败:", err)
		return
	}

	now := time.Now()
	clearDays := now.AddDate(0, 0, -clear)

	for _, file := range files {
		fullPath := filepath.Join(backupRoot, file.Name())
		info, err := file.Info()
		if err != nil {
			continue
		}

		// 如果文件的修改时间早于 clear 天前，则删除
		if info.ModTime().Before(clearDays) {
			err := os.RemoveAll(fullPath) // 使用 RemoveAll 可以递归删除文件夹
			if err != nil {
				fmt.Printf("删除过期备份失败 [%s]: %v\n", fullPath, err)
			} else {
				fmt.Printf("已成功清理过期备份: %s\n", fullPath)
			}
		}
	}
}
