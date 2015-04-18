package main

import (
	"encoding/json"
	"github.com/yqingp/lsearch/engine"
	"io/ioutil"
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

	if err := lsearch.MappingHandler(body); err != nil {
		result.Response = err.Error()
		result.ResponseCode = 500
	} else {
		result.ResponseCode = 200
		result.Response = "done"
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
}

func routes() {
	http.HandleFunc("/_mapping", mapping)
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
