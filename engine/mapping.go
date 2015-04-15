package engine

type MappingRequest struct {
}

/*
post structure
{
	action: "update||create||delete",
	name:"abc"
	fields:[
	{:name => "id"}unique
	{:name => fields1, :type => (string|text), :analyzer => xxx}
	...
	...
	]
}

*/
func (e *Engine) MappingHandler(action string, body []byte) {

}
