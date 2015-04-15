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
    e.Config = config.NewConfig()
    e.analyzer.Init()
    return nil
}

func (e *Engine) BindAddr() string {
    return e.Config.BindAddr()
}

func (e *Engine) DeleteIndex() {

}

func (e *Engine) Index() {

}

func (e *Engine) Search() {

}

func (e *Engine) Status() {

}
