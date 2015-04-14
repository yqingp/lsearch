package main

import (
    "github.com/yqingp/lsearch/engine"
    "io/ioutil"
    // "log"
    "net/http"
    "runtime"
)

var lsearch engine.Engine

func mapping(rw http.ResponseWriter, req *http.Request) {
    action := req.RequestURI[len("/_mapping/"):]
    body, _ := ioutil.ReadAll(req.Body)
    lsearch.MappingHandler(action, body)
    req.Body.Close()
}

func statusHandler(rw http.ResponseWriter, req *http.Request) {
    action := req.RequestURI[len("/_mapping/"):]
    body, _ := ioutil.ReadAll(req.Body)
    lsearch.MappingHandler(action, body)
    req.Body.Close()
}

func searchHandler(rw http.ResponseWriter, req *http.Request) {

}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
}

func routes() {
    http.HandleFunc("/_mapping/", mapping)
    http.HandleFunc("/_status/", statusHandler)
    http.HandleFunc("/search", searchHandler)
    http.HandleFunc("/index", indexHandler)
}

func main() {
    lsearch.Init()
    runtime.GOMAXPROCS(runtime.NumCPU())
    routes()
    http.ListenAndServe(lsearch.BindAddr(), nil)
}
