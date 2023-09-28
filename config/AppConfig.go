package config

type AppConfig struct {
	Tasks   map[string]TaskConfig `yaml:"tasks"`
	Setting SettingConfig         `yaml:"setting"`
}

type TaskConfig struct {
	StartProcess string            `yaml:"process"`
	ProcessArgs  string            `yaml:"args"`
	WorkingDir   string            `yaml:"working_dir"`
	Environment  map[string]string `yaml:"environment"`
	Delay        int               `yaml:"delay"`
	LogMaxSize   int               `yaml:"logmaxsize"`
	LogInterval  int               `yaml:"loginterval"`
}

type SettingConfig struct {
	LogDir      string `yaml:"logdir"`
	LogMaxSize  int    `yaml:"logmaxsize"`
	LogInterval int    `yaml:"loginterval"`
}
