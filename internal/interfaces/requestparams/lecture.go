package requestparams

type LectureListParams struct {
	CourseID int `json:"course_id" valid:"required"`
	// 	CurrentPage int    `json:"current_page" valid:"-"`
	// 	RowPerPage  int    `json:"row_per_page"`
	// 	Keyword     string `json:"keyword"`
}
