package store

import (
    . "github.com/yqingp/lsearch/store"
    "strconv"
    "testing"
    "time"
)

func TestDb(t *testing.T) {
    db, err := Open("./db", true)
    if err != nil {
        t.Error(err)
    }

    if db == nil {
        t.Error("db init fail")
    }

    for i := 0; i < 10000; i++ {
        ret, err := db.Set(-1, []byte("test"+strconv.Itoa(i)), []byte("test"))
        if err != nil {
            t.Error(err)
        }
        t.Log(ret)
    }
    time.Sleep(1000000000)

    // t.Log(db.)
}
