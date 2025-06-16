package requestparams

type CreateCourseParams struct {
	Title           string `json:"title" valid:"required"`
	Description     string `json:"description"`
	Thumbnail       string `json:"thumbnail"`
	Category        string `json:"category" valid:"required"`
	CreatedBy       int    `json:"created_by"`
	SkillKeywordIDs []int  `json:"skill_keyword_ids"`
}

type UpdateCourseParams struct {
	ID              int    `json:"course_id" valid:"required"`
	Title           string `json:"title" form:"course_title"`
	Description     string `json:"description" form:"course_description"`
	Thumbnail       string `json:"thumbnail" form:"course_thumbnail"`
	Category        string `json:"category" form:"course_category"`
	SkillKeywordIDs []int  `json:"skill_keyword_ids"`
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
