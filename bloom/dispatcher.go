/* *
 * Date: 2015.11.23
 * Author: @luxuan
 * Refer: https://github.com/tylertreat/BoomFilters/
 * */

package bloom

import (
    //"fmt"
    "sync"
    "math"
    "hash"
    "hash/fnv"
    "encoding/binary"

    //"time"
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
    return uint(float64(size) * 0.69 / float64(k) * 8);
}

func OptimalK(fpRate float64) uint32 {
    return uint32(math.Ceil(math.Log2(1 / fpRate)))
}

var hashfn hash.Hash64 = fnv.New64()
func Hash(b []byte) (uint32, uint32) {
    hashfn.Reset()
    hashfn.Write(b)
    //return hashfn.Sum64()
    sum := hashfn.Sum(nil)
    return binary.BigEndian.Uint32(sum[4:8]), binary.BigEndian.Uint32(sum[0:4])
}


// locker, dispatcher, adapter
type Dispatcher struct {
    nslot uint32
    rfs   []*RotateFilter
    locks []sync.Mutex
}

func NewDispatcher(erate float64, size uint) *Dispatcher {
    if size < 1024 * 1024 { // at least 1MB
        return nil
    }
    nslot := uint32(1024)
    ntable := uint32(4)

    dpc := &Dispatcher {
        nslot: nslot,
        locks: make([]sync.Mutex, nslot),
        rfs:   make([]*RotateFilter, nslot),
    }
    k := OptimalK(erate)
    nbit := uint32(size / uint(nslot) * 8)
    maxCount := uint32(OptimalN(size / uint(nslot), k))

    for i := uint32(0); i < nslot; i++ {
        dpc.rfs[i] = NewRotateFilter(k, nbit, maxCount, ntable)
    }
    return dpc
}

func (dpc *Dispatcher) Set(key string) bool {
    idx, code := Hash([]byte(key))
    idx %=  dpc.nslot
    dpc.locks[idx].Lock()
    existInCurrentTable := dpc.rfs[idx].Add(code)
    dpc.locks[idx].Unlock()
    return existInCurrentTable
}

func (dpc *Dispatcher) Get(key string) bool {
    idx, code := Hash([]byte(key))
    idx %=  dpc.nslot
    dpc.locks[idx].Lock()
    itable := dpc.rfs[idx].Check(code)
    dpc.locks[idx].Unlock()
    return itable >= 0
}

// return exist before add
func (dpc *Dispatcher) Incr(key string) bool {
    idx, code := Hash([]byte(key))
    idx %=  dpc.nslot
    dpc.locks[idx].Lock()
    itableBeforeAdd := dpc.rfs[idx].Incr(code)
    dpc.locks[idx].Unlock()
    return itableBeforeAdd >= 0
}
