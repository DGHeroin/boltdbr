package main

import (
    "flag"
    "fmt"
    "github.com/DGHeroin/boltdbr"
    "github.com/boltdb/bolt"
    "io/ioutil"
    "net/http"

    "net/rpc/jsonrpc"
    "os"
    "time"
)

var (
    address      = flag.String("addr", ":8192", "http serve address")
    testing      = flag.Bool("test", false, "is test mode")
    boltFilename = flag.String("f", "data.db", "boltdb filename")
    serveType    = flag.Int("k", 0, "serve type: [0 json_rpc] [1 http_rpc]")
    token        = flag.String("t", "", "token")
)

func main() {
    flag.Parse()
    var err error
    var boltFileName string
    if *testing {
        // create temp file
        f, _ := ioutil.TempFile("", "bolt-")
        f.Close()
        os.Remove(f.Name())
        boltFileName = f.Name()
    } else {
        boltFileName = *boltFilename
    }

    s := &boltdbr.BoltDBR{}
    s.DB, err = bolt.Open(boltFileName, 0600, &bolt.Options{Timeout: time.Second})
    if err != nil {
        fmt.Println(err)
        return
    }

    rpcServer, l, err := boltdbr.New(s, boltdbr.Options{Address: *address, Token: *token, Type:*serveType})
    if err != nil {
        fmt.Println(err)
        return
    }
    defer l.Close()
    if *serveType == 0 {
        for {
            conn, err := l.Accept()
            if err != nil {
                fmt.Println(err)
                return
            }
            go rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
        }
    }
    if *serveType == 1 {
        http.Serve(l, nil)
    }
}
