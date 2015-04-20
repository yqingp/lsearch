package engine

import (
    "github.com/yqingp/lsearch/index"
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
func (e *Engine) MappingHandler(body []byte) (*index.IndexMeta, error) {
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

    if action == "view" {
        results, err := e.ViewIndex(mapping)
        if err != nil {
            return nil, err
        }
        return results, nil
    }

    return nil, nil
}
