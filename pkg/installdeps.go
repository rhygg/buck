package pkg

import (
	"fmt"
	"net/http"
	"path/filepath"

	cp "github.com/otiai10/copy"

	// "path/filepath"
	"github.com/c4milo/unpackit"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/rhygg/buck/cache"
	"github.com/rhygg/buck/utils"
	"github.com/schollz/progressbar/v3"
	"github.com/tidwall/gjson"
	"io"
	"strings"
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
		name := key.String()
		var n string
		ver := strings.ReplaceAll(strings.ReplaceAll(value.String(), "^", ""), "~", "")
		if strings.Contains(name, "@") {
			n = strings.Split(name, "/")[1]
		} else {
			n = name
		}
		url := utils.NpmRegistryURL + name + "/-/" + n + "-" + ver + ".tgz"
		err := Download(url, name, dirs, ver, cache, count, lent)
		ds = append(ds, DependenciesField{
			Name:    name,
			Version: value.String(),
		})
   count ++
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

func Download(url string, name string, dirs UserDirectories, vsion string, cache cache.ModuleCache, count int, lent int) error {
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
    progressbar.OptionSetDescription(fmt.Sprintf("[cyan][%d/%d][reset] pkg: %s", count, lent, name)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	io.Copy(io.Writer(bar), resp.Body)
	_, err = unpackit.Unpack(resp.Body, filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, name))
	if err != nil {
		fmt.Println("[BUCK] download: " + err.Error())
	}
	chk, _ := cache.Cache.Get("module:" + name)
	if chk == true {
		err = cp.Copy(filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, name, "package"), filepath.Join(dirs.CWD, "node_modules", name))
		if err != nil {
			fmt.Println("[BUCK] download: " + err.Error())
		}
	} else {
		err = cache.Cache.Set("module:"+name, true)
		if err != nil {
			fmt.Print("Could not set module cache.")
		}
		err = cp.Copy(filepath.Join(dirs.HomeDir, dirs.ModuleCacheDir, name, "package"), filepath.Join(dirs.CWD, "node_modules", name))
		if err != nil {
			fmt.Println("[BUCK] download: " + err.Error())
		}
	}
	return nil
}
