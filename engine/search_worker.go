package engine

import (
    "runtime"
)

func (e *Engine) startSearchWorkers() {
    e.SearchRequests = make(chan *SearchRequest, runtime.NumCPU())

    cpuNum := runtime.NumCPU()
    for i := 0; i < cpuNum; i++ {
        go doSearch(e)
    }
}

func doSearch(e *Engine) {
    for {
        request := <-e.SearchRequests
        results, err := request.Index.Search(request.Query)
        request.Results = results
        request.Error = err
        request.Status <- true
    }
}
