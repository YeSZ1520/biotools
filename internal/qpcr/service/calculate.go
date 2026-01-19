package service

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/YeSZ1520/biotools/internal/qpcr/model"
	"github.com/sirupsen/logrus"
)

func Calculate(data map[string]map[string]*model.Experimental, baselineGene string, baseLineSamplePrefix string, dropCount int) (map[string][]*model.ExperimentalResult, error) {
	results := make(map[string][]*model.ExperimentalResult)
	baselineSamples, ok := data[baselineGene]
	if !ok {
		return nil, fmt.Errorf("基线基因 %s 不存在于数据中", baselineGene)
	}

	for gene, samples := range data {
		if gene == baselineGene {
			continue
		}

		var list []*model.ExperimentalResult

		// 收集共同存在的样本并计算 deltaCt
		for sampleID, targetExp := range samples {
			if baselineExp, exists := baselineSamples[sampleID]; exists {
				deltaCt := targetExp.MeanCt - baselineExp.MeanCt
				list = append(list, &model.ExperimentalResult{
					BaseLine: *baselineExp,
					Target:   *targetExp,
					DeltaCt:  deltaCt,
				})
			}
		}

		if len(list) == 0 {
			continue
		}

		meanDeltaCt, err := calculateBaseLineSampleCt(list, baseLineSamplePrefix, dropCount)
		if err != nil {
			logrus.Warnf("计算基线样本组 %s 均值 DeltaCt 失败: %v，跳过基因 %s", baseLineSamplePrefix, err, gene)
			continue
		}

		// 生成结果
		for i := range list {
			deltaDeltaCt := list[i].DeltaCt - meanDeltaCt
			relativeExpression := math.Pow(2, -deltaDeltaCt)
			list[i].Result = relativeExpression
		}
		sort.Slice(list, func(i, j int) bool {
			iHasPrefix := strings.HasPrefix(list[i].Target.SampleId, baseLineSamplePrefix)
			jHasPrefix := strings.HasPrefix(list[j].Target.SampleId, baseLineSamplePrefix)
			if iHasPrefix == jHasPrefix {
				return list[i].Target.SampleId < list[j].Target.SampleId
			}
			return iHasPrefix && !jHasPrefix
		})
		results[gene] = list
	}

	return results, nil
}


func calculateBaseLineSampleCt(list []*model.ExperimentalResult, baseLineSamplePrefix string, dropCount int) (float64, error) {
	var tagrgetSampleGroup []*model.ExperimentalResult

	for _, item := range list {
		if strings.HasPrefix(item.Target.SampleId, baseLineSamplePrefix) {
			tagrgetSampleGroup = append(tagrgetSampleGroup, item)
		}
	}
	if len(tagrgetSampleGroup) == 0 {
		return 0, fmt.Errorf("基线样本组 %s 在数据中不存在", baseLineSamplePrefix)
	}
	if len(tagrgetSampleGroup) <= 2*dropCount {
		return 0, fmt.Errorf("基线样本组 %s 样本数 %d 不足以删除 %d 个异常值", baseLineSamplePrefix, len(tagrgetSampleGroup), dropCount)
	}
	sort.Slice(tagrgetSampleGroup, func(i, j int) bool {
		return tagrgetSampleGroup[i].DeltaCt < tagrgetSampleGroup[j].DeltaCt
	})
	
	// 标记异常值
	for i := 0; i < dropCount; i++ {
		tagrgetSampleGroup[i].Drop = -1
	}
	for i := 0; i < dropCount; i++ {
		tagrgetSampleGroup[len(tagrgetSampleGroup)-1-i].Drop = 1
	}

	// 计算均值
	tagrgetSampleGroup = tagrgetSampleGroup[dropCount : len(tagrgetSampleGroup)-dropCount]
	var sumDeltaCt float64
	for _, v := range tagrgetSampleGroup {
		sumDeltaCt += v.DeltaCt
	}
	meanDeltaCt := sumDeltaCt / float64(len(tagrgetSampleGroup))

	return meanDeltaCt, nil
}

