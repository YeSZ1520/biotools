package main

import (
	"fmt"

	"github.com/YeSZ1520/biotools/internal/config"
	"github.com/YeSZ1520/biotools/internal/fatty-acid/service"

	"github.com/sirupsen/logrus"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		logrus.Fatalf("获取配置失败: %v", err)
	}

	data := service.CompareFattyAcid(config.FattyAcid.StandardFile, config.FattyAcid.ExperimentalFile, config.FattyAcid.OutputFile, config.FattyAcid.AreaThreshold)
	err = service.WriteExperimentalDataToExcel(data, config.FattyAcid.OutputFile)
	if err != nil {
		logrus.Fatalf("写入Excel文件失败: %v", err)
	}
	fmt.Println("数据已成功写入 " + config.FattyAcid.OutputFile)
	Wait()

}
func Wait() {
	fmt.Print("输入 Enter 继续...")
	fmt.Scanln()
}
