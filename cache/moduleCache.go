package cache

import (
	"github.com/c4milo/unpackit"
	"net/http"
  "io"
  "github.com/schollz/progressbar/v3"
)

func StoreCache(name string, url string, modCacheDir string) bool {
	resp, err := http.Get(url)
  bar := progressbar.DefaultBytes(
    resp.ContentLength,
    "downloading: "+name,
)
io.Copy(io.Writer(bar), resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		return false
	}
	dest, err := unpackit.Unpack(resp.Body, modCacheDir)

	if err != nil {
		return false
	}

	if dest == "" {
		return false
	}

	return true

}
