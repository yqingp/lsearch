package engine

import (
    "github.com/yqingp/lsearch/analyzer"
    "github.com/yqingp/lsearch/config"
    "github.com/yqingp/lsearch/index"
    "sync"
    // "github.com/yqingp/lsearch/search"
)

type Engine struct {
    analyzer        *analyzer.Analyzer
    Config          *config.Config
    version         string
    indexes         map[string]*index.Index
    indexWorkers    []chan IndexRequest
    searchWrokers   []chan SearchRequest
    analyzerWorkers []chan AnalyzerRequest
    mappingWorkers  []chan MappingRequest
    status          *Status
    mappingMutex    *sync.Mutex
}

func (e *Engine) Init() error {
    e.Config = config.New()
    e.analyzer.Init()
    e.initMutex()
    return nil
}

func (e *Engine) BindAddr() string {
    return e.Config.BindAddr()
}

func (e *Engine) initMutex() {
    e.mappingMutex = &sync.Mutex{}
}
