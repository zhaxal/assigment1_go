package main

import (
	"snippetbox/pkg/models"
)

type App struct {
	Database  *models.Database
	HTMLDir   string
	StaticDir string
}
