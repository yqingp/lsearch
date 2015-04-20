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
    Error        error
}

func (e *Engine) NewIndex(mapping *mapping.Mapping) error {
    e.mappingMutex.Lock()
    defer e.mappingMutex.Unlock()

    _, ok := e.indexes[mapping.Name]
    if ok {
        return errors.New("Index Exist")
    }

    e.indexes[mapping.Name] = index.New(mapping, e.Config.StorePath)
    e.indexes[mapping.Name].Analyzer = e.analyzer
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

    for _, v := range e.indexes {
        v.Analyzer = e.analyzer
    }
}

func (e *Engine) ViewIndex(name string) (*index.IndexMeta, error) {
    index, ok := e.indexes[name]

    if !ok {
        return nil, errors.New("Index Not Found")
    }

    return index.View(), nil
}

func (e *Engine) Index(body []byte) error {
    indexRequest, err := ParseIndexRequest(body)
    if err != nil {
        return err
    }

    indexRequest.RequestStart = time.Now()
    indexRequest.Status = make(chan bool)

    if err = indexRequest.Valid(); err != nil {
        return err
    }

    index, ok := e.indexes[indexRequest.Name]

    if !ok {
        return errors.New("Index Not Found")
    }
    indexRequest.Index = index

    e.IndexRequests <- indexRequest
    <-indexRequest.Status
    indexRequest.Duration = time.Now().Sub(indexRequest.RequestStart)

    return nil
}

func ParseIndexRequest(body []byte) (*IndexRequest, error) {
    request := &IndexRequest{}
    if err := json.Unmarshal(body, request); err != nil {
        return nil, errors.New("decode request error")
    }
    return request, nil
}

func (i *IndexRequest) Valid() error {
    if i.Name == "" || i.Action == "" || i.Documents == nil {
        return errors.New("Index Request Error")
    }

    i.Duration = time.Now().Sub(i.RequestStart)

    return nil
}
