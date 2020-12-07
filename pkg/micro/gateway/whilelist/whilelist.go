package whitelist

import (
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/reader"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/util/log"
	"sort"
	"strings"
	"sync"
)

type Whitelist struct {
	config.Config
	data []string
	rock sync.RWMutex
}

func (wl *Whitelist) update(value reader.Value) error {
	wl.clean()
	wl.rock.Lock()
	defer func() {
		wl.rock.Unlock()
	}()

	var whiteList []string
	err := value.Scan(&whiteList)
	if err != nil {
		return err
	}
	wl.data = whiteList

	log.Info("load whitelist: [", strings.Join(whiteList, ",")+"]")

	return nil
}

func (wl *Whitelist) In(url string) bool {
	wl.rock.Lock()
	defer func() {
		wl.rock.Unlock()
	}()
	if index := sort.SearchStrings(wl.data, url); index < 0 {
		return false
	}

	return true
}

func (wl *Whitelist) clean() {
	wl.rock.Lock()
	defer wl.rock.Unlock()
	wl.data = make([]string, 0)
}

func NewWhitelist(source source.Source, path ...string) (*Whitelist, error) {
	wl := &Whitelist{
		Config: config.NewConfig(),
		data:   make([]string, 0),
		rock:   sync.RWMutex{},
	}
	err := wl.Load(source)
	if err != nil {
		return nil, err
	}
	value := wl.Get(path...)
	if err := wl.update(value); err != nil {
		return nil, err
	}
	wl.enableAutoUpdate(path...)

	return wl, nil
}

func (wl *Whitelist) enableAutoUpdate(path ...string) error {
	w, err := wl.Config.Watch(path...)
	if err != nil {
		return err
	}
	go func() {
		for {
			v, err := w.Next()
			if err != nil {
				log.Error(err)
				continue
			}
			if err := wl.update(v); err != nil {
				log.Error(err)
			}
		}
	}()
	return nil
}
