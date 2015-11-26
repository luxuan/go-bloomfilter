package main

import (
    "fmt"
    "sync"
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


// go test -test.bench="."

const nroutine = 1000
const moretimes = 100
func BenchmarkAdd(b *testing.B) {  
    b.StopTimer()  

    keyPrefix := "set_key"
    var mcs [nroutine]*memcache.Client
    for i := 0; i < nroutine; i++ {
        mcs[i] = memcache.New("127.0.0.1:11211")

        // pre get because the lazy connecting strategy
        mcs[i].Get(keyPrefix) 
    }

    b.StartTimer()  


    var wg sync.WaitGroup
    nloop := b.N / nroutine * moretimes
    for i := 0; i < nroutine; i++ {  
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()

            for j := 0; j < nloop; j++ {  
                key := fmt.Sprintf("%s_%d_%d", keyPrefix, i, j)
                mcs[idx].Get(key)
            }
        }(i)
    }  
    wg.Wait()
}  
