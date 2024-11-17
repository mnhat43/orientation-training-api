package common

import (
	db "orientation-training-api/internal/platform/db"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
)

type BaseController struct {
	DB     *pg.DB
	Logger echo.Logger
}

func (repo *BaseController) Init(logger echo.Logger) {
	repo.Logger = logger
	repo.DB = db.Init(logger)
}
