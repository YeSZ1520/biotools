package service

import (
	"fmt"
	"strings"

	"github.com/YeSZ1520/biotools/internal/fatty-acid/model"
	"github.com/YeSZ1520/biotools/internal/utils"

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

	excelData,err:= utils.ReadTable(data_path,true)
	if err != nil {
		return nil, err
	}
	for _, rows := range excelData {
		for _, row := range rows {
			if len(row) >= 6 && utils.IsInteger(row[0]) {
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

	var experimentals []model.Experimental

	excelData, err:= utils.ReadTable(data_path,true)
	if err != nil {
		return nil, err
	}

	for sheetName, rows := range excelData {

		if len(rows) == 0 {
			logrus.Warnf("Sheet %s 清理后无数据，跳过", sheetName)
			continue
		}
		// 表头校验
		firstRow := rows[0]
		if len(firstRow)%4 != 0 {
			logrus.Warnf("Sheet %s 表头列数 %d 不正确，跳过", sheetName, len(firstRow))
			continue
		}
		// 解析数据行
		groupIndex := 1
		for i := 0; i < len(firstRow); i += 4 {
			for _, row := range rows {
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
