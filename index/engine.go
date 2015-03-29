package index

import (
	"errors"
	"fmt"
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

		indexPath := filepath.Join(storePath, fi.Name())

		index, _ := recoverIndex(indexPath)
		if index == nil {
			fmt.Println(indexPath)
			fmt.Println("blank dir")
			err := os.RemoveAll(indexPath)
			fmt.Println(err)
		}

	}
	return err
}
