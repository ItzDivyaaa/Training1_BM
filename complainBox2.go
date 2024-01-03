// complainBox/handlers.go

package complainBox

import (
	"complainBoxService/common"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sync"
)

var mu sync.Mutex

// Define your common functions here
func writeError(w http.ResponseWriter, errMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorMessage := map[string]string{"error": errMsg}
	json.NewEncoder(w).Encode(errorMessage)
}

func generateUniqueID() string {
	return fmt.Sprintf("%d", len(complaints)+1)
}

// Define your handlers here
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var credentials struct {
		SecretCode string `json:"secretCode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, exists := users[credentials.SecretCode]
	if !exists {
		writeError(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is POST, return a Method Not Allowed response otherwise
	if r.Method != http.MethodPost {
		writeError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var newUser User

	// Decode the JSON data from the request body into the newUser variable
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate email format
	if !isValidEmail(newUser.Email) {
		writeError(w, common.Str1, http.StatusBadRequest)
		return
	}

	// Validate name format (you can customize this validation as needed)
	if len(newUser.Name) < 3 {
		writeError(w, common.Str2, http.StatusBadRequest)
		return
	}

	if _, exists := users[newUser.SecretCode]; exists {
		writeError(w, "Secret code already in use", http.StatusBadRequest)
		return
	}

	newUser.ID = generateUniqueID()
	newUser.Complaints = []Complaint{}

	users[newUser.SecretCode] = newUser

	json.NewEncoder(w).Encode(newUser)
}

// isValidEmail checks if the given email is in a valid format
func isValidEmail(email string) bool {
	// You can implement a more sophisticated email validation logic if needed
	// This is a simple example
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func SubmitComplaintHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var newComplaint Complaint

	if err := json.NewDecoder(r.Body).Decode(&newComplaint); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, exists := users[newComplaint.SecretCode]
	if !exists {
		writeError(w, "User not found", http.StatusNotFound)
		return
	}

	newComplaint.ID = generateUniqueID()

	complaints[newComplaint.ID] = newComplaint

	user.Complaints = append(user.Complaints, newComplaint)
	users[newComplaint.SecretCode] = user

	w.WriteHeader(http.StatusCreated)
}

func GetAllComplaintsForUserHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var user struct {
		SecretCode string `json:"secretCode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	userDetails, exists := users[user.SecretCode]
	if !exists {
		writeError(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(userDetails.Complaints)
}

func GetAllComplaintsForAdminHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var adminCredentials struct {
		SecretCode string `json:"secretCode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&adminCredentials); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if adminCredentials.SecretCode != "admin" {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var allComplaints []Complaint
	for _, user := range users {
		allComplaints = append(allComplaints, user.Complaints...)
	}

	json.NewEncoder(w).Encode(allComplaints)
}

func ViewComplaintHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var complaint struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&complaint); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	complaintDetails, exists := complaints[complaint.ID]
	if !exists {
		writeError(w, "Complaint not found", http.StatusNotFound)
		return
	}

	_, exists = users[complaintDetails.SecretCode]
	if !exists {
		writeError(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(complaintDetails)
}

func ResolveComplaintHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var complaint struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&complaint); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var adminCredentials struct {
		SecretCode string `json:"secretCode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&adminCredentials); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if adminCredentials.SecretCode != "admin" {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	complaints :=
	complaintDetails, exists := complaints[complaint.ID]
	if !exists {
		writeError(w, "Complaint not found", http.StatusNotFound)
		return
	}

	complaintDetails.Resolved = true
	complaints[complaint.ID] = complaintDetails

	w.WriteHeader(http.StatusNoContent)
}
