package pkg

import (
	"errors"
	"path/filepath"

	cp "github.com/otiai10/copy"
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
  Packages        []Packages              `yaml:"packages"`
}

type Packages struct {
  Name string `yaml:"name"`
  Version string `yaml:"version"` 
}
type DependenciesField struct {
	Name    string `yaml:"name"`
	Version string `yaml: "version"`
}

type DevDependenciesField DependenciesField

func InstallPackage(mcache cache.ModuleCache, dirs UserDirectories, data string) error {
	latest := gjson.Get(data, "dist-tags.latest").String()
	name := gjson.Get(data, "name").String()
	install_url := utils.NpmRegistryURL + name + "/-/" + name + "-" + latest + ".tgz"
	check, _ := mcache.Cache.Get("module:" + name)
	if check == nil {
		mcache.Cache.Set("module:"+name, true)
		store := cache.StoreCache(name, install_url, filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, name))

		if store {
			err := cp.Copy(filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, name, "package"), filepath.Join(dirs.CWD, "node_modules", name))
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
	err := cp.Copy(filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, name, "package"), filepath.Join(dirs.CWD, "node_modules", name))
	if err != nil {
		return errors.New("could not download package")
	}

	if err != nil {
		return errors.New("could not install dependencies")
	}
	return nil

}
