package engine

import (
    "github.com/yqingp/lsearch/analyzer"
    "github.com/yqingp/lsearch/config"
    "github.com/yqingp/lsearch/index"
    "github.com/yqingp/lsearch/log"
    "github.com/yqingp/lsearch/search"
)

type Engine struct {
    analyzer        *analyzer.Analyzer
    config          *config.Config
    version         string
    indexes         map[string]*index.Index
    indexWorkers    []chan IndexRequest
    searchWrokers   []chan SearchRequest
    analyzerWorkers []chan AnalyzerRequest
    mappingWorkers  []chan MappingRequest
    status          *Status
}

func (e *Engine) Init() error {
    e.config = config.NewConfig()
    log.Init()
    e.analyzer.Init()
    return nil
}

func (e *Engine) BindIpAndPort() string {
    ipAndPort := ""
    if e.bindIP != "" {
        ipAndPort += e.bindIP
    }
    ipAndPort += ":"
    if e.bindPort != "" {
        ipAndPort += e.bindPort
    } else {
        ipAndPort += "8866"
    }

    return ipAndPort
}

func (e *Engine) NewIndexMapping() {

}

func (e *Engine) UpdateIndexMapping() {

}

func (e *Engine) DeleteIndex() {

}

func (e *Engine) Index() {

}

func (e *Engine) Search() {

}

func (e *Engine) Status() {

}
