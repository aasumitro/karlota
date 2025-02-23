package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed all:dist
var resource embed.FS

func SPA() (spa fs.FS) {
	var err error
	if spa, err = fs.Sub(resource, "dist"); err != nil {
		log.Fatalln("FS_ERR:", err.Error())
	}
	return spa
}
