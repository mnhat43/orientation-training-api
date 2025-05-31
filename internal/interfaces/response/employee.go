package response

// EmployeeOverview represents the employee information structure
type EmployeeOverview struct {
	UserID      int    `json:"user_id"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Avatar      string `json:"avatar"`
	Department  string `json:"department"`
	Status      string `json:"status"`
}

const (
	StatusNotAssigned = "Not Assigned"
	StatusInProgress  = "In Progress"
	StatusCompleted   = "Completed"
)
