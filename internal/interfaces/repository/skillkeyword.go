package repository

import (
	m "orientation-training-api/internal/models"
)

type SkillKeywordRepository interface {
	Create(skill *m.SkillKeyword) error
	GetByID(id int) (*m.SkillKeyword, error)
	List() ([]m.SkillKeyword, error)
	Update(skill *m.SkillKeyword) error
	Delete(id int) error
}
