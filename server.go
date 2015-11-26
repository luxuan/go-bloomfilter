/* *
 * Date: 2015.11.24
 * Author: @luxuan
 * */

package main

import (
    "github.com/luxuan/go-memcached-server"
    "github.com/luxuan/go-memcached-server/protocol"
    "github.com/luxuan/go-bloomfilter/bloom"
)

//implement: set/get(s)/incr/version
type Handler struct {
    *bloom.Dispatcher
}

func (h *Handler)BFSet(req *protocol.McRequest, res *protocol.McResponse) error {
    h.Set(req.Key)
    res.Response = "STORED"
    return nil
}

func (h *Handler)BFGet(req *protocol.McRequest, res *protocol.McResponse) error {
    var b []byte
    for _, key := range req.Keys {
        if exist := h.Get(key); exist {
            b = []byte("1")
        } else {
            b = []byte("0")
        }
        res.Values = append(res.Values, protocol.McValue{key, "0", b})
    }
    res.Response = "END"
    return nil
}

func (h *Handler)BFIncr(req *protocol.McRequest, res *protocol.McResponse) error {
    if exist := h.Incr(req.Key); exist {
        res.Response = "1"
    } else {
        res.Response = "0"
    }
    return nil
}

func (h *Handler)BFVersion(req *protocol.McRequest, res *protocol.McResponse) error {
    res.Response = "VERSION simple-memcached-0.1"
    return nil
}


func main() {
    server, err := memcached.NewServer("", nil)
    if err != nil {
        panic(err)
    }

    h := &Handler{
        bloom.NewDispatcher(0.01, 1024*1024*64),
    }
    // register handler
    server.RegisterFunc("set", h.BFSet)
    server.RegisterFunc("get", h.BFGet)
    server.RegisterFunc("gets", h.BFGet)
    server.RegisterFunc("incr", h.BFIncr)
    server.RegisterFunc("version", h.BFVersion)

    if err := server.ListenAndServe(); err != nil {
        panic(err)
    }
}
