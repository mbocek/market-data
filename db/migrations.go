package db

import "embed"

// Migrations is an embedded filesystem containing database migration files.
//
//go:embed migrations
var Migrations embed.FS
