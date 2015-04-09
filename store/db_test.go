package store

import (
    . "github.com/yqingp/lsearch/store"
    "testing"
)

func TestDb(t *testing.T) {
    db, err := Open("/Users/yanqingpei/Go/src/github.com/yqingp/lsearch/store/db", true)
    if err != nil {
        t.Error(err)
    }

    if db == nil {
        t.Error("db init fail")
    }

    ret, err := db.Set(-1, []byte("test"), []byte("test"))
    if err != nil {
        t.Error(err)
    }
    t.Log(ret)
    // t.Log(db.)
}
