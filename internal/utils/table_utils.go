package utils

import (

	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)


func ReadTable(filename string, cleanTable bool) (map[string][][]string, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("无法打开Excel文件: %v", err)
	}
	defer f.Close()

	result := make(map[string][][]string)
	sheetList := f.GetSheetList()

	for _, sheetName := range sheetList {
		mergedCells, err := f.GetMergeCells(sheetName)
		if err != nil {
			log.Printf("无法获取sheet %s 的合并单元格: %v\n", sheetName, err)
			continue
		}

		rows, err := f.GetRows(sheetName)
		if err != nil {
			log.Printf("无法获取sheet %s 的行数据: %v\n", sheetName, err)
			continue
		}

		for _, merge := range mergedCells {
			val, err := f.GetCellValue(sheetName, merge.GetStartAxis())
			if err != nil {
				log.Println(err)
				continue
			}

			startCol, startRow, _ := excelize.CellNameToCoordinates(merge.GetStartAxis())
			endCol, endRow, _ := excelize.CellNameToCoordinates(merge.GetEndAxis())

			for r := startRow - 1; r < endRow; r++ {
				for c := startCol - 1; c < endCol; c++ {
					for len(rows) <= r {
						rows = append(rows, []string{})
					}
					for len(rows[r]) <= c {
						rows[r] = append(rows[r], "")
					}
					rows[r][c] = val
				}
			}
		}

		if cleanTable {
			rows = CleanExcelData(rows)
		}
		result[sheetName] = rows
	}

	return result, nil
}