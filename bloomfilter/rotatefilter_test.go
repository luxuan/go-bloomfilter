package bloomfilter

import (
    //"fmt"
    "testing"
)

func TestSetAndGet(t *testing.T) {
    rf := NewRotateFilter(7, 1024, 10, 4)
    
    key := uint32(1111111111)
    if rf.Check(key) != -1 {
        t.Errorf("Get no-exist key error")
    }

    if rf.Add(key) {
        t.Errorf("First add should be empty")
    }

    if rf.Check(key) != 0 {
        t.Errorf("Get exist key error")
    }
}

func TestIncr(t *testing.T) {
    rf := NewRotateFilter(7, 1024, 10, 4)
    key := uint32(22222222)
    if rf.Incr(key) != -1 {
        t.Errorf("Incr no-exist key error")
    }

    if rf.Incr(key) != 0 {
        t.Errorf("Incr exist key error")
    }
}

func TestRotateTable(t *testing.T) {
    rf := NewRotateFilter(7, 1024, 10, 4)
    key0 := uint32(333777)
    key := key0
    for i := 0 ;i < 15; i++ {
        rf.Add(key)
        key += 789797
    }

    if idx := rf.Check(key0); idx != 1 {
        t.Errorf("Get rotated key error %d", idx)
    }

    if idx := rf.Check(key - 789797); idx != 0 {
        t.Errorf("Get no-rotate key error %d", idx)
    }
    t.Errorf("DEBUG")
}
