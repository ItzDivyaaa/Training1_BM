package complainBoxService

import (
	"complainBoxService/complainBox"
	"fmt"
	"net/http"
)

func main() {

	fmt.Println("Server is running on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}

	http.HandleFunc("/login", complainBox.LoginHandler)
	http.HandleFunc("/register", complainBox.RegisterHandler)
	http.HandleFunc("/submitComplaint", complainBox.SubmitComplaintHandler)
	http.HandleFunc("/getAllComplaintsForUser", complainBox.GetAllComplaintsForUserHandler)
	http.HandleFunc("/getAllComplaintsForAdmin", complainBox.GetAllComplaintsForAdminHandler)
	http.HandleFunc("/viewComplaint", complainBox.ViewComplaintHandler)
	http.HandleFunc("/resolveComplaint", complainBox.ResolveComplaintHandler)

}
