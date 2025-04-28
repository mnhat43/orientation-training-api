package requestparams

type CreateModuleItemParams struct {
	Title        string `json:"title" form:"title" valid:"required"`
	ItemType     string `json:"item_type" form:"item_type" valid:"required"`
	Resource     string `json:"resource" form:"resource" valid:"required"`
	Position     int    `json:"position" form:"position" `
	RequiredTime int    `json:"required_time" form:"required_time"`
	ModuleID     int    `json:"module_id" form:"module_id" valid:"required"`
}

type ModuleItemListParams struct {
	ModuleID int `json:"module_id" valid:"required"`
}

type ModuleItemIDParam struct {
	ModuleItemID int `json:"moduleItem_id" valid:"required"`
}
