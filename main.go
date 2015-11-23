package main

import (
	"github.com/luxuan/go-memcached-server"
)

func main() {
	server, err := memcached.NewServer("", nil)
	if err != nil {
		panic(err)
	}

	// register handler
	server.RegisterFunc("get", memcached.DefaultGet)
	server.RegisterFunc("set", memcached.DefaultSet)
	server.RegisterFunc("delete", memcached.DefaultDelete)
	server.RegisterFunc("incr", memcached.DefaultIncr)
	server.RegisterFunc("flush_all", memcached.DefaultFlushAll)
	server.RegisterFunc("version", memcached.DefaultVersion)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
