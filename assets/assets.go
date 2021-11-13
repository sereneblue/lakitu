package assets

import "embed"

//go:embed build/* build/_app/pages/* build/_app/assets/pages/*
var App embed.FS
