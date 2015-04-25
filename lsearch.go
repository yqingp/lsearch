package main

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/yqingp/lsearch/engine"
    "io/ioutil"
    "log"
    "net/http"
    "runtime"
)

var LSearch engine.Engine

var Router *gin.Engine

type ResponseResult struct {
    Results      interface{} `json:"results"`
    Response     string      `json:"response"`
    ResponseCode int         `json:"code"`
}

func mapping(rw http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    defer req.Body.Close()

    result := ResponseResult{}
    results, err := LSearch.MappingHandler(body)
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

    results, err := LSearch.ViewIndex(name)
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

func indexDocumentsHandler(rw http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    defer req.Body.Close()

    results := LSearch.Index(body)
    w, _ := json.Marshal(results)
    rw.WriteHeader(200)
    rw.Write(w)
}

func Routes() {
    // 创建索引
    Router.POST("/:name", func(c *gin.Context) {
    })

    // 更新索引
    Router.PUT("/:name", func(c *gin.Context) {

    })

    // 删除索引
    Router.DELETE("/:name", func(c *gin.Context) {

    })

    // 查看索引数据
    Router.GET("/:name", func(c *gin.Context) {

    })

    // 查看索引结构
    Router.GET("/:name/_mapping", func(c *gin.Context) {

    })

    // 查看索引状态统计
    Router.GET("/:name/_status", func(c *gin.Context) {

    })

    // 查看索引状态统计
    Router.GET("/_status", func(c *gin.Context) {

    })

    // 根据ID同时查询多条数据 多个索引下的数据
    //     "docs" : [
    //     {
    //         "index" : "test",
    //         "id" : "1"
    //     },
    //     {
    //         "index" : "test",
    //         "id" : "2"
    //     }
    // ]
    Router.GET("/_mget", func(c *gin.Context) {

    })

    // 某个索引下面的多个数据根据ID查询
    // "ids" : ["1", "2"]
    Router.GET("/:name/_mget", func(c *gin.Context) {

    })

    // 查看索引数据 单条
    Router.GET("/:name/:id", func(c *gin.Context) {

    })

    // 查看删除单挑数据
    Router.DELETE("/:name/:id", func(c *gin.Context) {

    })

    // 查看更新单条素据或者插入
    Router.PUT("/:name/:id", func(c *gin.Context) {

    })

    //插入数据
    Router.POST("/:name/:id", func(c *gin.Context) {

    })

    // 批量操作
    Router.POST("/_bulk", func(c *gin.Context) {

    })

    // 查询条件 删除
    Router.DELETE("/:name/_query", func(c *gin.Context) {

    })

    // 查询条件 更新
    Router.POST("/:name/_query", func(c *gin.Context) {

    })

    Router.POST("/:name/_search", func(c *gin.Context) {

    })

    Router.GET("/:name/_search", func(c *gin.Context) {

    })

    // http.HandleFunc("/_mapping", mapping)
    // http.HandleFunc("/_status/", statusHandler)
    // http.HandleFunc("/search", searchHandler)
    // http.HandleFunc("/index/view", viewIndexHandler)
    // http.HandleFunc("/index/documents", indexDocumentsHandler)
}

func main() {
    Router = gin.Default()
    LSearch.Init()
    runtime.GOMAXPROCS(runtime.NumCPU())
    Routes()
    Router.Run(LSearch.BindAddr())
}
