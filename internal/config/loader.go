package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const defaultConfigFile = "config.yaml"

var (
	configInstance *AppConfig
	configPath     string
)

func ReloadConfig(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	if err := config.validate(); err != nil {
		return err
	}
	config.setDefaults()
	config.Log.initLogger()
	configInstance = &config
	configPath = filePath
	return nil
}

func GetConfig() (*AppConfig, error) {
	if configInstance != nil {
		return configInstance, nil
	}
	if configPath == "" {
		configPath = defaultConfigFile
	}
	if err := ReloadConfig(configPath); err != nil {
		return nil, err
	}
	configJson, err := json.MarshalIndent(configInstance, "", "  ")
	if err == nil {
		logrus.Infof("配置加载成功: %s", string(configJson))
	}
	return configInstance, nil
}

func (c *AppConfig) validate() error {
	// 验证FattyAcidConfig
	if c.FattyAcid.StandardFile == "" {
		return fmt.Errorf("参考数据文件路径不能为空")
	}
	if c.FattyAcid.ExperimentalFile == "" {
		return fmt.Errorf("实验数据文件路径不能为空")
	}
	if c.FattyAcid.OutputFile == "" {
		return fmt.Errorf("输出文件路径不能为空")
	}
	return nil
}

func (c *AppConfig) setDefaults() {
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
}

func (c *LogConfig) initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	logPath := c.FilePath
	if logPath == "" {
		logrus.SetOutput(os.Stdout)
		logrus.Info("日志文件路径未设置，使用标准输出")
		return
	}

	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		logrus.SetOutput(os.Stdout)
		logrus.Errorf("创建日志目录失败: %v", err)
		return
	}

	if !filepath.IsAbs(logPath) {
		logPath = filepath.Join(logsDir, logPath)
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logrus.SetOutput(os.Stdout)
		logrus.Errorf("打开日志文件失败: %v", err)
	} else {
		logrus.SetOutput(io.MultiWriter(os.Stdout, f))
	}

	if lvl, err := logrus.ParseLevel(c.Level); err == nil {
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetReportCaller(true)
}
