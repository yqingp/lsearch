package engine

import (
	// "fmt"
	"github.com/yqingp/lsearch/analyzer"
	"github.com/yqingp/lsearch/config"
	"github.com/yqingp/lsearch/index"
	"log"
	"os"
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

var Logger *log.Logger = log.New(os.Stdout, "DEGUG", log.Llongfile|log.Ldate|log.Ltime)

func (e *Engine) Init() {
	e.Config = config.New()
	e.analyzer.Init()
	e.initMutex()
	e.RecoverIndexes()
	e.startIndexWorkers()
	e.startSearchWorkers()
	e.isInit = true
}

func (e *Engine) BindAddr() string {
	return e.Config.BindAddr()
}

func (e *Engine) initMutex() {
	e.mappingMutex = &sync.Mutex{}
}
