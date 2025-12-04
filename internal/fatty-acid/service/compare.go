package service

import (
	"fmt"
	"math"

	"os"

	"github.com/YeSZ1520/biotools/internal/fatty-acid/model"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

func CompareFattyAcid(standardDataPath, experimentalDataPath, outputFile string, threshold float64) []model.Experimental {
	standards, err := LoadStandardData(standardDataPath)
	if err != nil {
		logrus.Fatalf("加载标准数据集失败: %v", err)
	}
	printStandardsTable(standards)
	experimentals, err := LoadExperimentalData(experimentalDataPath)
	if err != nil {
		logrus.Fatalf("加载实验数据失败: %v", err)
	}
	for i, exp := range experimentals {
		if exp.AreaPct < threshold {
			continue
		}
		closestStandard := getClosestStandard(exp, standards)
		experimentals[i].Name = closestStandard.Name
		experimentals[i].RetentionTimeDeviation = math.Abs((exp.RetentionTime - closestStandard.RetentionTime) / closestStandard.RetentionTime * 100)
	}
	printExperimentalsTable(experimentals)
	return experimentals
}

func getClosestStandard(exp model.Experimental, standards []model.Standard) model.Standard {
	var closest model.Standard
	minDiff := 1e9

	for _, std := range standards {
		diff := math.Abs(exp.RetentionTime - std.RetentionTime)
		if diff < minDiff {
			minDiff = diff
			closest = std
		}
	}
	return closest
}

func printStandardsTable(standards []model.Standard) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"峰号", "编码", "名称", "保留时间", "面积", "面积百分比"})

	for _, s := range standards {
		row := []string{
			fmt.Sprintf("%d", s.ID),
			s.Code,
			s.Name,
			fmt.Sprintf("%.3f", s.RetentionTime),
			fmt.Sprintf("%.2f", s.Area),
			fmt.Sprintf("%.2f", s.AreaPct),
		}
		table.Append(row)
	}
	table.Render()
}

func printExperimentalsTable(experimentals []model.Experimental) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"表单", "组别", "峰号", "名称", "保留时间", "面积", "面积百分比", "相对偏差"})
	for _, e := range experimentals {
		row := []string{
			e.Sheet,
			fmt.Sprintf("%d", e.Group),
			fmt.Sprintf("%d", e.ID),
			e.Name,
			fmt.Sprintf("%.3f", e.RetentionTime),
			fmt.Sprintf("%.2f", e.Area),
			fmt.Sprintf("%.2f", e.AreaPct),
			fmt.Sprintf("%.2f", e.RetentionTimeDeviation),
		}
		table.Append(row)
	}
	table.Render()
}
