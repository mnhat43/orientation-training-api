package requestparams

// GetUserProgressParams defines parameters for fetching user progress
type GetUserProgressParams struct {
	CourseID int `json:"course_id" valid:"required"`
	UserID   int `json:"user_id"`
}

// UpdateUserProgressParams defines parameters for updating user progress
type UpdateUserProgressParams struct {
	CourseID           int    `json:"course_id" valid:"required"`
	UserID             int    `json:"user_id"`
	CoursePosition     int    `json:"course_position"`
	ModulePosition     int    `json:"module_position"`
	ModuleItemPosition int    `json:"module_item_position"`
	Completed          bool   `json:"completed"`
	CompletedDate      string `json:"completed_date"`
}

// CreateUserProgressParams defines parameters for creating new user progress
type CreateUserProgressParams struct {
	UserID    int   `json:"user_id" valid:"required"`
	CourseIDs []int `json:"course_ids" valid:"required"`
}

// GetAllUserProgressParams defines parameters for fetching all user progress
type GetAllUserProgressParams struct {
	UserID int `json:"user_id"`
}

// GetSingleCourseProgressParams defines parameters for fetching a single course's progress
type GetSingleCourseProgressParams struct {
	CourseID int `json:"course_id" valid:"required"`
	UserID   int `json:"user_id"`
}

type GetListTraineeByCourseIDParams struct {
	CourseID    int    `json:"course_id" valid:"required"`
	CurrentPage int    `json:"current_page" valid:"-"`
	RowPerPage  int    `json:"row_per_page"`
	Keyword     string `json:"keyword"`
}
type AddListTraineeToCourseParams struct {
	CourseID int   `json:"course_id" validate:"required"`
	Trainees []int `json:"trainees" validate:"required"`
}

// ReviewProgressParams defines parameters for reviewing user progress
type ReviewProgressParams struct {
	UserID             int     `json:"user_id" valid:"required"`
	CourseID           int     `json:"course_id" valid:"required"`
	PerformanceRating  float64 `json:"performance_rating" valid:"required"`
	PerformanceComment string  `json:"performance_comment" valid:"required"`
}
