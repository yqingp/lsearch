package main

// "github.com/yqingp/lsearch/"

import (
    "fmt"
    "github.com/yqingp/lsearch/index"
)

type LSearch struct {
    config *config
    engine *index.Engine
}

func NewLSearch(configFilePath string) *LSearch {
    return &LSearch{
        config: newConfig(configFilePath),
        engine: &index.Engine{},
    }
}

func (this *LSearch) Init() {
    fmt.Println("init")
    this.config.initStorePath()
    if err := this.engine.Init(this.config.storePath); err != nil {
        Logger.Fatal(err)
    }
}
