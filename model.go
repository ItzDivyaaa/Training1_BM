package model

// User represents a user record
type User struct {
	ID         string      `json:"id"`
	SecretCode string      `json:"secretCode"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Complaints []Complaint `json:"complaints"`
}

// Complaint represents a complaint record
type Complaint struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Severity int    `json:"severity"`
	Resolved bool   `json:"resolved"`
}

var users = make(map[string]User)
var complaints = make(map[string]Complaint)
