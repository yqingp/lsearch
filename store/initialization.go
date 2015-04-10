package store

import (
    "log"
    "os"
    "path/filepath"
)

func (self *DB) initKmap() error {
    var err error
    keyMapTrieFilePath := filepath.Join(self.baseDir, KeyMapTrieFileName)
    if self.keyMapTrie, err = OpenTrie(keyMapTrieFilePath); err != nil {
        return err
    }

    return nil
}

func (self *DB) initLogger() error {
    loggerFilePath := filepath.Join(self.baseDir, DbLogFileName)

    f, err := os.OpenFile(loggerFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)

    if err != nil {
        return err
    }
    self.logger = log.New(f, "[LSearch][DB]:", log.Llongfile|log.Ldate|log.Ltime)

    // self.logger = log.New(os.Stdout, "[LSearch][DB]:", log.Llongfile|log.Ldate|log.Ltime)
    return nil
}

func (self *DB) initDir() error {
    if err := os.MkdirAll(self.baseDir, 0755); err != nil {
        return err
    }

    return nil
}
