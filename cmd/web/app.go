package main

import (
	"awesomeProject/pkg/models"
	"github.com/alexedwards/scs"
)

type App struct {
	Database  *models.Database
	HTMLDir   string
	Sessions  *scs.Manager
	StaticDir string
}
