package bdd

import (
	"math"
	"math/big"
	"math/rand"
)

type cacheEntry struct {
	a     int32
	b     int32
	c     int32
	bdres *big.Int
	res   int32
}

func newCacheEntry() *cacheEntry {
	return &cacheEntry{a: -1}
}

func (c *cacheEntry) reset() {
	c.a = -1
}

type cache struct {
	table []*cacheEntry
}

func newCache(size int32) *cache {
	cache := &cache{}
	cache.resize(size)
	return cache
}

func (c *cache) resize(ns int32) {
	size := primeGte(int(ns))
	c.table = make([]*cacheEntry, size)
	for n := 0; n < size; n++ {
		c.table[n] = newCacheEntry()
	}
}

func (c *cache) reset() {
	for _, e := range c.table {
		e.reset()
	}
}

func (c *cache) lookup(hash int32) *cacheEntry {
	return c.table[int32(math.Abs(float64(hash%int32(len(c.table)))))]
}

const checktimes = 20

func primeGte(num int) int {
	if isEven(num) {
		num++
	}
	for !isPrime(num) {
		num += 2
	}
	return num
}

func primeLte(num int) int {
	if isEven(num) {
		num--
	}
	for !isPrime(num) {
		num -= 2
	}
	return num
}

func isEven(src int) bool {
	return (src & 0x1) == 0
}

func isPrime(src int) bool {
	return !hasEasyFactors(src) && isMillerRabinPrime(src)
}

func hasEasyFactors(src int) bool {
	return hasFactor(src, 3) || hasFactor(src, 5) || hasFactor(src, 7) || hasFactor(src, 11) || hasFactor(src, 13)
}

func hasFactor(src, n int) bool {
	return (src != n) && (src%n == 0)
}

func isMillerRabinPrime(src int) bool {
	for n := 0; n < checktimes; n++ {
		witness := random(src - 1)
		if isWitness(witness, src) {
			return false
		}
	}
	return true
}

func isWitness(witness, src int) bool {
	bitNum := numberOfBits(src-1) - 1
	d := 1
	for i := bitNum; i >= 0; i-- {
		x := d
		d = mulmod(d, d, src)
		if d == 1 && x != 1 && x != src-1 {
			return true
		}
		if bitIsSet(src-1, i) {
			d = mulmod(d, witness, src)
		}
	}
	return d != 1
}

func numberOfBits(src int) int {
	if src == 0 {
		return 0
	}
	for b := 31; b > 0; b-- {
		if bitIsSet(src, b) {
			return b + 1
		}
	}
	return 1
}

func bitIsSet(src, b int) bool {
	return (src & (1 << b)) != 0
}

func mulmod(a, b, c int) int {
	return (a * b) % c
}

func random(i int) int {
	return rand.Intn(i) + 1
}
