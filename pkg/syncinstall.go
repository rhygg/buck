package pkg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/rhygg/buck/cache"
	"github.com/rhygg/buck/utils"
	"github.com/shomali11/parallelizer"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

type InstalledMods struct {
	Dependencies   []DependenciesField
	DependencyPkgs []DependenciesField
}

func InstallPackagesSynchronously(packages []string, Mcache cache.ModuleCache, rcache cache.RequestData, count int, cwd string, homeDir string) error {
	// init the waitgroup and channel
	group := parallelizer.NewGroup(parallelizer.WithPoolSize(count * 2))
	var mods InstalledMods
	defer group.Close()
	start := time.Now()
	exists, _ := utils.Exists(filepath.Join(homeDir, "buck-cache"))
	ex, _ := utils.Exists(filepath.Join(cwd, "node_modules"))
	if !ex {
		err := utils.Create(filepath.Join(cwd, "node_modules"))
		if err != nil {
			log.Fatal("Could not create node-module directory")
		}
	}
	if !exists {
		err := utils.Create(filepath.Join(homeDir, "buck-cache"))
		if err != nil {
			log.Fatal("Could not create cache directory")
		}
	}
	for _, pkg := range packages {
		data, errorC := rcache.Cache.Get(pkg)
		d, _ := data.(string)
		name := gjson.Get(d, "name").String()
		version := gjson.Get(d, "dist-tags.latest").String()
		mods.Dependencies = append(mods.Dependencies, DependenciesField{Name: name, Version: version})
		if errorC != nil {
      fmt.Print("[Buck] cache: " + errorC.Error())
		}
		group.Add(func() {
			InstallPackage(Mcache, UserDirectories{
				CWD:            cwd,
				HomeDir:        homeDir,
				ModuleCacheDir: "buck-cache",
			}, d)
		})

		deps := InstallDeps(d, pkg, UserDirectories{
			CWD:            cwd,
			HomeDir:        homeDir,
			ModuleCacheDir: "buck-cache",
		}, Mcache)

		mods.DependencyPkgs = append(mods.DependencyPkgs, deps...)
	}
	var exDeps []DependenciesField
	var Lf LockFile
	ex, err := utils.Exists(filepath.Join(cwd, "buck-lock.yaml"))
	if err != nil {
		return errors.New("could not check if lockfile exists")
	}
	if ex == true {
		yfile, err := ioutil.ReadFile(filepath.Join(cwd, "buck-lock.yaml"))
		if err != nil {
			return errors.New("could not read lockfile")
		}
		err2 := yaml.Unmarshal(yfile, &Lf)
		if err2 != nil {
			return err2
		}

		exDeps = Lf.Dependencies

		mods.Dependencies = append(mods.Dependencies, exDeps...)

	}
	lfdata := LockFile{LockVersion: "1.0.0", Dependencies: mods.Dependencies}
	marsh, err := yaml.Marshal(&lfdata)
	if err != nil {
		return errors.New("could not create lockfile.")
	}
	err = ioutil.WriteFile("buck-lock.yaml", marsh, 0644)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	fmt.Println("🐕 Installed: " + strings.Join(packages, ", ") + " in " + elapsed.String())
	return nil
}
