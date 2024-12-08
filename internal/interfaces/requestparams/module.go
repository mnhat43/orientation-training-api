package requestparams

type CreateModuleParams struct {
	Title    string `json:"title" form:"module_title" valid:"required"`
	CourseID int    `json:"course_id" form:"course_id" valid:"required"`
}

type ModuleListParams struct {
	CourseID    int    `json:"course_id" valid:"required"`
	CurrentPage int    `json:"current_page" valid:"-"`
	RowPerPage  int    `json:"row_per_page"`
	Keyword     string `json:"keyword"`
}

type ModuleIDParam struct {
	ModuleID int `json:"module_id" valid:"required"`
}
