package engine

import (
    "runtime"
)

func (e *Engine) StartWorkers() {
    cpuNum := runtime.NumCPU()
    for i := 0; i < cpuNum; i++ {
        go IndexWork(e)
    }
}

func IndexWork(e *Engine) {
    for {
        request := <-e.IndexRequests
    }
}
