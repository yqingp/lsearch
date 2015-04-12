package engine

import (
    "github.com/yqingp/lsearch/log"
)

type Engine struct {
    basePath string
    bindIP   string
    bindPort string
}

func (e *Engine) Init() error {
    log.Init()
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
