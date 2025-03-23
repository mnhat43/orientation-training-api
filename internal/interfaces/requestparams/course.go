package requestparams

type CreateCourseParams struct {
	Title       string `json:"title" valid:"required"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	CreatedBy   int    `json:"created_by" valid:"required"`
}

type UpdateCourseParams struct {
	ID          int    `json:"course_id" valid:"required"`
	Title       string `json:"title" form:"course_title" valid:"required"`
	Description string `json:"description" form:"course_description"`
	Thumbnail   string `json:"thumbnail" form:"course_thumbnail"`
	CreatedBy   int    `json:"created_by" form:"created_by" valid:"required"`
}

type CourseIDParam struct {
	CourseID int `json:"course_id" valid:"required"`
}

type CourseListParams struct {
	CurrentPage int    `json:"current_page" valid:"-"`
	RowPerPage  int    `json:"row_per_page"`
	Keyword     string `json:"keyword"`
}

type UserCourseInfoRecords struct {
	UserId   int    `json:"user_id"`
	FullName string `json:"full_name"`
}
