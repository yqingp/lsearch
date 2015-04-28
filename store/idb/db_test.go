package idb

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
	_, err = db.Set(1, []byte("test"+strconv.Itoa(1)), []byte("test"))
	if err != nil {
		t.Error(err)
	}
	a, _ := db.GetByInternalId(1)
	t.Log(string(a))
	db.Add(1, []byte("test"))
	a, _ = db.GetByInternalId(1)
	t.Log(string(a))
	end := time.Now()
	t.Log("=============================")
	t.Log("10000 values insert spend:(ms)", end.Sub(start))
	db.Close()
}
