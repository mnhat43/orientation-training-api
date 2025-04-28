package requestparams

// GetUserProgressParams defines parameters for fetching user progress
type GetUserProgressParams struct {
	CourseID int `json:"course_id" valid:"required"`
}

// UpdateUserProgressParams defines parameters for updating user progress
type UpdateUserProgressParams struct {
	CourseID           int  `json:"course_id" valid:"required"`
	ModulePosition     int  `json:"module_position"`
	ModuleItemPosition int  `json:"module_item_position"`
	Completed          bool `json:"completed"`
}

// CreateUserProgressParams defines parameters for creating new user progress
type CreateUserProgressParams struct {
	CourseID           int  `json:"course_id" valid:"required"`
	ModulePosition     int  `json:"module_position" valid:"required"`
	ModuleItemPosition int  `json:"module_item_position" valid:"required"`
	Completed          bool `json:"completed"`
}

type GetListTraineeByCourseIDParams struct {
	CourseID    int    `json:"course_id" valid:"required"`
	CurrentPage int    `json:"current_page" valid:"-"`
	RowPerPage  int    `json:"row_per_page"`
	Keyword     string `json:"keyword"`
}
