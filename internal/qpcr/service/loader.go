package service

import (
	"fmt"

	"github.com/YeSZ1520/biotools/internal/qpcr/model"
	"github.com/YeSZ1520/biotools/internal/utils"
)

func LoadExperimentalData(data_path string) ([]model.ExperimentalRaw, error) {
	if data_path == "" {
		return nil, fmt.Errorf("实验数据路径不能为空")
	}

	var experimentals []model.ExperimentalRaw
	excelData, err := utils.ReadTable(data_path, true)
	if err != nil {
		return nil, err
	}
	for _, rows := range excelData {
		for _, row := range rows {
			if len(row) >= 4 && row[0] != "" && row[1] != "" && row[2] != "" && row[3] != "" {
				var item model.ExperimentalRaw
				item.Gene = row[0]
				item.SampleID = row[1]
				item.Well = row[2]
				fmt.Sscanf(row[3], "%f", &item.Ct)
				experimentals = append(experimentals, item)
			}
		}
	}
	return experimentals, nil
}

func FormatExperimentalData(raws []model.ExperimentalRaw) map[string]map[string]*model.Experimental {
	data := make(map[string]map[string]*model.Experimental)
	for _, raw := range raws {
		if _, ok := data[raw.Gene]; !ok {
			data[raw.Gene] = make(map[string]*model.Experimental)
		}
		if _, ok := data[raw.Gene][raw.SampleID]; !ok {
			data[raw.Gene][raw.SampleID] = &model.Experimental{
				Gene:     raw.Gene,
				SampleId: raw.SampleID,
				WellData: []model.WellData{},
			}
		}
		data[raw.Gene][raw.SampleID].WellData = append(data[raw.Gene][raw.SampleID].WellData, model.WellData{
			Well: raw.Well,
			Ct:   raw.Ct,
		})
	}
	for _, samples := range data {
		for _, exp := range samples {
			var sum float64
			count := 0
			for _, well := range exp.WellData {
				sum += well.Ct
				count++
			}
			if count > 0 {
				exp.MeanCt = sum / float64(count)
			}
		}
	}
	return data
}
