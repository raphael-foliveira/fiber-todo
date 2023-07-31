package common

type DatabaseConfig struct {
	Url string `json:"url"`
}

type AppConfig struct {
	Port int `json:"port"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	App      AppConfig      `json:"app"`
}
