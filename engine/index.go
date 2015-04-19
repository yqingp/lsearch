package engine

import (
    "errors"
    "github.com/yqingp/lsearch/index"
    "github.com/yqingp/lsearch/mapping"
)

type IndexRequest struct {
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
