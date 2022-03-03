package cache

import (
	"net/http"

	"github.com/c4milo/unpackit"
)

func StoreCache(name string, url string, modCacheDir string) bool {

	resp, err := http.Get(url)

	if err != nil {
		return false
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
