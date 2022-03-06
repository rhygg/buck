package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	// "path/filepath
	"io"
	"strings"

	"github.com/codeclysm/extract/v3"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/rhygg/buck/cache"
	"github.com/rhygg/buck/utils"
	"github.com/schollz/progressbar/v3"
	"github.com/tidwall/gjson"
)

type Deps struct {
	Name    string
	Version string
}

func InstallDeps(data string, name string, dirs UserDirectories, cache cache.ModuleCache) []DependenciesField {
	latest := gjson.Get(data, "dist-tags.latest")
	var ds []DependenciesField
	vsion := utils.Request(utils.NpmRegistryURL + name + "/" + latest.String())
	deps := gjson.Get(vsion, "dependencies")
	lent := len(deps.Array())
	fmt.Println(lent)
	count := 1
	deps.ForEach(func(key, value gjson.Result) bool {
		count++
		return true
	})
	c := 1
	deps.ForEach(func(key, value gjson.Result) bool {
		name := key.String()
		var n string
		ver := strings.ReplaceAll(strings.ReplaceAll(value.String(), "^", ""), "~", "")
		if strings.Contains(name, "@") {
			n = strings.Split(name, "/")[1]
		} else {
			n = name
		}
		url := utils.NpmRegistryURL + name + "/-/" + n + "-" + ver + ".tgz"
		ds = append(ds, DependenciesField{
			Name:    name,
			Version: value.String(),
		})
		err := Download(url, name, dirs, ver, cache, count, c)

		return err == nil
	})
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()
	magenta := color.New(color.FgHiMagenta).SprintFunc()
	fmt.Println(magenta("\n Dependencies:"))
	fmt.Println(green("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"))
	if len(ds) == 0 {
		fmt.Println("No Extra Dependencies")
	} else {
		for _, pkg := range ds {
			fmt.Print(green("\n" + pkg.Name))
			fmt.Print(" " + yellow(pkg.Version))
		}
	}
	fmt.Println(green("\n +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"))
	return ds
}

func Download(url string, name string, dirs UserDirectories, vsion string, cache cache.ModuleCache, count int, c int) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bar := progressbar.NewOptions(int(resp.ContentLength),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan][%d/%d][reset] pkg: %s", count, count, name)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	exists, err := utils.Exists(filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, filepath.Base(url)))
	if err != nil {
		fmt.Println("[BUCK] download(temp): " + err.Error())
	}

	if exists {
		err := os.Remove(filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, filepath.Base(url)))
		if err != nil {
			fmt.Println("[BUCK] download(temp): " + err.Error())
		}
	}

	f, _ := os.OpenFile(filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, filepath.Base(url)), os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	io.Copy(io.MultiWriter(bar, f), resp.Body)
	var shift = func(path string) string {
		parts := strings.Split(path, string(filepath.Separator))
		parts = parts[1:]
		return strings.Join(parts, string(filepath.Separator))
	}
	data, _ := ioutil.ReadFile(filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, filepath.Base(url)))
	buffer := bytes.NewBuffer(data)
	err = extract.Archive(context.Background(), buffer, filepath.Join(dirs.CWD, "node_modules", name), shift)
	if err != nil {
		fmt.Println("[BUCK] download(unpack): " + err.Error())
	}
	chk, _ := cache.Cache.Get("module:" + name)
	if chk == false {
		err = cache.Cache.Set("module:"+name, true)
		if err != nil {
			fmt.Print("Could not set module cache.")
		}
	}
	return nil
}
