package index

import (
	"errors"
	// "fmt"
	"os"
	"path/filepath"
)

type Engine struct {
	isInit         bool
	indexes        []*Index
	indexStorePath string
}

func (engine *Engine) Init(storePath string) error {
	if engine.isInit {
		return errors.New("inited")
	}
	f, err := os.Open(storePath)
	defer f.Close()

	if err != nil {
		return err
	}

	list, err := f.Readdir(-1)
	for _, fi := range list {
		if !fi.IsDir() {
			continue
		}

		indexPath, err := filepath.Abs(fi.Name())
		if err != nil {
			return err
		}

		recoverIndex(indexPath)

	}
	return err
}
