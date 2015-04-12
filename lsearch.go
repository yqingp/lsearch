package main

import (
    "github.com/yqingp/lsearch/engine"
    "net/http"
    "runtime"
)

var lsearch engine.Engine

func mappingHandler(rw http.ResponseWriter, req *http.Request) {

}

func statusHandler(rw http.ResponseWriter, req *http.Request) {

}

func searchHandler(rw http.ResponseWriter, req *http.Request) {

}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
}

func routes() {
    http.HandleFunc("/_mapping/*", mappingHandler)
    http.HandleFunc("/_status/*", statusHandler)
    http.HandleFunc("/search", searchHandler)
    http.HandleFunc("/index", indexHandler)
}

func main() {
    lsearch.Init()
    runtime.GOMAXPROCS(runtime.NumCPU())
    routes()
    http.ListenAndServe(lsearch.BindIpAndPort(), nil)
}
