package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var mu sync.Mutex

// WriteError writes an error response to the client.
func WriteError(w http.ResponseWriter, errMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorMessage := map[string]string{"error": errMsg}
	json.NewEncoder(w).Encode(errorMessage)
}

type Complaint

// GenerateUniqueID generates a unique ID.
func GenerateUniqueID(complaints map[string]Complaint) string {
	mu.Lock()
	defer mu.Unlock()
	return fmt.Sprintf("%d", len(complaints)+1)
}
