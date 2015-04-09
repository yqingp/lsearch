package store

import (
    "github.com/yqingp/lsearch/util"
    "log"
    "os"
    "path/filepath"
)

func (self *Db) initKmap() error {
    var err error
    kmapfileName := filepath.Join(self.basedir, "db.kmap")
    if self.kmap, err = util.Open(kmapfileName); err != nil {
        return err
    }

    return nil
}

func (self *Db) initLogger() error {
    loggerFileName := filepath.Join(self.basedir, "db.log")

    _, err := os.OpenFile(loggerFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)

    if err != nil {
        return err
    }

    self.logger = log.New(os.Stdout, "[LSearch][DB]:", log.Llongfile|log.Ldate|log.Ltime)
    return nil
}

func (self *Db) initDir() error {
    if err := os.MkdirAll(self.basedir, 0755); err != nil {
        return err
    }

    return nil
}
