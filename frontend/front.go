package frontend

import (
	"embed"
	"io/fs"
)

//go:embed dist
var distContent embed.FS

func GetDistFS() fs.FS {
	FsSub, err := fs.Sub(distContent, "dist")
	if err != nil {
		panic(err)
	}
	return FsSub
}
