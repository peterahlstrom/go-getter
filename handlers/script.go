package script

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

type ScriptError struct {
	Message	string
	HttpStatus	int
}

func (e *ScriptError) Error() string {
	return e.Message
}

var scriptLimiter chan struct{}
func InitScriptLimiter(max int) {
	scriptLimiter = make(chan struct{}, max)
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