package engine

import (
    "runtime"
)

func (e *Engine) StartWorkers() {
    cpuNum := runtime.NumCPU()
    for i := 0; i < cpuNum; i++ {
        go DoIndex(e)
    }
}

func DoIndex(e *Engine) {
    for {
        request := <-e.IndexRequests
        switch request.Action {
        case "create":
            {
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
