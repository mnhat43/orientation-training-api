package requestparams

type CreateModuleItemParams struct {
	Title    string `json:"title" form:"module_title" valid:"required"`
	ItemType string `json:"item_type" form:"module_item_type" valid:"required"`
	Url      string `json:"url" form:"module_url"`
	ModuleID int    `json:"module_id" form:"module_id" valid:"required"`
}

type ModuleItemListParams struct {
	ModuleID    int    `json:"module_id" valid:"required"`
	CurrentPage int    `json:"current_page" valid:"-"`
	RowPerPage  int    `json:"row_per_page"`
	Keyword     string `json:"keyword"`
}
