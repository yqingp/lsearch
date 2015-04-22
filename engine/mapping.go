package engine

import (
    "github.com/yqingp/lsearch/mapping"
)

type MappingRequest struct {
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
func (e *Engine) MappingHandler(body []byte) (interface{}, error) {
    mapping, err := mapping.New(body)
    if err != nil {
        return nil, err
    }

    action := mapping.Action

    if action == "create" {
        if err := e.NewIndex(mapping); err != nil {
            return nil, err
        }
    }

    if action == "delete" {
        if err := e.RemoveIndex(mapping); err != nil {
            return nil, err
        }
    }

    return nil, nil
}
