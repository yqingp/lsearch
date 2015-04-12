package engine

import (
    "log"
    "os"
)

var Logger = log.New(os.Stderr, "LSearch: ", log.Lshortfile|log.Ldate|log.Ltime)
