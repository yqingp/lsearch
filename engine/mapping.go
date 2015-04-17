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
func (e *Engine) MappingHandler(body []byte) error {
    mapping, err := mapping.New(body)
    if err != nil {
        return err
    }

    _ = mapping

    return nil
}
