/* *
 * Date: 2015.11.18
 * Author: @luxuan
 * Refer: https://github.com/reddragon/bloomfilter.go
 * Refer: https://github.com/tylertreat/BoomFilters/
 * */

package bloomfilter

import (
    "fmt"
    //"math"
    "hash"
    "hash/fnv"

    "time"
)

// TODO lock rf
type RotateFilter struct {
    k     uint32      // Number of hash functions
    nbit  uint32      // Number of bits in a filter
    maxCount uint32

    // TODO structed?
    ntable  uint32
    bytes   [][]byte    // The rotate filters bitmap
    counts  []uint32    // Number of elements in the rotate filters
    times []time.Time
}

// TODO refer to outer 
func NewRotateFilter(k, nbit, maxCount, ntable uint32) *RotateFilter {
    rf := &RotateFilter {
        k: k,
        nbit: nbit,
        maxCount: maxCount,

        ntable: ntable,
        bytes: make([][]byte, ntable),
        counts: make([]uint32, ntable),
        times: make([]time.Time, ntable),
    }

    for i := uint32(0); i < ntable; i++ {
        rf.bytes[i] = make([]byte, nbit / 8)
    }
    rf.times[0] = time.Now()
    return rf
}

func (rf *RotateFilter) rotateTable() {
    // rotate table, such as [0, 1, 2, 3] -> [3, 0, 1, 2]

    headCounts := rf.counts[rf.ntable - 1:]
    tailCounts := rf.counts[:rf.ntable - 1]
    rf.counts = append(headCounts, tailCounts...)

    headTimes := rf.times[rf.ntable - 1:]
    tailTimes := rf.times[:rf.ntable - 1]
    rf.times = append(headTimes, tailTimes...)

    headBytes := rf.bytes[rf.ntable - 1:]
    tailBytes := rf.bytes[:rf.ntable - 1]
    rf.bytes = append(headBytes, tailBytes...)

    // reset 0 for new table[0]
    rf.counts[0] = 0
    rf.times[0] = time.Now()
    nbyte := uint32(rf.nbit>>3)
    for i := uint32(0); i < nbyte; i++ {
        rf.bytes[0][i] = 0
    }
    //fmt.Println("rotate")
}

// only check exist in current table before add
func (rf *RotateFilter) Add(hashCode uint32) bool {
    // rotate if counter of last table has beyond the limit
    if rf.counts[0] >= rf.maxCount {
        rf.rotateTable()
    }

    exist := true  // exist key before insert
    delta := (hashCode >> 17) | (hashCode << 15);
    for i := uint32(0); i < rf.k; i++ {
        idx := hashCode % rf.nbit
        // if not exist then set
        if (rf.bytes[0][idx>>3] >> (idx & 7) & 1) == 0 {
            rf.bytes[0][idx>>3] |= (1 << (idx & 7))
            //fmt.Println("false")
            exist = false
        }
        hashCode += delta;
    }   
    
    if !exist {
        rf.counts[0]++
    }

    //fmt.Println(rf.counts)
    return exist
}

func (rf *RotateFilter) Check(hashCode uint32) int32 {
    fmt.Println(hashCode, rf.counts)
    delta := (hashCode >> 17) | (hashCode << 15);
    // check in all filters
    for i := uint32(0); i < rf.ntable; i++ {
        var j uint32
        for j = uint32(0); j < rf.k; j++ {
            idx := hashCode % rf.nbit
            if (rf.bytes[i][idx>>3] >> (idx & 7) & 1) == 0 {
                break; // not in this filter
            }
            hashCode += delta;
        }   
        if j == rf.k {
            return int32(i)
        }
    }
    return -1
}

// return idx of filters before add
func (rf *RotateFilter) Incr(hashCode uint32) int32 {
    if idx := rf.Check(hashCode); idx != 0 {
        // not exist or not in latest table
        rf.Add(hashCode)
        return idx
    }
    return 0
}

// Returns the current False Positive Rate of the Bloom Filter
//func (bf *BloomFilter) FalsePositiveRate() float64 {
//    return math.Pow((1 - math.Exp(-float64(bf.k*bf.n)/
//        float64(bf.m))), float64(bf.k))
//}
