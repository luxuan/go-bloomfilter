/* *
 * Date: 2015.11.18
 * Author: @luxuan
 * Refer: https://github.com/reddragon/bloomfilter.go
 * Refer: https://github.com/tylertreat/BoomFilters/
 * */

package bloomfilter

import (
    "math"
    "hash"
    "hash/fnv"

    "time"
)

var hashfn hash.Hash64 = fnv.New64()

func Hash(b []byte) (uint64) {
    hashfn.Reset()
    hashfn.Write(b)
    return hashfn.Sum64()
}


type RotateFilter struct {
    k      int         // Number of hash functions
    nbits  uint32      // Number of bits in a filter

    // TODO vector or structed?
    ntable  int
    bytes   [][]byte    // The rotate filters bitmap
    counts  []uint32    // Number of elements in the rotate filters
    rotateTimes []time.Time
}

func NewBloomFilter(numHashFuncs, bfSize int) *BloomFilter {
    bf := new(BloomFilter)
    bf.bitmap = make([]bool, bfSize)
    bf.k, bf.m = numHashFuncs, bfSize
    bf.n = 0
    bf.hashfn = fnv.New64()
    return bf
}

// rotate if counter of last table has beyond the limit
func (rf *RotateFilter) rotateIfNeeded() bool {
}

func (rf *RotateFilter) CheckAndAdd(hashCode uint32) bool {
    // TODO exist which table?
    if exist := Check(hashCode), existLastTable {
        return exist
    }
    Add(hashCode)
    return false
}

func (rf *RotateFilter) Add(hashCode uint32) bool {
    rf.rotateIfNeeded()

    // exists in which table?
    exist := true  // exist key before insert
    delta := (hashCode >> 17) | (hashCode << 15);
    for i := 0; i < rf.k; i++ {
        idx = hashCode % rf->nbit;
        // TODO n tables?
        if !bitGetCurrentTable(rf.bytes, idx) {
            bitSetCurrentTable(rf.bytes, idx)
            exist = false
        }
        hashCode += delta;
    }   
    
    if !exist {
        rf.count++
    }

    return exist
}

// -1 not exist
// 0 latest table
// ...
// n oldest table
func (rf *RotateFilter) Check(hashCode uint32) bool {
    delta := (hashCode >> 17) | (hashCode << 15);
    for i := 0; i < rf.k; i++ {
        idx = hashCode % rf->nbit;
        // TODO n tables
        if !bit_get(rf.bytes, idx) {
            return false
        }
        hashCode += delta;
    }   
    return true
}

// Returns the current False Positive Rate of the Bloom Filter
func (bf *BloomFilter) FalsePositiveRate() float64 {
    return math.Pow((1 - math.Exp(-float64(bf.k*bf.n)/
        float64(bf.m))), float64(bf.k))
}
