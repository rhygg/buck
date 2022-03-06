package cache

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rhygg/buck/utils"
	"github.com/schollz/progressbar/v3"
)

func StoreCache(name string, url string, modCacheDir string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	base := filepath.Base(url)
	exists, err := utils.Exists(filepath.Join(modCacheDir, base))
	if err != nil {
		return false
	}
	if exists {
		err := os.Remove(base)

		if err != nil {
			return false
		}
	}

	f, _ := os.OpenFile(filepath.Join(modCacheDir, base), os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading: "+name,
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		return false
	}

	if err != nil {
		return false
	}

	return true

}
