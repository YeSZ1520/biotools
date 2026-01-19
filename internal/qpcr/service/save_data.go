package service

import (
	"fmt"

	"github.com/YeSZ1520/biotools/internal/qpcr/model"
	"github.com/xuri/excelize/v2"
)

func SaveData(filename string, data map[string][]*model.ExperimentalResult) error {
	f := excelize.NewFile()

	// 创建黄色背景样式
	yellowStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFFF00"}, Pattern: 1},
	})
	if err != nil {
		return err
	}

	firstSheet := true
	for gene, results := range data {
		sheetName := gene
		// Excel sheet name max length is 31
		if len(sheetName) > 31 {
			sheetName = sheetName[:31]
		}

		if firstSheet {
			f.SetSheetName("Sheet1", sheetName)
			firstSheet = false
		} else {
			f.NewSheet(sheetName)
		}

		// 设置表头
		headers := []string{"样品组", "孔位", "基线Ct", "基线平均Ct", "", "样本孔位", "样本Ct", "样本平均Ct", "ΔCt", "结果"}
		for i, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, header)
		}

		currentRow := 2
		for _, res := range results {
			baseLen := len(res.BaseLine.WellData)
			targetLen := len(res.Target.WellData)
			maxRows := baseLen
			if targetLen > maxRows {
				maxRows = targetLen
			}
			if maxRows == 0 {
				maxRows = 1 // 至少占一行
			}

			// 写入合并列的数据 (A, D, H, I, J)
			// A: BaseLine.SampleId
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), res.BaseLine.SampleId)
			// D: BaseLine.MeanCt
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), res.BaseLine.MeanCt)
			// H: Target.MeanCt
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", currentRow), res.Target.MeanCt)
			// I: DeltaCt
			cellI := fmt.Sprintf("I%d", currentRow)
			f.SetCellValue(sheetName, cellI, res.DeltaCt)
			if res.Drop != 0 {
				f.SetCellStyle(sheetName, cellI, cellI, yellowStyle)
			}
			// J: Result
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", currentRow), res.Result)

			// 写入多行数据 (B, C, F, G)
			for i := 0; i < maxRows; i++ {
				row := currentRow + i
				// BaseLine WellData
				if i < baseLen {
					f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), res.BaseLine.WellData[i].Well)
					f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), res.BaseLine.WellData[i].Ct)
				}
				// Target WellData
				if i < targetLen {
					f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), res.Target.WellData[i].Well)
					f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), res.Target.WellData[i].Ct)
				}
			}

			// 合并单元格
			if maxRows > 1 {
				mergeCols := []string{"A", "D", "H", "I", "J"}
				for _, col := range mergeCols {
					startCell := fmt.Sprintf("%s%d", col, currentRow)
					endCell := fmt.Sprintf("%s%d", col, currentRow+maxRows-1)
					f.MergeCell(sheetName, startCell, endCell)
				}
			}

			currentRow += maxRows
		}
	}

	if err := f.SaveAs(filename); err != nil {
		return err
	}
	return nil
}
