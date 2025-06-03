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

// EmployeeDetail represents the complete employee detail information structure
type EmployeeDetail struct {
	UserInfo     UserInfo     `json:"userInfo"`
	ProcessStats ProcessStats `json:"processStats"`
	ProcessInfo  []CourseInfo `json:"processInfo"`
}

// UserInfo represents basic information about the employee
type UserInfo struct {
	ID          int    `json:"id"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Department  string `json:"department"`
	JoinedDate  string `json:"joinedDate"`
	Avatar      string `json:"avatar"`
}

// ProcessStats represents the overall statistics of the employee's training process
type ProcessStats struct {
	CompletedCourses int     `json:"completedCourses"`
	TotalCourses     int     `json:"totalCourses"`
	TotalScore       float64 `json:"totalScore"`
	UserScore        float64 `json:"userScore"`
	CompletedDate    string  `json:"completedDate,omitempty"`
}

// CourseInfo represents detailed information about a course in the employee's training process
type CourseInfo struct {
	CourseID       int         `json:"course_id"`
	CourseTitle    string      `json:"course_title"`
	Status         string      `json:"status"`
	UserScore      float64     `json:"userScore,omitempty"`
	TotalScore     float64     `json:"totalScore,omitempty"`
	CompletedDate  string      `json:"completedDate,omitempty"`
	HasAssessment  bool        `json:"hasAssessment"`
	Assessment     *Assessment `json:"assessment,omitempty"`
	Progress       int         `json:"progress,omitempty"`
	CurrentModule  string      `json:"currentModule,omitempty"`
	PendingReviews int         `json:"pendingReviews"`
}

// Assessment represents the performance assessment for a completed course
type Assessment struct {
	PerformanceRating  float64 `json:"performance_rating"`
	PerformanceComment string  `json:"performance_comment"`
	ReviewerName       string  `json:"reviewer_name"`
}
