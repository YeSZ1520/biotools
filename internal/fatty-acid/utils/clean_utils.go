package utils

import "strings"

// CleanExcelData 清理Excel数据，去除空行和空列，并填充空字符串
func CleanExcelData(rows [][]string) [][]string {
    if len(rows) == 0 {
        return rows
    }
    
    // 找出哪些列有数据
    colHasData := make(map[int]bool)
    var nonEmptyRows [][]string
    
    // 第一步：过滤空行，标记有效列
    for _, row := range rows {
        isEmpty := true
        for j, cell := range row {
            if strings.TrimSpace(cell) != "" {
                isEmpty = false
                colHasData[j] = true
            }
        }
        
        if !isEmpty {
            nonEmptyRows = append(nonEmptyRows, row)
        }
    }
    
    // 如果没有有效数据，返回空
    if len(nonEmptyRows) == 0 {
        return [][]string{}
    }
    
    // 找出最大列索引
    maxCol := 0
    for j := range colHasData {
        if j > maxCol {
            maxCol = j
        }
    }
    
    // 第二步：填充空字符串，移除空列
    var result [][]string
    for _, row := range nonEmptyRows {
        newRow := []string{}
        for j := 0; j <= maxCol; j++ {
            if colHasData[j] {
                // 填充空字符串
                if j < len(row) {
                    newRow = append(newRow, row[j])
                } else {
                    newRow = append(newRow, "") // 填充空字符串
                }
            }
        }
        result = append(result, newRow)
    }
    
    return result
}