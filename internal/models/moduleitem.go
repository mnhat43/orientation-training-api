package models

import (
	cm "orientation-training-api/internal/common"
)

type ModuleItem struct {
	cm.BaseModel

	ID       int    `pg:",pk"`
	Title    string `pg:",notnull"`
	ItemType string `pg:",notnull"`
	Resource string `pg:",notnull"`
	ModuleID int    `pg:",fk:module_id"`
	Position int    `pg:",notnull"`

	Module Module `pg:"rel:has-one"`
}
