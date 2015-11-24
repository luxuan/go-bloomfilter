package main

import (
    //"fmt"
    "testing"
    "github.com/bradfitz/gomemcache/memcache"
)

func TestSetAndGet(t *testing.T) {
    key0 := "set_and_get_key0"
    key1 := "set_and_get_key1"

    mc := memcache.New("127.0.0.1:11211")
    mc.Set(&memcache.Item{Key: key0, Value: []byte("")})

    item, err := mc.Get(key0)
    if err != nil || string(item.Value) != "1" {
        t.Errorf("Get exist key error", err)
    }

    item, err = mc.Get(key1)
    if err != nil || string(item.Value) != "0" {
        t.Errorf("Get no-exist key error", err)
    }
}

func TestIncr(t *testing.T) {
    key0 := "incr_key0"

    mc := memcache.New("127.0.0.1:11211")
    old, err := mc.Increment(key0, 0)
    if err != nil || old != 0 {
        t.Errorf("Incr no-exist key error", err)
    }

    old, err = mc.Increment(key0, 0)
    if err != nil || old != 1 {
        t.Errorf("Incr exist key error", err)
    }
}
