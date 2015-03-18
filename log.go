package lsearch

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "LSearch: ", log.Lshortfile|log.Ldate|log.Ltime)
