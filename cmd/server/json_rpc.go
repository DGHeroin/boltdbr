package main

import (
    "flag"
    "fmt"
    "github.com/DGHeroin/boltdbr"
    "github.com/boltdb/bolt"
    "io/ioutil"

    "net/rpc/jsonrpc"
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
    rpcServer, l, err := boltdbr.NewJson(s, *address)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer l.Close()

    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println("new conn", conn)
        //go func(conn net.Conn) {
        //    jsonrpc.ServeConn(conn)
        //}(conn)
        go rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
}
