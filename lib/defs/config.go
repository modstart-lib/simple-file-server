package defs

type Config struct {
	Debug bool `json:"debug"`

	Port     int    `json:"port"`
	ApiToken string `json:"apiToken"`

	TempDir string `json:"tempDir"`
	DataDir string `json:"dataDir"`
}
