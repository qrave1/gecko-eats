package postgres

import "embed"

//go:embed migrations/*
var MigrationsEmbed embed.FS
