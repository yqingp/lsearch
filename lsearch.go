package main

import (
    "encoding/json"
    "github.com/yqingp/lsearch/engine"
    "io/ioutil"
    "log"
    "net/http"
    "runtime"
)

var lsearch engine.Engine

type ResponseResult struct {
    Results      interface{} `json:"results"`
    Response     string      `json:"response"`
    ResponseCode int         `json:"code"`
}

func mapping(rw http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    defer req.Body.Close()

    result := ResponseResult{}
    results, err := lsearch.MappingHandler(body)
    if err != nil {
        result.Response = err.Error()
        result.ResponseCode = 500
    } else {
        result.ResponseCode = 200
        result.Response = "done"
        if results != nil {
            result.Results = results
        }
    }

    w, _ := json.Marshal(result)
    rw.WriteHeader(200)
    rw.Write(w)
}

func statusHandler(rw http.ResponseWriter, req *http.Request) {

}

func searchHandler(rw http.ResponseWriter, req *http.Request) {

}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
    log.Println(req.Form)
    // params :=

    // result := ResponseResult{}
    // w, _ := json.Marshal(result)
    // rw.WriteHeader(200)
    // rw.Write(w)
}

func viewIndexHandler(rw http.ResponseWriter, req *http.Request) {
    name := req.URL.Query()["name"][0]
    result := ResponseResult{}

    results, err := lsearch.ViewIndex(name)
    if err != nil {
        result.Response = err.Error()
        result.ResponseCode = 500
    } else {
        result.ResponseCode = 200
        result.Response = "done"
        if results != nil {
            result.Results = results
        }
    }

    w, _ := json.Marshal(result)
    rw.WriteHeader(200)
    rw.Write(w)
}

func routes() {
    http.HandleFunc("/_mapping", mapping)
    http.HandleFunc("/_status/", statusHandler)
    http.HandleFunc("/search", searchHandler)
    http.HandleFunc("/index/view", viewIndexHandler)
}

func main() {
    lsearch.Init()
    runtime.GOMAXPROCS(runtime.NumCPU())
    routes()
    http.ListenAndServe(lsearch.BindAddr(), nil)
}
