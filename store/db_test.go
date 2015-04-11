package store

import (
    // "os"
    "strconv"
    "testing"
    "time"
)

func TestGet(t *testing.T) {
    db, err := Open("./db", true)
    if err != nil {
        t.Error(err)
    }

    if db == nil {
        t.Error("db init fail")
    }

    // start := time.Now()
    ret, _ := db.Get([]byte("test0"))
    t.Log(string(ret))
    ret, _ = db.GetInternalId(1)
    t.Log(string(ret))
    db.Close()
    // end := time.Now()
}

func TestDb(t *testing.T) {
    db, err := Open("./db", true)
    if err != nil {
        t.Error(err)
    }

    if db == nil {
        t.Error("db init fail")
    }

    start := time.Now()
    for i := 0; i < 10000; i++ {
        _, err := db.Set(-1, []byte("test"+strconv.Itoa(i)), []byte("test"))
        if err != nil {
            t.Error(err)
        }
        // t.Log(ret)
    }
    end := time.Now()
    t.Log("=============================")
    t.Log("10000 values insert spend:(ms)", end.Sub(start))
    db.Close()
    // os.RemoveAll("db")
    // t.Log(db.)
}
