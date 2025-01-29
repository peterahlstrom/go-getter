package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	ScriptPath 	string `json:"scriptPath"`
	UrlPath 	string `json:"urlPath"`
	Port 		string `json:"port"`
	LogPath 	string `json:"logPath"`
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

	router := http.NewServeMux()
	router.HandleFunc(fmt.Sprintf("GET /%s", cfg.UrlPath), GetRequestHandler(cfg))

	addr := fmt.Sprintf(":%s", port)
	server := http.Server{
		Addr: addr,
		Handler: router,
	}
	log.Printf("INFO: Server starting. Listening to port %s\n", port)
	server.ListenAndServe()
}


func GetRequestHandler (cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logFile, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)

		start := time.Now()
		
		result, err := RunScript(cfg.ScriptPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get data: %v", err), http.StatusInternalServerError)
			log.Printf("ERROR: %s %s from %s - 500 Internal Server Error (%v)", 
				r.Method, r.URL, r.RemoteAddr, time.Since(start))
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		
		w.Write(result)
		log.Printf("INFO:  %s %s from %s - 200 OK (%v)", r.Method, r.URL, r.RemoteAddr, time.Since(start))
	}
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


func RunScript (path string) ([]byte, error) {
	cmd := exec.Command(path)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("the script resulted in an error: %v", err)
	}
	
	return output, nil
}
