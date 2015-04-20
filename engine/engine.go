package engine

import (
    // "fmt"
    "github.com/yqingp/lsearch/analyzer"
    "github.com/yqingp/lsearch/config"
    "github.com/yqingp/lsearch/index"
    "sync"
    // "github.com/yqingp/lsearch/search"
)

type Engine struct {
    analyzer         *analyzer.Analyzer
    Config           *config.Config
    version          string
    indexes          map[string]*index.Index
    IndexRequests    chan *IndexRequest
    SearchRequests   chan *SearchRequest
    AnalyzerRequests chan *AnalyzerRequest
    status           *Status
    mappingMutex     *sync.Mutex
    isInit           bool
}

func (e *Engine) Init() {
    e.Config = config.New()
    e.analyzer.Init()
    e.initMutex()
    e.RecoverIndexes()
    e.isInit = true
}

func (e *Engine) BindAddr() string {
    return e.Config.BindAddr()
}

func (e *Engine) initMutex() {
    e.mappingMutex = &sync.Mutex{}
}
