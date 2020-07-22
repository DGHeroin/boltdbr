package boltdbr

import (
    "errors"
    "github.com/boltdb/bolt"
    "net"
    "net/rpc"
)

var (
    ErrParamsEror     = errors.New("params error")
    ErrBucketNotFound = errors.New("bucket not found")
    ErrTokenNotMatch  = errors.New("token not match")
)

const (
    JsonRPC = iota
    HTTP
)

type BoltDBR struct {
    DB  *bolt.DB
    opt Options
}

type Query struct {
    Bucket []byte
    Key    []byte
    Value  []byte
    Token  string
}

type Response struct {
    Value []byte
    Error string
}
type Options struct {
    Address string
    Token   string
    Type    int
}

func New(s *BoltDBR, options Options) (*rpc.Server, net.Listener, error) {
    s.opt = options
    if options.Type == JsonRPC {
        srv := rpc.NewServer()
        if err := srv.Register(s); err != nil {
            return nil, nil, err
        }
        l, err := net.Listen("tcp", options.Address)
        return srv, l, err
    }
    if options.Type == HTTP {
        rpc.RegisterName("BoltDBR", s)
        rpc.HandleHTTP()
        l, err := net.Listen("tcp", options.Address)
        return nil, l, err
    }
    return nil, nil, ErrBucketNotFound
}

func (s *BoltDBR) checkAuth(q *Query) error {
    if s.opt.Token == "" {
        return nil
    }
    if s.opt.Token == q.Token {
        return nil
    }
    return ErrTokenNotMatch
}

func (s *BoltDBR) CreateBucket(q *Query, r *Response) error {
    if err := s.checkAuth(q); err != nil {
        return err
    }
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
    if err := s.checkAuth(q); err != nil {
        return err
    }
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
    if err := s.checkAuth(q); err != nil {
        return err
    }
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
    if err := s.checkAuth(q); err != nil {
        return err
    }
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
    if err := s.checkAuth(q); err != nil {
        return err
    }
    err := s.DB.Update(func(tx *bolt.Tx) error {
        return tx.DeleteBucket(q.Bucket)
    })
    return err
}
