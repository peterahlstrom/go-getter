package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Endpoint struct {
	ScriptPath 		string `json:"scriptPath"`
	RequireAuth		bool `json:"requireAuth"`
	ValidApiKeys	map[string]string `json:"apiKeys"`
}

type Config struct {
	LogPath					string `json:"logPath"`
	ConcurrentScriptsLimit 	int `json:"concurrentScriptsLimit"`
	Endpoints				map[string]Endpoint `json:"endpoints"`
}


func GetConfig (configPath string) (*Config, error) {
	var cfg Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse config-file: %v", err)
	}
	
	return &cfg, nil
} 
