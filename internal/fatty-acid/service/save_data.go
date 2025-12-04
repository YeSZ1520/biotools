package service

import (
	"fmt"

	"github.com/YeSZ1520/biotools/internal/fatty-acid/model"

	"github.com/xuri/excelize/v2"
)

func WriteExperimentalDataToExcel(data []model.Experimental, filename string) error {
	f := excelize.NewFile()

	sheetMap := make(map[string]bool)

	// 定义表头
	headers := []string{
		"组别", "峰号", "保留时间", "峰面积",
		"面积百分比", "名称", "相对偏差",
	}

	// 按Sheet名称分组数据
	groupedData := make(map[string][]model.Experimental)
	for _, item := range data {
		groupedData[item.Sheet] = append(groupedData[item.Sheet], item)
	}

	// 写入数据
	for sheetName, items := range groupedData {
		// 创建新Sheet
		_, err := f.NewSheet(sheetName)
		if err != nil {
			return fmt.Errorf("创建Sheet %s 失败: %v", sheetName, err)
		}
		sheetMap[sheetName] = true

		// 写入表头
		for col, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			f.SetCellValue(sheetName, cell, header)
		}

		// 设置表头样式
		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#D3D3D3"},
				Pattern: 1,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})

		headerRange, _ := excelize.CoordinatesToCellName(1, 1)
		lastCol, _ := excelize.ColumnNumberToName(len(headers))
		f.SetCellStyle(sheetName, headerRange, lastCol+"1", style)

		// 写入数据行
		for row, item := range items {
			rowIndex := row + 2

			// 设置每列的值
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), item.Group)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), item.ID)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), item.RetentionTime)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), item.Area)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), item.AreaPct)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIndex), item.Name)
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIndex), item.RetentionTimeDeviation)
		}
		// 列宽
		f.SetColWidth(sheetName, "A", "B", 5)
		f.SetColWidth(sheetName, "C", "E", 12)
		f.SetColWidth(sheetName, "F", "G", 20)

	}

	if !sheetMap["Sheet1"] {
		f.DeleteSheet("Sheet1")
	}

	// 保存文件
	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}
