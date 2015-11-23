/* *
 * Date: 2015.11.23
 * Author: @luxuan
 * Refer: https://github.com/tylertreat/BoomFilters/
 * */

package bloomfilter

import (
    "fmt"
    "sync"
    //"math"
    "hash"
    "hash/fnv"

    "time"
)

const fillRatio = 0.5

// size m, based on the number of items and the desired rate of false positives.
func OptimalM(n uint, fpRate float64) uint {
    return uint(math.Ceil(float64(n) / ((math.Log(fillRatio) *
        math.Log(1-fillRatio)) / math.Abs(math.Log(fpRate)))))
}

// nitem, based on mem size and optimalK
// TODO match the fillRate
func OptimalN(size uint, k uint32) uint {
    /* 0.69 =~ ln(2) ~= -ln(1 - erate^(1/log2(1/erate))) */
    return uint(size * 0.69 / k * 8);
}

func OptimalK(fpRate float) uint32 {
    return uint32(math.Ceil(math.Log2(1 / fpRate)))
}

var hashfn hash.Hash64 = fnv.New64()
func Hash(b []byte) (uint64) {
    hashfn.Reset()
    hashfn.Write(b)
    return hashfn.Sum64()
    // sum := hash.Sum(nil)
    // return binary.BigEndian.Uint32(sum[4:8]), binary.BigEndian.Uint32(sum[0:4])
}


// locker, dispatcher, adapter
type Dispatcher struct {
    nslot int
    rfs   []*RotateFilter
    locks []sync.Mutex
}

func NewDispatcher(erate float, size uint) *Dispatcher {
    if size < 1024 * 1024 { // at least 1MB
        return nil
    }
    nslot = 1024
    ntable = 4

    dispatcher := &Dispatcher {
        nslot: nslot
        locks: make(sync.Mutex, nslot),
        rfs:   make(*RotateFilter, nslot),
    }
    k := OptimalK(erate)
    nbit := size / nslot * 8
    maxCount := OptimalN(size / nslot, k)

    for i := 0; i < nslot; i++ {
        rfs[i] = NewRotateFilter(k, uint32(nbit), uint32(maxCount), uint32(ntable))
    }
    return dispatcher
}

func (dpc *Dispatcher) Set(key string) bool {
    i64 := Hash([]byte(key))
    // TODO cut bits
    i := i64 % dpc.nslot
    dpc.locks[i].Lock()
    existInCurrentTable := dpc.rfs[i].Add(uint32(i64))
    dpc.locks[i].Unlock()
    return existInCurrentTable >= 0
}

func (dpc *Dispatcher) Get(key string) bool {
    i64 := Hash([]byte(key))
    // TODO cut bits
    i := i64 % dpc.nslot
    dpc.locks[i].Lock()
    idx := dpc.rfs[i].Incr(uint32(i64))
    dpc.locks[i].Unlock()
    return idxBeforeAdd >= 0
}

// return exist before add
func (dpc *Dispatcher) Incr(key string) bool {
    i64 := Hash([]byte(key))
    // TODO cut bits
    i := i64 % dpc.nslot
    dpc.locks[i].Lock()
    idxBeforeAdd := dpc.rfs[i].Incr(uint32(i64))
    dpc.locks[i].Unlock()
    return idxBeforeAdd >= 0
}
