package main

import (
    "fmt"
    "github.com/DGHeroin/boltdbr"
    "net/rpc"
)

func main()  {
    client, err := rpc.DialHTTP("tcp", ":8192")
    if err != nil {
        fmt.Println(err)
        return
    }
    {
        q := &boltdbr.Query{Bucket: []byte("default")}
        r := new(boltdbr.Response)
        err = client.Call("BoltDBR.CreateBucket", q, &r)
        if err != nil {
            fmt.Println(err)
            return
        }
    }

    {
        q := &boltdbr.Query{Bucket: []byte("default"), Key: []byte("my-key"), Value: []byte("my-value")}
        r := new(boltdbr.Response)
        err = client.Call("BoltDBR.Set", q, &r)
        if err != nil {
            fmt.Println(err)
            return
        }
    }

    {
        q := &boltdbr.Query{Bucket: []byte("default"), Key: []byte("my-key")}
        r := new(boltdbr.Response)
        err = client.Call("BoltDBR.Get", q, &r)
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println("get", string(r.Value), r.Error)
    }

}
