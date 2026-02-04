package main

import (
	"github.com/YeSZ1520/biotools/internal/config"
	"github.com/YeSZ1520/biotools/internal/qpcr/service"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		logrus.Fatalf("获取配置失败: %v", err)
	}

	data, err := service.LoadExperimentalData(cfg.Qpcr.ExperimentalFile)
	if err != nil {
		logrus.Fatalf("加载实验数据失败: %v", err)
	}

	formatData := service.FormatExperimentalData(data)

	results, err := service.Calculate(formatData, cfg.Qpcr.BaseLineGene, cfg.Qpcr.BaseLineSamplePrefix, cfg.Qpcr.DropCount)
	if err != nil {
		logrus.Fatalf("计算结果失败: %v", err)
	}

	if err := service.SaveData(cfg.Qpcr.OutputFile, results); err != nil {
		logrus.Fatalf("保存结果失败: %v", err)
	}
}
