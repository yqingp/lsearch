# lsearch-store

##安装
  go get github.com/yqingp/lsearch


##使用
  ```go
  package main

  import (
      "fmt"
      . "github.com/yqingp/lsearch/store"
      "os"
      "strconv"
      "time"
  )

  func main() {
      db, err := Open("./db", true)
      if err != nil {
          fmt.Println(err)
          os.Exit(-1)
      }

      if db == nil {
          fmt.Println("db init fail")
          os.Exit(-1)
      }

      start := time.Now()
      for i := 0; i < 1000000; i++ {
          _, err := db.Set(-1, []byte("test"+strconv.Itoa(i)), []byte("test"+strconv.Itoa(i)))
          if err != nil {
              fmt.Println(err)
              os.Exit(-1)
          }
      }
      end := time.Now()
      fmt.Println("=============================")
      fmt.Println("1000000 values insert spend:(s)", end.Sub(start))

      ret, _ := db.Get([]byte("test0"))
      fmt.Println(string(ret))
      ret, _ = db.GetInternalId(1)
      fmt.Println(string(ret))

      db.Close()
  }
  ```