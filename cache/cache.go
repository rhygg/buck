package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bluele/gcache"
)

type RequestData struct {
	Cache     gcache.Cache
	CacheSize int
}

type ModuleCache struct {
	Cache gcache.Cache
	csize int
}

func Newmc(csize int, HomeDir string, ModuleCacheDir string) ModuleCache {
	return ModuleCache{
		Cache: gcache.New(csize).LRU().EvictedFunc(func(key, value interface{}) {
			res := strings.TrimLeft(key.(string), "module:")
			err := os.RemoveAll(filepath.Join(HomeDir, ModuleCacheDir, res))
			if err != nil {
				fmt.Print("[Error] " + "trace: error while remove module cache.")
			}
			fmt.Print("[cache] " + "evicted " + key.(string) + "from cache.")
		}).Build(),
		csize: csize,
	}
}
func Newrdc(csize int) RequestData {
	return RequestData{
		Cache: gcache.New(csize).LRU().EvictedFunc(func(key, value interface{}) {
			fmt.Println("[cache] " + "evicted " + key.(string) + "request data from cache.")
		}).Build(),
		CacheSize: csize,
	}
}

func (rdc RequestData) Store(key string, value interface{}) {
	rdc.Cache.Set(key, value)
}
