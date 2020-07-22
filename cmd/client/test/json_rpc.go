package main

import (
    "fmt"
    "github.com/DGHeroin/boltdbr"
    "net"
    "net/rpc/jsonrpc"
    "sync/atomic"
    "time"
)

var (
    address = ":8192"
    token   = "123"
)

func main() {
    conn, err := net.DialTimeout("tcp", address, time.Second*30)
    if err != nil {
        fmt.Println(err)
        return
    }
    client := jsonrpc.NewClient(conn)
    {
        q := &boltdbr.Query{Bucket: []byte("default")}
        q.Token = token
        r := new(boltdbr.Response)
        err = client.Call("BoltDBR.CreateBucket", q, &r)
        if err != nil {
            fmt.Println(err)
            return
        }
    }

    {
        q := &boltdbr.Query{Bucket: []byte("default"), Key: []byte("my-key"), Value: []byte("my-value")}
        q.Token = token
        r := new(boltdbr.Response)
        err = client.Call("BoltDBR.Set", q, &r)
        if err != nil {
            fmt.Println(err)
            return
        }
    }

    {
        q := &boltdbr.Query{Bucket: []byte("default"), Key: []byte("my-key")}
        q.Token = token
        r := new(boltdbr.Response)
        err = client.Call("BoltDBR.Get", q, &r)
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println("get", string(r.Value), r.Error)
    }

    testQps()
}

func testQps() {
    qps := int64(0)
    isRunning := true
    readCb := func() {
        conn, err := net.DialTimeout("tcp", address, time.Second*30)
        if err != nil {
            fmt.Println(err)
            return
        }
        client := jsonrpc.NewClient(conn)
        for isRunning {
            q := &boltdbr.Query{Bucket: []byte("default"), Key: []byte("my-key")}
            q.Token = token
            r := new(boltdbr.Response)
            err = client.Call("BoltDBR.Get", q, &r)
            if err != nil {
                fmt.Println(err)
                return
            }
            atomic.AddInt64(&qps, 1)
        }
    }

    for i := 0; i < 100; i++ {
        go readCb()
    }

    ticker := time.NewTicker(time.Second)
    for range ticker.C {
        q := atomic.LoadInt64(&qps)
        atomic.StoreInt64(&qps, 0)
        fmt.Println("qps:", q)
    }
}
