//go:build !production
// +build !production

package embed

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// Content serves content in the /embed directory.
var Content fs.FS

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Panicf("cannot get current directory: %v", err)
	}

	embedDir := filepath.Join(wd, "embed")
	stat, err := os.Stat(embedDir)
	if err != nil {
		log.Panicf("cannot stat embed directory: %v", err)
	}
	if !stat.IsDir() {
		log.Panicf("embed directory is not a directory: %s", embedDir)
	}

	log.Printf("[dev] serving embedded content from %s", embedDir)

	Content = os.DirFS(embedDir)
}
