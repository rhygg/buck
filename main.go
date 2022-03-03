package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/rhygg/buck/cache"
	"github.com/rhygg/buck/cmds"
	"github.com/rhygg/buck/utils"
	"github.com/urfave/cli/v2"
)

func DeleteModCache(dir string, mcache cache.ModuleCache) bool {
  files, _ := ioutil.ReadDir(dir)
    for _, f := range files {
        fi, err := os.Stat(filepath.Join(dir, f.Name()))
        if err != nil {
            fmt.Println(err)
        }
        currTime := time.Now()
        fTime := fi.ModTime()
        if  fTime.Sub(currTime).Hours()/24.00 >= 1 {
          fmt.Println("[BUCK] cache: removing "+f.Name()+" from cache.")
          err := os.Remove(filepath.Join(dir, f.Name()))
          if err != nil {
            return false
          }
          del := mcache.Cache.Remove(f.Name())
          return del == true
        }
    }
  return false
}
func main() {
	cwd, err := os.Getwd()
	if err != nil {
		utils.LogError(err.Error())
	}
	current, err := user.Current()

	if err != nil {
    fmt.Print("[BUCK] main: Could not find home directory location, exiting...")
	}

	homedir := current.HomeDir

	cacheR := cache.Newrdc(100)
	cacheM := cache.Newmc(50, homedir, "buck-cache")
  pos := DeleteModCache(filepath.Join(homedir, "buck-cache"), cacheM)
  if !pos {
    fmt.Print("[BUCK] cache: could not delete outdated cache.")
  }
	app := &cli.App{
		Name:      "Buck",
		Copyright: "(c) 2022 Adonis Tremblay",
		HelpName:  "Buck",
		Version:   "0.0.1",
		Usage:     "A NodeJS package manager, built for speed",
		UsageText: "buck [global options] command [command options] [arguments...]",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"install", "a", "i"},
				Usage:   "Install a/multiple package(s)",
				Action: func(c *cli.Context) error {
					cmds.Add(cmds.AddConfig{
						CWD:      cwd,
						Packages: c.Args().Slice(),
						HomeDir:  homedir,
					}, cacheR, cacheM)
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"uninstall", "r", "u"},
				Usage:   "Remove a/multiple package(s)",
				Action: func(c *cli.Context) error {
					fmt.Println("uninstall")
					return nil
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Update a package",
				Action: func(c *cli.Context) error {
					fmt.Println("update")
					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}
