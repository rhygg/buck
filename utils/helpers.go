package utils

import (
	"os"
	"github.com/tidwall/gjson"
)

func CheckPkg(pkg string) (bool, string) {
	data := Request(NpmRegistryURL + pkg)

	if data == "" {
		return false, ""
	}

	if gjson.Get(data, "error").Exists() {
		return false, ""
	}
	return true, data
}

//Create a folder/directory at a full qualified path
func Create(path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

//check if file exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

