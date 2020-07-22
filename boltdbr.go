package boltdbr

import (
    "errors"
    "github.com/boltdb/bolt"
    "net"
    "net/rpc"
)

var (
    ErrBucketNotFound = errors.New("bucket not found")
)

type BoltDBR struct {
    DB *bolt.DB
}

type Query struct {
    Bucket []byte
    Key    []byte
    Value  []byte
}
type Response struct {
    Value []byte
    Error string
}

func NewHttp(s *BoltDBR, address string) (net.Listener, error) {
    rpc.RegisterName("BoltDBR", s)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp", address)
    return l, err
}

func NewJson(s *BoltDBR, address string) (*rpc.Server, net.Listener, error) {
    srv := rpc.NewServer()

    if err := srv.Register(s); err != nil {
        return nil, nil, err
    }
    l, err := net.Listen("tcp", address)
    return srv, l, err
}

func (s *BoltDBR) CreateBucket(q *Query, r *Response) error {
    err := s.DB.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists(q.Bucket)
        return err
    })
    if err != nil {
        r.Error = err.Error()
    }

    return err
}

func (s *BoltDBR) Get(q *Query, r *Response) error {
    err := s.DB.View(func(tx *bolt.Tx) error {
        bk := tx.Bucket(q.Bucket)
        if bk == nil {
            return ErrBucketNotFound
        }
        r.Value = bk.Get(q.Key)
        return nil
    })

    return err
}
func (s *BoltDBR) Set(q *Query, r *Response) error {
    err := s.DB.Update(func(tx *bolt.Tx) error {
        bk, err := tx.CreateBucketIfNotExists(q.Bucket)
        if err != nil {
            return err
        }
        return bk.Put(q.Key, q.Value)
    })
    return err
}

func (s *BoltDBR) Delete(q *Query, r *Response) error {
    err := s.DB.Update(func(tx *bolt.Tx) error {
        bk, err := tx.CreateBucketIfNotExists(q.Bucket)
        if err != nil {
            return err
        }
        return bk.Delete(q.Key)
    })
    return err
}

func (s *BoltDBR) DeleteBucket(q *Query, r *Response) error {
    err := s.DB.Update(func(tx *bolt.Tx) error {
        return tx.DeleteBucket(q.Bucket)
    })
    return err
}

