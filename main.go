package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	script "github.com/peterahlstrom/go-getter/handlers"
)

type Endpoint struct {
	ScriptPath 	string `json:"scriptPath"`
	UrlPath 	string `json:"urlPath"`
}

type Config struct {
	LogPath					string `json:"logPath"`
	ConcurrentScriptsLimit 	int `json:"concurrentScriptsLimit"`
	Endpoints				[]Endpoint `json:"endpoints"`
}

var configPath = "config.json"

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: ./main <port>")
	}
	port := os.Args[1]

	cfg, err := GetConfig(configPath)
	if err != nil {
		log.Fatalf("ERROR: Config file %s: %v\n", configPath, err)
	}

	logFile, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	script.InitScriptLimiter(cfg.ConcurrentScriptsLimit)

	router := http.NewServeMux()
	
	for _, e := range cfg.Endpoints {
		router.HandleFunc(fmt.Sprintf("GET /%s", e.UrlPath), script.GetRequestHandler(e.ScriptPath))
	}

	addr := fmt.Sprintf(":%s", port)
	server := http.Server{
		Addr: addr,
		Handler: router,
	}
	log.Printf("INFO: Server starting. Listening to port %s\n", port)
	server.ListenAndServe()
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
