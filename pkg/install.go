package pkg

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/codeclysm/extract/v3"
	"github.com/rhygg/buck/cache"
	"github.com/rhygg/buck/utils"
	"github.com/tidwall/gjson"
)

type UserDirectories struct {
	CWD            string
	HomeDir        string
	ModuleCacheDir string
}

type LockFile struct {
	LockVersion     string                 `yaml:"lockFileVersion"`
	Dependencies    []DependenciesField    `yaml:"dependencies"`
	DevDependencies []DevDependenciesField `yaml:"devDependencies"`
	Packages        []Packages             `yaml:"packages"`
}

type Packages struct {
	Name         string              `yaml:"name"`
	Version      string              `yaml:"version"`
	Integrity    string              `yaml:"integrity"`
	Dev          bool                `yaml:"dev"`
	Engine       string              `yaml:"engine"`
	Dependencies []DependenciesField `yaml:"dependencies"`
	Resolved     string              `yaml:"resolved"`
}
type DependenciesField struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type DevDependenciesField DependenciesField

func InstallPackage(mcache cache.ModuleCache, dirs UserDirectories, data string) error {
	latest := gjson.Get(data, "dist-tags.latest").String()
	name := gjson.Get(data, "name").String()
	install_url := utils.NpmRegistryURL + name + "/-/" + name + "-" + latest + ".tgz"
	check, _ := mcache.Cache.Get("module:" + name)
	if check == nil {
		mcache.Cache.Set("module:"+name, true)
		store := cache.StoreCache(name, install_url, filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir))
		if store {
			var shift = func(path string) string {
				parts := strings.Split(path, string(filepath.Separator))
				parts = parts[1:]
				return strings.Join(parts, string(filepath.Separator))
			}
			data, _ := ioutil.ReadFile(filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, filepath.Base(install_url)))
			buffer := bytes.NewBuffer(data)
			err := extract.Archive(context.Background(), buffer, filepath.Join(dirs.CWD, "node_modules", name), shift)
			if err != nil {
				return errors.New("could not download package")
			}
			if err != nil {
				return errors.New("could not install dependencies")
			}
		} else {
			return errors.New("could not store package")

		}
	}
	var shift = func(path string) string {
		parts := strings.Split(path, string(filepath.Separator))
		parts = parts[1:]
		return strings.Join(parts, string(filepath.Separator))
	}
	data2, _ := ioutil.ReadFile(filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, filepath.Base(install_url)))
	buffer2 := bytes.NewBuffer(data2)

	err := extract.Archive(context.Background(), buffer2, filepath.Join(dirs.CWD, "node_modules", name), shift)
	if err != nil {
		return errors.New("could not download package")
	}

	if err != nil {
		return errors.New("could not download package")
	}

	if err != nil {
		return errors.New("could not install dependencies")
	}
	return nil

}
