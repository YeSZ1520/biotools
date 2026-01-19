package model

// ExperimentalRaw qPCR 实验数据原始数据
type ExperimentalRaw struct {
	Gene     string  // 基因名称
	SampleID string  // 样本编号
	Well     string  // 孔位
	Ct       float64 // Ct值
}

// Experimental 孔位数据
type WellData struct {
	Well string  // 孔位
	Ct   float64 // Ct值
}

// Experimental qPCR 实验数据处理数据
type Experimental struct {
	Gene     string     // 基因名称
	SampleId string     // 样本编号
	WellData []WellData // 孔位数据
	MeanCt   float64    // 平均Ct值
}

type ExperimentalResult struct {
	BaseLine Experimental // 基线样本数据
	Target   Experimental // 目标样本数据
	Result   float64      // 结果值
	DeltaCt  float64	  // DeltaCt值
	Drop    int          // 是否为异常值，0否，1过大，-1过小
}
