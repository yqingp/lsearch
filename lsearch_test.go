package main

import (
	"github.com/yqingp/lsearch/engine"
	// "log"
	// "os"
	// "encoding/json"
	"testing"
)

var mappingText []byte = []byte(`
	{
		"action":"create",
		"name":"weibo",
		"fields":[
			{
				"name":"id",
				"type":0,
				"is_index":true
			},
			{
				"name":"text",
				"type":1,
				"is_index":true
			}
		]
	}
`)

var indexText []byte = []byte(`
	{"name":"weibo","action":"create","documents":[{"id":"weiboid1","values":{"text":"测试索引"}},{"id":"weiboid2","values":{"text":"测试索引微博索引"}}]}
`)

var searchText []byte = []byte(`
    {
        "name":"weibo","query":{"text":"测试"}
    }
`)

func TestIndex(t *testing.T) {
	var lsearch engine.Engine
	lsearch.Init()
	_ = mappingText
	_ = indexText
	_ = searchText
	// lsearch.MappingHandler(mappingText)

	// lsearch.Index(indexText)

	// a, _ := json.Marshal(lsearch.Search(searchText))
	// log.Println(string(a))

	// os.RemoveAll("dbpath")
}
