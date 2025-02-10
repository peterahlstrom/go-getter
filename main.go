package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/peterahlstrom/go-getter/auth"
	"github.com/peterahlstrom/go-getter/config"
	"github.com/peterahlstrom/go-getter/handlers/script"
)


var configPath = "config.json"

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: ./main <port>")
	}
	port := os.Args[1]

	cfg, err := config.GetConfig(configPath)
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
	
	for k, e := range cfg.Endpoints {
		router.HandleFunc(fmt.Sprintf("GET %s", k), script.GetRequestHandler(e.ScriptPath))
	}

	secureHandler := auth.ApiKeyMiddleWare(cfg.Endpoints)(router)

	addr := fmt.Sprintf(":%s", port)
	server := http.Server{
		Addr: addr,
		Handler: secureHandler,
	}
	fmt.Printf("INFO: Server starting. Listening to port %s\n", port)
	server.ListenAndServe()
}
