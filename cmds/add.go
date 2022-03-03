package cmds

import (
	"errors"
	"fmt"
	"time"

	"github.com/rhygg/buck/cache"
	"github.com/rhygg/buck/pkg"
	"github.com/rhygg/buck/utils"
	"github.com/theckman/yacspin"
)

type AddConfig struct {
	CWD      string
	Packages []string
	HomeDir  string
}

func Add(c AddConfig, cacheR cache.RequestData, cacheM cache.ModuleCache) error {
	// counts the amount of packages that are added (for waitgroup functionality)
	count := 0
	//check if all the packages are valid
	for pkg := range c.Packages {
		cfg := yacspin.Config{
			Frequency:       100 * time.Millisecond,
			CharSet:         yacspin.CharSets[18],
			Suffix:          "  Checking validation of package " + c.Packages[pkg],
			SuffixAutoColon: true,
			StopColors:      []string{"fgGreen"},
		}
		spinner, err := yacspin.New(cfg)

		if err != nil {
			fmt.Println(utils.LogError(err.Error()))
		}

		err = spinner.Start()

		check, data := utils.CheckPkg(c.Packages[pkg])
		if !check {
			fmt.Println("\n"+(utils.LogError("[Buck]")+" Could not find "+c.Packages[pkg]+" Exiting now... \n"), err)
			return errors.New("Could not find package " + c.Packages[pkg])
		} else {
			count++
			cacheR.Store(c.Packages[pkg], data)
		}
		spinner.Stop()
	}
	fmt.Print(utils.LogInfo("✓") + "  Package check is completed ... \n")
	fmt.Print(utils.LogInfo("✍") + "  Installing packages ... \n")
	pkg.InstallPackagesSynchronously(c.Packages, cacheM, cacheR, count, c.CWD, c.HomeDir)

	return nil
}
