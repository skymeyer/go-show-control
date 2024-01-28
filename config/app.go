package config

type App struct {
	LogFile   string              `yaml:"log"`
	IODrivers map[string]IODriver `yaml:"io"`
}

type IODriver struct {
	Enabled bool   `yaml:"enabled"`
	Device  string `yaml:"device"`
}
