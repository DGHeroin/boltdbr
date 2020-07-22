package main

import (
    "flag"
    "fmt"
    "github.com/DGHeroin/boltdbr"
    "github.com/boltdb/bolt"
    "io/ioutil"
    "net/http"
    "os"
    "time"
)

var (
    address = flag.String("addr", ":8192", "http serve address")
)

func main()  {
    var err error
    // create temp file
    f, _ := ioutil.TempFile("", "bolt-")
    f.Close()
    os.Remove(f.Name())
    boltFileName := f.Name()

    s := &boltdbr.BoltDBR{}
    s.DB, err = bolt.Open(boltFileName, 0600, &bolt.Options{Timeout: time.Second})
    l, err := boltdbr.NewHttp(s, *address)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer l.Close()
    http.Serve(l, nil)
}
