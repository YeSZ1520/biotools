package model

// Experimental 实验样本中的脂肪酸数据
type Experimental struct {
	Sheet                  string  // 表单名称
	Group                  int     // 样本组别
	ID                     int     // 峰号
	RetentionTime          float64 // 保留时间
	Area                   float64 // 峰面积
	AreaPct                float64 // 面积百分比
	Name                   string  // 名称
	RetentionTimeDeviation float64 // 保留时间相对偏差
}
