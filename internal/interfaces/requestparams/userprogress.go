package requestparams

// GetUserProgressParams defines parameters for fetching user progress
type GetUserProgressParams struct {
	UserID   int `json:"user_id" valid:"required"`
	CourseID int `json:"course_id" valid:"required"`
}

// UpdateUserProgressParams defines parameters for updating user progress
type UpdateUserProgressParams struct {
	UserID            int  `json:"user_id"`
	CourseID          int  `json:"course_id" valid:"required"`
	ModulePosition    int  `json:"module_position" valid:"required"`
	ModuleItemPosition int  `json:"module_item_position" valid:"required"`
}
