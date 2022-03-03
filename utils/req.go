package utils

import (
	"io/ioutil"
	"net/http"
)

func Request(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}
