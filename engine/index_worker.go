package engine

import (
    "runtime"
)

func (e *Engine) initIndexWorkers() {
    e.IndexRequests = make(chan *IndexRequest, runtime.NumCPU())
}

func (e *Engine) startIndexWorkers() {
    cpuNum := runtime.NumCPU()
    for i := 0; i < cpuNum; i++ {
        go doIndex(e)
    }
}

func doIndex(e *Engine) {
    for {
        request := <-e.IndexRequests
        switch request.Action {
        case "create":
            {
                // Logger.Println(request.Documents)
                results, err := request.Index.AddDocuments(request.Documents)
                request.Results = results
                request.Error = err
            }
        case "delete":
            request.Index.DeleteDocuments(request.Documents)
        case "update":
            request.Index.UpdateDocuments(request.Documents)
        }
        request.Status <- true
    }
}
