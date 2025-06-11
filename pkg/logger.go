package remotelist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var logMu sync.Mutex

func logOperation(op, detail string) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Operation: op,
		Details:   detail,
	}

	logMu.Lock()
	defer logMu.Unlock()

	logPath := filepath.Join("./src", "log.json")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		println("[ERROR] Failed to write log:", err.Error())
		return
	}
	defer f.Close()

	data, _ := json.Marshal(entry)
	f.Write(append(data, '\n'))
}
