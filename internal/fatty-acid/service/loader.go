package service

import (
	"fmt"
	"strings"

	"github.com/YeSZ1520/biotools/internal/fatty-acid/model"
	"github.com/YeSZ1520/biotools/internal/fatty-acid/utils"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

func LoadStandardData(data_path string) ([]model.Standard, error) {
	if data_path == "" {
		return nil, fmt.Errorf("参考数据路径不能为空")
	}
	f, err := excelize.OpenFile(data_path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	logrus.Infof("找到 Sheet： %v", sheetList)
	var standards []model.Standard
	for _, sheetName := range sheetList {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			logrus.Errorf("获取行失败 %s: %v", sheetName, err)
			continue
		}

		for _, row := range rows {
			if len(row) >= 6 && row[0] != "" && utils.IsInteger(row[0]) {
				var item model.Standard
				fmt.Sscanf(row[0], "%d", &item.ID)
				item.Code = row[1]
				item.Name = row[2]
				fmt.Sscanf(row[3], "%f", &item.RetentionTime)
				fmt.Sscanf(row[4], "%f", &item.AreaPct)
				fmt.Sscanf(row[5], "%f", &item.Area)
				standards = append(standards, item)
			}
		}
	}

	return standards, nil
}

func LoadExperimentalData(data_path string) ([]model.Experimental, error) {
	if data_path == "" {
		return nil, fmt.Errorf("实验数据路径不能为空")
	}
	f, err := excelize.OpenFile(data_path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	var experimentals []model.Experimental

	for _, sheetName := range sheetList {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			logrus.Errorf("获取行失败 %s: %v", sheetName, err)
			continue
		}
		// 去除空行和空列
		cleanedRows := utils.CleanExcelData(rows)
		if len(cleanedRows) == 0 {
			logrus.Warnf("Sheet %s 清理后无数据，跳过", sheetName)
			continue
		}
		// 表头校验
		firstRow := cleanedRows[0]
		if len(firstRow)%4 != 0 {
			logrus.Warnf("Sheet %s 表头列数 %d 不正确，跳过", sheetName, len(firstRow))
			continue
		}
		headValid := true
		for i := 0; i < len(firstRow); i += 4 {
			if strings.TrimSpace(firstRow[i]) != "峰号" && strings.TrimSpace(firstRow[i+1]) != "保留时间" && strings.TrimSpace(firstRow[i+2]) != "峰面积" && strings.TrimSpace(firstRow[i+3]) != "面积百分比" {
				logrus.Warnf("Sheet %s 表头格式不正确，跳过", sheetName)
				headValid = false
				break
			}
		}
		if !headValid {
			continue
		}
		// 解析数据行
		groupIndex := 1
		for i := 0; i < len(firstRow); i += 4 {
			for _, row := range cleanedRows {
				id := row[i]
				if id == "" || strings.TrimSpace(id) == "总计" || !utils.IsInteger(id) {
					continue
				}
				var item model.Experimental
				item.Sheet = sheetName
				item.Group = groupIndex
				fmt.Sscanf(id, "%d", &item.ID)
				fmt.Sscanf(row[i+1], "%f", &item.RetentionTime)
				fmt.Sscanf(row[i+2], "%f", &item.AreaPct)
				fmt.Sscanf(row[i+3], "%f", &item.Area)
				experimentals = append(experimentals, item)
			}
			groupIndex++
		}

	}
	return experimentals, nil
}
