package requestparams

// CreateTemplatePathParams defines parameters for creating a new template path
type CreateTemplatePathParams struct {
	Name        string `json:"name" valid:"required"`
	Description string `json:"description"`
	CourseIds   []int  `json:"course_ids" valid:"required"`
	Duration    int    `json:"duration"`
}

// UpdateTemplatePathParams defines parameters for updating an existing template path
type UpdateTemplatePathParams struct {
	TempPathID  int    `json:"id" valid:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CourseIds   []int  `json:"course_ids"`
	Duration    int    `json:"duration"`
}

// DeleteTemplatePathParams defines parameters for deleting a template path
type TempPathIDParam struct {
	TempPathID int `json:"id" valid:"required"`
}
