package store

import (
    // "os"
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
