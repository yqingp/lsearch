package mapping

import (
    "encoding/json"
    "errors"
    "github.com/yqingp/lsearch/field"
    "log"
)

const (
    MaxFiledNum = 100
)

type Mapping struct {
    Action string        `json:"action,omitempty"`
    Name   string        `json:"name,omitempty"`
    Fields []field.Filed `json:"fields,omitempty"`
}

/*
post structure
{
	action: "update||create||delete",
	name:"abc"
	fields:[
	{:name => "id"}unique
	{:name => fields1, :type => 0,1}
	...
	...
	]
}
*/
func New(body []byte) (*Mapping, error) {
    mapping := &Mapping{}
    if err := json.Unmarshal(body, mapping); err != nil {
        return nil, errors.New("decode mapping error")
    }

    if err := mapping.validate(); err != nil {
        return nil, err
    }

    log.Println(mapping)

    return mapping, nil
}

func (m *Mapping) validate() error {
    if m.Action != "create" && m.Action != "delete" {
        return errors.New("mapping action error")
    }

    if m.Name == "" {
        return errors.New("mapping name error")
    }

    if m.Action == "delete" {
        return nil
    }

    if !m.validateFields() {
        return errors.New("mapping fields error")
    }

    return nil
}

func (m *Mapping) validateFields() bool {
    names := make(map[string]string)
    for _, v := range m.Fields {
        if _, ok := names[v.Name]; ok {
            return false
        } else {
            names[v.Name] = ""
        }

        if !v.Valid() {
            return false
        }
    }

    if _, ok := names["id"]; !ok {
        return false
    }

    if len(m.Fields) > MaxFiledNum {
        return false
    }

    return true
}
