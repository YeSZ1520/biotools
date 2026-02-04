package config

type AppConfig struct {
	FattyAcid FattyAcidConfig `yaml:"fatty_acid" json:"fatty_acid"`
	Qpcr      QpcrConfig      `yaml:"qpcr" json:"qpcr"`
	Log       LogConfig       `yaml:"log" json:"log"`
}

type FattyAcidConfig struct {
	StandardFile     string  `yaml:"standard_file" json:"standard_file"`
	ExperimentalFile string  `yaml:"experimental_file" json:"experimental_file"`
	OutputFile       string  `yaml:"output_file" json:"output_file"`
	AreaThreshold    float64 `yaml:"area_threshold" json:"area_threshold"`
}

type QpcrConfig struct {
	ExperimentalFile       string `yaml:"experimental_file" json:"experimental_file"`
	OutputFile             string `yaml:"output_file" json:"output_file"`
	BaseLineGene           string `yaml:"baseline_gene" json:"baseline_gene"`
	BaseLineSamplePrefix   string `yaml:"baseline_sample_prefix" json:"baseline_sample_prefix"`
	DropCount              int    `yaml:"drop_count" json:"drop_count"`
}

type LogConfig struct {
	Level    string `yaml:"level" json:"level"`
	FilePath string `yaml:"file_path" json:"file_path"`
}
