package model

// Standard 标准样本中的脂肪酸数据
type Standard struct {
	ID            int     // 记录ID
	Code          string  // 脂肪酸编码，如 C20:2
	Name          string  // 名称
	RetentionTime float64 // 保留时间
	Area          float64 // 峰面积
	AreaPct       float64 // 面积百分比
}
