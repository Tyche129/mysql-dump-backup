package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Databases  []Database `yaml:"databases"`
	BackupPath string     `yaml:"backup_path"`
	MySQL      MySQL      `yaml:"mysql"`
	Cron       string     `yaml:"cron"`
	Des        string     `yaml:"des"`
	Clear      int        `yaml:"clear"`
	Init       bool       `yaml:"init"`
}

type Database struct {
	Name string `yaml:"name"`
}

type MySQL struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
