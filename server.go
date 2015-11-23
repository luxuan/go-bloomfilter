package memcached

import (
	//"fmt"
	//"os"
	"github.com/luxuan/go-memcached-server/protocol"
	"strconv"
)

var dpc = NewDispatcher(0.01, 1024*1024*64) // 64MB

//implement: set/get/incr/version

func BFSet(req *protocol.McRequest, res *protocol.McResponse) error {
    dpc.Set(req.Key)
	res.Response = "STORED"
	return nil
}

func BFGet(req *protocol.McRequest, res *protocol.McResponse) error {
	for _, key := range req.Keys {
        exist := dpc.Get(key)
		res.Values = append(res.Values, protocol.McValue{key, "0", []byte(exist)})
	}
	res.Response = "END"
	return nil
}

func BFIncr(req *protocol.McRequest, res *protocol.McResponse) error {
    exist := dpc.Incr(req.Key)
	res.Response = string(exist)
	return nil
}

func DefaultVersion(req *protocol.McRequest, res *protocol.McResponse) error {
	res.Response = "VERSION simple-memcached-0.1"
	return nil
}

