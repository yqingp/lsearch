package engine

import (
    "github.com/yqingp/lsearch/analyzer"
    "github.com/yqingp/lsearch/config"
    "github.com/yqingp/lsearch/index"
    "github.com/yqingp/lsearch/log"
    "github.com/yqingp/lsearch/search"
)

type Engine struct {
    basePath string
    bindIP   string
    bindPort string
    analyzer *analyzer.Analyzer
    config   *config.Config
    version  string
    indexes  []*index.Index
}

func (e *Engine) Init() error {
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
