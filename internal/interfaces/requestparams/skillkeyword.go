package requestparams

type CreateSkillKeywordRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateSkillKeywordRequest struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type DeleteSkillKeywordRequest struct {
	ID int `json:"id" validate:"required"`
}
