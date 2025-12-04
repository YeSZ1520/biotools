package config

type AppConfig struct {
	FattyAcid FattyAcidConfig `yaml:"fatty_acid" json:"fatty_acid"`
	Log       LogConfig       `yaml:"log" json:"log"`
}

type FattyAcidConfig struct {
	StandardFile     string  `yaml:"standard_file" json:"standard_file"`
	ExperimentalFile string  `yaml:"experimental_file" json:"experimental_file"`
	OutputFile       string  `yaml:"output_file" json:"output_file"`
	AreaThreshold    float64 `yaml:"area_threshold" json:"area_threshold"`
}

type LogConfig struct {
	Level    string `yaml:"level" json:"level"`
	FilePath string `yaml:"file_path" json:"file_path"`
}
