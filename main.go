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

type Endpoint struct {
	ScriptPath 	string `json:"scriptPath"`
	UrlPath 	string `json:"urlPath"`
}

type Config struct {
	LogPath					string `json:"logPath"`
	ConcurrentScriptsExec 	int `json:"concurrentScriptsExec"`
	Endpoints				[]Endpoint `json:"endpoints"`
}

type ScriptError struct {
	Message	string
	HttpStatus	int
}

func (e *ScriptError) Error() string {
	return e.Message
}

var scriptLimiter chan struct{}

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

	scriptLimiter = make(chan struct{}, cfg.ConcurrentScriptsExec)

	router := http.NewServeMux()
	
	for _, e := range cfg.Endpoints {
		router.HandleFunc(fmt.Sprintf("GET /%s", e.UrlPath), GetRequestHandler(e.ScriptPath))
	}

	addr := fmt.Sprintf(":%s", port)
	server := http.Server{
		Addr: addr,
		Handler: router,
	}
	log.Printf("INFO: Server starting. Listening to port %s\n", port)
	server.ListenAndServe()
}


func GetRequestHandler (scriptPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		result, err := RunScript(scriptPath)
		if err != nil {
			e := err.(*ScriptError)
			http.Error(w, e.Message, e.HttpStatus)
			log.Printf("ERROR: %s %s from %s - %d %s (%v)", 
				r.Method, r.URL, r.RemoteAddr, e.HttpStatus, http.StatusText(e.HttpStatus), time.Since(start))
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		
		w.Write(*result)
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


func RunScript (path string) (*[]byte, error) {
	select {
	case scriptLimiter <- struct{}{}:
		defer func() { <- scriptLimiter	}()
	case <- time.After(time.Second * 5):
		return nil, &ScriptError{Message: "Script server busy, try again later.",
		HttpStatus: http.StatusGatewayTimeout}
	}

	cmd := exec.Command(path)
		output, err := cmd.Output()
		if err != nil {
			return nil, &ScriptError{Message: fmt.Sprintf("The script resulted in an error: %v", err),
				HttpStatus: http.StatusInternalServerError,}
		}
		return &output, nil
}
