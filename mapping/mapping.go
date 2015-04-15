package mapping

import (
    "encoding/json"
    "errors"
    "github.com/yqingp/lsearch/field"
    "log"
)

type Mapping struct {
    Action string `json:"action"`
    Name   string `json:"name"`
    Fields []field.Filed
}

func NewMapping(body []byte) (*Mapping, error) {
    mapping := &Mapping{}
    if err := json.Unmarshal(body, mapping); err != nil {
        return nil, errors.New("decode mapping error")
    }

    log.Println(mapping)
    return mapping, nil
}
