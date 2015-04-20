package engine

import (
    "encoding/json"
    "errors"
    "github.com/yqingp/lsearch/document"
    "github.com/yqingp/lsearch/index"
    "github.com/yqingp/lsearch/mapping"
    "time"
)

type IndexRequest struct {
    Name         string              `json:"name"`
    Action       string              `json:"action"`
    Documents    []document.Document `json:"documents,omitempty"`
    RequestStart time.Time
    Status       chan bool
    Results      interface{}
    Duration     time.Duration
    Index        *index.Index
}

func (e *Engine) NewIndex(mapping *mapping.Mapping) error {
    e.mappingMutex.Lock()
    defer e.mappingMutex.Unlock()

    _, ok := e.indexes[mapping.Name]
    if ok {
        return errors.New("Index Exist")
    }

    index.New(mapping, e.Config.StorePath)

    return nil
}

func (e *Engine) RemoveIndex(mapping *mapping.Mapping) error {
    e.mappingMutex.Lock()
    defer e.mappingMutex.Unlock()

    index, ok := e.indexes[mapping.Name]
    if !ok {
        return errors.New("Index Not Found")
    }

    index.Remove()
    delete(e.indexes, mapping.Name)
    return nil
}

func (e *Engine) RecoverIndexes() {
    if e.isInit {
        return
    }
    e.indexes = index.RecoverIndexes(e.Config.StorePath)
}

func (e *Engine) ViewIndex(name string) (*index.IndexMeta, error) {
    index, ok := e.indexes[name]

    if !ok {
        return nil, errors.New("Index Not Found")
    }

    return index.View(), nil
}

func (e *Engine) Index(body []byte) {
    indexRequest := &IndexRequest{
        Body:         body,
        RequestStart: time.Now(),
        Status:       make(chan bool),
    }

    e.IndexRequests <- indexRequest
    <-indexRequest.Status
    indexRequest.Duration = time.Now().Sub(indexRequest.RequestStart)
}

func ParseIndexRequest(body []byte) (*IndexRequest, error) {
    request := &IndexRequest{}
    if err := json.Unmarshal(body, request); err != nil {
        return nil, errors.New("decode request error")
    }

    return request, nil
}
