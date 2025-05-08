package requestparams

// GetTemplatePathParams defines parameters for retrieving a template path
type GetTemplatePathParams struct {
	TempPathID int `json:"id" valid:"required"`
}

// CreateTemplatePathParams defines parameters for creating a new template path
type CreateTemplatePathParams struct {
	Name        string `json:"name" valid:"required"`
	Description string `json:"description"`
	Courses     []int  `json:"courses" valid:"required"`
}

// UpdateTemplatePathParams defines parameters for updating an existing template path
type UpdateTemplatePathParams struct {
	TempPathID  int    `json:"id" valid:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Courses     []int  `json:"courses"`
}

// DeleteTemplatePathParams defines parameters for deleting a template path
type DeleteTemplatePathParams struct {
	TempPathID int `json:"id" valid:"required"`
}
