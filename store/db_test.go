package store

import (
    "os"
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

    start := time.Now().Nanosecond()
    for i := 0; i < 100; i++ {
        ret, err := db.Set(-1, []byte("test"+strconv.Itoa(i)), []byte("test"))
        if err != nil {
            t.Error(err)
        }
        t.Log(ret)
    }
    end := time.Now().Nanosecond()
    t.Log("=============================")
    t.Log("10000 times insert spend:(ms)", (end-start)/1000/1000)
    db.Close()
    os.RemoveAll("db")
    // t.Log(db.)
}
